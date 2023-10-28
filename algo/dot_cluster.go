package algo

const clusterSize = 2

func generateDotCluster(dotSettings DotSettings) []SquareMatrix[byte] {
	regularDot := newDotBuilder(dotSettings.Size, false)
	invertedDot := newDotBuilder(dotSettings.Size, true)

	dotMatricies := make([]SquareMatrix[byte], 0, dotSettings.Size*dotSettings.Size)

	maxSettings := dotSettings.MaxValue
	minSettings := dotSettings.MinValue

	centralPoint := min(minSettings+(maxSettings-minSettings)/2+1, maxSettings)

	for y := 0; y < clusterSize; y++ {
		for x := 0; x < clusterSize; x++ {
			isInverted := (x+y)%2 != 0
			if isInverted {
				maxValue := max((centralPoint - 1), minSettings)
				minValue := minSettings

				dotMatricies = append(dotMatricies, invertedDot.generateMatrix(minValue, maxValue))
			} else {
				maxValue := maxSettings
				minValue := centralPoint

				dotMatricies = append(dotMatricies, regularDot.generateMatrix(minValue, maxValue))
			}

		}
	}

	return dotMatricies
}

func CreateThresholdMatrix(dotSettings DotSettings) SquareMatrix[byte] {
	cluster := generateDotCluster(dotSettings)

	return concatMatricies[byte](clusterSize, cluster)
}
