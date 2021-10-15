package services

type Array struct {
	Data []uint8
	Size int
}

type Matrix struct {
	Data [][]uint8
	Size [2]int
}

func NewArray(size int) Array {
	return Array{
		Data: make([]uint8, size),
		Size: size,
	}
}

func NewMatrix(size [2]int) Matrix {
	m := Matrix{
		Data: make([][]uint8, size[0]),
		Size: size,
	}
	for i := range m.Data {
		m.Data[i] = NewArray(size[1]).Data
	}
	return m
}


func (a *Array) At(x int) uint8 {
	return a.Data[x]
}

func (a *Array) Set(x int, value uint8) {
	a.Data[x] = value
}

func (m *Matrix) At(y, x int) uint8 {
	return m.Data[y][x]
}

func (m *Matrix) Set(y, x int, value uint8) {
	m.Data[y][x] = value
}

func (m *Matrix) SliceRow(y, x, size int) Array {
	slice := NewArray(size)
	for k := 0; k < size; k++ {
		slice.Data[k] = m.Data[y][x+k]
	}
	return slice
}

//func (m *Matrix) Slice(i0, i1, j0, j1 int) Matrix {
//
//}

