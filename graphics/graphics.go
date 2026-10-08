// Package graphics provides image transforms, filters, and file encoders.
package graphics

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
)

var edfX = [3][3]int8{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
var edfY = [3][3]int8{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}

// Invert inverts the premultiplied RGB channels and preserves alpha and bounds.
func Invert(img image.Image) (image.Image, error) {
	return EachPixelOfImage(img, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		return 255 - r, 255 - g, 255 - b, a
	})
}

// EachPixel applies a function to each pixel of an image
func EachPixel(file *os.File, f func(uint8, uint8, uint8, uint8) (uint8, uint8, uint8, uint8)) (image.Image, error) {
	src, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	return EachPixelOfImage(src, f)
}

// EachPixelOfImage transforms 8-bit premultiplied RGBA channels. It preserves
// bounds and leaves src unchanged. The callback runs once per pixel in row order.
func EachPixelOfImage(src image.Image, f func(uint8, uint8, uint8, uint8) (uint8, uint8, uint8, uint8)) (image.Image, error) {
	bounds := src.Bounds()
	out := image.NewRGBA(bounds)
	draw.Draw(out, bounds, src, bounds.Min, draw.Src)
	for y := 0; y < bounds.Dy(); y++ {
		row := out.Pix[y*out.Stride:][:bounds.Dx()*4]
		for x := 0; x < len(row); x += 4 {
			row[x], row[x+1], row[x+2], row[x+3] = f(row[x], row[x+1], row[x+2], row[x+3])
		}
	}
	return out, nil
}

// ResizeNearestNeighbor resizes an image using the nearest neighbor algorithm
func ResizeNearestNeighbor(file *os.File, newWidth, newHeight int) (*image.NRGBA, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return NearestNeighbor(img, newWidth, newHeight)
}

// NearestNeighbor resizes to positive dimensions, using floor-scaled source
// coordinates. The result starts at (0, 0); source offsets and alpha are preserved.
// Empty sources and nonpositive or overflowing output dimensions return an error.
func NearestNeighbor(img image.Image, newWidth, newHeight int) (*image.NRGBA, error) {
	if img == nil || newWidth <= 0 || newHeight <= 0 {
		return nil, fmt.Errorf("graphics: source and positive dimensions are required")
	}
	bounds := img.Bounds()
	if bounds.Empty() {
		return nil, fmt.Errorf("graphics: cannot resize an empty image")
	}
	maxInt := int(^uint(0) >> 1)
	if newWidth > maxInt/4 || newHeight > maxInt/(4*newWidth) {
		return nil, fmt.Errorf("graphics: output dimensions overflow")
	}
	out := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight))
	for y := 0; y < newHeight; y++ {
		sy := bounds.Min.Y + y*bounds.Dy()/newHeight
		for x := 0; x < newWidth; x++ {
			sx := bounds.Min.X + x*bounds.Dx()/newWidth
			pixel := color.NRGBAModel.Convert(img.At(sx, sy)).(color.NRGBA)
			out.SetNRGBA(x, y, pixel)
		}
	}
	return out, nil
}

// Grayscale converts using color.GrayModel, preserving source bounds.
// It uses direct pixel access for RGBA sources, including subimages.
func Grayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	out := image.NewGray(bounds)
	if src, ok := img.(*image.RGBA); ok {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			source := src.Pix[src.PixOffset(bounds.Min.X, y):][:bounds.Dx()*4]
			target := out.Pix[(y-bounds.Min.Y)*out.Stride:][:bounds.Dx()]
			for x := range target {
				i := 4 * x
				r, g, b := uint32(source[i])*0x101, uint32(source[i+1])*0x101, uint32(source[i+2])*0x101
				target[x] = uint8((19595*r + 38470*g + 7471*b + (1 << 15)) >> 24)
			}
		}
		return out
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			out.SetGray(x, y, color.GrayModel.Convert(img.At(x, y)).(color.Gray))
		}
	}
	return out
}

// Edgedetect is the legacy Sobel filter, rounding edge magnitudes upward.
// EdgeDetect rounds downward. Both clamp magnitudes to 255 and leave borders black.
func Edgedetect(img image.Image) image.Image { return edgeDetect(img, true) }

// Threshold returns a black/white image, preserving bounds. Grayscale values
// strictly greater than threshold become white; other values become black.
func Threshold(img image.Image, threshold uint8) image.Image {
	gray := Grayscale(img).(*image.Gray)
	out := image.NewGray(gray.Bounds())
	for i, value := range gray.Pix {
		if value > threshold {
			out.Pix[i] = 255
		}
	}
	return out
}

// GrayAt returns the grayscale value of a pixel using color.GrayModel.
func GrayAt(img image.Image, x, y int) uint8 {
	return color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y
}

// SaveImage writes PNG, JPEG, or GIF according to the filename extension
// (case-insensitive). Other extensions use JPEG. It reports encoding and close errors.
func SaveImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	var encodeErr error
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".png":
		encodeErr = png.Encode(f, img)
	case ".gif":
		encodeErr = gif.Encode(f, img, nil)
	default:
		encodeErr = jpeg.Encode(f, img, nil)
	}
	closeErr := f.Close()
	if encodeErr != nil {
		return encodeErr
	}
	return closeErr
}

