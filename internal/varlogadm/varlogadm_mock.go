// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kakao/varlog/internal/varlogadm (interfaces: ClusterMetadataView,StorageNodeManager,MetadataRepositoryManager,StorageNodeWatcher,StatRepository)

// Package varlogadm is a generated GoMock package.
package varlogadm

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	snc "github.com/kakao/varlog/pkg/snc"
	types "github.com/kakao/varlog/pkg/types"
	mrpb "github.com/kakao/varlog/proto/mrpb"
	snpb "github.com/kakao/varlog/proto/snpb"
	varlogpb "github.com/kakao/varlog/proto/varlogpb"
	vmspb "github.com/kakao/varlog/proto/vmspb"
)

// MockClusterMetadataView is a mock of ClusterMetadataView interface.
type MockClusterMetadataView struct {
	ctrl     *gomock.Controller
	recorder *MockClusterMetadataViewMockRecorder
}

// MockClusterMetadataViewMockRecorder is the mock recorder for MockClusterMetadataView.
type MockClusterMetadataViewMockRecorder struct {
	mock *MockClusterMetadataView
}

// NewMockClusterMetadataView creates a new mock instance.
func NewMockClusterMetadataView(ctrl *gomock.Controller) *MockClusterMetadataView {
	mock := &MockClusterMetadataView{ctrl: ctrl}
	mock.recorder = &MockClusterMetadataViewMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterMetadataView) EXPECT() *MockClusterMetadataViewMockRecorder {
	return m.recorder
}

// ClusterMetadata mocks base method.
func (m *MockClusterMetadataView) ClusterMetadata(arg0 context.Context) (*varlogpb.MetadataDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClusterMetadata", arg0)
	ret0, _ := ret[0].(*varlogpb.MetadataDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClusterMetadata indicates an expected call of ClusterMetadata.
func (mr *MockClusterMetadataViewMockRecorder) ClusterMetadata(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClusterMetadata", reflect.TypeOf((*MockClusterMetadataView)(nil).ClusterMetadata), arg0)
}

// StorageNode mocks base method.
func (m *MockClusterMetadataView) StorageNode(arg0 context.Context, arg1 types.StorageNodeID) (*varlogpb.StorageNodeDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorageNode", arg0, arg1)
	ret0, _ := ret[0].(*varlogpb.StorageNodeDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorageNode indicates an expected call of StorageNode.
func (mr *MockClusterMetadataViewMockRecorder) StorageNode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorageNode", reflect.TypeOf((*MockClusterMetadataView)(nil).StorageNode), arg0, arg1)
}

// MockStorageNodeManager is a mock of StorageNodeManager interface.
type MockStorageNodeManager struct {
	ctrl     *gomock.Controller
	recorder *MockStorageNodeManagerMockRecorder
}

// MockStorageNodeManagerMockRecorder is the mock recorder for MockStorageNodeManager.
type MockStorageNodeManagerMockRecorder struct {
	mock *MockStorageNodeManager
}

// NewMockStorageNodeManager creates a new mock instance.
func NewMockStorageNodeManager(ctrl *gomock.Controller) *MockStorageNodeManager {
	mock := &MockStorageNodeManager{ctrl: ctrl}
	mock.recorder = &MockStorageNodeManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageNodeManager) EXPECT() *MockStorageNodeManagerMockRecorder {
	return m.recorder
}

// AddLogStream mocks base method.
func (m *MockStorageNodeManager) AddLogStream(arg0 context.Context, arg1 *varlogpb.LogStreamDescriptor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLogStream", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLogStream indicates an expected call of AddLogStream.
func (mr *MockStorageNodeManagerMockRecorder) AddLogStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLogStream", reflect.TypeOf((*MockStorageNodeManager)(nil).AddLogStream), arg0, arg1)
}

// AddLogStreamReplica mocks base method.
func (m *MockStorageNodeManager) AddLogStreamReplica(arg0 context.Context, arg1 types.StorageNodeID, arg2 types.TopicID, arg3 types.LogStreamID, arg4 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLogStreamReplica", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLogStreamReplica indicates an expected call of AddLogStreamReplica.
func (mr *MockStorageNodeManagerMockRecorder) AddLogStreamReplica(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLogStreamReplica", reflect.TypeOf((*MockStorageNodeManager)(nil).AddLogStreamReplica), arg0, arg1, arg2, arg3, arg4)
}

// AddStorageNode mocks base method.
func (m *MockStorageNodeManager) AddStorageNode(arg0 snc.StorageNodeManagementClient) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddStorageNode", arg0)
}

// AddStorageNode indicates an expected call of AddStorageNode.
func (mr *MockStorageNodeManagerMockRecorder) AddStorageNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddStorageNode", reflect.TypeOf((*MockStorageNodeManager)(nil).AddStorageNode), arg0)
}

// Close mocks base method.
func (m *MockStorageNodeManager) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageNodeManagerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorageNodeManager)(nil).Close))
}

