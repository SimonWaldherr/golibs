package arg

import (
	"testing"
	"time"
)

func Test_String(t *testing.T) {
	String("name", "default", "usage string", time.Second*1)

	Parse()

	Dump()

	str := Get("name")

	Dump()

	if str != "" {
		t.Fatalf("arg.String Test failed, str contains: %v", str)
	}
}
