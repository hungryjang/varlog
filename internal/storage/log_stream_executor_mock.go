// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/log_stream_executor.go

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	varlog "github.com/kakao/varlog/pkg/varlog"
	types "github.com/kakao/varlog/pkg/varlog/types"
	snpb "github.com/kakao/varlog/proto/snpb"
	varlogpb "github.com/kakao/varlog/proto/varlogpb"
)

// MockSealer is a mock of Sealer interface.
type MockSealer struct {
	ctrl     *gomock.Controller
	recorder *MockSealerMockRecorder
}

// MockSealerMockRecorder is the mock recorder for MockSealer.
type MockSealerMockRecorder struct {
	mock *MockSealer
}

// NewMockSealer creates a new mock instance.
func NewMockSealer(ctrl *gomock.Controller) *MockSealer {
	mock := &MockSealer{ctrl: ctrl}
	mock.recorder = &MockSealerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSealer) EXPECT() *MockSealerMockRecorder {
	return m.recorder
}

// Seal mocks base method.
func (m *MockSealer) Seal(lastCommittedGLSN types.GLSN) (varlogpb.LogStreamStatus, types.GLSN) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seal", lastCommittedGLSN)
	ret0, _ := ret[0].(varlogpb.LogStreamStatus)
	ret1, _ := ret[1].(types.GLSN)
	return ret0, ret1
}

// Seal indicates an expected call of Seal.
func (mr *MockSealerMockRecorder) Seal(lastCommittedGLSN interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seal", reflect.TypeOf((*MockSealer)(nil).Seal), lastCommittedGLSN)
}

// MockUnsealer is a mock of Unsealer interface.
type MockUnsealer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsealerMockRecorder
}

// MockUnsealerMockRecorder is the mock recorder for MockUnsealer.
type MockUnsealerMockRecorder struct {
	mock *MockUnsealer
}

// NewMockUnsealer creates a new mock instance.
func NewMockUnsealer(ctrl *gomock.Controller) *MockUnsealer {
	mock := &MockUnsealer{ctrl: ctrl}
	mock.recorder = &MockUnsealerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsealer) EXPECT() *MockUnsealerMockRecorder {
	return m.recorder
}

// Unseal mocks base method.
func (m *MockUnsealer) Unseal() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unseal")
	ret0, _ := ret[0].(error)
	return ret0
}

// Unseal indicates an expected call of Unseal.
func (mr *MockUnsealerMockRecorder) Unseal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unseal", reflect.TypeOf((*MockUnsealer)(nil).Unseal))
}

// MockSyncer is a mock of Syncer interface.
type MockSyncer struct {
	ctrl     *gomock.Controller
	recorder *MockSyncerMockRecorder
}

// MockSyncerMockRecorder is the mock recorder for MockSyncer.
type MockSyncerMockRecorder struct {
	mock *MockSyncer
}

// NewMockSyncer creates a new mock instance.
func NewMockSyncer(ctrl *gomock.Controller) *MockSyncer {
	mock := &MockSyncer{ctrl: ctrl}
	mock.recorder = &MockSyncerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSyncer) EXPECT() *MockSyncerMockRecorder {
	return m.recorder
}

// Sync mocks base method.
func (m *MockSyncer) Sync(ctx context.Context, replica Replica, lastGLSN types.GLSN) (*SyncTaskStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", ctx, replica, lastGLSN)
	ret0, _ := ret[0].(*SyncTaskStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sync indicates an expected call of Sync.
func (mr *MockSyncerMockRecorder) Sync(ctx, replica, lastGLSN interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockSyncer)(nil).Sync), ctx, replica, lastGLSN)
}

// SyncReplicate mocks base method.
func (m *MockSyncer) SyncReplicate(ctx context.Context, first, last, current snpb.SyncPosition, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncReplicate", ctx, first, last, current, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncReplicate indicates an expected call of SyncReplicate.
func (mr *MockSyncerMockRecorder) SyncReplicate(ctx, first, last, current, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncReplicate", reflect.TypeOf((*MockSyncer)(nil).SyncReplicate), ctx, first, last, current, data)
}

// MockLogStreamExecutor is a mock of LogStreamExecutor interface.
type MockLogStreamExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockLogStreamExecutorMockRecorder
}

// MockLogStreamExecutorMockRecorder is the mock recorder for MockLogStreamExecutor.
type MockLogStreamExecutorMockRecorder struct {
	mock *MockLogStreamExecutor
}

// NewMockLogStreamExecutor creates a new mock instance.
func NewMockLogStreamExecutor(ctrl *gomock.Controller) *MockLogStreamExecutor {
	mock := &MockLogStreamExecutor{ctrl: ctrl}
	mock.recorder = &MockLogStreamExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogStreamExecutor) EXPECT() *MockLogStreamExecutorMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockLogStreamExecutor) Run(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockLogStreamExecutorMockRecorder) Run(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockLogStreamExecutor)(nil).Run), ctx)
}

// Close mocks base method.
func (m *MockLogStreamExecutor) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockLogStreamExecutorMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockLogStreamExecutor)(nil).Close))
}

// LogStreamID mocks base method.
func (m *MockLogStreamExecutor) LogStreamID() types.LogStreamID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogStreamID")
	ret0, _ := ret[0].(types.LogStreamID)
	return ret0
}

