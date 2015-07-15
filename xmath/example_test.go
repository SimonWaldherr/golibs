package xmath_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/xmath"
)

var float = []float64{.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 1.3e3}

func round(v float64) int {
	return xmath.Round(v)
}

func ExampleArithmetic() {
	fmt.Printf("%#v => %v", float, round(xmath.Arithmetic(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 1300} => 173
}

func ExampleMedian() {
	fmt.Printf("%#v => %v", float, round(xmath.Median(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 1300} => 13
}

func ExampleHarmonic() {
	fmt.Printf("%#v => %v", float, round(xmath.Harmonic(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 1300} => 2
}

func ExampleGeometric() {
	fmt.Printf("%#v => %v", float, round(xmath.Geometric(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 1300} => 10
}
