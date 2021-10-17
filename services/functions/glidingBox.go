package functions

import (
	"glidingBox/services/buffers"
	"unsafe"
)

const (
	tile      = 56 // The size of each row tile to iterate.
	realTile  = tile * 3
	innerTile = tile - 1 // -1 is to consider the center overlap
)

type ComputeTile struct {
	row        *[realTile]uint8
	center     *[realTile]uint8
	kernelSize int
}

func pixelInsideBox(t *ComputeTile, rowIdx, centerIdx int) uint8 {
	k := pIdx(rowIdx)
	c := pIdx(centerIdx)
	d0 := AbsDiff(t.row[k+0], t.center[c+0])
	d1 := AbsDiff(t.row[k+1], t.center[c+1])
	d2 := AbsDiff(t.row[k+2], t.center[c+2])

	if PixelMax(d0, d1, d2) <= uint8(t.kernelSize) {
		return 1
	}
	return 0
}

func computeCenterDiff(m buffers.RawImage, ry, cy, tx, tl, kernelSize int) buffers.Vector {
	tileSize := tl
	radius := kernelSize / 2

	result := buffers.NewVector(tileSize)

	if tileSize <= radius {
		return result
	}

	cTile := &ComputeTile{
		row:        (*[realTile]uint8)(unsafe.Pointer(&m.Data[ry][tx*3])),
		center:     (*[realTile]uint8)(unsafe.Pointer(&m.Data[cy][tx*3])),
		kernelSize: kernelSize,
	}

	// First idx in center (c = 0) and jumps the first column
	for k := 1; k < radius+1; k++ {
		mass := pixelInsideBox(cTile, k, 0)
		result.Add(0, mass)
	}

	// For each incomplete center box at left
	for c := 1; c < radius; c++ {
		for k := 0; k < min(radius+1+c, tileSize); k++ {
			mass := pixelInsideBox(cTile, k, c)
			result.Add(c, mass)
		}
	}

	// For each complete center box
	for c := radius; c < tileSize-radius; c++ {
		for k := -radius; k < radius+1; k++ {
			mass := pixelInsideBox(cTile, k+c, c)
			result.Add(c, mass)
		}
	}

	// For each incomplete center box at right
	for c := max(radius, tileSize-radius); c < tileSize; c++ {
		for k := c - radius; k < tileSize; k++ {
			mass := pixelInsideBox(cTile, k, c)
			result.Add(c, mass)
		}
	}

	return result
}

// GlidingBox convoluted the image with a box of size radius
func GlidingBox(m buffers.RawImage, diameter int) []int32 {
	height, width := m.GetShape()

	// Create the temp results vector
	results := make([]int32, width)
	occurrences := make([]int32, diameter*diameter)

	radius := diameter / 2
	numberOfTiles := ceil(floatDivision(width, innerTile))

	// For each line
	for y := 0; y < height-diameter+1; y++ {
		ClearArray(results)

		// For each row of the box
		for l := 0; l < diameter; l++ {

			// For each tile in a row
			for t := 0; t < numberOfTiles; t++ {
				// Compute the tile properties
				tileStart := t * (innerTile)
				tileEnd := min(tileStart+tile, width)
				tileSize := tileEnd - tileStart

				// Compute the result
				resultBuffer := computeCenterDiff(m, y+l, y+radius, tileStart, tileSize, diameter)

				// Save the results (start after any overlap grater than one)
				for x := tileStart; x < tileEnd; x++ {
					resultIdx := x - tileStart
					results[x] += int32(resultBuffer.At(resultIdx))
				}
			}
		}

		// Group results (ignoring the left and right edges)
		for i := radius; i < width-radius; i++ {
			occurrences[results[i]-1]++
		}
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
