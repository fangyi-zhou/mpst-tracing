package mpstconformancemonitoringprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
)

const (
	typeStr component.Type = "mpstconformancemonitoring"
)

func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesProcessor(createTraceProcessor, component.StabilityLevelDevelopment))
}

func createDefaultConfig() component.ProcessorConfig {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
	}
}

func createTraceProcessor(_ context.Context,
	params component.ProcessorCreateSettings,
	cfg component.ProcessorConfig,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	return newMpstConformanceProcessor(params.Logger, cfg.(*Config), nextConsumer)
}
