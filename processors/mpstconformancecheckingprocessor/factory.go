package mpstconformancecheckingprocessor

import (
	"context"
	"errors"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	typeStr configmodels.Type = "mpstconformancechecking"
)

func NewFactory() component.ProcessorFactory {
	return processorhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		processorhelper.WithTraces(createTraceProcessor))
}

func createDefaultConfig() configmodels.Processor {
	return nil
}

func createTraceProcessor(ctx context.Context, params component.ProcessorCreateParams, cfg configmodels.Processor, nextConsumer consumer.TraceConsumer) (component.TraceProcessor, error) {
	return nil, errors.New("not Implemented")
}