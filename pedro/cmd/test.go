//+build ignore

package main

import (
	"fmt"
	"github.com/fangyi-zhou/mpst-tracing/pedro"
	"log"
	"os"
)

// go run cmd/test.go ~/repos/Pedro/_build/default/src/pedrolib.so ~/repos/Pedro/examples/proto.pdr
func main() {
	if len(os.Args) < 3 {
		panic("test.go PEDRO_SHARED_OBJECT PEDRO_FILE")
	}
	pedroHandle, err := pedro.LoadRuntime(os.Args[1])
	if err != nil {
		log.Panicf("unable to load Pedro, err %s", err)
	}
	err = pedroHandle.LoadFromFile(os.Args[2])
	if err != nil {
		log.Panicf("%s", err)
	}
	transitions1 := pedroHandle.GetEnabledTransitions()
	fmt.Println("Enabled transitions:")
	for _, t := range transitions1 {
		fmt.Println(t)
	}
	err = pedroHandle.DoTransition("P!Q<m1>");
	if err != nil {
		log.Panicf("%s", err)
	}
	fmt.Println("Successfully performed transition P!Q<m1>")
	transitions2 := pedroHandle.GetEnabledTransitions()
	fmt.Println("Enabled transitions:")
	for _, t := range transitions2 {
		fmt.Println(t)
	}
	defer pedroHandle.Close()
}
