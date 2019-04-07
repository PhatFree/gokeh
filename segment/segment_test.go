package segment

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"testing"
)

func TestSegment(t *testing.T) {
	println("beging test")
	dest := "/Users/jeffdev/go/src/github.com/PhatFree/gokeh/images/digal_hiker_small.jpg"
	testImg := loadImage(dest)

	mask, err := segmentImage(testImg)

	if err != nil {
		fmt.Printf("There was an error: %s", err)
	}

	if mask.Opaque() {
		t.Errorf("the mask is empty")
	}

	println("ending test")

}

func loadImage(file string) (img image.Image) {
	reader, err := os.Open(file)
	if err != nil {
		println(file)
		log.Fatalf("error loading %s: %s", file, err)
	}
	img, _, err = image.Decode(reader)
	if err != nil {
		log.Fatalf("error loading %s: %s", file, err)
	}
	return img
}
