package metadata_repository

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kakao/varlog/pkg/varlog"
	types "github.com/kakao/varlog/pkg/varlog/types"
	"github.com/kakao/varlog/pkg/varlog/util/testutil"
	pb "github.com/kakao/varlog/proto/metadata_repository"
	snpb "github.com/kakao/varlog/proto/storage_node"
	varlogpb "github.com/kakao/varlog/proto/varlog"
	"github.com/kakao/varlog/vtesting"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type metadataRepoCluster struct {
	nrRep             int
	peers             []string
	nodes             []*RaftMetadataRepository
	reporterClientFac ReporterClientFactory
	logger            *zap.Logger
}

var testSnapCount uint64

func newMetadataRepoCluster(n, nrRep int, increseUncommit bool) *metadataRepoCluster {
	peers := make([]string, n)
	nodes := make([]*RaftMetadataRepository, n)

	for i := range peers {
		peers[i] = fmt.Sprintf("http://127.0.0.1:%d", 10000+i)
	}

	clus := &metadataRepoCluster{
		nrRep:             nrRep,
		peers:             peers,
		nodes:             nodes,
		reporterClientFac: NewDummyReporterClientFactory(!increseUncommit),
		logger:            zap.L(),
	}

	for i := range clus.peers {
		clus.clear(i)
		clus.createMetadataRepo(i, false)
	}

	return clus
}

func (clus *metadataRepoCluster) clear(idx int) {
	if idx < 0 || idx >= len(clus.nodes) {
		return
	}

	url, _ := url.Parse(clus.peers[idx])
	nodeID := types.NewNodeID(url.Host)

	os.RemoveAll(fmt.Sprintf("raft-%d", nodeID))
	os.RemoveAll(fmt.Sprintf("raft-%d-snap", nodeID))

	return
}

func (clus *metadataRepoCluster) createMetadataRepo(idx int, join bool) error {
	if idx < 0 || idx >= len(clus.nodes) {
		return errors.New("out of range")
	}

	url, _ := url.Parse(clus.peers[idx])
	nodeID := types.NewNodeID(url.Host)

	options := &MetadataRepositoryOptions{
		ClusterID:         types.ClusterID(1),
		NodeID:            nodeID,
		Join:              join,
		SnapCount:         testSnapCount,
		RaftTick:          vtesting.TestRaftTick(),
		RPCTimeout:        vtesting.TimeoutAccordingToProcCnt(DefaultRPCTimeout),
		NumRep:            clus.nrRep,
		PeerList:          *cli.NewStringSlice(clus.peers...),
		RPCBindAddress:    ":0",
		ReporterClientFac: clus.reporterClientFac,
		Logger:            clus.logger,
	}

	clus.nodes[idx] = NewRaftMetadataRepository(options)
	return nil
}

func (clus *metadataRepoCluster) appendMetadataRepo() error {
	idx := len(clus.nodes)
	clus.peers = append(clus.peers, fmt.Sprintf("http://127.0.0.1:%d", 10000+idx))
	clus.nodes = append(clus.nodes, nil)

	clus.clear(idx)

	return clus.createMetadataRepo(idx, true)
}

func (clus *metadataRepoCluster) start(idx int) error {
	if idx < 0 || idx >= len(clus.nodes) {
		return errors.New("out of range")
	}

	clus.nodes[idx].Run()

	return nil
}

func (clus *metadataRepoCluster) Start() {
	clus.logger.Info("cluster start")
	for i := range clus.nodes {
		clus.start(i)
	}
	clus.logger.Info("cluster complete")
}

func (clus *metadataRepoCluster) stop(idx int) error {
	if idx < 0 || idx >= len(clus.nodes) {
		return errors.New("out or range")
	}

	return clus.nodes[idx].Close()
}

func (clus *metadataRepoCluster) restart(idx int) error {
	if idx < 0 || idx >= len(clus.nodes) {
		return errors.New("out of range")
	}

	clus.stop(idx)
	clus.createMetadataRepo(idx, false)
	return clus.start(idx)
}

func (clus *metadataRepoCluster) close(idx int) error {
	if idx < 0 || idx >= len(clus.nodes) {
		return errors.New("out or range")
	}

	err := clus.stop(idx)
	clus.clear(idx)

	clus.logger.Info("cluster node stop", zap.Int("idx", idx))

	return err
}

func (clus *metadataRepoCluster) hasSnapshot(idx int) (bool, error) {
	if idx < 0 || idx >= len(clus.nodes) {
		return false, errors.New("out or range")
	}

	return clus.nodes[idx].raftNode.loadSnapshot() != nil, nil
}

// Close closes all cluster nodes
func (clus *metadataRepoCluster) Close() error {
	var err error
	for i := range clus.peers {
		if erri := clus.close(i); erri != nil {
			err = erri
		}
	}
	return err
}

func (clus *metadataRepoCluster) leader() int {
	leader := -1
	for i, n := range clus.nodes {
		if n.isLeader() {
			leader = i
			break
		}
	}

	return leader
}

func (clus *metadataRepoCluster) leaderElected() bool {
	//return clus.leader() >= 0

	for _, n := range clus.nodes {
		if !n.hasLeader() {
			return false
		}
	}

	return true
}

func (clus *metadataRepoCluster) leaderFail() bool {
	leader := clus.leader()
	if leader < 0 {
		return false
	}

	clus.stop(leader)
	return true
}

func (clus *metadataRepoCluster) closeNoErrors(t *testing.T) {
	clus.logger.Info("cluster stop")
	if err := clus.Close(); err != nil {
		t.Log(err)
	}
	clus.logger.Info("cluster stop complete")
}

func makeLocalLogStream(snID types.StorageNodeID, knownHighWatermark types.GLSN, lsID types.LogStreamID, offset types.LLSN, length uint64) *snpb.LocalLogStreamDescriptor {
	lls := &snpb.LocalLogStreamDescriptor{
		StorageNodeID: snID,
		HighWatermark: knownHighWatermark,
	}
	ls := &snpb.LocalLogStreamDescriptor_LogStreamUncommitReport{
		LogStreamID:           lsID,
		UncommittedLLSNOffset: offset,
		UncommittedLLSNLength: length,
	}
	lls.Uncommit = append(lls.Uncommit, ls)

	return lls
}

