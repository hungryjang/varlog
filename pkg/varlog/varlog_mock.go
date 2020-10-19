// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kakao/varlog/pkg/varlog (interfaces: ManagementClient)

// Package varlog is a generated GoMock package.
package varlog

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/kakao/varlog/pkg/varlog/types"
	snpb "github.com/kakao/varlog/proto/snpb"
	varlogpb "github.com/kakao/varlog/proto/varlogpb"
)

// MockManagementClient is a mock of ManagementClient interface.
type MockManagementClient struct {
	ctrl     *gomock.Controller
	recorder *MockManagementClientMockRecorder
}

// MockManagementClientMockRecorder is the mock recorder for MockManagementClient.
type MockManagementClientMockRecorder struct {
	mock *MockManagementClient
}

// NewMockManagementClient creates a new mock instance.
func NewMockManagementClient(ctrl *gomock.Controller) *MockManagementClient {
	mock := &MockManagementClient{ctrl: ctrl}
	mock.recorder = &MockManagementClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManagementClient) EXPECT() *MockManagementClientMockRecorder {
	return m.recorder
}

// AddLogStream mocks base method.
func (m *MockManagementClient) AddLogStream(arg0 context.Context, arg1 types.LogStreamID, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLogStream", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLogStream indicates an expected call of AddLogStream.
func (mr *MockManagementClientMockRecorder) AddLogStream(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLogStream", reflect.TypeOf((*MockManagementClient)(nil).AddLogStream), arg0, arg1, arg2)
}

// Close mocks base method.
func (m *MockManagementClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockManagementClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockManagementClient)(nil).Close))
}

// GetMetadata mocks base method.
func (m *MockManagementClient) GetMetadata(arg0 context.Context, arg1 snpb.MetadataType) (*varlogpb.StorageNodeMetadataDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetadata", arg0, arg1)
	ret0, _ := ret[0].(*varlogpb.StorageNodeMetadataDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetadata indicates an expected call of GetMetadata.
func (mr *MockManagementClientMockRecorder) GetMetadata(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadata", reflect.TypeOf((*MockManagementClient)(nil).GetMetadata), arg0, arg1)
}

// PeerAddress mocks base method.
func (m *MockManagementClient) PeerAddress() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeerAddress")
	ret0, _ := ret[0].(string)
	return ret0
}

// PeerAddress indicates an expected call of PeerAddress.
func (mr *MockManagementClientMockRecorder) PeerAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeerAddress", reflect.TypeOf((*MockManagementClient)(nil).PeerAddress))
}

// PeerStorageNodeID mocks base method.
func (m *MockManagementClient) PeerStorageNodeID() types.StorageNodeID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeerStorageNodeID")
	ret0, _ := ret[0].(types.StorageNodeID)
	return ret0
}

// PeerStorageNodeID indicates an expected call of PeerStorageNodeID.
func (mr *MockManagementClientMockRecorder) PeerStorageNodeID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeerStorageNodeID", reflect.TypeOf((*MockManagementClient)(nil).PeerStorageNodeID))
}

// RemoveLogStream mocks base method.
func (m *MockManagementClient) RemoveLogStream(arg0 context.Context, arg1 types.LogStreamID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveLogStream", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveLogStream indicates an expected call of RemoveLogStream.
func (mr *MockManagementClientMockRecorder) RemoveLogStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLogStream", reflect.TypeOf((*MockManagementClient)(nil).RemoveLogStream), arg0, arg1)
}

// Seal mocks base method.
func (m *MockManagementClient) Seal(arg0 context.Context, arg1 types.LogStreamID, arg2 types.GLSN) (varlogpb.LogStreamStatus, types.GLSN, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seal", arg0, arg1, arg2)
	ret0, _ := ret[0].(varlogpb.LogStreamStatus)
	ret1, _ := ret[1].(types.GLSN)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Seal indicates an expected call of Seal.
func (mr *MockManagementClientMockRecorder) Seal(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seal", reflect.TypeOf((*MockManagementClient)(nil).Seal), arg0, arg1, arg2)
}

// Sync mocks base method.
func (m *MockManagementClient) Sync(arg0 context.Context, arg1 types.LogStreamID, arg2 types.StorageNodeID, arg3 string, arg4 types.GLSN) (snpb.SyncState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(snpb.SyncState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sync indicates an expected call of Sync.
func (mr *MockManagementClientMockRecorder) Sync(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockManagementClient)(nil).Sync), arg0, arg1, arg2, arg3, arg4)
}

// Unseal mocks base method.
func (m *MockManagementClient) Unseal(arg0 context.Context, arg1 types.LogStreamID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unseal", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unseal indicates an expected call of Unseal.
func (mr *MockManagementClientMockRecorder) Unseal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unseal", reflect.TypeOf((*MockManagementClient)(nil).Unseal), arg0, arg1)
}
