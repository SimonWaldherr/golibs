// manipulating images pixel by pixel
package graphics

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func EachPixel(file *os.File, f func(uint8, uint8, uint8, uint8) (uint8, uint8, uint8, uint8)) (image.Image, error) {
	src, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bsrc := src.Bounds()
	img := image.NewRGBA(bsrc)
	draw.Draw(img, bsrc, src, bsrc.Min, draw.Src)

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
	return img, nil
}

func ResizeNearestNeighbor(file *os.File, newWidth, newHeight int) (*image.NRGBA, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	newimg := image.NewNRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{newWidth, newHeight}})

	xn := (w<<16)/newWidth + 1
	yn := (h<<16)/newHeight + 1

	for yo := 0; yo < newHeight; yo++ {
		y := (yo * yn) >> 16
		for xo := 0; xo < newWidth; xo++ {
			x := (xo * xn) >> 16
			newimg.Set(xo, yo, img.At(x, y))
		}
	}
	return newimg, nil
}
