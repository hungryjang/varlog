// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/storage.go

// Package storage is a generated GoMock package.
package storage

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	varlog "github.com/kakao/varlog/pkg/varlog"
	types "github.com/kakao/varlog/pkg/varlog/types"
)

// MockScanner is a mock of Scanner interface.
type MockScanner struct {
	ctrl     *gomock.Controller
	recorder *MockScannerMockRecorder
}

// MockScannerMockRecorder is the mock recorder for MockScanner.
type MockScannerMockRecorder struct {
	mock *MockScanner
}

// NewMockScanner creates a new mock instance.
func NewMockScanner(ctrl *gomock.Controller) *MockScanner {
	mock := &MockScanner{ctrl: ctrl}
	mock.recorder = &MockScannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScanner) EXPECT() *MockScannerMockRecorder {
	return m.recorder
}

// Next mocks base method.
func (m *MockScanner) Next() (varlog.LogEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].(varlog.LogEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Next indicates an expected call of Next.
func (mr *MockScannerMockRecorder) Next() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockScanner)(nil).Next))
}

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Read mocks base method.
func (m *MockStorage) Read(glsn types.GLSN) (varlog.LogEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", glsn)
	ret0, _ := ret[0].(varlog.LogEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockStorageMockRecorder) Read(glsn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockStorage)(nil).Read), glsn)
}

// ReadByLLSN mocks base method.
func (m *MockStorage) ReadByLLSN(llsn types.LLSN) (varlog.LogEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadByLLSN", llsn)
	ret0, _ := ret[0].(varlog.LogEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadByLLSN indicates an expected call of ReadByLLSN.
func (mr *MockStorageMockRecorder) ReadByLLSN(llsn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadByLLSN", reflect.TypeOf((*MockStorage)(nil).ReadByLLSN), llsn)
}

// Scan mocks base method.
func (m *MockStorage) Scan(begin, end types.GLSN) (Scanner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scan", begin, end)
	ret0, _ := ret[0].(Scanner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Scan indicates an expected call of Scan.
func (mr *MockStorageMockRecorder) Scan(begin, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockStorage)(nil).Scan), begin, end)
}

// Write mocks base method.
func (m *MockStorage) Write(llsn types.LLSN, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", llsn, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write.
func (mr *MockStorageMockRecorder) Write(llsn, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockStorage)(nil).Write), llsn, data)
}

// Commit mocks base method.
func (m *MockStorage) Commit(llsn types.LLSN, glsn types.GLSN) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", llsn, glsn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockStorageMockRecorder) Commit(llsn, glsn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockStorage)(nil).Commit), llsn, glsn)
}

// DeleteCommitted mocks base method.
func (m *MockStorage) DeleteCommitted(glsn types.GLSN) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCommitted", glsn)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCommitted indicates an expected call of DeleteCommitted.
func (mr *MockStorageMockRecorder) DeleteCommitted(glsn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCommitted", reflect.TypeOf((*MockStorage)(nil).DeleteCommitted), glsn)
}

// DeleteUncommitted mocks base method.
func (m *MockStorage) DeleteUncommitted(llsn types.LLSN) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUncommitted", llsn)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUncommitted indicates an expected call of DeleteUncommitted.
func (mr *MockStorageMockRecorder) DeleteUncommitted(llsn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUncommitted", reflect.TypeOf((*MockStorage)(nil).DeleteUncommitted), llsn)
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}
