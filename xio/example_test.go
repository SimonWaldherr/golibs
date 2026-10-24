package xio_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"simonwaldherr.de/go/golibs/xio"
)

func ExampleReadLimit() {
	data, err := xio.ReadLimit(strings.NewReader("hello world"), 5)
	fmt.Println(string(data), errors.Is(err, xio.ErrLimitExceeded))
	// Output: hello true
}

func ExampleWriteFull() {
	var output bytes.Buffer
	n, err := xio.WriteFull(&output, []byte("hello"))
	fmt.Println(n, err, output.String())
	// Output: 5 <nil> hello
}
