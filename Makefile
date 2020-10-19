export

MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIR := $(dir $(MAKEFILE_PATH))
BUILD_DIR := $(MAKEFILE_DIR)/build

GO := go
GOPATH := $(shell $(GO) env GOPATH)
LDFLAGS :=
GOFLAGS := -race
GCFLAGS := -gcflags=all='-N -l'

PKG_SRCS := $(abspath $(shell find . -name '*.go'))
INTERNAL_SRCS := $(abspath $(shell find . -name '*.go'))

PROTOC_VERSION := 3.12.3
PROTOC_HOME := $(BUILD_DIR)/protoc
PROTOC := protoc
PROTO_INCS := -I ${GOPATH}/src -I ${MAKEFILE_DIR}/proto -I ${MAKEFILE_DIR}/vendor -I .

TEST_COUNT := 1
TEST_FLAGS := -count $(TEST_COUNT) -p 1

ifneq ($(TEST_CPU),)
	TEST_FLAGS := $(TEST_FLAGS) -cpu $(TEST_CPU)
endif

ifneq ($(TEST_TIMEOUT),)
	TEST_FLAGS := $(TEST_FLAGS) -timeout $(TEST_TIMEOUT)
endif

ifneq ($(TEST_PARALLEL),)
	TEST_FLAGS := $(TEST_FLAGS) -parallel $(TEST_PARALLEL)
endif

TEST_COVERAGE := 0
ifeq ($(TEST_COVERAGE),1)
	TEST_FLAGS := $(TEST_FLAGS) -coverprofile=$(BUILD_DIR)/reports/coverage.out
endif

TEST_FAILFAST := 0
ifeq ($(TEST_FAILFAST),1)
	TEST_FLAGS := $(TEST_FLAGS) -failfast
endif

TEST_VERBOSE := 1
ifeq ($(TEST_VERBOSE),1)
	TEST_FLAGS := $(TEST_FLAGS) -v
endif

TEST_DIRS := $(sort $(dir $(shell find . -name '*_test.go')))

all : proto libvarlog storagenode metadata_repository vms

VARLOGPB_PROTO := proto/varlogpb
SNPB_PROTO := proto/snpb
MRPB_PROTO := proto/mrpb
VMSPB_PROTO := proto/vmspb

PROTO := $(VARLOGPB_PROTO) $(SNPB_PROTO) $(MRPB_PROTO) $(VMSPB_PROTO)

proto : check_protoc gogoproto $(PROTO)

STORAGE_NODE := cmd/storagenode
storagenode : $(STORAGE_NODE_PROTO) $(STORAGE_NODE)

LIBVARLOG := pkg/varlog
libvarlog : $(LIBVARLOG)

METADATA_REPOSITORY := cmd/metadata_repository
metadata_repository : $(METADATA_REPOSITORY_PROTO) $(METADATA_REPOSITORY)

VMS := cmd/vms
vms : $(PROTO) $(VMS)

SUBDIRS := $(PROTO) $(STORAGE_NODE) $(LIBVARLOG) $(METADATA_REPOSITORY) $(VMS)
subdirs : $(SUBDIRS)

$(SUBDIRS) :
	$(MAKE) -C $@

mockgen: \
	pkg/varlog/varlog_mock.go \
	internal/vms/vms_mock.go \
	internal/storage/storage_node_mock.go \
	internal/storage/storage_mock.go \
	internal/storage/log_stream_executor_mock.go \
	internal/storage/log_stream_reporter_mock.go \
	internal/storage/log_stream_reporter_client_mock.go \
	internal/storage/replicator_mock.go \
	internal/storage/replicator_client_mock.go \
	proto/snpb/mock/replicator_mock.go \
	proto/snpb/mock/log_io_mock.go \
	proto/snpb/mock/log_stream_reporter_mock.go \
	proto/snpb/mock/management_mock.go \
	proto/mrpb/mock/management_mock.go \
	proto/mrpb/mock/metadata_repository_mock.go

pkg/varlog/varlog_mock.go: pkg/varlog/management_client.go
	mockgen -build_flags -mod=vendor \
		-self_package github.com/kakao/varlog/pkg/varlog \
		-package varlog \
		-destination $@ \
		github.com/kakao/varlog/pkg/varlog \
		ManagementClient

internal/vms/vms_mock.go: internal/vms/cluster_manager.go
	mockgen -build_flags -mod=vendor \
		-self_package github.com/kakao/varlog/internal/vms \
		-package vms \
		-destination $@ \
		github.com/kakao/varlog/internal/vms \
		ClusterMetadataView


internal/storage/storage_node_mock.go: internal/storage/storage_node.go
	mockgen -self_package github.com/kakao/varlog/internal/storage \
		-package storage \
		-source $< \
		-destination $@

internal/storage/storage_mock.go: internal/storage/storage.go
	mockgen -self_package github.com/kakao/varlog/internal/storage \
		-package storage \
		-source $< \
		-destination $@

