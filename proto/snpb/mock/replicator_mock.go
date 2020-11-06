// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kakao/varlog/proto/snpb (interfaces: ReplicatorServiceClient,ReplicatorServiceServer,ReplicatorService_ReplicateClient,ReplicatorService_ReplicateServer)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"

	snpb "github.com/kakao/varlog/proto/snpb"
)

// MockReplicatorServiceClient is a mock of ReplicatorServiceClient interface.
type MockReplicatorServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockReplicatorServiceClientMockRecorder
}

// MockReplicatorServiceClientMockRecorder is the mock recorder for MockReplicatorServiceClient.
type MockReplicatorServiceClientMockRecorder struct {
	mock *MockReplicatorServiceClient
}

// NewMockReplicatorServiceClient creates a new mock instance.
func NewMockReplicatorServiceClient(ctrl *gomock.Controller) *MockReplicatorServiceClient {
	mock := &MockReplicatorServiceClient{ctrl: ctrl}
	mock.recorder = &MockReplicatorServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReplicatorServiceClient) EXPECT() *MockReplicatorServiceClientMockRecorder {
	return m.recorder
}

// Replicate mocks base method.
func (m *MockReplicatorServiceClient) Replicate(arg0 context.Context, arg1 ...grpc.CallOption) (snpb.ReplicatorService_ReplicateClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Replicate", varargs...)
	ret0, _ := ret[0].(snpb.ReplicatorService_ReplicateClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Replicate indicates an expected call of Replicate.
func (mr *MockReplicatorServiceClientMockRecorder) Replicate(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replicate", reflect.TypeOf((*MockReplicatorServiceClient)(nil).Replicate), varargs...)
}

// SyncReplicate mocks base method.
func (m *MockReplicatorServiceClient) SyncReplicate(arg0 context.Context, arg1 *snpb.SyncReplicateRequest, arg2 ...grpc.CallOption) (*snpb.SyncReplicateResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SyncReplicate", varargs...)
	ret0, _ := ret[0].(*snpb.SyncReplicateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SyncReplicate indicates an expected call of SyncReplicate.
func (mr *MockReplicatorServiceClientMockRecorder) SyncReplicate(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncReplicate", reflect.TypeOf((*MockReplicatorServiceClient)(nil).SyncReplicate), varargs...)
}

// MockReplicatorServiceServer is a mock of ReplicatorServiceServer interface.
type MockReplicatorServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockReplicatorServiceServerMockRecorder
}

// MockReplicatorServiceServerMockRecorder is the mock recorder for MockReplicatorServiceServer.
type MockReplicatorServiceServerMockRecorder struct {
	mock *MockReplicatorServiceServer
}

// NewMockReplicatorServiceServer creates a new mock instance.
func NewMockReplicatorServiceServer(ctrl *gomock.Controller) *MockReplicatorServiceServer {
	mock := &MockReplicatorServiceServer{ctrl: ctrl}
	mock.recorder = &MockReplicatorServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReplicatorServiceServer) EXPECT() *MockReplicatorServiceServerMockRecorder {
	return m.recorder
}

// Replicate mocks base method.
func (m *MockReplicatorServiceServer) Replicate(arg0 snpb.ReplicatorService_ReplicateServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replicate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Replicate indicates an expected call of Replicate.
func (mr *MockReplicatorServiceServerMockRecorder) Replicate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replicate", reflect.TypeOf((*MockReplicatorServiceServer)(nil).Replicate), arg0)
}

// SyncReplicate mocks base method.
func (m *MockReplicatorServiceServer) SyncReplicate(arg0 context.Context, arg1 *snpb.SyncReplicateRequest) (*snpb.SyncReplicateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncReplicate", arg0, arg1)
	ret0, _ := ret[0].(*snpb.SyncReplicateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SyncReplicate indicates an expected call of SyncReplicate.
func (mr *MockReplicatorServiceServerMockRecorder) SyncReplicate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncReplicate", reflect.TypeOf((*MockReplicatorServiceServer)(nil).SyncReplicate), arg0, arg1)
}

// MockReplicatorService_ReplicateClient is a mock of ReplicatorService_ReplicateClient interface.
type MockReplicatorService_ReplicateClient struct {
	ctrl     *gomock.Controller
	recorder *MockReplicatorService_ReplicateClientMockRecorder
}

// MockReplicatorService_ReplicateClientMockRecorder is the mock recorder for MockReplicatorService_ReplicateClient.
type MockReplicatorService_ReplicateClientMockRecorder struct {
	mock *MockReplicatorService_ReplicateClient
}

// NewMockReplicatorService_ReplicateClient creates a new mock instance.
func NewMockReplicatorService_ReplicateClient(ctrl *gomock.Controller) *MockReplicatorService_ReplicateClient {
	mock := &MockReplicatorService_ReplicateClient{ctrl: ctrl}
	mock.recorder = &MockReplicatorService_ReplicateClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReplicatorService_ReplicateClient) EXPECT() *MockReplicatorService_ReplicateClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockReplicatorService_ReplicateClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockReplicatorService_ReplicateClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockReplicatorService_ReplicateClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockReplicatorService_ReplicateClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockReplicatorService_ReplicateClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockReplicatorService_ReplicateClient)(nil).Context))
}

