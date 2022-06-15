package admin

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/kakao/varlog/internal/admin/admerrors"
	"github.com/kakao/varlog/internal/admin/snwatcher"
	"github.com/kakao/varlog/pkg/types"
	"github.com/kakao/varlog/pkg/util/netutil"
	"github.com/kakao/varlog/pkg/verrors"
	"github.com/kakao/varlog/proto/mrpb"
	"github.com/kakao/varlog/proto/snpb"
	"github.com/kakao/varlog/proto/varlogpb"
	"github.com/kakao/varlog/proto/vmspb"
)

const numLogStreamMutex = 512

type Admin struct {
	config

	closed       bool
	lis          net.Listener
	server       *grpc.Server
	serverAddr   string
	healthServer *health.Server

	// single large lock
	mu                sync.RWMutex
	muLogStreamStatus [numLogStreamMutex]sync.Mutex

	snw     *snwatcher.StorageNodeWatcher
	lsidGen *LogStreamIDGenerator
	tpidGen *TopicIDGenerator
}

// New creates an Admin.
// It should be created only once.
// It returns an error if cluster metadata cannot be fetched.
func New(ctx context.Context, opts ...Option) (*Admin, error) {
	cfg, err := newConfig(opts)
	if err != nil {
		return nil, err
	}

	cmView := cfg.mrmgr.ClusterMetadataView()

	logStreamIDGen, err := NewLogStreamIDGenerator(ctx, cmView)
	if err != nil {
		return nil, err
	}

	topicIDGen, err := NewTopicIDGenerator(ctx, cmView)
	if err != nil {
		return nil, err
	}

	cm := &Admin{
		config:       cfg,
		lsidGen:      logStreamIDGen,
		tpidGen:      topicIDGen,
		server:       grpc.NewServer(),
		healthServer: health.NewServer(),
	}
	cm.snw, err = snwatcher.New(append(
		cm.snwatcherOpts,
		snwatcher.WithClusterMetadataView(cmView),
		snwatcher.WithStorageNodeManager(cm.snmgr),
		snwatcher.WithStatisticsRepository(cm.statRepository),
		snwatcher.WithStorageNodeWatcherHandler(cm),
		snwatcher.WithLogger(cm.logger),
	)...)
	return cm, err
}

// Address returns the bound address of the admin server.
func (adm *Admin) Address() string {
	adm.mu.RLock()
	defer adm.mu.RUnlock()
	return adm.serverAddr
}

// Serve accepts incoming RPC calls on the server.
// This method blocks calling goroutine.
func (adm *Admin) Serve() error {
	adm.mu.Lock()
	if adm.lis != nil {
		adm.mu.Unlock()
		return errors.New("admin: already serving")
	}
	lis, err := net.Listen("tcp", adm.listenAddress)
	if err != nil {
		adm.mu.Unlock()
		return errors.WithStack(err)
	}
	adm.lis = lis
	addrs, _ := netutil.GetListenerAddrs(lis.Addr())
	adm.serverAddr = addrs[0]

	vmspb.RegisterClusterManagerServer(adm.server, &server{admin: adm})
	grpc_health_v1.RegisterHealthServer(adm.server, adm.healthServer)
	adm.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// SN Watcher
	if err := adm.snw.Start(); err != nil {
		adm.mu.Unlock()
		return err
	}
	adm.mu.Unlock()

	return adm.server.Serve(lis)
}

// Close closes the admin.
// This method closes the gRPC server immediately.
func (adm *Admin) Close() (err error) {
	adm.mu.Lock()
	defer adm.mu.Unlock()
	if adm.closed {
		return nil
	}
	adm.closed = true

	adm.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	// SN Watcher
	err = adm.snw.Stop()
	err = multierr.Combine(err, adm.snmgr.Close(), adm.mrmgr.Close())
	adm.server.Stop()
	return err
}

// Metadata returns the metadata of cluster.
//
// Deprecated: Only integration test code calls this method, but they are not allowed call this method directly.
func (adm *Admin) Metadata(ctx context.Context) (*varlogpb.MetadataDescriptor, error) {
	adm.mu.RLock()
	defer adm.mu.RUnlock()
	return adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
}

func (adm *Admin) getStorageNode(ctx context.Context, snid types.StorageNodeID) (*vmspb.StorageNodeMetadata, error) {
	snm, ok := adm.statRepository.GetStorageNode(snid)
	if !ok {
		return nil, errors.WithMessagef(admerrors.ErrNoSuchStorageNode, "get storage node %d", int32(snid))
	}
	return snm, nil
}

