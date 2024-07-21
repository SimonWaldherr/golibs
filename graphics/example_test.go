package graphics_test

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"simonwaldherr.de/go/golibs/graphics"
	//"../graphics"
)

func ExampleRNN() {
	img := createTestImage()

	img_rnn, _ := graphics.NearestNeighbor(img, 80, 24)
	saveImage("rnn_example.png", img_rnn)
	fmt.Println("Inverted image saved to rnn_example.png")
	// Output: Inverted image saved to rnn_example.png
}

func ExampleInvert() {
	img := createTestImage()

	invertedImg, _ := graphics.Invert(img)
	saveImage("invert_example.png", invertedImg)
	fmt.Println("Inverted image saved to invert_example.png")
	// Output: Inverted image saved to invert_example.png
}

func ExampleGrayscale() {
	img := createTestImage()

	grayImg := graphics.Grayscale(img)
	saveImage("grayscale_example.png", grayImg)
	fmt.Println("Grayscale image saved to grayscale_example.png")
	// Output: Grayscale image saved to grayscale_example.png
}

func ExampleGaussianBlur() {
	img := createTestImage()

	blurredImg := graphics.GaussianBlur(img)
	saveImage("gaussian_blur_example.png", blurredImg)
	fmt.Println("Gaussian blurred image saved to gaussian_blur_example.png")
	// Output: Gaussian blurred image saved to gaussian_blur_example.png
}

func ExampleSharpen() {
	img := createTestImage()

	sharpenedImg := graphics.Sharpen(img)
	saveImage("sharpen_example.png", sharpenedImg)
	fmt.Println("Sharpened image saved to sharpen_example.png")
	// Output: Sharpened image saved to sharpen_example.png
}

func ExampleEdgeDetect() {
	img := createTestImage()

	edgeImg := graphics.EdgeDetect(img)
	saveImage("edge_detect_example.png", edgeImg)
	fmt.Println("Edge detected image saved to edge_detect_example.png")
	// Output: Edge detected image saved to edge_detect_example.png
}

func ExampleThreshold() {
	img := createTestImage()

	threshImg := graphics.Threshold(img, 128)
	saveImage("threshold_example.png", threshImg)
	fmt.Println("Thresholded image saved to threshold_example.png")
	// Output: Thresholded image saved to threshold_example.png
}

func ExampleAdjustBrightness() {
	img := createTestImage()

	brightImg := graphics.AdjustBrightness(img, 30)
	saveImage("brightness_example.png", brightImg)
	fmt.Println("Brightness adjusted image saved to brightness_example.png")
	// Output: Brightness adjusted image saved to brightness_example.png
}

func ExampleAdjustContrast() {
	img := createTestImage()

	contrastImg := graphics.AdjustContrast(img, 1.5)
	saveImage("contrast_example.png", contrastImg)
	fmt.Println("Contrast adjusted image saved to contrast_example.png")
	// Output: Contrast adjusted image saved to contrast_example.png
}

func ExampleSepia() {
	img := createTestImage()

	sepiaImg := graphics.Sepia(img)
	saveImage("sepia_example.png", sepiaImg)
	fmt.Println("Sepia image saved to sepia_example.png")
	// Output: Sepia image saved to sepia_example.png
}

func ExampleRotate90() {
	img := createTestImage()

	rotatedImg := graphics.Rotate90(img)
	saveImage("rotate90_example.png", rotatedImg)
	fmt.Println("Rotated image saved to rotate90_example.png")
	// Output: Rotated image saved to rotate90_example.png
}

func ExampleFlipHorizontal() {
	img := createTestImage()

	flippedImg := graphics.FlipHorizontal(img)
	saveImage("flip_horizontal_example.png", flippedImg)
	fmt.Println("Flipped image saved to flip_horizontal_example.png")
	// Output: Flipped image saved to flip_horizontal_example.png
}

// createTestImage creates a simple test image
func createTestImage() image.Image {
	img, _ := graphics.LoadImage("img.jpg")

	return img
}

// saveImage saves an image to a file
func saveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		fmt.Println("Failed to encode image:", err)
	}
}
