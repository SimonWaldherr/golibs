package graphics

import (
	"image"
	"testing"
)

var benchmarkImage image.Image

func BenchmarkImageFilters(b *testing.B) {
	img := makeTestImage(256, 256)
	for _, filter := range []struct {
		name string
		fn   func(image.Image) image.Image
	}{
		{"grayscale", Grayscale},
		{"gaussian", GaussianBlur},
		{"each_pixel", func(img image.Image) image.Image {
			out, err := EachPixelOfImage(img, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) { return b, g, r, a })
			if err != nil {
				panic(err)
			}
			return out
		}},
	} {
		b.Run(filter.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchmarkImage = filter.fn(img)
			}
		})
	}
}
