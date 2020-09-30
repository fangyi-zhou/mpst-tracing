package main

type A struct {
	ba chan int
	ca chan int
	b  *B
	c  *C
}
type B struct {
	ab chan int
	cb chan int
	a  *A
	c  *C
}
type C struct {
	ac chan int
	bc chan int
	a  *A
	b  *B
}

func (a *A) sendB(v int) {
	a.b.ab <- v
}

func (a *A) sendC(v int) {
	a.c.ac <- v
}

func (b *B) sendA(v int) {
	b.a.ba <- v
}

func (b *B) sendC(v int) {
	b.c.bc <- v
}

func (c *C) sendA(v int) {
	c.a.ca <- v
}

func (c *C) sendB(v int) {
	c.b.cb <- v
}
