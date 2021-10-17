package functions

import (
	"glidingBox/services/buffers"
	"math"
)

func floatDivision(a, b int) float64 {
	return float64(a) / float64(b)
}

func ceil(a float64) int {
	return int(math.Ceil(a))
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

// Legacy
func pIdx(a int) int {
	return a * 3
}

func PixelMax(a, b, c uint8) uint8 {
	return buffers.MaxUi8(a, buffers.MaxUi8(b, c))
}
