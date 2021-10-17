package buffers

type Matrix struct {
	Data  [][]uint8
	Shape [2]int
}

func NewMatrix(size [2]int) Matrix {
	m := Matrix{
		Data:  make([][]uint8, size[0]),
		Shape: size,
	}

	for i := range m.Data {
		m.Data[i] = NewVector(size[1]).Data
	}

	return m
}

func (m *Matrix) GetShape() (int, int) {
	return m.Shape[0], m.Shape[1]
}

func (m *Matrix) At(y, x int) uint8 {
	return m.Data[y][x]
}

func (m *Matrix) Set(y, x int, value uint8) {
	m.Data[y][x] = value
}

func (m *Matrix) SliceRow(y, x, size int) Vector {
	slice := NewVector(size)
	slice.Data = m.Data[y][x : x+size]
	return slice
}
