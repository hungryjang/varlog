package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/kakao/varlog/internal/metadata_repository"
	"github.com/kakao/varlog/internal/storagenode"
	"github.com/kakao/varlog/internal/vms"
	"github.com/kakao/varlog/pkg/logc"
	"github.com/kakao/varlog/pkg/rpc"
	"github.com/kakao/varlog/pkg/types"
	"github.com/kakao/varlog/pkg/util/testutil"
	"github.com/kakao/varlog/pkg/util/testutil/ports"
	"github.com/kakao/varlog/pkg/varlog"
	"github.com/kakao/varlog/pkg/verrors"
	"github.com/kakao/varlog/proto/varlogpb"
	"github.com/kakao/varlog/vtesting"
)

const varlogClusterPortBase = 10000
const varlogClusterManagerPortBase = 999

type VarlogClusterOptions struct {
	NrRep                 int
	NrMR                  int
	SnapCount             int
	CollectorName         string
	UnsafeNoWal           bool
	ReporterClientFac     metadata_repository.ReporterClientFactory
	SNManagementClientFac metadata_repository.StorageNodeManagementClientFactory
	VMSOpts               *vms.Options
}

type VarlogCluster struct {
	VarlogClusterOptions
	MRs                []*metadata_repository.RaftMetadataRepository
	SNs                map[types.StorageNodeID]*storagenode.StorageNode
	snWG               sync.WaitGroup
	volumes            map[types.StorageNodeID]storagenode.Volume
	snRPCEndpoints     map[types.StorageNodeID]string
	CM                 vms.ClusterManager
	CMCli              varlog.ClusterManagerClient
	MRPeers            []string
	MRRPCEndpoints     []string
	MRIDs              []types.NodeID
	snID               types.StorageNodeID
	lsID               types.LogStreamID
	issuedLogStreamIDs []types.LogStreamID
	ClusterID          types.ClusterID
	logger             *zap.Logger

	portLease *ports.Lease
}

func NewVarlogCluster(opts VarlogClusterOptions) *VarlogCluster {
	portLease, err := ports.ReserveWeaklyWithRetry(varlogClusterPortBase)
	if err != nil {
		panic(err)
	}

	mrPeers := make([]string, opts.NrMR)
	mrRPCEndpoints := make([]string, opts.NrMR)
	MRs := make([]*metadata_repository.RaftMetadataRepository, opts.NrMR)
	mrIDs := make([]types.NodeID, opts.NrMR)

	for i := range mrPeers {
		raftPort := i*2 + portLease.Base()
		rpcPort := i*2 + 1 + portLease.Base()

		mrPeers[i] = fmt.Sprintf("http://127.0.0.1:%d", raftPort)
		mrRPCEndpoints[i] = fmt.Sprintf("127.0.0.1:%d", rpcPort)
	}

	clus := &VarlogCluster{
		VarlogClusterOptions: opts,
		logger:               zap.L(),
		MRPeers:              mrPeers,
		MRRPCEndpoints:       mrRPCEndpoints,
		MRIDs:                mrIDs,
		MRs:                  MRs,
		SNs:                  make(map[types.StorageNodeID]*storagenode.StorageNode),
		volumes:              make(map[types.StorageNodeID]storagenode.Volume),
		snRPCEndpoints:       make(map[types.StorageNodeID]string),
		ClusterID:            types.ClusterID(1),
		portLease:            portLease,
	}

	for i := range clus.MRPeers {
		clus.clearMR(i)
		clus.createMR(i, false, clus.UnsafeNoWal, false)
	}

	return clus
}

func (clus *VarlogCluster) clearMR(idx int) {
	if idx < 0 || idx >= len(clus.MRs) {
		return
	}

	url, _ := url.Parse(clus.MRPeers[idx])
	nodeID := types.NewNodeID(url.Host)

	os.RemoveAll(fmt.Sprintf("%s/wal/%d", vtesting.TestRaftDir(), nodeID))
	os.RemoveAll(fmt.Sprintf("%s/snap/%d", vtesting.TestRaftDir(), nodeID))
	os.RemoveAll(fmt.Sprintf("%s/sml/%d", vtesting.TestRaftDir(), nodeID))
}

