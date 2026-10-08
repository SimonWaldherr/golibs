package graphics

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type genericImage struct{ image.Image }

func offsetImage() *image.RGBA {
	parent := image.NewRGBA(image.Rect(-4, 2, 5, 8))
	img := parent.SubImage(image.Rect(-2, 4, 1, 6)).(*image.RGBA)
	for y := 4; y < 6; y++ {
		for x := -2; x < 1; x++ {
			v := uint8((y-4)*3 + x + 3)
			img.SetRGBA(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
		}
	}
	return img
}

func TestTransformsWithOffsetAndStride(t *testing.T) {
	img := offsetImage()
	before := append([]byte(nil), img.Pix...)
	for _, test := range []struct {
		name   string
		output image.Image
		want   []uint8
	}{
		{"flip", FlipHorizontal(img), []uint8{3, 2, 1, 6, 5, 4}},
		{"rotate", Rotate90(img), []uint8{4, 1, 5, 2, 6, 3}},
	} {
		t.Run(test.name, func(t *testing.T) {
			bounds := test.output.Bounds()
			var values []uint8
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, _, _, _ := test.output.At(x, y).RGBA()
					values = append(values, uint8(r>>8))
				}
			}
			if !reflect.DeepEqual(values, test.want) {
				t.Fatalf("pixels = %v, want %v", values, test.want)
			}
		})
	}
	if FlipHorizontal(img).Bounds() != img.Bounds() || Rotate90(img).Bounds() != image.Rect(0, 0, 2, 3) {
		t.Fatal("transform bounds incorrect")
	}
	if !reflect.DeepEqual(img.Pix, before) {
		t.Fatal("transform modified source")
	}
	out, err := NearestNeighbor(img, 6, 4)
	if err != nil {
		t.Fatal(err)
	}
	for y := 0; y < 4; y++ {
		for x := 0; x < 6; x++ {
			want := uint8(y/2*3 + x/2 + 1)
			if got := out.NRGBAAt(x, y).R; got != want {
				t.Fatalf("resized pixel (%d,%d) = %d, want %d", x, y, got, want)
			}
		}
	}
}

func TestGrayscaleMatchesStandardModel(t *testing.T) {
	parent := image.NewRGBA(image.Rect(-10, -10, 10, 10))
	for y := -10; y < 10; y++ {
		for x := -10; x < 10; x++ {
			parent.SetRGBA(x, y, color.RGBA{uint8((x + 10) * 7), uint8((y + 10) * 11), uint8((x + y + 20) * 5), 255})
		}
	}
	img := parent.SubImage(image.Rect(-7, -3, 6, 8)).(*image.RGBA)
	nrgba := image.NewNRGBA(img.Bounds())
	nrgba.SetNRGBA(-7, -3, color.NRGBA{100, 180, 200, 100})
	rgba64 := image.NewRGBA64(img.Bounds())
	rgba64.SetRGBA64(-7, -3, color.RGBA64{12345, 23456, 34567, 65535})
	for _, input := range []image.Image{img, genericImage{img}, nrgba, rgba64} {
		out := Grayscale(input)
		if out.Bounds() != input.Bounds() {
			t.Fatal("grayscale bounds changed")
		}
		for y := input.Bounds().Min.Y; y < input.Bounds().Max.Y; y++ {
			for x := input.Bounds().Min.X; x < input.Bounds().Max.X; x++ {
				want := color.GrayModel.Convert(input.At(x, y)).(color.Gray)
				if got := out.At(x, y).(color.Gray); got != want || GrayAt(input, x, y) != want.Y {
					t.Fatalf("gray pixel (%d,%d) = %v, want %v", x, y, got, want)
				}
			}
		}
	}
}

func TestPixelFilters(t *testing.T) {
	img := image.NewRGBA(image.Rect(3, 5, 4, 6))
	img.SetRGBA(3, 5, color.RGBA{10, 100, 250, 255})
	for _, test := range []struct {
		name   string
		output image.Image
		want   color.RGBA
	}{
		{"brightness up", AdjustBrightness(img, 20), color.RGBA{30, 120, 255, 255}},
		{"brightness down", AdjustBrightness(img, -20), color.RGBA{0, 80, 230, 255}},
		{"contrast zero", AdjustContrast(img, 0), color.RGBA{128, 128, 128, 255}},
		{"contrast identity", AdjustContrast(img, 1), color.RGBA{10, 100, 250, 255}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.output.Bounds() != img.Bounds() || test.output.At(3, 5).(color.RGBA) != test.want {
				t.Fatalf("output = %v, want %v", test.output.At(3, 5), test.want)
			}
		})
	}
	blue := image.NewRGBA(image.Rect(0, 0, 1, 1))
	blue.SetRGBA(0, 0, color.RGBA{0, 0, 255, 255})
	if Sepia(blue).At(0, 0).(color.RGBA) != (color.RGBA{48, 42, 33, 255}) {
		t.Fatal("sepia coefficients changed")
	}
	inverted, err := Invert(blue)
	if err != nil || inverted.At(0, 0).(color.RGBA) != (color.RGBA{255, 255, 0, 255}) {
		t.Fatal("inversion failed")
	}
}

func TestEachPixelAlphaAndInputPreservation(t *testing.T) {
	img := offsetImage()
	img.SetRGBA(-2, 4, color.RGBA{20, 30, 40, 50})
	before := append([]byte(nil), img.Pix...)
	calls := 0
	out, err := EachPixelOfImage(img, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) { calls++; return b, g, r, a })
	if err != nil || calls != 6 || out.At(-2, 4).(color.RGBA) != (color.RGBA{40, 30, 20, 50}) {
		t.Fatal("pixel callback changed alpha, channels, or count")
	}
	if !reflect.DeepEqual(img.Pix, before) {
		t.Fatal("pixel callback modified source")
	}
}

