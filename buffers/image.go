package buffers

// Planar
// ========
// RRRRRRRR
// RRRRRRRR
// RRRRRRRR
// RRRRRRRR
// GGGGGGGG
// GGGGGGGG
// GGGGGGGG
// GGGGGGGG
// BBBBBBBB
// BBBBBBBB
// BBBBBBBB
// BBBBBBBB

// PlanarAlternate
// ========
// RRRRRRRR
// GGGGGGGG
// BBBBBBBB
// RRRRRRRR
// GGGGGGGG
// BBBBBBBB
// RRRRRRRR
// GGGGGGGG
// BBBBBBBB
// RRRRRRRR
// GGGGGGGG
// BBBBBBBB

// Interleaved
// ========
// RGBRGBRGBRGBRGBRGBRGBRGB
// RGBRGBRGBRGBRGBRGBRGBRGB
// RGBRGBRGBRGBRGBRGBRGBRGB
// RGBRGBRGBRGBRGBRGBRGBRGB

const (
	Planar = iota
	PlanarAlternate
	Interleaved
)

type Image struct {
	Matrix
	Layout   int
	Channels int
	StrideX  int
	StrideY  int
}

type Pixel struct {
	r uint8
	g uint8
	b uint8
}

// NewImage create a new image using a layout pattern given the shape
func NewInterleavedImage(shape [2]int) Image {
	m := NewMatrix(
		[2]int{shape[0], shape[1] * 3},
	)

	m.Shape = shape
	return Image{
		m,
		Interleaved,
		3,
		3,
		1,
	}
}

func (img *Image) GetPixel(y, x int) Pixel {
	dim0 := y * img.StrideY
	dim1 := x * img.StrideX

	return Pixel{
		img.At(dim0, dim1+0),
		img.At(dim0, dim1+1),
		img.At(dim0, dim1+2),
	}
}

func (img *Image) SetPixel(y, x int, r, g, b uint32) {
	interleavedX := x * 3
	img.Set(y, interleavedX+0, uint8(r))
	img.Set(y, interleavedX+1, uint8(g))
	img.Set(y, interleavedX+2, uint8(b))
}

func (p Pixel) SubAbs(a Pixel) Pixel {
	return Pixel{
		minusAbs(p.r, a.r),
		minusAbs(p.g, a.g),
		minusAbs(p.b, a.b),
	}
}

func (p Pixel) Max() uint8 {
	return maxUi8(p.r, maxUi8(p.g, p.b))
}

func minusAbs(a, b uint8) uint8 {
	if a < b {
		return b - a
	}
	return a - b
}

func maxUi8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}
