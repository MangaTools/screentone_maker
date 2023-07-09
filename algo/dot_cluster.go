package algo

import (
	"image"
	"sort"
)

type DotCluster struct {
	dots        [][]Dot
	clusterSize int
	dotSize     int
	// offsetMatrix is the order matrix in which dots pixel matrix should enables
	offsetMatrix [][]uint
}

func NewDotCluster(clusterSize int, dotSize int) *DotCluster {
	useMix := true
	if clusterSize == 0 {
		clusterSize = 2
		useMix = false
	}

	dots := create2DMatrix[Dot](clusterSize)
	offsetMatrix := bayerOrderMatrix(clusterSize)

	for x := 0; x < clusterSize; x++ {
		for y := 0; y < clusterSize; y++ {
			dots[x][y] = *NewDot(dotSize, (x+y)%2 != 0)
		}
	}

	cluster := &DotCluster{
		dots:         dots,
		clusterSize:  clusterSize,
		dotSize:      dotSize,
		offsetMatrix: offsetMatrix,
	}

	if useMix {
		cluster.ditherMixDots()
	}

	return cluster
}

func (d *DotCluster) IsPixelBlack(x, y int, color byte) bool {
	clusterPixelSize := d.dotSize * d.clusterSize
	clusterPixelX := x % clusterPixelSize
	clusterPixelY := y % clusterPixelSize

	dotPixelX := clusterPixelX % d.dotSize
	dotPixelY := clusterPixelY % d.dotSize

	dotIndexX := clusterPixelX / d.dotSize
	dotIndexY := clusterPixelY / d.dotSize

	return d.dots[dotIndexX][dotIndexY].IsPixelBlack(dotPixelX, dotPixelY, color)
}

// TODO (ShaDream): rewrite this shit
func (d *DotCluster) ditherMixDots() {
	clusterDots := d.clusterSize * d.clusterSize
	dotsPixels := d.dotSize * d.dotSize

	clusterPixelsCount := dotsPixels * clusterDots
	clusterPixelsCountPerGroup := clusterPixelsCount / 2

	orderOfMatrixSteps := make([]pointWithValue[uint], 0, d.clusterSize*d.clusterSize)
	for x := 0; x < d.clusterSize; x++ {
		for y := 0; y < d.clusterSize; y++ {
			orderOfMatrixSteps = append(orderOfMatrixSteps, pointWithValue[uint]{
				Point: image.Point{
					X: x,
					Y: y,
				},
				Value: d.offsetMatrix[x][y],
			})
		}
	}

	sort.Slice(orderOfMatrixSteps, func(i, j int) bool {
		return orderOfMatrixSteps[i].Value < orderOfMatrixSteps[j].Value
	})

	dotMatrix := generateDotMatrix(d.dotSize)
	dotOrderPoints := make([]pointWithValue[float64], 0, dotsPixels)
	for x := 0; x < d.dotSize; x++ {
		for y := 0; y < d.dotSize; y++ {
			dotOrderPoints = append(dotOrderPoints, pointWithValue[float64]{
				Point: image.Point{
					X: x,
					Y: y,
				},
				Value: dotMatrix[x][y],
			})
		}
	}

	sort.Slice(dotOrderPoints, func(i, j int) bool {
		return dotOrderPoints[i].Value > dotOrderPoints[j].Value
	})

	currentDotPointMatrix := create2DMatrix[uint](d.clusterSize)

	blackStep := float64(maxPixelValue-centerValue) / float64(clusterPixelsCountPerGroup-1)
	blackBase := maxPixelValue
	whiteIndex := 0

	whiteStep := float64(centerValue-2) / float64(clusterPixelsCountPerGroup-1)
	whiteBase := 1
	blackIndex := 0

	for i := 0; i < clusterPixelsCount; i++ {
		clusterDot := orderOfMatrixSteps[i%clusterDots]
		isWhite := (clusterDot.X+clusterDot.Y)%2 != 0

		orderPoint := currentDotPointMatrix[clusterDot.X][clusterDot.Y]
		currentDotPointMatrix[clusterDot.X][clusterDot.Y]++
		dotPoint := dotOrderPoints[orderPoint]

		if isWhite {
			value := float64(whiteBase) + (float64(whiteStep) * float64(whiteIndex))
			whiteIndex++
			d.dots[clusterDot.X][clusterDot.Y].PixelThresholdPoints[dotPoint.X][dotPoint.Y] = byte(value)
		} else {
			value := float64(blackBase) - (float64(blackStep) * float64(blackIndex))
			blackIndex++
			d.dots[clusterDot.X][clusterDot.Y].PixelThresholdPoints[dotPoint.X][dotPoint.Y] = byte(value)
		}
	}
}

type pointWithValue[T any] struct {
	image.Point
	Value T
}
