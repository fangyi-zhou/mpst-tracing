//+build ignore

package main

import (
	"fmt"
	"github.com/fangyi-zhou/mpst-tracing/pedro"
	"log"
	"os"
)

// go run cmd/test.go ~/repos/Pedro/_build/default/src/pedrolib.so ~/repos/Pedro/examples/Simple.scr TwoBuyer
func main() {
	if len(os.Args) < 3 {
		panic("test.go PEDRO_SHARED_OBJECT PEDRO_FILE")
	}
	pedroHandle, err := pedro.LoadRuntime(os.Args[1])
	if err != nil {
		log.Panicf("unable to load Pedro, err %s", err)
	}
	err = pedroHandle.ImportNuscrFile(os.Args[2], os.Args[3])
	if err != nil {
		log.Panicf("%s", err)
	}
	transitions1 := pedroHandle.GetEnabledTransitions()
	fmt.Println("Enabled transitions:")
	for _, t := range transitions1 {
		fmt.Println(t)
	}
	err = pedroHandle.DoTransition("C!A<share>");
	if err != nil {
		log.Panicf("%s", err)
	}
	fmt.Println("Successfully performed transition C!A<share>")
	transitions2 := pedroHandle.GetEnabledTransitions()
	fmt.Println("Enabled transitions:")
	for _, t := range transitions2 {
		fmt.Println(t)
	}
	err = pedroHandle.DoTransition("C?A<share>");
	if err != nil {
		log.Panicf("%s", err)
	}
	fmt.Println("Successfully performed transition C?A<share>")
	transitions3 := pedroHandle.GetEnabledTransitions()
	fmt.Println("Enabled transitions:")
	for _, t := range transitions3 {
		fmt.Println(t)
	}
	pedroHandle.Close()
}
