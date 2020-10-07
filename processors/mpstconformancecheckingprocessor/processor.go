package mpstconformancecheckingprocessor

import (
	"context"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
)

type mpstConformanceProcessor struct {
	logger *zap.Logger
}

func (m mpstConformanceProcessor) ProcessTraces(ctx context.Context, traces pdata.Traces) (pdata.Traces, error) {
	spans := traces.ResourceSpans()
	for i := 0; i < spans.Len(); i++ {
		span := spans.At(i)
		if span.IsNil() {
			continue
		}
		serviceName := processor.ServiceNameForResource(span.Resource())
		m.logger.Info("Found trace for service", zap.String("serviceName", serviceName))
	}
	return traces, nil
}

func newMpstConformanceProcessor(logger *zap.Logger, nextConsumer consumer.TraceConsumer, cfg *Config) (mpstConformanceProcessor, error) {
	return mpstConformanceProcessor{
		logger: logger,
	}, nil
}
