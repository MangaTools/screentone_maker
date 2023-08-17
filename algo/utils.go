package algo

// create 2d square matrix with variable size
func create2DMatrix[T any](size int) [][]T {
	result := make([][]T, size)

	for i := 0; i < size; i++ {
		result[i] = make([]T, size)
	}

	return result
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func toPoint2DValueArray[T any](slice2D [][]T) []PointValue2D[int, T] {
	size := len(slice2D)
	result := make([]PointValue2D[int, T], 0, size*size)

	for x := range slice2D {
		slice := slice2D[x]
		for y := range slice {
			result = append(result, PointValue2D[int, T]{
				Point2D: Point2D[int]{
					X: x,
					Y: y,
				},
				Value: slice[y],
			})
		}
	}

	return result
}
