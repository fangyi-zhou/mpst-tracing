package main

import (
	"fmt"
	"math/rand"
	"sync"
)

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

func (a *A) run(wg *sync.WaitGroup) {
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
func (b *B) run(wg *sync.WaitGroup) {
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

func (c *C) run(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Running C")
	// Receive a quote
	var quote = <-c.bc
	// Propose a share
	var share = quote/2 + rand.Intn(10) - 5
	fmt.Println("C: Proposing share", share)
	c.sendA(share)
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
	var wg sync.WaitGroup
	var a, b, c = spawn()
	wg.Add(3)
	go a.run(&wg)
	go b.run(&wg)
	go c.run(&wg)
	wg.Wait()
}
