package set

import "testing"

var (
	benchmarkFound bool
	benchmarkSet   Set[int]
)

func BenchmarkSetHas(b *testing.B) {
	values := make([]int, 1024)
	for i := range values {
		values[i] = i
	}
	s := New(values...)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkFound = s.Has(i & 1023)
	}
}

func BenchmarkSetUnion(b *testing.B) {
	var left, right Set[int]
	for i := 0; i < 512; i++ {
		left.Add(i)
		right.Add(i + 256)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSet = left.Union(right)
	}
}

func BenchmarkSetUnionEmpty(b *testing.B) {
	benchmarkEmptyOperand(b, true)
}

func BenchmarkSetDifferenceEmpty(b *testing.B) {
	benchmarkEmptyOperand(b, false)
}

func benchmarkEmptyOperand(b *testing.B, union bool) {
	values := make([]int, 1024)
	for i := range values {
		values[i] = i
	}
	populated := New(values...)
	var empty Set[int]
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if union {
			benchmarkSet = empty.Union(populated)
		} else {
			benchmarkSet = populated.Difference(empty)
		}
	}
}
