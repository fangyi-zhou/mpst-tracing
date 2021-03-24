package twobuyer

import (
	"context"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagators"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"log"
	"time"
)

// https://github.com/open-telemetry/opentelemetry-go/blob/master/example/namedtracer/main.go
func InitStdoutTracer() func() {
	var err error
	exp, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Panicf("failed to initialize stdout exporter %v\n", err)
		return func() {}
	}
	bsp := trace.NewBatchSpanProcessor(exp)
	tp := trace.NewTracerProvider(
		trace.WithConfig(
			trace.Config{
				DefaultSampler: trace.AlwaysSample(),
			},
		),
		trace.WithSpanProcessor(bsp),
	)
	global.SetTracerProvider(tp)
	return func() {}
}

// https://github.com/open-telemetry/opentelemetry-go/blob/master/example/jaeger/main.go
func InitJaegerTracer() func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "TwoBuyer",
			Tags:        []label.KeyValue{},
		}),
		jaeger.WithSDK(&trace.Config{DefaultSampler: trace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}

	return flush
}

func InitOtlpTracer() func() {
	// https://github.com/open-telemetry/opentelemetry-go/blob/master/example/otel-collector/main.go
	exp, err := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress("localhost:55680"),
	)
	if err != nil {
		log.Panicf("Failed to create exporter, %v\n", err)
		return nil
	}

	bsp := trace.NewBatchSpanProcessor(exp)
	tracerProvider := trace.NewTracerProvider(
		trace.WithConfig(trace.Config{DefaultSampler: trace.AlwaysSample()}),
		trace.WithResource(resource.New(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("TwoBuyer"),
		)),
		trace.WithSpanProcessor(bsp),
	)

	pusher := push.New(
		basic.New(
			simple.NewWithExactDistribution(),
			exp,
		),
		exp,
		push.WithPeriod(2*time.Second),
	)

	global.SetTextMapPropagator(propagators.TraceContext{})
	global.SetTracerProvider(tracerProvider)
	global.SetMeterProvider(pusher.MeterProvider())
	pusher.Start()

	return func() {
		bsp.Shutdown() // shutdown the processor
		err := exp.Shutdown(context.Background())
		if err != nil {
			log.Panicf("Failed to stop exporter %v\n", err)
		}
		pusher.Stop() // pushes any last exports to the receiver
	}
}
