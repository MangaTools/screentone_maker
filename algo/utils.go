package algo

func create2DMatrix[T any](size int) [][]T {
	result := make([][]T, size)

	for i := 0; i < size; i++ {
		result[i] = make([]T, size)
	}

	return result
}
