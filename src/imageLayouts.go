package src

import (
	"glidingBox/buffers"
	"image"
	"image/color"
)

// EncodeInterleave Return an interleaved RGB image as an uint8 matrix.
func EncodeInterleave(img image.Image) buffers.Image {
	size := img.Bounds().Size()
	colorModel := color.NRGBAModel
	m := buffers.NewInterleavedImage(
		[2]int{size.Y, size.X},
	)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			pixel := img.At(x, y)
			result := colorModel.Convert(pixel)
			r, g, b, _ := result.RGBA()

			m.SetPixel(y, x, r, g, b)
		}
	}

	return m
}
