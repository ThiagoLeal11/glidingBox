package src

import "glidingBox/buffers"

func labelClusters(b buffers.Matrix, c buffers.Matrix, x, y, label int) int {
	if c.At(y, x) == 0 && b.At(y, x) != 0 {
		height, width := b.GetShape()
		c.Set(y, x, uint8(label))

		canLookUp := (x - 1) >= 0
		canLookRight := (y + 1) < width
		canLookLeft := (y - 1) >= 0
		canLookDown := (x + 1) < height

		count := 0

		if canLookUp {
			count += labelClusters(b, c, x-1, y, label)
		}
		if canLookLeft {
			count += labelClusters(b, c, x, y-1, label)
		}
		if canLookRight {
			count += labelClusters(b, c, x, y+1, label)
		}
		if canLookDown {
			count += labelClusters(b, c, x+1, y, label)
		}

		return count + 1
	}
	return 0
}

func clusterize(b buffers.Matrix, c buffers.Matrix) (int, int) {
	height, width := b.GetShape()
	currentLabel := 1

	max := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if b.At(y, x) == 0 {
				continue
			}

			if c.At(y, x) == 0 {
				candidate := labelClusters(b, c, x, y, currentLabel)
				if candidate > max {
					max = candidate
				}

				currentLabel++
			}
		}
	}

	return currentLabel - 1, max
}

func PercolationSimple(m buffers.Image, diameter int) LocalPercolationData {
	height, width := m.GetShape()
	radius := diameter / 2

	boxCount := (width - diameter + 1) * (height - diameter + 1)

	clusterCounts := make([]int, boxCount)
	biggestClusterAreas := make([]int, boxCount)
	percolationBoxCounters := 0

	region := buffers.NewMatrix([2]int{diameter, diameter})
	clusters := buffers.NewMatrix([2]int{diameter, diameter})

	idx := 0

	boxArea := diameter * diameter

	for y := radius; y < height-radius; y++ {
		for x := radius; x < width-radius; x++ {
			center := m.GetPixel(y, x)
			percolatingPixels := 0

			for j := y - radius; j < y+radius+1; j++ {
				for i := x - radius; i < x+radius+1; i++ {
					focus := m.GetPixel(j, i)

					mi := i - x + radius
					mj := j - y + radius

					region.Set(mj, mi, 0)
					clusters.Set(mj, mi, 0)

					if focus.SubAbs(center).Max() <= uint8(diameter) {
						percolatingPixels += 1
						region.Set(mj, mi, 1)
					}
				}
			}
			clustersOnBox, biggestClusterSize := clusterize(region, clusters)
			clusterCounts[idx] = clustersOnBox
			biggestClusterAreas[idx] = biggestClusterSize / boxArea
			idx++

			if float64(percolatingPixels/boxArea) >= 0.59275 {
				percolationBoxCounters += 1
			}

		}
	}

	// Compute the statistics
	boxAverageClusterCount := mean(clusterCounts)
	boxAverageBiggestClusterArea := mean(biggestClusterAreas)
	boxPercolation := float64(percolationBoxCounters) / float64(boxCount)

	return LocalPercolationData{
		AverageClusterCount:       boxAverageClusterCount,
		AverageBiggestClusterArea: boxAverageBiggestClusterArea,
		Percolation:               boxPercolation,
		KernelSize:                diameter,
	}
}

type LocalPercolationData struct {
	AverageClusterCount       float64 `json:"average_cluster_count"`
	AverageBiggestClusterArea float64 `json:"average_biggest_cluster_area"`
	Percolation               float64 `json:"percolation"`
	KernelSize                int     `json:"kernel_size"`
}

func mean(v []int) float64 {
	acc := 0
	for _, value := range v {
		acc += value
	}

	return float64(acc) / float64(len(v))
}
