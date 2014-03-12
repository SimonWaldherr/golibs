package xmath

import "math"
import "strconv"

func Sqrt(n int) int {
	var t uint
	var b uint
	var r uint
	t = uint(n)
	p := uint(1 << 30)
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
	return int(r)
}

func Prime(n int) int {
	var primeList = []int{2}
	var isPrime int = 1
	var num int = 3
	var sqrtNum int = 0
	for len(primeList) < n {
		sqrtNum = Sqrt(num)
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

func pythagoras(i1 float64, i2 float64, i3 float64) float64 {
		switch "?" {
		case i1:
			b, _ := strconv.ParseFloat(i2, 0)
			c, _ := strconv.ParseFloat(i3, 0)
			return math.Sqrt(q(float64(c)) - q(float64(b)))
		case i2:
			a, _ := strconv.ParseFloat(i1, 0)
			c, _ := strconv.ParseFloat(i3, 0)
			return math.Sqrt(q(float64(c)) - q(float64(a)))
		case i3:
			a, _ := strconv.ParseFloat(i1, 0)
			b, _ := strconv.ParseFloat(i2, 0)
			return math.Sqrt(q(float64(a)) + q(float64(b)))
		default:
			return 0
		}
	}
}
