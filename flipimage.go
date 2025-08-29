/*
 * Authors: Ayush Sharma and Dhruv Patel
 * File: flipimage.go
 * Description: Parallel vertical flip (upside-down) transformation for images.
 *              Compatible with image processing pipelines.
 */
package main

import (
        "image"
        "sync"
)

/*
 * Take a chunk of pixels from an image and flip them onto a new image.
 */
func setFlippedPixels(
	og image.Image,
	width int,
	height int,
	dst *image.RGBA,
	index int,
	maxIndex int,
) {
	for y := index; y < height; y += maxIndex {
		for x := 0; x < width; x++ {
			flippedy := height - y - 1
                        dst.Set(x, flippedy, og.At(x, y))
                }
        }
}

/*
 * Flip an image upside down.
 */
func FlipUpsideDownParallel(img image.Image) (image.Image, error) {
        bounds := img.Bounds()
        flipped := image.NewRGBA(bounds)
        width, height := bounds.Dx(), bounds.Dy()

	pool := GetGlobalWorkers()
        workers := pool.NumWorkers

	var wg sync.WaitGroup
	wg.Add(workers)

        // Process each row in parallel
        for worker := 0; worker < workers; worker++ {
                func(workerID int) {
			pool.Submit(func() {
                        	defer wg.Done()
				setFlippedPixels(
					img,
					width,
					height,
					flipped,
					workerID,
					workers,
				)
			})
                }(worker)
        }

        wg.Wait()
        return flipped, nil
}
