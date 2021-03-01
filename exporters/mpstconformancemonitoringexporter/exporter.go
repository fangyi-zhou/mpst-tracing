package mpstconformancemonitoringexporter

import (
	"context"
	"fmt"
	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/globaltype"
	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/causalorder"
	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/types"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.uber.org/zap"
	"strings"
)

type mpstConformanceMonitoringExporter struct {
	logger *zap.Logger
	gtype  globaltype.GlobalType
}

func (m mpstConformanceMonitoringExporter) extractLocalTraces(traces pdata.Traces) map[string]causalorder.LocalTrace {
	var processedTraces = map[string]causalorder.LocalTrace{}
	spans := traces.ResourceSpans()
	for i := 0; i < spans.Len(); i++ {
		span := spans.At(i)
		spanSlices := span.InstrumentationLibrarySpans()
		for j := 0; j < spanSlices.Len(); j++ {
			slice := spanSlices.At(j)
			library := slice.InstrumentationLibrary()
			libraryName := library.Name()
			currentEndpoint := getEndpointFromLibraryName(libraryName)
			innerSpans := slice.Spans()
			for k := 0; k < innerSpans.Len(); k++ {
				innerSpan := innerSpans.At(k)
				attributes := innerSpan.Attributes()
				if hasMpstMetadata(attributes) {
					partner_, _ := attributes.Get(partnerKey)
					msgLabel_, _ := attributes.Get(msgLabelKey)
					action_, _ := attributes.Get(actionKey)
					partner := partner_.StringVal()
					action := action_.StringVal()
					label := msgLabel_.StringVal()
					if action == "send" {
						message := types.Message{
							Label:  label,
							Origin: currentEndpoint,
							Dest:   partner,
							Action: "send",
						}
						processedTraces[currentEndpoint] = append(processedTraces[currentEndpoint], message)
					} else if action == "recv" {
						message := types.Message{
							Label:  label,
							Origin: partner,
							Dest:   currentEndpoint,
							Action: "recv",
						}
						processedTraces[currentEndpoint] = append(processedTraces[currentEndpoint], message)
					} else {
						m.logger.Warn("Invalid action", zap.String("action", action))
					}
				}
			}
		}
	}
	return processedTraces
}

func (m mpstConformanceMonitoringExporter) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (m mpstConformanceMonitoringExporter) Shutdown(ctx context.Context) error {
	return nil
}

func (m mpstConformanceMonitoringExporter) ConsumeTraces(ctx context.Context, td pdata.Traces) error {
	localTraces := m.extractLocalTraces(td)
	err := checkSendRecvMatching(localTraces)
	if err != nil {
		return err
	}
	causalOrder := causalorder.Construct(localTraces)
	err = causalOrder.CheckProtocolConformance(globaltype.TwoBuyer())
	return err
}

var (
	actionKey   = "mpst/action"
	msgLabelKey = "mpst/msgLabel"
	partnerKey  = "mpst/partner"
)

type participantPair struct {
	from string
	to   string
}

func checkSendRecvMatching(traces map[string]causalorder.LocalTrace) error {
	sendQueues := map[participantPair][]types.Message{}
	recvQueues := map[participantPair][]types.Message{}
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

func hasMpstMetadata(attributes pdata.AttributeMap) bool {
	_, hasAction := attributes.Get(actionKey)
	_, hasLabel := attributes.Get(msgLabelKey)
	_, hasPartner := attributes.Get(partnerKey)
	return hasAction && hasLabel && hasPartner
}

func recvWithoutSendErr(orig string, dest string, msg types.Message) error {
	return fmt.Errorf("message labelled %s received by %s without matching send, allegedly from %s", msg.Label, dest, orig)
}

func mismatchLabelErr(orig string, dest string, sendMsg types.Message, recvMsg types.Message) error {
	return fmt.Errorf("message label mismatch, sent from %s to %s, label %s is sent, but %s is received", orig, dest, sendMsg.Label, recvMsg.Label)
}

func missingRecvMessageErr(orig string, dest string, msg types.Message) error {
	return fmt.Errorf("message labelled %s sent from %s is not received by %s", msg.Label, orig, dest)
}

func getEndpointFromLibraryName(libraryName string) string {
	separated := strings.Split(libraryName, "/")
	return separated[len(separated)-1]
}
func newMpstConformanceExporter(logger *zap.Logger, cfg *Config) (*mpstConformanceMonitoringExporter, error) {
	// gtype, err := globaltype.LoadFromSexpFile(cfg.Protocol)
	// if err != nil {
	//	return nil, err
	// }
	return &mpstConformanceMonitoringExporter{
		logger: logger,
		gtype:  globaltype.TwoBuyer(), // Hard code to Two Buyer for now
	}, nil
}