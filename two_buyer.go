package main

import (
	"fmt"
	"sync"
)

type A struct {}
type B struct {}
type C struct {}

func (*A) run (wg *sync.WaitGroup, b *B, c *C) {
	defer wg.Done()
	fmt.Println("Running A")
}

func (*B) run (wg *sync.WaitGroup, a *A, c *C) {
	defer wg.Done()
	fmt.Println("Running B")
}

func (*C) run (wg *sync.WaitGroup, a *A, b *B) {
	defer wg.Done()
	fmt.Println("Running C")
}

func runAll() {
	var wg sync.WaitGroup
	var a = A {}
	var b = B {}
	var c = C {}
	wg.Add(3)
	go a.run(&wg, &b, &c)
	go b.run(&wg, &a, &c)
	go c.run(&wg, &a, &b)
	wg.Wait()
}
