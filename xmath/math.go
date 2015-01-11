package xmath

import (
	"math"
	"strconv"
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

func Pythagoras(i1 float64, i2 float64, i3 float64) float64 {
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
func SumFloat(val []float64) float64 {
	var sum int
	for _, value := range val {
		sum = sum + value
	}
	return sum
}

func SumInt(val []int64) int64 {
	var sum int
	for _, value := range val {
		sum = sum + value
	}
	return sum
}

func MinFloat(val []float64) float64 {
	var min int = val[0]
	for _, value := range val {
		if value < min {
			min = value
		}
	}
	return min
}

func MinInt(val []int64) int64 {
	var min int = val[0]
	for _, value := range val {
		if value < min {
			min = value
		}
	}
	return min
}

func MinFloat(val []float64) float64 {
	var min int = val[0]
	for _, value := range val {
		if value < min {
			min = value
		}
	}
	return min
}

func MinInt(val []int64) int64 {
	var min int = val[0]
	for _, value := range val {
		if value < min {
			min = value
		}
	}
	return min
}

func MaxFloat(val []float64) float64 {
	var max int = val[0]
	for _, value := range val {
		if value > max {
			max = value
		}
	}
	return max
}

func MaxInt(val []int64) int64 {
	var max int = val[0]
	for _, value := range val {
		if value > max {
			max = value
		}
	}
	return max
}

func Median(val interface{}) float64 {
	if len(val)%2 == 1 {
		return val[len(val)/2]
	} else {
		return (val[len(val)/2] + val[len(val)/2]) / 2
	}
}