func makeLogStream(lsID types.LogStreamID, snIDs []types.StorageNodeID) *varlogpb.LogStreamDescriptor {
	ls := &varlogpb.LogStreamDescriptor{
		LogStreamID: lsID,
		Status:      varlogpb.LogStreamStatusRunning,
	}

	for _, snID := range snIDs {
		r := &varlogpb.ReplicaDescriptor{
			StorageNodeID: snID,
		}

		ls.Replicas = append(ls.Replicas, r)
	}

	return ls
}

func TestMRApplyReport(t *testing.T) {

	Convey("Report Should not be applied if not register LogStream", t, func(ctx C) {
		rep := 2
		clus := newMetadataRepoCluster(1, rep, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		mr := clus.nodes[0]

		snIds := make([]types.StorageNodeID, rep)
		for i := range snIds {
			snIds[i] = types.StorageNodeID(i)

			sn := &varlogpb.StorageNodeDescriptor{
				StorageNodeID: snIds[i],
			}

			err := mr.storage.registerStorageNode(sn)
			So(err, ShouldBeNil)
		}
		lsId := types.LogStreamID(0)
		notExistSnID := types.StorageNodeID(rep)

		lls := makeLocalLogStream(snIds[0], types.InvalidGLSN, lsId, types.MinLLSN, 2)
		mr.applyReport(&pb.Report{LogStream: lls})

		for _, snId := range snIds {
			r := mr.storage.LookupLocalLogStreamReplica(lsId, snId)
			So(r, ShouldBeNil)
		}

		Convey("LocalLogStream should register when register LogStream", func(ctx C) {
			ls := makeLogStream(lsId, snIds)
			err := mr.storage.registerLogStream(ls)
			So(err, ShouldBeNil)

			for _, snId := range snIds {
				r := mr.storage.LookupLocalLogStreamReplica(lsId, snId)
				So(r, ShouldNotBeNil)
			}

			Convey("Report should not apply if snID is not exist in LocalLogStream", func(ctx C) {
				lls := makeLocalLogStream(notExistSnID, types.InvalidGLSN, lsId, types.MinLLSN, 2)
				mr.applyReport(&pb.Report{LogStream: lls})

				r := mr.storage.LookupLocalLogStreamReplica(lsId, notExistSnID)
				So(r, ShouldBeNil)
			})

			Convey("Report should apply if snID is exist in LocalLogStream", func(ctx C) {
				snId := snIds[0]
				lls := makeLocalLogStream(snId, types.InvalidGLSN, lsId, types.MinLLSN, 2)
				mr.applyReport(&pb.Report{LogStream: lls})

				r := mr.storage.LookupLocalLogStreamReplica(lsId, snId)
				So(r, ShouldNotBeNil)
				So(r.UncommittedLLSNEnd(), ShouldEqual, types.MinLLSN+types.LLSN(2))

				Convey("Report which have bigger END LLSN Should be applied", func(ctx C) {
					lls := makeLocalLogStream(snId, types.InvalidGLSN, lsId, types.MinLLSN, 3)
					mr.applyReport(&pb.Report{LogStream: lls})

					r := mr.storage.LookupLocalLogStreamReplica(lsId, snId)
					So(r, ShouldNotBeNil)
					So(r.UncommittedLLSNEnd(), ShouldEqual, types.MinLLSN+types.LLSN(3))
				})

				Convey("Report which have smaller END LLSN Should Not be applied", func(ctx C) {
					lls := makeLocalLogStream(snId, types.InvalidGLSN, lsId, types.MinLLSN, 1)
					mr.applyReport(&pb.Report{LogStream: lls})

					r := mr.storage.LookupLocalLogStreamReplica(lsId, snId)
					So(r, ShouldNotBeNil)
					So(r.UncommittedLLSNEnd(), ShouldNotEqual, types.MinLLSN+types.LLSN(1))
				})
			})
		})
	})
}

func TestMRCalculateCommit(t *testing.T) {
	Convey("Calculate commit", t, func(ctx C) {
		clus := newMetadataRepoCluster(1, 2, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		mr := clus.nodes[0]

		snIds := make([]types.StorageNodeID, 2)
		for i := range snIds {
			snIds[i] = types.StorageNodeID(i)
			sn := &varlogpb.StorageNodeDescriptor{
				StorageNodeID: snIds[i],
			}

			err := mr.storage.registerStorageNode(sn)
			So(err, ShouldBeNil)
		}
		lsId := types.LogStreamID(0)
		ls := makeLogStream(lsId, snIds)
		err := mr.storage.registerLogStream(ls)
		So(err, ShouldBeNil)

		Convey("LogStream which all reports have not arrived cannot be commit", func(ctx C) {
			lls := makeLocalLogStream(snIds[0], types.InvalidGLSN, lsId, types.MinLLSN, 2)
			mr.applyReport(&pb.Report{LogStream: lls})

			replicas := mr.storage.LookupLocalLogStream(lsId)
			_, minHWM, nrCommit := mr.calculateCommit(replicas)
			So(nrCommit, ShouldEqual, 0)
			So(minHWM, ShouldEqual, types.InvalidGLSN)
		})

		Convey("LogStream which all reports are disjoint cannot be commit", func(ctx C) {
			lls := makeLocalLogStream(snIds[0], types.GLSN(10), lsId, types.MinLLSN+types.LLSN(5), 1)
			mr.applyReport(&pb.Report{LogStream: lls})

			lls = makeLocalLogStream(snIds[1], types.GLSN(7), lsId, types.MinLLSN+types.LLSN(3), 2)
			mr.applyReport(&pb.Report{LogStream: lls})

			replicas := mr.storage.LookupLocalLogStream(lsId)
			knownHWM, minHWM, nrCommit := mr.calculateCommit(replicas)
			So(nrCommit, ShouldEqual, 0)
			So(knownHWM, ShouldEqual, types.GLSN(10))
			So(minHWM, ShouldEqual, types.GLSN(7))
		})

		Convey("LogStream Should be commit where replication is completed", func(ctx C) {
			lls := makeLocalLogStream(snIds[0], types.GLSN(10), lsId, types.MinLLSN+types.LLSN(3), 3)
			mr.applyReport(&pb.Report{LogStream: lls})

			lls = makeLocalLogStream(snIds[1], types.GLSN(9), lsId, types.MinLLSN+types.LLSN(3), 2)
			mr.applyReport(&pb.Report{LogStream: lls})

			replicas := mr.storage.LookupLocalLogStream(lsId)
			knownHWM, minHWM, nrCommit := mr.calculateCommit(replicas)
			So(nrCommit, ShouldEqual, 2)
			So(minHWM, ShouldEqual, types.GLSN(9))
			So(knownHWM, ShouldEqual, types.GLSN(10))
		})
	})
}

func TestMRGlobalCommit(t *testing.T) {
	Convey("Calculate commit", t, func(ctx C) {
		rep := 2
		clus := newMetadataRepoCluster(1, rep, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})

		clus.Start()
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		mr := clus.nodes[0]

		snIds := make([][]types.StorageNodeID, 2)
		for i := range snIds {
			snIds[i] = make([]types.StorageNodeID, rep)
			for j := range snIds[i] {
				snIds[i][j] = types.StorageNodeID(i*2 + j)

				sn := &varlogpb.StorageNodeDescriptor{
					StorageNodeID: snIds[i][j],
				}

				err := mr.storage.registerStorageNode(sn)
				So(err, ShouldBeNil)
			}
		}

		lsIds := make([]types.LogStreamID, 2)
		for i := range lsIds {
			lsIds[i] = types.LogStreamID(i)
		}

		for i, lsId := range lsIds {
			ls := makeLogStream(lsId, snIds[i])
			err := mr.storage.registerLogStream(ls)
			So(err, ShouldBeNil)
		}

		Convey("global commit", func(ctx C) {
			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[0][0], types.InvalidGLSN, lsIds[0], types.MinLLSN, 2)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[0][1], types.InvalidGLSN, lsIds[0], types.MinLLSN, 2)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[1][0], types.InvalidGLSN, lsIds[1], types.MinLLSN, 4)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[1][1], types.InvalidGLSN, lsIds[1], types.MinLLSN, 3)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			// global commit (2, 3) highest glsn: 5
			So(testutil.CompareWaitN(10, func() bool {
				return mr.storage.GetHighWatermark() == types.GLSN(5)
			}), ShouldBeTrue)

			Convey("LogStream should be dedup", func(ctx C) {
				So(testutil.CompareWaitN(10, func() bool {
					lls := makeLocalLogStream(snIds[0][0], types.InvalidGLSN, lsIds[0], types.MinLLSN, 3)
					return mr.proposeReport(lls) == nil
				}), ShouldBeTrue)

				So(testutil.CompareWaitN(10, func() bool {
					lls := makeLocalLogStream(snIds[0][1], types.InvalidGLSN, lsIds[0], types.MinLLSN, 2)
					return mr.proposeReport(lls) == nil
				}), ShouldBeTrue)

				time.Sleep(vtesting.TimeoutUnitTimesFactor(1))

				So(testutil.CompareWaitN(50, func() bool {
					return mr.storage.GetHighWatermark() == types.GLSN(5)
				}), ShouldBeTrue)
			})

			Convey("LogStream which have wrong GLSN but have uncommitted should commit", func(ctx C) {
				So(testutil.CompareWaitN(10, func() bool {
					lls := makeLocalLogStream(snIds[0][0], types.InvalidGLSN, lsIds[0], types.MinLLSN, 6)
					return mr.proposeReport(lls) == nil
				}), ShouldBeTrue)

				So(testutil.CompareWaitN(10, func() bool {
					lls := makeLocalLogStream(snIds[0][1], types.InvalidGLSN, lsIds[0], types.MinLLSN, 6)
					return mr.proposeReport(lls) == nil
				}), ShouldBeTrue)

				So(testutil.CompareWaitN(10, func() bool {
					return mr.storage.GetHighWatermark() == types.GLSN(9)
				}), ShouldBeTrue)
			})
		})
	})
}

