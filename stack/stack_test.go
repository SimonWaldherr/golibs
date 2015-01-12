package stack

import (
	"fmt"
	"github.com/simonwaldherr/golibs/as"
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
}

func Test_Ring(t *testing.T) {
	var ret string
	ring := Ring()
	ring.Init()
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
}
