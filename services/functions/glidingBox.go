package functions

import (
	"fmt"
	"glidingBox/services/buffers"
)

const (
	tile      = 56       // The size of each row tile to iterate.
	innerTile = tile - 1 // -1 is to consider the center overlap
)

// Fast absDiff
func absDiff(a, b uint8) uint8 {
	res := int16(a) - int16(b)
	signBit := res >> ((2 << 3) - 1)
	value := (res ^ signBit) + (signBit & 1)
	return uint8(value)
}

func pixelInsideBox(d0, d1, d2 uint8, kernelSize int) uint8 {
	if PixelMax(d0, d1, d2) <= uint8(kernelSize) {
		return 1
	}
	return 0
}

func computeCenterDiff(rowBuffer, centerBuffer buffers.Vector, kernelSize int) buffers.Vector {
	tileSize := rowBuffer.Shape / 3
	radius := kernelSize / 2

	row := rowBuffer.Data
	center := centerBuffer.Data
	result := buffers.NewVector(tileSize)

	if tileSize <= radius {
		return result
	}

	// First idx in center (c = 0) and jumps the first column
	for k := 1; k < radius+1; k++ {
		K := pIdx(k)
		d0 := absDiff(row[K+0], center[0])
		d1 := absDiff(row[K+1], center[1])
		d2 := absDiff(row[K+2], center[2])
		result.Add(0, pixelInsideBox(d0, d1, d2, kernelSize))
	}

	// For each incomplete center box at left
	for c := 1; c < radius; c++ {
		C := pIdx(c)
		for k := 0; k < min(radius+1+c, tileSize); k++ {
			K := pIdx(k)
			d0 := absDiff(row[K+0], center[C+0])
			d1 := absDiff(row[K+1], center[C+1])
			d2 := absDiff(row[K+2], center[C+2])
			result.Add(c, pixelInsideBox(d0, d1, d2, kernelSize))
		}
	}

	// For each complete center box
	for c := radius; c < tileSize-radius; c++ {
		C := pIdx(c)
		if c == 52 {
			print("")
		}
		for k := -radius; k < radius+1; k++ {
			K := C + pIdx(k)
			d0 := absDiff(row[K+0], center[C+0])
			d1 := absDiff(row[K+1], center[C+1])
			d2 := absDiff(row[K+2], center[C+2])
			result.Add(c, pixelInsideBox(d0, d1, d2, kernelSize))
		}
	}

	// For each incomplete center box at right
	for c := max(radius, tileSize-radius); c < tileSize; c++ {
		C := pIdx(c)
		for k := c - radius; k < tileSize; k++ {
			K := pIdx(k)
			d0 := absDiff(row[K+0], center[C+0])
			d1 := absDiff(row[K+1], center[C+1])
			d2 := absDiff(row[K+2], center[C+2])
			result.Add(c, pixelInsideBox(d0, d1, d2, kernelSize))
		}
	}

	return result
}

// GlidingBox convoluted the image with a box of size radius
func GlidingBox(m buffers.RawImage, diameter int) []int32 {
	height, width := m.GetShape()

	// Create the temp results vector
	occurrences := make([]int32, diameter*diameter)

	radius := diameter / 2
	numberOfTiles := ceil(floatDivision(width, innerTile))

	var sim int32

	// For each line
	for y := 0; y < height-diameter; y++ {
		results := make([]int32, width)

		// For each row of the box
		for l := 0; l < diameter; l++ {

			// For each tile in a row
			for t := 0; t < numberOfTiles; t++ {
				// Compute the tile properties
				tileStart := t * (innerTile)
				tileEnd := min(tileStart+tile, width)
				tileSize := tileEnd - tileStart

				if t == 14 {
					print("")
				}

				// Start the buffers
				compBuffer := m.ImageSliceRow(y+l, tileStart, tileSize)
				centerBuffer := m.ImageSliceRow(y+radius, tileStart, tileSize)

				// Compute the result
				resultBuffer := computeCenterDiff(compBuffer, centerBuffer, diameter)

				// Save the results (start after any overlap grater than one)
				for x := tileStart; x < tileEnd; x++ {
					if x == 770 && y == 3 {
						fmt.Printf("")
					}
					resultIdx := x - tileStart
					value := int32(resultBuffer.At(resultIdx))
					results[x] += value
					r := results[x]
					sim = r
				}
			}
		}

		for i := radius; i < width-radius; i++ {
			if results[i]-1 == 4 {
				//fmt.Printf("%d ", i)
			}
		}

		// Group results (ignoring the left and right edges)
		for i := radius; i < width-radius; i++ {
			//fmt.Printf("%d ", results[i]-1)
			occurrences[results[i]-1]++
		}
		fmt.Print("")
	}

	fmt.Println(sim)

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
