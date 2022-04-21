// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kakao/varlog/pkg/varlog (interfaces: Admin)

// Package varlog is a generated GoMock package.
package varlog

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	types "github.com/kakao/varlog/pkg/types"
	snpb "github.com/kakao/varlog/proto/snpb"
	varlogpb "github.com/kakao/varlog/proto/varlogpb"
	vmspb "github.com/kakao/varlog/proto/vmspb"
)

// MockAdmin is a mock of Admin interface.
type MockAdmin struct {
	ctrl     *gomock.Controller
	recorder *MockAdminMockRecorder
}

// MockAdminMockRecorder is the mock recorder for MockAdmin.
type MockAdminMockRecorder struct {
	mock *MockAdmin
}

// NewMockAdmin creates a new mock instance.
func NewMockAdmin(ctrl *gomock.Controller) *MockAdmin {
	mock := &MockAdmin{ctrl: ctrl}
	mock.recorder = &MockAdminMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdmin) EXPECT() *MockAdminMockRecorder {
	return m.recorder
}

// AddLogStream mocks base method.
func (m *MockAdmin) AddLogStream(arg0 context.Context, arg1 types.TopicID, arg2 []*varlogpb.ReplicaDescriptor) (*varlogpb.LogStreamDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLogStream", arg0, arg1, arg2)
	ret0, _ := ret[0].(*varlogpb.LogStreamDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddLogStream indicates an expected call of AddLogStream.
func (mr *MockAdminMockRecorder) AddLogStream(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLogStream", reflect.TypeOf((*MockAdmin)(nil).AddLogStream), arg0, arg1, arg2)
}

// AddMRPeer mocks base method.
func (m *MockAdmin) AddMRPeer(arg0 context.Context, arg1, arg2 string) (types.NodeID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMRPeer", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.NodeID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMRPeer indicates an expected call of AddMRPeer.
func (mr *MockAdminMockRecorder) AddMRPeer(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMRPeer", reflect.TypeOf((*MockAdmin)(nil).AddMRPeer), arg0, arg1, arg2)
}

// AddStorageNode mocks base method.
func (m *MockAdmin) AddStorageNode(arg0 context.Context, arg1 string) (*snpb.StorageNodeMetadataDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddStorageNode", arg0, arg1)
	ret0, _ := ret[0].(*snpb.StorageNodeMetadataDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddStorageNode indicates an expected call of AddStorageNode.
func (mr *MockAdminMockRecorder) AddStorageNode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddStorageNode", reflect.TypeOf((*MockAdmin)(nil).AddStorageNode), arg0, arg1)
}

// AddTopic mocks base method.
func (m *MockAdmin) AddTopic(arg0 context.Context) (varlogpb.TopicDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTopic", arg0)
	ret0, _ := ret[0].(varlogpb.TopicDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTopic indicates an expected call of AddTopic.
func (mr *MockAdminMockRecorder) AddTopic(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTopic", reflect.TypeOf((*MockAdmin)(nil).AddTopic), arg0)
}

// Close mocks base method.
func (m *MockAdmin) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockAdminMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAdmin)(nil).Close))
}

// DescribeTopic mocks base method.
func (m *MockAdmin) DescribeTopic(arg0 context.Context, arg1 types.TopicID) (*vmspb.DescribeTopicResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeTopic", arg0, arg1)
	ret0, _ := ret[0].(*vmspb.DescribeTopicResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeTopic indicates an expected call of DescribeTopic.
func (mr *MockAdminMockRecorder) DescribeTopic(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeTopic", reflect.TypeOf((*MockAdmin)(nil).DescribeTopic), arg0, arg1)
}

// GetMRMembers mocks base method.
func (m *MockAdmin) GetMRMembers(arg0 context.Context) (*vmspb.GetMRMembersResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMRMembers", arg0)
	ret0, _ := ret[0].(*vmspb.GetMRMembersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMRMembers indicates an expected call of GetMRMembers.
func (mr *MockAdminMockRecorder) GetMRMembers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMRMembers", reflect.TypeOf((*MockAdmin)(nil).GetMRMembers), arg0)
}

// GetStorageNodes mocks base method.
func (m *MockAdmin) GetStorageNodes(arg0 context.Context) (map[types.StorageNodeID]*snpb.StorageNodeMetadataDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStorageNodes", arg0)
	ret0, _ := ret[0].(map[types.StorageNodeID]*snpb.StorageNodeMetadataDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStorageNodes indicates an expected call of GetStorageNodes.
func (mr *MockAdminMockRecorder) GetStorageNodes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStorageNodes", reflect.TypeOf((*MockAdmin)(nil).GetStorageNodes), arg0)
}

// RemoveLogStreamReplica mocks base method.
func (m *MockAdmin) RemoveLogStreamReplica(arg0 context.Context, arg1 types.StorageNodeID, arg2 types.TopicID, arg3 types.LogStreamID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveLogStreamReplica", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveLogStreamReplica indicates an expected call of RemoveLogStreamReplica.
func (mr *MockAdminMockRecorder) RemoveLogStreamReplica(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLogStreamReplica", reflect.TypeOf((*MockAdmin)(nil).RemoveLogStreamReplica), arg0, arg1, arg2, arg3)
}

// RemoveMRPeer mocks base method.
func (m *MockAdmin) RemoveMRPeer(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveMRPeer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveMRPeer indicates an expected call of RemoveMRPeer.
func (mr *MockAdminMockRecorder) RemoveMRPeer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveMRPeer", reflect.TypeOf((*MockAdmin)(nil).RemoveMRPeer), arg0, arg1)
}

// Seal mocks base method.
func (m *MockAdmin) Seal(arg0 context.Context, arg1 types.TopicID, arg2 types.LogStreamID) (*vmspb.SealResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seal", arg0, arg1, arg2)
	ret0, _ := ret[0].(*vmspb.SealResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Seal indicates an expected call of Seal.
func (mr *MockAdminMockRecorder) Seal(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seal", reflect.TypeOf((*MockAdmin)(nil).Seal), arg0, arg1, arg2)
}

// Sync mocks base method.
func (m *MockAdmin) Sync(arg0 context.Context, arg1 types.TopicID, arg2 types.LogStreamID, arg3, arg4 types.StorageNodeID) (*snpb.SyncStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*snpb.SyncStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sync indicates an expected call of Sync.
func (mr *MockAdminMockRecorder) Sync(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockAdmin)(nil).Sync), arg0, arg1, arg2, arg3, arg4)
}

// Topics mocks base method.
func (m *MockAdmin) Topics(arg0 context.Context) ([]varlogpb.TopicDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Topics", arg0)
	ret0, _ := ret[0].([]varlogpb.TopicDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Topics indicates an expected call of Topics.
func (mr *MockAdminMockRecorder) Topics(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Topics", reflect.TypeOf((*MockAdmin)(nil).Topics), arg0)
}

// Trim mocks base method.
func (m *MockAdmin) Trim(arg0 context.Context, arg1 types.TopicID, arg2 types.GLSN) (map[types.LogStreamID]map[types.StorageNodeID]error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trim", arg0, arg1, arg2)
	ret0, _ := ret[0].(map[types.LogStreamID]map[types.StorageNodeID]error)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Trim indicates an expected call of Trim.
func (mr *MockAdminMockRecorder) Trim(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trim", reflect.TypeOf((*MockAdmin)(nil).Trim), arg0, arg1, arg2)
}

// UnregisterLogStream mocks base method.
func (m *MockAdmin) UnregisterLogStream(arg0 context.Context, arg1 types.TopicID, arg2 types.LogStreamID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterLogStream", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterLogStream indicates an expected call of UnregisterLogStream.
func (mr *MockAdminMockRecorder) UnregisterLogStream(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterLogStream", reflect.TypeOf((*MockAdmin)(nil).UnregisterLogStream), arg0, arg1, arg2)
}

// UnregisterStorageNode mocks base method.
func (m *MockAdmin) UnregisterStorageNode(arg0 context.Context, arg1 types.StorageNodeID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterStorageNode", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterStorageNode indicates an expected call of UnregisterStorageNode.
func (mr *MockAdminMockRecorder) UnregisterStorageNode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterStorageNode", reflect.TypeOf((*MockAdmin)(nil).UnregisterStorageNode), arg0, arg1)
}

// UnregisterTopic mocks base method.
func (m *MockAdmin) UnregisterTopic(arg0 context.Context, arg1 types.TopicID) (*vmspb.UnregisterTopicResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterTopic", arg0, arg1)
	ret0, _ := ret[0].(*vmspb.UnregisterTopicResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnregisterTopic indicates an expected call of UnregisterTopic.
func (mr *MockAdminMockRecorder) UnregisterTopic(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterTopic", reflect.TypeOf((*MockAdmin)(nil).UnregisterTopic), arg0, arg1)
}

// Unseal mocks base method.
func (m *MockAdmin) Unseal(arg0 context.Context, arg1 types.TopicID, arg2 types.LogStreamID) (*varlogpb.LogStreamDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unseal", arg0, arg1, arg2)
	ret0, _ := ret[0].(*varlogpb.LogStreamDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unseal indicates an expected call of Unseal.
func (mr *MockAdminMockRecorder) Unseal(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unseal", reflect.TypeOf((*MockAdmin)(nil).Unseal), arg0, arg1, arg2)
}

// UpdateLogStream mocks base method.
func (m *MockAdmin) UpdateLogStream(arg0 context.Context, arg1 types.TopicID, arg2 types.LogStreamID, arg3, arg4 *varlogpb.ReplicaDescriptor) (*varlogpb.LogStreamDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLogStream", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*varlogpb.LogStreamDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateLogStream indicates an expected call of UpdateLogStream.
func (mr *MockAdminMockRecorder) UpdateLogStream(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLogStream", reflect.TypeOf((*MockAdmin)(nil).UpdateLogStream), arg0, arg1, arg2, arg3, arg4)
}
