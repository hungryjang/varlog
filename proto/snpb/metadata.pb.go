// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/snpb/metadata.proto

package snpb

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"

	github_com_kakao_varlog_pkg_types "github.com/kakao/varlog/pkg/types"
	varlogpb "github.com/kakao/varlog/proto/varlogpb"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// StorageNodeMetadataDescriptor represents the metadata of stroage node.
type StorageNodeMetadataDescriptor struct {
	// ClusterID is the identifier of the cluster.
	ClusterID github_com_kakao_varlog_pkg_types.ClusterID `protobuf:"varint,1,opt,name=cluster_id,json=clusterId,proto3,casttype=github.com/kakao/varlog/pkg/types.ClusterID" json:"clusterId"`
	// StorageNode is detailed information about the storage node.
	StorageNode *varlogpb.StorageNodeDescriptor `protobuf:"bytes,2,opt,name=storage_node,json=storageNode,proto3" json:"storageNode"`
	// LogStreams are the list of metadata for log stream replicas.
	LogStreamReplicas []LogStreamReplicaMetadataDescriptor `protobuf:"bytes,3,rep,name=log_stream_replicas,json=logStreamReplicas,proto3" json:"logStreams"`
	// CreatedTime is the creation time of the storage node.
	// Note that the CreatedTime is immutable after the metadata repository sets.
	// TODO: Currently the storage node has no responsibility to persist
	// CreatedTime. How can we tell a recovered storage node the CreatedTime?
	//
	// Deprecated:
	CreatedTime time.Time `protobuf:"bytes,4,opt,name=created_time,json=createdTime,proto3,stdtime" json:"createdTime"`
	// UpdatedTime
	//
	// Deprecated:
	UpdatedTime time.Time `protobuf:"bytes,5,opt,name=updated_time,json=updatedTime,proto3,stdtime" json:"updatedTime"`
}

func (m *StorageNodeMetadataDescriptor) Reset()         { *m = StorageNodeMetadataDescriptor{} }
func (m *StorageNodeMetadataDescriptor) String() string { return proto.CompactTextString(m) }
func (*StorageNodeMetadataDescriptor) ProtoMessage()    {}
func (*StorageNodeMetadataDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_b0d7c3885ca513ae, []int{0}
}
func (m *StorageNodeMetadataDescriptor) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StorageNodeMetadataDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StorageNodeMetadataDescriptor.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StorageNodeMetadataDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StorageNodeMetadataDescriptor.Merge(m, src)
}
func (m *StorageNodeMetadataDescriptor) XXX_Size() int {
	return m.ProtoSize()
}
func (m *StorageNodeMetadataDescriptor) XXX_DiscardUnknown() {
	xxx_messageInfo_StorageNodeMetadataDescriptor.DiscardUnknown(m)
}

var xxx_messageInfo_StorageNodeMetadataDescriptor proto.InternalMessageInfo

func (m *StorageNodeMetadataDescriptor) GetClusterID() github_com_kakao_varlog_pkg_types.ClusterID {
	if m != nil {
		return m.ClusterID
	}
	return 0
}

func (m *StorageNodeMetadataDescriptor) GetStorageNode() *varlogpb.StorageNodeDescriptor {
	if m != nil {
		return m.StorageNode
	}
	return nil
}

func (m *StorageNodeMetadataDescriptor) GetLogStreamReplicas() []LogStreamReplicaMetadataDescriptor {
	if m != nil {
		return m.LogStreamReplicas
	}
	return nil
}

func (m *StorageNodeMetadataDescriptor) GetCreatedTime() time.Time {
	if m != nil {
		return m.CreatedTime
	}
	return time.Time{}
}

func (m *StorageNodeMetadataDescriptor) GetUpdatedTime() time.Time {
	if m != nil {
		return m.UpdatedTime
	}
	return time.Time{}
}

