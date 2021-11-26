package src

import (
	"glidingBox/buffers"
)

// GlidingBox convoluted the image with a box of size radius
func GlidingBox(img buffers.Image, diameter int) []int {
	height, width := img.GetShape()

	// Create the temp results vector
	results := buffers.NewVector(width)
	occurrences := make([]int, diameter*diameter)

	radius := diameter / 2

	// For each image row
	for y := radius; y < height-radius; y++ {
		results.Clean()
		// For each box row
		for l := -radius; l < radius+1; l++ {

			computePartialCenterDiff(img, diameter, radius, width, y, l, results)

		}

		// Group results (ignoring the left and right edges)
		for i := radius; i < width-radius; i++ {
			occurrences[results.At(i)-1]++
		}
	}

	return occurrences
}

func computePartialCenterDiff(img buffers.Image, diameter int, radius int, width int, y int, l int, results buffers.Vector) {
	// For each complete center box
	for c := radius; c < width-radius; c++ {
		center := img.GetPixel(y, c)
		// For each focus point in the box
		for f := -radius; f < radius+1; f++ {
			mass := img.GetPixel(y+l, c+f).SubAbs(center).Max()
			// mass is valid if is greater than the kernelSize
			if mass <= uint8(diameter) {
				results.AddAt(c, 1)
			}
		}
	}
}

func GlidingBoxSimple(m buffers.Image, diameter int) []float64 {
	height, width := m.GetShape()
	occurrences := make([]int32, diameter*diameter)
	radius := diameter / 2

	for y := radius; y < height-radius; y++ {
		for x := radius; x < width-radius; x++ {
			center := m.GetPixel(y, x)
			mass := -1
			for j := y - radius; j < y+radius+1; j++ {
				for i := x - radius; i < x+radius+1; i++ {
					focus := m.GetPixel(j, i)
					if focus.SubAbs(center).Max() <= uint8(diameter) {
						mass += 1
					}
				}
			}
			occurrences[mass] += 1
		}
	}

	boxInsideImg := (width - diameter + 1) * (height - diameter + 1)
	var probability []float64

	for _, c := range occurrences {
		probability = append(probability, float64(c)/float64(boxInsideImg))
	}
	return probability
}
