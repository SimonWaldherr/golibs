package xmath

import "testing"

var f = []float64{.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.000003}
var i = []int64{3, 1, 4, 1, 5, 9, 2, 6, 5}
var x = []interface{}{float64(3.3), int8(5), string("42"), int64(1)}

func Test_Sqrt(t *testing.T) {
	if Sqrt(9) != 3 {
		t.Fatalf("Sqrt Test failed")
	}
}

func Test_Prime(t *testing.T) {
	if Prime(99) != 523 {
		t.Fatalf("Prime Test failed")
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
	if Min(f) != 0.5 {
		t.Fatalf("MinFloat Test failed")
	}
	if Min(i) != 1 {
		t.Fatalf("MinInt Test failed")
	}
}

func Test_Max(t *testing.T) {
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
	if (float64(Median(x)) + float64(Arithmetic(x))) != 17.825 {
		t.Fatalf("AverageInterface Test failed")
	}
	if Harmonic(f) != 1.9887784588749662 {
		t.Fatalf("Harmonic Test failed")
	}
	if Geometric(f) != 5.124640442022734 {
		t.Fatalf("Geometric Test failed")
	}
}
