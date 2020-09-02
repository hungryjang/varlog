// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/log_stream_reporter.go

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/kakao/varlog/pkg/varlog/types"
)

// MockLogStreamReporter is a mock of LogStreamReporter interface.
type MockLogStreamReporter struct {
	ctrl     *gomock.Controller
	recorder *MockLogStreamReporterMockRecorder
}

// MockLogStreamReporterMockRecorder is the mock recorder for MockLogStreamReporter.
type MockLogStreamReporterMockRecorder struct {
	mock *MockLogStreamReporter
}

// NewMockLogStreamReporter creates a new mock instance.
func NewMockLogStreamReporter(ctrl *gomock.Controller) *MockLogStreamReporter {
	mock := &MockLogStreamReporter{ctrl: ctrl}
	mock.recorder = &MockLogStreamReporterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogStreamReporter) EXPECT() *MockLogStreamReporterMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockLogStreamReporter) Run(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", ctx)
}

// Run indicates an expected call of Run.
func (mr *MockLogStreamReporterMockRecorder) Run(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockLogStreamReporter)(nil).Run), ctx)
}

// Close mocks base method.
func (m *MockLogStreamReporter) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockLogStreamReporterMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockLogStreamReporter)(nil).Close))
}

// StorageNodeID mocks base method.
func (m *MockLogStreamReporter) StorageNodeID() types.StorageNodeID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorageNodeID")
	ret0, _ := ret[0].(types.StorageNodeID)
	return ret0
}

// StorageNodeID indicates an expected call of StorageNodeID.
func (mr *MockLogStreamReporterMockRecorder) StorageNodeID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorageNodeID", reflect.TypeOf((*MockLogStreamReporter)(nil).StorageNodeID))
}

// GetReport mocks base method.
func (m *MockLogStreamReporter) GetReport(ctx context.Context) (types.GLSN, []UncommittedLogStreamStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReport", ctx)
	ret0, _ := ret[0].(types.GLSN)
	ret1, _ := ret[1].([]UncommittedLogStreamStatus)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetReport indicates an expected call of GetReport.
func (mr *MockLogStreamReporterMockRecorder) GetReport(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReport", reflect.TypeOf((*MockLogStreamReporter)(nil).GetReport), ctx)
}

// Commit mocks base method.
func (m *MockLogStreamReporter) Commit(ctx context.Context, highWatermark, prevHighWatermark types.GLSN, commitResults []CommittedLogStreamStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", ctx, highWatermark, prevHighWatermark, commitResults)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockLogStreamReporterMockRecorder) Commit(ctx, highWatermark, prevHighWatermark, commitResults interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockLogStreamReporter)(nil).Commit), ctx, highWatermark, prevHighWatermark, commitResults)
}