// Contains mocks base method.
func (m *MockStorageNodeManager) Contains(arg0 types.StorageNodeID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Contains", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Contains indicates an expected call of Contains.
func (mr *MockStorageNodeManagerMockRecorder) Contains(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Contains", reflect.TypeOf((*MockStorageNodeManager)(nil).Contains), arg0)
}

// ContainsAddress mocks base method.
func (m *MockStorageNodeManager) ContainsAddress(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainsAddress", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ContainsAddress indicates an expected call of ContainsAddress.
func (mr *MockStorageNodeManagerMockRecorder) ContainsAddress(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainsAddress", reflect.TypeOf((*MockStorageNodeManager)(nil).ContainsAddress), arg0)
}

// GetMetadata mocks base method.
func (m *MockStorageNodeManager) GetMetadata(arg0 context.Context, arg1 types.StorageNodeID) (*snpb.StorageNodeMetadataDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetadata", arg0, arg1)
	ret0, _ := ret[0].(*snpb.StorageNodeMetadataDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetadata indicates an expected call of GetMetadata.
func (mr *MockStorageNodeManagerMockRecorder) GetMetadata(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadata", reflect.TypeOf((*MockStorageNodeManager)(nil).GetMetadata), arg0, arg1)
}

// GetMetadataByAddr mocks base method.
func (m *MockStorageNodeManager) GetMetadataByAddr(arg0 context.Context, arg1 string) (snc.StorageNodeManagementClient, *snpb.StorageNodeMetadataDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetadataByAddr", arg0, arg1)
	ret0, _ := ret[0].(snc.StorageNodeManagementClient)
	ret1, _ := ret[1].(*snpb.StorageNodeMetadataDescriptor)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMetadataByAddr indicates an expected call of GetMetadataByAddr.
func (mr *MockStorageNodeManagerMockRecorder) GetMetadataByAddr(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadataByAddr", reflect.TypeOf((*MockStorageNodeManager)(nil).GetMetadataByAddr), arg0, arg1)
}

// Refresh mocks base method.
func (m *MockStorageNodeManager) Refresh(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh.
func (mr *MockStorageNodeManagerMockRecorder) Refresh(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockStorageNodeManager)(nil).Refresh), arg0)
}

// RemoveLogStream mocks base method.
func (m *MockStorageNodeManager) RemoveLogStream(arg0 context.Context, arg1 types.StorageNodeID, arg2 types.TopicID, arg3 types.LogStreamID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveLogStream", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveLogStream indicates an expected call of RemoveLogStream.
func (mr *MockStorageNodeManagerMockRecorder) RemoveLogStream(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLogStream", reflect.TypeOf((*MockStorageNodeManager)(nil).RemoveLogStream), arg0, arg1, arg2, arg3)
}

// RemoveStorageNode mocks base method.
func (m *MockStorageNodeManager) RemoveStorageNode(arg0 types.StorageNodeID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveStorageNode", arg0)
}

// RemoveStorageNode indicates an expected call of RemoveStorageNode.
func (mr *MockStorageNodeManagerMockRecorder) RemoveStorageNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveStorageNode", reflect.TypeOf((*MockStorageNodeManager)(nil).RemoveStorageNode), arg0)
}

// Seal mocks base method.
func (m *MockStorageNodeManager) Seal(arg0 context.Context, arg1 types.TopicID, arg2 types.LogStreamID, arg3 types.GLSN) ([]snpb.LogStreamReplicaMetadataDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seal", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]snpb.LogStreamReplicaMetadataDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Seal indicates an expected call of Seal.
func (mr *MockStorageNodeManagerMockRecorder) Seal(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seal", reflect.TypeOf((*MockStorageNodeManager)(nil).Seal), arg0, arg1, arg2, arg3)
}

// Sync mocks base method.
func (m *MockStorageNodeManager) Sync(arg0 context.Context, arg1 types.TopicID, arg2 types.LogStreamID, arg3, arg4 types.StorageNodeID, arg5 types.GLSN) (*snpb.SyncStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*snpb.SyncStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sync indicates an expected call of Sync.
func (mr *MockStorageNodeManagerMockRecorder) Sync(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockStorageNodeManager)(nil).Sync), arg0, arg1, arg2, arg3, arg4, arg5)
}

// Trim mocks base method.
func (m *MockStorageNodeManager) Trim(arg0 context.Context, arg1 types.TopicID, arg2 types.GLSN) ([]vmspb.TrimResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trim", arg0, arg1, arg2)
	ret0, _ := ret[0].([]vmspb.TrimResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Trim indicates an expected call of Trim.
func (mr *MockStorageNodeManagerMockRecorder) Trim(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trim", reflect.TypeOf((*MockStorageNodeManager)(nil).Trim), arg0, arg1, arg2)
}

// Unseal mocks base method.
func (m *MockStorageNodeManager) Unseal(arg0 context.Context, arg1 types.TopicID, arg2 types.LogStreamID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unseal", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unseal indicates an expected call of Unseal.
func (mr *MockStorageNodeManagerMockRecorder) Unseal(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unseal", reflect.TypeOf((*MockStorageNodeManager)(nil).Unseal), arg0, arg1, arg2)
}

// MockMetadataRepositoryManager is a mock of MetadataRepositoryManager interface.
type MockMetadataRepositoryManager struct {
	ctrl     *gomock.Controller
	recorder *MockMetadataRepositoryManagerMockRecorder
}

// MockMetadataRepositoryManagerMockRecorder is the mock recorder for MockMetadataRepositoryManager.
type MockMetadataRepositoryManagerMockRecorder struct {
	mock *MockMetadataRepositoryManager
}

// NewMockMetadataRepositoryManager creates a new mock instance.
func NewMockMetadataRepositoryManager(ctrl *gomock.Controller) *MockMetadataRepositoryManager {
	mock := &MockMetadataRepositoryManager{ctrl: ctrl}
	mock.recorder = &MockMetadataRepositoryManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetadataRepositoryManager) EXPECT() *MockMetadataRepositoryManagerMockRecorder {
	return m.recorder
}

// AddPeer mocks base method.
func (m *MockMetadataRepositoryManager) AddPeer(arg0 context.Context, arg1 types.NodeID, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPeer", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPeer indicates an expected call of AddPeer.
func (mr *MockMetadataRepositoryManagerMockRecorder) AddPeer(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPeer", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).AddPeer), arg0, arg1, arg2, arg3)
}

// Close mocks base method.
func (m *MockMetadataRepositoryManager) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockMetadataRepositoryManagerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).Close))
}

// ClusterMetadataView mocks base method.
func (m *MockMetadataRepositoryManager) ClusterMetadataView() ClusterMetadataView {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClusterMetadataView")
	ret0, _ := ret[0].(ClusterMetadataView)
	return ret0
}

// ClusterMetadataView indicates an expected call of ClusterMetadataView.
func (mr *MockMetadataRepositoryManagerMockRecorder) ClusterMetadataView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClusterMetadataView", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).ClusterMetadataView))
}

// GetClusterInfo mocks base method.
func (m *MockMetadataRepositoryManager) GetClusterInfo(arg0 context.Context) (*mrpb.ClusterInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusterInfo", arg0)
	ret0, _ := ret[0].(*mrpb.ClusterInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusterInfo indicates an expected call of GetClusterInfo.
func (mr *MockMetadataRepositoryManagerMockRecorder) GetClusterInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterInfo", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).GetClusterInfo), arg0)
}

// NumberOfMR mocks base method.
func (m *MockMetadataRepositoryManager) NumberOfMR() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NumberOfMR")
	ret0, _ := ret[0].(int)
	return ret0
}

// NumberOfMR indicates an expected call of NumberOfMR.
func (mr *MockMetadataRepositoryManagerMockRecorder) NumberOfMR() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NumberOfMR", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).NumberOfMR))
}

