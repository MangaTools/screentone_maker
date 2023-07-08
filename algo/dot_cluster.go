package algo

type DotCluster struct {
	dots         [][]Dot
	clusterSize  int
	dotSize      int
	offsetMatrix [][]int8
}

func NewDotCluster(clusterSize int, dotSize int) *DotCluster {
	dots := create2DMatrix[Dot](clusterSize)

	offsetMatrix := bayerDitherMatrix(clusterSize)

	for x := 0; x < clusterSize; x++ {
		for y := 0; y < clusterSize; y++ {
			dots[x][y] = *NewDot(dotSize)
			// Use offset to change dot actuation point in different places
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

	return d.dots[dotIndexX][dotIndexY].IsPixelBlack(dotPixelX, dotPixelY, color, isDotAnti)
}
