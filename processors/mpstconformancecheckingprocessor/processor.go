package mpstconformancecheckingprocessor

import (
	"context"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.uber.org/zap"
)

type mpstConformanceProcessor struct {
}

func (m mpstConformanceProcessor) ProcessTraces(ctx context.Context, traces pdata.Traces) (pdata.Traces, error) {
	panic("implement me")
}

func newMpstConformanceProcessor(logger *zap.Logger, nextConsumer consumer.TraceConsumer, cfg *Config) (mpstConformanceProcessor, error) {
	return mpstConformanceProcessor{}, nil
}
