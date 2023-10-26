package algo

import (
	"math"
	"sort"
)

func bayerOrderMatrix(clusterSize int) SquareMatrix[uint] {
	current := bayer2x2
	for clusterSize > len(current) {
		current = expandBayerDitherMatrix(current)
	}

	return orderMatrix(current)
}

// use bayer values to count from smallest to biggest
func orderMatrix(matrix [][]float64) SquareMatrix[uint] {
	points := toPoint2DValueArray[float64](matrix)
	sort.Slice(points, func(i, j int) bool {
		return points[i].Value < points[j].Value
	})

	result := newEmptyPixelMatrix[uint](len(matrix))

	value := uint(0)
	for _, point := range points {
		result.Set(point.X, point.Y, value)
		value++
	}

	return result
}

// expandBayerDitherMatrix expands previous bayer matrix by 2
func expandBayerDitherMatrix(previous [][]float64) [][]float64 {
	newSize := len(previous) * 2
	result := create2DMatrix[float64](newSize)

	twoSizePowerTwo := math.Pow(float64(2)*float64(newSize), 2)

	for x := 0; x < newSize; x++ {
		for y := 0; y < newSize; y++ {
			subMatrixX := x % 2
			subMatrixY := y % 2

			subMatrixIndexX := x / 2
			subMatrixIndexY := y / 2

			previousValue := float64(previous[subMatrixIndexX][subMatrixIndexY])
			bayerValue := float64(bayer2x2[subMatrixX][subMatrixY])
			value := (twoSizePowerTwo*previousValue + bayerValue) / twoSizePowerTwo

			result[x][y] = value
		}
	}

	return result
}

var bayer2x2 = [][]float64{
	{0, 2},
	{3, 1},
}
