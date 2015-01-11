package stack

import (
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