func (clus *VarlogCluster) createMR(idx int, join, unsafeNoWal, recoverFromSML bool) error {
	if idx < 0 || idx >= len(clus.MRs) {
		return errors.New("out of range")
	}

	url, _ := url.Parse(clus.MRPeers[idx])
	nodeID := types.NewNodeID(url.Host)

	opts := &metadata_repository.MetadataRepositoryOptions{
		RaftOptions: metadata_repository.RaftOptions{
			Join:        join,
			UnsafeNoWal: unsafeNoWal,
			EnableSML:   unsafeNoWal,
			SnapCount:   uint64(clus.SnapCount),
			RaftTick:    vtesting.TestRaftTick(),
			RaftDir:     vtesting.TestRaftDir(),
			Peers:       clus.MRPeers,
		},

		ClusterID:                      clus.ClusterID,
		RaftAddress:                    clus.MRPeers[idx],
		RPCTimeout:                     vtesting.TimeoutAccordingToProcCnt(metadata_repository.DefaultRPCTimeout),
		NumRep:                         clus.NrRep,
		RecoverFromSML:                 recoverFromSML,
		RPCBindAddress:                 clus.MRRPCEndpoints[idx],
		ReporterClientFac:              clus.ReporterClientFac,
		StorageNodeManagementClientFac: clus.SNManagementClientFac,
		Logger:                         clus.logger,
	}

	opts.CollectorName = "nop"
	if clus.CollectorName != "" {
		opts.CollectorName = clus.CollectorName
	}
	opts.CollectorEndpoint = "localhost:55680"

	clus.MRIDs[idx] = nodeID
	clus.MRs[idx] = metadata_repository.NewRaftMetadataRepository(opts)

	return nil
}

func (clus *VarlogCluster) AppendMR() error {
	idx := len(clus.MRs)
	raftPort := 2*idx + clus.portLease.Base()
	rpcPort := 2*idx + 1 + clus.portLease.Base()
	clus.MRPeers = append(clus.MRPeers, fmt.Sprintf("http://127.0.0.1:%d", raftPort))
	clus.MRRPCEndpoints = append(clus.MRRPCEndpoints, fmt.Sprintf("127.0.0.1:%d", rpcPort))
	clus.MRIDs = append(clus.MRIDs, types.InvalidNodeID)
	clus.MRs = append(clus.MRs, nil)

	clus.clearMR(idx)

	return clus.createMR(idx, true, clus.UnsafeNoWal, false)
}

func (clus *VarlogCluster) StartMR(idx int) error {
	if idx < 0 || idx >= len(clus.MRs) {
		return errors.New("out of range")
	}

	clus.MRs[idx].Run()

	return nil
}

func (clus *VarlogCluster) Start() {
	clus.logger.Info("cluster start")
	for i := range clus.MRs {
		clus.StartMR(i)
	}
	clus.logger.Info("cluster complete")
}

func (clus *VarlogCluster) StopMR(idx int) error {
	if idx < 0 || idx >= len(clus.MRs) {
		return errors.New("out or range")
	}

	return clus.MRs[idx].Close()
}

func (clus *VarlogCluster) RestartMR(idx int) error {
	if idx < 0 || idx >= len(clus.MRs) {
		return errors.New("out of range")
	}

	clus.StopMR(idx)
	clus.createMR(idx, false, clus.UnsafeNoWal, clus.UnsafeNoWal)
	return clus.StartMR(idx)
}

func (clus *VarlogCluster) CloseMR(idx int) error {
	if idx < 0 || idx >= len(clus.MRs) {
		return errors.New("out or range")
	}

	err := clus.StopMR(idx)
	clus.clearMR(idx)

	clus.logger.Info("cluster node stop", zap.Int("idx", idx))

	return err
}

