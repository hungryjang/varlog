// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kakao/varlog/internal/storagenode (interfaces: ReplicatorClient)

// Package storagenode is a generated GoMock package.
package storagenode

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	types "github.com/kakao/varlog/pkg/types"
	snpb "github.com/kakao/varlog/proto/snpb"
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

// Replicate mocks base method.
func (m *MockReplicatorClient) Replicate(arg0 context.Context, arg1 types.LLSN, arg2 []byte) <-chan error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replicate", arg0, arg1, arg2)
	ret0, _ := ret[0].(<-chan error)
	return ret0
}

// Replicate indicates an expected call of Replicate.
func (mr *MockReplicatorClientMockRecorder) Replicate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replicate", reflect.TypeOf((*MockReplicatorClient)(nil).Replicate), arg0, arg1, arg2)
}

// Run mocks base method.
func (m *MockReplicatorClient) Run(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockReplicatorClientMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockReplicatorClient)(nil).Run), arg0)
}

// SyncReplicate mocks base method.
func (m *MockReplicatorClient) SyncReplicate(arg0 context.Context, arg1 types.LogStreamID, arg2, arg3, arg4 snpb.SyncPosition, arg5 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncReplicate", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncReplicate indicates an expected call of SyncReplicate.
func (mr *MockReplicatorClientMockRecorder) SyncReplicate(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncReplicate", reflect.TypeOf((*MockReplicatorClient)(nil).SyncReplicate), arg0, arg1, arg2, arg3, arg4, arg5)
}