// Header mocks base method.
func (m *MockReplicatorService_ReplicateClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockReplicatorService_ReplicateClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockReplicatorService_ReplicateClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockReplicatorService_ReplicateClient) Recv() (*snpb.ReplicationResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*snpb.ReplicationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockReplicatorService_ReplicateClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockReplicatorService_ReplicateClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockReplicatorService_ReplicateClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockReplicatorService_ReplicateClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockReplicatorService_ReplicateClient)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockReplicatorService_ReplicateClient) Send(arg0 *snpb.ReplicationRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockReplicatorService_ReplicateClientMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockReplicatorService_ReplicateClient)(nil).Send), arg0)
}

// SendMsg mocks base method.
func (m *MockReplicatorService_ReplicateClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockReplicatorService_ReplicateClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockReplicatorService_ReplicateClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockReplicatorService_ReplicateClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockReplicatorService_ReplicateClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockReplicatorService_ReplicateClient)(nil).Trailer))
}

// MockReplicatorService_ReplicateServer is a mock of ReplicatorService_ReplicateServer interface.
type MockReplicatorService_ReplicateServer struct {
	ctrl     *gomock.Controller
	recorder *MockReplicatorService_ReplicateServerMockRecorder
}

// MockReplicatorService_ReplicateServerMockRecorder is the mock recorder for MockReplicatorService_ReplicateServer.
type MockReplicatorService_ReplicateServerMockRecorder struct {
	mock *MockReplicatorService_ReplicateServer
}

// NewMockReplicatorService_ReplicateServer creates a new mock instance.
func NewMockReplicatorService_ReplicateServer(ctrl *gomock.Controller) *MockReplicatorService_ReplicateServer {
	mock := &MockReplicatorService_ReplicateServer{ctrl: ctrl}
	mock.recorder = &MockReplicatorService_ReplicateServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReplicatorService_ReplicateServer) EXPECT() *MockReplicatorService_ReplicateServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockReplicatorService_ReplicateServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockReplicatorService_ReplicateServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockReplicatorService_ReplicateServer)(nil).Context))
}

// Recv mocks base method.
func (m *MockReplicatorService_ReplicateServer) Recv() (*snpb.ReplicationRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*snpb.ReplicationRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockReplicatorService_ReplicateServerMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockReplicatorService_ReplicateServer)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockReplicatorService_ReplicateServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockReplicatorService_ReplicateServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockReplicatorService_ReplicateServer)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockReplicatorService_ReplicateServer) Send(arg0 *snpb.ReplicationResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockReplicatorService_ReplicateServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockReplicatorService_ReplicateServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockReplicatorService_ReplicateServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockReplicatorService_ReplicateServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockReplicatorService_ReplicateServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m *MockReplicatorService_ReplicateServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockReplicatorService_ReplicateServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockReplicatorService_ReplicateServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method.
func (m *MockReplicatorService_ReplicateServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockReplicatorService_ReplicateServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockReplicatorService_ReplicateServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockReplicatorService_ReplicateServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockReplicatorService_ReplicateServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockReplicatorService_ReplicateServer)(nil).SetTrailer), arg0)
}
