package graphics

import (
	_ "code.google.com/p/vp8-go/webp"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func EachPixel(file *os.File, f func(uint8, uint8, uint8, uint8) (uint8, uint8, uint8, uint8)) image.Image {

	src, _, err := image.Decode(file)

	bsrc := src.Bounds()
	img := image.NewRGBA(bsrc)
	draw.Draw(img, bsrc, src, bsrc.Min, draw.Src)

	if err != nil {
		return nil
	}

	b := img.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			oldPixel := img.At(x, y)
			r1, g1, b1, a1 := oldPixel.RGBA()
			r2, g2, b2, a2 := f(uint8(r1), uint8(g1), uint8(b1), uint8(a1))
			pixel := color.RGBA{uint8(r2), uint8(g2), uint8(b2), uint8(a2)}
			img.Set(x, y, pixel)
		}
	}
	return img
}