// Close closes all cluster MRs
func (clus *VarlogCluster) Close() error {
	var err error

	if clus.CM != nil {
		clus.CM.Close()
	}

	for i := range clus.MRPeers {
		err = multierr.Append(err, clus.CloseMR(i))
	}

	err = multierr.Append(err, os.RemoveAll(vtesting.TestRaftDir()))

	for _, sn := range clus.SNs {
		// TODO (jun): remove temporary directories
		snmeta, erri := sn.GetMetadata(context.TODO())
		if erri != nil {
			err = multierr.Append(err, erri)
			clus.logger.Warn("could not get meta", zap.Error(erri))
		}
		for _, storage := range snmeta.GetStorageNode().GetStorages() {
			dbpath := storage.GetPath()
			if dbpath == ":memory:" {
				continue
			}
			/* comment out for test
			if err := os.RemoveAll(dbpath); err != nil {
				clus.logger.Warn("could not remove dbpath", zap.String("path", dbpath), zap.Error(err))
			}
			*/
		}

		// TODO:: sn.Close() does not close connect
		sn.Close()
	}
	clus.snWG.Wait()

	clus.issuedLogStreamIDs = nil

	return multierr.Combine(err, clus.portLease.Release())
}

func (clus *VarlogCluster) HealthCheck() bool {
	for _, endpoint := range clus.MRRPCEndpoints {
		conn, err := rpc.NewConn(context.TODO(), endpoint)
		if err != nil {
			return false
		}
		defer conn.Close()

		healthClient := grpc_health_v1.NewHealthClient(conn.Conn)
		if _, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{}); err != nil {
			return false
		}
	}

	return true
}

func (clus *VarlogCluster) Leader() int {
	leader := -1
	for i, n := range clus.MRs {
		cinfo, _ := n.GetClusterInfo(context.TODO(), clus.ClusterID)
		if cinfo.GetLeader() != types.InvalidNodeID && clus.MRIDs[i] == cinfo.GetLeader() {
			leader = i
			break
		}
	}

	return leader
}

func (clus *VarlogCluster) LeaderElected() bool {
	for _, n := range clus.MRs {
		if cinfo, _ := n.GetClusterInfo(context.TODO(), clus.ClusterID); cinfo.GetLeader() == types.InvalidNodeID {
			return false
		}
	}

	return true
}

func (clus *VarlogCluster) LeaderFail() bool {
	leader := clus.Leader()
	if leader < 0 {
		return false
	}

	clus.StopMR(leader)
	return true
}

func (clus *VarlogCluster) AddSN() (types.StorageNodeID, error) {
	snID := clus.snID
	clus.snID += types.StorageNodeID(1)

	datadir, err := ioutil.TempDir("", "test_*")
	if err != nil {
		return types.StorageNodeID(0), err
	}
	volume, err := storagenode.NewVolume(datadir)
	if err != nil {
		return types.StorageNodeID(0), err
	}

	sn, err := storagenode.New(context.TODO(),
		storagenode.WithListenAddress("127.0.0.1:0"),
		storagenode.WithClusterID(clus.ClusterID),
		storagenode.WithStorageNodeID(snID),
		storagenode.WithVolumes(volume),
	)
	if err != nil {
		return types.StorageNodeID(0), err
	}
	clus.snWG.Add(1)
	go func() {
		defer clus.snWG.Done()
		_ = sn.Run()
	}()

	meta, err := sn.GetMetadata(context.TODO())
	if err != nil {
		return types.StorageNodeID(0), err
	}

	clus.SNs[snID] = sn
	clus.volumes[snID] = volume
	clus.snRPCEndpoints[snID] = meta.StorageNode.Address

	err = clus.MRs[0].RegisterStorageNode(context.TODO(), meta.GetStorageNode())
	return snID, err
}

