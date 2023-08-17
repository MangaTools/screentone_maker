package executor

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
	"github.com/shadream/screentone_maker/algo"
)

var acceptedImageExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

func RunExecution(settings ExecutionSettings) error {
	clusterSettings := algo.ClusterSettings{
		Size: int(settings.ClusterSize),
		DotSettings: algo.DotSettings{
			MinValue: byte(settings.Black),
			MaxValue: byte(settings.White),
			Size:     int(settings.DotSize),
		},
	}

	cluster := algo.NewDotCluster(clusterSettings)
	imagePaths, err := getImagesPaths(settings.InputPath)
	if err != nil {
		return err
	}

	if settings.OutPath == "" {
		settings.OutPath = settings.InputPath
	}

	err = os.MkdirAll(settings.OutPath, 0o755)
	if err != nil {
		return fmt.Errorf("невозможно создать папку с результатом: %w", err)
	}

	imageChan := make(chan string)
	bar := progressbar.Default(int64(len(imagePaths)), "Загрузка и обработка...")

	doneWaitGroup := &sync.WaitGroup{}

	for i := 0; i < int(settings.Threads); i++ {
		doneWaitGroup.Add(1)
		go worker(bar, imageChan, doneWaitGroup, cluster, settings.OutPath)
	}

	for _, image := range imagePaths {
		imageChan <- image
		bar.Add(1)
	}

	bar.Close()
	close(imageChan)

	doneWaitGroup.Wait()

	fmt.Println("Готово!")

	return nil
}

func worker(bar *progressbar.ProgressBar, imagePathChan <-chan string, waitGroup *sync.WaitGroup, cluster *algo.DotCluster, outPath string) {
	for imagePath := range imagePathChan {
		executeOnFile(imagePath, outPath, cluster)
	}

	bar.Add(1)
	waitGroup.Done()
}

func executeOnFile(imagePath string, outPath string, cluster *algo.DotCluster) {
	parsedImage, err := readImageToStruct(imagePath)
	if err != nil {
		fmt.Println(err)

		return
	}
	imageSize := parsedImage.Bounds().Max

	resultImage := image.NewGray(image.Rect(0, 0, imageSize.X, imageSize.Y))
	grayColorModel := color.GrayModel

	for x := 0; x < imageSize.X; x++ {
		for y := 0; y < imageSize.Y; y++ {
			grayColor := grayColorModel.Convert(parsedImage.At(x, y)).(color.Gray)

			isBlack := cluster.IsPixelBlack(x, y, grayColor.Y)
			var resultColor byte
			if !isBlack {
				resultColor = 255
			}

			resultImage.SetGray(x, y, color.Gray{Y: resultColor})
		}
	}

	filename := strings.TrimSuffix(filepath.Base(imagePath), filepath.Ext(imagePath))
	filename += ".png"

	resultImagePath := filepath.Join(outPath, filename)

	resultFile, err := os.OpenFile(resultImagePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0o755)
	if err != nil {
		fmt.Printf("Невозможно создать файл %s: %s\n", filepath.Base(imagePath), err)

		return
	}
	defer resultFile.Close()

	err = png.Encode(resultFile, resultImage)
	if err != nil {
		fmt.Printf("Невозможно закодировать файл %s: %s\n", filepath.Base(imagePath), err)

		return
	}
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

func getImagesPaths(folder string) ([]string, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("Невозможно прочитать директорию: %w", err)
	}

	result := make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(folder, file.Name())
		fileExt := filepath.Ext(filePath)
		if _, ok := acceptedImageExtensions[fileExt]; !ok {
			continue
		}

		result = append(result, filePath)
	}

	return result, nil
}