// LogStreamID indicates an expected call of LogStreamID.
func (mr *MockLogStreamExecutorMockRecorder) LogStreamID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogStreamID", reflect.TypeOf((*MockLogStreamExecutor)(nil).LogStreamID))
}

// Status mocks base method.
func (m *MockLogStreamExecutor) Status() varlogpb.LogStreamStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(varlogpb.LogStreamStatus)
	return ret0
}

// Status indicates an expected call of Status.
func (mr *MockLogStreamExecutorMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockLogStreamExecutor)(nil).Status))
}

// Read mocks base method.
func (m *MockLogStreamExecutor) Read(ctx context.Context, glsn types.GLSN) (varlog.LogEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", ctx, glsn)
	ret0, _ := ret[0].(varlog.LogEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockLogStreamExecutorMockRecorder) Read(ctx, glsn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockLogStreamExecutor)(nil).Read), ctx, glsn)
}

// Subscribe mocks base method.
func (m *MockLogStreamExecutor) Subscribe(ctx context.Context, begin, end types.GLSN) (<-chan ScanResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx, begin, end)
	ret0, _ := ret[0].(<-chan ScanResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockLogStreamExecutorMockRecorder) Subscribe(ctx, begin, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockLogStreamExecutor)(nil).Subscribe), ctx, begin, end)
}

// Append mocks base method.
func (m *MockLogStreamExecutor) Append(ctx context.Context, data []byte, backups ...Replica) (types.GLSN, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, data}
	for _, a := range backups {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Append", varargs...)
	ret0, _ := ret[0].(types.GLSN)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Append indicates an expected call of Append.
func (mr *MockLogStreamExecutorMockRecorder) Append(ctx, data interface{}, backups ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, data}, backups...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Append", reflect.TypeOf((*MockLogStreamExecutor)(nil).Append), varargs...)
}

// Trim mocks base method.
func (m *MockLogStreamExecutor) Trim(ctx context.Context, glsn types.GLSN) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trim", ctx, glsn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Trim indicates an expected call of Trim.
func (mr *MockLogStreamExecutorMockRecorder) Trim(ctx, glsn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trim", reflect.TypeOf((*MockLogStreamExecutor)(nil).Trim), ctx, glsn)
}

// Replicate mocks base method.
func (m *MockLogStreamExecutor) Replicate(ctx context.Context, llsn types.LLSN, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replicate", ctx, llsn, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Replicate indicates an expected call of Replicate.
func (mr *MockLogStreamExecutorMockRecorder) Replicate(ctx, llsn, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replicate", reflect.TypeOf((*MockLogStreamExecutor)(nil).Replicate), ctx, llsn, data)
}

// GetReport mocks base method.
func (m *MockLogStreamExecutor) GetReport() UncommittedLogStreamStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReport")
	ret0, _ := ret[0].(UncommittedLogStreamStatus)
	return ret0
}

// GetReport indicates an expected call of GetReport.
func (mr *MockLogStreamExecutorMockRecorder) GetReport() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReport", reflect.TypeOf((*MockLogStreamExecutor)(nil).GetReport))
}

// Commit mocks base method.
func (m *MockLogStreamExecutor) Commit(ctx context.Context, commitResult CommittedLogStreamStatus) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Commit", ctx, commitResult)
}

// Commit indicates an expected call of Commit.
func (mr *MockLogStreamExecutorMockRecorder) Commit(ctx, commitResult interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockLogStreamExecutor)(nil).Commit), ctx, commitResult)
}

// Seal mocks base method.
func (m *MockLogStreamExecutor) Seal(lastCommittedGLSN types.GLSN) (varlogpb.LogStreamStatus, types.GLSN) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seal", lastCommittedGLSN)
	ret0, _ := ret[0].(varlogpb.LogStreamStatus)
	ret1, _ := ret[1].(types.GLSN)
	return ret0, ret1
}

// Seal indicates an expected call of Seal.
func (mr *MockLogStreamExecutorMockRecorder) Seal(lastCommittedGLSN interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seal", reflect.TypeOf((*MockLogStreamExecutor)(nil).Seal), lastCommittedGLSN)
}

// Unseal mocks base method.
func (m *MockLogStreamExecutor) Unseal() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unseal")
	ret0, _ := ret[0].(error)
	return ret0
}

// Unseal indicates an expected call of Unseal.
func (mr *MockLogStreamExecutorMockRecorder) Unseal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unseal", reflect.TypeOf((*MockLogStreamExecutor)(nil).Unseal))
}

// Sync mocks base method.
func (m *MockLogStreamExecutor) Sync(ctx context.Context, replica Replica, lastGLSN types.GLSN) (*SyncTaskStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", ctx, replica, lastGLSN)
	ret0, _ := ret[0].(*SyncTaskStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sync indicates an expected call of Sync.
func (mr *MockLogStreamExecutorMockRecorder) Sync(ctx, replica, lastGLSN interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockLogStreamExecutor)(nil).Sync), ctx, replica, lastGLSN)
}

// SyncReplicate mocks base method.
func (m *MockLogStreamExecutor) SyncReplicate(ctx context.Context, first, last, current snpb.SyncPosition, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncReplicate", ctx, first, last, current, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncReplicate indicates an expected call of SyncReplicate.
func (mr *MockLogStreamExecutorMockRecorder) SyncReplicate(ctx, first, last, current, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncReplicate", reflect.TypeOf((*MockLogStreamExecutor)(nil).SyncReplicate), ctx, first, last, current, data)
}
