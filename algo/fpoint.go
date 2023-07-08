package algo

import "math"

type FPoint struct {
	X float64
	Y float64
}

func (p *FPoint) Distance(x, y float64) float64 {
	return math.Sqrt(math.Pow(x-p.X, 2) + math.Pow(y-p.Y, 2))
}
