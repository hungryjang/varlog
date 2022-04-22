package executor

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/kakao/varlog/internal/storagenode_deprecated/id"
	"github.com/kakao/varlog/internal/storagenode_deprecated/reportcommitter"
	"github.com/kakao/varlog/internal/storagenode_deprecated/storage"
	"github.com/kakao/varlog/internal/storagenode_deprecated/telemetry"
	"github.com/kakao/varlog/pkg/types"
	"github.com/kakao/varlog/proto/snpb"
	"github.com/kakao/varlog/proto/varlogpb"
)

func TestLogStreamReporter(t *testing.T) {
	const (
		snid    = types.StorageNodeID(1)
		lsid1   = types.LogStreamID(1)
		lsid2   = types.LogStreamID(2)
		topicID = 1
	)

	defer goleak.VerifyNone(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	strg1, err := storage.NewStorage(storage.WithPath(t.TempDir()))
	require.NoError(t, err)
	lse1, err := New(
		WithStorageNodeID(snid),
		WithLogStreamID(lsid1),
		WithStorage(strg1),
		WithMetrics(telemetry.NewMetrics()),
	)
	require.NoError(t, err)
	defer func() { require.NoError(t, lse1.Close()) }()

	status, sealedGLSN, err := lse1.Seal(context.TODO(), types.InvalidGLSN)
	require.NoError(t, err)
	require.Equal(t, types.InvalidGLSN, sealedGLSN)
	require.Equal(t, varlogpb.LogStreamStatusSealed, status)
	require.NoError(t, lse1.Unseal(context.TODO(), []varlogpb.LogStreamReplica{
		{
			StorageNode: varlogpb.StorageNode{
				StorageNodeID: lse1.storageNodeID,
			},
			TopicLogStream: varlogpb.TopicLogStream{
				LogStreamID: lse1.logStreamID,
			},
		},
	}))
	require.Equal(t, varlogpb.LogStreamStatusRunning, lse1.Metadata().Status)

	strg2, err := storage.NewStorage(storage.WithPath(t.TempDir()))
	require.NoError(t, err)
	lse2, err := New(
		WithStorageNodeID(snid),
		WithLogStreamID(lsid2),
		WithStorage(strg2),
		WithMetrics(telemetry.NewMetrics()),
	)
	require.NoError(t, err)
	defer func() { require.NoError(t, lse2.Close()) }()

	status, sealedGLSN, err = lse2.Seal(context.TODO(), types.InvalidGLSN)
	require.NoError(t, err)
	require.Equal(t, types.InvalidGLSN, sealedGLSN)
	require.Equal(t, varlogpb.LogStreamStatusSealed, status)
	require.NoError(t, lse2.Unseal(context.TODO(), []varlogpb.LogStreamReplica{
		{
			StorageNode: varlogpb.StorageNode{
				StorageNodeID: lse2.storageNodeID,
			},
			TopicLogStream: varlogpb.TopicLogStream{
				LogStreamID: lse2.logStreamID,
			},
		},
	}))
	require.Equal(t, varlogpb.LogStreamStatusRunning, lse2.Metadata().Status)

	rcg := reportcommitter.NewMockGetter(ctrl)
	rcg.EXPECT().ReportCommitter(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ types.TopicID, lsid types.LogStreamID) (reportcommitter.ReportCommitter, bool) {
			switch lsid {
			case lsid1:
				return lse1, true
			case lsid2:
				return lse2, true
			default:
				return nil, false
			}
		},
	).AnyTimes()
	rcg.EXPECT().GetReports(gomock.Any(), gomock.Any()).DoAndReturn(
		func(rsp *snpb.GetReportResponse, f func(reportcommitter.ReportCommitter, *snpb.GetReportResponse)) {
			f(lse1, rsp)
			f(lse2, rsp)
		},
	).AnyTimes()

	_, err = New()
	require.Error(t, err)

	getter := id.NewMockStorageNodeIDGetter(ctrl)
	getter.EXPECT().StorageNodeID().Return(snid).AnyTimes()

	lsr := reportcommitter.New(
		reportcommitter.WithStorageNodeIDGetter(getter),
		reportcommitter.WithReportCommitterGetter(rcg),
	)

	err = lsr.Commit(context.TODO(), snpb.LogStreamCommitResult{})
	require.Error(t, err)

	var (
		reports []snpb.LogStreamUncommitReport
		wg      sync.WaitGroup
	)

	// Reports
	rsp := snpb.GetReportResponse{}
	err = lsr.GetReport(context.TODO(), &rsp)
	require.NoError(t, err)

	reports = rsp.UncommitReports
	require.Len(t, reports, 2)
	for _, report := range reports {
		switch report.GetLogStreamID() {
		case lsid1:
			require.Equal(t, types.Version(0), report.GetVersion())
			require.Equal(t, types.LLSN(1), report.GetUncommittedLLSNOffset())
			require.EqualValues(t, 0, report.GetUncommittedLLSNLength())
		case lsid2:
			require.Equal(t, types.Version(0), report.GetVersion())
			require.Equal(t, types.LLSN(1), report.GetUncommittedLLSNOffset())
			require.EqualValues(t, 0, report.GetUncommittedLLSNLength())
		}
	}

	// Append
	// LSE1
	// LSE2
	wg.Add(3)
	go func() {
		defer wg.Done()

		res, err := lse1.Append(context.TODO(), [][]byte{[]byte("foo")})
		require.NoError(t, err)
		require.Equal(t, types.GLSN(1), res[0].Meta.GLSN)
	}()
	go func() {
		defer wg.Done()

		res, err := lse2.Append(context.TODO(), [][]byte{[]byte("foo")})
		require.NoError(t, err)
		require.Equal(t, types.GLSN(2), res[0].Meta.GLSN)
	}()
	go func() {
		defer wg.Done()

		require.Eventually(t, func() bool {
			err := lsr.GetReport(context.TODO(), &rsp)
			require.NoError(t, err)

			reports := rsp.UncommitReports
			return reports[0].UncommittedLLSNLength > 0 && reports[1].UncommittedLLSNLength > 0
		}, time.Second, time.Millisecond)

		require.NoError(t, lsr.Commit(context.TODO(), snpb.LogStreamCommitResult{
			LogStreamID:         lsid1,
			Version:             1,
			CommittedGLSNOffset: 1,
			CommittedGLSNLength: 1,
			CommittedLLSNOffset: 1,
		}))

		require.NoError(t, lsr.Commit(context.TODO(), snpb.LogStreamCommitResult{
			LogStreamID:         lsid2,
			Version:             1,
			CommittedGLSNOffset: 2,
			CommittedGLSNLength: 1,
			CommittedLLSNOffset: 1,
		}))
	}()
	wg.Wait()

	require.Eventually(t, func() bool {
		err := lsr.GetReport(context.TODO(), &rsp)
		require.NoError(t, err)

		reports := rsp.UncommitReports
		return reports[0].UncommittedLLSNLength == 0 && reports[1].UncommittedLLSNLength == 0
	}, time.Second, time.Millisecond)

	require.NoError(t, lsr.Close())
	require.NoError(t, lsr.Close())

	err = lsr.GetReport(context.TODO(), &rsp)
	require.Error(t, err)

	require.Error(t, lsr.Commit(context.TODO(), snpb.LogStreamCommitResult{
		LogStreamID:         lsid1,
		Version:             2,
		CommittedGLSNOffset: 3,
		CommittedGLSNLength: 1,
		CommittedLLSNOffset: 2,
	},
	))

	require.Error(t, lsr.Commit(context.TODO(), snpb.LogStreamCommitResult{
		LogStreamID:         lsid2,
		Version:             2,
		CommittedGLSNOffset: 4,
		CommittedGLSNLength: 1,
		CommittedLLSNOffset: 2,
	}))
}