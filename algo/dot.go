package algo

import "math"

type DotMatrixBuilder struct {
	Size                int
	PointThresholdOrder []Point2D[int]
}

func newDotBuilder(size int, inverted bool) *DotMatrixBuilder {
	order := generateCirclePointOrder(size, inverted)

	return &DotMatrixBuilder{
		Size:                size,
		PointThresholdOrder: order,
	}
}

func (d DotMatrixBuilder) generateMatrix(min, max byte) SquareMatrix[byte] {
	resultMatrix := newEmptyPixelMatrix[byte](d.Size)

	// NOTE(ShaDream): example: min = 1, max = 3. points are 2. first point = 1, second point = 3. step is equal 2. (3-1)/(2-x) = 2, x is 1.
	stepValue := float64(max-min) / float64(len(d.PointThresholdOrder)-1)

	for i, point := range d.PointThresholdOrder {
		newValue := math.Round(float64(max) - (stepValue * float64(i)))

		resultMatrix.set(point.X, point.Y, byte(newValue))
	}

	return resultMatrix
}
