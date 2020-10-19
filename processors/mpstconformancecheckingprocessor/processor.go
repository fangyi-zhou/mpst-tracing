package mpstconformancecheckingprocessor

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
	"strings"
)

type mpstConformanceProcessor struct {
	logger *zap.Logger
}

type message struct {
	label string
	//	origin string
}

func (m mpstConformanceProcessor) ProcessTraces(ctx context.Context, traces pdata.Traces) (pdata.Traces, error) {
	spans := traces.ResourceSpans()
	sendQueues := map[string]map[string][]message{}
	recvQueues := map[string]map[string][]message{}
	for i := 0; i < spans.Len(); i++ {
		span := spans.At(i)
		if span.IsNil() {
			continue
		}
		serviceName := processor.ServiceNameForResource(span.Resource())
		m.logger.Info("Found trace for service", zap.String("serviceName", serviceName))
		spanSlices := span.InstrumentationLibrarySpans()
		for j := 0; j < spanSlices.Len(); j++ {
			slice := spanSlices.At(j)
			if slice.IsNil() {
				continue
			}
			library := slice.InstrumentationLibrary()
			if library.IsNil() {
				m.logger.Warn("Cannot get instrumentation library, skipping")
				continue
			}
			libraryName := library.Name()
			currentEndpoint := getEndpointFromLibraryName(libraryName)
			innerSpans := slice.Spans()
			for k := 0; k < innerSpans.Len(); k++ {
				innerSpan := innerSpans.At(k)
				if innerSpan.IsNil() {
					continue
				}
				spanName := innerSpan.Name()
				if strings.HasPrefix(spanName, "Send") {
					separated := strings.Split(spanName, " ")
					partner := separated[1]
					label := separated[2]
					if sendQueues[currentEndpoint] == nil {
						sendQueues[currentEndpoint] = map[string][]message{}
					}
					sendQueues[currentEndpoint][partner] = append(sendQueues[currentEndpoint][partner], message{label})
				} else if strings.HasPrefix(spanName, "Recv") {
					separated := strings.Split(spanName, " ")
					partner := separated[1]
					label := separated[2]
					if recvQueues[partner] == nil {
						recvQueues[partner] = map[string][]message{}
					}
					recvQueues[partner][currentEndpoint] = append(recvQueues[partner][currentEndpoint], message{label})
				} else {
					m.logger.Info("Skipping unknown inner span name", zap.String("spanName", innerSpan.Name()))
				}
			}
		}
	}

	var errs []error
	for orig, sendPartialQueue := range sendQueues {
		for dest, sendQueue := range sendPartialQueue {
			for _, sendMsg := range sendQueue {
				if len(recvQueues[orig][dest]) == 0 {
					errs = append(errs, missingRecvMessageErr(orig, dest, sendMsg))
					continue
				}
				recvMsg := recvQueues[orig][dest][0]
				if sendMsg.label != recvMsg.label {
					errs = append(errs, mismatchLabelErr(orig, dest, sendMsg, recvMsg))
					continue
				}
				recvQueues[orig][dest] = recvQueues[orig][dest][1:]
			}
		}
	}

	for orig, recvPartialQueue := range recvQueues {
		for dest, recvQueue := range recvPartialQueue {
			for _, recvMsg := range recvQueue {
				errs = append(errs, recvWithoutSendErr(orig, dest, recvMsg))
			}
		}
	}

	if len(errs) != 0 {
		return traces, componenterror.CombineErrors(errs)
	}
	return traces, nil
}

func recvWithoutSendErr(orig string, dest string, msg message) error {
	return fmt.Errorf("message labelled %s received by %s without matching send, allegedly from %s", msg.label, dest, orig)
}

func mismatchLabelErr(orig string, dest string, sendMsg message, recvMsg message) error {
	return fmt.Errorf("message label mismatch, sent from %s to %s, label %s is sent, but %s is received", orig, dest, sendMsg.label, recvMsg.label)
}

func missingRecvMessageErr(orig string, dest string, msg message) error {
	return fmt.Errorf("message labelled %s sent from %s is not received by %s", msg.label, orig, dest)
}

func getEndpointFromLibraryName(libraryName string) string {
	separated := strings.Split(libraryName, "/")
	return separated[len(separated)-1]
}

func newMpstConformanceProcessor(logger *zap.Logger, nextConsumer consumer.TraceConsumer, cfg *Config) (mpstConformanceProcessor, error) {
	return mpstConformanceProcessor{
		logger: logger,
	}, nil
}
