/*
Author: Ayush Sharma and Dhruv Patel
Grayscaling image
*/
package main

import (
        "image"
        "image/color"
        // "image/png"
        // "os"
        "sync"
)

/*
 * Set a chunk of pixels from an image grayscale and save them
 * onto a new image.
 */
func setGrayPixels(
	src image.Image,
	bounds image.Rectangle,
	dst *image.Gray,
	index int,
	maxIndex int,
) {
	for y := bounds.Min.Y + index; y < bounds.Max.Y; y += maxIndex {
        	for x := bounds.Min.X; x < bounds.Max.X; x++ {
                	r, g, b, _ := src.At(x, y).RGBA()
                        gray := uint8((r>>8 + g>>8 + b>>8) / 3)
                        dst.Set(x, y, color.Gray{Y: gray})
                }
        }
}

// GrayscaleParallel converts an image to grayscale using Goroutines.
// Returns (image.Image, error) to fit pipeline-style processing.
func GrayscaleParallel(img image.Image) (image.Image, error) {
        bounds := img.Bounds()
        grayImg := image.NewGray(bounds)

	pool := GetGlobalWorkers()
        numWorkers := pool.NumWorkers

        var wg sync.WaitGroup
	wg.Add(numWorkers)

        for i := 0; i < numWorkers; i++ {
                func(workerID int) {
			pool.Submit(func() {
                        	defer wg.Done()
				setGrayPixels(
					img,
					bounds,
					grayImg,
					workerID,
					numWorkers,
				)
			})
                }(i)
        }
        wg.Wait() // Wait for all Goroutines to finish
        return grayImg, nil // No error in this case, but signature matches pipeline
}

/*
func main() {
        // Example pipeline usage:
        img, err := loadImage("input.png")
        if err != nil {
                panic(err)
        }

        //Process IMage
        grayImg, err := grayscaleParallel(img)
        if err != nil {
                panic(err)
        }

        // Save output
        err = saveImage(grayImg, "output_gray.png")
        if err != nil {
                panic(err)
        }
}

// Helper function to load an image (returns pipeline signature)
func loadImage(filename string) (image.Image, error) {
        file, err := os.Open(filename)
        if err != nil {
                return nil, err
        }
        defer file.Close()
        return png.Decode(file)
}

// Helper function to save an image (returns error)
func saveImage(img image.Image, filename string) error {
        file, err := os.Create(filename)
        if err != nil {
                return err
        }
        defer file.Close()
        return png.Encode(file, img)
}
*/
