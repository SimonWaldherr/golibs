package semaphore

import (
	"context"
	"testing"
)

func BenchmarkAcquireRelease(b *testing.B) {
	s, err := New(1)
	if err != nil {
		b.Fatal(err)
	}
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := s.Acquire(ctx); err != nil {
			b.Fatal(err)
		}
		if err := s.Release(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAcquireReleaseParallel(b *testing.B) {
	s, err := New(4)
	if err != nil {
		b.Fatal(err)
	}
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(worker *testing.PB) {
		for worker.Next() {
			if err := s.Acquire(ctx); err != nil {
				b.Error(err)
				return
			}
			if err := s.Release(); err != nil {
				b.Error(err)
				return
			}
		}
	})
}