// RegisterLogStream mocks base method.
func (m *MockMetadataRepositoryManager) RegisterLogStream(arg0 context.Context, arg1 *varlogpb.LogStreamDescriptor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterLogStream", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterLogStream indicates an expected call of RegisterLogStream.
func (mr *MockMetadataRepositoryManagerMockRecorder) RegisterLogStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterLogStream", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).RegisterLogStream), arg0, arg1)
}

// RegisterStorageNode mocks base method.
func (m *MockMetadataRepositoryManager) RegisterStorageNode(arg0 context.Context, arg1 *varlogpb.StorageNodeDescriptor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterStorageNode", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterStorageNode indicates an expected call of RegisterStorageNode.
func (mr *MockMetadataRepositoryManagerMockRecorder) RegisterStorageNode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterStorageNode", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).RegisterStorageNode), arg0, arg1)
}

// RegisterTopic mocks base method.
func (m *MockMetadataRepositoryManager) RegisterTopic(arg0 context.Context, arg1 types.TopicID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterTopic", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterTopic indicates an expected call of RegisterTopic.
func (mr *MockMetadataRepositoryManagerMockRecorder) RegisterTopic(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterTopic", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).RegisterTopic), arg0, arg1)
}

// RemovePeer mocks base method.
func (m *MockMetadataRepositoryManager) RemovePeer(arg0 context.Context, arg1 types.NodeID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemovePeer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePeer indicates an expected call of RemovePeer.
func (mr *MockMetadataRepositoryManagerMockRecorder) RemovePeer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePeer", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).RemovePeer), arg0, arg1)
}

// Seal mocks base method.
func (m *MockMetadataRepositoryManager) Seal(arg0 context.Context, arg1 types.LogStreamID) (types.GLSN, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seal", arg0, arg1)
	ret0, _ := ret[0].(types.GLSN)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Seal indicates an expected call of Seal.
func (mr *MockMetadataRepositoryManagerMockRecorder) Seal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seal", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).Seal), arg0, arg1)
}

// UnregisterLogStream mocks base method.
func (m *MockMetadataRepositoryManager) UnregisterLogStream(arg0 context.Context, arg1 types.LogStreamID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterLogStream", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterLogStream indicates an expected call of UnregisterLogStream.
func (mr *MockMetadataRepositoryManagerMockRecorder) UnregisterLogStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterLogStream", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).UnregisterLogStream), arg0, arg1)
}

// UnregisterStorageNode mocks base method.
func (m *MockMetadataRepositoryManager) UnregisterStorageNode(arg0 context.Context, arg1 types.StorageNodeID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterStorageNode", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterStorageNode indicates an expected call of UnregisterStorageNode.
func (mr *MockMetadataRepositoryManagerMockRecorder) UnregisterStorageNode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterStorageNode", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).UnregisterStorageNode), arg0, arg1)
}

// UnregisterTopic mocks base method.
func (m *MockMetadataRepositoryManager) UnregisterTopic(arg0 context.Context, arg1 types.TopicID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterTopic", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterTopic indicates an expected call of UnregisterTopic.
func (mr *MockMetadataRepositoryManagerMockRecorder) UnregisterTopic(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterTopic", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).UnregisterTopic), arg0, arg1)
}

// Unseal mocks base method.
func (m *MockMetadataRepositoryManager) Unseal(arg0 context.Context, arg1 types.LogStreamID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unseal", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unseal indicates an expected call of Unseal.
func (mr *MockMetadataRepositoryManagerMockRecorder) Unseal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unseal", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).Unseal), arg0, arg1)
}

// UpdateLogStream mocks base method.
func (m *MockMetadataRepositoryManager) UpdateLogStream(arg0 context.Context, arg1 *varlogpb.LogStreamDescriptor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLogStream", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLogStream indicates an expected call of UpdateLogStream.
func (mr *MockMetadataRepositoryManagerMockRecorder) UpdateLogStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLogStream", reflect.TypeOf((*MockMetadataRepositoryManager)(nil).UpdateLogStream), arg0, arg1)
}

// MockStorageNodeWatcher is a mock of StorageNodeWatcher interface.
type MockStorageNodeWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockStorageNodeWatcherMockRecorder
}

// MockStorageNodeWatcherMockRecorder is the mock recorder for MockStorageNodeWatcher.
type MockStorageNodeWatcherMockRecorder struct {
	mock *MockStorageNodeWatcher
}

// NewMockStorageNodeWatcher creates a new mock instance.
func NewMockStorageNodeWatcher(ctrl *gomock.Controller) *MockStorageNodeWatcher {
	mock := &MockStorageNodeWatcher{ctrl: ctrl}
	mock.recorder = &MockStorageNodeWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageNodeWatcher) EXPECT() *MockStorageNodeWatcherMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorageNodeWatcher) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageNodeWatcherMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorageNodeWatcher)(nil).Close))
}

