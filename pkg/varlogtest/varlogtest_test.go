package varlogtest_test

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/kakao/varlog/pkg/types"
	"github.com/kakao/varlog/pkg/util/container/set"
	"github.com/kakao/varlog/pkg/varlogtest"
	"github.com/kakao/varlog/proto/varlogpb"
)

func TestVarlogTest(t *testing.T) {
	defer goleak.VerifyNone(t)

	const (
		clusterID         = types.ClusterID(1)
		replicationFactor = 3

		numStorageNodes       = 5
		numTopics             = 10
		avgLogStreamsPerTopic = 10
		numLogStreams         = numTopics * avgLogStreamsPerTopic / replicationFactor
		avgLogsPerTopic       = 100
		numLogs               = numTopics * avgLogsPerTopic
	)

	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	vt := varlogtest.New(clusterID, replicationFactor)
	admin := vt.Admin()
	varlog := vt.Log()
	defer func() {
		require.NoError(t, varlog.Close())
		require.NoError(t, admin.Close())
	}()

	var (
		topicIDs     = make([]types.TopicID, 0, numTopics)
		logStreamIDs = make([]types.LogStreamID, 0, numLogStreams)
		hwms         = make(map[types.TopicID]types.GLSN, numTopics)
	)

	// Add topics
	for i := 0; i < numTopics; i++ {
		topicDesc, err := admin.AddTopic(context.Background())
		require.NoError(t, err)
		require.Equal(t, varlogpb.TopicStatusRunning, topicDesc.Status)
		require.Empty(t, topicDesc.LogStreams)
		require.NotContains(t, topicIDs, topicDesc.TopicID)
		topicIDs = append(topicIDs, topicDesc.TopicID)
	}

	// Append logs, but no log stream
	for i := 0; i < numLogs; i++ {
		tpID := topicIDs[rng.Intn(numTopics)]
		_, err := varlog.Append(context.Background(), tpID, nil)
		require.Error(t, err)
	}

	// Add log streams, but no storage node
	for i := 0; i < numLogStreams; i++ {
		tpID := topicIDs[rng.Intn(numTopics)]
		_, err := admin.AddLogStream(context.Background(), tpID, nil)
		require.Error(t, err)
	}

	// Add storage nodes
	for i := 0; i < numStorageNodes; i++ {
		addr := fmt.Sprintf("sn%03d", i+1)
		snMetaDesc, err := admin.AddStorageNode(context.Background(), addr)
		require.NoError(t, err)
		require.Equal(t, clusterID, snMetaDesc.ClusterID)
		require.Empty(t, snMetaDesc.LogStreams)
		snDesc := snMetaDesc.StorageNode
		require.Equal(t, varlogpb.StorageNodeStatusRunning, snDesc.Status)
		require.Equal(t, addr, snDesc.Address)
		require.NotEmpty(t, snDesc.Storages)
	}

	// Add log streams
	addLogStream := func(tpID types.TopicID) {
		lsDesc, err := admin.AddLogStream(context.Background(), tpID, nil)
		require.NoError(t, err)
		require.Equal(t, tpID, lsDesc.TopicID)
		require.Equal(t, varlogpb.LogStreamStatusRunning, lsDesc.Status)
		require.Len(t, lsDesc.Replicas, replicationFactor)

		snIDSet := set.New(replicationFactor)
		for _, replicaDesc := range lsDesc.Replicas {
			require.NotContains(t, snIDSet, replicaDesc.StorageNodeID)
			require.NotEmpty(t, replicaDesc.Path)
			snIDSet.Add(replicaDesc.StorageNodeID)
		}
		require.Len(t, snIDSet, replicationFactor)

		logStreamIDs = append(logStreamIDs, lsDesc.LogStreamID)
	}
	for i := 0; i < numTopics; i++ {
		addLogStream(topicIDs[i])
	}
	for i := 0; i < numLogStreams-numTopics; i++ {
		tpID := topicIDs[rng.Intn(numTopics)]
		addLogStream(tpID)
	}

	// Append logs
	appendLog := func(tpID types.TopicID) {
		hwms[tpID]++
		data := []byte(fmt.Sprintf("%d,%d", tpID, hwms[tpID]))
		actualGLSN, err := varlog.Append(context.Background(), tpID, data)
		require.NoError(t, err)
		require.Equal(t, hwms[tpID], actualGLSN)
	}
	for i := 0; i < numTopics; i++ {
		appendLog(topicIDs[i])
	}
	for i := 0; i < numLogs-numTopics; i++ {
		tpID := topicIDs[rng.Intn(numTopics)]
		appendLog(tpID)
	}

	// Subscribe
	subscribe := func(tpID types.TopicID, begin, end types.GLSN) {
		expectedGLSN := begin
		onNext := func(logEntry varlogpb.LogEntry, err error) {
			if err != nil {
				require.ErrorIs(t, err, io.EOF)
				return
			}
			require.Equal(t, expectedGLSN, logEntry.GLSN)
			require.Equal(t, []byte(fmt.Sprintf("%d,%d", tpID, expectedGLSN)), logEntry.Data)
			expectedGLSN++
		}
		closer, err := varlog.Subscribe(context.Background(), tpID, types.MinGLSN, hwms[tpID]+1, onNext)
		require.NoError(t, err)
		closer()
		require.Equal(t, end, expectedGLSN)
	}
	for i := 0; i < numTopics; i++ {
		tpID := topicIDs[i]
		subscribe(tpID, types.MinGLSN, hwms[tpID]+1)
	}
}