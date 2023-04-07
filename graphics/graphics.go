// manipulating images pixel by pixel
package graphics

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
)

var edfX = [3][3]int8{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
var edfY = [3][3]int8{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}

// Invert inverts an image
func Invert(img image.Image) (image.Image, error) {
	b := img.Bounds()
	newImg := image.NewRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pixel := img.At(x, y)
			r, g, b, a := pixel.RGBA()
			newImg.Set(x, y, color.RGBA{uint8(255 - r), uint8(255 - g), uint8(255 - b), uint8(a)})
		}
	}
	return newImg, nil
}

// EachPixel applies a function to each pixel of an image
func EachPixel(file *os.File, f func(uint8, uint8, uint8, uint8) (uint8, uint8, uint8, uint8)) (image.Image, error) {
	src, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	return EachPixelOfImage(src, f)
}

// EachPixelOfImage applies a function to each pixel of an image
func EachPixelOfImage(src image.Image, f func(uint8, uint8, uint8, uint8) (uint8, uint8, uint8, uint8)) (image.Image, error) {
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

// ResizeNearestNeighbor resizes an image using the nearest neighbor algorithm
func ResizeNearestNeighbor(file *os.File, newWidth, newHeight int) (*image.NRGBA, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return NearestNeighbor(img, newWidth, newHeight)
}

// NearestNeighbor resizes an image using the nearest neighbor algorithm
func NearestNeighbor(img image.Image, newWidth, newHeight int) (*image.NRGBA, error) {
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

// Grayscale converts an image to grayscale
func Grayscale(img image.Image) image.Image {
	min := img.Bounds().Min
	max := img.Bounds().Max

	grayImg := image.NewGray(image.Rect(max.X, max.Y, min.X, min.Y))
	for x := 0; x < max.X; x++ {
		for y := 0; y < max.Y; y++ {
			grayPixel := color.GrayModel.Convert(img.At(x, y))
			grayImg.Set(x, y, grayPixel)
		}
	}
	return grayImg
}

// Edgedetect detects edges in an image
func Edgedetect(img image.Image) image.Image {
	img = Grayscale(img)
	max := img.Bounds().Max
	min := img.Bounds().Min

	edImg := image.NewGray(image.Rect(max.X, max.Y, min.X, min.Y))

	width := max.X
	height := max.Y
	var pixel color.Color
	for x := 1; x < width-1; x++ {
		for y := 1; y < height-1; y++ {
			pX, pY := detectEdgeAroundPixel(img, x, y)
			val := uint32(math.Ceil(math.Sqrt(pX*pX + pY*pY)))
			pixel = color.Gray{Y: uint8(val)}
			edImg.SetGray(x, y, pixel.(color.Gray))
		}
	}
	return edImg
}

// detectEdgeAroundPixel detects the edge around a pixel
func detectEdgeAroundPixel(img image.Image, x int, y int) (float64, float64) {
	var pX, pY int
	curX := x - 1
	curY := y - 1

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r, _, _, _ := img.At(curX, curY).RGBA()
			pX += int(edfX[i][j]) * int(uint8(r))
			pY += int(edfY[i][j]) * int(uint8(r))
			curX++
		}
		curX = x - 1
		curY++
	}

	return math.Abs(float64(pX)), math.Abs(float64(pY))
}

// Threshold thresholds an image
func Threshold(img image.Image, threshold uint8) image.Image {
	img = Grayscale(img)
	max := img.Bounds().Max
	min := img.Bounds().Min

	thImg := image.NewGray(image.Rect(max.X, max.Y, min.X, min.Y))

	width := max.X
	height := max.Y
	var pixel color.Color
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			p := img.At(x, y)
			if r, _, _, _ := p.RGBA(); uint8(r) > threshold {
				pixel = color.Gray{Y: 255}
			} else {
				pixel = color.Gray{Y: 0}
			}
			thImg.SetGray(x, y, pixel.(color.Gray))
		}
	}
	return thImg
}

// GrayAt returns the grayscale value of a pixel
func GrayAt(img image.Image, x, y int) uint8 {
	p := img.At(x, y)
	r, _, _, _ := p.RGBA()
	return uint8(r)
}

// SaveImage saves an image to a file
func SaveImage(img image.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)
}

// LoadImage loads an image from a file
func LoadImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
