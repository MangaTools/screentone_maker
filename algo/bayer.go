package algo

import (
	"math"
	"sort"
)

func bayerOrderMatrix(clusterSize int) [][]uint {
	current := bayer2x2
	for clusterSize > len(current) {
		current = expandBayerDitherMatrix(current)
	}

	result := orderMatrix(current)

	return result
}

// use bayer values to count from smallest to biggest
func orderMatrix(matrix [][]float64) [][]uint {
	result := create2DMatrix[uint](len(matrix))

	elements := len(matrix) * len(matrix)

	values := make([]struct {
		x, y  int
		value float64
	}, 0, elements)

	for x := 0; x < len(matrix); x++ {
		for y := 0; y < len(matrix); y++ {
			values = append(values, struct {
				x     int
				y     int
				value float64
			}{
				x:     x,
				y:     y,
				value: matrix[x][y],
			})
		}
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i].value < values[j].value
	})

	startValue := uint(0)

	for i := 0; i < len(values); i++ {
		result[values[i].x][values[i].y] = startValue
		startValue++
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