func TestMRSimpleReportNCommit(t *testing.T) {
	Convey("Uncommitted LocalLogStream should be committed", t, func(ctx C) {
		clus := newMetadataRepoCluster(1, 1, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		clus.Start()
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		snID := types.StorageNodeID(0)
		snIDs := make([]types.StorageNodeID, 1)
		snIDs = append(snIDs, snID)

		lsID := types.LogStreamID(snID)

		sn := &varlogpb.StorageNodeDescriptor{
			StorageNodeID: snID,
		}

		rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
		defer cancel()

		err := clus.nodes[0].RegisterStorageNode(rctx, sn)
		So(err, ShouldBeNil)

		So(testutil.CompareWaitN(50, func() bool {
			return clus.reporterClientFac.(*DummyReporterClientFactory).lookupClient(snID) != nil
		}), ShouldBeTrue)

		ls := makeLogStream(lsID, snIDs)
		rctx, cancel = context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
		defer cancel()
		err = clus.nodes[0].RegisterLogStream(rctx, ls)
		So(err, ShouldBeNil)

		reporterClient := clus.reporterClientFac.(*DummyReporterClientFactory).lookupClient(snID)
		reporterClient.increaseUncommitted()

		So(testutil.CompareWaitN(50, func() bool {
			return reporterClient.numUncommitted() == 0
		}), ShouldBeTrue)
	})
}

func TestMRRequestMap(t *testing.T) {
	Convey("requestMap should have request when wait ack", t, func(ctx C) {
		clus := newMetadataRepoCluster(1, 1, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		mr := clus.nodes[0]

		sn := &varlogpb.StorageNodeDescriptor{
			StorageNodeID: types.StorageNodeID(0),
		}

		requestNum := atomic.LoadUint64(&mr.requestNum)

		var wg sync.WaitGroup
		var st sync.WaitGroup

		st.Add(1)
		wg.Add(1)
		go func() {
			defer wg.Done()
			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(1))
			defer cancel()
			st.Done()
			mr.RegisterStorageNode(rctx, sn)
		}()

		st.Wait()
		So(testutil.CompareWaitN(1, func() bool {
			_, ok := mr.requestMap.Load(requestNum + 1)
			return ok
		}), ShouldBeTrue)

		wg.Wait()
	})

	Convey("requestMap should ignore request that have different nodeIndex", t, func(ctx C) {
		clus := newMetadataRepoCluster(1, 1, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		mr := clus.nodes[0]

		sn := &varlogpb.StorageNodeDescriptor{
			StorageNodeID: types.StorageNodeID(0),
		}

		var st sync.WaitGroup
		var wg sync.WaitGroup
		st.Add(1)
		wg.Add(1)
		go func() {
			defer wg.Done()
			st.Done()

			testutil.CompareWaitN(50, func() bool {
				_, ok := mr.requestMap.Load(uint64(1))
				return ok
			})

			dummy := &committedEntry{
				entry: &pb.RaftEntry{
					NodeIndex:    2,
					RequestIndex: uint64(1),
				},
			}
			mr.commitC <- dummy
		}()

		st.Wait()
		rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(2))
		defer cancel()
		err := mr.RegisterStorageNode(rctx, sn)

		wg.Wait()
		So(err, ShouldNotBeNil)
	})

	Convey("requestMap should delete request when context timeout", t, func(ctx C) {
		clus := newMetadataRepoCluster(1, 1, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		mr := clus.nodes[0]

		sn := &varlogpb.StorageNodeDescriptor{
			StorageNodeID: types.StorageNodeID(0),
		}

		rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(1))
		defer cancel()

		requestNum := atomic.LoadUint64(&mr.requestNum)
		err := mr.RegisterStorageNode(rctx, sn)
		So(err, ShouldNotBeNil)

		_, ok := mr.requestMap.Load(requestNum + 1)
		So(ok, ShouldBeFalse)
	})

	Convey("requestMap should delete after ack", t, func(ctx C) {
		clus := newMetadataRepoCluster(1, 1, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		clus.Start()
		So(testutil.CompareWaitN(50, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		mr := clus.nodes[0]

		So(testutil.CompareWaitN(50, func() bool {
			return mr.storage.LookupEndpoint(mr.nodeID) != ""
		}), ShouldBeTrue)

		sn := &varlogpb.StorageNodeDescriptor{
			StorageNodeID: types.StorageNodeID(0),
		}

		requestNum := atomic.LoadUint64(&mr.requestNum)
		err := mr.RegisterStorageNode(context.TODO(), sn)
		So(err, ShouldBeNil)

		_, ok := mr.requestMap.Load(requestNum + 1)
		So(ok, ShouldBeFalse)
	})
}

func TestMRGetLastCommitted(t *testing.T) {
	Convey("getLastCommitted", t, func(ctx C) {
		rep := 2
		clus := newMetadataRepoCluster(1, rep, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		clus.Start()
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		mr := clus.nodes[0]

		snIds := make([][]types.StorageNodeID, 2)
		for i := range snIds {
			snIds[i] = make([]types.StorageNodeID, rep)
			for j := range snIds[i] {
				snIds[i][j] = types.StorageNodeID(i*2 + j)

				sn := &varlogpb.StorageNodeDescriptor{
					StorageNodeID: snIds[i][j],
				}

				err := mr.storage.registerStorageNode(sn)
				So(err, ShouldBeNil)
			}
		}

		lsIds := make([]types.LogStreamID, 2)
		for i := range lsIds {
			lsIds[i] = types.LogStreamID(i)
		}

		for i, lsId := range lsIds {
			ls := makeLogStream(lsId, snIds[i])
			err := mr.storage.registerLogStream(ls)
			So(err, ShouldBeNil)
		}

		Convey("getLastCommitted should return last committed GLSN", func(ctx C) {
			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[0][0], types.InvalidGLSN, lsIds[0], types.MinLLSN, 2)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[0][1], types.InvalidGLSN, lsIds[0], types.MinLLSN, 2)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[1][0], types.InvalidGLSN, lsIds[1], types.MinLLSN, 4)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[1][1], types.InvalidGLSN, lsIds[1], types.MinLLSN, 3)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			// global commit (2, 3) highest glsn: 5
			So(testutil.CompareWaitN(10, func() bool {
				return mr.storage.GetHighWatermark() == types.GLSN(5)
			}), ShouldBeTrue)

			So(mr.getLastCommitted(lsIds[0]), ShouldEqual, types.GLSN(2))
			So(mr.getLastCommitted(lsIds[1]), ShouldEqual, types.GLSN(5))

			Convey("getLastCommitted should return same if not committed", func(ctx C) {
				for i := 0; i < 10; i++ {
					So(testutil.CompareWaitN(10, func() bool {
						lls := makeLocalLogStream(snIds[1][0], types.InvalidGLSN, lsIds[1], types.MinLLSN, uint64(4+i))
						return mr.proposeReport(lls) == nil
					}), ShouldBeTrue)

					So(testutil.CompareWaitN(10, func() bool {
						lls := makeLocalLogStream(snIds[1][1], types.InvalidGLSN, lsIds[1], types.MinLLSN, uint64(4+i))
						return mr.proposeReport(lls) == nil
					}), ShouldBeTrue)

					So(testutil.CompareWaitN(50, func() bool {
						return mr.storage.GetHighWatermark() == types.GLSN(6+i)
					}), ShouldBeTrue)

					So(mr.getLastCommitted(lsIds[0]), ShouldEqual, types.GLSN(2))
					So(mr.getLastCommitted(lsIds[1]), ShouldEqual, types.GLSN(6+i))
				}
			})

			Convey("getLastCommitted should return same for sealed LS", func(ctx C) {
				rctx, _ := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
				_, err := mr.Seal(rctx, lsIds[1])
				So(err, ShouldBeNil)

				for i := 0; i < 10; i++ {
					So(testutil.CompareWaitN(10, func() bool {
						lls := makeLocalLogStream(snIds[0][0], types.InvalidGLSN, lsIds[0], types.MinLLSN, uint64(3+i))
						return mr.proposeReport(lls) == nil
					}), ShouldBeTrue)

					So(testutil.CompareWaitN(10, func() bool {
						lls := makeLocalLogStream(snIds[0][1], types.InvalidGLSN, lsIds[0], types.MinLLSN, uint64(3+i))
						return mr.proposeReport(lls) == nil
					}), ShouldBeTrue)

					So(testutil.CompareWaitN(10, func() bool {
						lls := makeLocalLogStream(snIds[1][0], types.InvalidGLSN, lsIds[1], types.MinLLSN, uint64(4+i))
						return mr.proposeReport(lls) == nil
					}), ShouldBeTrue)

					So(testutil.CompareWaitN(10, func() bool {
						lls := makeLocalLogStream(snIds[1][1], types.InvalidGLSN, lsIds[1], types.MinLLSN, uint64(4+i))
						return mr.proposeReport(lls) == nil
					}), ShouldBeTrue)

					So(testutil.CompareWaitN(10, func() bool {
						return mr.storage.GetHighWatermark() == types.GLSN(6+i)
					}), ShouldBeTrue)

					So(mr.getLastCommitted(lsIds[0]), ShouldEqual, types.GLSN(6+i))
					So(mr.getLastCommitted(lsIds[1]), ShouldEqual, types.GLSN(5))
				}
			})
		})
	})
}

func TestMRSeal(t *testing.T) {
	Convey("seal", t, func(ctx C) {
		rep := 2
		clus := newMetadataRepoCluster(1, rep, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		clus.Start()
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		mr := clus.nodes[0]

		snIds := make([][]types.StorageNodeID, 2)
		for i := range snIds {
			snIds[i] = make([]types.StorageNodeID, rep)
			for j := range snIds[i] {
				snIds[i][j] = types.StorageNodeID(i*2 + j)

				sn := &varlogpb.StorageNodeDescriptor{
					StorageNodeID: snIds[i][j],
				}

				err := mr.storage.registerStorageNode(sn)
				So(err, ShouldBeNil)
			}
		}

		lsIds := make([]types.LogStreamID, 2)
		for i := range lsIds {
			lsIds[i] = types.LogStreamID(i)
		}

		for i, lsId := range lsIds {
			ls := makeLogStream(lsId, snIds[i])
			err := mr.storage.registerLogStream(ls)
			So(err, ShouldBeNil)
		}

		Convey("Seal should commit and return last committed", func(ctx C) {
			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[0][0], types.InvalidGLSN, lsIds[0], types.MinLLSN, 2)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[0][1], types.InvalidGLSN, lsIds[0], types.MinLLSN, 2)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[1][0], types.InvalidGLSN, lsIds[1], types.MinLLSN, 4)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			So(testutil.CompareWaitN(10, func() bool {
				lls := makeLocalLogStream(snIds[1][1], types.InvalidGLSN, lsIds[1], types.MinLLSN, 3)
				return mr.proposeReport(lls) == nil
			}), ShouldBeTrue)

			rctx, _ := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			lc, err := mr.Seal(rctx, lsIds[1])
			So(err, ShouldBeNil)
			So(lc, ShouldEqual, types.GLSN(5))

			Convey("Seal should return same last committed", func(ctx C) {
				for i := 0; i < 10; i++ {
					rctx, _ := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
					lc, err := mr.Seal(rctx, lsIds[1])
					So(err, ShouldBeNil)
					So(lc, ShouldEqual, types.GLSN(5))
				}
			})
		})
	})
}

