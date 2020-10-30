package twobuyer

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"math/rand"
	"sync"
)

func (a *A) run(wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := context.Background()
	var span trace.Span
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
func (b *B) run(wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := context.Background()
	var span trace.Span
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

func (c *C) run(wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := context.Background()
	var span trace.Span
	ctx, span = c.tracer.Start(ctx, "TwoBuyer Endpoint C")
	defer span.End()
	// Receive a quote
	var quote = c.recvB(ctx, "quote")
	// Propose a share
	var share = quote/2 + rand.Intn(10) - 5
	fmt.Println("C: Proposing share", share)
	c.sendA(ctx, "share", share)
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
	shutdown := initOtlpTracer()
	defer shutdown()

	var wg sync.WaitGroup
	var a, b, c = spawn()
	wg.Add(3)
	go a.run(&wg)
	go b.run(&wg)
	go c.run(&wg)
	wg.Wait()
}
