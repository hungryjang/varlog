// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/replicator_client.go

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/kakao/varlog/pkg/varlog/types"
)

// MockReplicatorClient is a mock of ReplicatorClient interface.
type MockReplicatorClient struct {
	ctrl     *gomock.Controller
	recorder *MockReplicatorClientMockRecorder
}

// MockReplicatorClientMockRecorder is the mock recorder for MockReplicatorClient.
type MockReplicatorClientMockRecorder struct {
	mock *MockReplicatorClient
}

// NewMockReplicatorClient creates a new mock instance.
func NewMockReplicatorClient(ctrl *gomock.Controller) *MockReplicatorClient {
	mock := &MockReplicatorClient{ctrl: ctrl}
	mock.recorder = &MockReplicatorClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReplicatorClient) EXPECT() *MockReplicatorClientMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockReplicatorClient) Run(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockReplicatorClientMockRecorder) Run(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockReplicatorClient)(nil).Run), ctx)
}

// Close mocks base method.
func (m *MockReplicatorClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockReplicatorClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockReplicatorClient)(nil).Close))
}

// Replicate mocks base method.
func (m *MockReplicatorClient) Replicate(ctx context.Context, llsn types.LLSN, data []byte) <-chan error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replicate", ctx, llsn, data)
	ret0, _ := ret[0].(<-chan error)
	return ret0
}

// Replicate indicates an expected call of Replicate.
func (mr *MockReplicatorClientMockRecorder) Replicate(ctx, llsn, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replicate", reflect.TypeOf((*MockReplicatorClient)(nil).Replicate), ctx, llsn, data)
}

// PeerStorageNodeID mocks base method.
func (m *MockReplicatorClient) PeerStorageNodeID() types.StorageNodeID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeerStorageNodeID")
	ret0, _ := ret[0].(types.StorageNodeID)
	return ret0
}

// PeerStorageNodeID indicates an expected call of PeerStorageNodeID.
func (mr *MockReplicatorClientMockRecorder) PeerStorageNodeID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeerStorageNodeID", reflect.TypeOf((*MockReplicatorClient)(nil).PeerStorageNodeID))
}
