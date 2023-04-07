package graphics_test

import (
	"simonwaldherr.de/go/golibs/graphics"
	//"../graphics"
)

func Example() {
	img := graphics.LoadImage("img.jpg")
	img_gray := graphics.Grayscale(img)

	img_rnn, _ := graphics.NearestNeighbor(img, 80, 24)

	img_each, _ := graphics.EachPixelOfImage(img, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		return g, b, r, a
	})

	img_edge := graphics.Edgedetect(img)

	img_invert, _ := graphics.Invert(img)

	img_threshold := graphics.Threshold(img, 128)

	graphics.SaveImage(img_gray, "img_gray.png")
	graphics.SaveImage(img_rnn, "img_rnn.png")
	graphics.SaveImage(img_each, "img_each.png")
	graphics.SaveImage(img_edge, "img_edge.png")
	graphics.SaveImage(img_invert, "img_invert.png")
	graphics.SaveImage(img_threshold, "img_threshold.png")

	// Output:
}