// LogStreamReplicaMetadataDescriptor represents the metadata of log stream
// replica.
type LogStreamReplicaMetadataDescriptor struct {
	varlogpb.LogStreamReplica `protobuf:"bytes,1,opt,name=log_stream_replica,json=logStreamReplica,proto3,embedded=log_stream_replica" json:"log_stream_replica"`
	// Status is the status of the log stream replica.
	//
	// TODO: Use a separate type to represent the status of the log stream replica
	// rather than `varlogpb.LogStreamStatus` that is shared with the metadata
	// repository.
	Status varlogpb.LogStreamStatus `protobuf:"varint,2,opt,name=status,proto3,enum=varlog.varlogpb.LogStreamStatus" json:"status,omitempty"`
	// Version is the latest version of the commit received from the metadata
	// repository.
	Version github_com_kakao_varlog_pkg_types.Version `protobuf:"varint,3,opt,name=version,proto3,casttype=github.com/kakao/varlog/pkg/types.Version" json:"version,omitempty"`
	// GlobalHighWatermark is the latest high watermark received from the metadata
	// repository.
	GlobalHighWatermark github_com_kakao_varlog_pkg_types.GLSN `protobuf:"varint,4,opt,name=global_high_watermark,json=globalHighWatermark,proto3,casttype=github.com/kakao/varlog/pkg/types.GLSN" json:"globalHighWatermark"`
	// LocalLowWatermark is the first log sequence number in the log stream
	// replica.
	// The LocalLowWatermark becomes higher when the log is truncated by prefix
	// trimming.
	LocalLowWatermark varlogpb.LogSequenceNumber `protobuf:"bytes,5,opt,name=local_low_watermark,json=localLowWatermark,proto3" json:"localLowWatermark"`
	// LocalHighWatermark is the last log sequence number in the log stream
	// replica.
	LocalHighWatermark varlogpb.LogSequenceNumber `protobuf:"bytes,6,opt,name=local_high_watermark,json=localHighWatermark,proto3" json:"localHighWatermark"`
	// Path is the directory where the data for the log stream replica is stored.
	Path string `protobuf:"bytes,7,opt,name=path,proto3" json:"path,omitempty"`
	// CreatedTime
	//
	// FIXME: StartTime or UpTime
	CreatedTime time.Time `protobuf:"bytes,8,opt,name=created_time,json=createdTime,proto3,stdtime" json:"createdTime"`
	// UpdatedTime
	//
	// Deprecated:
	UpdatedTime time.Time `protobuf:"bytes,9,opt,name=updated_time,json=updatedTime,proto3,stdtime" json:"updatedTime"`
}

func (m *LogStreamReplicaMetadataDescriptor) Reset()         { *m = LogStreamReplicaMetadataDescriptor{} }
func (m *LogStreamReplicaMetadataDescriptor) String() string { return proto.CompactTextString(m) }
func (*LogStreamReplicaMetadataDescriptor) ProtoMessage()    {}
func (*LogStreamReplicaMetadataDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_b0d7c3885ca513ae, []int{1}
}
func (m *LogStreamReplicaMetadataDescriptor) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LogStreamReplicaMetadataDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LogStreamReplicaMetadataDescriptor.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LogStreamReplicaMetadataDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogStreamReplicaMetadataDescriptor.Merge(m, src)
}
func (m *LogStreamReplicaMetadataDescriptor) XXX_Size() int {
	return m.ProtoSize()
}
func (m *LogStreamReplicaMetadataDescriptor) XXX_DiscardUnknown() {
	xxx_messageInfo_LogStreamReplicaMetadataDescriptor.DiscardUnknown(m)
}

var xxx_messageInfo_LogStreamReplicaMetadataDescriptor proto.InternalMessageInfo

func (m *LogStreamReplicaMetadataDescriptor) GetStatus() varlogpb.LogStreamStatus {
	if m != nil {
		return m.Status
	}
	return varlogpb.LogStreamStatusRunning
}

func (m *LogStreamReplicaMetadataDescriptor) GetVersion() github_com_kakao_varlog_pkg_types.Version {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *LogStreamReplicaMetadataDescriptor) GetGlobalHighWatermark() github_com_kakao_varlog_pkg_types.GLSN {
	if m != nil {
		return m.GlobalHighWatermark
	}
	return 0
}

func (m *LogStreamReplicaMetadataDescriptor) GetLocalLowWatermark() varlogpb.LogSequenceNumber {
	if m != nil {
		return m.LocalLowWatermark
	}
	return varlogpb.LogSequenceNumber{}
}

