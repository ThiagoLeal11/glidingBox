package buffers

type RawImage struct {
	Matrix
	c int
}

func NewRawImage(shape [2]int) RawImage {
	m := NewMatrix(
		[2]int{shape[0], shape[1] * 3},
	)

	m.Shape = shape
	return RawImage{
		m,
		3,
	}
}

func (m *RawImage) ImageSliceRow(y, x, size int) Vector {
	return m.SliceRow(y, x*3, size*3)
}

type Pixel struct {
	r uint8
	b uint8
	g uint8
}

func (m *RawImage) GetPixel(y, x int) Pixel {
	xStart := x * 3
	return Pixel{
		m.At(y, xStart+0),
		m.At(y, xStart+1),
		m.At(y, xStart+2),
	}
}

func (p Pixel) MinusAbs(a Pixel) Pixel {
	return Pixel{
		minusAbs(p.r, a.r),
		minusAbs(p.g, a.g),
		minusAbs(p.b, a.b),
	}
}

func minusAbs(a, b uint8) uint8 {
	if a < b {
		return b - a
	}
	return a - b
}

func MaxPixel(p Pixel) uint8 {
	return MaxUi8(p.r, MaxUi8(p.g, p.b))
}
