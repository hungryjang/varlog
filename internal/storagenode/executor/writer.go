package executor

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"

	"github.com/kakao/varlog/internal/storagenode/storage"
	"github.com/kakao/varlog/pkg/types"
	"github.com/kakao/varlog/pkg/util/mathutil"
	"github.com/kakao/varlog/pkg/util/runner"
	"github.com/kakao/varlog/pkg/verrors"
)

type writerConfig struct {
	queueSize  int
	batchSize  int
	strg       storage.Storage
	lsc        *logStreamContext
	committer  committer
	replicator replicator
	state      stateProvider
}

func (c writerConfig) validate() error {
	if c.queueSize <= 0 {
		return errors.Wrap(verrors.ErrInvalid, "writer: zero or negative queue size")
	}
	if c.batchSize <= 0 {
		return errors.Wrap(verrors.ErrInvalid, "writer: zero or negative batch size")
	}
	if c.strg == nil {
		return errors.Wrap(verrors.ErrInvalid, "writer: no storage")
	}
	if c.lsc == nil {
		return errors.Wrap(verrors.ErrInvalid, "writer: no log stream context")
	}
	if c.committer == nil {
		return errors.Wrap(verrors.ErrInvalid, "writer: no committer")
	}
	if c.replicator == nil {
		return errors.Wrap(verrors.ErrInvalid, "writer: no replicator")
	}
	if c.state == nil {
		return errors.Wrap(verrors.ErrInvalid, "writer: no state provider")
	}
	return nil
}

type WriterOption interface {
	applyWriter(*writerConfig)
}

/*
type writerState int32

const (
	writerStateInit writerState = iota
	writerStateRun
	writerStateClosed
)
*/
const (
	writerStateInit   = 0
	writerStateRun    = 1
	writerStateClosed = 2
)

type writer interface {
	send(ctx context.Context, tb *appendTask) error
	stop()
	waitForDrainage(ctx context.Context) error
}

type writerImpl struct {
	writerConfig

	q writeQueue

	dispatcher struct {
		runner *runner.Runner
		cancel context.CancelFunc
	}

	inflight int64

	// NOTE: atomic variable can be used to stop writer.
	// Can it avoid leak of appendTask absolutely?
	// ws int32

	// NOTE: Mutex is more expensive than atomic variable.
	running struct {
		val bool
		mu  sync.RWMutex
	}

	popCv struct {
		cv *sync.Cond
		mu sync.Mutex
	}

	writeTaskBatch     []*appendTask
	replicateTaskBatch []*replicateTask
}

var _ writer = (*writerImpl)(nil)

func newWriter(cfg writerConfig) (*writerImpl, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	w := &writerImpl{
		writerConfig:       cfg,
		writeTaskBatch:     make([]*appendTask, 0, cfg.batchSize),
		replicateTaskBatch: make([]*replicateTask, 0, cfg.batchSize),
	}
	if err := w.init(); err != nil {
		return nil, err
	}
	// w.ws = writerStateRun
	return w, nil
}

func (w *writerImpl) init() error {
	q, err := newWriteQueue(w.queueSize)
	if err != nil {
		return err
	}
	w.q = q

	w.popCv.cv = sync.NewCond(&w.popCv.mu)

	r := runner.New("writer", nil)
	cancel, err := r.Run(w.writeLoop)
	if err != nil {
		return err
	}
	w.dispatcher.runner = r
	w.dispatcher.cancel = cancel
	w.running.val = true

	return nil
}

func (w *writerImpl) send(ctx context.Context, tb *appendTask) error {
	if tb == nil {
		return errors.WithStack(verrors.ErrInvalid)
	}

	/*
		if atomic.LoadInt32(&w.ws) != writerStateRun {
			return errors.WithStack(verrors.ErrClosed)
		}
	*/

	w.running.mu.RLock()
	defer w.running.mu.RUnlock()
	if !w.running.val {
		return errors.WithStack(verrors.ErrClosed)
	}

	if err := w.state.mutableWithBarrier(); err != nil {
		return err
	}
	defer w.state.releaseBarrier()

	atomic.AddInt64(&w.inflight, 1)

	tb.ctime = time.Now()
	return w.q.pushWithContext(ctx, tb)
}

func (w *writerImpl) writeLoop(ctx context.Context) {
	for ctx.Err() == nil {
		w.resetBatch()

		if err := w.writeLoopInternal(ctx); err != nil {
			w.state.setSealing()
		}

		w.popCv.cv.L.Lock()
		w.popCv.cv.Signal()
		w.popCv.cv.L.Unlock()
	}
}

func (w *writerImpl) writeLoopInternal(ctx context.Context) error {
	oldLLSN, newLLSN, numPopped, err := w.ready(ctx)
	defer func() {
		atomic.AddInt64(&w.inflight, -numPopped)
	}()
	if err != nil {
		w.batchError(err)
		return err
	}

	if err := w.fanout(ctx); err != nil {
		// NOTE: ownership of write/replicate tasks is moved.
		return err
	}

	if !w.lsc.uncommittedLLSNEnd.CompareAndSwap(oldLLSN, newLLSN) {
		// NOTE: If this CAS operation fails, it means other goroutine changes
		// uncommittedLLSNEnd of LogStreamContext.
		panic(errors.Errorf(
			"uncommittedLLSNEnd swap failure: current=%d old=%d new=%d: more than one writer?",
			w.lsc.uncommittedLLSNEnd.Load(), oldLLSN, newLLSN,
		))
	}
	return nil
}