// Run mocks base method.
func (m *MockStorageNodeWatcher) Run() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run")
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockStorageNodeWatcherMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockStorageNodeWatcher)(nil).Run))
}

// MockStatRepository is a mock of StatRepository interface.
type MockStatRepository struct {
	ctrl     *gomock.Controller
	recorder *MockStatRepositoryMockRecorder
}

// MockStatRepositoryMockRecorder is the mock recorder for MockStatRepository.
type MockStatRepositoryMockRecorder struct {
	mock *MockStatRepository
}

// NewMockStatRepository creates a new mock instance.
func NewMockStatRepository(ctrl *gomock.Controller) *MockStatRepository {
	mock := &MockStatRepository{ctrl: ctrl}
	mock.recorder = &MockStatRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatRepository) EXPECT() *MockStatRepositoryMockRecorder {
	return m.recorder
}

// GetAppliedIndex mocks base method.
func (m *MockStatRepository) GetAppliedIndex() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAppliedIndex")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetAppliedIndex indicates an expected call of GetAppliedIndex.
func (mr *MockStatRepositoryMockRecorder) GetAppliedIndex() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAppliedIndex", reflect.TypeOf((*MockStatRepository)(nil).GetAppliedIndex))
}

// GetLogStream mocks base method.
func (m *MockStatRepository) GetLogStream(arg0 types.LogStreamID) *LogStreamStat {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogStream", arg0)
	ret0, _ := ret[0].(*LogStreamStat)
	return ret0
}

