package arg_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/arg"
	"time"
)

func Example() {
	arg.String("name", "default", "usage string", time.Second*5)
	arg.Parse()
	fmt.Println("\ninput:", arg.Get("name"))
}
