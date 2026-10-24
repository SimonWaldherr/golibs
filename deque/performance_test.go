package deque

import "testing"

var (
	benchmarkValue int
	benchmarkDeque Deque[int]
)

// Each iteration removes and reuses one element of a warm 1,024-element ring.
func BenchmarkDequeCycle(b *testing.B) {
	var values Deque[int]
	for i := 0; i < 1024; i++ {
		values.PushBack(i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value, _ := values.PopFront()
		values.PushBack(value)
		benchmarkValue = value
	}
}

func BenchmarkDequeGrow(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var values Deque[int]
		for item := 0; item < 1024; item++ {
			values.PushBack(item)
		}
		benchmarkDeque = values
	}
}
