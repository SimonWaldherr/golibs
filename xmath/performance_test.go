package xmath

import "testing"

var benchmarkNumber float64

func BenchmarkNumericAggregates(b *testing.B) {
	floats := make([]float64, 1024)
	integers := make([]int64, 1024)
	for i := range floats {
		floats[i] = float64(i + 1)
		integers[i] = int64(i + 1)
	}
	for _, input := range []struct {
		name string
		data interface{}
	}{{"float64", floats}, {"int64", integers}} {
		for _, operation := range []struct {
			name string
			fn   func(interface{}) float64
		}{{"sum", Sum}, {"arithmetic", Arithmetic}, {"min", Min}, {"rms", Rootmeansquare}, {"median", Median}} {
			b.Run(input.name+"/"+operation.name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					benchmarkNumber = operation.fn(input.data)
				}
			})
		}
	}
}
