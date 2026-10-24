package xio

import (
	"bytes"
	"io"
	"testing"
)

var (
	benchmarkData    []byte
	benchmarkWritten int
)

func BenchmarkReadLimit(b *testing.B) {
	data := bytes.Repeat([]byte("x"), 4096)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkData, err = ReadLimit(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteFull(b *testing.B) {
	data := bytes.Repeat([]byte("x"), 4096)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkWritten, err = WriteFull(io.Discard, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
