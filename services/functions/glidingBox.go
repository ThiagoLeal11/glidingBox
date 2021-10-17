package functions

import (
	"glidingBox/services/buffers"
)

const (
	tile      = 56       // The size of each row tile to iterate.
	innerTile = tile - 1 // -1 is to consider the center overlap
)

func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

func computeCenterDiff(rowBuffer, centerBuffer buffers.Vector, kernelSize int) buffers.Vector {
	result := buffers.NewVector(rowBuffer.Shape / 3)

	row := rowBuffer.Data
	center := centerBuffer.Data
	radius := kernelSize / 2

	// First idx in center (c = 0) and jumps the first column
	for k := 1; k < radius+1; k++ {
		K := pIdx(k)
		d0 := absDiff(row[K+0], center[0])
		d1 := absDiff(row[K+1], center[1])
		d2 := absDiff(row[K+2], center[2])
		result.Set(0, PixelMax(d0, d1, d2))
	}

	// For each incomplete center box at left
	for c := 1; c < radius; c++ {
		C := pIdx(c)
		for k := 0; k < radius+1+c; k++ {
			K := pIdx(k)
			d0 := absDiff(row[K+0], center[C+0])
			d1 := absDiff(row[K+1], center[C+1])
			d2 := absDiff(row[K+2], center[C+2])
			result.Set(c, PixelMax(d0, d1, d2))
		}
	}

	// For each complete center box
	for c := radius; c < tile-radius; c++ {
		C := pIdx(c)
		for k := -radius; k < radius+1; k++ {
			K := C + pIdx(k)
			d0 := absDiff(row[K+0], center[C+0])
			d1 := absDiff(row[K+1], center[C+1])
			d2 := absDiff(row[K+2], center[C+2])
			result.Set(c, PixelMax(d0, d1, d2))
		}
	}

	// For each incomplete center box at right
	for c := tile - radius; c < tile; c++ {
		C := pIdx(c)
		for k := c - radius; k < tile; k++ {
			K := pIdx(k)
			d0 := absDiff(row[K+0], center[C+0])
			d1 := absDiff(row[K+1], center[C+1])
			d2 := absDiff(row[K+2], center[C+2])
			result.Set(c, PixelMax(d0, d1, d2))
		}
	}

	// Normalize result vector
	for i := 0; i < result.Shape; i++ {
		value := uint8(1)
		if result.At(i) > uint8(kernelSize) {
			value = uint8(0)
		}
		result.Set(i, value)
	}

	return result
}

// GlidingBox convoluted the image with a box of size radius
func GlidingBox(m buffers.RawImage, diameter int) []int32 {
	height, width := m.GetShape()

	// Create the temp results vector
	results := make([][]int32, height)
	for i := range results {
		results[i] = make([]int32, width)
	}

	radius := diameter / 2
	numberOfTiles := ceil(floatDivision(width, innerTile))

	// For each line
	for y := 0; y < height-diameter; y++ {

		// For each row of the box
		for l := 0; l < diameter; l++ {

			// For each tile in a row
			for t := 0; t < numberOfTiles; t++ {
				// Compute the start of the tile (the tile is always completely inside the image)
				expectedStart := t * (innerTile)
				start := min(expectedStart, width-tile)

				// Start the buffers
				compBuffer := m.ImageSliceRow(y+l, start, tile)
				centerBuffer := m.ImageSliceRow(y+radius, start, tile)

				// Compute the result
				resultBuffer := computeCenterDiff(compBuffer, centerBuffer, diameter)

				// Save the results (start after any overlap grater than one)
				for x := expectedStart; x < start+tile; x++ {
					resultIdx := x - (innerTile * t)
					results[y][x] += int32(resultBuffer.At(resultIdx))
				}
			}
		}
	}

	// Clean results
	resultSize := (width - diameter + 1) * (height - diameter + 1)
	innerResult := make([]int32, resultSize)
	for y := radius; y < height-radius; y++ {
		for x := radius; x < width-radius; x++ {
			innerIdx := y + x - (radius * 2)
			innerResult[innerIdx] = results[y][x]
		}

	}

	// groupResults
	occurrences := make([]int32, diameter*diameter)
	for i := 0; i < resultSize; i++ {
		occurrences[innerResult[i]]++
	}

	return occurrences
}

func GlidingBoxSimple(m buffers.RawImage, diameter int) []int32 {
	height, width := m.GetShape()
	occurrences := make([]int32, diameter*diameter)
	radius := diameter / 2

	for y := radius; y < height-radius; y++ {
		for x := radius; x < width-radius; x++ {
			central := m.GetPixel(y, x)
			mass := 0
			for j := y - radius; j < y+radius+1; j++ {
				for i := x - radius; i < x+radius+1; i++ {
					if j == y && i == x {
						continue
					}
					pixel := m.GetPixel(j, i)
					if buffers.MaxPixel(pixel.MinusAbs(central)) <= uint8(diameter) {
						mass += 1
					}
				}
			}
			occurrences[mass] += 1
		}
	}

	return occurrences
}
