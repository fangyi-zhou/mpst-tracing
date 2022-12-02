package mpstmetadatataggingprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
)

const (
	typeStr component.Type = "mpstmetadatatagging"
)

func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesProcessor(createTracesProcessor, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.ProcessorConfig {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
	}
}

func createTracesProcessor(
	_ context.Context,
	settings component.ProcessorCreateSettings,
	config component.ProcessorConfig,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	return newMpstMetadataTaggingProcessor(settings.Logger, config.(*Config), nextConsumer)
}
