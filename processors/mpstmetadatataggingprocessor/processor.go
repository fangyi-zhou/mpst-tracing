package mpstmetadatataggingprocessor

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.uber.org/zap"
)

type MpstMetadataTaggingProcessor struct {
	Logger *zap.Logger
}

func (m MpstMetadataTaggingProcessor) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (m MpstMetadataTaggingProcessor) Shutdown(ctx context.Context) error {
	return nil
}

func (m MpstMetadataTaggingProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (m MpstMetadataTaggingProcessor) ConsumeTraces(ctx context.Context, td pdata.Traces) error {
	panic("implement me")
}

func newMpstMetadataTaggingProcessor(
	logger *zap.Logger,
	config *Config,
) (component.TracesProcessor, error) {
	return &MpstMetadataTaggingProcessor{
		Logger: logger,
	}, nil
}
