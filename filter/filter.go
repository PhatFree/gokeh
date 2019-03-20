package filter

import (
	"image"
	"image/color"

	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/fcolor"
	"github.com/anthonynsimon/bild/parallel"
)

func ApplyBlur(in image.RGBA64, mask image.Gray16) (out *image.RGBA64) {
	// Copy the input image into the output image
	out = image.NewRGBA64(in.Bounds())

	// Look at each pixel in the mask
	bounds := mask.Bounds()
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; x < bounds.Dy(); y++ {
			maskPixel := mask.At(x, y).(color.Gray16)
			// Make sure that mask pixel isn't black
			// TODO: maybe also check if it's just really dim
			if maskPixel.Y != 0 {
				// accumulate a shifted copy of the color image weighted by the mask
				multLayer := multiply(in, maskPixel)
				out = offsetAdd(out, multLayer, x, y)
			}

		}
	}

	return out
}

func offsetAdd(bg image.Image, fg image.Image, xOffset int, yOffset int) (out *image.RGBA64) {
	bgBounds := bg.Bounds()
	width := bgBounds.Dx()
	height := bgBounds.Dy()

	bgSrc := clone.AsRGBA(bg)
	fgSrc := clone.AsRGBA(fg)
	out = image.NewRGBA64(image.Rect(0, 0, width, height))

	// Add the two images, offseting the fg image
	parallel.Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			// Make ser we don't go off the bottom
			if y+yOffset > end {
				break
			}

			for x := 0; x < width; x++ {
				// make sure we don't go off the right
				if x+xOffset > width {
					break
				}

				// take the background pixel and the shifted fg pixel and add them
				bgPos := y*bgSrc.Stride + x*4
				fgPos := (y+yOffset)*fgSrc.Stride + (x+xOffset)*4
				result := add(
					fcolor.NewRGBAF64(bgSrc.Pix[bgPos+0], bgSrc.Pix[bgPos+1], bgSrc.Pix[bgPos+2], bgSrc.Pix[bgPos+3]),
					fcolor.NewRGBAF64(fgSrc.Pix[fgPos+0], fgSrc.Pix[fgPos+1], fgSrc.Pix[fgPos+2], fgSrc.Pix[fgPos+3]))

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

func multiply(in image.RGBA64, mult color.Color) (out *image.RGBA) {
	bounds := in.Bounds()
	r, g, b, a := mult.RGBA()
	multColor := fcolor.NewRGBAF64(uint8(r), uint8(g), uint8(b), uint8(a))

	width := bounds.Dx()
	height := bounds.Dy()

	inSrc := clone.AsRGBA(&in)
	out = image.NewRGBA(bounds)

	parallel.Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				pos := y*inSrc.Stride + x*4
				result := colorMult(
					fcolor.NewRGBAF64(inSrc.Pix[pos+0], inSrc.Pix[pos+1], inSrc.Pix[pos+2], inSrc.Pix[pos+3]),
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

// Multiply combines the foreground and background images by multiplying their
// normalized values and returns the resulting image.
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