func (adm *Admin) listStorageNodes(ctx context.Context) (snmds map[types.StorageNodeID]*vmspb.StorageNodeMetadata, err error) {
	adm.mu.RLock()
	defer adm.mu.RUnlock()
	return adm.statRepository.ListStorageNodes(), nil
}

// addStorageNode adds a new storage node to the cluster.
// It is idempotent, that is, adding an already added storage node is okay.
//
// It could not add a storage node under the following conditions:
//  - It could not fetch metadata from the storage node.
//  - It is rejected by the metadata repository.
func (adm *Admin) addStorageNode(ctx context.Context, snid types.StorageNodeID, addr string) (*snpb.StorageNodeMetadataDescriptor, error) {
	adm.mu.Lock()
	defer adm.mu.Unlock()

	snmd, err := adm.snmgr.GetMetadataByAddress(ctx, snid, addr)
	if err != nil {
		return nil, err
	}

	md, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return nil, err
	}
	// If there is a storage node whose ID and address are the same as the
	// arguments snid and addr, it will succeed. If only one of them is the
	// same, it will fail by the metadata repository.
	if snd := md.GetStorageNode(snid); snd != nil && snd.Address == addr {
		return snmd, nil
	}

	now := time.Now().UTC()
	snd := snmd.ToStorageNodeDescriptor()
	snd.Status = varlogpb.StorageNodeStatusRunning
	snd.CreateTime = now
	if err = adm.mrmgr.RegisterStorageNode(ctx, snd); err != nil {
		return nil, err
	}

	adm.snmgr.AddStorageNode(ctx, snmd.StorageNode.StorageNodeID, addr)
	adm.statRepository.Report(ctx, snmd, now)
	return snmd, err
}

func (adm *Admin) unregisterStorageNode(ctx context.Context, snid types.StorageNodeID) error {
	adm.mu.Lock()
	defer adm.mu.Unlock()

	clusmeta, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return err
	}

	if clusmeta.GetStorageNode(snid) == nil {
		return nil
	}

	// TODO (jun): Use helper function
	for _, lsdesc := range clusmeta.GetLogStreams() {
		for _, replica := range lsdesc.GetReplicas() {
			if replica.GetStorageNodeID() == snid {
				return errors.WithMessagef(admerrors.ErrNotIdleReplicas, "unregister storage node")
			}
		}
	}

	if err := adm.mrmgr.UnregisterStorageNode(ctx, snid); err != nil {
		return err
	}

	adm.snmgr.RemoveStorageNode(snid)
	adm.statRepository.RemoveStorageNode(snid)
	return nil
}

func (adm *Admin) getTopic(ctx context.Context, tpid types.TopicID) (*varlogpb.TopicDescriptor, error) {
	md, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return nil, err
	}
	td := md.GetTopic(tpid)
	if td == nil {
		return nil, errors.WithMessagef(admerrors.ErrNoSuchTopic, "get topic %d", int32(tpid))
	}
	return td, nil
}

func (adm *Admin) listTopics(ctx context.Context) ([]varlogpb.TopicDescriptor, error) {
	adm.mu.Lock()
	defer adm.mu.Unlock()

	md, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil || len(md.Topics) == 0 {
		return nil, err
	}

	tds := make([]varlogpb.TopicDescriptor, len(md.Topics))
	for idx := range md.Topics {
		tds[idx] = *md.Topics[idx]
	}
	return tds, nil
}

func (adm *Admin) addTopic(ctx context.Context) (varlogpb.TopicDescriptor, error) {
	adm.mu.Lock()
	defer adm.mu.Unlock()

	topicID := adm.tpidGen.Generate()
	// Note that the metadata repository accepts redundant RegisterTopic
	// RPC only if the topic has no log streams.
	if err := adm.mrmgr.RegisterTopic(ctx, topicID); err != nil {
		return varlogpb.TopicDescriptor{}, err
	}

	return varlogpb.TopicDescriptor{TopicID: topicID}, nil
}

func (adm *Admin) unregisterTopic(ctx context.Context, tpid types.TopicID) error {
	adm.mu.Lock()
	defer adm.mu.Unlock()

	clusmeta, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return err
	}

	// TODO: Should it returns an error when removing the topic that has already been deleted or does not exist?
	topicdesc, err := clusmeta.MustHaveTopic(tpid)
	if err != nil {
		return err
	}

	// TODO: Should it returns an error when removing the topic that has already been deleted or does not exist?
	status := topicdesc.GetStatus()
	if status.Deleted() {
		return errors.Errorf("invalid topic status: %s", status)
	}

	// TODO: Can we remove the topic that has active log streams?
	// TODO:: seal logStreams and refresh metadata
	return adm.mrmgr.UnregisterTopic(ctx, tpid)
}

