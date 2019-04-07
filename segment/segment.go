package segment

import (
	"fmt"
	"image"
	_ "image/jpeg"

	"gocv.io/x/gocv"
)

func segmentImage(in image.Image) (image.Gray, error) {
	mat, _ := gocv.ImageToMatRGB(in)

	/*
		Do nerual net stuff here
	*/

	temp, _ := mat.ToImage()

	mask, ok := temp.(*image.Gray)

	if !ok {
		return image.Gray{}, fmt.Errorf("not a Gray mat")
	}
	return *mask, nil
}
