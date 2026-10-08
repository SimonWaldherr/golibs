package xtime

import (
	"testing"
	"time"
)

var benchmarkTime string

func BenchmarkStrfTime(b *testing.B) {
	tm := time.Date(2026, 10, 8, 12, 34, 56, 0, time.UTC)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkTime = StrfTime("%Y-%m-%d %H:%M:%S %Z", tm)
	}
}
