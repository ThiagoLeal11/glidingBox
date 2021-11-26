package buffers

// Vector is a container for a 2D tensor
type Vector struct {
	Data  []uint8
	Shape int
}

// NewVector allocates a matrix with the given size
func NewVector(size int) Vector {
	return Vector{
		Data:  make([]uint8, size),
		Shape: size,
	}
}

// At returns the value at the given index x
func (a *Vector) At(x int) uint8 {
	return a.Data[x]
}

// Set persists a value for the given index x
func (a *Vector) Set(x int, value uint8) {
	a.Data[x] = value
}

// AddAt sum the value to the value at the given index x
func (a *Vector) AddAt(x int, value uint8) {
	a.Data[x] += value
}

func (a *Vector) Clean() {
	for i := range a.Data {
		a.Data[i] = 0
	}
}
