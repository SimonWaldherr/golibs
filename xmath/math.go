// Package xmath provides numeric aggregates and elementary math functions.
// Besides a few standard formulas it contains various mean algorithms.
package xmath

import (
	"math"
	"reflect"
	"sort"
)

// Sqrt returns the integer part of the square root of n.
// It panics if n is negative.
func Sqrt(n int64) int64 {
	if n < 0 {
		panic("xmath: square root of a negative integer")
	}
	var t int64
	var b int64
	var r int64
	t = int64(n)
	p := int64(1 << 62)
	for p > t {
		p >>= 2
	}
	for ; p != 0; p >>= 2 {
		b = r | p
		r >>= 1
		if t >= b {
			t -= b
			r |= p
		}
	}
	return int64(r)
}

// Prime returns the nth prime number as int.
func Prime(n int) int {
	var primeList = []int{2}
	isPrime := 1
	num := 3
	sqrtNum := 0
	for len(primeList) < n {
		sqrtNum = int(Sqrt(int64(num)))
		for i := 0; i < len(primeList); i++ {
			if num%primeList[i] == 0 {
				isPrime = 0
			}
			if primeList[i] > sqrtNum {
				i = len(primeList)
			}
		}
		if isPrime == 1 {
			primeList = append(primeList, num)
		} else {
			isPrime = 1
		}
		num = num + 2
	}
	return primeList[n-1]
}

// Deg2Rad returns the rad of a deg.
func Deg2Rad(deg float64) float64 {
	return (deg * math.Pi) / 180
}

// Rad2Deg returns the deg of a rad.
func Rad2Deg(rad float64) float64 {
	return (rad * 180) / math.Pi
}

// Round returns a rounded int from a float64.
// It rounds via "Round half away from zero".
func Round(v float64) int {
	if v < 0 {
		return int(math.Ceil(v - 0.5))
	}
	return int(math.Floor(v + 0.5))
}

// FloatRound returns a rounded float64 from a float64
// with d digits after the point. It rounds via
// "Round half away from zero".
func FloatRound(v float64, d int) float64 {
	pow := math.Pow(10, float64(d))
	return math.Round(v*pow) / pow
}

// Count returns the length of any slice (like len()).
func Count(val interface{}) int {
	slice := reflect.ValueOf(val)
	return slice.Len()
}

// Sum returns the sum of the values as float64. Empty inputs return zero.
// Built-in numeric slices avoid allocation; other slices and arrays use as.Float.
func Sum(val interface{}) float64 { return reduce(val, sumReduction) }

// Min returns the smallest converted value. It panics for an empty input.
func Min(val interface{}) float64 { return reduce(val, minReduction) }

// Max returns the greatest converted value. It panics for an empty input.
func Max(val interface{}) float64 { return reduce(val, maxReduction) }

// Meantype selects an algorithm for Mean.
type Meantype int

const (
	// ArithmeticMean selects the sum divided by the number of values.
	ArithmeticMean Meantype = iota
	// GeometricMean selects the nth root of the product.
	GeometricMean
	// HarmonicMean selects the reciprocal of the average reciprocals.
	HarmonicMean
	// MedianMean selects the middle value of the sorted input.
	MedianMean
	// RmsMean selects the root mean square.
	RmsMean
	// Default selects ArithmeticMean.
	Default
)

// Mean returns the requested mean of a slice or array. Unknown modes use Arithmetic.
func Mean(val interface{}, t Meantype) float64 {
	switch t {
	case ArithmeticMean:
		return Arithmetic(val)
	case GeometricMean:
		return Geometric(val)
	case HarmonicMean:
		return Harmonic(val)
	case MedianMean:
		return Median(val)
	case RmsMean:
		return Rootmeansquare(val)
	}
	return Arithmetic(val)
}

// Median returns the median from a slice of Values as float64.
// The median is the numerical value separating the higher half
// of a data sample from the lower half. The median of a list of
// numbers can be found by arranging all the observations from
// lowest value to highest value and picking the middle one. The input is not
// modified. Median panics for an empty input.
// It uses "as" (simonwaldherr.de/go/golibs/as) to
// convert given values to floats.
func Median(val interface{}) float64 {
	out := numericCopy(val)
	sort.Float64s(out)
	c := len(out)
	if c%2 == 1 {
		return out[c/2]
	}
	return (out[c/2] + out[c/2-1]) / 2
}

// Arithmetic returns the arithmetic mean of a slice or array. Empty inputs return NaN.
// The arithmetic mean or simply the mean or average when the context is clear,
// is the sum of a list of numbers divided by the number of numbers
// in the list.
// It uses "as" (simonwaldherr.de/go/golibs/as) to
// convert given values to floats.
func Arithmetic(val interface{}) float64 { return reduce(val, arithmeticReduction) }

// Rootmeansquare returns the square root of the mean of squared values.
// Empty inputs return NaN. Built-in numeric slices are processed without allocation.
func Rootmeansquare(val interface{}) float64 { return reduce(val, rmsReduction) }

// Harmonic returns the harmonic mean of the converted values.
// Empty inputs return NaN. A zero value yields zero for positive inputs.
func Harmonic(val interface{}) float64 { return reduce(val, harmonicReduction) }

// Geometric returns the geometric mean of the converted values.
// It uses the product of the values and returns 1 for an empty input.
func Geometric(val interface{}) float64 { return reduce(val, geometricReduction) }

// Even tells if a number is even
func Even(number int) bool {
	return number%2 == 0
}

// Odd tells if a number is odd
func Odd(number int) bool {
	return !Even(number)
}
