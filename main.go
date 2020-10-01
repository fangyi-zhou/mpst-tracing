package main

import (
	"fmt"
	"math/rand"
	"mpst-tracing/app"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Two Buyer Protocol:")
	app.RunAll()
}
