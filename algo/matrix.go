package algo

type SquareMatrix[T any] struct {
	Size   int
	Matrix []T
}

func (pm *SquareMatrix[T]) Get(x, y int) T {
	return pm.Matrix[pm.getIndex(x, y)]
}

func (pm *SquareMatrix[T]) set(x, y int, value T) {
	index := pm.getIndex(x, y)
	pm.Matrix[index] = value
}

func (pm *SquareMatrix[T]) getIndex(x, y int) int {
	return x + (y * pm.Size)
}

func newEmptyPixelMatrix[T any](size int) SquareMatrix[T] {
	matrix := SquareMatrix[T]{
		Size:   size,
		Matrix: make([]T, size*size),
	}

	return matrix
}

func concatMatricies[T any](size int, matricies []SquareMatrix[T]) SquareMatrix[T] {
	matrixSize := matricies[0].Size
	resultMatrixSize := matrixSize * size
	resultMatrix := SquareMatrix[T]{
		Size:   resultMatrixSize,
		Matrix: make([]T, 0, resultMatrixSize*resultMatrixSize),
	}

	// NOTE(amogilnikov): algorithm copies values in result matrix by lines in each matrix (copy entire Y line of small matrix instead one element at time).
	for y := 0; y < resultMatrixSize; y++ {
		matrixY := y / matrixSize
		insideMatrixIndex := (y % matrixSize) * matrixSize
		for matrixXIndex := 0; matrixXIndex < size; matrixXIndex++ {
			matrixIndex := (matrixY * size) + matrixXIndex

			values := matricies[matrixIndex].Matrix[insideMatrixIndex : insideMatrixIndex+matrixSize]

			resultMatrix.Matrix = append(resultMatrix.Matrix, values...)
		}
	}

	return resultMatrix
}
