package mpstmetadatataggingprocessor

import (
	"context"

	"github.com/fangyi-zhou/mpst-tracing/labels"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
	"go.uber.org/zap"
)

type roleName string
type messageName string

type MpstMetadataTaggingProcessor struct {
	logger        *zap.Logger
	nextConsumer  consumer.Traces
	roleLookup    map[string]roleName
	messageLookup map[roleName]map[string]messageName
}

func (m MpstMetadataTaggingProcessor) Start(_ context.Context, _ component.Host) error {
	return nil
}

func (m MpstMetadataTaggingProcessor) Shutdown(_ context.Context) error {
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
		var roleNameExists = false
		var role roleName
		if serviceNameExists {
			service := serviceName.StringVal()
			role, roleNameExists = m.roleLookup[service]
		}
		ils := rs.InstrumentationLibrarySpans()
		for j := 0; j < ils.Len(); j++ {
			il := ils.At(j)
			if roleNameExists {
				// Update role via instrumentation library name, as currently is done.
				attachCurrentRoleTag(il, string(role))
			} else {
				m.logger.Warn("Unable to find role name from trace", zap.String("identifier", serviceName.StringVal()))
				continue
			}
			spans := il.Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				traceName := span.Name()
				if roleNameExists {
					message, messageNameExists := m.messageLookup[role][traceName]
					if messageNameExists {
						// TODO: attach appropriate metadata (partner, action) here
						// span.Attributes().Range(func(key string, val pdata.AttributeValue) bool {
						// 	m.logger.Info("attribute", zap.String("key", key), zap.String("val", val.StringVal()))
						// 	return true
						// })
						span.Attributes().InsertString(labels.MsgLabelKey, string(message))
						m.logger.Info(
							"Attached label to trace",
							zap.String("label", string(message)),
						)
					}
				}
			}
		}
	}
	return m.nextConsumer.ConsumeTraces(ctx, td)
}

func attachCurrentRoleTag(il pdata.InstrumentationLibrarySpans, role string) {
	spans := il.Spans()
	for i := 0; i < spans.Len(); i++ {
		span := spans.At(i)
		span.Attributes().InsertString(labels.CurrentRoleKey, role)
	}
}

func newMpstMetadataTaggingProcessor(
	logger *zap.Logger,
	config *Config,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	roleLookup := make(map[string]roleName)
	messageLookup := make(map[roleName]map[string]messageName)
	for role, roleData := range config.Roles {
		role := roleName(role)
		roleLookup[roleData.Name] = role
		messageLookup[role] = make(map[string]messageName)
		for message, messageData := range roleData.Messages {
			messageLookup[role][messageData.Name] = messageName(message)
		}
	}
	return &MpstMetadataTaggingProcessor{
		logger:        logger,
		nextConsumer:  nextConsumer,
		roleLookup:    roleLookup,
		messageLookup: messageLookup,
	}, nil
}
