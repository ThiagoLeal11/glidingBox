package src

import "math"

func floatDivision(a, b int) float64 {
	return float64(a) / float64(b)
}

func ceil(a float64) int {
	return int(math.Ceil(a))
}
