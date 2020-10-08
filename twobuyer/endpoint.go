package twobuyer

import (
	"context"
	"go.opentelemetry.io/otel/api/trace"
)

type A struct {
	ba     chan int
	ca     chan int
	b      *B
	c      *C
	tracer trace.Tracer
}
type B struct {
	ab     chan int
	cb     chan int
	a      *A
	c      *C
	tracer trace.Tracer
}
type C struct {
	ac     chan int
	bc     chan int
	a      *A
	b      *B
	tracer trace.Tracer
}

func (a *A) sendB(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = a.tracer.Start(ctx, "Send B "+label)
	defer span.End()
	a.b.ab <- v
}

func (a *A) recvB(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = a.tracer.Start(ctx, "Recv B "+label)
	defer span.End()
	return <-a.ba
}

func (a *A) sendC(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = a.tracer.Start(ctx, "Send C "+label)
	defer span.End()
	a.c.ac <- v
}

func (a *A) recvC(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = a.tracer.Start(ctx, "Recv C "+label)
	defer span.End()
	return <-a.ca
}

func (b *B) sendA(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = b.tracer.Start(ctx, "Send A "+label)
	defer span.End()
	b.a.ba <- v
}

func (b *B) recvA(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = b.tracer.Start(ctx, "Recv A "+label)
	defer span.End()
	return <-b.ab
}

func (b *B) sendC(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = b.tracer.Start(ctx, "Send C "+label)
	defer span.End()
	b.c.bc <- v
}

func (b *B) recvC(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = b.tracer.Start(ctx, "Recv C "+label)
	defer span.End()
	return <-b.cb
}

func (c *C) sendA(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = c.tracer.Start(ctx, "Send A "+label)
	defer span.End()
	c.a.ca <- v
}

func (c *C) recvA(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = c.tracer.Start(ctx, "Recv A "+label)
	defer span.End()
	return <-c.ac
}

func (c *C) sendB(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = c.tracer.Start(ctx, "Send B "+label)
	defer span.End()
	c.b.cb <- v
}

func (c *C) recvB(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = c.tracer.Start(ctx, "Recv B "+label)
	defer span.End()
	return <-c.bc
}