func (m *LogStreamReplicaMetadataDescriptor) GetLocalHighWatermark() varlogpb.LogSequenceNumber {
	if m != nil {
		return m.LocalHighWatermark
	}
	return varlogpb.LogSequenceNumber{}
}

func (m *LogStreamReplicaMetadataDescriptor) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *LogStreamReplicaMetadataDescriptor) GetCreatedTime() time.Time {
	if m != nil {
		return m.CreatedTime
	}
	return time.Time{}
}

func (m *LogStreamReplicaMetadataDescriptor) GetUpdatedTime() time.Time {
	if m != nil {
		return m.UpdatedTime
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*StorageNodeMetadataDescriptor)(nil), "varlog.snpb.StorageNodeMetadataDescriptor")
	proto.RegisterType((*LogStreamReplicaMetadataDescriptor)(nil), "varlog.snpb.LogStreamReplicaMetadataDescriptor")
}

func init() { proto.RegisterFile("proto/snpb/metadata.proto", fileDescriptor_b0d7c3885ca513ae) }

var fileDescriptor_b0d7c3885ca513ae = []byte{
	// 662 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0x4f, 0x6f, 0xd3, 0x30,
	0x14, 0x6f, 0x58, 0xe9, 0x56, 0x77, 0xfc, 0x99, 0x07, 0x5a, 0x57, 0x44, 0x5d, 0x7a, 0x40, 0xbd,
	0x2c, 0x15, 0x05, 0xa1, 0x89, 0x63, 0x98, 0x34, 0x26, 0x95, 0x1d, 0xdc, 0x69, 0x88, 0x5d, 0x22,
	0xb7, 0x31, 0x6e, 0xd4, 0xa4, 0x0e, 0xb6, 0xb3, 0x69, 0x37, 0x3e, 0xc2, 0x2e, 0xdc, 0xf7, 0x3d,
	0xf8, 0x02, 0x3b, 0xee, 0xc8, 0x29, 0x48, 0xeb, 0x05, 0xf5, 0x23, 0xec, 0x84, 0xe2, 0x24, 0x6b,
	0x69, 0x87, 0xd8, 0x24, 0x38, 0xc5, 0x7e, 0xef, 0xf7, 0x7e, 0xef, 0xcf, 0xef, 0x59, 0x01, 0xeb,
	0x81, 0xe0, 0x8a, 0x37, 0xe5, 0x30, 0xe8, 0x36, 0x7d, 0xaa, 0x88, 0x43, 0x14, 0x31, 0xb5, 0x0d,
	0x96, 0x0e, 0x89, 0xf0, 0x38, 0x33, 0x63, 0x5f, 0x65, 0x83, 0xb9, 0xaa, 0x1f, 0x76, 0xcd, 0x1e,
	0xf7, 0x9b, 0x8c, 0x33, 0xde, 0xd4, 0x98, 0x6e, 0xf8, 0x49, 0xdf, 0x12, 0x92, 0xf8, 0x94, 0xc4,
	0x56, 0x9e, 0x30, 0xce, 0x99, 0x47, 0x27, 0x28, 0xea, 0x07, 0xea, 0x38, 0x75, 0xa2, 0x59, 0xa7,
	0x72, 0x7d, 0x2a, 0x15, 0xf1, 0x83, 0x14, 0xb0, 0x96, 0x64, 0x9e, 0x2b, 0xa9, 0xfe, 0x35, 0x0f,
	0x9e, 0x76, 0x14, 0x17, 0x84, 0xd1, 0x5d, 0xee, 0xd0, 0xf7, 0xa9, 0x77, 0x8b, 0xca, 0x9e, 0x70,
	0x03, 0xc5, 0x05, 0x94, 0x00, 0xf4, 0xbc, 0x50, 0x2a, 0x2a, 0x6c, 0xd7, 0x29, 0x1b, 0x35, 0xa3,
	0x71, 0xcf, 0xda, 0xbb, 0x88, 0x50, 0xf1, 0x6d, 0x62, 0xdd, 0xd9, 0x1a, 0x47, 0xa8, 0x98, 0x42,
	0x76, 0x9c, 0xcb, 0x08, 0xbd, 0x4e, 0x3b, 0x73, 0x48, 0xe8, 0x0f, 0xc8, 0x80, 0x70, 0xdd, 0x63,
	0x52, 0x41, 0xf6, 0x09, 0x06, 0xac, 0xa9, 0x8e, 0x03, 0x2a, 0xcd, 0x2b, 0x1a, 0x3c, 0x21, 0x81,
	0x07, 0x60, 0x59, 0x26, 0x55, 0xd9, 0x43, 0xee, 0xd0, 0xf2, 0x9d, 0x9a, 0xd1, 0x28, 0xb5, 0x9e,
	0x9b, 0xe9, 0x00, 0xb3, 0x6e, 0xcc, 0xa9, 0xd2, 0x27, 0x25, 0x5b, 0x0f, 0xc6, 0x11, 0x2a, 0xc9,
	0x89, 0x0b, 0x4f, 0x5f, 0xa0, 0x04, 0xab, 0x1e, 0x67, 0xb6, 0x54, 0x82, 0x12, 0xdf, 0x16, 0x34,
	0xf0, 0xdc, 0x1e, 0x91, 0xe5, 0x85, 0xda, 0x42, 0xa3, 0xd4, 0x6a, 0x9a, 0x53, 0x1a, 0x99, 0x6d,
	0xce, 0x3a, 0x1a, 0x86, 0x13, 0xd4, 0xfc, 0x78, 0x2c, 0x78, 0x16, 0xa1, 0xdc, 0x38, 0x42, 0xc0,
	0xcb, 0xb0, 0x12, 0xaf, 0x78, 0x33, 0x71, 0x12, 0xee, 0x83, 0xe5, 0x9e, 0xa0, 0x44, 0x51, 0xc7,
	0x8e, 0xb5, 0x29, 0xe7, 0x75, 0x43, 0x15, 0x33, 0x11, 0xce, 0xcc, 0x84, 0x33, 0xf7, 0x32, 0xe1,
	0xac, 0xb5, 0x94, 0xb8, 0x94, 0xc6, 0xc5, 0x9e, 0x93, 0x1f, 0xc8, 0xc0, 0xd3, 0x86, 0x98, 0x37,
	0x0c, 0x9c, 0x09, 0xef, 0xdd, 0x9b, 0xf3, 0xa6, 0x71, 0x13, 0xde, 0x29, 0x43, 0xfd, 0x5b, 0x01,
	0xd4, 0xff, 0xde, 0x3d, 0xfc, 0x08, 0xe0, 0xfc, 0x2c, 0xf5, 0x92, 0x94, 0x5a, 0xcf, 0xe6, 0xd4,
	0x9a, 0x25, 0xb4, 0x96, 0xe2, 0x5a, 0xce, 0x23, 0x64, 0xe0, 0x87, 0xb3, 0x23, 0x83, 0x9b, 0xa0,
	0x20, 0x15, 0x51, 0xa1, 0xd4, 0xe2, 0xdf, 0x6f, 0xd5, 0xfe, 0x4c, 0xd7, 0xd1, 0x38, 0x9c, 0xe2,
	0x21, 0x06, 0x8b, 0x87, 0x54, 0x48, 0x97, 0x0f, 0xcb, 0x0b, 0x35, 0xa3, 0x91, 0xb7, 0x36, 0x2f,
	0x23, 0xf4, 0xea, 0x56, 0x4b, 0xb9, 0x9f, 0xc4, 0xe3, 0x8c, 0x08, 0x7e, 0x31, 0xc0, 0x63, 0xe6,
	0xf1, 0x2e, 0xf1, 0xec, 0xbe, 0xcb, 0xfa, 0xf6, 0x11, 0x51, 0x54, 0xf8, 0x44, 0x0c, 0xb4, 0x92,
	0x79, 0xab, 0x3d, 0x8e, 0xd0, 0x6a, 0x02, 0x78, 0xe7, 0xb2, 0xfe, 0x87, 0xcc, 0x7d, 0x19, 0xa1,
	0x17, 0xb7, 0xca, 0xbc, 0xdd, 0xee, 0xec, 0xe2, 0xeb, 0x98, 0xa0, 0x1f, 0xef, 0x6d, 0x8f, 0x78,
	0xb6, 0xc7, 0x8f, 0xa6, 0xf2, 0x27, 0x8a, 0xd7, 0xaf, 0x9d, 0x0e, 0xfd, 0x1c, 0xd2, 0x61, 0x8f,
	0xee, 0x86, 0x7e, 0x97, 0x0a, 0x6b, 0x3d, 0x55, 0x7e, 0x45, 0xd3, 0xb4, 0xf9, 0xd1, 0x15, 0x37,
	0x9e, 0x37, 0xc1, 0x00, 0x3c, 0x4a, 0xd2, 0xcd, 0xf4, 0x5b, 0xb8, 0x71, 0xbe, 0x4a, 0x9a, 0x0f,
	0x6a, 0x9e, 0xdf, 0x9a, 0xc1, 0xd7, 0xd8, 0x20, 0x04, 0xf9, 0x80, 0xa8, 0x7e, 0x79, 0xb1, 0x66,
	0x34, 0x8a, 0x58, 0x9f, 0xe7, 0xde, 0xcd, 0xd2, 0x7f, 0x7a, 0x37, 0xc5, 0x7f, 0xf3, 0x6e, 0xde,
	0xe4, 0x7f, 0x9e, 0x22, 0xc3, 0xda, 0x3e, 0xbb, 0xa8, 0x1a, 0xe7, 0x17, 0x55, 0xe3, 0x64, 0x54,
	0xcd, 0x9d, 0x8e, 0xaa, 0xc6, 0xf9, 0xa8, 0x9a, 0xfb, 0x3e, 0xaa, 0xe6, 0x0e, 0x36, 0x6e, 0xb2,
	0x0c, 0x57, 0x3f, 0x90, 0x6e, 0x41, 0x9f, 0x5f, 0xfe, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x6b, 0x4f,
	0x1d, 0xee, 0x55, 0x06, 0x00, 0x00,
}

