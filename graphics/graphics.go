// manipulating images pixel by pixel
package graphics

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
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
			r, g, b, a := img.At(x, y).RGBA()
			newImg.Set(x, y, color.RGBA{uint8(255 - r>>8), uint8(255 - g>>8), uint8(255 - b>>8), uint8(a >> 8)})
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

// SaveImage saves an image to a file.
func SaveImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	switch {
	case filename[len(filename)-4:] == ".png":
		return png.Encode(f, img)
	case filename[len(filename)-4:] == ".jpg", filename[len(filename)-5:] == ".jpeg":
		return jpeg.Encode(f, img, nil)
	case filename[len(filename)-4:] == ".gif":
		return gif.Encode(f, img, nil)
	default:
		return jpeg.Encode(f, img, nil)
	}
}

// LoadImage loads an image from a file.
func LoadImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// AdjustBrightness adjusts the brightness of an image
func AdjustBrightness(img image.Image, adjustment float64) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			nr := clamp(float64(r>>8) + adjustment)
			ng := clamp(float64(g>>8) + adjustment)
			nb := clamp(float64(b>>8) + adjustment)
			newImg.Set(x, y, color.RGBA{uint8(nr), uint8(ng), uint8(nb), uint8(a >> 8)})
		}
	}
	return newImg
}

// AdjustContrast adjusts the contrast of an image
func AdjustContrast(img image.Image, factor float64) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			nr := clamp(((float64(r>>8) - 128) * factor) + 128)
			ng := clamp(((float64(g>>8) - 128) * factor) + 128)
			nb := clamp(((float64(b>>8) - 128) * factor) + 128)
			newImg.Set(x, y, color.RGBA{uint8(nr), uint8(ng), uint8(nb), uint8(a >> 8)})
		}
	}
	return newImg
}

// clamp ensures the value is between 0 and 255
func clamp(value float64) float64 {
	return math.Min(math.Max(value, 0), 255)
}

// Blur applies a simple blur filter to an image
func Blur(img image.Image) image.Image {
	kernel := [][]float64{
		{1 / 9.0, 1 / 9.0, 1 / 9.0},
		{1 / 9.0, 1 / 9.0, 1 / 9.0},
		{1 / 9.0, 1 / 9.0, 1 / 9.0},
	}
	return applyKernel(img, kernel)
}

// Sepia applies a sepia tone to an image
func Sepia(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			tr := clamp(0.393*float64(r>>8) + 0.769*float64(g>>8) + 0.189*float64(b>>8))
			tg := clamp(0.349*float64(r>>8) + 0.686*float64(g>>8) + 0.168*float64(b>>8))
			tb := clamp(0.272*float64(r>>8) + 0.534*float64(g>>8) + 0.131*float64(b>>8))
			newImg.Set(x, y, color.RGBA{uint8(tr), uint8(tg), uint8(tb), uint8(a >> 8)})
		}
	}
	return newImg
}

// Rotate90 rotates an image 90 degrees clockwise
func Rotate90(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(bounds.Max.Y-y-1, x, img.At(x, y))
		}
	}
	return newImg
}

// FlipHorizontal flips an image horizontally
func FlipHorizontal(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(bounds.Max.X-x-1, y, img.At(x, y))
		}
	}
	return newImg
}

// EdgeDetect detects edges in a grayscale image.
func EdgeDetect(img image.Image) image.Image {
	grayImg := Grayscale(img)
	gray, ok := grayImg.(*image.Gray)
	if !ok {
		log.Fatalf("Expected *image.Gray but got %T", grayImg)
	}
	b := gray.Bounds()
	edgeImg := image.NewGray(b)

	for y := 1; y < b.Max.Y-1; y++ {
		for x := 1; x < b.Max.X-1; x++ {
			gx, gy := applySobel(gray, x, y)
			val := uint8(math.Sqrt(float64(gx*gx + gy*gy)))
			edgeImg.SetGray(x, y, color.Gray{Y: val})
		}
	}
	return edgeImg
}

func applySobel(img *image.Gray, x, y int) (int, int) {
	var gx, gy int
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			pixel := img.GrayAt(x+j, y+i).Y
			gx += int(edfX[i+1][j+1]) * int(pixel)
			gy += int(edfY[i+1][j+1]) * int(pixel)
		}
	}
	return gx, gy
}

// Convolution kernel for Gaussian blur
var gaussianKernel = [][]float64{
	{1, 4, 7, 4, 1},
	{4, 16, 26, 16, 4},
	{7, 26, 41, 26, 7},
	{4, 16, 26, 16, 4},
	{1, 4, 7, 4, 1},
}

// Sharpen kernel
var sharpenKernel = [][]float64{
	{0, -1, 0},
	{-1, 5, -1},
	{0, -1, 0},
}

// normalizeKernel normalizes a convolution kernel
func normalizeKernel(kernel [][]float64) [][]float64 {
	var sum float64
	for _, row := range kernel {
		for _, val := range row {
			sum += val
		}
	}

	if sum == 0 {
		return kernel
	}

	normalizedKernel := make([][]float64, len(kernel))
	for i := range kernel {
		normalizedKernel[i] = make([]float64, len(kernel[i]))
		for j := range kernel[i] {
			normalizedKernel[i][j] = kernel[i][j] / sum
		}
	}
	return normalizedKernel
}

// applyKernel applies a convolution kernel to an image
func applyKernel(img image.Image, kernel [][]float64) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	normalizedKernel := normalizeKernel(kernel)
	kernelSize := len(normalizedKernel)
	offset := kernelSize / 2

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rSum, gSum, bSum, aSum float64

			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					px := x + kx - offset
					py := y + ky - offset

					if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
						r, g, b, a := img.At(px, py).RGBA()
						weight := normalizedKernel[ky][kx]
						rSum += float64(r>>8) * weight
						gSum += float64(g>>8) * weight
						bSum += float64(b>>8) * weight
						aSum += float64(a>>8) * weight
					}
				}
			}

			newImg.Set(x, y, color.RGBA{
				uint8(math.Min(math.Max(rSum, 0), 255)),
				uint8(math.Min(math.Max(gSum, 0), 255)),
				uint8(math.Min(math.Max(bSum, 0), 255)),
				uint8(math.Min(math.Max(aSum, 0), 255)),
			})
		}
	}
	return newImg
}

// GaussianBlur applies a Gaussian blur to an image
func GaussianBlur(img image.Image) image.Image {
	return applyKernel(img, gaussianKernel)
}

// Sharpen applies a sharpening filter to an image
func Sharpen(img image.Image) image.Image {
	return applyKernel(img, sharpenKernel)
}
