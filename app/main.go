package main

import (
	"fmt"
	"github.com/fangyi-zhou/mpst-tracing/app"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Two Buyer Protocol:")
	app.RunAll()
}
