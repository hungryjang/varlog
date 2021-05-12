package executor

import (
	"sync"
	"time"

	"github.com/kakao/varlog/pkg/types"
)

var commitTaskPool = sync.Pool{
	New: func() interface{} {
		return &commitTask{}
	},
}

type commitTask struct {
	highWatermark      types.GLSN
	prevHighWatermark  types.GLSN
	committedGLSNBegin types.GLSN
	committedGLSNEnd   types.GLSN

	ctime time.Time
}

func newCommitTask() *commitTask {
	return commitTaskPool.Get().(*commitTask)
}

func (t *commitTask) release() {
	t.highWatermark = types.InvalidGLSN
	t.prevHighWatermark = types.InvalidGLSN
	t.committedGLSNBegin = types.InvalidGLSN
	t.committedGLSNEnd = types.InvalidGLSN
	commitTaskPool.Put(t)
}