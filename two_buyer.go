package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/sdk/trace"
	"log"
	"math/rand"
	"sync"
)

import "go.opentelemetry.io/otel/exporters/stdout"

func (a *A) run(wg *sync.WaitGroup, ctx *context.Context) {
	defer wg.Done()
	fmt.Println("Running A")
	// Send query to B
	var query = rand.Intn(100)
	fmt.Println("A: Sending query", query)
	a.sendB(query)
	// Receive a quote
	var quote = <-a.ba
	var otherShare = <-a.ca
	if otherShare*2 >= quote {
		// 1 stands for ok
		a.sendB(1)
	} else {
		a.sendB(0)
	}
}
func (b *B) run(wg *sync.WaitGroup, ctx *context.Context) {
	defer wg.Done()
	fmt.Println("Running B")
	// Receive a query
	var query = <-b.ab
	// Send a quote
	var quote = query * 2
	fmt.Println("B: Sending quote", quote)
	b.sendA(quote)
	b.sendC(quote)
	var decision = <-b.ab
	if decision == 1 {
		fmt.Println("Succeed!")
	} else {
		fmt.Println("Failed to succeed!")
	}
}

func (c *C) run(wg *sync.WaitGroup, ctx *context.Context) {
	defer wg.Done()
	fmt.Println("Running C")
	// Receive a quote
	var quote = <-c.bc
	// Propose a share
	var share = quote/2 + rand.Intn(10) - 5
	fmt.Println("C: Proposing share", share)
	c.sendA(share)
}

var tp *trace.TracerProvider

// https://github.com/open-telemetry/opentelemetry-go/blob/master/example/namedtracer/main.go
func initTracer() func() {
	exp, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Panicf("failed to initialise stdout exported %v\n", err)
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

func spawn() (*A, *B, *C) {
	var a = A{
		make(chan int, 1),
		make(chan int, 1),
		nil,
		nil,
	}
	var b = B{
		make(chan int, 1),
		make(chan int, 1),
		nil,
		nil,
	}
	var c = C{
		make(chan int, 1),
		make(chan int, 1),
		nil,
		nil,
	}
	b.a = &a
	c.a = &a
	a.b = &b
	c.b = &b
	a.c = &c
	b.c = &c
	return &a, &b, &c
}

func runAll() {
	shutdown := initTracer()
	defer shutdown()

	var wg sync.WaitGroup
	tracer := tp.Tracer("TwoBuyer")
	ctx := context.Background()
	var a, b, c = spawn()
	ctx, span := tracer.Start(ctx, "TwoBuyer")
	defer span.End()
	wg.Add(3)
	go a.run(&wg, &ctx)
	go b.run(&wg, &ctx)
	go c.run(&wg, &ctx)
	wg.Wait()
}
