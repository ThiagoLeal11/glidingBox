package buffers

type Vector struct {
	Data  []uint8
	Shape int
}

func NewVector(size int) Vector {
	return Vector{
		Data:  make([]uint8, size),
		Shape: size,
	}
}

func (a *Vector) At(x int) uint8 {
	return a.Data[x]
}

func (a *Vector) Set(x int, value uint8) {
	a.Data[x] = value
}

func (a *Vector) Add(x int, value uint8) {
	a.Data[x] += value
}
