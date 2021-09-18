package telemetry

//go:generate mockgen -build_flags -mod=vendor -self_package github.com/kakao/varlog/internal/storagenode/telemetry -package telemetry -destination telemetry_mock.go . Measurable

import (
	"context"

	"github.com/kakao/varlog/pkg/types"
	"github.com/kakao/varlog/pkg/util/telemetry"
	"github.com/kakao/varlog/pkg/util/telemetry/metric"
	"github.com/kakao/varlog/pkg/util/telemetry/trace"
)

const (
	serviceName = "sn"
)

type Measurable interface {
	Stub() *Stub
}

type Stub struct {
	tm telemetry.Telemetry
	tr trace.Tracer
	mt metric.Meter
	mb *MetricsBag
}

func NewTelemetryStub(ctx context.Context, telemetryType string, storageNodeID types.StorageNodeID, endpoint string) (*Stub, error) {
	tm, err := telemetry.New(ctx, serviceName, storageNodeID.String(),
		telemetry.WithExporterType(telemetryType),
		telemetry.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, err
	}
	stub := &Stub{
		tm: tm,
		tr: telemetry.Tracer(telemetry.ServiceNamespace + "." + serviceName),
		mt: telemetry.Meter(telemetry.ServiceNamespace + "." + serviceName),
	}
	stub.mb = newMetricsBag(stub)
	return stub, nil
}

func NewNopTelmetryStub() *Stub {
	tm := telemetry.NewNopTelemetry()
	stub := &Stub{
		tm: tm,
		tr: telemetry.Tracer(telemetry.ServiceNamespace + "." + serviceName),
		mt: telemetry.Meter(telemetry.ServiceNamespace + "." + serviceName),
	}
	stub.mb = newMetricsBag(stub)
	return stub
}

func (ts *Stub) close(ctx context.Context) {
	ts.tm.Close(ctx)
}

func (ts *Stub) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return ts.tr.Start(ctx, name, opts...)
}

func (ts *Stub) Metrics() *MetricsBag {
	return ts.mb
}

func (ts *Stub) Meter() metric.Meter {
	return ts.mt
}