func (w *writerImpl) ready(ctx context.Context) (types.LLSN, types.LLSN, int64, error) {
	numPopped := int64(0)

	t, err := w.q.popWithContext(ctx)
	if err != nil {
		return types.InvalidLLSN, types.InvalidLLSN, numPopped, err
	}
	numPopped++

	oldLLSN := w.lsc.uncommittedLLSNEnd.Load()
	newLLSN := oldLLSN

	primary := t.primary
	backup := !t.primary

	if err := w.fillBatch(t, newLLSN); err != nil {
		return types.InvalidLLSN, types.InvalidLLSN, numPopped, err
	}
	newLLSN++

	popSize := mathutil.MinInt(w.batchSize-1, w.q.size())
	for i := 0; i < popSize; i++ {
		t := w.q.pop()
		numPopped++

		if err := w.fillBatch(t, newLLSN); err != nil {
			return types.InvalidLLSN, types.InvalidLLSN, numPopped, err
		}

		newLLSN++

		primary = primary || t.primary
		backup = backup || !t.primary
	}
	// NB: Need to debug
	if primary && backup {
		panic("mixed primary and non-primary requests")
	}
	return oldLLSN, newLLSN, numPopped, nil
}

func (w *writerImpl) fillBatch(t *appendTask, llsn types.LLSN) error {
	w.writeTaskBatch = append(w.writeTaskBatch, t)

	//priamry
	if t.primary {
		// assign llsn
		t.llsn = llsn
		// put into replicateTaskBatch
		rt := newReplicateTask()
		rt.llsn = t.llsn
		rt.data = t.data
		rt.replicas = t.replicas
		w.replicateTaskBatch = append(w.replicateTaskBatch, rt)
		return nil
	}

	// backup: check llsn
	if t.llsn != llsn {
		return errors.Errorf("llsn mismatch: %d != %d", llsn, t.llsn)
	}
	return nil
}

func (w *writerImpl) writePipeline(ctx context.Context) error {
	if err := w.write(); err != nil {
		return err
	}
	return w.sendToCommitter(ctx)
}

func (w *writerImpl) write() (err error) {
	batch := w.strg.NewWriteBatch()

	defer func() {
		err = multierr.Append(err, batch.Close())
		if err != nil {
			w.writeTaskBatchError(w.writeTaskBatch, err)
		}
	}()

	for idx := range w.writeTaskBatch {
		t := w.writeTaskBatch[idx]
		err = batch.Put(t.llsn, t.data)
		if err != nil {
			return err
		}
	}

	err = batch.Apply()
	return
}

func (w *writerImpl) sendToCommitter(ctx context.Context) (err error) {
	idx := 0
	for idx < len(w.writeTaskBatch) {
		tb := w.writeTaskBatch[idx]
		err = w.committer.sendCommitWaitTask(ctx, tb)
		if err != nil {
			break
		}
		idx++
	}
	w.writeTaskBatchError(w.writeTaskBatch[idx:], err)
	return err
}

func (w *writerImpl) sendToReplicator(ctx context.Context) (err error) {
	idx := 0
	for idx < len(w.replicateTaskBatch) {
		rtb := w.replicateTaskBatch[idx]
		err = w.replicator.send(ctx, rtb)
		if err != nil {
			break
		}
		idx++
	}
	w.releaseReplicateTaskBatch(w.replicateTaskBatch[idx:])
	return err
}

func (w *writerImpl) fanout(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return w.writePipeline(ctx)
	})
	grp.Go(func() error {
		return w.sendToReplicator(ctx)
	})
	return grp.Wait()
}

func (w *writerImpl) stop() {
	/*
		if !atomic.CompareAndSwapInt32(&w.ws, writerStateRun, writerStateClosed) {
			return
		}
	*/
	w.running.mu.Lock()
	w.running.val = false
	w.running.mu.Unlock()

	w.dispatcher.cancel()
	w.dispatcher.runner.Stop()

	w.drainQueue(errors.WithStack(verrors.ErrClosed))
	w.resetBatch()
}

func (w *writerImpl) drainQueue(err error) {
	for w.q.size() > 0 {
		tb := w.q.pop()
		tb.err = err
		tb.wg.Done()
	}
}

func (w *writerImpl) waitForDrainage(ctx context.Context) error {
	done := make(chan struct{})
	var wg sync.WaitGroup

	defer func() {
		close(done)
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			w.popCv.cv.L.Lock()
			w.popCv.cv.Signal()
			w.popCv.cv.L.Unlock()
		case <-done:
		}
	}()

	w.popCv.cv.L.Lock()
	defer w.popCv.cv.L.Unlock()
	for atomic.LoadInt64(&w.inflight) > 0 && ctx.Err() == nil {
		w.popCv.cv.Wait()
	}
	if w.q.size() == 0 {
		return nil
	}
	return ctx.Err()
}

func (w *writerImpl) resetBatch() {
	w.writeTaskBatch = w.writeTaskBatch[0:0]
	w.replicateTaskBatch = w.replicateTaskBatch[0:0]
}

func (w *writerImpl) batchError(err error) {
	w.writeTaskBatchError(w.writeTaskBatch, err)
	w.releaseReplicateTaskBatch(w.replicateTaskBatch)
}

func (w *writerImpl) writeTaskBatchError(batch []*appendTask, err error) {
	for i := 0; i < len(batch); i++ {
		tb := batch[i]
		tb.err = err
		tb.wg.Done()
	}
}

func (w *writerImpl) releaseReplicateTaskBatch(batch []*replicateTask) {
	for i := 0; i < len(batch); i++ {
		rtb := batch[i]
		rtb.release()
	}
}