package mpstmetadatataggingprocessor

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.uber.org/zap"
)

type MpstMetadataTaggingProcessor struct {
	logger       *zap.Logger
	nextConsumer consumer.Traces
	roleLookup   map[string]string
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
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		serviceName, serviceNameExists := rs.Resource().Attributes().Get("service.name")
		var roleNameExists bool = false
		var roleName string
		if serviceNameExists {
			service := serviceName.StringVal()
			roleName, roleNameExists = m.roleLookup[service]
		}
		ils := rs.InstrumentationLibrarySpans()
		for j := 0; j < ils.Len(); j++ {
			il := ils.At(j)
			if roleNameExists {
				// Update role via instrumentation library name, as currently is done.
				il.InstrumentationLibrary().SetName(roleName)
			} else {
				m.logger.Warn("Unable to find role name from trace", zap.String("identifier", serviceName.StringVal()))
			}
			spans := il.Spans()
			for k := 0; k < spans.Len(); k++ {
				// TODO: attach appropriate metadata here
				//span := spans.At(k)
				//m.logger.Info(
				//	"Found span",
				//	zap.String("traceName", span.Name()),
				//	zap.String("traceId", span.SpanID().HexString()),
				//)
			}
		}
	}
	return m.nextConsumer.ConsumeTraces(ctx, td)
}

func newMpstMetadataTaggingProcessor(
	logger *zap.Logger,
	config *Config,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	roleLookup := make(map[string]string)
	for role, roleData := range config.Roles {
		roleLookup[roleData.Name] = role
	}
	return &MpstMetadataTaggingProcessor{
		logger:       logger,
		nextConsumer: nextConsumer,
		roleLookup:   roleLookup,
	}, nil
}