func TestMRUnseal(t *testing.T) {
	Convey("unseal", t, func(ctx C) {
		rep := 2
		clus := newMetadataRepoCluster(1, rep, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		clus.Start()
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		mr := clus.nodes[0]

		snIds := make([][]types.StorageNodeID, 2)
		for i := range snIds {
			snIds[i] = make([]types.StorageNodeID, rep)
			for j := range snIds[i] {
				snIds[i][j] = types.StorageNodeID(i*2 + j)

				sn := &varlogpb.StorageNodeDescriptor{
					StorageNodeID: snIds[i][j],
				}

				err := mr.storage.registerStorageNode(sn)
				So(err, ShouldBeNil)
			}
		}

		lsIds := make([]types.LogStreamID, 2)
		for i := range lsIds {
			lsIds[i] = types.LogStreamID(i)
		}

		for i, lsId := range lsIds {
			ls := makeLogStream(lsId, snIds[i])
			err := mr.storage.registerLogStream(ls)
			So(err, ShouldBeNil)
		}

		So(testutil.CompareWaitN(10, func() bool {
			lls := makeLocalLogStream(snIds[0][0], types.InvalidGLSN, lsIds[0], types.MinLLSN, 2)
			return mr.proposeReport(lls) == nil
		}), ShouldBeTrue)

		So(testutil.CompareWaitN(10, func() bool {
			lls := makeLocalLogStream(snIds[0][1], types.InvalidGLSN, lsIds[0], types.MinLLSN, 2)
			return mr.proposeReport(lls) == nil
		}), ShouldBeTrue)

		So(testutil.CompareWaitN(10, func() bool {
			lls := makeLocalLogStream(snIds[1][0], types.InvalidGLSN, lsIds[1], types.MinLLSN, 4)
			return mr.proposeReport(lls) == nil
		}), ShouldBeTrue)

		So(testutil.CompareWaitN(10, func() bool {
			lls := makeLocalLogStream(snIds[1][1], types.InvalidGLSN, lsIds[1], types.MinLLSN, 3)
			return mr.proposeReport(lls) == nil
		}), ShouldBeTrue)

		rctx, _ := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
		lc, err := mr.Seal(rctx, lsIds[1])
		So(err, ShouldBeNil)
		So(lc, ShouldEqual, types.GLSN(5))

		Convey("Unealed LS should update report", func(ctx C) {
			rctx, _ := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			err := mr.Unseal(rctx, lsIds[1])
			So(err, ShouldBeNil)

			for i := 0; i < 10; i++ {
				So(testutil.CompareWaitN(10, func() bool {
					lls := makeLocalLogStream(snIds[1][0], types.InvalidGLSN, lsIds[1], types.MinLLSN, uint64(4+i))
					return mr.proposeReport(lls) == nil
				}), ShouldBeTrue)

				So(testutil.CompareWaitN(10, func() bool {
					lls := makeLocalLogStream(snIds[1][1], types.InvalidGLSN, lsIds[1], types.MinLLSN, uint64(4+i))
					return mr.proposeReport(lls) == nil
				}), ShouldBeTrue)

				So(testutil.CompareWaitN(50, func() bool {
					return mr.storage.GetHighWatermark() == types.GLSN(6+i)
				}), ShouldBeTrue)

				So(mr.getLastCommitted(lsIds[1]), ShouldEqual, types.GLSN(6+i))
			}
		})
	})
}

func TestMRFailoverLeaderElection(t *testing.T) {
	Convey("Given MR cluster", t, func(ctx C) {
		nrRep := 1
		nrNode := 3

		clus := newMetadataRepoCluster(nrNode, nrRep, true)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		clus.Start()
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		snIDs := make([]types.StorageNodeID, nrRep)
		for i := range snIDs {
			snIDs[i] = types.StorageNodeID(i)

			sn := &varlogpb.StorageNodeDescriptor{
				StorageNodeID: snIDs[i],
			}

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			err := clus.nodes[0].RegisterStorageNode(rctx, sn)
			So(err, ShouldBeNil)
		}

		lsID := types.LogStreamID(0)

		ls := makeLogStream(lsID, snIDs)
		rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
		defer cancel()
		err := clus.nodes[0].RegisterLogStream(rctx, ls)
		So(err, ShouldBeNil)

		reporterClient := clus.reporterClientFac.(*DummyReporterClientFactory).lookupClient(snIDs[0])
		So(testutil.CompareWaitN(50, func() bool {
			return !reporterClient.getKnownHighWatermark().Invalid()
		}), ShouldBeTrue)

		Convey("When node fail", func(ctx C) {
			leader := clus.leader()

			So(clus.leaderFail(), ShouldBeTrue)

			Convey("Then MR Cluster should elect", func(ctx C) {
				So(testutil.CompareWaitN(50, func() bool {
					return leader != clus.leader()
				}), ShouldBeTrue)

				prev := reporterClient.getKnownHighWatermark()

				So(testutil.CompareWaitN(50, func() bool {
					return reporterClient.getKnownHighWatermark() > prev
				}), ShouldBeTrue)
			})
		})
	})
}

func TestMRFailoverJoinNewNode(t *testing.T) {
	Convey("Given MR cluster", t, func(ctx C) {
		nrRep := 1
		nrNode := 3

		clus := newMetadataRepoCluster(nrNode, nrRep, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		clus.Start()
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		snIDs := make([]types.StorageNodeID, nrRep)
		for i := range snIDs {
			snIDs[i] = types.StorageNodeID(i)

			sn := &varlogpb.StorageNodeDescriptor{
				StorageNodeID: snIDs[i],
			}

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			err := clus.nodes[0].RegisterStorageNode(rctx, sn)
			So(err, ShouldBeNil)
		}

		lsID := types.LogStreamID(0)

		ls := makeLogStream(lsID, snIDs)
		rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
		defer cancel()
		err := clus.nodes[0].RegisterLogStream(rctx, ls)
		So(err, ShouldBeNil)

		Convey("When new node join", func(ctx C) {
			newNode := nrNode
			So(clus.appendMetadataRepo(), ShouldBeNil)

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			So(clus.nodes[0].AddPeer(rctx,
				types.ClusterID(0),
				clus.nodes[newNode].nodeID,
				clus.peers[newNode]), ShouldBeNil)
			So(clus.start(newNode), ShouldBeNil)
			nrNode += 1

			Convey("Then getMeta from new node should be success", func(ctx C) {
				So(testutil.CompareWaitN(50, func() bool {
					cinfo, err := clus.nodes[newNode].GetClusterInfo(context.TODO(), types.ClusterID(0))
					if err != nil {
						return false
					}

					if len(cinfo.Members) != nrNode {
						return false
					}

					for _, member := range cinfo.Members {
						if member.Endpoint == "" {
							return false
						}
					}

					return true
				}), ShouldBeTrue)

				Convey("Register to new node should be success", func(ctx C) {
					snID := snIDs[nrRep-1] + types.StorageNodeID(1)

					sn := &varlogpb.StorageNodeDescriptor{
						StorageNodeID: snID,
					}

					rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
					defer cancel()

					err := clus.nodes[newNode].RegisterStorageNode(rctx, sn)
					if rctx.Err() != nil {
						clus.logger.Info("complete with ctx error", zap.String("err", rctx.Err().Error()))
					}
					So(err, ShouldBeNil)

					meta, err := clus.nodes[newNode].GetMetadata(context.TODO())
					So(err, ShouldBeNil)
					So(meta.GetStorageNode(snID), ShouldNotBeNil)
				})
			})
		})

		Convey("When new nodes joining", func(ctx C) {
			newNode := nrNode
			So(clus.appendMetadataRepo(), ShouldBeNil)

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			So(clus.nodes[0].AddPeer(rctx,
				types.ClusterID(0),
				clus.nodes[newNode].nodeID,
				clus.peers[newNode]), ShouldBeNil)
			So(clus.start(newNode), ShouldBeNil)
			nrNode += 1

			Convey("Then proposal should be operated", func(ctx C) {
				snID := snIDs[nrRep-1] + types.StorageNodeID(1)

				sn := &varlogpb.StorageNodeDescriptor{
					StorageNodeID: snID,
				}

				rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
				defer cancel()

				err := clus.nodes[newNode].RegisterStorageNode(rctx, sn)
				if rctx.Err() != nil {
					clus.logger.Info("complete with ctx error", zap.String("err", rctx.Err().Error()))
				}
				So(err, ShouldBeNil)
			})
		})

		Convey("When new nodes started but not yet joined", func(ctx C) {
			newNode := nrNode
			So(clus.appendMetadataRepo(), ShouldBeNil)
			So(clus.start(newNode), ShouldBeNil)
			nrNode += 1

			time.Sleep(10 * time.Second)

			Convey("Then it should not have member info", func(ctx C) {
				cinfo, err := clus.nodes[newNode].GetClusterInfo(context.TODO(), types.ClusterID(0))
				So(err, ShouldResemble, varlog.ErrNotMember)

				cinfo, err = clus.nodes[0].GetClusterInfo(context.TODO(), types.ClusterID(0))
				So(err, ShouldBeNil)
				So(len(cinfo.Members), ShouldBeLessThan, nrNode)

				Convey("After joining, it should have member info", func(ctx C) {
					rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
					defer cancel()

					So(clus.nodes[0].AddPeer(rctx,
						types.ClusterID(0),
						clus.nodes[newNode].nodeID,
						clus.peers[newNode]), ShouldBeNil)

					So(testutil.CompareWaitN(50, func() bool {
						cinfo, _ := clus.nodes[newNode].GetClusterInfo(context.TODO(), types.ClusterID(0))
						return len(cinfo.GetMembers()) == nrNode
					}), ShouldBeTrue)

					Convey("Then proposal should be operated", func(ctx C) {
						snID := snIDs[nrRep-1] + types.StorageNodeID(1)

						sn := &varlogpb.StorageNodeDescriptor{
							StorageNodeID: snID,
						}

						rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(80))
						defer cancel()

						err := clus.nodes[newNode].RegisterStorageNode(rctx, sn)
						if rctx.Err() != nil {
							clus.logger.Info("complete with ctx error", zap.String("err", rctx.Err().Error()))
						}
						So(err, ShouldBeNil)
					})
				})
			})
		})
	})
}

func TestMRFailoverLeaveNode(t *testing.T) {
	Convey("Given MR cluster", t, func(ctx C) {
		nrRep := 1
		nrNode := 3

		clus := newMetadataRepoCluster(nrNode, nrRep, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		clus.Start()
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		leader := clus.leader()
		So(leader, ShouldBeGreaterThan, -1)

		Convey("When follower leave", func(ctx C) {
			leaveNode := (leader + 1) % nrNode
			checkNode := (leaveNode + 1) % nrNode

			So(testutil.CompareWaitN(50, func() bool {
				cinfo, _ := clus.nodes[checkNode].GetClusterInfo(context.TODO(), 0)
				return len(cinfo.Members) == nrNode
			}), ShouldBeTrue)

			clus.stop(leaveNode)

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			So(clus.nodes[checkNode].RemovePeer(rctx,
				types.ClusterID(0),
				clus.nodes[leaveNode].nodeID), ShouldBeNil)

			Convey("Then GetMembership should return 2 peers", func(ctx C) {
				cinfo, err := clus.nodes[checkNode].GetClusterInfo(context.TODO(), 0)
				So(err, ShouldBeNil)
				So(len(cinfo.Members), ShouldEqual, nrNode-1)
			})
		})

		Convey("When leader leave", func(ctx C) {
			leaveNode := leader
			checkNode := (leaveNode + 1) % nrNode

			So(testutil.CompareWaitN(50, func() bool {
				cinfo, _ := clus.nodes[checkNode].GetClusterInfo(context.TODO(), 0)
				return len(cinfo.Members) == nrNode
			}), ShouldBeTrue)

			clus.stop(leaveNode)

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			So(clus.nodes[checkNode].RemovePeer(rctx,
				types.ClusterID(0),
				clus.nodes[leaveNode].nodeID), ShouldBeNil)

			Convey("Then GetMembership should return 2 peers", func(ctx C) {
				cinfo, err := clus.nodes[checkNode].GetClusterInfo(context.TODO(), 0)
				So(err, ShouldBeNil)
				So(len(cinfo.Members), ShouldEqual, nrNode-1)
			})
		})
	})
}

func TestMRFailoverRestart(t *testing.T) {
	Convey("Given MR cluster with 5 peers", t, func(ctx C) {
		nrRep := 1
		nrNode := 5

		clus := newMetadataRepoCluster(nrNode, nrRep, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		clus.Start()
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		leader := clus.leader()
		So(leader, ShouldBeGreaterThan, -1)

		Convey("When follower restart and new node joined", func(ctx C) {
			restartNode := (leader + 1) % nrNode
			clus.restart(restartNode)

			newNode := nrNode
			So(clus.appendMetadataRepo(), ShouldBeNil)

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			So(clus.nodes[leader].AddPeer(rctx,
				types.ClusterID(0),
				clus.nodes[newNode].nodeID,
				clus.peers[newNode]), ShouldBeNil)
			So(clus.start(newNode), ShouldBeNil)
			nrNode += 1

			Convey("Then GetMembership should return 6 peers", func(ctx C) {
				So(testutil.CompareWaitN(50, func() bool {
					cinfo, err := clus.nodes[restartNode].GetClusterInfo(context.TODO(), 0)
					if err != nil {
						return false
					}
					return len(cinfo.Members) == nrNode
				}), ShouldBeTrue)
			})
		})

		Convey("When follower restart and some node leave", func(ctx C) {
			restartNode := (leader + 1) % nrNode
			leaveNode := (leader + 2) % nrNode

			clus.stop(restartNode)
			clus.stop(leaveNode)

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			So(clus.nodes[leader].RemovePeer(rctx,
				types.ClusterID(0),
				clus.nodes[leaveNode].nodeID), ShouldBeNil)

			nrNode -= 1

			clus.restart(restartNode)

			Convey("Then GetMembership should return 4 peers", func(ctx C) {
				So(testutil.CompareWaitN(50, func() bool {
					cinfo, err := clus.nodes[restartNode].GetClusterInfo(context.TODO(), 0)
					if err != nil {
						return false
					}
					return len(cinfo.Members) == nrNode
				}), ShouldBeTrue)
			})
		})
	})
}

func TestMRLoadSnapshop(t *testing.T) {
	Convey("Given MR cluster which have snapshot", t, func(ctx C) {
		testSnapCount = 10
		defer func() { testSnapCount = 0 }()
		nrRep := 1
		nrNode := 3

		clus := newMetadataRepoCluster(nrNode, nrRep, false)
		clus.Start()
		Reset(func() {
			clus.closeNoErrors(t)
		})
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		leader := clus.leader()
		So(leader, ShouldBeGreaterThan, -1)
		restartNode := (leader + 1) % nrNode

		snIDs := make([]types.StorageNodeID, nrRep)
		for i := range snIDs {
			snIDs[i] = types.StorageNodeID(i)

			sn := &varlogpb.StorageNodeDescriptor{
				StorageNodeID: snIDs[i],
			}

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			err := clus.nodes[leader].RegisterStorageNode(rctx, sn)
			So(err, ShouldBeNil)
		}

		So(testutil.CompareWaitN(50, func() bool {
			snap := clus.nodes[restartNode].raftNode.loadSnapshot()
			return snap != nil
		}), ShouldBeTrue)

		Convey("When follower restart", func(ctx C) {
			clus.restart(restartNode)
			So(testutil.CompareWaitN(50, func() bool {
				cinfo, err := clus.nodes[restartNode].GetClusterInfo(context.TODO(), 0)
				if err != nil {
					return false
				}
				return len(cinfo.Members) == nrNode
			}), ShouldBeTrue)

			Convey("Then GetMembership should recover metadata", func(ctx C) {
				So(testutil.CompareWaitN(50, func() bool {
					meta, err := clus.nodes[restartNode].GetMetadata(context.TODO())
					if err != nil {
						return false
					}

					return meta.GetStorageNode(snIDs[0]) != nil
				}), ShouldBeTrue)
			})
		})
	})
}

func TestMRRemoteSnapshot(t *testing.T) {
	Convey("Given MR cluster which have snapshot", t, func(ctx C) {
		testSnapCount = 10
		defer func() { testSnapCount = 0 }()
		nrRep := 1
		nrNode := 3

		clus := newMetadataRepoCluster(nrNode, nrRep, false)
		clus.Start()
		Reset(func() {
			clus.closeNoErrors(t)
		})
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		leader := clus.leader()
		So(leader, ShouldBeGreaterThan, -1)

		snIDs := make([]types.StorageNodeID, nrRep)
		for i := range snIDs {
			snIDs[i] = types.StorageNodeID(i)

			sn := &varlogpb.StorageNodeDescriptor{
				StorageNodeID: snIDs[i],
			}

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			err := clus.nodes[leader].RegisterStorageNode(rctx, sn)
			So(err, ShouldBeNil)
		}

		So(testutil.CompareWaitN(50, func() bool {
			snap := clus.nodes[leader].raftNode.loadSnapshot()
			return snap != nil
		}), ShouldBeTrue)

		Convey("When new node join", func(ctx C) {
			newNode := nrNode
			So(clus.appendMetadataRepo(), ShouldBeNil)

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			So(clus.nodes[leader].AddPeer(rctx,
				types.ClusterID(0),
				clus.nodes[newNode].nodeID,
				clus.peers[newNode]), ShouldBeNil)
			So(clus.start(newNode), ShouldBeNil)
			nrNode += 1

			Convey("Then GetMembership should recover metadata", func(ctx C) {
				So(testutil.CompareWaitN(50, func() bool {
					cinfo, err := clus.nodes[newNode].GetClusterInfo(context.TODO(), 0)
					if err != nil {
						return false
					}
					return len(cinfo.Members) == nrNode
				}), ShouldBeTrue)

				Convey("Then replication should be operate", func(ctx C) {
					for i := range snIDs {
						rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
						defer cancel()

						err := clus.nodes[leader].UnregisterStorageNode(rctx, snIDs[i])
						So(err, ShouldBeNil)
					}
				})
			})
		})
	})
}

func TestMRRemoteSnapshotFail(t *testing.T) {
	Convey("Given MR cluster which have snapshot", t, func(ctx C) {
		Convey("When new node join but sendSnapshot fail", func(ctx C) {
			Convey("Then replication should be operate", func(ctx C) {
			})
		})
	})
}

func TestMRFailoverRestartWithSnapshot(t *testing.T) {
	Convey("Given MR cluster with 5 peers", t, func(ctx C) {
		nrRep := 1
		nrNode := 5
		testSnapCount = 10
		defer func() { testSnapCount = 0 }()

		clus := newMetadataRepoCluster(nrNode, nrRep, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		clus.Start()
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		leader := clus.leader()
		So(leader, ShouldBeGreaterThan, -1)

		Convey("When follower restart with snapshot and some node leave", func(ctx C) {
			restartNode := (leader + 1) % nrNode
			leaveNode := (leader + 2) % nrNode

			So(testutil.CompareWaitN(50, func() bool {
				hasSnapshot, _ := clus.hasSnapshot(restartNode)
				return hasSnapshot
			}), ShouldBeTrue)

			clus.stop(restartNode)
			clus.stop(leaveNode)

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()

			So(clus.nodes[leader].RemovePeer(rctx,
				types.ClusterID(0),
				clus.nodes[leaveNode].nodeID), ShouldBeNil)

			nrNode -= 1

			clus.restart(restartNode)

			Convey("Then GetMembership should return 4 peers", func(ctx C) {
				So(testutil.CompareWaitN(50, func() bool {
					cinfo, err := clus.nodes[restartNode].GetClusterInfo(context.TODO(), 0)
					if err != nil {
						return false
					}
					return len(cinfo.Members) == nrNode
				}), ShouldBeTrue)
			})
		})
	})
}

func TestMRFailoverRecoverReportCollector(t *testing.T) {
	Convey("Given MR cluster with 3 peers, 5 StorageNodes", t, func(ctx C) {
		nrRep := 1
		nrNode := 3
		nrStorageNode := 3
		testSnapCount = 10

		clus := newMetadataRepoCluster(nrNode, nrRep, false)
		Reset(func() {
			clus.closeNoErrors(t)
			testSnapCount = 0
		})
		clus.Start()
		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		leader := clus.leader()
		So(leader, ShouldBeGreaterThan, -1)

		snIds := make([][]types.StorageNodeID, nrStorageNode)
		for i := range snIds {
			snIds[i] = make([]types.StorageNodeID, nrRep)
			for j := range snIds[i] {
				snIds[i][j] = types.StorageNodeID(i*nrStorageNode + j)

				sn := &varlogpb.StorageNodeDescriptor{
					StorageNodeID: snIds[i][j],
				}

				rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
				defer cancel()

				err := clus.nodes[leader].RegisterStorageNode(rctx, sn)
				So(err, ShouldBeNil)
			}
		}

		Convey("When follower restart with snapshot", func(ctx C) {
			restartNode := (leader + 1) % nrNode

			So(testutil.CompareWaitN(50, func() bool {
				hasSnapshot, _ := clus.hasSnapshot(restartNode)
				return hasSnapshot
			}), ShouldBeTrue)

			clus.stop(restartNode)
			clus.restart(restartNode)

			Convey("Then ReportCollector should recover", func(ctx C) {
				So(testutil.CompareWaitN(50, func() bool {
					clus.nodes[restartNode].reportCollector.mu.RLock()
					defer clus.nodes[restartNode].reportCollector.mu.RUnlock()

					return len(clus.nodes[restartNode].reportCollector.executors) == nrStorageNode
				}), ShouldBeTrue)
			})
		})
	})
}

func TestMRProposeTimeout(t *testing.T) {
	Convey("Given MR which is not running", t, func(ctx C) {
		clus := newMetadataRepoCluster(1, 1, false)
		Reset(func() {
			clus.closeNoErrors(t)
		})
		mr := clus.nodes[0]

		Convey("When cli register SN with timeout", func(ctx C) {
			snID := types.StorageNodeID(time.Now().UnixNano())
			sn := &varlogpb.StorageNodeDescriptor{
				StorageNodeID: types.StorageNodeID(snID),
			}

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()
			err := mr.RegisterStorageNode(rctx, sn)
			Convey("Then it should be timed out", func(ctx C) {
				So(err, ShouldResemble, context.DeadlineExceeded)
			})
		})
	})
}

func TestMRProposeRetry(t *testing.T) {
	Convey("Given MR", t, func(ctx C) {
		clus := newMetadataRepoCluster(3, 1, false)
		clus.Start()
		Reset(func() {
			clus.closeNoErrors(t)
		})

		So(testutil.CompareWaitN(10, func() bool {
			return clus.leaderElected()
		}), ShouldBeTrue)

		Convey("When cli register SN & transfer leader for dropping propose", func(ctx C) {
			leader := clus.leader()
			clus.nodes[leader].raftNode.transferLeadership(false)

			snID := types.StorageNodeID(time.Now().UnixNano())
			sn := &varlogpb.StorageNodeDescriptor{
				StorageNodeID: types.StorageNodeID(snID),
			}

			rctx, cancel := context.WithTimeout(context.Background(), vtesting.TimeoutUnitTimesFactor(50))
			defer cancel()
			err := clus.nodes[leader].RegisterStorageNode(rctx, sn)

			Convey("Then it should be success", func(ctx C) {
				So(err, ShouldBeNil)
				//So(atomic.LoadUint64(&clus.nodes[leader].requestNum), ShouldBeGreaterThan, 1)
			})
		})
	})
}
