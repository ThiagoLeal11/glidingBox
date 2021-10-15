package services

import (
	"glidingBox/services/buffers"
	"math"
)

const (
	tile      = 64       // The size of each row tile to iterate.
	innerTile = tile - 1 // -1 is to consider the center overlap
)

func ceilDivision(a, b int) int {
	return int(math.Ceil(float64(a) / float64(b)))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Range struct {
	start int
	stop  int
}

func pIdx(a int) int {
	return a * 3
}

func PixelMax(a, b, c uint8) uint8 {
	return buffers.MaxUi8(a, buffers.MaxUi8(b, c))
}

// GlidingBox convoluted the image with a box of size radius
func GlidingBox(m buffers.RawImage, diameter int) []int64 {
	height, width := m.GetShape()

	numberOfBox := width * height
	results := make([]int64, numberOfBox)

	radius := diameter / 2
	numberOfTiles := ceilDivision(width, innerTile)

	// For each line
	for y := 0; y < height-diameter; y++ {

		// For each row of the box
		for l := 0; l < diameter; l++ {

			// For each tile in a row
			for t := 0; t < numberOfTiles; t++ {
				// Compute the start of the tile (the tile is always completely inside the image)
				expectedStart := t * (innerTile)
				startX := min(expectedStart, width-tile)

				// Start the buffers
				compBuffer := m.SliceRow(y+l, startX, tile)
				centerBuffer := m.SliceRow(y+radius, startX, tile)
				resultBuffer := buffers.NewVector(tile)

				// For each incomplete center box at left
				for c := 0; c < radius; c++ {
					for k := 0; k < radius+1+c; k++ {
						diff0 := compBuffer.At(pIdx(k)+0) - centerBuffer.At(pIdx(c)+0)
						diff1 := compBuffer.At(pIdx(k)+1) - centerBuffer.At(pIdx(c)+1)
						diff2 := compBuffer.At(pIdx(k)+2) - centerBuffer.At(pIdx(c)+2)
						resultBuffer.Set(c, PixelMax(diff0, diff1, diff2))
					}
				}

				// For each complete center box
				for c := radius; c < tile-radius; c++ {
					// For each box
					for k := -radius; k < radius+1; k++ {
						diff0 := compBuffer.At(pIdx(c+k)+0) - centerBuffer.At(pIdx(c)+0)
						diff1 := compBuffer.At(pIdx(c+k)+1) - centerBuffer.At(pIdx(c)+1)
						diff2 := compBuffer.At(pIdx(c+k)+2) - centerBuffer.At(pIdx(c)+2)
						resultBuffer.Set(c, PixelMax(diff0, diff1, diff2))
					}
				}

				// For each incomplete center box at right
				for c := tile - radius; c < tile; c++ {
					// don't compute the last column of the box, to prevent same column overlap
					for k := c - radius; k < tile && !(k == tile-1 && c == tile-1); k++ {
						diff0 := compBuffer.At(pIdx(k)+0) - centerBuffer.At(pIdx(c)+0)
						diff1 := compBuffer.At(pIdx(k)+1) - centerBuffer.At(pIdx(c)+1)
						diff2 := compBuffer.At(pIdx(k)+2) - centerBuffer.At(pIdx(c)+2)
						resultBuffer.Set(c, PixelMax(diff0, diff1, diff2))
					}
				}

				// Normalize result vector
				for i := 0; i < resultBuffer.Shape; i++ {
					value := uint8(0)
					if resultBuffer.At(i) > uint8(diameter) {
						value = uint8(1)
					}
					resultBuffer.Set(i, value)
				}

				//Save results (start after the incomplete overlap from start)
				lastEndX := (t-1)*innerTile + tile
				resultStart := max(startX, lastEndX) // Start after the overlap and radius (prevent half box)
				resultEnd := startX + tile           // Stop after the radius (prevent half box)
				for i := resultStart - (innerTile * t) - 1; i < resultEnd-(innerTile*t); i++ {
					resultIndex := y*width + i
					results[resultIndex] += int64(resultBuffer.At(i))
				}
			}
		}
	}

	return results
}
