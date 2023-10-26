package algo

import (
	"math"
	"sort"
)

type SquareMatrix[T any] struct {
	Size   int
	Matrix []T
}

func (pm *SquareMatrix[T]) Get(x, y int) T {
	return pm.Matrix[pm.getIndex(x, y)]
}

func (pm *SquareMatrix[T]) Set(x, y int, value T) {
	index := pm.getIndex(x, y)
	pm.Matrix[index] = value
}

func (pm *SquareMatrix[T]) getIndex(x, y int) int {
	return x + (y * pm.Size)
}

func (pm *SquareMatrix[T]) getPosition(index int) (int, int) {
	return index % pm.Size, index / pm.Size
}

func (pm *SquareMatrix[T]) Change(chageFunc func(previous T, x, y int) T) {
	for x := 0; x < pm.Size; x++ {
		for y := 0; y < pm.Size; y++ {
			index := pm.getIndex(x, y)
			pm.Matrix[index] = chageFunc(pm.Matrix[index], x, y)
		}
	}
}

func newEmptyPixelMatrix[T any](size int) SquareMatrix[T] {
	matrix := SquareMatrix[T]{
		Size:   size,
		Matrix: make([]T, size*size),
	}

	return matrix
}

func NewMatrixFrom2DSlices[T any](values [][]T) SquareMatrix[T] {
	size := len(values)

	matrix := newEmptyPixelMatrix[T](size)

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			matrix.Set(x, y, values[x][y])
		}
	}

	return matrix
}

func NewMatrixWithSetupFunc[T any](size int, setupFunc func(x, y int) T) SquareMatrix[T] {
	matrix := newEmptyPixelMatrix[T](size)

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			matrix.Set(x, y, setupFunc(x, y))
		}
	}

	return matrix
}

func NewDotPixelMatrix(size int, inverted bool, min, max byte) SquareMatrix[byte] {
	points := generateDotPoints(size)
	matrix := newEmptyPixelMatrix[byte](size)

	sortPoints(points, inverted)
	equalDistributionValues(points, &matrix, min, max)

	return matrix
}

// sortPoints according to value.
func sortPoints(points []PointValue2D[int, float64], inverted bool) {
	if inverted {
		sort.Slice(points, func(i, j int) bool { return !(points[i].Value < points[j].Value) })
	} else {
		sort.Slice(points, func(i, j int) bool { return points[i].Value < points[j].Value })
	}
}

// dotMatrixEqualDistribution sets values in matrix one by one by step.
func equalDistributionValues(points []PointValue2D[int, float64], matrix *SquareMatrix[byte], minValue, maxValue byte) {
	// NOTE(ShaDream): example: min = 1, max = 3. points are 2. first point = 1, second point = 3. step is equal 2. (3-1)/(2-x) = 2, x is 1.
	stepValue := float64(maxValue-minValue) / float64(len(points)-1)

	baseValue := float64(minValue)

	for i := range points {
		newValue := math.Round(baseValue + (stepValue * float64(i)))

		point := points[i]
		setValue := min(int(newValue), int(maxValue))
		matrix.Set(point.X, point.Y, byte(setValue))
	}
}

// generateDotPoints generates array of points with it's weight. points are not ordered.
func generateDotPoints(size int) []PointValue2D[int, float64] {
	centralValue := float64(size-1) / 2
	centerPoint := Point2D[float64]{X: centralValue + .1, Y: centralValue + .15}

	result := make([]PointValue2D[int, float64], 0, size*size)

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			value := getDotPixelLevel(size, centerPoint, x, y)
			point := PointValue2D[int, float64]{
				Point2D: Point2D[int]{
					X: x,
					Y: y,
				},
				Value: value,
			}

			result = append(result, point)
		}
	}

	return result
}

// getDotPixelLevel calculate how far pixel from center point.
func getDotPixelLevel(size int, point Point2D[float64], pixelX, pixelY int) float64 {
	distance := point.Distance(float64(pixelX), float64(pixelY))
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

func ConcatMatricies[T any](size int, matricies []SquareMatrix[T]) SquareMatrix[T] {
	matrixSize := matricies[0].Size
	resultMatrixSize := matrixSize * size
	resultMatrix := newEmptyPixelMatrix[T](resultMatrixSize)

	for y := 0; y < resultMatrixSize; y++ {
		matrixY := y / matrixSize
		insideY := y % matrixSize
		for x := 0; x < resultMatrixSize; x++ {
			matrixX := x / matrixSize
			matrixIndex := matrixY*size + matrixX

			insideX := x % matrixSize

			value := matricies[matrixIndex].Get(insideX, insideY)
			resultMatrix.Set(x, y, value)
		}
	}

	return resultMatrix
}
