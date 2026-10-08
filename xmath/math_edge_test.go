package xmath

import (
	"math"
	"testing"
)

func TestSqrtFullInt64Range(t *testing.T) {
	tests := []struct {
		n, want int64
	}{
		{0, 0}, {1, 1}, {2, 1}, {3, 1}, {4, 2}, {9, 3},
		{1<<32 - 1, 65535}, {1 << 32, 65536}, {1<<32 + 1, 65536},
		{1<<60 - 1, 1<<30 - 1}, {1 << 60, 1 << 30}, {1<<60 + 1, 1 << 30},
		{3037000499 * 3037000499, 3037000499}, {math.MaxInt64, 3037000499},
	}
	for _, tt := range tests {
		if got := Sqrt(tt.n); got != tt.want {
			t.Errorf("Sqrt(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestSqrtNegativePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Sqrt(-1) did not panic")
		}
	}()
	Sqrt(-1)
}

func TestFloatRoundWithoutIntegerOverflow(t *testing.T) {
	tests := []struct {
		value  float64
		digits int
		want   float64
	}{
		{1e20, 0, 1e20}, {-1e20, 0, -1e20}, {1e20, 2, 1e20},
		{1.25, 1, 1.3}, {-1.25, 1, -1.3}, {125, -1, 130}, {-125, -1, -130},
	}
	for _, tt := range tests {
		if got := FloatRound(tt.value, tt.digits); got != tt.want {
			t.Errorf("FloatRound(%g, %d) = %g, want %g", tt.value, tt.digits, got, tt.want)
		}
	}
}

func TestFloatRoundNonFiniteValues(t *testing.T) {
	for _, value := range []float64{math.Inf(1), math.Inf(-1)} {
		if got := FloatRound(value, 2); got != value {
			t.Errorf("FloatRound(%v, 2) = %v", value, got)
		}
	}
	if got := FloatRound(math.NaN(), 2); !math.IsNaN(got) {
		t.Errorf("FloatRound(NaN, 2) = %v", got)
	}
}
