package services

import (
	"fmt"
	"glidingBox/services/buffers"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func OpenImage(path string) (image.Image, error) {
	var err error

	// Open file
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("can't open the file: %v", err)
		return nil, err
	}
	defer file.Close()

	// Ger filename and suffix
	fileMetadata, _ := file.Stat()
	fileName := fileMetadata.Name()
	parts := strings.Split(fileName, ".")
	suffix := parts[len(parts)-1]

	// Decode file based on suffix
	var loadedImg image.Image
	if strings.ToLower(suffix) == "png" {
		loadedImg, err = png.Decode(file)
	} else if strings.ToLower(suffix) == "jpeg" || strings.ToLower(suffix) == "jpg" {
		loadedImg, err = jpeg.Decode(file)
	} else {
		fmt.Printf("image format unknow. Tried png, jpg and jpeg")
	}

	if err != nil {
		fmt.Printf("decoding error: %v", err.Error())
		return nil, err
	}

	return loadedImg, nil
}

// EncodeInterleave Return an interleaved RGB image as an uint8 array.
func EncodeInterleave(img image.Image) buffers.RawImage {
	// The number of color channels in the image
	channels := 3

	size := img.Bounds().Size()
	colorModel := color.NRGBAModel
	m := buffers.NewRawImage(
		[2]int{size.Y, size.X},
	)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			// Reshape tensor (Y, X, channels into: Y, X * channels)
			pixel := img.At(x, y)
			result := colorModel.Convert(pixel)
			r, g, b, _ := result.RGBA()

			X := x * channels
			m.Set(y, X+0, uint8(r))
			m.Set(y, X+1, uint8(g))
			m.Set(y, X+2, uint8(b))
		}
	}

	return m
}
