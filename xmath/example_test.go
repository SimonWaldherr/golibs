package xmath_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/xmath"
)

var float = []float64{.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 1.3e2}
var ints = []int{10, 12, 14, 20}

func round(v float64) int {
	return xmath.Round(v)
}

func round2(v float64) float64 {
	return float64(xmath.Round(v*100)) / 100
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

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 130} => 9
}

func ExampleHarmonic() {
	fmt.Printf("%#v => %v", float, round(xmath.Harmonic(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 130} => 2
}

func ExampleGeometric() {
	fmt.Printf("%#v => %v", float, round(xmath.Geometric(float)))

	// Output: []float64{0.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.00031, 130} => 8
}

func ExampleMean() {
	fmt.Printf("Arithmetic: %#v => %v\n", ints, round2(xmath.Mean(ints, xmath.ArithmeticMean)))
	fmt.Printf("Geometric: %#v => %v\n", ints, round2(xmath.Mean(ints, xmath.GeometricMean)))
	fmt.Printf("Harmonic: %#v => %v\n", ints, round2(xmath.Mean(ints, xmath.HarmonicMean)))
	fmt.Printf("Median: %#v => %v\n", ints, round2(xmath.Mean(ints, xmath.MedianMean)))
	fmt.Printf("Rootmeansquare: %#v => %v\n", ints, round2(xmath.Mean(ints, xmath.RmsMean)))
	fmt.Printf("Default: %#v => %v\n", ints, round2(xmath.Mean(ints, xmath.Default)))

	// Output:
	// Arithmetic: []int{10, 12, 14, 20} => 14
	// Geometric: []int{10, 12, 14, 20} => 13.54
	// Harmonic: []int{10, 12, 14, 20} => 13.13
	// Median: []int{10, 12, 14, 20} => 13
	// Rootmeansquare: []int{10, 12, 14, 20} => 14.49
	// Default: []int{10, 12, 14, 20} => 14
}
