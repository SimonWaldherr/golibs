package main

import (
	"fmt"
	"simonwaldherr.de/go/golibs/arg"
	"time"
)

func main() {
	arg.StringArg("name1", "default", "enter string", time.Second*5)
	arg.StringArg("name2", "default", "enter string", time.Second*5)
	arg.StringArg("name3", "default", "enter string", time.Second*0)
	arg.Parse()
	fmt.Println("\ninput 1:", arg.Get("name1"))
	fmt.Println("\ninput 2:", arg.Get("name2"))
	fmt.Println("\ninput 3:", arg.Get("name3"))
}
