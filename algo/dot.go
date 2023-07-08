package algo

type Dot struct {
	PixelThresholdPoints PixelMatrix
	Size                 int
}

func NewDot(size int) *Dot {
	return &Dot{
		PixelThresholdPoints: NewPixelMatrix(size),
		Size:                 size,
	}
}

// IsPixelBlack returns true if black and false when white.
func (d *Dot) IsPixelBlack(x, y int, grayColor byte, anti bool) bool {
	value := d.PixelThresholdPoints.Get(x, y)

	if anti {
		return grayColor <= maxPixelValue-value
	}
	return grayColor < value
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

const (
	maxPixelValue = 255
	centerValue   = 128
)
