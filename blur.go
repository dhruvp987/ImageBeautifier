/*
 * Authors: Ayush Sharma and Dhruv Patel
 * File: blur.go
 * Description: Parallel convolution-based blur with pipeline compatibility
 */

package main

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// BlurParallel applies Gaussian blur using all available CPU cores
// Returns:
//   - image.Image: The blurred image (as *image.RGBA)
//   - error: Always nil in current implementation (maintained for pipeline compatibility)
func BlurParallel(img image.Image) (image.Image, error) {
	// Define 3x3 Gaussian kernel
	kernel := [][]float64{
		{1 / 16.0, 2 / 16.0, 1 / 16.0},
		{2 / 16.0, 4 / 16.0, 2 / 16.0},
		{1 / 16.0, 2 / 16.0, 1 / 16.0},
	}

	bounds := img.Bounds()
	blurred := image.NewRGBA(bounds)

	pool := GetGlobalWorkers()
	workers := pool.NumWorkers

	var wg sync.WaitGroup
	wg.Add(workers)

	// Process image in parallel strips
	for i := 0; i < workers; i++ {
		func(workerID int) {
			pool.Submit(func() {
				defer wg.Done()
				applyConvolutionWorker(
					img,
					blurred,
					kernel,
					bounds,
					workerID,
					workers,
				)
			})
		}(i)
	}

	wg.Wait()
	return blurred, nil
}

// Worker function for parallel convolution
func applyConvolutionWorker(
	src image.Image,
	dst *image.RGBA,
	kernel [][]float64,
	bounds image.Rectangle,
	workerID int,
	workers int,
) {
	kernelSize := len(kernel)
	radius := kernelSize / 2

	// Calculate this worker's strip of the image
	stripHeight := bounds.Dy() / workers
	yStart := bounds.Min.Y + workerID*stripHeight
	yEnd := yStart + stripHeight
	if workerID == workers-1 { // Last worker gets remainder rows
		yEnd = bounds.Max.Y
	}

	for y := yStart; y < yEnd; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a float64

			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					px := clamp(x+kx-radius, bounds.Min.X, bounds.Max.X-1)
					py := clamp(y+ky-radius, bounds.Min.Y, bounds.Max.Y-1)

					pixel := src.At(px, py)
					pr, pg, pb, pa := pixel.RGBA()
					weight := kernel[ky][kx]

					r += float64(pr>>8) * weight
					g += float64(pg>>8) * weight
					b += float64(pb>>8) * weight
					a += float64(pa>>8) * weight
				}
			}

			dst.Set(x, y, color.RGBA{
				R: uint8(math.Min(255, math.Max(0, r))),
				G: uint8(math.Min(255, math.Max(0, g))),
				B: uint8(math.Min(255, math.Max(0, b))),
				A: uint8(math.Min(255, math.Max(0, a))),
			})
		}
	}
}

// clamp ensures pixel coordinates stay within bounds
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
