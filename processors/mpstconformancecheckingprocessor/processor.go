package mpstconformancecheckingprocessor

import (
	"context"
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
	label  string
	origin string
}

func (m mpstConformanceProcessor) ProcessTraces(ctx context.Context, traces pdata.Traces) (pdata.Traces, error) {
	spans := traces.ResourceSpans()
	messageQueue := map[string][]message{}
	_ = messageQueue // Make code compile
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
					message := message{label: label, origin: currentEndpoint}
					// TODO

					// Make code compile
					_ = partner
					_ = message
				} else if strings.HasPrefix(spanName, "Recv") {
					separated := strings.Split(spanName, " ")
					partner := separated[1]
					label := separated[2]
					// TODO

					// Make code compile
					_ = partner
					_ = label
				} else {
					m.logger.Info("Skipping unknown inner span name", zap.String("spanName", innerSpan.Name()))
				}
			}
		}
	}
	return traces, nil
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
