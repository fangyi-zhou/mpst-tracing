package app

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/propagation"
	trace2 "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagators"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"log"
	"math/rand"
	"sync"
	"time"
)

import "go.opentelemetry.io/otel/exporters/stdout"

func (a *A) run(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	var span trace2.Span
	ctx, span = a.tracer.Start(ctx, "TwoBuyer Endpoint A")
	defer span.End()
	// Send query to B
	var query = rand.Intn(100)
	fmt.Println("A: Sending query", query)
	a.sendB(ctx, "query", query)
	// Receive a quote
	var quote = a.recvB(ctx, "quote")
	var otherShare = a.recvC(ctx, "share")
	if otherShare*2 >= quote {
		// 1 stands for ok
		a.sendB(ctx, "buy", 1)
	} else {
		a.sendB(ctx, "buy", 0)
	}
}
func (b *B) run(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	var span trace2.Span
	ctx, span = b.tracer.Start(ctx, "TwoBuyer Endpoint B")
	defer span.End()
	// Receive a query
	var query = b.recvA(ctx, "query")
	// Send a quote
	var quote = query * 2
	fmt.Println("B: Sending quote", quote)
	b.sendA(ctx, "quote", quote)
	b.sendC(ctx, "quote", quote)
	var decision = b.recvA(ctx, "buy")
	if decision == 1 {
		fmt.Println("Succeed!")
	} else {
		fmt.Println("Failed to succeed!")
	}
}

func (c *C) run(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	var span trace2.Span
	ctx, span = c.tracer.Start(ctx, "TwoBuyer Endpoint C")
	defer span.End()
	// Receive a quote
	var quote = c.recvB(ctx, "quote")
	// Propose a share
	var share = quote/2 + rand.Intn(10) - 5
	fmt.Println("C: Proposing share", share)
	c.sendA(ctx, "share", share)
}

var tp *trace.TracerProvider

// https://github.com/open-telemetry/opentelemetry-go/blob/master/example/namedtracer/main.go
func initStdoutTracer() func() {
	exp, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Panicf("failed to initialise jaeger exporter %v\n", err)
		return nil
	}
	bsp := trace.NewBatchSpanProcessor(exp)
	tp = trace.NewTracerProvider(
		trace.WithConfig(
			trace.Config{
				DefaultSampler: trace.AlwaysSample(),
			}),
		trace.WithSpanProcessor(bsp))
	global.SetTracerProvider(tp)
	return bsp.Shutdown
}

// https://github.com/open-telemetry/opentelemetry-go/blob/master/example/jaeger/main.go
func initJaegerTracer() func() {
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

	return func() {
		flush()
	}
}

func initOtlpTracer() func() {
	// https://github.com/open-telemetry/opentelemetry-go/blob/master/example/otel-collector/main.go
	exp, err := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress("localhost:30080"),
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

	tcPropagator := propagators.TraceContext{}
	props := propagation.New(propagation.WithExtractors(tcPropagator),
		propagation.WithInjectors(tcPropagator))
	global.SetPropagators(props)
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

func spawn() (*A, *B, *C) {
	var a = A{
		make(chan int, 1),
		make(chan int, 1),
		nil,
		nil,
		global.Tracer("TwoBuyer/A"),
	}
	var b = B{
		make(chan int, 1),
		make(chan int, 1),
		nil,
		nil,
		global.Tracer("TwoBuyer/B"),
	}
	var c = C{
		make(chan int, 1),
		make(chan int, 1),
		nil,
		nil,
		global.Tracer("TwoBuyer/C"),
	}
	b.a = &a
	c.a = &a
	a.b = &b
	c.b = &b
	a.c = &c
	b.c = &c
	return &a, &b, &c
}

func RunAll() {
	shutdown := initJaegerTracer()
	defer shutdown()

	var wg sync.WaitGroup
	tracer := global.Tracer("TwoBuyer")
	ctx := context.Background()
	var a, b, c = spawn()
	ctx, span := tracer.Start(ctx, "TwoBuyer")
	defer span.End()
	wg.Add(3)
	go a.run(&wg, ctx)
	go b.run(&wg, ctx)
	go c.run(&wg, ctx)
	wg.Wait()
}
