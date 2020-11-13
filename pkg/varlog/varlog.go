package varlog

import (
	"context"
	"io"

	"go.uber.org/zap"

	"github.com/kakao/varlog/pkg/logc"
	"github.com/kakao/varlog/pkg/mrc/mrconnector"
	"github.com/kakao/varlog/pkg/types"
	"github.com/kakao/varlog/pkg/util/syncutil/atomicutil"
	"github.com/kakao/varlog/proto/varlogpb"
)

// Varlog is a log interface with thread-safety. Many goroutines can share the same varlog object.
type Varlog interface {
	io.Closer

	Append(ctx context.Context, data []byte, opts ...AppendOption) (types.GLSN, error)

	AppendTo(ctx context.Context, logStreamID types.LogStreamID, data []byte, opts ...AppendOption) (types.GLSN, error)

	Read(ctx context.Context, logStreamID types.LogStreamID, glsn types.GLSN) ([]byte, error)

	Subscribe(ctx context.Context, glsn types.GLSN, onNextFunc OnNext) error

	Trim(ctx context.Context, glsn types.GLSN) error
}

type OnNext func(glsn types.GLSN, err error)

type varlog struct {
	clusterID         types.ClusterID
	refresher         *metadataRefresher
	lsSelector        LogStreamSelector
	replicasRetriever ReplicasRetriever
	allowlist         Allowlist

	logCLManager logc.LogClientManager
	logger       *zap.Logger
	opts         *options

	closed atomicutil.AtomicBool
}

var _ Varlog = (*varlog)(nil)

// Open creates new logs or opens an already created logs.
func Open(clusterID types.ClusterID, mrAddrs []string, opts ...Option) (Varlog, error) {
	logOpts := defaultOptions()
	for _, opt := range opts {
		opt.apply(&logOpts)
	}
	logOpts.logger = logOpts.logger.Named("varlog").With(zap.Any("cid", clusterID))

	v := &varlog{
		clusterID: clusterID,
		logger:    logOpts.logger,
		opts:      &logOpts,
	}

	// mr connector
	connectorOpts := []mrconnector.Option{
		mrconnector.WithClusterID(clusterID),
		mrconnector.WithConnectionTimeout(v.opts.openTimeout),
		mrconnector.WithLogger(v.opts.logger),
	}
	connector, err := mrconnector.New(context.Background(), mrAddrs, connectorOpts...)
	if err != nil {
		return nil, err
	}

	// allowlist
	allowlist, err := newTransientAllowlist(v.opts.denyTTL, v.logger)
	if err != nil {
		return nil, err
	}
	v.allowlist = allowlist

	// log stream selector
	v.lsSelector = newAppendableLogStreamSelector(allowlist)

	// replicas retriever
	replicasRetriever := &renewableReplicasRetriever{}
	v.replicasRetriever = replicasRetriever

	// metadata refresher
	// NOTE: metadata refresher has ownership of mrconnector and allow/denylist
	ctx, cancel := context.WithTimeout(context.Background(), v.opts.openTimeout)
	defer cancel()
	refresher, err := newMetadataRefresher(ctx, connector, allowlist, replicasRetriever, v.opts.metadataRefreshInterval, v.logger)
	if err != nil {
		return nil, err
	}
	v.refresher = refresher

	// logcl manager
	// TODO (jun): metadataRefresher should implement ClusterMetadataView
	metadata := refresher.metadata.Load().(*varlogpb.MetadataDescriptor)
	logCLManager, err := logc.NewLogClientManager(metadata, v.logger)
	if err != nil {
		return nil, err
	}
	v.logCLManager = logCLManager

	return v, nil
}

func (v *varlog) Append(ctx context.Context, data []byte, opts ...AppendOption) (types.GLSN, error) {
	return v.append(ctx, 0, data, opts...)
}

func (v *varlog) AppendTo(ctx context.Context, logStreamID types.LogStreamID, data []byte, opts ...AppendOption) (types.GLSN, error) {
	opts = append(opts, withoutSelectLogStream())
	return v.append(ctx, logStreamID, data, opts...)
}

func (v *varlog) Read(ctx context.Context, logStreamID types.LogStreamID, glsn types.GLSN) ([]byte, error) {
	logEntry, err := v.read(ctx, logStreamID, glsn)
	if err != nil {
		return nil, err
	}
	return logEntry.Data, nil
}

func (v *varlog) Subscribe(ctx context.Context, glsn types.GLSN, onNextFunc OnNext) error {
	panic("not implemented")
}

func (v *varlog) Trim(ctx context.Context, glsn types.GLSN) error {
	panic("not implemented")
}

func (v *varlog) Close() (err error) {
	if v.closed.Load() {
		return
	}

	if e := v.logCLManager.Close(); e != nil {
		v.logger.Warn("error while closing logcl manager", zap.Error(err))
		err = e
	}
	if e := v.refresher.Close(); e != nil {
		v.logger.Warn("error while closing metadata refresher", zap.Error(err))
		err = e
	}
	v.closed.Store(true)
	return
}
