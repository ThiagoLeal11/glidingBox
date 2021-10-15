package services


type Pixel struct {
	R, G, B uint8
}

type RawImage struct {
	Matrix
	c int
}


func MakeRawImage(size [2]int) RawImage {
	m := NewMatrix(
		[2]int{size[0], size[1] * 3},
	)

	m.Size = size
	return RawImage{
		m,
		3,
	}
}


func (m *RawImage) PixelAt(y, x int) Pixel {
	// Correct x value
	x = x * 3

	return Pixel{
		m.Data[y][x+0],
		m.Data[y][x+1],
		m.Data[y][x+2],
	}
}

func pixelMax(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func Max_uint8(a, b, c uint8) uint8 {
	return pixelMax(a, pixelMax(b, c))
}

func (m *RawImage) PixelSet(y, x int, value Pixel) {
	// Correct x value
	x = x * 3

	m.Data[y][x+0] = value.R
	m.Data[y][x+1] = value.G
	m.Data[y][x+2] = value.B
}

func (m *RawImage) PixelSliceRow(y, x, size int) Array {
	return m.SliceRow(y, x * 3, size * 3)
}
