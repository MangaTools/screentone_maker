package executor_test

import (
	"image"
	"image/color"
	"runtime"
	"sync"
	"testing"
)

const (
	maxX = 800
	maxY = 1200

	value = 128
)

var colorValue = color.Gray{
	Y: value,
}

func BenchmarkOneThread(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultImage := image.NewGray(image.Rect(0, 0, maxX, maxY))

		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				resultImage.Set(x, y, colorValue)
			}
		}
	}
}

func BenchmarkMultipleThreads(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultImage := image.NewGray(image.Rect(0, 0, maxX, maxY))

		rowChan := make(chan int, runtime.NumCPU())
		wg := sync.WaitGroup{}

		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for y := range rowChan {
					for x := 0; x < maxX; x++ {
						resultImage.Set(x, y, colorValue)
					}
					wg.Done()
				}
			}()
		}

		for y := 0; y < maxY; y++ {
			wg.Add(1)
			rowChan <- y
		}

		close(rowChan)
		wg.Wait()
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

func BenchmarkRawSetMultipleThreads(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultImage := image.NewGray(image.Rect(0, 0, maxX, maxY))

		rowChan := make(chan int, runtime.NumCPU())
		wg := sync.WaitGroup{}

		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for y := range rowChan {
					lineIndex := y * maxX
					for x := 0; x < maxX; x++ {
						index := lineIndex + x
						resultImage.Pix[index] = value
					}

					wg.Done()
				}
			}()
		}

		for i := 0; i < b.N; i++ {
			for y := 0; y < maxY; y++ {
				wg.Add(1)
				rowChan <- y
			}
		}

		close(rowChan)
		wg.Wait()
	}
}
