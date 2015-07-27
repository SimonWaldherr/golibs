package graphics

import (
	"image/png"
	"os"
	"testing"
)

func Test_EachPixel(t *testing.T) {
	var xr, xg, xb int = 0, 0, 0
	file, err := os.Open("./img.png")
	file1, err1 := os.Open("./graphics.go")
	defer func() {
		file.Close()
		file1.Close()
	}()

	if err != nil {
		t.Fatalf("EachPixel Test failed: %v", err)
	}

	_, err = EachPixel(file1, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		return g, b, r, a
	})
	if err1 != nil || err == nil {
		t.Fatalf("EachPixel Test failed")
	}

	img, _ := EachPixel(file, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		if r > g && r > b {
			xr++
		}
		if g > r && g > b {
			xg++
		}
		if b > g && b > r {
			xb++
		}
		return g, b, r, a
	})
	fd, err := os.Create("./inv.png")
	if err != nil {
		t.Fatalf("EachPixel Test failed: %v", err)
	}

	err = png.Encode(fd, img)
	if err != nil {
		t.Fatalf("EachPixel Test failed: %v", err)
	}

	err = fd.Close()
	if err != nil {
		t.Fatalf("EachPixel Test failed: %v", err)
	}
	if xr != 11409 || xg != 25162 || xb != 94726 {
		t.Fatalf("EachPixel Test failed")
	}
}

func Test_ResizeNearestNeighbor(t *testing.T) {
	file, err := os.Open("./img.png")
	file1, err1 := os.Open("./img.err")
	defer func() {
		file.Close()
		file1.Close()
	}()

	if err1 == nil {
		t.Fatalf("ResizeNearestNeighbor Test failed: %v", err)
	}

	img, _ := ResizeNearestNeighbor(file, 200, 150)
	fd, err := os.Create("./rnn.png")
	if err != nil {
		t.Fatalf("ResizeNearestNeighbor Test failed: %v", err)
	}

	err = png.Encode(fd, img)
	if err != nil {
		t.Fatalf("ResizeNearestNeighbor Test failed: %v", err)
	}

	err = fd.Close()
	if err != nil {
		t.Fatalf("ResizeNearestNeighbor Test failed: %v", err)
	}
}

func Test_ImageDecodeFail(t *testing.T) {
	file, err := os.Open("../test.txt")
	defer func() {
		file.Close()
	}()

	_, err = ResizeNearestNeighbor(file, 200, 150)

	if err == nil {
		t.Fatalf("ImageDecodeFail Test failed: %v", err)
	}
}
