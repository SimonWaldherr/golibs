package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var benchmarkError error

func BenchmarkRetrySuccess(b *testing.B) {
	ctx := context.Background()
	fn := func(context.Context) error { return nil }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkError = Do(ctx, 3, 0, fn)
	}
}

func BenchmarkRetryExhausted(b *testing.B) {
	ctx := context.Background()
	err := errors.New("retryable")
	fn := func(context.Context) error { return err }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkError = Do(ctx, 3, 0, fn)
	}
}

func BenchmarkRetryDelayed(b *testing.B) {
	ctx := context.Background()
	wantErr := errors.New("retryable")
	fn := func(context.Context) error { return wantErr }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkError = Do(ctx, 5, time.Nanosecond, fn)
		if benchmarkError != wantErr {
			b.Fatal(benchmarkError)
		}
	}
}
