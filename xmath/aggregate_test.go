package xmath

import (
	"math"
	"reflect"
	"testing"
)

type formattedNumber int

func (formattedNumber) String() string { return "7" }

func TestAggregateConversions(t *testing.T) {
	inputs := []interface{}{
		[]float64{2, 4}, []float32{2, 4}, []int{2, 4}, []int8{2, 4},
		[]int16{2, 4}, []int32{2, 4}, []int64{2, 4}, []uint{2, 4},
		[]uint8{2, 4}, []uint16{2, 4}, []uint32{2, 4}, []uint64{2, 4},
		[2]int{2, 4}, []string{"2", "4"}, []interface{}{2, "4"},
	}
	operations := []struct {
		name string
		fn   func(interface{}) float64
		want float64
	}{
		{"sum", Sum, 6}, {"min", Min, 2}, {"max", Max, 4}, {"median", Median, 3},
		{"arithmetic", Arithmetic, 3}, {"rms", Rootmeansquare, math.Sqrt(10)},
		{"harmonic", Harmonic, 8.0 / 3}, {"geometric", Geometric, math.Sqrt(8)},
	}
	for _, input := range inputs {
		t.Run(reflect.TypeOf(input).String(), func(t *testing.T) {
			for _, operation := range operations {
				if got := operation.fn(input); math.Abs(got-operation.want) > 1e-12 {
					t.Errorf("%s = %v, want %v", operation.name, got, operation.want)
				}
			}
		})
	}
	if got := Sum([]formattedNumber{123}); got != 7 {
		t.Fatalf("named Stringer conversion = %v, want 7", got)
	}
}

func TestAggregateEmptyInputs(t *testing.T) {
	var values []float64
	if Sum(values) != 0 || Geometric(values) != 1 {
		t.Fatal("empty sum/product changed")
	}
	for _, fn := range []func(interface{}) float64{Arithmetic, Harmonic, Rootmeansquare} {
		if !math.IsNaN(fn(values)) {
			t.Fatal("empty mean must return NaN")
		}
	}
	for name, fn := range map[string]func(interface{}) float64{"min": Min, "max": Max, "median": Median} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("empty input did not panic")
				}
			}()
			fn(values)
		})
	}
}

func TestAggregateNaNAndNegativeValues(t *testing.T) {
	values := []float64{-3, -1, -2}
	if Sum(values) != -6 || Min(values) != -3 || Max(values) != -1 || Median(values) != -2 {
		t.Fatal("negative aggregates changed")
	}
	for _, input := range [][]float64{{math.NaN(), 1}, {1, math.NaN()}} {
		if !math.IsNaN(Sum(input)) {
			t.Error("sum must propagate NaN")
		}
	}
	if !math.IsNaN(Min([]float64{math.NaN(), 1})) || Min([]float64{1, math.NaN()}) != 1 {
		t.Fatal("minimum NaN comparison behavior changed")
	}
	if !math.IsNaN(Max([]float64{math.NaN(), 1})) || Max([]float64{1, math.NaN()}) != 1 {
		t.Fatal("maximum NaN comparison behavior changed")
	}
}

func TestMedianDoesNotMutateInput(t *testing.T) {
	values := []float64{4, 1, 3, 2}
	original := append([]float64(nil), values...)
	if Median(values) != 2.5 || !reflect.DeepEqual(values, original) {
		t.Fatal("median result or input preservation failed")
	}
	if Median([]float64{9, 1, 3}) != 3 {
		t.Fatal("odd median failed")
	}
}

func TestNumericReductionsDoNotAllocate(t *testing.T) {
	var values interface{} = []float64{1, 2, 3, 4}
	for name, fn := range map[string]func(interface{}) float64{
		"sum": Sum, "min": Min, "max": Max, "arithmetic": Arithmetic,
		"rms": Rootmeansquare, "harmonic": Harmonic, "geometric": Geometric,
	} {
		t.Run(name, func(t *testing.T) {
			if allocations := testing.AllocsPerRun(100, func() { benchmarkNumber = fn(values) }); allocations != 0 {
				t.Errorf("allocations = %v, want 0", allocations)
			}
		})
	}
}