func (clus *VarlogCluster) AddSNByVMS() (types.StorageNodeID, error) {
	snID := clus.snID
	clus.snID += types.StorageNodeID(1)

	datadir, err := ioutil.TempDir("", "test_*")
	if err != nil {
		return types.StorageNodeID(0), err
	}
	volume, err := storagenode.NewVolume(datadir)
	if err != nil {
		return types.StorageNodeID(0), err
	}

	sn, err := storagenode.New(context.TODO(),
		storagenode.WithListenAddress("127.0.0.1:0"),
		storagenode.WithClusterID(clus.ClusterID),
		storagenode.WithStorageNodeID(snID),
		storagenode.WithVolumes(volume),
	)
	if err != nil {
		return types.StorageNodeID(0), err
	}
	clus.snWG.Add(1)
	go func() {
		defer clus.snWG.Done()
		_ = sn.Run()
	}()

	var meta *varlogpb.StorageNodeMetadataDescriptor
	meta, err = sn.GetMetadata(context.TODO())
	if err != nil {
		goto err_out
	}

	_, err = clus.CM.AddStorageNode(context.TODO(), meta.StorageNode.Address)
	if err != nil {
		goto err_out
	}

	clus.SNs[snID] = sn
	clus.volumes[snID] = volume
	clus.snRPCEndpoints[snID] = meta.StorageNode.Address
	return meta.StorageNode.StorageNodeID, nil

err_out:
	sn.Close()
	return types.StorageNodeID(0), err
}

func (clus *VarlogCluster) RecoverSN(snID types.StorageNodeID) (*storagenode.StorageNode, error) {
	volume, ok := clus.volumes[snID]
	if !ok {
		return nil, errors.New("no volume")
	}

	addr, _ := clus.snRPCEndpoints[snID]

	sn, err := storagenode.New(context.TODO(),
		storagenode.WithClusterID(clus.ClusterID),
		storagenode.WithStorageNodeID(snID),
		storagenode.WithListenAddress(addr),
		storagenode.WithVolumes(volume),
	)
	if err != nil {
		return nil, err
	}
	clus.snWG.Add(1)
	go func() {
		defer clus.snWG.Done()
		_ = sn.Run()
	}()

	clus.SNs[snID] = sn

	return sn, nil
}

func (clus *VarlogCluster) AddLS() (types.LogStreamID, error) {
	if len(clus.SNs) < clus.NrRep {
		return types.LogStreamID(0), verrors.ErrInvalid
	}

	lsID := clus.lsID
	clus.lsID += types.LogStreamID(1)

	snIDs := make([]types.StorageNodeID, 0, len(clus.SNs))
	for snID := range clus.SNs {
		snIDs = append(snIDs, snID)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(snIDs), func(i, j int) { snIDs[i], snIDs[j] = snIDs[j], snIDs[i] })

	replicas := make([]*varlogpb.ReplicaDescriptor, 0, clus.NrRep)
	for i := 0; i < clus.NrRep; i++ {
		storageNodeID := snIDs[i]
		storageNode := clus.SNs[storageNodeID]
		snmeta, err := storageNode.GetMetadata(context.TODO())
		if err != nil {
			return types.LogStreamID(0), err
		}
		targetpath := snmeta.GetStorageNode().GetStorages()[0].Path
		replicas = append(replicas, &varlogpb.ReplicaDescriptor{StorageNodeID: snIDs[i], Path: targetpath})
	}

	for _, r := range replicas {
		sn, _ := clus.SNs[r.StorageNodeID]
		if _, err := sn.AddLogStream(context.TODO(), lsID, r.Path); err != nil {
			return types.LogStreamID(0), err
		}
	}

	ls := &varlogpb.LogStreamDescriptor{
		LogStreamID: lsID,
		Replicas:    replicas,
	}

	if err := clus.MRs[0].RegisterLogStream(context.TODO(), ls); err != nil {
		return types.LogStreamID(0), err
	}

	// seal: assume that there is no vms
	for _, replica := range replicas {
		snid := replica.GetStorageNodeID()
		_, _, err := clus.SNs[snid].Seal(context.TODO(), lsID, types.InvalidGLSN)
		if err != nil {
			return types.LogStreamID(0), err
		}
	}

	// wait for sealed
	sealed := testutil.CompareWaitN(50, func() bool {
		for _, replica := range replicas {
			snid := replica.GetStorageNodeID()
			snmd, err := clus.SNs[snid].GetMetadata(context.TODO())
			if err != nil {
				return false
			}
			lsmd, ok := snmd.GetLogStream(lsID)
			if !ok {
				return false
			}
			if lsmd.GetStatus() != varlogpb.LogStreamStatusSealed {
				return false
			}
		}
		return true
	})
	if !sealed {
		return types.LogStreamID(0), errors.New("invalid status, expected=LogStreamStatusSealed")
	}

	// unseal
	for _, replica := range replicas {
		snid := replica.GetStorageNodeID()
		err := clus.SNs[snid].Unseal(context.TODO(), lsID)
		if err != nil {
			return types.LogStreamID(0), err
		}
	}

	running := testutil.CompareWaitN(50, func() bool {
		for _, replica := range replicas {
			snid := replica.GetStorageNodeID()
			snmd, err := clus.SNs[snid].GetMetadata(context.TODO())
			if err != nil {
				return false
			}
			lsmd, ok := snmd.GetLogStream(lsID)
			if !ok {
				return false
			}
			if !lsmd.GetStatus().Running() {
				return false
			}
		}
		return true
	})
	if !running {
		return types.LogStreamID(0), errors.New("invalid status, expected=LogStreamStatusRunning")
	}

	return lsID, nil
}

