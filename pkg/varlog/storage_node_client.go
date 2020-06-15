package varlog

import (
	"context"
	"io"

	pb "github.com/kakao/varlog/proto/storage_node"
)

type StorageNode struct {
	Id   StorageNodeID
	Addr string
}

type StorageNodeClient interface {
	Append(ctx context.Context, logStreamID LogStreamID, data []byte, backups ...StorageNode) (uint64, error)
	Read(ctx context.Context, logStreamID LogStreamID, glsn uint64) ([]byte, error)
	Subscribe(ctx context.Context, glsn uint64) (<-chan []byte, error)
	Trim(ctx context.Context, glsn uint64) (uint64, error)
	Close() error
}

type storageNodeClient struct {
	rpcConn   *rpcConn
	rpcClient pb.StorageNodeServiceClient
}

func NewStorageNodeClient(address string) (StorageNodeClient, error) {
	rpcConn, err := newRpcConn(address)
	if err != nil {
		return nil, err
	}
	return NewStorageNodeClientFromRpcConn(rpcConn)
}

func NewStorageNodeClientFromRpcConn(rpcConn *rpcConn) (StorageNodeClient, error) {
	return &storageNodeClient{
		rpcConn:   rpcConn,
		rpcClient: pb.NewStorageNodeServiceClient(rpcConn.conn),
	}, nil
}

func (c *storageNodeClient) Append(ctx context.Context, logStreamID LogStreamID, data []byte, backups ...StorageNode) (uint64, error) {
	req := &pb.AppendRequest{
		Payload: data,
	}
	rsp, err := c.rpcClient.Append(ctx, req)
	if err != nil {
		return 0, err
	}
	return rsp.GetGlsn(), nil
}

func (c *storageNodeClient) Read(ctx context.Context, logStreamID LogStreamID, glsn uint64) ([]byte, error) {
	req := &pb.ReadRequest{
		Glsn: glsn,
	}
	rsp, err := c.rpcClient.Read(ctx, req)
	if err != nil {
		return nil, err
	}
	return rsp.GetPayload(), nil
}

func (c *storageNodeClient) Subscribe(ctx context.Context, glsn uint64) (<-chan []byte, error) {
	req := &pb.SubscribeRequest{Glsn: glsn}
	stream, err := c.rpcClient.Subscribe(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte)
	go func() {
		defer close(ch)
		for {
			rsp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic("not yet implemented")
			}
			ch <- rsp.GetPayload()
		}
	}()
	return ch, nil
}

func (c *storageNodeClient) Trim(ctx context.Context, glsn uint64) (uint64, error) {
	req := &pb.TrimRequest{
		Glsn:  glsn,
		Async: false,
	}
	rsp, err := c.rpcClient.Trim(ctx, req)
	if err != nil {
		return 0, err
	}
	return rsp.GetNumTrimmed(), nil
}

func (c *storageNodeClient) Close() error {
	return c.rpcConn.close()
}