package xmath

import (
	"github.com/simonwaldherr/golibs/as"
	"math"
	"reflect"
	"sort"
)

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

func Prime(n int) int {
	var primeList = []int{2}
	var isPrime int = 1
	var num int = 3
	var sqrtNum int = 0
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

func Deg2Rad(deg float64) float64 {
	return (deg * math.Pi) / 180
}

func Rad2Deg(rad float64) float64 {
	return (rad * 180) / math.Pi
}

func Round(v float64) int {
	if v < 0 {
		return int(math.Ceil(v - 0.5))
	}
	return int(math.Floor(v + 0.5))
}

func Count(val interface{}) int {
	slice := reflect.ValueOf(val)
	return slice.Len()
}

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

func Min(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	out := make([]float64, c)
	for i := 0; i < c; i++ {
		out[i] = as.Float(slice.Index(i).Interface())
	}
	var min float64 = out[0]
	for _, value := range out {
		if value < min {
			min = value
		}
	}
	return min
}

func Max(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	out := make([]float64, c)
	for i := 0; i < c; i++ {
		out[i] = as.Float(slice.Index(i).Interface())
	}
	var max float64 = out[0]
	for _, value := range out {
		if value > max {
			max = value
		}
	}
	return max
}

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
	} else {
		return (out[c/2] + out[c/2]) / 2
	}
}

func Arithmetic(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	out := make([]float64, c)
	for i := 0; i < c; i++ {
		out[i] = as.Float(slice.Index(i).Interface())
	}
	return (Sum(out) / float64(len(out)))
}

func Harmonic(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	var s float64 = 0
	for i := 0; i < c; i++ {
		s = s + 1/as.Float(slice.Index(i).Interface())
	}
	return (float64(c) * 1 / s)
}

func Geometric(val interface{}) float64 {
	slice := reflect.ValueOf(val)
	c := slice.Len()
	var m float64 = 1
	for i := 0; i < c; i++ {
		m = m * as.Float(slice.Index(i).Interface())
	}
	return float64(math.Pow(float64(m), 1/float64(c)))
}
