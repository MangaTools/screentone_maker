package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	flag "github.com/spf13/pflag"

	"github.com/shadream/screentone_maker/algo"
)

var acceptedImageExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

var (
	dotSize     = 5
	clusterSize = 6
	folderPath  = ""
	outFolder   = ""
)

func init() {
	flag.IntVarP(&dotSize, "dot_size", "d", 5, "max dot size in pixels")
	flag.IntVarP(&clusterSize, "cluster_size", "c", 6, "max dots in one cluster(matrix size)")
	flag.StringVarP(&folderPath, "input_folder", "i", "", "path to folder with images to use")
	flag.StringVarP(&outFolder, "out_folder", "o", "", "path to folder with result images")

	flag.Parse()
}

func main() {
	if folderPath == "" {
		log.Fatal("use -f 'folder path' to use program")
	}

	cluster := algo.NewDotCluster(clusterSize, dotSize)
	imagePaths := getImagesPaths(folderPath)

	err := os.MkdirAll(outFolder, 0o755)
	if err != nil {
		log.Fatalf("can not create result dir: %s", err)
	}

	imageChan := make(chan string)

	doneWaitGroup := &sync.WaitGroup{}

	for i := 0; i < runtime.NumCPU(); i++ {
		doneWaitGroup.Add(1)
		go worker(imageChan, doneWaitGroup, cluster)
	}

	for _, image := range imagePaths {
		imageChan <- image
	}

	close(imageChan)

	doneWaitGroup.Wait()

	fmt.Println("Done!")
}

func worker(imagePathChan <-chan string, waitGroup *sync.WaitGroup, cluster *algo.DotCluster) {
	for imagePath := range imagePathChan {
		executeOnFile(imagePath, cluster)
	}

	waitGroup.Done()
}

func executeOnFile(imagePath string, cluster *algo.DotCluster) {
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

	resultImagePath := filepath.Join(outFolder, filename)

	resultFile, err := os.OpenFile(resultImagePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0o755)
	if err != nil {
		fmt.Printf("can not create file for %s: %s\n", filepath.Base(imagePath), err)

		return
	}
	defer resultFile.Close()

	err = png.Encode(resultFile, resultImage)
	if err != nil {
		fmt.Printf("can not encode result file for %s: %s\n", filepath.Base(imagePath), err)

		return
	}
}

func readImageToStruct(imagePath string) (image.Image, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("can not open file %s: %w", filepath.Base(imagePath), err)
	}
	defer imageFile.Close()

	parsedImage, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, fmt.Errorf("can not decode file %s: %w", filepath.Base(imagePath), err)
	}

	return parsedImage, nil
}

func getImagesPaths(folder string) []string {
	files, err := os.ReadDir(folder)
	if err != nil {
		log.Fatal(fmt.Sprintf("can not read dir: %s", err))
	}

	result := make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(folderPath, file.Name())
		fileExt := filepath.Ext(filePath)
		if _, ok := acceptedImageExtensions[fileExt]; !ok {
			continue
		}

		result = append(result, filePath)
	}

	return result
}
