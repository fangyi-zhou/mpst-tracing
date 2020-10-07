package mpstconformancecheckingprocessor

import (
	"context"
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
	return &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			NameVal: "MPSTConformanceChecking",
			TypeVal: typeStr,
		},
	}
}

func createTraceProcessor(ctx context.Context, params component.ProcessorCreateParams, cfg configmodels.Processor, nextConsumer consumer.TraceConsumer) (component.TraceProcessor, error) {
	tp, err := newMpstConformanceProcessor(params.Logger, nextConsumer, cfg.(*Config))
	if err != nil {
		return nil, err
	}
	return processorhelper.NewTraceProcessor(cfg, nextConsumer, tp)
}