func TestResizeValidationAndAlpha(t *testing.T) {
	img := image.NewNRGBA(image.Rect(3, 4, 4, 5))
	img.SetNRGBA(3, 4, color.NRGBA{100, 180, 200, 100})
	out, err := NearestNeighbor(img, 2, 2)
	if err != nil || out.NRGBAAt(1, 1) != (color.NRGBA{100, 180, 200, 100}) {
		t.Fatal("resize changed straight-alpha channels")
	}
	for _, test := range []struct {
		image         image.Image
		width, height int
	}{
		{nil, 1, 1}, {img, 0, 1}, {img, 1, -1}, {image.NewRGBA(image.Rect(0, 0, 0, 0)), 1, 1}, {img, int(^uint(0) >> 1), 2},
	} {
		if result, err := NearestNeighbor(test.image, test.width, test.height); err == nil || result != nil {
			t.Fatal("invalid resize did not return an error")
		}
	}
}

func TestThresholdAndSaturatedEdges(t *testing.T) {
	img := image.NewGray(image.Rect(3, 5, 5, 6))
	img.SetGray(3, 5, color.Gray{100})
	img.SetGray(4, 5, color.Gray{101})
	out := Threshold(img, 100).(*image.Gray)
	if out.GrayAt(3, 5).Y != 0 || out.GrayAt(4, 5).Y != 255 {
		t.Fatal("threshold boundary incorrect")
	}
	step := image.NewGray(image.Rect(3, 5, 8, 10))
	for y := 5; y < 10; y++ {
		for x := 5; x < 8; x++ {
			step.SetGray(x, y, color.Gray{255})
		}
	}
	for _, output := range []image.Image{EdgeDetect(step), Edgedetect(step)} {
		if output.Bounds() != step.Bounds() || output.At(5, 7).(color.Gray).Y != 255 || output.At(3, 5).(color.Gray).Y != 0 {
			t.Fatal("Sobel magnitude overflow or bounds changed")
		}
	}
}

func TestConvolutionFilters(t *testing.T) {
	img := image.NewRGBA(image.Rect(-2, 4, 3, 9))
	img.SetRGBA(0, 6, color.RGBA{100, 100, 100, 100})
	gaussian := GaussianBlur(img).(*image.RGBA)
	if got := gaussian.RGBAAt(0, 6); got != (color.RGBA{15, 15, 15, 15}) {
		t.Fatalf("Gaussian impulse = %v, want 15 per channel", got)
	}
	if got := Blur(img).(*image.RGBA).RGBAAt(0, 6); got != (color.RGBA{11, 11, 11, 11}) {
		t.Fatalf("box impulse = %v, want 11 per channel", got)
	}
	if got := Sharpen(img).(*image.RGBA).RGBAAt(0, 6); got != (color.RGBA{255, 255, 255, 255}) {
		t.Fatalf("sharpen impulse = %v", got)
	}
	generic := GaussianBlur(genericImage{img}).(*image.RGBA)
	if !reflect.DeepEqual(gaussian.Pix, generic.Pix) {
		t.Fatal("generic convolution differs from RGBA path")
	}
	kernel := [][]float64{{-1, 0, 1}}
	if !reflect.DeepEqual(normalizeKernel(kernel), kernel) {
		t.Fatal("zero-sum kernel changed")
	}
}

func TestImageFiles(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 3, 2))
	for _, test := range []struct{ name, format string }{{"image.PNG", "png"}, {"image.jpg", "jpeg"}, {"image.jpeg", "jpeg"}, {"image.gif", "gif"}, {"i", "jpeg"}} {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), test.name)
			if err := SaveImage(img, path); err != nil {
				t.Fatal(err)
			}
			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			_, format, err := image.Decode(file)
			file.Close()
			if err != nil || format != test.format {
				t.Fatalf("format = %q, %v", format, err)
			}
			loaded, err := LoadImage(path)
			if err != nil || loaded.Bounds() != img.Bounds() {
				t.Fatalf("LoadImage: %v", err)
			}
			file, err = os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			resized, err := ResizeNearestNeighbor(file, 2, 1)
			if err != nil || resized.Bounds() != image.Rect(0, 0, 2, 1) {
				t.Fatalf("ResizeNearestNeighbor: %v", err)
			}
			file, err = os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			if _, err := EachPixel(file, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) { return r, g, b, a }); err != nil {
				t.Fatal(err)
			}
		})
	}
	if err := SaveImage(img, t.TempDir()); err == nil {
		t.Fatal("directory accepted as image file")
	}
	if _, err := LoadImage(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("missing image accepted")
	}
	invalid := filepath.Join(t.TempDir(), "invalid")
	if err := os.WriteFile(invalid, []byte("invalid image"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadImage(invalid); err == nil {
		t.Fatal("invalid image accepted")
	}
	file, err := os.Open(invalid)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if _, err := EachPixel(file, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) { return r, g, b, a }); err == nil {
		t.Fatal("invalid pixel image accepted")
	}
	if err := SaveImage(image.NewRGBA(image.Rectangle{}), filepath.Join(t.TempDir(), "empty.png")); err == nil {
		t.Fatal("empty PNG accepted")
	}
}