internal/storage/log_stream_executor_mock.go: internal/storage/log_stream_executor.go
	mockgen -self_package github.com/kakao/varlog/internal/storage \
		-package storage \
		-source $< \
		-destination $@

internal/storage/log_stream_reporter_mock.go: internal/storage/log_stream_reporter.go
	mockgen -self_package github.com/kakao/varlog/internal/storage \
		-package storage \
		-source $< \
		-destination $@

internal/storage/log_stream_reporter_client_mock.go: internal/storage/log_stream_reporter_client.go
	mockgen -self_package github.com/kakao/varlog/internal/storage \
		-package storage \
		-source $< \
		-destination $@

internal/storage/replicator_mock.go: internal/storage/replicator.go
	mockgen -self_package github.com/kakao/varlog/internal/storage \
		-package storage \
		-source $< \
		-destination $@

internal/storage/replicator_client_mock.go: internal/storage/replicator_client.go
	mockgen -self_package github.com/kakao/varlog/internal/storage \
		-package storage \
		-source $< \
		-destination $@

proto/snpb/mock/replicator_mock.go: $(PROTO) proto/snpb/replicator.pb.go
	mockgen -build_flags -mod=vendor \
		-package mock \
		-destination $@ \
		github.com/kakao/varlog/proto/snpb \
		ReplicatorServiceClient,ReplicatorServiceServer,ReplicatorService_ReplicateClient,ReplicatorService_ReplicateServer

proto/snpb/mock/log_io_mock.go: $(PROTO) proto/snpb/log_io.pb.go
	mockgen -build_flags -mod=vendor \
		-package mock \
		-destination $@ \
		github.com/kakao/varlog/proto/snpb \
		LogIOClient,LogIOServer,LogIO_SubscribeClient,LogIO_SubscribeServer

proto/snpb/mock/log_stream_reporter_mock.go: $(PROTO) proto/snpb/log_stream_reporter.pb.go
	mockgen -build_flags -mod=vendor \
		-package mock \
		-destination $@ \
		github.com/kakao/varlog/proto/snpb \
		LogStreamReporterServiceClient,LogStreamReporterServiceServer

proto/snpb/mock/management_mock.go: $(PROTO) proto/snpb/management.pb.go
	mockgen -build_flags -mod=vendor \
		-package mock \
		-destination $@ \
		github.com/kakao/varlog/proto/snpb \
		ManagementClient,ManagementServer

proto/mrpb/mock/management_mock.go: $(PROTO) proto/mrpb/management.pb.go
	mockgen -build_flags -mod=vendor \
		-package mock \
		-destination $@ \
		github.com/kakao/varlog/proto/mrpb \
		ManagementClient,ManagementServer

proto/mrpb/mock/metadata_repository_mock.go: $(PROTO) proto/mrpb/metadata_repository.pb.go
	mockgen -build_flags -mod=vendor \
		-package mock \
		-destination $@ \
		github.com/kakao/varlog/proto/mrpb \
		MetadataRepositoryServiceClient,MetadataRepositoryServiceServer

.PHONY: test test_report coverage_report

test:
	tmpfile=$$(mktemp); \
	(TERM=sh $(GO) test $(GOFLAGS) $(GCFLAGS) $(TEST_FLAGS) ./... 2>&1; echo $$? > $$tmpfile) | \
	tee $(BUILD_DIR)/reports/test_output.txt; \
	ret=$$(cat $$tmpfile); \
	rm -f $$tmpfile; \
	exit $$ret

test_report:
	cat $(BUILD_DIR)/reports/test_output.txt | \
		go-junit-report > $(BUILD_DIR)/reports/report.xml
	rm $(BUILD_DIR)/reports/test_output.txt

coverage_report:
	gocov convert $(BUILD_DIR)/reports/coverage.out | gocov-xml > $(BUILD_DIR)/reports/coverage.xml

clean :
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir clean; \
	done

.PHONY: check_protoc
check_protoc:
ifneq ($(wildcard $(PROTOC_HOME)/bin/protoc),)
	$(info Installed protoc: $(shell $(PROTOC_HOME)/bin/protoc --version))
else
ifneq ($(shell which protoc),)
	$(info Installed protoc: $(shell protoc --version))
else
	$(error Install protoc. You can use '"make protoc"'.)
endif
endif

.PHONY: protoc
protoc: $(PROTOC_HOME)

$(PROTOC_HOME):
	PROTOC_HOME=$(PROTOC_HOME) PROTOC_VERSION=$(PROTOC_VERSION) scripts/install_protoc.sh

GOGOPROTO_SRC := $(GOPATH)/src/github.com/gogo/protobuf

.PHONY: gogoproto
gogoproto: $(GOGOPROTO_SRC)

$(GOGOPROTO_SRC):
	$(GO) get -u github.com/gogo/protobuf/protoc-gen-gogo
	$(GO) get -u github.com/gogo/protobuf/gogoproto
	$(GO) get -u github.com/gogo/protobuf/proto
	$(GO) get -u github.com/gogo/protobuf/jsonpb

.PHONY : all clean subdirs $(SUBDIRS) mockgen test
