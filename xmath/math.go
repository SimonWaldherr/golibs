// xmath is a package with math functions.
// Besides a few standard formulas it contains various mean algorithms.
package xmath

import (
	"math"
	"reflect"
	"simonwaldherr.de/go/golibs/as"
	"sort"
)

// Sqrt calculates the square root of n.
func Sqrt(n int64) int64 {
	var t int64
	var b int64
	var r int64
	t = int64(n)
	p := int64(1 << 30)
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

// Round returns a rounded float64 from a float64
// with d digits after the point. It rounds via
// "Round half away from zero".
func FloatRound(v float64, d int) float64 {
	pow := math.Pow(10, float64(d))
	return float64(Round(v*pow)) / pow
}

// Count returns the length of any slice (like len()).
func Count(val interface{}) int {
	slice := reflect.ValueOf(val)
	return slice.Len()
}

// Sum returns the sum from a slice of Values as float64.
// It uses "as" (simonwaldherr.de/go/golibs/as) to
// convert given values to floats.
func Sum(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	out := make([]float64, c)
	for i := 0; i < c; i++ {
		out[i] = as.Float(slice.Index(i).Interface())
	}
	var sum float64
	for _, value := range out {
		sum = sum + value
	}
	return sum
}

// Min returns the smallest number from a slice of Values as float64.
// It uses "as" (simonwaldherr.de/go/golibs/as) to
// convert given values to floats.
func Min(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	out := make([]float64, c)
	for i := 0; i < c; i++ {
		out[i] = as.Float(slice.Index(i).Interface())
	}
	min := out[0]
	for _, value := range out {
		if value < min {
			min = value
		}
	}
	return min
}

// Max returns the greatest number from a slice of Values as float64.
// It uses "as" (simonwaldherr.de/go/golibs/as) to
// convert given values to floats.
func Max(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	out := make([]float64, c)
	for i := 0; i < c; i++ {
		out[i] = as.Float(slice.Index(i).Interface())
	}
	max := out[0]
	for _, value := range out {
		if value > max {
			max = value
		}
	}
	return max
}

type Meantype int

const (
	ArithmeticMean Meantype = iota
	GeometricMean
	HarmonicMean
	MedianMean
	RmsMean
	Default
)

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
// lowest value to highest value and picking the middle one.
// It uses "as" (simonwaldherr.de/go/golibs/as) to
// convert given values to floats.
func Median(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	out := make([]float64, c)
	for i := 0; i < c; i++ {
		out[i] = as.Float(slice.Index(i).Interface())
	}
	sort.Float64s(out)
	if c%2 == 1 {
		return out[c/2]
	}
	return (out[c/2] + out[c/2-1]) / 2
}

// Arithmetic returns the arithmetic mean from a slice of Values as float64.
// The arithmetic mean or simply the mean or average when the context is clear,
// is the sum of a list of numbers divided by the number of numbers
// in the list.
// It uses "as" (simonwaldherr.de/go/golibs/as) to
// convert given values to floats.
func Arithmetic(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	out := make([]float64, c)
	for i := 0; i < c; i++ {
		out[i] = as.Float(slice.Index(i).Interface())
	}
	return (Sum(out) / float64(len(out)))
}

// Rootmeansquare returns the root mean square from a slice of Values as float64.
// The root mean square is the root value of the sum of the squared value of a
// list of numbers divided by the number of numbers in the list.
// It uses "as" (simonwaldherr.de/go/golibs/as) to
// convert given values to floats.
func Rootmeansquare(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	out := make([]float64, c)
	for i := 0; i < c; i++ {
		out[i] = math.Pow(as.Float(slice.Index(i).Interface()), 2)
	}
	return math.Sqrt(Sum(out) / float64(len(out)))
}

// Harmonic returns the harmonic mean from a slice of Values as float64.
// It uses "as" (simonwaldherr.de/go/golibs/as) to
// convert given values to floats.
func Harmonic(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	s := float64(0)
	for i := 0; i < c; i++ {
		s = s + 1/as.Float(slice.Index(i).Interface())
	}
	return (float64(c) * 1 / s)
}

// Geometric returns the geometric mean from a slice of Values as float64.
// The geometric mean is a type of mean or average, which indicates the central
// tendency or typical value of a set of numbers by using the product of their
// values (as opposed to the arithmetic mean which uses their sum). The
// geometric mean is defined as the nth root of the product of n numbers.
// It uses "as" (simonwaldherr.de/go/golibs/as) to
// convert given values to floats.
func Geometric(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	var m float64 = 1
	for i := 0; i < c; i++ {
		m = m * as.Float(slice.Index(i).Interface())
	}
	return float64(math.Pow(float64(m), 1/float64(c)))
}

// Even tells if a number is even
func Even(number int) bool {
	return number%2 == 0
}

// Odd tells if a number is odd
func Odd(number int) bool {
	return !Even(number)
}
