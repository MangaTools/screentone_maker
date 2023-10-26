package executor

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/shadream/screentone_maker/algo"
	"github.com/sourcegraph/conc/pool"
)

var acceptedImageExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

type Executor struct {
	thresholdMatrix algo.SquareMatrix[byte]
	settings        ExecutionSettings
}

func NewExecutor(settings ExecutionSettings) *Executor {
	clusterSettings := algo.ClusterSettings{
		Size: int(settings.ClusterSize),
		DotSettings: algo.DotSettings{
			MinValue: byte(settings.Black),
			MaxValue: byte(settings.White),
			Size:     int(settings.DotSize),
		},
	}

	thresholdMatrix := algo.CreateThresholdMatrix(clusterSettings)

	return &Executor{
		thresholdMatrix: thresholdMatrix,
		settings:        settings,
	}
}

func (e *Executor) ExecuteFolder(inputFolder, outputFolder string, recursive bool) error {
	inputFolder, err := filepath.Abs(inputFolder)
	if err != nil {
		return fmt.Errorf("Невозможно получить полный путь входной папки: %w", err)
	}

	images, err := getImagesPaths(inputFolder, acceptedImageExtensions, recursive)
	if err != nil {
		return fmt.Errorf("Невозможно получить пути до изображений: %w", err)
	}

	outputFolder, err = filepath.Abs(outputFolder)
	if err != nil {
		return fmt.Errorf("Невозможно получить полный путь выходной папки: %w", err)
	}

	items := len(images)

	wgPool := pool.New().WithMaxGoroutines(int(e.settings.Threads))
	bar := progressbar.Default(int64(items), "Загрузка и обработка...")

	for _, imageFile := range images {
		imageInputPath := filepath.Join(inputFolder, imageFile)
		imageOutputPath := filepath.Join(outputFolder, imageFile)
		err := os.MkdirAll(filepath.Dir(imageOutputPath), 0o755)
		if err != nil {
			return fmt.Errorf("Невозможно создать выходную папку %s: %w", filepath.Dir(imageOutputPath), err)
		}

		wgPool.Go(func() {
			err = e.ExecuteFile(imageInputPath, imageOutputPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Невозможно обработать файл %s: %s", imageInputPath, err.Error())
			}

			bar.Add(1)
		})
	}

	wgPool.Wait()

	return nil
}

func (e *Executor) ExecuteFile(inputFilePath, outputFilePath string) error {
	input, err := os.ReadFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("Невозможно прочитать файл %s: %w", inputFilePath, err)
	}

	resultData, err := e.Execute(bytes.NewBuffer(input))
	if err != nil {
		return fmt.Errorf("Попытка обработки файла %s: %w", inputFilePath, err)
	}

	if err := os.WriteFile(outputFilePath, resultData, 0o755); err != nil {
		return fmt.Errorf("Запись обработанного файла %s в %s: %w", inputFilePath, outputFilePath, err)
	}

	return nil
}

func (e *Executor) Execute(input io.Reader) ([]byte, error) {
	parsedImage, _, err := image.Decode(input)
	if err != nil {
		return nil, fmt.Errorf("Невозможно декодировать входные данные: %w", err)
	}

	imageSize := parsedImage.Bounds().Max

	resultImage := image.NewGray(image.Rect(0, 0, imageSize.X, imageSize.Y))
	grayColorModel := color.GrayModel

	for y := 0; y < imageSize.Y; y++ {
		lineIndex := y * imageSize.X

		matrixYBaseIndex := (y % e.thresholdMatrix.Size) * e.thresholdMatrix.Size
		for x := 0; x < imageSize.X; x++ {
			grayColor := grayColorModel.Convert(parsedImage.At(x, y)).(color.Gray)

			matrixX := x % e.thresholdMatrix.Size
			matrixIndex := matrixYBaseIndex + matrixX

			isBlack := e.thresholdMatrix.Matrix[matrixIndex] < grayColor.Y
			var resultColor byte
			if !isBlack {
				resultColor = 255
			}

			index := lineIndex + x
			resultImage.Pix[index] = resultColor
		}
	}

	result := bytes.Buffer{}
	if err := png.Encode(&result, resultImage); err != nil {
		return nil, fmt.Errorf("Невозможно закодировать итоговое изображение: %w", err)
	}

	return result.Bytes(), nil
}

func readImageToStruct(imagePath string) (image.Image, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("Невозможно открыть файл %s: %w", filepath.Base(imagePath), err)
	}
	defer imageFile.Close()

	parsedImage, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, fmt.Errorf("Невозможно декодировать файл %s: %w", filepath.Base(imagePath), err)
	}

	return parsedImage, nil
}

func getImagesPaths(folder string, extensions map[string]bool, recursive bool) ([]string, error) {
	dirs := []string{folder}

	result := make([]string, 0)

	folderCut := folder + string(filepath.Separator)

	for len(dirs) != 0 {
		currentFolder := dirs[0]

		files, err := os.ReadDir(currentFolder)
		if err != nil {
			return nil, fmt.Errorf("Невозможно прочитать директорию: %w", err)
		}

		dirs = dirs[1:]

		for _, file := range files {
			if file.IsDir() {
				dirs = append(dirs, filepath.Join(currentFolder, file.Name()))

				continue
			}

			filePath := filepath.Join(currentFolder, file.Name())
			fileExt := filepath.Ext(filePath)
			if _, ok := acceptedImageExtensions[fileExt]; !ok {
				continue
			}

			relativePath, _ := strings.CutPrefix(filePath, folderCut)
			result = append(result, relativePath)
		}
	}

	return result, nil
}