func (adm *Admin) getLogStream(ctx context.Context, tpid types.TopicID, lsid types.LogStreamID) (*varlogpb.LogStreamDescriptor, error) {
	md, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get log stream")
	}
	td := md.GetTopic(tpid)
	if td == nil {
		return nil, errors.WithMessagef(admerrors.ErrNoSuchTopic, "get log stream: tpid %d", int32(tpid))
	}
	if !td.HasLogStream(lsid) {
		return nil, errors.WithMessagef(admerrors.ErrNoSuchLogStream, "get log stream: no log stream in topic: lsid %d", int32(lsid))
	}
	lsd := md.GetLogStream(lsid)
	if lsd == nil {
		return nil, errors.WithMessagef(admerrors.ErrNoSuchLogStream, "get log stream: lsid %d", int32(lsid))
	}
	if lsd.TopicID != tpid {
		return nil, fmt.Errorf("get log stream: unexpected topic: expected %d, actual %d", int32(tpid), int32(lsd.TopicID))
	}
	return lsd, nil
}

func (adm *Admin) listLogStreams(ctx context.Context, tpid types.TopicID) ([]*varlogpb.LogStreamDescriptor, error) {
	md, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "list log streams")
	}
	td := md.GetTopic(tpid)
	if td == nil {
		return nil, errors.WithMessagef(admerrors.ErrNoSuchTopic, "list log streams: tpid %d", int32(tpid))
	}
	lsds := make([]*varlogpb.LogStreamDescriptor, 0, len(td.LogStreams))
	for _, lsid := range td.LogStreams {
		lsd := md.GetLogStream(lsid)
		if lsd == nil {
			continue
		}
		if lsd.TopicID != tpid {
			return nil, fmt.Errorf("list log streams: unexpected topic: expected %d, actual %d", int32(tpid), int32(lsd.TopicID))
		}
		lsds = append(lsds, lsd)
	}
	return lsds, nil
}

func (adm *Admin) describeTopic(ctx context.Context, tpid types.TopicID) (td varlogpb.TopicDescriptor, lsds []varlogpb.LogStreamDescriptor, err error) {
	adm.mu.Lock()
	defer adm.mu.Unlock()

	md, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil || len(md.Topics) == 0 {
		return
	}

	tdPtr := md.GetTopic(tpid)
	if tdPtr == nil {
		err = errors.Wrapf(verrors.ErrNotExist, "no such topic (topicID=%d)", tpid)
		return
	}
	td = *proto.Clone(tdPtr).(*varlogpb.TopicDescriptor)
	lsds = make([]varlogpb.LogStreamDescriptor, 0, len(td.LogStreams))
	for _, lsID := range td.LogStreams {
		lsdPtr := md.GetLogStream(lsID)
		if lsdPtr == nil {
			continue
		}
		lsd := *proto.Clone(lsdPtr).(*varlogpb.LogStreamDescriptor)
		lsds = append(lsds, lsd)
	}

	return td, lsds, nil
}

// addLogStream adds a new log stream in the topic specified by the argument
// tpid.
// It returns an error if the topic does not exist.
// The argument replicas can be empty, then, this method selects proper
// replicas.
// If the argument replicas are defined, each storage node for the replica must
// exist.
func (adm *Admin) addLogStream(ctx context.Context, tpid types.TopicID, replicas []*varlogpb.ReplicaDescriptor) (*varlogpb.LogStreamDescriptor, error) {
	lsdesc, err := adm.addLogStreamInternal(ctx, tpid, replicas)
	if err != nil {
		return lsdesc, err
	}

	if adm.enableAutoUnseal {
		if err := adm.waitUnsealed(ctx, lsdesc.LogStreamID); err != nil {
			return nil, err
		}
		// FIXME: Failure of Seal does not mean failure of AddLogStream.
		clusmeta, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
		if err != nil {
			return nil, err
		}
		return clusmeta.GetLogStream(lsdesc.LogStreamID), nil
	}

	// FIXME: Failure of Seal does not mean failure of AddLogStream.
	err = adm.waitSealed(ctx, lsdesc.LogStreamID)
	if err != nil {
		return lsdesc, err
	}

	// FIXME: Failure of Unseal does not mean failure of AddLogStream.
	return adm.unseal(ctx, tpid, lsdesc.LogStreamID)
}

