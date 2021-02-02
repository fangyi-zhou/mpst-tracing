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
	// Send query to S
	var query = rand.Intn(100)
	fmt.Println("A: Sending query", query)
	a.sendS(ctx, "query", query)
	// Receive a quote
	var quote = a.recvS(ctx, "quote")
	var otherShare = a.recvB(ctx, "share")
	if otherShare*2 >= quote {
		// 1 stands for ok
		a.sendS(ctx, "buy", 1)
	} else {
		a.sendS(ctx, "buy", 0)
	}
}
func (s *S) run(wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := context.Background()
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "TwoBuyer Endpoint S")
	defer span.End()
	// Receive a query
	var query = s.recvA(ctx, "query")
	// Send a quote
	var quote = query * 2
	fmt.Println("S: Sending quote", quote)
	s.sendA(ctx, "quote", quote)
	s.sendB(ctx, "quote", quote)
	var decision = s.recvA(ctx, "buy")
	if decision == 1 {
		fmt.Println("Succeed!")
	} else {
		fmt.Println("Failed to succeed!")
	}
}

func (b *B) run(wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := context.Background()
	var span trace.Span
	ctx, span = b.tracer.Start(ctx, "TwoBuyer Endpoint B")
	defer span.End()
	// Receive a quote
	var quote = b.recvS(ctx, "quote")
	// Propose a share
	var share = quote/2 + rand.Intn(10) - 5
	fmt.Println("B: Proposing share", share)
	b.sendA(ctx, "share", share)
}

func spawn() (*A, *B, *S) {
	var a = A{
		make(chan int, 1),
		make(chan int, 1),
		nil,
		nil,
		global.Tracer("TwoBuyer/A"),
	}
	var s = S{
		make(chan int, 1),
		make(chan int, 1),
		nil,
		nil,
		global.Tracer("TwoBuyer/S"),
	}
	var b = B{
		make(chan int, 1),
		make(chan int, 1),
		nil,
		nil,
		global.Tracer("TwoBuyer/B"),
	}
	s.a = &a
	b.a = &a
	a.s = &s
	b.s = &s
	a.b = &b
	s.b = &b
	return &a, &b, &s
}

func RunAll() {
	shutdown := initOtlpTracer()
	defer shutdown()

	var wg sync.WaitGroup
	var a, b, s = spawn()
	wg.Add(3)
	go a.run(&wg)
	go b.run(&wg)
	go s.run(&wg)
	wg.Wait()
}
