package executor_test

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/shadream/screentone_maker/executor"
	"github.com/stretchr/testify/require"
)

const (
	maxX = 800
	maxY = 1200

	value = 128

	testDataPath     = "./testdata"
	rainbowImageName = "rainbow.png"
)

var (
	colorValue = color.Gray{
		Y: value,
	}

	defaultExecutorSettings = executor.ExecutionSettings{
		DotSize: 6,
		White:   255,
		Threads: uint(runtime.NumCPU()),
	}
)

func getRainbowData(t testing.TB) []byte {
	t.Helper()

	rainbowImagePath := filepath.Join(testDataPath, rainbowImageName)

	file, err := os.Open(rainbowImagePath)
	require.NoError(t, err)
	defer file.Close()

	data, err := io.ReadAll(file)
	require.NoError(t, err)

	return data
}

func BenchmarkExecutorProcessImage(b *testing.B) {
	rainbowData := getRainbowData(b)
	executor := executor.NewExecutor(defaultExecutorSettings)

	for i := 0; i < b.N; i++ {
		executor.Execute(bytes.NewReader(rainbowData))
	}
}

func BenchmarkSetByFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultImage := image.NewGray(image.Rect(0, 0, maxX, maxY))

		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				resultImage.Set(x, y, colorValue)
			}
		}
	}
}

func BenchmarkRawSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultImage := image.NewGray(image.Rect(0, 0, maxX, maxY))

		for y := 0; y < maxY; y++ {
			lineIndex := y * maxX
			for x := 0; x < maxX; x++ {
				index := lineIndex + x
				resultImage.Pix[index] = value
			}
		}
	}
}
