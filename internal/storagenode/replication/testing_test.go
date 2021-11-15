package replication

import (
	gomock "github.com/golang/mock/gomock"

	"github.com/kakao/varlog/internal/storagenode/telemetry"
)

func NewTestMeasurable(ctrl *gomock.Controller) *telemetry.MockMeasurable {
	m := telemetry.NewMockMeasurable(ctrl)
	nop := telemetry.NewNopTelmetryStub()
	m.EXPECT().Stub().Return(nop).AnyTimes()
	return m
}