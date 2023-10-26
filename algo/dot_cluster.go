package algo

type DotCluster struct {
	dots SquareMatrix[Dot]

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

func CreateThresholdMatrix(clusterSettings ClusterSettings) SquareMatrix[byte] {
	cluster := NewDotCluster(clusterSettings)

	matricies := []SquareMatrix[byte]{}
	for y := 0; y < cluster.clusterSize; y++ {
		for x := 0; x < cluster.clusterSize; x++ {
			matricies = append(matricies, cluster.dots.Get(x, y).PixelThresholdPoints)
		}
	}

	return ConcatMatricies[byte](cluster.clusterSize, matricies)
}
