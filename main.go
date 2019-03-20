package main

import (
	"flag"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {

	var fileName string
	flag.StringVar(&fileName, "file", "images/digital_car_small.jpg", "the path to the image")

	flag.Parse()

	var err error

	if fileName != "" {
		fileName, err = filepath.Abs(fileName)
		if err != nil {
			panic(err)
		}
	}
	img := loadImage(fileName)
	if img != nil {
		println("Sucess!")
	}
}

func loadImage(file string) image.Image {
	reader, err := os.Open(file)
	if err != nil {
		log.Fatalf("error loading %s: %s", file, err)
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatalf("error loading %s: %s", file, err)
	}
	return img
}
