package algo

type Dot struct {
	PixelThresholdPoints Matrix[byte]
	Size                 int
}

func NewDot(size int, inverted bool, globalMin, globalMax byte) *Dot {
	currentMin := max(globalMin, centerValue)
	currentMax := min(globalMax, maxPixelValue)
	if inverted {
		currentMin = max(globalMin, 0)
		currentMax = min(globalMax, centerValue-1)
	}

	matrix := NewDotPixelMatrix(size, inverted, currentMin, currentMax)

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
