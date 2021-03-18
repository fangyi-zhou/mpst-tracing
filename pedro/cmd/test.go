//+build ignore

package main

import (
	"github.com/fangyi-zhou/mpst-tracing/pedro"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		panic("test.go PEDRO_SHARED_OBJECT PEDRO_FILE")
	}
	pedroHandle, err := pedro.LoadRuntime(os.Args[1])
	if err != nil {
		log.Panicf("unable to load Pedro, err %s", err)
	}
	pedroHandle.RunMain(os.Args[2])
	defer pedroHandle.Close()
}