func (adm *Admin) addLogStreamInternal(ctx context.Context, tpid types.TopicID, replicas []*varlogpb.ReplicaDescriptor) (*varlogpb.LogStreamDescriptor, error) {
	adm.mu.Lock()
	defer adm.mu.Unlock()

	var err error

	if len(replicas) == 0 {
		replicas, err = adm.snSelector.Select(ctx)
		if err != nil {
			return nil, err
		}
	}

	// See https://github.com/kakao/varlog/pull/198#discussion_r215602
	logStreamID := adm.lsidGen.Generate()

	// duplicated by verifyLogStream
	/*
		if err := clusmeta.MustNotHaveLogStream(logStreamID); err != nil {
			if e := adm.lsidGen.Refresh(ctx); e != nil {
				err = multierr.Append(err, e)
				adm.logger.Panic("could not refresh LogStreamIDGenerator", zap.Error(err))
			}
			return nil, err
		}
	*/

	logStreamDesc := &varlogpb.LogStreamDescriptor{
		TopicID:     tpid,
		LogStreamID: logStreamID,
		Status:      varlogpb.LogStreamStatusSealing,
		Replicas:    replicas,
	}

	clusmeta, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return nil, err
	}
	if err := adm.verifyLogStream(clusmeta, logStreamDesc); err != nil {
		return nil, err
	}

	// TODO: Choose the primary - e.g., shuffle logStreamReplicaMetas
	if err := adm.snmgr.AddLogStream(ctx, logStreamDesc); err != nil {
		return nil, err
	}

	// NB: RegisterLogStream returns nil if the logstream already exists.
	return logStreamDesc, adm.mrmgr.RegisterLogStream(ctx, logStreamDesc)
}

func (adm *Admin) verifyLogStream(clusmeta *varlogpb.MetadataDescriptor, lsdesc *varlogpb.LogStreamDescriptor) error {
	replicas := lsdesc.GetReplicas()
	// the number of logstream replica
	if uint(len(replicas)) != adm.replicationFactor {
		return errors.Errorf("invalid number of log stream replicas: %d", len(replicas))
	}
	// storagenode existence
	for _, replica := range replicas {
		if _, err := clusmeta.MustHaveStorageNode(replica.GetStorageNodeID()); err != nil {
			return err
		}
	}
	// logstream existence
	if err := clusmeta.MustNotHaveLogStream(lsdesc.GetLogStreamID()); err != nil {
		_ = adm.lsidGen.Refresh(context.TODO())
		return err
	}
	return nil
}

func (adm *Admin) waitSealed(ctx context.Context, lsid types.LogStreamID) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			lsStat := adm.statRepository.GetLogStream(lsid).Copy()
			if lsStat.Status() == varlogpb.LogStreamStatusSealed {
				return nil
			}
		}
	}
}

func (adm *Admin) waitUnsealed(ctx context.Context, lsid types.LogStreamID) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			md, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
			if err != nil {
				return err
			}
			if md.GetLogStream(lsid).Status != varlogpb.LogStreamStatusRunning {
				continue
			}
			lsStat := adm.statRepository.GetLogStream(lsid).Copy()
			if lsStat.Status() == varlogpb.LogStreamStatusRunning {
				return nil
			}
		}
	}
}

func (adm *Admin) updateLogStream(ctx context.Context, lsid types.LogStreamID, poppedReplica, pushedReplica *varlogpb.ReplicaDescriptor) (*varlogpb.LogStreamDescriptor, error) {
	// NOTE (jun): Name of the method - updateLogStream can be confused.
	// updateLogStream can change only replicas. To update status, use Seal or Unseal.
	adm.mu.Lock()
	defer adm.mu.Unlock()

	adm.lockLogStreamStatus(lsid)
	defer adm.unlockLogStreamStatus(lsid)

	clusmeta, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return nil, err
	}

	oldLSDesc, err := clusmeta.MustHaveLogStream(lsid)
	if err != nil {
		return nil, err
	}

	status := oldLSDesc.GetStatus()
	if status.Running() || status.Deleted() {
		return nil, errors.Errorf("invalid log stream status: %s", status)
	}

	if poppedReplica == nil {
		// TODO: Choose laggy replica
		selector := newVictimSelector(adm.snmgr, lsid, oldLSDesc.GetReplicas())
		victims, err := selector.Select(ctx)
		if err != nil {
			return nil, err
		}
		poppedReplica = victims[0]
	}

	if pushedReplica == nil {
		oldReplicas := oldLSDesc.GetReplicas()
		denylist := make([]types.StorageNodeID, len(oldReplicas))
		for i, replica := range oldReplicas {
			denylist[i] = replica.GetStorageNodeID()
		}

		selector, err := newRandomReplicaSelector(adm.mrmgr.ClusterMetadataView(), 1, denylist...)
		if err != nil {
			return nil, err
		}
		candidates, err := selector.Select(ctx)
		if err != nil {
			return nil, err
		}
		pushedReplica = candidates[0]
	}

	replace := false
	newLSDesc := proto.Clone(oldLSDesc).(*varlogpb.LogStreamDescriptor)
	for i := range newLSDesc.Replicas {
		// TODO - fix? poppedReplica can ignore path.
		if newLSDesc.Replicas[i].GetStorageNodeID() == poppedReplica.GetStorageNodeID() {
			newLSDesc.Replicas[i] = pushedReplica
			replace = true
			break
		}
	}
	if !replace {
		adm.logger.Panic("logstream push/pop error")
	}

	if err := adm.snmgr.AddLogStreamReplica(ctx, pushedReplica.GetStorageNodeID(), newLSDesc.TopicID, lsid, pushedReplica.GetPath()); err != nil {
		return nil, err
	}

	// To reset the status of the log stream, set it as LogStreamStatusRunning
	defer func() {
		adm.statRepository.SetLogStreamStatus(lsid, varlogpb.LogStreamStatusRunning)
	}()

	if err := adm.mrmgr.UpdateLogStream(ctx, newLSDesc); err != nil {
		return nil, err
	}

	return newLSDesc, nil
}

