package xmath

import "testing"

var f = []float64{.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.000003}
var i = []int64{3, 1, 4, 1, 5, 9, 2, 6, 5}
var x = []interface{}{float64(3.3), int8(5), string("42"), int64(1)}
var ints = []int{10, 12, 14, 20}

func round2(v float64) float64 {
	return float64(Round(v*100)) / 100
}

func Test_Sqrt(t *testing.T) {
	if Sqrt(9) != 3 {
		t.Fatalf("Sqrt Test failed")
	}
}

func Benchmark_Sqrt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sqrt(int64(i))
	}
}

func Test_Prime(t *testing.T) {
	if Prime(99) != 523 {
		t.Fatalf("Prime Test failed")
	}
}

func Test_Deg2Rad(t *testing.T) {
	if int64(Deg2Rad(3.44)*100) != int64(6) {
		t.Fatalf("Deg2Rad Test failed")
	}
}

func Test_Rad2Deg(t *testing.T) {
	if int64(Rad2Deg(0.20944)) != int64(12) {
		t.Fatalf("Rad2Deg Test failed")
	}
}

func Test_DegAndRad(t *testing.T) {
	if Rad2Deg(Deg2Rad(42)) != float64(42) {
		t.Fatalf("DegAndRad Test failed")
	}
}

func Test_Round(t *testing.T) {
	if Round(99.999) != 100 {
		t.Fatalf("Round Test failed")
	}
	if Round(0.001) != 0 {
		t.Fatalf("Round Test failed")
	}
	if Round(-0.001) != 0 {
		t.Fatalf("Round Test failed")
	}
	if Round(-0.6) != -1 {
		t.Fatalf("Round Test failed")
	}
	i := 3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211706798214808651328230664709384460955058223172535940812848111745028410270193852110555964462294895493038196442881097566593344612847564823378678316527120190914564
	if float64(Round(i*100))/100 != 3.14 {
		t.Fatalf("Round Test failed")
	}
}

func Test_FloatRound(t *testing.T) {
	var f float64 = FloatRound(-1.15, 1)
	f += FloatRound(-0.33, 1)
	f += FloatRound(1.54, 1)
	f += FloatRound(1.55, 1)
	f += FloatRound(1.4, 0)

	if f != 2.6 {
		t.Fatalf("FloatRound Test failed")
	}
}

func Test_Count(t *testing.T) {
	if Count(f) != 7 {
		t.Fatalf("Count Test failed")
	}
}

func Test_Sum(t *testing.T) {
	if Sum(f) != 87.270003 {
		t.Fatalf("SumFloat Test failed")
	}
	if Sum(i) != 36 {
		t.Fatalf("SumInt Test failed")
	}
	if Sum(x) != 51.3 {
		t.Fatalf("SumInterface Test failed")
	}
}

func Test_Min(t *testing.T) {
	if Min([]int{3, 5, 7, -9, 11}) != -9 {
		t.Fatalf("Min Test failed")
	}
	if Min(f) != 0.5 {
		t.Fatalf("MinFloat Test failed")
	}
	if Min(i) != 1 {
		t.Fatalf("MinInt Test failed")
	}
}

func Test_Max(t *testing.T) {
	if Max([]int{1, 3, 5, 42}) != 42 {
		t.Fatalf("Max Test failed")
	}
	if Max(f) != 42.000003 {
		t.Fatalf("MaxFloat Test failed")
	}
	if Max(i) != 9 {
		t.Fatalf("MaxInt Test failed")
	}
}

func Test_Average(t *testing.T) {
	if (Median(f) + Arithmetic(f)) != 16.457143285714288 {
		t.Fatalf("AverageFloat Test failed")
	}
	if (float64(Median(i)) + float64(Arithmetic(i))) != 8 {
		t.Fatalf("AverageInt Test failed")
	}
	if (float64(Median(x)) + float64(Arithmetic(x))) != 16.975 {
		t.Fatalf("AverageInterface Test failed")
	}
	if Harmonic(f) != 1.9887784588749662 {
		t.Fatalf("Harmonic Test failed")
	}
	if Geometric(f) != 5.124640442022734 {
		t.Fatalf("Geometric Test failed")
	}
}

func Test_Mean(t *testing.T) {
	if round2(Mean(ints, ArithmeticMean)) != 14 {
		t.Fatalf("ArithmeticMean Test failed")
	}
	if round2(Mean(ints, GeometricMean)) != 13.54 {
		t.Fatalf("GeometricMean Test failed")
	}
	if round2(Mean(ints, HarmonicMean)) != 13.13 {
		t.Fatalf("HarmonicMean Test failed")
	}
	if round2(Mean(ints, MedianMean)) != 13 {
		t.Fatalf("MedianMean Test failed")
	}
	if round2(Mean(ints, RmsMean)) != 14.49 {
		t.Fatalf("RmsMean Test failed")
	}
	if round2(Mean(ints, Default)) != 14 {
		t.Fatalf("Default Mean Test failed")
	}
}

func Benchmark_Median(b *testing.B) {
	var prev int64 = 0
	for i := 0; i < b.N; i++ {
		Median([]int64{int64(i), prev})
		prev = int64(i)
	}
}

func Benchmark_Arithmetic(b *testing.B) {
	var prev int64 = 0
	for i := 0; i < b.N; i++ {
		Arithmetic([]int64{int64(i), prev})
		prev = int64(i)
	}
}

func Benchmark_Harmonic(b *testing.B) {
	var prev int64 = 0
	for i := 0; i < b.N; i++ {
		Harmonic([]int64{int64(i), prev})
		prev = int64(i)
	}
}

func Benchmark_Geometric(b *testing.B) {
	var prev int64 = 0
	for i := 0; i < b.N; i++ {
		Geometric([]int64{int64(i), prev})
		prev = int64(i)
	}
}

func Test_Odd_Even(t *testing.T) {
	if Odd(3) != true {
		t.Fatalf("Odd Test 1 failed")
	}
	if Odd(4) != false {
		t.Fatalf("Odd Test 2 failed")
	}
	if Even(3) != false {
		t.Fatalf("Even Test 1 failed")
	}
	if Even(4) != true {
		t.Fatalf("Even Test 2 failed")
	}
}
