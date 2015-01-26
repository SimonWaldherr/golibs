package graphics

import (
	"image/png"
	"os"
	"testing"
)

func Test_EachPixel(t *testing.T) {
	var xr, xg, xb int = 0, 0, 0
	file, err := os.Open("./img.png")
	defer file.Close()

	if err != nil {
		t.Fatalf("EachPixel Test failed: %v", err)
	}
	img := EachPixel(file, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
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