func (adm *Admin) unregisterLogStream(ctx context.Context, tpid types.TopicID, lsid types.LogStreamID) error {
	adm.mu.Lock()
	defer adm.mu.Unlock()

	adm.lockLogStreamStatus(lsid)
	defer adm.unlockLogStreamStatus(lsid)

	clusmeta, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return err
	}

	lsdesc, err := clusmeta.MustHaveLogStream(lsid)
	if err != nil {
		return err
	}

	status := lsdesc.GetStatus()
	// TODO (jun): Check whether status.Deleted means unregistered.
	// If so, is status.Deleted okay or not?
	if status.Running() || status.Deleted() {
		return errors.Errorf("invalid log stream status: %s", status)
	}

	// TODO (jun): test if the log stream has no logs

	return adm.mrmgr.UnregisterLogStream(ctx, lsid)
}

func (adm *Admin) removeLogStreamReplica(ctx context.Context, snid types.StorageNodeID, tpid types.TopicID, lsid types.LogStreamID) error {
	adm.mu.Lock()
	defer adm.mu.Unlock()

	adm.lockLogStreamStatus(lsid)
	defer adm.unlockLogStreamStatus(lsid)

	clusmeta, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return err
	}

	if err := adm.removableLogStreamReplica(clusmeta, snid, lsid); err != nil {
		return err
	}

	return adm.snmgr.RemoveLogStreamReplica(ctx, snid, tpid, lsid)
}

func (adm *Admin) removableLogStreamReplica(clusmeta *varlogpb.MetadataDescriptor, snid types.StorageNodeID, lsid types.LogStreamID) error {
	lsdesc := clusmeta.GetLogStream(lsid)
	if lsdesc == nil {
		// unregistered LS or garbage
		return nil
	}

	replicas := lsdesc.GetReplicas()
	for _, replica := range replicas {
		if replica.GetStorageNodeID() == snid {
			return errors.Wrap(verrors.ErrState, "running log stream is not removable")
		}
	}
	return nil
}

func (adm *Admin) lockLogStreamStatus(lsid types.LogStreamID) {
	adm.muLogStreamStatus[lsid%numLogStreamMutex].Lock()
}

func (adm *Admin) unlockLogStreamStatus(lsid types.LogStreamID) {
	adm.muLogStreamStatus[lsid%numLogStreamMutex].Unlock()
}

// seal seals the log stream identified by the argument tpid and lsid.
// FIXME (jun): Define the specification of the seal more concretely.
func (adm *Admin) seal(ctx context.Context, tpid types.TopicID, lsid types.LogStreamID) ([]snpb.LogStreamReplicaMetadataDescriptor, types.GLSN, error) {
	adm.lockLogStreamStatus(lsid)
	defer adm.unlockLogStreamStatus(lsid)

	return adm.sealInternal(ctx, tpid, lsid)
}

func (adm *Admin) sealInternal(ctx context.Context, tpid types.TopicID, lsid types.LogStreamID) ([]snpb.LogStreamReplicaMetadataDescriptor, types.GLSN, error) {
	adm.statRepository.SetLogStreamStatus(lsid, varlogpb.LogStreamStatusSealing)

	lastGLSN, err := adm.mrmgr.Seal(ctx, lsid)
	if err != nil {
		adm.statRepository.SetLogStreamStatus(lsid, varlogpb.LogStreamStatusRunning)
		return nil, types.InvalidGLSN, err
	}

	result, err := adm.snmgr.Seal(ctx, tpid, lsid, lastGLSN)
	if err != nil {
		adm.statRepository.SetLogStreamStatus(lsid, varlogpb.LogStreamStatusRunning)
	}

	return result, lastGLSN, err
}

