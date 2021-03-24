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

func (a *A) SendS(ctx context.Context, label string, v int) {
	var span trace.Span
	_, span = a.tracer.Start(ctx, "Send S "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("S"), actionKey.String("Send"))
	defer span.End()
	a.s.as <- v
}

func (a *A) RecvS(ctx context.Context, label string) int {
	var span trace.Span
	_, span = a.tracer.Start(ctx, "Recv S "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("S"), actionKey.String("Recv"))
	defer span.End()
	return <-a.sa
}

func (a *A) SendB(ctx context.Context, label string, v int) {
	var span trace.Span
	_, span = a.tracer.Start(ctx, "Send B "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("B"), actionKey.String("Send"))
	defer span.End()
	a.b.ab <- v
}

func (a *A) RecvB(ctx context.Context, label string) int {
	var span trace.Span
	_, span = a.tracer.Start(ctx, "Recv B "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("B"), actionKey.String("Recv"))
	defer span.End()
	return <-a.ba
}

func (s *S) SendA(ctx context.Context, label string, v int) {
	var span trace.Span
	_, span = s.tracer.Start(ctx, "Send A "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("A"), actionKey.String("Send"))
	defer span.End()
	s.a.sa <- v
}

func (s *S) RecvA(ctx context.Context, label string) int {
	var span trace.Span
	_, span = s.tracer.Start(ctx, "Recv A "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("A"), actionKey.String("Recv"))
	defer span.End()
	return <-s.as
}

func (s *S) SendB(ctx context.Context, label string, v int) {
	var span trace.Span
	_, span = s.tracer.Start(ctx, "Send B "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("B"), actionKey.String("Send"))
	defer span.End()
	s.b.sb <- v
}

func (s *S) RecvB(ctx context.Context, label string) int {
	var span trace.Span
	_, span = s.tracer.Start(ctx, "Recv B "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("B"), actionKey.String("Recv"))
	defer span.End()
	return <-s.bs
}

func (b *B) SendA(ctx context.Context, label string, v int) {
	var span trace.Span
	_, span = b.tracer.Start(ctx, "Send A "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("A"), actionKey.String("Send"))
	defer span.End()
	b.a.ba <- v
}

func (b *B) RecvA(ctx context.Context, label string) int {
	var span trace.Span
	_, span = b.tracer.Start(ctx, "Recv A "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("A"), actionKey.String("Recv"))
	defer span.End()
	return <-b.ab
}

func (b *B) SendS(ctx context.Context, label string, v int) {
	var span trace.Span
	_, span = b.tracer.Start(ctx, "Send S "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("S"), actionKey.String("Send"))
	defer span.End()
	b.s.bs <- v
}

func (b *B) RecvS(ctx context.Context, label string) int {
	var span trace.Span
	_, span = b.tracer.Start(ctx, "Recv S "+label)
	span.SetAttributes(msgLabelKey.String(label), partnerKey.String("S"), actionKey.String("Recv"))
	defer span.End()
	return <-b.sb
}
