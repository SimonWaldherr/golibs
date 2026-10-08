package graphics

import (
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"testing"
)

// makeTestImage returns a small deterministic RGBA image for fast tests.
func makeTestImage(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), uint8((x + y) % 256), 255})
		}
	}
	return img
}

func Test_EachPixel(t *testing.T) {
	t.Parallel()
	img := makeTestImage(80, 60)

	out, err := EachPixelOfImage(img, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		return g, b, r, a
	})
	if err != nil {
		t.Fatalf("EachPixelOfImage failed: %v", err)
	}
	if out == nil {
		t.Fatalf("EachPixelOfImage returned nil image")
	}
	if out.Bounds() != img.Bounds() {
		t.Fatalf("output bounds mismatch: got %v want %v", out.Bounds(), img.Bounds())
	}
}

func Test_ResizeNearestNeighbor(t *testing.T) {
	t.Parallel()
	img := makeTestImage(40, 30)
	newimg, err := NearestNeighbor(img, 20, 10)
	if err != nil {
		t.Fatalf("NearestNeighbor failed: %v", err)
	}
	if newimg.Bounds().Dx() != 20 || newimg.Bounds().Dy() != 10 {
		t.Fatalf("unexpected resized dimensions: %v", newimg.Bounds())
	}
}

func Test_Grayscale(t *testing.T) {
	t.Parallel()
	img := makeTestImage(32, 16)
	g := Grayscale(img)
	if g == nil {
		t.Fatalf("Grayscale returned nil")
	}
	if g.Bounds() != img.Bounds() {
		t.Fatalf("Grayscale bounds mismatch: got %v want %v", g.Bounds(), img.Bounds())
	}
}

func Test_Edgedetect(t *testing.T) {
	t.Parallel()
	img := makeTestImage(48, 48)
	ed := EdgeDetect(img)
	if ed == nil {
		t.Fatalf("EdgeDetect returned nil")
	}
	if ed.Bounds() != img.Bounds() {
		t.Fatalf("EdgeDetect bounds mismatch: got %v want %v", ed.Bounds(), img.Bounds())
	}
}

func Test_ImageDecodeFail(t *testing.T) {
	// create a temporary file with invalid image data
	tmp, err := ioutil.TempFile(".", "img_invalid_*.txt")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}
	name := tmp.Name()
	defer func() { _ = os.Remove(name) }()
	if _, err := tmp.WriteString("not an image"); err != nil {
		t.Fatalf("unable to write temp file: %v", err)
	}
	if err := tmp.Close(); err != nil {
		t.Fatalf("unable to close temp file: %v", err)
	}

	f, err := os.Open(name)
	if err != nil {
		t.Fatalf("unable to open temp file: %v", err)
	}
	defer func() { _ = f.Close() }()

	_, err = ResizeNearestNeighbor(f, 200, 150)
	if err == nil {
		t.Fatalf("expected decode error for invalid image data")
	}
}