// unseal unseals the log stream replicas corresponded with the given logStreamID.
func (adm *Admin) unseal(ctx context.Context, tpid types.TopicID, lsid types.LogStreamID) (*varlogpb.LogStreamDescriptor, error) {
	adm.lockLogStreamStatus(lsid)
	defer adm.unlockLogStreamStatus(lsid)

	return adm.unsealInternal(ctx, tpid, lsid)
}

func (adm *Admin) unsealInternal(ctx context.Context, tpid types.TopicID, lsid types.LogStreamID) (*varlogpb.LogStreamDescriptor, error) {
	var err error
	var clusmeta *varlogpb.MetadataDescriptor
	adm.statRepository.SetLogStreamStatus(lsid, varlogpb.LogStreamStatusUnsealing)

	if err = adm.snmgr.Unseal(ctx, tpid, lsid); err != nil {
		goto errOut
	}

	if err = adm.mrmgr.Unseal(ctx, lsid); err != nil {
		goto errOut
	}

	if clusmeta, err = adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx); err != nil {
		goto errOut
	}

	return clusmeta.GetLogStream(lsid), nil

errOut:
	adm.statRepository.SetLogStreamStatus(lsid, varlogpb.LogStreamStatusRunning)
	return nil, err
}

// sync copies the log entries of the src to the dst. sync may be long-running, thus it
// returns immediately without waiting for the completion of sync. Callers of sync
// periodically can call sync, and get the current state of the sync progress.
// SyncState is one of SyncStateError, SyncStateInProgress, or SyncStateComplete. If sync
// returns SyncStateComplete, all the log entries were copied well. If it returns
// SyncStateInProgress, it is still progressing. Otherwise, if it returns SyncStateError,
// it is stopped by an error.
// To start sync, the log stream status of the src must be LogStreamStatusSealed and the log
// stream status of the dst must be LogStreamStatusSealing. If either of the statuses is not
// correct, sync returns ErrSyncInvalidStatus.
func (adm *Admin) sync(ctx context.Context, tpid types.TopicID, lsid types.LogStreamID, srcID, dstID types.StorageNodeID) (*snpb.SyncStatus, error) {
	adm.lockLogStreamStatus(lsid)
	defer adm.unlockLogStreamStatus(lsid)

	return adm.syncInternal(ctx, tpid, lsid, srcID, dstID)
}

func (adm *Admin) syncInternal(ctx context.Context, tpid types.TopicID, lsid types.LogStreamID, srcID, dstID types.StorageNodeID) (*snpb.SyncStatus, error) {
	lastGLSN, err := adm.mrmgr.Seal(ctx, lsid)
	if err != nil {
		return nil, err
	}
	return adm.snmgr.Sync(ctx, tpid, lsid, srcID, dstID, lastGLSN)
}

// trim removes log entries from the log streams in a topic.
// The argument tpid is the topic ID of the topic to be trimmed.
// The argument lastGLSN is the last global sequence number of the log stream to be trimmed.
func (adm *Admin) trim(ctx context.Context, tpid types.TopicID, lastGLSN types.GLSN) ([]vmspb.TrimResult, error) {
	adm.mu.Lock()
	defer adm.mu.Unlock()
	return adm.snmgr.Trim(ctx, tpid, lastGLSN)
}

func (adm *Admin) getMetadataRepositoryNode(ctx context.Context, nid types.NodeID) (*varlogpb.MetadataRepositoryNode, error) {
	ci, err := adm.mrmgr.GetClusterInfo(ctx)
	if err != nil {
		return nil, err
	}
	member, ok := ci.GetMembers()[nid]
	if !ok {
		return nil, errors.New("admin: no such mr")
	}
	return &varlogpb.MetadataRepositoryNode{
		NodeID:  nid,
		RaftURL: member.Peer,
		RPCAddr: member.Endpoint,
		Leader:  ci.Leader == nid,
		Learner: member.Learner,
	}, nil
}

func (adm *Admin) listMetadataRepositoryNodes(ctx context.Context) ([]*varlogpb.MetadataRepositoryNode, error) {
	ci, err := adm.mrmgr.GetClusterInfo(ctx)
	if err != nil {
		return nil, err
	}
	nodes := make([]*varlogpb.MetadataRepositoryNode, 0, len(ci.GetMembers()))
	for nid, member := range ci.GetMembers() {
		nodes = append(nodes, &varlogpb.MetadataRepositoryNode{
			NodeID:  nid,
			RaftURL: member.Peer,
			RPCAddr: member.Endpoint,
			Leader:  ci.Leader == nid,
			Learner: member.Learner,
		})
	}
	return nodes, nil
}

