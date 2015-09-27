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

func ExampleStack2() {
	var array *stack.Stack
	var x interface{}

	for i := 0; i < 2; i++ {
		if i == 0 {
			fmt.Println("Create Last In First Out Stack")
			array = stack.Lifo()
		} else {
			fmt.Println("Create First In First Out Stack")
			array = stack.Fifo()
		}

		fmt.Printf("Len: %v, Empty: %v\n", array.Len(), array.IsEmpty())
		array.Push(1)
		fmt.Printf("Len: %v, Empty: %v\n", array.Len(), array.IsEmpty())
		array.Push(2)
		fmt.Printf("Len: %v, Empty: %v\n", array.Len(), array.IsEmpty())
		x = array.Pop()
		fmt.Printf("Len: %v, Empty: %v, Value: %#v\n", array.Len(), array.IsEmpty(), x)
		x = array.Pop()
		fmt.Printf("Len: %v, Empty: %v, Value: %#v\n", array.Len(), array.IsEmpty(), x)
		x = array.Pop()
		fmt.Printf("Len: %v, Empty: %v, Value: %#v\n", array.Len(), array.IsEmpty(), x)
		array.Push(3)
		fmt.Printf("Len: %v, Empty: %v\n", array.Len(), array.IsEmpty())
		array.Push(false)
		array.Push("Lorem Ipsum Dolor sit Amet")
		fmt.Printf("Len: %v, Empty: %v\n", array.Len(), array.IsEmpty())

		for array.Len() > 0 {
			fmt.Println(as.String(array.Pop()))
		}
		fmt.Println()
	}

	// Output:
	// Create Last In First Out Stack
	// Len: 0, Empty: true
	// Len: 1, Empty: false
	// Len: 2, Empty: false
	// Len: 1, Empty: false, Value: 2
	// Len: 0, Empty: true, Value: 1
	// Len: 0, Empty: true, Value: ""
	// Len: 1, Empty: false
	// Len: 3, Empty: false
	// Lorem Ipsum Dolor sit Amet
	// false
	// 3
	//
	// Create First In First Out Stack
	// Len: 0, Empty: true
	// Len: 1, Empty: false
	// Len: 2, Empty: false
	// Len: 1, Empty: false, Value: 1
	// Len: 0, Empty: true, Value: 2
	// Len: 0, Empty: true, Value: ""
	// Len: 1, Empty: false
	// Len: 3, Empty: false
	// 3
	// false
	// Lorem Ipsum Dolor sit Amet
	//
}

func ExampleToFifo() {
	array1 := stack.Lifo()
	array1.Push(1)
	array1.Push(2)
	array1.Push(3)
	array1.Push(4)
	array2 := array1.ToFifo()
	for array2.Len() > 0 {
		fmt.Println(as.String(array2.Pop()))
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
}

func ExampleToFifoFromFifo() {
	array1 := stack.Fifo()
	array1.Push(1)
	array1.Push(2)
	array1.Push(3)
	array1.Push(4)
	array2 := array1.ToFifo()
	for array2.Len() > 0 {
		fmt.Println(as.String(array2.Pop()))
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
}

func ExampleToLifo() {
	array1 := stack.Fifo()
	array1.Push(1)
	array1.Push(2)
	array1.Push(3)
	array1.Push(4)
	array2 := array1.ToLifo()
	for array2.Len() > 0 {
		fmt.Println(as.String(array2.Pop()))
	}

	// Output:
	// 4
	// 3
	// 2
	// 1
}

func ExampleToLifoFromLifo() {
	array1 := stack.Lifo()
	array1.Push(1)
	array1.Push(2)
	array1.Push(3)
	array1.Push(4)
	array2 := array1.ToLifo()
	for array2.Len() > 0 {
		fmt.Println(as.String(array2.Pop()))
	}

	// Output:
	// 4
	// 3
	// 2
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