func (clus *VarlogCluster) AddLSByVMS() (logStreamID types.LogStreamID, err error) {
	if len(clus.SNs) < clus.NrRep {
		return types.LogStreamID(0), verrors.ErrInvalid
	}

	// FIXME: Use client api to add log stream since it is an integration test.
	logStreamDesc, err := clus.CM.AddLogStream(context.TODO(), nil)
	if err != nil {
		return types.LogStreamID(0), err
	}

	logStreamID = logStreamDesc.GetLogStreamID()
	clus.lsID = logStreamID + types.LogStreamID(1)

	// wait for sealed
	sealed := testutil.CompareWaitN(50, func() bool {
		time.Sleep(100 * time.Millisecond) // backoff

		for _, rd := range logStreamDesc.GetReplicas() {
			var snmd *varlogpb.StorageNodeMetadataDescriptor
			snid := rd.GetStorageNodeID()
			snmd, err = clus.SNs[snid].GetMetadata(context.TODO())
			if err != nil {
				return false
			}
			lsmd, ok := snmd.GetLogStream(logStreamID)
			if !ok {
				err = errors.Errorf("no such LogStreamID: %d", logStreamID)
				return false
			}
			if lsmd.GetStatus() != varlogpb.LogStreamStatusSealed {
				err = errors.Errorf("invalid status: expected=%v actual=%v",
					varlogpb.LogStreamStatusSealed,
					lsmd.GetStatus(),
				)
				return false
			}
		}

		var md *varlogpb.MetadataDescriptor
		md, err = clus.CM.Metadata(context.TODO())
		if err != nil {
			return false
		}

		var lsd *varlogpb.LogStreamDescriptor
		lsd, err = md.HaveLogStream(logStreamID)
		if err != nil {
			return false
		}

		if lsd.GetStatus() != varlogpb.LogStreamStatusSealed {
			err = errors.Errorf("invalid status: expected=%v actual=%v",
				varlogpb.LogStreamStatusSealed,
				lsd.GetStatus(),
			)
		}
		return err == nil

	})
	if !sealed {
		return types.LogStreamID(0), err
	}

	running := testutil.CompareWaitN(50, func() bool {
		time.Sleep(100 * time.Millisecond) // backoff

		// unseal
		_, err = clus.CMCli.Unseal(context.TODO(), logStreamID)
		if err != nil {
			return false
		}

		for _, rd := range logStreamDesc.GetReplicas() {
			var snmd *varlogpb.StorageNodeMetadataDescriptor
			snid := rd.GetStorageNodeID()
			snmd, err = clus.SNs[snid].GetMetadata(context.TODO())
			if err != nil {
				return false
			}

			lsmd, ok := snmd.GetLogStream(logStreamID)
			if !ok {
				err = errors.Errorf("no such LogStreamID: %d", logStreamID)
				return false
			}

			if !lsmd.GetStatus().Running() {
				err = errors.Errorf("invalid status: expected=%v actual=%v",
					varlogpb.LogStreamStatusRunning,
					lsmd.GetStatus(),
				)
				return false
			}
		}

		var md *varlogpb.MetadataDescriptor
		md, err = clus.CM.Metadata(context.TODO())
		if err != nil {
			return false

		}
		var lsd *varlogpb.LogStreamDescriptor
		lsd, err = md.HaveLogStream(logStreamID)
		if err != nil {
			return false
		}
		if lsd.GetStatus() != varlogpb.LogStreamStatusRunning {
			err = errors.Errorf("invalid status: expected=%v actual=%v",
				varlogpb.LogStreamStatusRunning,
				lsd.GetStatus(),
			)
		}

		return err == nil
	})
	if !running {
		return types.LogStreamID(0), err
	}

	clus.issuedLogStreamIDs = append(clus.issuedLogStreamIDs, logStreamID)
	return logStreamID, nil
}

