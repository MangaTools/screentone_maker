package algo

import "sort"

const (
	horizontalOffset = .1
	verticalOffset   = .15
)

// generateCirclePointOrder generates array of points with it's weight. points are not ordered.
func generateCirclePointOrder(size int, inverted bool) []Point2D[int] {
	centralValue := float64(size-1) / 2
	centerPoint := Point2D[float64]{X: centralValue + horizontalOffset, Y: centralValue + verticalOffset}

	valuePoints := make([]PointValue2D[int, float64], 0, size*size)

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			value := getDotPixelLevel(size, centerPoint, float64(x)+.5, float64(y)+.5)
			point := PointValue2D[int, float64]{
				Point2D: Point2D[int]{
					X: x,
					Y: y,
				},
				Value: value,
			}

			valuePoints = append(valuePoints, point)
		}
	}

	sortPoints(valuePoints, inverted)

	resultPoints := make([]Point2D[int], 0, len(valuePoints))
	for _, point := range valuePoints {
		resultPoints = append(resultPoints, point.Point2D)
	}

	return resultPoints
}

// getDotPixelLevel calculate how far pixel from center point.
func getDotPixelLevel(size int, point Point2D[float64], pixelX, pixelY float64) float64 {
	distance := point.Distance(pixelX, pixelY)
	circleRadius := getMaxDistance(point, size)

	return float64(1) - distance/circleRadius
}

// getMaxDistance calculates farthes distance between center of point and square verticies.
func getMaxDistance(point Point2D[float64], size int) float64 {
	max := float64(0)

	maxSize := float64(size - 1)

	corners := []Point2D[float64]{{X: 0, Y: 0}, {X: 0, Y: maxSize}, {X: maxSize, Y: maxSize}, {X: maxSize, Y: 0}}
	for _, corner := range corners {
		value := point.Distance(corner.X, corner.Y)
		if max < value {
			max = value
		}
	}

	return max
}

// sortPoints according to value.
func sortPoints(points []PointValue2D[int, float64], inverted bool) {
	if inverted {
		sort.Slice(points, func(i, j int) bool { return points[i].Value < points[j].Value })
	} else {
		sort.Slice(points, func(i, j int) bool { return !(points[i].Value < points[j].Value) })
	}
}
