package stack

import (
	"fmt"
	"simonwaldherr.de/go/golibs/as"
	"strings"
	"testing"
)

func Test_Stack(t *testing.T) {
	array := Lifo()
	array.Push("1")
	array.Push(float64(13.37))
	array.Push(false)
	array.Push("Lorem Ipsum Dolor sit Amet")
	var str string
	for array.Len() > 0 {
		str += as.String(array.Pop()) + " "
	}
	str = strings.Trim(str, "\t\r\n ")
	if str != "Lorem Ipsum Dolor sit Amet false 13.37 1" {
		t.Fatalf("Stack Test failed")
	}
	if array.Pop() != "" {
		t.Fatalf("Stack Test failed")
	}
}

func Test_FiFo_IsEmpty(t *testing.T) {
	array := Fifo()
	if array.IsEmpty() != true {
		t.Fatalf("FiFo IsEmpty failed")
	}
}

func Test_FiFo_IsNotEmpty(t *testing.T) {
	array := Fifo()
	array.Push(1)
	if array.IsEmpty() == true {
		t.Fatalf("FiFo IsNotEmpty failed")
	}
}

func Test_FiFo_IsEmptyAgain(t *testing.T) {
	array := Fifo()
	array.Push(1)
	array.Pop()
	if array.IsEmpty() != true {
		t.Fatalf("FiFo IsEmptyAgain failed")
	}
}

func Test_LiFo_IsEmpty(t *testing.T) {
	array := Lifo()
	if array.IsEmpty() != true {
		t.Fatalf("LiFo IsEmpty failed")
	}
}

func Test_LiFo_IsNotEmpty(t *testing.T) {
	array := Lifo()
	array.Push(1)
	if array.IsEmpty() == true {
		t.Fatalf("LiFo IsNotEmpty failed")
	}
}

func Test_LiFo_IsEmptyAgain(t *testing.T) {
	array := Lifo()
	array.Push(1)
	array.Pop()
	if array.IsEmpty() != true {
		t.Fatalf("LiFo IsEmptyAgain failed")
	}
}

func Test_Convert(t *testing.T) {
	var str string
	array := Fifo()
	array.Push("1")
	array.Push("2")
	array.Push("3")

	str = fmt.Sprintf("%v", array.Val())
	if str != "[1 2 3]" {
		t.Fatalf("Convert 1 failed: %v", str)
	}

	array = array.ToFifo()

	str = fmt.Sprintf("%v", array.Val())
	if str != "[1 2 3]" {
		t.Fatalf("Convert 2 failed: %v", str)
	}

	array = array.ToLifo()

	str = fmt.Sprintf("%v", array.Val())
	if str != "[3 2 1]" {
		t.Fatalf("Convert 3 failed: %v", str)
	}

	array = array.ToLifo()

	str = fmt.Sprintf("%v", array.Val())
	if str != "[3 2 1]" {
		t.Fatalf("Convert 4 failed: %v", str)
	}

	array = array.ToFifo()

	str = fmt.Sprintf("%v", array.Val())
	if str != "[1 2 3]" {
		t.Fatalf("Convert 5 failed: %v", str)
	}
}

func Test_Stack_AddGet(t *testing.T) {
	array := Fifo()
	array.Add(1)
	if array.Get() != 1 {
		t.Fatalf("Stack Add/Get failed")
	}
}

func Test_Stack_Unset(t *testing.T) {
	array := Fifo()
	array.Add(true)
	array.Unset()
	if array.Len() != -1 {
		t.Fatalf("Stack Unset 1 failed")
	}
	if array.Pop() != "" {
		t.Fatalf("Stack Unset 2 failed")
	}
}

func Test_Ring(t *testing.T) {
	var ret string
	ring := Ring()
	ring.Init(8)
	ring.SetSize(9)
	for i := 0; i < 16; i++ {
		ring.Push(fmt.Sprintf("%v", i))
	}
	vals := ring.Get(2)
	for _, v := range vals {
		ret += string(v)
	}
	if ret != "121314" {
		t.Fatalf("Ring Test failed")
	}
	if ring.GetSize() != 9 {
		t.Fatalf("Ring Test failed")
	}
	a, b := ring.Pos()
	if (a + b) != 20 {
		t.Fatalf("Ring Test failed")
	}
	ret = ""
	vals2 := ring.Get(8)
	for _, v2 := range vals2 {
		ret += string(v2)
	}
	if ret != "891011121314" {
		t.Fatalf("Ring Test failed")
	}
}
