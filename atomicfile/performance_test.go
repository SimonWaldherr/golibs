package atomicfile

import (
	"bytes"
	"path/filepath"
	"testing"
)

func BenchmarkWriteFile(b *testing.B) {
	path := filepath.Join(b.TempDir(), "output")
	data := bytes.Repeat([]byte("x"), 4096)
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := WriteFile(path, data, 0600); err != nil {
			b.Fatal(err)
		}
	}
}
