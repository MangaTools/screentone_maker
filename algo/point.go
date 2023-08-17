package algo

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Point2D[T constraints.Integer | constraints.Float] struct {
	X, Y T
}

// Distance calculate distance between two points
func (p *Point2D[T]) Distance(x, y T) float64 {
	distanceX := float64(x) - float64(p.X)
	distanceY := float64(y) - float64(p.Y)

	pow2 := math.Pow(distanceX, 2) + math.Pow(distanceY, 2)

	return math.Sqrt(pow2)
}

type PointValue2D[T constraints.Integer | constraints.Float, V any] struct {
	Point2D[T]
	Value V
}
