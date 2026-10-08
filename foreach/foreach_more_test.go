package foreach

import "testing"

func Test_ArrayIterates(t *testing.T) {
	type S struct {
		A int
		B string
		C float64
	}
	s := S{A: 1, B: "x", C: 3.14}
	calls := 0
	Array(s, func(i int, value interface{}, depth int) {
		calls++
	})
	if calls != 3 {
		t.Fatalf("Array did not iterate struct fields, calls=%d", calls)
	}
}
