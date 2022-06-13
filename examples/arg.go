// +build local

package main

import (
	"fmt"
	"time"
	
	"simonwaldherr.de/go/golibs/arg"
)

func main() {
	arg.String("name", "default", "usage string", time.Second*5)
	arg.Parse()
	fmt.Println("\ninput:", arg.Get("name"))
}
