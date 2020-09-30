package main

import "go.opentelemetry.io/otel/api/trace"

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

func (a *A) sendB(label string, v int) {
	a.b.ab <- v
}

func (a *A) sendC(label string, v int) {
	a.c.ac <- v
}

func (b *B) sendA(label string, v int) {
	b.a.ba <- v
}

func (b *B) sendC(label string, v int) {
	b.c.bc <- v
}

func (c *C) sendA(label string, v int) {
	c.a.ca <- v
}

func (c *C) sendB(label string, v int) {
	c.b.cb <- v
}
