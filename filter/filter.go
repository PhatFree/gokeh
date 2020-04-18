package filter

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/fcolor"
	"github.com/anthonynsimon/bild/parallel"
)

func ApplyBlur(in *image.RGBA, mask *image.Gray16) (out *image.RGBA) {
	// Copy the input image into the output image
	out = clone.AsRGBA(in)

	// Look at each pixel in the mask
	bounds := mask.Bounds()
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			maskPixel := mask.At(x, y).(color.Gray16)
			// Make sure that mask pixel isn't black
			// TODO: maybe also check if it's just really dim
			if maskPixel.Y != 0 {
				// accumulate a shifted copy of the color image weighted by the mask
				multLayer := multiply(in, maskPixel)
				adjMult := multiply(multLayer, color.Gray{Y: 255 / 2})
				adjIn := multiply(out, color.Gray{Y: 255 / 2})
				newOut := offsetAdd(adjIn, adjMult, x, y)
				out = offsetAdd(out, newOut, 0, 0)

				// outputFile is a File type which satisfies Writer interface
				path := fmt.Sprintf("output.%d-%d.png", x, y)
				outputFile, err := os.Create(path)
				if err != nil {
					panic("Failed to output image")
				}

				// Encode takes a writer interface and an image interface
				// We pass it the File and the RGBA
				png.Encode(outputFile, out)

				// Don't forget to close files
				outputFile.Close()
			}
		}
	}
	// out = multiply(out, color.Gray16{Y: maskTotal / 3})

	return out
}

func offsetAdd(bg *image.RGBA, fg *image.RGBA, xOffset int, yOffset int) (out *image.RGBA) {
	bgBounds := bg.Bounds()
	width := bgBounds.Dx()
	height := bgBounds.Dy()

	out = image.NewRGBA(image.Rect(0, 0, width, height))

	// Add the two images, offseting the fg image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Wrap around if we go over
			adjY := y + yOffset
			adjX := x + xOffset
			if adjY >= height {
				adjY = adjY - height
			}
			if adjX >= width {
				adjX = adjX - width
			}

			// take the background pixel and the shifted fg pixel and add them
			bgPos := y*bg.Stride + x*4
			fgPos := adjY*fg.Stride + adjX*4
			result := add(
				fcolor.NewRGBAF64(bg.Pix[bgPos+0], bg.Pix[bgPos+1], bg.Pix[bgPos+2], bg.Pix[bgPos+3]),
				fcolor.NewRGBAF64(fg.Pix[fgPos+0], fg.Pix[fgPos+1], fg.Pix[fgPos+2], fg.Pix[fgPos+3]))

			result.Clamp()
			outPos := y*out.Stride + x*4
			out.Pix[outPos+0] = uint8(result.R * 255)
			out.Pix[outPos+1] = uint8(result.G * 255)
			out.Pix[outPos+2] = uint8(result.B * 255)
			out.Pix[outPos+3] = uint8(result.A * 255)
		}
	}

	return out
}

func multiply(in *image.RGBA, mult color.Color) (out *image.RGBA) {
	bounds := in.Bounds()
	r, g, b, a := mult.RGBA()
	multColor := fcolor.NewRGBAF64(uint8(r), uint8(g), uint8(b), uint8(a))

	width := bounds.Dx()
	height := bounds.Dy()

	out = image.NewRGBA(bounds)

	parallel.Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				pos := y*in.Stride + x*4
				result := colorMult(
					fcolor.NewRGBAF64(in.Pix[pos+0], in.Pix[pos+1], in.Pix[pos+2], in.Pix[pos+3]),
					multColor)

				result.Clamp()
				outPos := y*out.Stride + x*4
				out.Pix[outPos+0] = uint8(result.R * 255)
				out.Pix[outPos+1] = uint8(result.G * 255)
				out.Pix[outPos+2] = uint8(result.B * 255)
				out.Pix[outPos+3] = uint8(result.A * 255)
			}
		}
	})

	return out
}

func add(c0, c1 fcolor.RGBAF64) fcolor.RGBAF64 {
	r := c0.R + c1.R
	g := c0.G + c1.G
	b := c0.B + c1.B

	c2 := fcolor.RGBAF64{R: r, G: g, B: b, A: c1.A}
	return alphaComp(c0, c2)
}

func colorMult(c0, c1 fcolor.RGBAF64) fcolor.RGBAF64 {
	r := c0.R * c1.R
	g := c0.G * c1.G
	b := c0.B * c1.B

	c2 := fcolor.RGBAF64{R: r, G: g, B: b, A: c1.A}
	return alphaComp(c0, c2)
}

// alphaComp returns a new color after compositing the two colors
// based on the foreground's alpha channel.
func alphaComp(bg, fg fcolor.RGBAF64) fcolor.RGBAF64 {
	fg.Clamp()
	fga := fg.A

	r := (fg.R * fga / 1) + ((1 - fga) * bg.R / 1)
	g := (fg.G * fga / 1) + ((1 - fga) * bg.G / 1)
	b := (fg.B * fga / 1) + ((1 - fga) * bg.B / 1)
	a := bg.A + fga

	return fcolor.RGBAF64{R: r, G: g, B: b, A: a}
}
