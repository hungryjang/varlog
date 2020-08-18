// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kakao/varlog/proto/storage_node (interfaces: LogIOClient,LogIOServer,LogIO_SubscribeClient,LogIO_SubscribeServer)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage_node "github.com/kakao/varlog/proto/storage_node"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockLogIOClient is a mock of LogIOClient interface.
type MockLogIOClient struct {
	ctrl     *gomock.Controller
	recorder *MockLogIOClientMockRecorder
}

// MockLogIOClientMockRecorder is the mock recorder for MockLogIOClient.
type MockLogIOClientMockRecorder struct {
	mock *MockLogIOClient
}

// NewMockLogIOClient creates a new mock instance.
func NewMockLogIOClient(ctrl *gomock.Controller) *MockLogIOClient {
	mock := &MockLogIOClient{ctrl: ctrl}
	mock.recorder = &MockLogIOClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogIOClient) EXPECT() *MockLogIOClientMockRecorder {
	return m.recorder
}

// Append mocks base method.
func (m *MockLogIOClient) Append(arg0 context.Context, arg1 *storage_node.AppendRequest, arg2 ...grpc.CallOption) (*storage_node.AppendResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Append", varargs...)
	ret0, _ := ret[0].(*storage_node.AppendResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Append indicates an expected call of Append.
func (mr *MockLogIOClientMockRecorder) Append(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Append", reflect.TypeOf((*MockLogIOClient)(nil).Append), varargs...)
}

// Read mocks base method.
func (m *MockLogIOClient) Read(arg0 context.Context, arg1 *storage_node.ReadRequest, arg2 ...grpc.CallOption) (*storage_node.ReadResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Read", varargs...)
	ret0, _ := ret[0].(*storage_node.ReadResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockLogIOClientMockRecorder) Read(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockLogIOClient)(nil).Read), varargs...)
}

// Subscribe mocks base method.
func (m *MockLogIOClient) Subscribe(arg0 context.Context, arg1 *storage_node.SubscribeRequest, arg2 ...grpc.CallOption) (storage_node.LogIO_SubscribeClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Subscribe", varargs...)
	ret0, _ := ret[0].(storage_node.LogIO_SubscribeClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockLogIOClientMockRecorder) Subscribe(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockLogIOClient)(nil).Subscribe), varargs...)
}

// Trim mocks base method.
func (m *MockLogIOClient) Trim(arg0 context.Context, arg1 *storage_node.TrimRequest, arg2 ...grpc.CallOption) (*storage_node.TrimResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Trim", varargs...)
	ret0, _ := ret[0].(*storage_node.TrimResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Trim indicates an expected call of Trim.
func (mr *MockLogIOClientMockRecorder) Trim(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trim", reflect.TypeOf((*MockLogIOClient)(nil).Trim), varargs...)
}

// MockLogIOServer is a mock of LogIOServer interface.
type MockLogIOServer struct {
	ctrl     *gomock.Controller
	recorder *MockLogIOServerMockRecorder
}

// MockLogIOServerMockRecorder is the mock recorder for MockLogIOServer.
type MockLogIOServerMockRecorder struct {
	mock *MockLogIOServer
}

// NewMockLogIOServer creates a new mock instance.
func NewMockLogIOServer(ctrl *gomock.Controller) *MockLogIOServer {
	mock := &MockLogIOServer{ctrl: ctrl}
	mock.recorder = &MockLogIOServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogIOServer) EXPECT() *MockLogIOServerMockRecorder {
	return m.recorder
}

// Append mocks base method.
func (m *MockLogIOServer) Append(arg0 context.Context, arg1 *storage_node.AppendRequest) (*storage_node.AppendResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Append", arg0, arg1)
	ret0, _ := ret[0].(*storage_node.AppendResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Append indicates an expected call of Append.
func (mr *MockLogIOServerMockRecorder) Append(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Append", reflect.TypeOf((*MockLogIOServer)(nil).Append), arg0, arg1)
}

// Read mocks base method.
func (m *MockLogIOServer) Read(arg0 context.Context, arg1 *storage_node.ReadRequest) (*storage_node.ReadResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0, arg1)
	ret0, _ := ret[0].(*storage_node.ReadResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockLogIOServerMockRecorder) Read(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockLogIOServer)(nil).Read), arg0, arg1)
}

// Subscribe mocks base method.
func (m *MockLogIOServer) Subscribe(arg0 *storage_node.SubscribeRequest, arg1 storage_node.LogIO_SubscribeServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockLogIOServerMockRecorder) Subscribe(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockLogIOServer)(nil).Subscribe), arg0, arg1)
}