func (adm *Admin) mrInfos(ctx context.Context) (*mrpb.ClusterInfo, error) {
	adm.mu.RLock()
	defer adm.mu.RUnlock()
	return adm.mrmgr.GetClusterInfo(ctx)
}

func (adm *Admin) addMetadataRepositoryNode(ctx context.Context, raftURL, rpcAddr string) (*varlogpb.MetadataRepositoryNode, error) {
	nid := types.NewNodeIDFromURL(raftURL)
	if nid == types.InvalidNodeID {
		return nil, errors.Wrap(verrors.ErrInvalid, "raft address")
	}

	if err := adm.mrmgr.AddPeer(ctx, nid, raftURL, rpcAddr); err != nil {
		if !errors.Is(err, verrors.ErrAlreadyExists) {
			return nil, err
		}
	}

	return &varlogpb.MetadataRepositoryNode{
		NodeID:  nid,
		RaftURL: raftURL,
		RPCAddr: rpcAddr,
		// TODO: Fill these fields.
		// Leader:
		// Learner:
	}, nil
}

func (adm *Admin) addMRPeer(ctx context.Context, raftURL, rpcAddr string) (types.NodeID, error) {
	nodeID := types.NewNodeIDFromURL(raftURL)
	if nodeID == types.InvalidNodeID {
		return nodeID, errors.Wrap(verrors.ErrInvalid, "raft address")
	}

	adm.mu.RLock()
	defer adm.mu.RUnlock()

	err := adm.mrmgr.AddPeer(ctx, nodeID, raftURL, rpcAddr)
	if err != nil {
		if !errors.Is(err, verrors.ErrAlreadyExists) {
			return types.InvalidNodeID, err
		}
	}

	return nodeID, nil
}

func (adm *Admin) deleteMetadataRepositoryNode(ctx context.Context, nid types.NodeID) error {
	if nid == types.InvalidNodeID {
		return errors.Wrap(verrors.ErrInvalid, "raft address")
	}

	if err := adm.mrmgr.RemovePeer(ctx, nid); err != nil {
		// TODO: What does it mean that RemovePeer returns an error of verrors.ErrAlreadyExists?
		if !errors.Is(err, verrors.ErrAlreadyExists) {
			return err
		}
	}

	return nil
}

func (adm *Admin) removeMRPeer(ctx context.Context, raftURL string) error {
	nodeID := types.NewNodeIDFromURL(raftURL)
	if nodeID == types.InvalidNodeID {
		return errors.Wrap(verrors.ErrInvalid, "raft address")
	}

	adm.mu.RLock()
	defer adm.mu.RUnlock()

	err := adm.mrmgr.RemovePeer(ctx, nodeID)
	if err != nil {
		if !errors.Is(err, verrors.ErrAlreadyExists) {
			return err
		}
	}

	return nil
}

func (adm *Admin) HandleHeartbeatTimeout(ctx context.Context, snid types.StorageNodeID) {
	meta, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return
	}

	//TODO: store sn status
	for _, ls := range meta.GetLogStreams() {
		if ls.IsReplica(snid) {
			adm.logger.Debug("seal due to heartbeat timeout", zap.Any("snid", snid), zap.Any("lsid", ls.LogStreamID))
			adm.seal(ctx, ls.TopicID, ls.LogStreamID)
		}
	}
}

func (adm *Admin) checkLogStreamStatus(ctx context.Context, tpid types.TopicID, lsid types.LogStreamID, mrStatus, replicaStatus varlogpb.LogStreamStatus) {
	adm.lockLogStreamStatus(lsid)
	defer adm.unlockLogStreamStatus(lsid)

	lsStat := adm.statRepository.GetLogStream(lsid).Copy()

	switch lsStat.Status() {
	case varlogpb.LogStreamStatusRunning:
		if mrStatus.Sealed() || replicaStatus.Sealed() {
			adm.logger.Info("seal due to status mismatch", zap.Any("lsid", lsid))
			adm.sealInternal(ctx, tpid, lsid)
		}

	case varlogpb.LogStreamStatusSealing:
		for _, r := range lsStat.Replicas() {
			if r.Status != varlogpb.LogStreamStatusSealed {
				adm.logger.Info("seal due to status", zap.Any("lsid", lsid))
				adm.sealInternal(ctx, tpid, lsid)
				return
			}
		}
		adm.logger.Info("sealed", zap.Any("lsid", lsid))
		adm.statRepository.SetLogStreamStatus(lsid, varlogpb.LogStreamStatusSealed)

	case varlogpb.LogStreamStatusSealed:
		for _, r := range lsStat.Replicas() {
			if r.Status != varlogpb.LogStreamStatusSealed {
				adm.statRepository.SetLogStreamStatus(lsid, varlogpb.LogStreamStatusSealing)
				return
			}
		}

	case varlogpb.LogStreamStatusUnsealing:
		for _, r := range lsStat.Replicas() {
			if r.Status == varlogpb.LogStreamStatusRunning {
				continue
			} else if r.Status == varlogpb.LogStreamStatusSealed {
				return
			} else if r.Status == varlogpb.LogStreamStatusSealing {
				adm.logger.Info("seal due to unexpected status", zap.Any("lsid", lsid))
				adm.sealInternal(ctx, tpid, lsid)
				return
			}
		}
		adm.statRepository.SetLogStreamStatus(lsid, varlogpb.LogStreamStatusRunning)
	}
}

