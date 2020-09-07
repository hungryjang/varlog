package varlog

import (
	"context"

	"github.com/kakao/varlog/pkg/varlog/types"
	"github.com/kakao/varlog/pkg/varlog/util/stringsutil"
	snpb "github.com/kakao/varlog/proto/storage_node"
	varlogpb "github.com/kakao/varlog/proto/varlog"
)

type ManagementClient interface {
	GetMetadata(ctx context.Context, clusterID types.ClusterID, metadataType snpb.MetadataType) (*varlogpb.StorageNodeMetadataDescriptor, error)
	AddLogStream(ctx context.Context, clusterID types.ClusterID,
		storageNodeID types.StorageNodeID, logStreamID types.LogStreamID, path string) error
	RemoveLogStream(ctx context.Context, clusterID types.ClusterID,
		storageNodeID types.StorageNodeID, logStreamID types.LogStreamID) error
	Seal(ctx context.Context, clusterID types.ClusterID,
		storageNodeID types.StorageNodeID, logStreamID types.LogStreamID,
		lastCommittedGLSN types.GLSN) (varlogpb.LogStreamStatus, types.GLSN, error)
	Unseal(ctx context.Context, clusterID types.ClusterID, storageNodeID types.StorageNodeID,
		logStreamID types.LogStreamID) error
}

type managementClient struct {
	rpcConn   *RpcConn
	rpcClient snpb.ManagementClient
}

func NewManagementClient(address string) (ManagementClient, error) {
	rpcConn, err := NewRpcConn(address)
	if err != nil {
		return nil, err
	}
	return &managementClient{
		rpcConn:   rpcConn,
		rpcClient: snpb.NewManagementClient(rpcConn.Conn),
	}, nil
}

func (c *managementClient) GetMetadata(ctx context.Context, clusterID types.ClusterID, metadataType snpb.MetadataType) (*varlogpb.StorageNodeMetadataDescriptor, error) {
	rsp, err := c.rpcClient.GetMetadata(ctx, &snpb.GetMetadataRequest{
		ClusterID:    clusterID,
		MetadataType: metadataType,
	})
	if err != nil {
		return nil, err
	}
	return rsp.GetStorageNodeMetadata(), nil
}

func (c *managementClient) AddLogStream(ctx context.Context, cid types.ClusterID, snid types.StorageNodeID, lsid types.LogStreamID, path string) error {
	// TODO(jun): Check ranges CID, SNID and LSID
	if stringsutil.Empty(path) {
		return ErrInvalid // FIXME: ErrInvalid ErrInvalidArgument
	}
	// FIXME(jun): Does the return value of AddLogStream need?
	_, err := c.rpcClient.AddLogStream(ctx, &snpb.AddLogStreamRequest{
		ClusterID:     cid,
		StorageNodeID: snid,
		LogStreamID:   lsid,
		Storage: &varlogpb.StorageDescriptor{
			Path: path,
		},
	})
	return err

}

func (c *managementClient) RemoveLogStream(ctx context.Context, cid types.ClusterID, snid types.StorageNodeID, lsid types.LogStreamID) error {
	// TODO(jun): Check ranges CID, SNID and LSID
	_, err := c.rpcClient.RemoveLogStream(ctx, &snpb.RemoveLogStreamRequest{
		ClusterID:     cid,
		StorageNodeID: snid,
		LogStreamID:   lsid,
	})
	return err
}

func (c *managementClient) Seal(ctx context.Context, cid types.ClusterID, snid types.StorageNodeID, lsid types.LogStreamID, lastCommittedGLSN types.GLSN) (varlogpb.LogStreamStatus, types.GLSN, error) {
	// TODO(jun): Check ranges CID, SNID and LSID
	rsp, err := c.rpcClient.Seal(ctx, &snpb.SealRequest{
		ClusterID:         cid,
		StorageNodeID:     snid,
		LogStreamID:       lsid,
		LastCommittedGLSN: lastCommittedGLSN,
	})
	return rsp.GetStatus(), rsp.GetLastCommittedGLSN(), err
}

func (c *managementClient) Unseal(ctx context.Context, cid types.ClusterID, snid types.StorageNodeID, lsid types.LogStreamID) error {
	// TODO(jun): Check ranges CID, SNID and LSID
	_, err := c.rpcClient.Unseal(ctx, &snpb.UnsealRequest{
		ClusterID:     cid,
		StorageNodeID: snid,
		LogStreamID:   lsid,
	})
	return err
}

func (c *managementClient) Sync(ctx context.Context) error {
	panic("not yet implemented")
}
