package twobuyer

import (
	"context"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

var (
	actionKey   = label.Key("mpst/action")
	msgLabelKey = label.Key("mpst/msgLabel")
	partnerKey  = label.Key("mpst/partner")
)

type A struct {
	sa     chan int
	ba     chan int
	s      *S
	b      *B
	tracer trace.Tracer
}
type S struct {
	as     chan int
	bs     chan int
	a      *A
	b      *B
	tracer trace.Tracer
}
type B struct {
	ab     chan int
	sb     chan int
	a      *A
	s      *S
	tracer trace.Tracer
}

func (a *A) sendS(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = a.tracer.Start(ctx, "Send S "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("S"), actionKey.String("send"))
	defer span.End()
	a.s.as <- v
}

func (a *A) recvS(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = a.tracer.Start(ctx, "Recv S "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("S"), actionKey.String("recv"))
	defer span.End()
	return <-a.sa
}

func (a *A) sendB(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = a.tracer.Start(ctx, "Send B "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("B"), actionKey.String("send"))
	defer span.End()
	a.b.ab <- v
}

func (a *A) recvB(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = a.tracer.Start(ctx, "Recv B "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("B"), actionKey.String("recv"))
	defer span.End()
	return <-a.ba
}

func (s *S) sendA(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "Send A "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("A"), actionKey.String("send"))
	defer span.End()
	s.a.sa <- v
}

func (s *S) recvA(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "Recv A "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("A"), actionKey.String("recv"))
	defer span.End()
	return <-s.as
}

func (s *S) sendB(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "Send B "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("B"), actionKey.String("send"))
	defer span.End()
	s.b.sb <- v
}

func (s *S) recvB(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = s.tracer.Start(ctx, "Recv B "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("B"), actionKey.String("recv"))
	defer span.End()
	return <-s.bs
}

func (b *B) sendA(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = b.tracer.Start(ctx, "Send A "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("A"), actionKey.String("send"))
	defer span.End()
	b.a.ba <- v
}

func (b *B) recvA(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = b.tracer.Start(ctx, "Recv A "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("A"), actionKey.String("recv"))
	defer span.End()
	return <-b.ab
}

func (b *B) sendS(ctx context.Context, label string, v int) {
	var span trace.Span
	ctx, span = b.tracer.Start(ctx, "Send S "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("S"), actionKey.String("send"))
	defer span.End()
	b.s.bs <- v
}

func (b *B) recvS(ctx context.Context, label string) int {
	var span trace.Span
	ctx, span = b.tracer.Start(ctx, "Recv S "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("S"), actionKey.String("recv"))
	defer span.End()
	return <-b.sb
}
