package segment

import (
	"fmt"
	"image"
	_ "image/jpeg"

	"gocv.io/x/gocv"
)

func segmentImage(in image.Image) (image.Gray, error) {

	//load image into mat for openCV
	//currently only wrriten to take RGB, grayscale would be easy to add
	matIn, err := gocv.ImageToMatRGB(in)
	if err != nil {
		return image.Gray{}, err
	}
	matDes := gocv.NewMat()

	net :=
		gocv.ReadNetFromTensorflow("$GOPATH/src/github.com/PhatFree/gokeh/segment/tfModels/monodepth/model_city2eigen.pb")

	//inputs needed
	//

	/*
		Do nerual net stuff here
	*/

	depSz := matDes.Size()

	matDpth := gocv.NewMat()
	matTemp := gocv.NewMatWithSizeFromScalar(gocv.NewScalar(0.3128, 0.3128, 0.3128, 0.3128), depSz[0], depSz[1], matDes.Type())
	matDes.AddFloat(0.00001)
	gocv.Divide(matTemp, matDes, &matDpth)

	temp, _ := matDpth.ToImage()

	mask, ok := temp.(*image.Gray)

	if !ok {
		return image.Gray{}, fmt.Errorf("not a Gray mat")
	}
	return *mask, nil
}
