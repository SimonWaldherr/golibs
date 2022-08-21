package bitmask_test

import (
	"fmt"

	"simonwaldherr.de/go/golibs/bitmask"
)

func Example() {
	i := bitmask.New(0b11111111)

	i.Set(0, false)
	fmt.Println(i)

	i.Set(3, false)
	fmt.Println(i)

	i.Set(0, true)
	fmt.Println(i)

	// Output:
	// 11111110
	// 11110110
	// 11110111
}
