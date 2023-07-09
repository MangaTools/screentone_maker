package algo

import "math"

type Dot struct {
	PixelThresholdPoints PixelMatrix
	Size                 int
}

func NewDot(size int, isAnti bool) *Dot {
	matrix := NewPixelMatrix(size)
	if isAnti {
		matrix.Change(func(previous byte, x, y int) byte {
			return byte(math.Max(float64(maxPixelValue-previous), 1))
		})
	}
	return &Dot{
		PixelThresholdPoints: matrix,
		Size:                 size,
	}
}

// IsPixelBlack returns true if black and false when white.
func (d *Dot) IsPixelBlack(x, y int, grayColor byte) bool {
	value := d.PixelThresholdPoints.Get(x, y)

	return grayColor < value
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

const (
	maxPixelValue = 255
	centerValue   = 128
)
