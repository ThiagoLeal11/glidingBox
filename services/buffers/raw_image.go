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
