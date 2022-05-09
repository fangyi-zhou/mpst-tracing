package mpstconformancemonitoringprocessor

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"strings"

	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/globaltype"
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/mixedstateglobaltype"
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/pedro"

	"github.com/fangyi-zhou/mpst-tracing/labels"
	"github.com/pkg/errors"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type mpstMetadata struct {
	action      string
	label       string
	partner     string
	currentRole string
}

func getEndpointFromLibraryName(libraryName string) string {
	separated := strings.Split(libraryName, "/")
	return separated[len(separated)-1]
}

func extractMpstMetadata(attributes pcommon.Map) (mpstMetadata, error) {
	action, hasAction := attributes.Get(labels.ActionKey)
	label, hasLabel := attributes.Get(labels.MsgLabelKey)
	partner, hasPartner := attributes.Get(labels.PartnerKey)
	currentRole, hasCurrentRole := attributes.Get(labels.CurrentRoleKey)
	if hasAction && hasLabel && hasPartner && hasCurrentRole {
		return mpstMetadata{
			action:      action.StringVal(),
			label:       label.StringVal(),
			partner:     partner.StringVal(),
			currentRole: getEndpointFromLibraryName(currentRole.StringVal()),
		}, nil
	} else {
		return mpstMetadata{}, errors.New("No Mpst Metadata present")
	}
}

type mpstConformanceMonitoringProcessor struct {
	logger       *zap.Logger
	nextConsumer consumer.Traces
	modelFactory model.ModelFactory
	model        model.Model
}

func (m *mpstConformanceMonitoringProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (m *mpstConformanceMonitoringProcessor) processLocalTraces(traces ptrace.Traces) error {
	m.logger.Info("Processing Traces", zap.Int("count", traces.SpanCount()))
	processedTraces := make(map[string][]model.Action)
	spans := traces.ResourceSpans()
	for i := 0; i < spans.Len(); i++ {
		span := spans.At(i)
		spanSlices := span.ScopeSpans()
		for j := 0; j < spanSlices.Len(); j++ {
			slice := spanSlices.At(j)
			innerSpans := slice.Spans()
			for k := 0; k < innerSpans.Len(); k++ {
				innerSpan := innerSpans.At(k)
				attributes := innerSpan.Attributes()
				if metadata, err := extractMpstMetadata(attributes); err == nil {
					message := model.Action{
						Label: metadata.label,
					}
					ok := false
					if metadata.action == "Send" {
						message.Src = metadata.currentRole
						message.Dest = metadata.partner
						message.IsSend = true
						ok = true
					} else if metadata.action == "Recv" {
						message.Dest = metadata.currentRole
						message.Src = metadata.partner
						message.IsSend = false
						ok = true
					} else {
						m.logger.Warn("Invalid action", zap.String("action", metadata.action))
					}
					if ok {
						// See: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/batchpersignal/batchpersignal.go
						trace := ptrace.NewTraces()
						newRS := trace.ResourceSpans().AppendEmpty()
						span.Resource().CopyTo(newRS.Resource())
						newSS := newRS.ScopeSpans().AppendEmpty()
						slice.Scope().CopyTo(newSS.Scope())
						target := newRS.ScopeSpans().At(0).Spans().AppendEmpty()
						innerSpan.CopyTo(target)
						done := make(chan bool)
						message.Done = done
						go func() {
							validated := <-done
							//m.logger.Info("Done", zap.String("action", message.String()))
							span := trace.ResourceSpans().
								At(0).
								ScopeSpans().
								At(0).
								Spans().
								At(0)
							attr := span.Attributes()
							attr.InsertBool(labels.ValidatedKey, validated)
							if !validated && span.Status().Code() != ptrace.StatusCodeError {
								span.Status().SetCode(ptrace.StatusCodeError)
							}
							err := m.nextConsumer.ConsumeTraces(context.TODO(), trace)
							if err != nil {
								m.logger.Fatal("error", zap.Error(err))
							}
						}()
						processedTraces[metadata.currentRole] = append(
							processedTraces[metadata.currentRole],
							message,
						)
					}
				}
			}
		}
	}
	for endpoint, localTrace := range processedTraces {
		m.model.AcceptTrace(endpoint, localTrace)
	}
	return nil
}

func (m *mpstConformanceMonitoringProcessor) Start(_ context.Context, _ component.Host) error {
	return nil
}

func (m *mpstConformanceMonitoringProcessor) Shutdown(_ context.Context) error {
	m.model.Shutdown()
	return nil
}

func (m *mpstConformanceMonitoringProcessor) ConsumeTraces(
	ctx context.Context,
	td ptrace.Traces,
) error {
	err := m.processLocalTraces(td)
	if err != nil {
		return err
	}
	return nil
	//return m.nextConsumer.ConsumeTraces(ctx, td)
}

func newMpstConformanceProcessor(
	logger *zap.Logger,
	cfg *Config,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	var factory model.ModelFactory
	var err error
	switch cfg.SemanticModelType {
	case "gtype_lts":
		factory, err = globaltype.CreateGlobalTypeModelFactory(
			cfg.GlobalTypeSexpFileName,
			cfg.GlobalTypeProtobufFileName,
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load global type")
		}
		logger.Info("Loaded global type")
	case "gtype_pedro":
		factory, err = pedro.CreatePedroModelFactory(
			cfg.PedroSoFileName,
			cfg.ProtocolFileName,
			cfg.ProtocolName,
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load pedro semantic model")
		}
		logger.Info("Loaded petri net")
	case "gtype_mixed_state":
		factory, err = mixedstateglobaltype.CreateMixedStateGlobalTypeModelFactory(
			cfg.GlobalTypeSexpFileName,
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load global type")
		}
		logger.Info("Loaded global type (with mixed states)")
	default:
		return nil, fmt.Errorf("unknown semantic model type %s", cfg.SemanticModelType)
	}
	m, err := factory.MakeModelWithLogger(logger)
	if err != nil {
		return nil, err
	}
	return &mpstConformanceMonitoringProcessor{
		logger:       logger,
		nextConsumer: nextConsumer,
		modelFactory: factory,
		model:        m,
	}, nil
}
