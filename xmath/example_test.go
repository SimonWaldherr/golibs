package xmath_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/xmath"
)

var float = []float64{.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 1.3e2}

func round(v float64) int {
	return xmath.Round(v)
}

func ExampleCount() {
	fmt.Printf("%#v => %v", float, xmath.Count(float))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 130} => 8
}

func ExampleMax() {
	fmt.Printf("%#v => %v", float, round(xmath.Max(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 130} => 130
}

func ExampleMin() {
	fmt.Printf("%#v => %v", float, round(xmath.Min(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 130} => 1
}

func ExampleSum() {
	fmt.Printf("%#v => %v", float, round(xmath.Sum(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 130} => 217
}

func ExamplePrime() {
	fmt.Printf("%#v => %v", 42, xmath.Prime(42))

	// Output: 42 => 181
}

func ExampleArithmetic() {
	fmt.Printf("%#v => %v", float, round(xmath.Arithmetic(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 130} => 27
}

func ExampleMedian() {
	fmt.Printf("%#v => %v", float, round(xmath.Median(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 130} => 13
}

func ExampleHarmonic() {
	fmt.Printf("%#v => %v", float, round(xmath.Harmonic(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 130} => 2
}

func ExampleGeometric() {
	fmt.Printf("%#v => %v", float, round(xmath.Geometric(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 130} => 8
}