func (clus *VarlogCluster) UpdateLS(lsID types.LogStreamID, oldsn, newsn types.StorageNodeID) error {
	sn := clus.LookupSN(newsn)
	_, err := sn.AddLogStream(context.TODO(), lsID, "path")
	if err != nil {
		return err
	}

	meta, err := clus.GetMR().GetMetadata(context.TODO())
	if err != nil {
		return err
	}

	oldLSDesc := meta.GetLogStream(lsID)
	if oldLSDesc == nil {
		return errors.New("logStream is not exist")
	}
	newLSDesc := proto.Clone(oldLSDesc).(*varlogpb.LogStreamDescriptor)

	exist := false
	for _, r := range newLSDesc.Replicas {
		if r.StorageNodeID == oldsn {
			r.StorageNodeID = newsn
			exist = true
		}
	}

	if !exist {
		return errors.New("invalid victim")
	}

	return clus.GetMR().UpdateLogStream(context.TODO(), newLSDesc)
}

func (clus *VarlogCluster) UpdateLSByVMS(lsID types.LogStreamID, oldsn, newsn types.StorageNodeID) error {
	sn := clus.LookupSN(newsn)

	snmeta, err := sn.GetMetadata(context.TODO())
	if err != nil {
		return err
	}
	path := snmeta.GetStorageNode().GetStorages()[0].GetPath()
	/*
		_, err = sn.AddLogStream(clus.ClusterID, newsn, lsID, path)
		if err != nil {
			return err
		}
	*/

	newReplica := &varlogpb.ReplicaDescriptor{
		StorageNodeID: newsn,
		Path:          path,
	}
	oldReplica := &varlogpb.ReplicaDescriptor{
		StorageNodeID: oldsn,
	}

	/*
		meta, err := clus.CM.Metadata(context.TODO())
		if err != nil {
			return err
		}

		oldLSDesc := meta.GetLogStream(lsID)
		if oldLSDesc == nil {
			return errors.New("logStream is not exist")
		}
		newLSDesc := proto.Clone(oldLSDesc).(*varlogpb.LogStreamDescriptor)

		exist := false
		for _, r := range newLSDesc.Replicas {
			if r.StorageNodeID == oldsn {
				r.StorageNodeID = newsn
				exist = true
			}
		}

		if !exist {
			return errors.New("invalid victim")
		}
	*/

	_, err = clus.CM.UpdateLogStream(context.TODO(), lsID, oldReplica, newReplica)
	return err
}

func (clus *VarlogCluster) LookupSN(snID types.StorageNodeID) *storagenode.StorageNode {
	sn, _ := clus.SNs[snID]
	return sn
}

func (clus *VarlogCluster) GetMR() *metadata_repository.RaftMetadataRepository {
	if len(clus.MRs) == 0 {
		return nil
	}

	return clus.MRs[0]
}

func (clus *VarlogCluster) LookupMR(nodeID types.NodeID) (*metadata_repository.RaftMetadataRepository, bool) {
	for idx, mrID := range clus.MRIDs {
		if nodeID == mrID {
			return clus.MRs[idx], true
		}
	}
	return nil, false
}

func (clus *VarlogCluster) GetVMS() vms.ClusterManager {
	return clus.CM
}

