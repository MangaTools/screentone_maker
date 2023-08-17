package algo

type DotCluster struct {
	dots Matrix[Dot]

	clusterSize int
	dotSize     int
}

func NewDotCluster(clusterSettings ClusterSettings) *DotCluster {
	// NOTE(ShaDream) min size is 2, but 0 turns off bayer dither across dots
	useBayer := true
	if clusterSettings.Size == 0 {
		clusterSettings.Size = 2
		useBayer = false
	}

	offsetMatrix := newEmptyPixelMatrix[uint](clusterSettings.Size)
	if useBayer {
		offsetMatrix = bayerOrderMatrix(clusterSettings.Size)
	}

	setupDotFunc := func(x, y int) Dot {
		maxValue := max(clusterSettings.DotSettings.MaxValue-byte(offsetMatrix.Get(x, y)), clusterSettings.DotSettings.MinValue)
		dot := *NewDot(clusterSettings.DotSettings.Size, (x+y)%2 != 0, clusterSettings.DotSettings.MinValue, maxValue)

		return dot
	}

	dotMatrix := NewMatrixWithSetupFunc[Dot](clusterSettings.Size, setupDotFunc)

	cluster := &DotCluster{
		dots:        dotMatrix,
		clusterSize: clusterSettings.Size,
		dotSize:     clusterSettings.DotSettings.Size,
	}

	return cluster
}

func (d *DotCluster) IsPixelBlack(x, y int, color byte) bool {
	clusterPixelSize := d.dotSize * d.clusterSize

	x = x % clusterPixelSize
	y = y % clusterPixelSize

	clusterPoint := Point2D[int]{
		X: x / d.dotSize,
		Y: y / d.dotSize,
	}

	dotPoint := Point2D[int]{
		X: x % d.dotSize,
		Y: y % d.dotSize,
	}

	dot := d.dots.Get(clusterPoint.X, clusterPoint.Y)

	return dot.IsPixelBlack(dotPoint.X, dotPoint.Y, color)
}

// // TODO (ShaDream): rewrite this shit
// func (d *DotCluster) ditherMixDots() {
// 	clusterDots := d.clusterSize * d.clusterSize
// 	dotsPixels := d.dotSize * d.dotSize

// 	clusterPixelsCount := dotsPixels * clusterDots
// 	clusterPixelsCountPerGroup := clusterPixelsCount / 2

// 	orderOfMatrixSteps := make([]pointWithValue[uint], 0, d.clusterSize*d.clusterSize)
// 	for x := 0; x < d.clusterSize; x++ {
// 		for y := 0; y < d.clusterSize; y++ {
// 			orderOfMatrixSteps = append(orderOfMatrixSteps, pointWithValue[uint]{
// 				Point: image.Point{
// 					X: x,
// 					Y: y,
// 				},
// 				Value: d.offsetMatrix[x][y],
// 			})
// 		}
// 	}

// 	sort.Slice(orderOfMatrixSteps, func(i, j int) bool {
// 		return orderOfMatrixSteps[i].Value < orderOfMatrixSteps[j].Value
// 	})

// 	dotMatrix := generateDotMatrix(d.dotSize)
// 	dotOrderPoints := make([]pointWithValue[float64], 0, dotsPixels)
// 	for x := 0; x < d.dotSize; x++ {
// 		for y := 0; y < d.dotSize; y++ {
// 			dotOrderPoints = append(dotOrderPoints, pointWithValue[float64]{
// 				Point: image.Point{
// 					X: x,
// 					Y: y,
// 				},
// 				Value: dotMatrix[x][y],
// 			})
// 		}
// 	}

// 	sort.Slice(dotOrderPoints, func(i, j int) bool {
// 		return dotOrderPoints[i].Value > dotOrderPoints[j].Value
// 	})

// 	currentDotPointMatrix := create2DMatrix[uint](d.clusterSize)

// 	blackStep := float64(maxPixelValue-centerValue) / float64(clusterPixelsCountPerGroup-1)
// 	blackBase := maxPixelValue
// 	whiteIndex := 0

// 	whiteStep := float64(centerValue-2) / float64(clusterPixelsCountPerGroup-1)
// 	whiteBase := 1
// 	blackIndex := 0

// 	for i := 0; i < clusterPixelsCount; i++ {
// 		clusterDot := orderOfMatrixSteps[i%clusterDots]
// 		isWhite := (clusterDot.X+clusterDot.Y)%2 != 0

// 		orderPoint := currentDotPointMatrix[clusterDot.X][clusterDot.Y]
// 		currentDotPointMatrix[clusterDot.X][clusterDot.Y]++
// 		dotPoint := dotOrderPoints[orderPoint]

// 		if isWhite {
// 			value := float64(whiteBase) + (float64(whiteStep) * float64(whiteIndex))
// 			whiteIndex++
// 			d.dots[clusterDot.X][clusterDot.Y].PixelThresholdPoints[dotPoint.X][dotPoint.Y] = byte(value)
// 		} else {
// 			value := float64(blackBase) - (float64(blackStep) * float64(blackIndex))
// 			blackIndex++
// 			d.dots[clusterDot.X][clusterDot.Y].PixelThresholdPoints[dotPoint.X][dotPoint.Y] = byte(value)
// 		}
// 	}
// }

// type pointWithValue[T any] struct {
// 	image.Point
// 	Value T
// }