// GetLogStream indicates an expected call of GetLogStream.
func (mr *MockStatRepositoryMockRecorder) GetLogStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogStream", reflect.TypeOf((*MockStatRepository)(nil).GetLogStream), arg0)
}

// GetStorageNode mocks base method.
func (m *MockStatRepository) GetStorageNode(arg0 types.StorageNodeID) *varlogpb.StorageNodeDescriptor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStorageNode", arg0)
	ret0, _ := ret[0].(*varlogpb.StorageNodeDescriptor)
	return ret0
}

// GetStorageNode indicates an expected call of GetStorageNode.
func (mr *MockStatRepositoryMockRecorder) GetStorageNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStorageNode", reflect.TypeOf((*MockStatRepository)(nil).GetStorageNode), arg0)
}

// Refresh mocks base method.
func (m *MockStatRepository) Refresh(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Refresh", arg0)
}

// Refresh indicates an expected call of Refresh.
func (mr *MockStatRepositoryMockRecorder) Refresh(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockStatRepository)(nil).Refresh), arg0)
}

// Report mocks base method.
func (m *MockStatRepository) Report(arg0 context.Context, arg1 *snpb.StorageNodeMetadataDescriptor) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Report", arg0, arg1)
}

// Report indicates an expected call of Report.
func (mr *MockStatRepositoryMockRecorder) Report(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Report", reflect.TypeOf((*MockStatRepository)(nil).Report), arg0, arg1)
}

// SetLogStreamStatus mocks base method.
func (m *MockStatRepository) SetLogStreamStatus(arg0 types.LogStreamID, arg1 varlogpb.LogStreamStatus) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLogStreamStatus", arg0, arg1)
}

// SetLogStreamStatus indicates an expected call of SetLogStreamStatus.
func (mr *MockStatRepositoryMockRecorder) SetLogStreamStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLogStreamStatus", reflect.TypeOf((*MockStatRepository)(nil).SetLogStreamStatus), arg0, arg1)
}
