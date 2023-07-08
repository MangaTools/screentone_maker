package algo

import (
	"math"
	"math/rand"
	"sort"
)

type DotCluster struct {
	dots         [][]Dot
	clusterSize  int
	dotSize      int
	offsetMatrix [][]int8
}

func NewDotCluster(clusterSize int, dotSize int) *DotCluster {
	dots := create2DMatrix[Dot](clusterSize)

	offsetMatrix := BayerDitherMatrix(clusterSize)

	for x := 0; x < clusterSize; x++ {
		for y := 0; y < clusterSize; y++ {
			dots[x][y] = *NewDot(dotSize)
			dots[x][y].PixelThresholdPoints.SetOffset(offsetMatrix[x][y])
		}
	}

	return &DotCluster{
		dots:         dots,
		clusterSize:  clusterSize,
		dotSize:      dotSize,
		offsetMatrix: offsetMatrix,
	}
}

func (d *DotCluster) IsPixelBlack(x, y int, color byte) bool {
	clusterPixelSize := d.dotSize * d.clusterSize
	clusterPixelX := x % clusterPixelSize
	clusterPixelY := y % clusterPixelSize

	dotPixelX := clusterPixelX % d.dotSize
	dotPixelY := clusterPixelY % d.dotSize

	dotIndexX := clusterPixelX / d.dotSize
	dotIndexY := clusterPixelY / d.dotSize

	isDotAnti := (dotIndexX+dotIndexY)%2 == 0

	// newValue := int32(color) + int32(d.offsetMatrix[x%d.clusterSize][y%d.clusterSize])
	// clippedColor := byte(math.Max(math.Min(maxPixelValue, float64(newValue)), 0))

	return d.dots[dotIndexX][dotIndexY].IsPixelBlack(dotPixelX, dotPixelY, color, isDotAnti)
}

func (d *DotCluster) calculateDot() {
}

func createOffsetMatrixBySize(size int) [][]int8 {
	matrix := create2DMatrix[int8](size)
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			matrix[x][y] = int8((rand.Int() % 5) - 2)
		}
	}

	return matrix
}

func BayerDitherMatrix(clusterSize int) [][]int8 {
	current := bayer2x2
	for clusterSize > len(current) {
		current = expandBayerDitherMatrix(current)
	}

	result := orderMatrix(current)

	for x := 0; x < len(result); x++ {
		for y := 0; y < len(result); y++ {
			result[x][y] = result[x][y] % (int8(clusterSize) / 2)
		}
	}

	return result
}

func orderMatrix(matrix [][]float64) [][]int8 {
	result := create2DMatrix[int8](len(matrix))

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

	startValue := int8(-elements / 2)

	for i := 0; i < len(values); i++ {
		result[values[i].x][values[i].y] = startValue
		startValue++
	}

	return result
}

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

var matrixOffset = [][]byte{
	{1, 4, 5, 5, 4, 1},
	{4, 2, 4, 3, 2, 4},
	{5, 3, 1, 1, 4, 5},
	{5, 4, 1, 1, 3, 5},
	{4, 2, 3, 4, 2, 4},
	{1, 4, 5, 5, 4, 1},
}
