package algo

type Dot struct {
	PixelThresholdPoints SquareMatrix[byte]
	min, max             byte
}

func NewDot(size int, inverted bool, globalMin, globalMax byte) *Dot {
	currentMin := max(globalMin, (globalMax-globalMin)/2+1)
	currentMax := min(globalMax, maxPixelValue)
	if inverted {
		currentMin = max(globalMin, 0)
		currentMax = min(globalMax, (globalMax-globalMin)/2)
	}

	matrix := NewDotPixelMatrix(size, inverted, currentMin, currentMax)

	return &Dot{
		PixelThresholdPoints: matrix,
		min:                  currentMin,
		max:                  currentMax,
	}
}

// IsPixelBlack returns true if black and false when white.
func (d *Dot) IsPixelBlack(x, y int, grayColor byte) bool {
	value := d.PixelThresholdPoints.Get(x, y)

	return grayColor < value
}
