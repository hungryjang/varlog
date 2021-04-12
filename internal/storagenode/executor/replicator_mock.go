// Code generated by MockGen. DO NOT EDIT.
// Source: replicator.go

// Package executor is a generated GoMock package.
package executor

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockReplicator is a mock of replicator interface.
type MockReplicator struct {
	ctrl     *gomock.Controller
	recorder *MockReplicatorMockRecorder
}

// MockReplicatorMockRecorder is the mock recorder for MockReplicator.
type MockReplicatorMockRecorder struct {
	mock *MockReplicator
}

// NewMockReplicator creates a new mock instance.
func NewMockReplicator(ctrl *gomock.Controller) *MockReplicator {
	mock := &MockReplicator{ctrl: ctrl}
	mock.recorder = &MockReplicatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReplicator) EXPECT() *MockReplicatorMockRecorder {
	return m.recorder
}

// drainQueue mocks base method.
func (m *MockReplicator) drainQueue() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "drainQueue")
}

// drainQueue indicates an expected call of drainQueue.
func (mr *MockReplicatorMockRecorder) drainQueue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "drainQueue", reflect.TypeOf((*MockReplicator)(nil).drainQueue))
}

// send mocks base method.
func (m *MockReplicator) send(ctx context.Context, t *replicateTask) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "send", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// send indicates an expected call of send.
func (mr *MockReplicatorMockRecorder) send(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "send", reflect.TypeOf((*MockReplicator)(nil).send), ctx, t)
}

// stop mocks base method.
func (m *MockReplicator) stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "stop")
}

// stop indicates an expected call of stop.
func (mr *MockReplicatorMockRecorder) stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "stop", reflect.TypeOf((*MockReplicator)(nil).stop))
}