// LoadImage loads an image from a file.
func LoadImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// AdjustBrightness adds adjustment to each premultiplied RGB channel,
// clamping to [0,255] and preserving alpha and bounds.
func AdjustBrightness(img image.Image, adjustment float64) image.Image {
	out, _ := EachPixelOfImage(img, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		return uint8(clamp(float64(r) + adjustment)), uint8(clamp(float64(g) + adjustment)), uint8(clamp(float64(b) + adjustment)), a
	})
	return out
}

// AdjustContrast scales premultiplied RGB values around 128, clamping to
// [0,255] and preserving alpha and bounds. A factor of 1 preserves the image.
func AdjustContrast(img image.Image, factor float64) image.Image {
	out, _ := EachPixelOfImage(img, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		return uint8(clamp((float64(r)-128)*factor + 128)), uint8(clamp((float64(g)-128)*factor + 128)), uint8(clamp((float64(b)-128)*factor + 128)), a
	})
	return out
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

// Sepia applies a sepia tone to premultiplied RGB, preserving alpha and bounds.
func Sepia(img image.Image) image.Image {
	out, _ := EachPixelOfImage(img, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		return uint8(clamp(0.393*float64(r) + 0.769*float64(g) + 0.189*float64(b))),
			uint8(clamp(0.349*float64(r) + 0.686*float64(g) + 0.168*float64(b))),
			uint8(clamp(0.272*float64(r) + 0.534*float64(g) + 0.131*float64(b))), a
	})
	return out
}

// Rotate90 rotates clockwise into an image starting at (0,0).
func Rotate90(img image.Image) image.Image {
	src := rgbaSource(img)
	bounds := src.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			from := src.PixOffset(x, y)
			to := out.PixOffset(bounds.Max.Y-y-1, x-bounds.Min.X)
			copy(out.Pix[to:to+4], src.Pix[from:from+4])
		}
	}
	return out
}

// FlipHorizontal mirrors an image horizontally, preserving its bounds.
func FlipHorizontal(img image.Image) image.Image {
	src := rgbaSource(img)
	bounds := src.Bounds()
	out := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			from := src.PixOffset(x, y)
			to := out.PixOffset(bounds.Min.X+bounds.Max.X-x-1, y)
			copy(out.Pix[to:to+4], src.Pix[from:from+4])
		}
	}
	return out
}

// EdgeDetect applies a Sobel filter, clamping edge magnitude to [0,255].
// It preserves bounds and leaves the outermost pixel border black.
func EdgeDetect(img image.Image) image.Image { return edgeDetect(img, false) }

func edgeDetect(img image.Image, roundUp bool) image.Image {
	gray := Grayscale(img).(*image.Gray)
	bounds := gray.Bounds()
	out := image.NewGray(bounds)
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			gx, gy := applySobel(gray, x, y)
			magnitude := math.Sqrt(float64(gx*gx + gy*gy))
			if roundUp {
				magnitude = math.Ceil(magnitude)
			}
			out.SetGray(x, y, color.Gray{Y: uint8(clamp(magnitude))})
		}
	}
	return out
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

// rgbaSource borrows RGBA pixels or converts once. Callers must not modify them.
func rgbaSource(img image.Image) *image.RGBA {
	if src, ok := img.(*image.RGBA); ok {
		return src
	}
	bounds := img.Bounds()
	out := image.NewRGBA(bounds)
	draw.Draw(out, bounds, img, bounds.Min, draw.Src)
	return out
}

// applyKernel uses zero padding and the normalized full kernel at image edges.
func applyKernel(img image.Image, kernel [][]float64) image.Image {
	src := rgbaSource(img)
	bounds := src.Bounds()
	out := image.NewRGBA(bounds)
	weights := normalizeKernel(kernel)
	offset := len(weights) / 2
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rSum, gSum, bSum, aSum float64
			for ky, row := range weights {
				py := y + ky - offset
				if py < bounds.Min.Y || py >= bounds.Max.Y {
					continue
				}
				for kx, weight := range row {
					px := x + kx - offset
					if px < bounds.Min.X || px >= bounds.Max.X {
						continue
					}
					i := src.PixOffset(px, py)
					rSum += float64(src.Pix[i]) * weight
					gSum += float64(src.Pix[i+1]) * weight
					bSum += float64(src.Pix[i+2]) * weight
					aSum += float64(src.Pix[i+3]) * weight
				}
			}
			i := out.PixOffset(x, y)
			out.Pix[i], out.Pix[i+1], out.Pix[i+2], out.Pix[i+3] = uint8(clamp(rSum)), uint8(clamp(gSum)), uint8(clamp(bSum)), uint8(clamp(aSum))
		}
	}
	return out
}

// GaussianBlur applies a Gaussian blur to an image
func GaussianBlur(img image.Image) image.Image {
	return applyKernel(img, gaussianKernel)
}

// Sharpen applies a sharpening filter to an image
func Sharpen(img image.Image) image.Image {
	return applyKernel(img, sharpenKernel)
}
