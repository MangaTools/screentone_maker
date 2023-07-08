package algo

import (
	"math"
	"sort"
)

type PixelMatrix [][]byte

func (pm PixelMatrix) Get(x, y int) byte {
	return pm[x][y]
}

func (pm *PixelMatrix) SetOffset(value int8) {
	size := len(*pm)

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			newValue := float64(((*pm)[y][x])) - float64(value)
			clippedValue := math.Max(math.Min(maxPixelValue, newValue), centerValue)
			(*pm)[x][y] = byte(clippedValue)
		}
	}
}

func NewPixelMatrix(size int) PixelMatrix {
	return dotMatrixEqualDistribution(generateDotMatrix(size, centerValue, maxPixelValue), centerValue, maxPixelValue)
}

func generateDotMatrix(size int, min, max byte) [][]byte {
	centralValue := float64(size-1) / 2
	centerPoint := FPoint{X: centralValue + .1, Y: centralValue + .2}

	result := create2DMatrix[byte](size)

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			result[y][x] = getDotPixelLevel(size, centerPoint, x, y, min, max)
		}
	}

	return result
}

func dotMatrixEqualDistribution(matrix [][]byte, min, max byte) [][]byte {
	matrixSize := len(matrix)

	step := float64(max-min) / float64(matrixSize*matrixSize)

	order := make([]struct {
		x, y  int
		value byte
	}, 0, matrixSize*matrixSize)

	for x := 0; x < matrixSize; x++ {
		for y := 0; y < matrixSize; y++ {
			order = append(order, struct {
				x, y  int
				value byte
			}{
				x:     x,
				y:     y,
				value: matrix[x][y],
			})
		}
	}

	sort.Slice(order, func(i, j int) bool {
		return order[i].value < order[j].value
	})

	currentValue := float64(min)

	for i := 0; i < len(order); i++ {
		orderValue := order[i]
		matrix[orderValue.x][orderValue.y] = byte(currentValue)

		currentValue += step
	}

	return matrix
}

// func lerpDotMatrix(matrix [][]byte) [][]byte {
// 	var min byte = 255
// 	var max byte = 0

// 	matrixSize := len(matrix)

// 	for x := 0; x < matrixSize; x++ {
// 		for y := 0; y < matrixSize; y++ {
// 			value := matrix[x][y]
// 			if value > max {
// 				max = value
// 			}
// 			if value < min {
// 				min = value
// 			}
// 		}
// 	}

// 	max = max - min

// 	for x := 0; x < matrixSize; x++ {
// 		for y := 0; y < matrixSize; y++ {
// 			t := float64((matrix[x][y] - min)) / float64(max)

// 			newValue := byte(math.Round(lerp(centerValue, maxPixelValue, t)))
// 			matrix[x][y] = newValue
// 		}
// 	}

// 	return matrix
// }

func getDotPixelLevel(size int, point FPoint, pixelX, pixelY int, min, max byte) byte {
	distance := point.Distance(float64(pixelX), float64(pixelY))
	circleRadius := getMaxDistance(point, size)

	tLerp := distance / circleRadius

	byteFloat := lerp(float64(max), float64(min), tLerp)

	return byte(math.Round(byteFloat))
}

func getMaxDistance(point FPoint, size int) float64 {
	max := float64(0)

	maxSize := float64(size - 1)

	corners := []FPoint{{X: 0, Y: 0}, {X: 0, Y: maxSize}, {X: maxSize, Y: maxSize}, {X: maxSize, Y: 0}}
	for _, corner := range corners {
		value := point.Distance(corner.X, corner.Y)
		if max < value {
			max = value
		}
	}

	return max
}
