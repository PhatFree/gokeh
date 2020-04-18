package filter

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"testing"
	"time"
)

func init() {
	LoadTestImages()
}

var testImage *image.RGBA
var testMask *image.Gray16

func LoadTestImages() {
	// Read image from file that already exists
	testImageFile, err := os.Open("test_data/tree.png")
	maskImageFile, err := os.Open("test_data/mask.png")
	if err != nil {
		panic("Can't start test, can't read test data")
	}
	defer testImageFile.Close()
	defer maskImageFile.Close()

	// Alternatively, since we know it is a png already
	// we can call png.Decode() directly
	testImageRaw, err := png.Decode(testImageFile)
	testImageBounds := testImageRaw.Bounds()
	testImage = image.NewRGBA(testImageBounds)
	testMaskRaw, err := png.Decode(maskImageFile)
	testMaskBounds := testMaskRaw.Bounds()
	testMask = image.NewGray16(testMaskBounds)
	if err != nil {
		panic("Can't start test, can't decode test data")
	}

	for y := 0; y < testImageBounds.Max.Y; y++ {
		for x := 0; x < testImageBounds.Max.X; x++ {
			testImage.Set(x, y, color.RGBAModel.Convert(testImageRaw.At(x, y)))
		}
	}
	for y := 0; y < testMaskBounds.Max.Y; y++ {
		for x := 0; x < testMaskBounds.Max.X; x++ {
			testMask.Set(x, y, color.Gray16Model.Convert(testMaskRaw.At(x, y)))
		}
	}
}

func TestApplyBlur(t *testing.T) {
	blurred := ApplyBlur(testImage, testMask)

	// outputFile is a File type which satisfies Writer interface
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	outputFile, err := os.Create("output." + timestamp + ".png")
	if err != nil {
		t.Errorf("Failed to output image")
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, blurred)

	// Don't forget to close files
	outputFile.Close()
}
