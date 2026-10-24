package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	benchmarkText string
	benchmarkPos  int
)

func BenchmarkReadUntil(b *testing.B) {
	for _, size := range []int{4096, 16384} {
		name := "4KiB"
		if size == 16384 {
			name = "16KiB"
		}
		b.Run(name, func(b *testing.B) {
			path := filepath.Join(b.TempDir(), "input.txt")
			if err := os.WriteFile(path, []byte(strings.Repeat("a", size)+"\nignored"), 0600); err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.SetBytes(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				text, delimiter, position, err := ReadUntil(path, []string{"\n", "\r"})
				if err != nil || delimiter != "\n" || position != size {
					b.Fatalf("ReadUntil: delimiter %q, position %d, error %v", delimiter, position, err)
				}
				benchmarkText, benchmarkPos = text, position
			}
		})
	}
}

func BenchmarkReadBlocks(b *testing.B) {
	path := filepath.Join(b.TempDir(), "input.txt")
	input := strings.Repeat(strings.Repeat("a", 63)+"\n", 256)
	if err := os.WriteFile(path, []byte(input), 0600); err != nil {
		b.Fatal(err)
	}
	callback := func(block string) (string, error) { return block, nil }
	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkText, err = ReadBlocks(path, []string{"\n"}, callback)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteLines(b *testing.B) {
	path := filepath.Join(b.TempDir(), "output.txt")
	lines := make([]string, 2048)
	for i := range lines {
		lines[i] = strings.Repeat("a", 127) + "\n"
	}
	b.ReportAllocs()
	b.SetBytes(int64(len(lines) * len(lines[0])))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := WriteLines(path, lines, false); err != nil {
			b.Fatal(err)
		}
	}
}
