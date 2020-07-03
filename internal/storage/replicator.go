package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/kakao/varlog/pkg/varlog/types"
)

type Replica struct {
	StorageNodeID types.StorageNodeID
	LogStreamID   types.LogStreamID
	Address       string
}

type replicateTask struct {
	llsn     types.LLSN
	data     []byte
	replicas []Replica

	errC chan<- error
}

const (
	replicateCSize = 0
)

type Replicator struct {
	rcm        map[types.StorageNodeID]ReplicatorClient
	mtxRcm     sync.RWMutex
	replicateC chan *replicateTask
	once       sync.Once
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

func NewReplicator() *Replicator {
	return &Replicator{
		rcm:        make(map[types.StorageNodeID]ReplicatorClient),
		replicateC: make(chan *replicateTask, replicateCSize),
	}
}

func (r *Replicator) Run(ctx context.Context) {
	r.once.Do(func() {
		ctx, cancel := context.WithCancel(ctx)
		r.cancel = cancel
		r.wg.Add(1)
		go r.dispatchReplicateC(ctx)
	})
}

func (r *Replicator) Close() {
	if r.cancel != nil {
		r.cancel()
	}
	r.wg.Wait()
}

func (r *Replicator) Replicate(ctx context.Context, llsn types.LLSN, data []byte, replicas []Replica) <-chan error {
	errC := make(chan error, 1)
	if len(replicas) == 0 {
		errC <- fmt.Errorf("no replicas")
		close(errC)
		return errC
	}
	task := &replicateTask{
		llsn:     llsn,
		data:     data,
		replicas: replicas,
		errC:     errC,
	}
	select {
	case r.replicateC <- task:
	case <-ctx.Done():
		errC <- ctx.Err()
		close(errC)
	}
	return errC
}

func (r *Replicator) dispatchReplicateC(ctx context.Context) {
	defer r.wg.Done()
	for {
		select {
		case t := <-r.replicateC:
			r.replicate(ctx, t)
		case <-ctx.Done():
			return
		}
	}
}

func (r *Replicator) replicate(ctx context.Context, t *replicateTask) {
	errCs := make([]<-chan error, len(t.replicas))
	for i, replica := range t.replicas {
		rc, err := r.getOrConnect(ctx, replica)
		if err != nil {
			t.errC <- err
			close(t.errC)
			return
		}
		errCs[i] = rc.Replicate(ctx, t.llsn, t.data)
	}
	for _, errC := range errCs {
		select {
		case err := <-errC:
			if err != nil {
				t.errC <- err
				close(t.errC)
				return
			}
		case <-ctx.Done():
			t.errC <- ctx.Err()
			close(t.errC)
			return
		}
	}
	t.errC <- nil
	close(t.errC)
}

func (r *Replicator) getOrConnect(ctx context.Context, replica Replica) (ReplicatorClient, error) {
	r.mtxRcm.RLock()
	rc, ok := r.rcm[replica.StorageNodeID]
	r.mtxRcm.RUnlock()
	if ok {
		return rc, nil
	}

	r.mtxRcm.Lock()
	defer r.mtxRcm.Unlock()
	rc, ok = r.rcm[replica.StorageNodeID]
	if ok {
		return rc, nil
	}
	rc, err := NewReplicatorClient(replica.Address)
	if err != nil {
		return nil, err
	}
	if err = rc.Run(ctx); err != nil {
		return nil, err
	}
	r.rcm[replica.StorageNodeID] = rc
	return rc, nil
}