// Trim mocks base method.
func (m *MockLogIOServer) Trim(arg0 context.Context, arg1 *storage_node.TrimRequest) (*storage_node.TrimResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trim", arg0, arg1)
	ret0, _ := ret[0].(*storage_node.TrimResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Trim indicates an expected call of Trim.
func (mr *MockLogIOServerMockRecorder) Trim(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trim", reflect.TypeOf((*MockLogIOServer)(nil).Trim), arg0, arg1)
}

// MockLogIO_SubscribeClient is a mock of LogIO_SubscribeClient interface.
type MockLogIO_SubscribeClient struct {
	ctrl     *gomock.Controller
	recorder *MockLogIO_SubscribeClientMockRecorder
}

// MockLogIO_SubscribeClientMockRecorder is the mock recorder for MockLogIO_SubscribeClient.
type MockLogIO_SubscribeClientMockRecorder struct {
	mock *MockLogIO_SubscribeClient
}

// NewMockLogIO_SubscribeClient creates a new mock instance.
func NewMockLogIO_SubscribeClient(ctrl *gomock.Controller) *MockLogIO_SubscribeClient {
	mock := &MockLogIO_SubscribeClient{ctrl: ctrl}
	mock.recorder = &MockLogIO_SubscribeClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogIO_SubscribeClient) EXPECT() *MockLogIO_SubscribeClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockLogIO_SubscribeClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockLogIO_SubscribeClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockLogIO_SubscribeClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockLogIO_SubscribeClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockLogIO_SubscribeClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockLogIO_SubscribeClient)(nil).Context))
}

// Header mocks base method.
func (m *MockLogIO_SubscribeClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockLogIO_SubscribeClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockLogIO_SubscribeClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockLogIO_SubscribeClient) Recv() (*storage_node.SubscribeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*storage_node.SubscribeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockLogIO_SubscribeClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockLogIO_SubscribeClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockLogIO_SubscribeClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockLogIO_SubscribeClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockLogIO_SubscribeClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method.
func (m *MockLogIO_SubscribeClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockLogIO_SubscribeClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockLogIO_SubscribeClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockLogIO_SubscribeClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockLogIO_SubscribeClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockLogIO_SubscribeClient)(nil).Trailer))
}

// MockLogIO_SubscribeServer is a mock of LogIO_SubscribeServer interface.
type MockLogIO_SubscribeServer struct {
	ctrl     *gomock.Controller
	recorder *MockLogIO_SubscribeServerMockRecorder
}

// MockLogIO_SubscribeServerMockRecorder is the mock recorder for MockLogIO_SubscribeServer.
type MockLogIO_SubscribeServerMockRecorder struct {
	mock *MockLogIO_SubscribeServer
}

// NewMockLogIO_SubscribeServer creates a new mock instance.
func NewMockLogIO_SubscribeServer(ctrl *gomock.Controller) *MockLogIO_SubscribeServer {
	mock := &MockLogIO_SubscribeServer{ctrl: ctrl}
	mock.recorder = &MockLogIO_SubscribeServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogIO_SubscribeServer) EXPECT() *MockLogIO_SubscribeServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockLogIO_SubscribeServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockLogIO_SubscribeServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockLogIO_SubscribeServer)(nil).Context))
}

// RecvMsg mocks base method.
func (m *MockLogIO_SubscribeServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockLogIO_SubscribeServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockLogIO_SubscribeServer)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockLogIO_SubscribeServer) Send(arg0 *storage_node.SubscribeResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockLogIO_SubscribeServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockLogIO_SubscribeServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockLogIO_SubscribeServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockLogIO_SubscribeServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockLogIO_SubscribeServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m *MockLogIO_SubscribeServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockLogIO_SubscribeServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockLogIO_SubscribeServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method.
func (m *MockLogIO_SubscribeServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockLogIO_SubscribeServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockLogIO_SubscribeServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockLogIO_SubscribeServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockLogIO_SubscribeServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockLogIO_SubscribeServer)(nil).SetTrailer), arg0)
}