func (clus *VarlogCluster) getSN(lsID types.LogStreamID, idx int) (*storagenode.StorageNode, error) {
	if len(clus.MRs) == 0 {
		return nil, verrors.ErrInvalid
	}

	var meta *varlogpb.MetadataDescriptor
	var err error
	for _, mr := range clus.MRs {
		meta, err = mr.GetMetadata(context.TODO())
		if meta != nil {
			break
		}
	}

	if meta == nil {
		return nil, err
	}

	ls := meta.GetLogStream(lsID)
	if ls == nil {
		return nil, verrors.ErrNotExist
	}

	if len(ls.Replicas) < idx+1 {
		return nil, verrors.ErrInvalid
	}

	sn := clus.LookupSN(ls.Replicas[idx].StorageNodeID)
	if sn == nil {
		return nil, verrors.ErrInternal
	}

	return sn, nil
}

func (clus *VarlogCluster) GetPrimarySN(lsID types.LogStreamID) (*storagenode.StorageNode, error) {
	return clus.getSN(lsID, 0)
}

func (clus *VarlogCluster) GetBackupSN(lsID types.LogStreamID, idx int) (*storagenode.StorageNode, error) {
	return clus.getSN(lsID, idx)
}

func (clus *VarlogCluster) NewLogIOClient(lsID types.LogStreamID) (logc.LogIOClient, error) {
	sn, err := clus.GetPrimarySN(lsID)
	if err != nil {
		return nil, err
	}

	snMeta, err := sn.GetMetadata(context.TODO())
	if err != nil {
		return nil, err
	}

	return logc.NewLogIOClient(context.TODO(), snMeta.StorageNode.Address)
}

func (clus *VarlogCluster) RunClusterManager(mrAddrs []string, opts *vms.Options) (vms.ClusterManager, error) {
	if clus.VarlogClusterOptions.NrRep < 1 {
		return nil, verrors.ErrInvalidArgument
	}

	if opts == nil {
		vmOpts := vms.DefaultOptions()
		opts = &vmOpts
		opts.Logger = clus.logger
	}

	opts.ListenAddress = fmt.Sprintf("127.0.0.1:%d", clus.portLease.Base()+varlogClusterManagerPortBase)
	opts.ClusterID = clus.ClusterID
	opts.MetadataRepositoryAddresses = mrAddrs
	opts.ReplicationFactor = uint(clus.VarlogClusterOptions.NrRep)

	cm, err := vms.NewClusterManager(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	if err := cm.Run(); err != nil {
		return nil, err
	}
	clus.CM = cm
	return cm, nil
}

func (clus *VarlogCluster) NewClusterManagerClient() (varlog.ClusterManagerClient, error) {
	addr := clus.CM.Address()
	cmcli, err := varlog.NewClusterManagerClient(context.TODO(), addr)
	if err != nil {
		return nil, err
	}
	clus.CMCli = cmcli
	return cmcli, err
}

func (clus *VarlogCluster) GetClusterManagerClient() varlog.ClusterManagerClient {
	return clus.CMCli
}

func (clus *VarlogCluster) Logger() *zap.Logger {
	return clus.logger
}

func (clus *VarlogCluster) IssuedLogStreams() []types.LogStreamID {
	ret := make([]types.LogStreamID, len(clus.issuedLogStreamIDs))
	copy(ret, clus.issuedLogStreamIDs)
	return ret
}

func WithTestCluster(opts VarlogClusterOptions, f func(env *VarlogCluster)) func() {
	return func() {
		env := NewVarlogCluster(opts)
		env.Start()

		So(testutil.CompareWaitN(10, func() bool {
			return env.HealthCheck()
		}), ShouldBeTrue)

		mr := env.GetMR()
		So(testutil.CompareWaitN(10, func() bool {
			return mr.GetServerAddr() != ""
		}), ShouldBeTrue)
		mrAddr := mr.GetServerAddr()

		// VMS Server
		_, err := env.RunClusterManager([]string{mrAddr}, opts.VMSOpts)
		So(err, ShouldBeNil)

		// VMS Client
		cmCli, err := env.NewClusterManagerClient()
		So(err, ShouldBeNil)

		Reset(func() {
			env.Close()
			cmCli.Close()
			testutil.GC()
		})

		f(env)
	}
}