func (this *LogStreamReplicaMetadataDescriptor) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*LogStreamReplicaMetadataDescriptor)
	if !ok {
		that2, ok := that.(LogStreamReplicaMetadataDescriptor)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.LogStreamReplica.Equal(&that1.LogStreamReplica) {
		return false
	}
	if this.Status != that1.Status {
		return false
	}
	if this.Version != that1.Version {
		return false
	}
	if this.GlobalHighWatermark != that1.GlobalHighWatermark {
		return false
	}
	if !this.LocalLowWatermark.Equal(&that1.LocalLowWatermark) {
		return false
	}
	if !this.LocalHighWatermark.Equal(&that1.LocalHighWatermark) {
		return false
	}
	if this.Path != that1.Path {
		return false
	}
	if !this.CreatedTime.Equal(that1.CreatedTime) {
		return false
	}
	if !this.UpdatedTime.Equal(that1.UpdatedTime) {
		return false
	}
	return true
}
func (m *StorageNodeMetadataDescriptor) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StorageNodeMetadataDescriptor) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StorageNodeMetadataDescriptor) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.UpdatedTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintMetadata(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x2a
	n2, err2 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.CreatedTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedTime):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintMetadata(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x22
	if len(m.LogStreamReplicas) > 0 {
		for iNdEx := len(m.LogStreamReplicas) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LogStreamReplicas[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMetadata(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.StorageNode != nil {
		{
			size, err := m.StorageNode.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMetadata(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.ClusterID != 0 {
		i = encodeVarintMetadata(dAtA, i, uint64(m.ClusterID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *LogStreamReplicaMetadataDescriptor) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LogStreamReplicaMetadataDescriptor) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LogStreamReplicaMetadataDescriptor) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n4, err4 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.UpdatedTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedTime):])
	if err4 != nil {
		return 0, err4
	}
	i -= n4
	i = encodeVarintMetadata(dAtA, i, uint64(n4))
	i--
	dAtA[i] = 0x4a
	n5, err5 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.CreatedTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedTime):])
	if err5 != nil {
		return 0, err5
	}
	i -= n5
	i = encodeVarintMetadata(dAtA, i, uint64(n5))
	i--
	dAtA[i] = 0x42
	if len(m.Path) > 0 {
		i -= len(m.Path)
		copy(dAtA[i:], m.Path)
		i = encodeVarintMetadata(dAtA, i, uint64(len(m.Path)))
		i--
		dAtA[i] = 0x3a
	}
	{
		size, err := m.LocalHighWatermark.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintMetadata(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size, err := m.LocalLowWatermark.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintMetadata(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.GlobalHighWatermark != 0 {
		i = encodeVarintMetadata(dAtA, i, uint64(m.GlobalHighWatermark))
		i--
		dAtA[i] = 0x20
	}
	if m.Version != 0 {
		i = encodeVarintMetadata(dAtA, i, uint64(m.Version))
		i--
		dAtA[i] = 0x18
	}
	if m.Status != 0 {
		i = encodeVarintMetadata(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x10
	}
	{
		size, err := m.LogStreamReplica.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintMetadata(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintMetadata(dAtA []byte, offset int, v uint64) int {
	offset -= sovMetadata(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *StorageNodeMetadataDescriptor) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ClusterID != 0 {
		n += 1 + sovMetadata(uint64(m.ClusterID))
	}
	if m.StorageNode != nil {
		l = m.StorageNode.ProtoSize()
		n += 1 + l + sovMetadata(uint64(l))
	}
	if len(m.LogStreamReplicas) > 0 {
		for _, e := range m.LogStreamReplicas {
			l = e.ProtoSize()
			n += 1 + l + sovMetadata(uint64(l))
		}
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedTime)
	n += 1 + l + sovMetadata(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedTime)
	n += 1 + l + sovMetadata(uint64(l))
	return n
}

func (m *LogStreamReplicaMetadataDescriptor) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.LogStreamReplica.ProtoSize()
	n += 1 + l + sovMetadata(uint64(l))
	if m.Status != 0 {
		n += 1 + sovMetadata(uint64(m.Status))
	}
	if m.Version != 0 {
		n += 1 + sovMetadata(uint64(m.Version))
	}
	if m.GlobalHighWatermark != 0 {
		n += 1 + sovMetadata(uint64(m.GlobalHighWatermark))
	}
	l = m.LocalLowWatermark.ProtoSize()
	n += 1 + l + sovMetadata(uint64(l))
	l = m.LocalHighWatermark.ProtoSize()
	n += 1 + l + sovMetadata(uint64(l))
	l = len(m.Path)
	if l > 0 {
		n += 1 + l + sovMetadata(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedTime)
	n += 1 + l + sovMetadata(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedTime)
	n += 1 + l + sovMetadata(uint64(l))
	return n
}

func sovMetadata(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMetadata(x uint64) (n int) {
	return sovMetadata(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *StorageNodeMetadataDescriptor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMetadata
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: StorageNodeMetadataDescriptor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StorageNodeMetadataDescriptor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterID", wireType)
			}
			m.ClusterID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ClusterID |= github_com_kakao_varlog_pkg_types.ClusterID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StorageNode", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMetadata
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.StorageNode == nil {
				m.StorageNode = &varlogpb.StorageNodeDescriptor{}
			}
			if err := m.StorageNode.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LogStreamReplicas", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMetadata
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LogStreamReplicas = append(m.LogStreamReplicas, LogStreamReplicaMetadataDescriptor{})
			if err := m.LogStreamReplicas[len(m.LogStreamReplicas)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMetadata
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.CreatedTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdatedTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMetadata
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.UpdatedTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMetadata(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMetadata
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *LogStreamReplicaMetadataDescriptor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMetadata
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: LogStreamReplicaMetadataDescriptor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LogStreamReplicaMetadataDescriptor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LogStreamReplica", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMetadata
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LogStreamReplica.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= varlogpb.LogStreamStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			m.Version = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Version |= github_com_kakao_varlog_pkg_types.Version(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GlobalHighWatermark", wireType)
			}
			m.GlobalHighWatermark = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GlobalHighWatermark |= github_com_kakao_varlog_pkg_types.GLSN(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LocalLowWatermark", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMetadata
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LocalLowWatermark.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LocalHighWatermark", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMetadata
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LocalHighWatermark.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Path", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMetadata
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Path = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMetadata
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.CreatedTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdatedTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMetadata
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.UpdatedTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMetadata(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMetadata
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMetadata(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMetadata
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMetadata
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMetadata
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMetadata
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMetadata        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMetadata          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMetadata = fmt.Errorf("proto: unexpected end of group")
)
