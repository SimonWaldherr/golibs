package xmath

import "testing"

var f = []float64{.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.000003}
var i = []int64{3, 1, 4, 1, 5, 9, 2, 6, 5}

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

func Test_SumFloat(t *testing.T) {
	if SumFloat(f) != 87.270003 {
		t.Fatalf("SumFloat Test failed")
	}
}

func Test_SumInt(t *testing.T) {
	if SumInt(i) != 36 {
		t.Fatalf("SumInt Test failed")
	}
}

func Test_MinFloat(t *testing.T) {
	if MinFloat(f) != 0.5 {
		t.Fatalf("MinFloat Test failed")
	}
}

func Test_MinInt(t *testing.T) {
	if MinInt(i) != 1 {
		t.Fatalf("MinInt Test failed")
	}
}

func Test_MaxFloat(t *testing.T) {
	if MaxFloat(f) != 42.000003 {
		t.Fatalf("MaxFloat Test failed")
	}
}

func Test_MaxInt(t *testing.T) {
	if MaxInt(i) != 9 {
		t.Fatalf("MaxInt Test failed")
	}
}

func Test_AverageFloat(t *testing.T) {
	if (MedianFloat(f) + AvgFloat(f)) != 16.457143285714288 {
		t.Fatalf("AverageFloat Test failed")
	}
}

func Test_AverageInt(t *testing.T) {
	if (float64(MedianInt(i)) + float64(AvgInt(i))) != 9 {
		t.Fatalf("AverageInt Test failed")
	}
}
