package logio

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/kakao/varlog/internal/storagenode_deprecated/rpcserver"
	"github.com/kakao/varlog/pkg/verrors"
	"github.com/kakao/varlog/proto/snpb"
	"github.com/kakao/varlog/proto/varlogpb"
)

type Server interface {
	snpb.LogIOServer
	rpcserver.Registrable
}

type server struct {
	config
}

func NewServer(opts ...Option) Server {
	cfg := newConfig(opts)
	return &server{config: cfg}
}

var _ Server = (*server)(nil)

func (s *server) Register(server *grpc.Server) {
	s.logger.Info("register to rpcserver server")
	snpb.RegisterLogIOServer(server, s)
}

func (s *server) withTelemetry(ctx context.Context, spanName string, req interface{}, h rpcserver.Handler) (rsp interface{}, err error) {
	rsp, err = h(ctx, req)
	if err != nil {
		s.logger.Error(spanName,
			zap.Error(err),
			zap.Stringer("request", req.(fmt.Stringer)),
		)
	}
	return rsp, err
}

func (s *server) Append(ctx context.Context, req *snpb.AppendRequest) (*snpb.AppendResponse, error) {
	code := codes.Internal
	rspI, err := s.withTelemetry(ctx, "varlog.snpb.LogIO/Append", req,
		func(ctx context.Context, reqI interface{}) (interface{}, error) {
			startTime := time.Now()
			defer func() {
				dur := time.Since(startTime)
				s.metrics.RPCServerAppendDuration.Record(
					ctx,
					float64(dur.Microseconds())/1000.0,
				)
			}()

			var rsp *snpb.AppendResponse
			lse, ok := s.readWriterGetter.ReadWriter(req.GetTopicID(), req.GetLogStreamID())
			if !ok {
				code = codes.NotFound
				return rsp, errors.WithStack(verrors.ErrInvalid)
			}

			backups := make([]varlogpb.LogStreamReplica, 0, len(req.Backups))
			for i := range req.Backups {
				backups = append(backups, varlogpb.LogStreamReplica{
					StorageNode: varlogpb.StorageNode{
						StorageNodeID: req.Backups[i].GetStorageNodeID(),
						Address:       req.Backups[i].GetAddress(),
					},
					TopicLogStream: varlogpb.TopicLogStream{
						TopicID:     req.GetTopicID(),
						LogStreamID: req.GetLogStreamID(),
					},
				})
			}

			res, err := lse.Append(ctx, req.GetPayload(), backups...)
			if err != nil {
				code = codes.Internal
				return rsp, err
			}
			return &snpb.AppendResponse{Results: res}, nil
		},
	)
	return rspI.(*snpb.AppendResponse), verrors.ToStatusErrorWithCode(err, code)
}

func (s *server) Read(ctx context.Context, req *snpb.ReadRequest) (*snpb.ReadResponse, error) {
	code := codes.Internal
	rspI, err := s.withTelemetry(ctx, "varlog.snpb.LogIO/Read", req,
		func(ctx context.Context, reqI interface{}) (interface{}, error) {
			var rsp *snpb.ReadResponse
			lse, ok := s.readWriterGetter.ReadWriter(req.GetTopicID(), req.GetLogStreamID())
			if !ok {
				code = codes.NotFound
				return rsp, errors.WithStack(verrors.ErrInvalid)
			}

			logEntry, err := lse.Read(ctx, req.GetGLSN())
			if err != nil {
				// TODO: Check whether these are safe.
				switch errors.Cause(err) {
				case verrors.ErrNoEntry:
					code = codes.NotFound
				case verrors.ErrTrimmed:
					code = codes.OutOfRange
				case verrors.ErrUndecidable:
					code = codes.Unavailable
				default:
					code = codes.Internal
				}
				/*
					if errors.Is(err, verrors.ErrNoEntry) {
						code = codes.NotFound
					} else if errors.Is(err, verrors.ErrTrimmed) {
						code = codes.OutOfRange
					} else if errors.Is(err, verrors.ErrUndecidable) {
						// TODO (jun): consider codes.FailedPrecondition
						code = codes.Unavailable
					} else {
						code = codes.Internal
					}
				*/
				return rsp, errors.Wrap(err, "storagenode")
			}
			return &snpb.ReadResponse{
				Payload: logEntry.Data,
				GLSN:    req.GetGLSN(),
				LLSN:    logEntry.LLSN,
			}, nil
		},
	)
	return rspI.(*snpb.ReadResponse), verrors.ToStatusErrorWithCode(err, code)
}

