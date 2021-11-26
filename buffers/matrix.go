package buffers

// Matrix is a container for a 2D tensor encoded as a 1D tensor the shape describes respectively (height, width)
type Matrix struct {
	Data  [][]uint8
	Shape [2]int
}

// NewMatrix allocates a matrix with the given shape in [2]int{height, width}
func NewMatrix(size [2]int) Matrix {
	m := Matrix{
		Data:  make([][]uint8, size[0]),
		Shape: size,
	}

	for i := range m.Data {
		m.Data[i] = make([]uint8, size[1])
	}

	return m
}

// GetShape returns the height and the width of a matrix.
func (m *Matrix) GetShape() (int, int) {
	return m.Shape[0], m.Shape[1]
}

// At returns the value for the line y and column x
func (m *Matrix) At(y, x int) uint8 {
	return m.Data[y][x]
}

// Set persists a value for the line y and column x
func (m *Matrix) Set(y, x int, value uint8) {
	m.Data[y][x] = value
}
