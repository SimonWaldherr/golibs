package arg_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/arg"
	"time"
)

func ExampleString() {
	arg.String("1", "default", "usage string", time.Second*5)
	arg.String("", "", "", time.Second*5)
	arg.Parse()
	fmt.Printf("1: %v\n", arg.Get("1"))

	// Output: # usage string: # : 1:
}
