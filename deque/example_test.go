package deque_test

import (
	"fmt"

	"simonwaldherr.de/go/golibs/deque"
)

func Example() {
	var queue deque.Deque[string]
	queue.PushBack("second")
	queue.PushFront("first")
	queue.PushBack("third")

	copy := queue.Clone()
	copy.PopBack()
	fmt.Println(copy.Values())
	for queue.Len() > 0 {
		value, _ := queue.PopFront()
		fmt.Println(value)
	}
	// Output:
	// [first second]
	// first
	// second
	// third
}
