package mpstmetadatataggingprocessor

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	typeStr config.Type = "mpstmetadatatagging"
)

func NewFactory() component.ProcessorFactory {
	return processorhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		processorhelper.WithTraces(createTracesProcessor),
	)
}

func createDefaultConfig() config.Processor {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewID(typeStr)),
	}
}

func createTracesProcessor(
	ctx context.Context,
	settings component.ProcessorCreateSettings,
	config config.Processor,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	return newMpstMetadataTaggingProcessor(settings.Logger, config.(*Config), nextConsumer)
}
