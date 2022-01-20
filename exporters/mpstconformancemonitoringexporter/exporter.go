package mpstconformancemonitoringexporter

import (
	"context"
	"fmt"
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/globaltype"
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/mixedstateglobaltype"
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/pedro"
	"strings"

	"github.com/fangyi-zhou/mpst-tracing/labels"
	"github.com/pkg/errors"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
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

func extractMpstMetadata(attributes pdata.AttributeMap) (mpstMetadata, error) {
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

type mpstConformanceMonitoringExporter struct {
	logger *zap.Logger
	model  *model.Model
}

func (m mpstConformanceMonitoringExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (m mpstConformanceMonitoringExporter) processLocalTraces(traces pdata.Traces) error {
	m.logger.Info("Processing Traces", zap.Int("count", traces.SpanCount()))
	processedTraces := make(map[string][]model.Action)
	spans := traces.ResourceSpans()
	for i := 0; i < spans.Len(); i++ {
		span := spans.At(i)
		spanSlices := span.InstrumentationLibrarySpans()
		for j := 0; j < spanSlices.Len(); j++ {
			slice := spanSlices.At(j)
			innerSpans := slice.Spans()
			for k := 0; k < innerSpans.Len(); k++ {
				innerSpan := innerSpans.At(k)
				attributes := innerSpan.Attributes()
				if metadata, err := extractMpstMetadata(attributes); err == nil {
					if metadata.action == "Send" {
						message := model.Action{
							Label:  metadata.label,
							Src:    metadata.currentRole,
							Dest:   metadata.partner,
							IsSend: true,
						}
						processedTraces[metadata.currentRole] = append(
							processedTraces[metadata.currentRole],
							message,
						)
					} else if metadata.action == "Recv" {
						message := model.Action{
							Label:  metadata.label,
							Src:    metadata.partner,
							Dest:   metadata.currentRole,
							IsSend: false,
						}
						processedTraces[metadata.currentRole] = append(processedTraces[metadata.currentRole], message)
					} else {
						m.logger.Warn("Invalid action", zap.String("action", metadata.action))
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

func (m mpstConformanceMonitoringExporter) Start(_ context.Context, _ component.Host) error {
	return nil
}

func (m mpstConformanceMonitoringExporter) Shutdown(_ context.Context) error {
	m.model.Shutdown()
	return nil
}

func (m mpstConformanceMonitoringExporter) ConsumeTraces(
	_ context.Context,
	td pdata.Traces,
) error {
	err := m.processLocalTraces(td)
	return err
	/*
		err := checkSendRecvMatching(localTraces)
		if err != nil {
			return err
		}
		causalOrder := causalorder.Construct(m.logger, localTraces)
		err = causalOrder.CheckProtocolConformance(m.gtype)
		return err
	*/
}

/*
type participantPair struct {
	from string
	to   string
}

func checkSendRecvMatching(traces map[string]causalorder.LocalTrace) error {
	sendQueues := map[participantPair][]globaltype.Message{}
	recvQueues := map[participantPair][]globaltype.Message{}
	var errs []error
	for endpoint, localTrace := range traces {
		for _, message := range localTrace {
			if message.Action == "send" {
				sendQueues[participantPair{endpoint, message.Dest}] = append(sendQueues[participantPair{endpoint, message.Dest}], message)
			} else {
				recvQueues[participantPair{message.Origin, endpoint}] = append(recvQueues[participantPair{message.Origin, endpoint}], message)
			}
		}
	}
	for ppair, sendQueue := range sendQueues {
		for _, sendMsg := range sendQueue {
			if len(recvQueues[ppair]) == 0 {
				errs = append(errs, missingRecvMessageErr(ppair.from, ppair.to, sendMsg))
				continue
			}
			recvMsg := recvQueues[ppair][0]
			if sendMsg.Label != recvMsg.Label {
				errs = append(errs, mismatchLabelErr(ppair.from, ppair.to, sendMsg, recvMsg))
				continue
			}
			recvQueues[ppair] = recvQueues[ppair][1:]
		}
	}

	for ppair, recvQueue := range recvQueues {
		for _, recvMsg := range recvQueue {
			errs = append(errs, recvWithoutSendErr(ppair.from, ppair.to, recvMsg))
		}
	}

	if len(errs) != 0 {
		return componenterror.CombineErrors(errs)
	}
	return nil
}
*/

/*
func recvWithoutSendErr(orig string, dest string, msg globaltype.Message) error {
	return fmt.Errorf("message labelled %s received by %s without matching send, allegedly from %s", msg.Label, dest, orig)
}

func mismatchLabelErr(orig string, dest string, sendMsg globaltype.Message, recvMsg globaltype.Message) error {
	return fmt.Errorf("message label mismatch, sent from %s to %s, label %s is sent, but %s is received", orig, dest, sendMsg.Label, recvMsg.Label)
}

func missingRecvMessageErr(orig string, dest string, msg globaltype.Message) error {
	return fmt.Errorf("message labelled %s sent from %s is not received by %s", msg.Label, orig, dest)
}
*/

func newMpstConformanceExporter(
	logger *zap.Logger,
	cfg *Config,
) (*mpstConformanceMonitoringExporter, error) {
	var m model.Model
	switch cfg.SemanticModelType {
	case "gtype_lts":
		gtypeModel, err := globaltype.CreateGlobalTypeSemanticModel(
			cfg.GlobalTypeSexpFileName,
			logger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load global type")
		}
		logger.Info("Loaded global type")
		m = model.MakeModelWithLogger(gtypeModel, logger)
	case "gtype_pedro":
		pedroModel, err := pedro.CreatePedroSemanticModel(
			cfg.PedroSoFileName,
			cfg.ProtocolFileName,
			cfg.ProtocolName,
			logger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load pedro semantic model")
		}
		logger.Info("Loaded petri net model")
		m = model.MakeModelWithLogger(pedroModel, logger)
	case "gtype_mixed_state":
		gtypeModel, err := mixedstateglobaltype.CreateMixedStateGlobalTypeSemanticModel(
			cfg.GlobalTypeSexpFileName,
			logger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load global type")
		}
		logger.Info("Loaded global type (with mixed states)")
		m = model.MakeModelWithLogger(gtypeModel, logger)
	default:
		return nil, fmt.Errorf("unknown semantic model type %s", cfg.SemanticModelType)
	}
	return &mpstConformanceMonitoringExporter{
		logger: logger,
		model:  &m,
	}, nil
}
