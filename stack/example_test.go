package stack_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/as"
	"simonwaldherr.de/go/golibs/stack"
)

func ExampleStack() {
	array := stack.Lifo()
	array.Push("1")
	array.Push(float64(13.37))
	array.Push(false)
	array.Push("Lorem Ipsum Dolor sit Amet")

	for array.Len() > 0 {
		fmt.Println(as.String(array.Pop()))
	}

	// Output:
	// Lorem Ipsum Dolor sit Amet
	// false
	// 13.37
	// 1
}

func ExampleRings() {
	ring := stack.Ring()
	ring.Init(4)

	for i := 0; i < 16; i++ {
		ring.Push(fmt.Sprintf("%v", i))
	}

	vals := ring.Get(2)
	for _, v := range vals {
		fmt.Println(v)
	}

	// Output:
	// 12
	// 13
	// 14
}