func (s *server) Subscribe(req *snpb.SubscribeRequest, stream snpb.LogIO_SubscribeServer) error {
	code := codes.Internal
	_, err := s.withTelemetry(stream.Context(), "varlog.snpb.LogIO/Subscribe", req,
		func(ctx context.Context, reqI interface{}) (interface{}, error) {
			if req.GetGLSNBegin() >= req.GetGLSNEnd() {
				code = codes.InvalidArgument
				return nil, errors.New("storagenode: invalid subscription range")
			}
			reader, ok := s.readWriterGetter.ReadWriter(req.GetTopicID(), req.GetLogStreamID())
			if !ok {
				code = codes.NotFound
				return nil, errors.WithStack(verrors.ErrInvalid)
			}

			subEnv, err := reader.Subscribe(ctx, req.GetGLSNBegin(), req.GetGLSNEnd())
			if err != nil {
				return nil, err
			}
			// FIXME: monitor stream's context, and stop subEnv if the context is canceled.
			defer subEnv.Stop()

			for sr := range subEnv.ScanResultC() {
				if err := stream.Send(&snpb.SubscribeResponse{
					GLSN:    sr.LogEntry.GLSN,
					LLSN:    sr.LogEntry.LLSN,
					Payload: sr.LogEntry.Data,
				}); err != nil {
					return nil, errors.WithStack(err)
				}
			}
			// FIXME: if the subscribe is finished without critical error (i.e., other than io.EOF), Err()
			// should return nil.
			err = subEnv.Err()
			if err == io.EOF {
				err = nil
			}
			return nil, err
		},
	)
	return verrors.ToStatusErrorWithCode(err, code)
}

func (s *server) SubscribeTo(req *snpb.SubscribeToRequest, stream snpb.LogIO_SubscribeToServer) error {
	code := codes.Internal
	_, err := s.withTelemetry(stream.Context(), "varlog.snpb.LogIO/Subscribe", req,
		func(ctx context.Context, _ interface{}) (interface{}, error) {
			if req.LLSNBegin >= req.LLSNEnd {
				code = codes.InvalidArgument
				return nil, errors.New("storagenode: invalid subscription range")
			}

			reader, ok := s.readWriterGetter.ReadWriter(req.TopicID, req.LogStreamID)
			if !ok {
				code = codes.NotFound
				return nil, errors.WithStack(verrors.ErrInvalid)
			}

			subEnv, err := reader.SubscribeTo(ctx, req.LLSNBegin, req.LLSNEnd)
			if err != nil {
				return nil, err
			}
			// FIXME: monitor stream's context, and stop subEnv if the context is canceled.
			defer subEnv.Stop()

			for sr := range subEnv.ScanResultC() {
				if err := stream.Send(&snpb.SubscribeToResponse{
					LogEntry: varlogpb.LogEntry{
						LogEntryMeta: varlogpb.LogEntryMeta{
							TopicID:     req.TopicID,
							LogStreamID: req.LogStreamID,
							GLSN:        sr.LogEntry.GLSN,
							LLSN:        sr.LogEntry.LLSN,
						},
						Data: sr.LogEntry.Data,
					},
				}); err != nil {
					return nil, errors.WithStack(err)
				}
			}
			// FIXME: if the subscribe is finished without critical error (i.e., other than io.EOF), Err()
			// should return nil.
			err = subEnv.Err()
			if err == io.EOF {
				err = nil
			}
			return nil, err
		},
	)
	return verrors.ToStatusErrorWithCode(err, code)
}

func (s *server) TrimDeprecated(ctx context.Context, req *snpb.TrimDeprecatedRequest) (*pbtypes.Empty, error) {
	code := codes.Internal
	rspI, err := s.withTelemetry(ctx, "varlog.snpb.LogIO/TrimDeprecated", req,
		func(ctx context.Context, reqI interface{}) (interface{}, error) {
			trimGLSN := req.GetGLSN()

			// TODO
			var wg sync.WaitGroup
			var err error
			var mu sync.Mutex
			s.readWriterGetter.ForEachReadWriters(func(rw ReadWriter) {
				readWriter := rw
				wg.Add(1)
				go func() {
					defer wg.Done()
					cerr := readWriter.Trim(ctx, trimGLSN)
					mu.Lock()
					err = multierr.Append(err, cerr)
					mu.Unlock()
				}()
			})
			wg.Wait()
			return &pbtypes.Empty{}, nil
		},
	)
	return rspI.(*pbtypes.Empty), verrors.ToStatusErrorWithCode(err, code)
}

func (s *server) LogStreamMetadata(ctx context.Context, req *snpb.LogStreamMetadataRequest) (*snpb.LogStreamMetadataResponse, error) {
	code := codes.Internal
	rspI, err := s.withTelemetry(ctx, "varlog.snpb.LogIO/LogStreamMetadata", req,
		func(ctx context.Context, _ interface{}) (interface{}, error) {
			var rsp *snpb.LogStreamMetadataResponse
			lse, exist := s.readWriterGetter.ReadWriter(req.TopicID, req.LogStreamID)
			if !exist {
				code = codes.NotFound
				return rsp, errors.WithStack(verrors.ErrInvalid)
			}

			lsd, err := lse.LogStreamMetadata()
			if err != nil {
				return rsp, err
			}

			return &snpb.LogStreamMetadataResponse{LogStreamDescriptor: lsd}, nil
		},
	)
	return rspI.(*snpb.LogStreamMetadataResponse), verrors.ToStatusErrorWithCode(err, code)
}