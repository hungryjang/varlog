package metadata_repository

import (
	varlogpb "github.com/kakao/varlog/proto/varlog"
)

type MetadataRepository interface {
	RegisterStorageNode(*varlogpb.StorageNodeDescriptor) error
	CreateLogStream(*varlogpb.LogStreamDescriptor) error
	GetMetadata() (*varlogpb.MetadataDescriptor, error)
	Close() error
}
