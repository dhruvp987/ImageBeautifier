/*
 * Authors: Dhruv Patel and Ayush Sharma
 * File: resize.go
 * Description:
 *   Resize an image.
 */

package main

import (
    "errors"
    "image"
    "image/color"
    "math"
    "sync"
)

/*
 * A four-dimensional vector to perform vector operations on.
 */
type vector4 struct {
    X float64;
    Y float64;
    Z float64;
    A float64;
}

/*
 * Convert a Color into a Vector containing RGBA values.
 */
func colorToVector4(c color.Color) vector4 {
    r, g, b, a := c.RGBA()
    return vector4{
	// The RGBA values are 16-bit, but we need only 8 bits,
	// so get rid of the least significant 8 bits.
        float64(r >> 8),
	float64(g >> 8),
	float64(b >> 8),
	float64(a >> 8),
    }
}

/*
 * Convert a Vector containing RGBA values into a color.RGBA struct.
 */
func (v vector4) toRGBA() color.RGBA {
    return color.RGBA{
        uint8(v.X),
	uint8(v.Y),
	uint8(v.Z),
	uint8(v.A),
    }
}

/*
 * Multiply a Vector by a scalar.
 */
func (v vector4) scalarMult(sc float64) vector4 {
    return vector4{v.X * sc, v.Y * sc, v.Z * sc, v.A * sc}
}

/*
 * Add a Vector with another Vector.
 */
func (v0 vector4) add(v1 vector4) vector4 {
    return vector4{
        v0.X + v1.X,
	v0.Y + v1.Y,
	v0.Z + v1.Z,
	v0.A + v1.A,
    }
}

/*
 * Set a pixel's color using bilinear interpolation.
 */
func pixelBilinear(
    newX int,
    newY int,
    ogImg image.Image,
    newImg *image.RGBA,
    factor float64,
) {
    scaledX := float64(newX) * (1.0 / factor)
    scaledY := float64(newY) * (1.0 / factor)

    x0 := math.Floor(scaledX)
    x1 := x0 + 1
    y0 := math.Floor(scaledY)
    y1 := y0 + 1

    x0Int := int(x0)
    x1Int := int(x1)
    y0Int := int(y0)
    y1Int := int(y1)

    // If the scaled values land exactly on a pixel, just return
    // that pixel's color.
    if x0 == x1 && y0 == y1 {
	    r, g, b, a := ogImg.At(x0Int, y0Int).RGBA()
	    newImg.Set(
		newX,
	        newY,
		color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)},
	    )
	    return
    }

    // Get the colors of the four neighboring pixels.
    leftTop := colorToVector4(ogImg.At(x0Int, y0Int))
    leftBot := colorToVector4(ogImg.At(x0Int, y1Int))
    rightTop := colorToVector4(ogImg.At(x1Int, y0Int))
    rightBot := colorToVector4(ogImg.At(x1Int, y1Int))

    leftXWeight := (x1 - scaledX) / (x1 - x0)
    rightXWeight := (scaledX - x0) / (x1 - x0)

    // Do two horizontal linear interpolations.
    topLinClr := leftTop.scalarMult(leftXWeight).add(rightTop.scalarMult(rightXWeight))
    botLinClr := leftBot.scalarMult(leftXWeight).add(rightBot.scalarMult(rightXWeight))

    topYWeight := (y1 - scaledY) / (y1 - y0)
    botYWeight := (scaledY - y0) / (y1 - y0)

    // Do a vertical linear interpolation with the two previous results.
    clr := topLinClr.scalarMult(topYWeight).add(botLinClr.scalarMult(botYWeight))

    newImg.Set(newX, newY, clr.toRGBA())
}

/*
 * Resample an image concurrently using bilinear interpolation.
 */
func concurBilinear(img image.Image, factor float64) (image.Image, error) {
    if factor <= 0 {
	return nil, errors.New("resize: factor must be greater than 0")
    }

    oldBounds := img.Bounds()
    newBounds := image.Rect(
        0,
	0,
	int(float64(oldBounds.Max.X)*factor),
	int(float64(oldBounds.Max.Y)*factor),
    )
    newWidth := newBounds.Dx()
    newHeight := newBounds.Dy()

    newImg := image.NewRGBA(newBounds)

    pool := GetGlobalWorkers()

    // Dedicate chunks of pixels to Goroutines to resample in parallel.
    numWorkers := pool.NumWorkers
    chunk := (newWidth * newHeight) / pool.NumWorkers
    if chunk < 1 {
	// If the pool's number of workers is greater than the number
	// of pixels, chunk would be set to 0, meaning that no work
	// would be done.
	numWorkers = newWidth * newHeight
	chunk = 1
    }

    var wg sync.WaitGroup
    wg.Add(numWorkers)

    for i := 0; i < numWorkers; i++ {
        func(index int) {
	    pool.Submit(func() {
	        defer wg.Done()

                start := index * chunk
	        end := (index + 1) * chunk
	        if i == numWorkers - 1 {
                    end = newWidth * newHeight - 1
	        }

	        for j := start; j < end; j++ {
                    newX := j % newWidth
		    newY := j / newWidth
                    pixelBilinear(newX, newY, img, newImg, factor)
	        }
            })
	}(i)
    }
    wg.Wait()

    return newImg, nil
}

/*
 * Get a function that will resize any image by the given factor.
 */
func ResizeT(factor float64) func (image.Image) (image.Image, error) {
    return func(src image.Image) (image.Image, error) {
        return concurBilinear(src, factor)
    }
}
