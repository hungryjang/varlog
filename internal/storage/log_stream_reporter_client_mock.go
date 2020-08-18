// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/log_stream_reporter_client.go

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	storage_node "github.com/kakao/varlog/proto/storage_node"
	reflect "reflect"
)

// MockLogStreamReporterClient is a mock of LogStreamReporterClient interface
type MockLogStreamReporterClient struct {
	ctrl     *gomock.Controller
	recorder *MockLogStreamReporterClientMockRecorder
}

// MockLogStreamReporterClientMockRecorder is the mock recorder for MockLogStreamReporterClient
type MockLogStreamReporterClientMockRecorder struct {
	mock *MockLogStreamReporterClient
}

// NewMockLogStreamReporterClient creates a new mock instance
func NewMockLogStreamReporterClient(ctrl *gomock.Controller) *MockLogStreamReporterClient {
	mock := &MockLogStreamReporterClient{ctrl: ctrl}
	mock.recorder = &MockLogStreamReporterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogStreamReporterClient) EXPECT() *MockLogStreamReporterClientMockRecorder {
	return m.recorder
}

// GetReport mocks base method
func (m *MockLogStreamReporterClient) GetReport(arg0 context.Context) (*storage_node.LocalLogStreamDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReport", arg0)
	ret0, _ := ret[0].(*storage_node.LocalLogStreamDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReport indicates an expected call of GetReport
func (mr *MockLogStreamReporterClientMockRecorder) GetReport(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReport", reflect.TypeOf((*MockLogStreamReporterClient)(nil).GetReport), arg0)
}

// Commit mocks base method
func (m *MockLogStreamReporterClient) Commit(arg0 context.Context, arg1 *storage_node.GlobalLogStreamDescriptor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit
func (mr *MockLogStreamReporterClientMockRecorder) Commit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockLogStreamReporterClient)(nil).Commit), arg0, arg1)
}

// Close mocks base method
func (m *MockLogStreamReporterClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockLogStreamReporterClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockLogStreamReporterClient)(nil).Close))
}