package storage

import (
	"context"
	"sync"

	"github.com/kakao/varlog/pkg/varlog"
	"github.com/kakao/varlog/pkg/varlog/types"
	pb "github.com/kakao/varlog/proto/storage_node"
	"google.golang.org/grpc"
)

type StorageNodeService struct {
	pb.UnimplementedStorageNodeServiceServer

	storageNodeID types.StorageNodeID
	lseM          map[types.LogStreamID]LogStreamExecutor
	m             sync.RWMutex
}

func NewStorageNodeService(storageNodeID types.StorageNodeID) *StorageNodeService {
	return &StorageNodeService{
		storageNodeID: storageNodeID,
		lseM:          make(map[types.LogStreamID]LogStreamExecutor),
	}
}

func (s *StorageNodeService) Register(server *grpc.Server) {
	pb.RegisterStorageNodeServiceServer(server, s)
}

func (s *StorageNodeService) AddLogStream(ctx context.Context, req *pb.AddLogStreamRequest) (*pb.AddLogStreamResponse, error) {
	s.m.Lock()
	defer s.m.Unlock()

	id := req.GetLogStreamID()
	_, ok := s.lseM[id]
	if ok {
		// already exists
		return nil, varlog.ErrInvalid
	}
	storage := NewInMemoryStorage()
	lse, err := NewLogStreamExecutor(id, storage)
	if err != nil {
		return nil, err
	}
	s.lseM[id] = lse
	return &pb.AddLogStreamResponse{}, nil
}

func (s *StorageNodeService) getLogStreamExecutor(logStreamID types.LogStreamID) (LogStreamExecutor, bool) {
	s.m.RLock()
	defer s.m.RUnlock()
	lse, ok := s.lseM[logStreamID]
	return lse, ok
}

func (s *StorageNodeService) Append(ctx context.Context, req *pb.AppendRequest) (*pb.AppendResponse, error) {
	lse, ok := s.getLogStreamExecutor(req.GetLogStreamID())
	if !ok {
		return nil, varlog.ErrInvalid
	}
	// TODO: create child context by using operation timeout
	// TODO: create replicas by using request
	glsn, err := lse.Append(ctx, req.GetPayload())
	if err != nil {
		return nil, err
	}
	return &pb.AppendResponse{GLSN: glsn}, nil
}

func (s *StorageNodeService) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	lse, ok := s.getLogStreamExecutor(req.GetLogStreamID())
	if !ok {
		return nil, varlog.ErrInvalid
	}

	// TODO: create child context by using operation timeout
	data, err := lse.Read(ctx, req.GetGLSN())
	if err != nil {
		return nil, err
	}
	return &pb.ReadResponse{Payload: data, GLSN: req.GetGLSN()}, nil
}

func (s *StorageNodeService) Subscribe(*pb.SubscribeRequest, pb.StorageNodeService_SubscribeServer) error {
	panic("not yet implemented")
}

func (s *StorageNodeService) Trim(ctx context.Context, req *pb.TrimRequest) (*pb.TrimResponse, error) {
	s.m.RLock()
	targetLSEs := make([]LogStreamExecutor, len(s.lseM))
	i := 0
	for _, lse := range s.lseM {
		targetLSEs[i] = lse
		i++
	}
	s.m.RUnlock()

	// NOTE: subtle case
	// If the trim operation will remove very large GLSN that is not stored yet, current LSEs
	// remove all log entries. After replied the trim operation, log entries within the scope
	// of removing will be saved again.

	// TODO: create child context by using operation timeout
	type result struct {
		num uint64
		err error
	}

	c := make(chan result, len(s.lseM))
	var wg sync.WaitGroup
	wg.Add(len(targetLSEs))
	for _, lse := range targetLSEs {
		go func(lse LogStreamExecutor) {
			defer wg.Done()
			cnt, err := lse.Trim(ctx, req.GetGLSN(), req.GetAsync())
			c <- result{cnt, err}
		}(lse)
	}
	wg.Wait()
	close(c)
	rsp := &pb.TrimResponse{}
	for res := range c {
		rsp.NumTrimmed += res.num
		if res.err != nil {
			return nil, res.err
		}
	}
	return rsp, nil
}