func (adm *Admin) syncLogStream(ctx context.Context, topicID types.TopicID, logStreamID types.LogStreamID) {
	adm.lockLogStreamStatus(logStreamID)
	defer adm.unlockLogStreamStatus(logStreamID)

	min, max := types.MaxVersion, types.InvalidVersion
	var src, tgt types.StorageNodeID

	lsStat := adm.statRepository.GetLogStream(logStreamID).Copy()

	if !lsStat.Status().Sealed() {
		return
	}

	snIDs := make([]types.StorageNodeID, 0, len(lsStat.Replicas()))
	for snID := range lsStat.Replicas() {
		snIDs = append(snIDs, snID)
	}
	sort.Slice(snIDs, func(i, j int) bool { return snIDs[i] < snIDs[j] })

	sealed := true
	for i, snID := range snIDs {
		r, _ := lsStat.Replica(snID)

		if !r.Status.Sealed() {
			return
		}

		if r.Status != varlogpb.LogStreamStatusSealed {
			sealed = false
		}

		if r.Status == varlogpb.LogStreamStatusSealing && (i == 0 || r.Version < min) {
			min = r.Version
			tgt = snID
		}

		if r.Status == varlogpb.LogStreamStatusSealed && (i == 0 || r.Version > max) {
			max = r.Version
			src = snID
		}
	}

	// TODO (jun): Extract this block to separate method.
	if sealed && adm.enableAutoUnseal {
		lsd, err := adm.unsealInternal(context.Background(), topicID, logStreamID)
		if err != nil {
			adm.logger.Error("could not unseal",
				zap.Int32("tpid", int32(topicID)),
				zap.Int32("lsid", int32(logStreamID)),
				zap.Error(err),
			)
			return
		}
		adm.logger.Debug("unseal",
			zap.Int32("tpid", int32(topicID)),
			zap.Int32("lsid", int32(logStreamID)),
			zap.String("lsd", lsd.String()),
		)
		return
	}

	// FIXME (jun): Since there is no invalid identifier for the storage
	// node, it cannot check whether a source or target is selected. It
	// thus checks that min and max are in a valid range.
	if src != tgt && !max.Invalid() && min != types.MaxVersion {
		logger := adm.logger.With(
			zap.Int32("tpid", int32(topicID)),
			zap.Int32("lsid", int32(logStreamID)),
			zap.Int32("src", int32(src)),
			zap.Int32("dst", int32(tgt)),
		)
		status, err := adm.syncInternal(ctx, topicID, logStreamID, src, tgt)
		if err != nil {
			logger.Error("could not sync", zap.Error(err))
			return
		}
		logger.Debug("sync", zap.String("status", status.String()))
	}
}

func (adm *Admin) HandleReport(ctx context.Context, snm *snpb.StorageNodeMetadataDescriptor) {
	meta, err := adm.mrmgr.ClusterMetadataView().ClusterMetadata(ctx)
	if err != nil {
		return
	}

	// Sync LogStreamStatus
	for _, ls := range snm.GetLogStreamReplicas() {
		mls := meta.GetLogStream(ls.LogStreamID)
		if mls != nil {
			adm.checkLogStreamStatus(ctx, ls.TopicID, ls.LogStreamID, mls.Status, ls.Status)
			continue
		}
		if time.Since(ls.CreatedTime) > adm.logStreamGCTimeout {
			adm.removeLogStreamReplica(ctx, snm.StorageNode.StorageNodeID, ls.TopicID, ls.LogStreamID)
		}
	}

	// Sync LogStream
	if adm.disableAutoLogStreamSync {
		return
	}
	for _, ls := range snm.GetLogStreamReplicas() {
		if ls.Status.Sealed() {
			adm.syncLogStream(ctx, ls.TopicID, ls.LogStreamID)
		}
	}
}
