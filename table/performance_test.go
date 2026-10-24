package table

import (
	"fmt"
	"strings"
	"testing"
)

type benchmarkRow struct {
	Name  string
	Count int
	Price float64
}

var (
	benchmarkOutput string
	benchmarkRows   []benchmarkRow
)

func BenchmarkRender(b *testing.B) {
	data := make([]benchmarkRow, 100)
	for i := range data {
		data[i] = benchmarkRow{Name: fmt.Sprintf("row-%03d", i), Count: i, Price: float64(i) / 10}
	}
	for _, format := range []string{"html", "markdown", "ascii"} {
		b.Run(format, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				output, err := Render(data, TableOption{Format: format})
				if err != nil {
					b.Fatal(err)
				}
				benchmarkOutput = output
			}
		})
	}
}

func BenchmarkParseCSVInput(b *testing.B) {
	var input strings.Builder
	input.WriteString("Name,Count,Price\n")
	for i := 0; i < 100; i++ {
		fmt.Fprintf(&input, "row-%03d,%d,%g\n", i, i, float64(i)/10)
	}
	data := input.String()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var rows []benchmarkRow
		if err := ParseCSVInput(strings.NewReader(data), &rows); err != nil {
			b.Fatal(err)
		}
		benchmarkRows = rows
	}
}
