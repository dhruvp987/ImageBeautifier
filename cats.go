/*
 * Dhruv Patel and Ayush Sharma
 * File: cats.go
 * Description:
 *   Cry tears of joy by putting cats on your image.
 */

package main

import (
    "image"
    "image/draw"
    "image/png"
    "math/rand"
    "os"
    "sync"
)

const MAX_CAT_IMGS = 5

var CAT_IMG_PATHS = [2]string{
    "assets/cat1.png",
    "assets/cat2.png",
}

/*
 * A representation of a rectangle.
 */
type rect struct {
    X0 int;
    Y0 int;
    Width int;
    Height int;
}

/*
 * Check whether a rectangular area overlaps other rectangular areas.
 */
func overlaps(newArea rect, areas []rect) bool {
    numAreas := len(areas)
    for i := 0; i < numAreas; i++ {
	curArea := areas[i]
        if !(newArea.X0 + newArea.Width <= curArea.X0)  &&
	   !(newArea.X0 >= curArea.X0 + curArea.Width)  &&
	   !(newArea.Y0 + newArea.Height <= curArea.Y0) &&
	   !(newArea.Y0 >= curArea.Y0 + curArea.Height) {
            return true
	}
    }
    return false
}

/*
 * Load a PNG from disk into a usable image data structure.
 */
func decodePng(path string) (image.Image, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    img, err := png.Decode(file)
    if err != nil {
        return nil, err
    }

    return img, nil
}

/*
 * Get the requested image either by using whatever has already been loaded in
 * the store, or by loading a new one from disk and putting it in the store.
 */
func getImage(path string, store *map[string]image.Image) (image.Image, error) {
    derefStore := *store
    img, prs := derefStore[path]
    if !prs {
        newImg, err := decodePng(path)
	if err != nil {
            return nil, err
	}
        derefStore[path] = newImg
	return newImg, nil
    }
    return img, nil
}

/*
 * Draw an image on another image as a Goroutine.
 */
func drawImageRoutine(dst draw.Image, src image.Image, point image.Point, wg *sync.WaitGroup) {
    defer wg.Done()
    draw.Draw(dst, dst.Bounds().Add(point), src, image.Point{X: 0, Y: 0}, draw.Over)
}

/*
 * Randomly draw cat images on another image.
 */
func CatImages(img image.Image) (image.Image, error) {
    numCatImgs := len(CAT_IMG_PATHS)
    // Used to cache cat images.
    store := make(map[string]image.Image)
    // Used to track claimed areas and prevent overlapping
    claimedAreas := make([]rect, 0)

    baseBounds := img.Bounds()
    newImg := image.NewRGBA(baseBounds)

    // Draw the original image on the new image.
    draw.Draw(newImg, baseBounds, img, image.Point{X: 0, Y: 0}, draw.Src);

    numImgsDraw := (rand.Int() % MAX_CAT_IMGS) + 1

    pool := GetGlobalWorkers()

    var wg sync.WaitGroup
    for i := 0; i < numImgsDraw; i++ {
        catImg, err := getImage(CAT_IMG_PATHS[rand.Intn(numCatImgs)], &store)
	if err != nil {
	    // Just move on to the next iteration and try again.
            continue
	}
	catBounds := catImg.Bounds()

	randX := rand.Intn(baseBounds.Max.X)
	randY := rand.Intn(baseBounds.Max.Y)

        catArea := rect{randX, randY, catBounds.Dx(), catBounds.Dy()}
	if overlaps(catArea, claimedAreas) {
	    // Just move on to the next iteration and try again.
	    continue
	}
	claimedAreas = append(claimedAreas, catArea)

        wg.Add(1)
	func(catImage image.Image, x int, y int) {
	    pool.Submit(func() {
                defer wg.Done()
                draw.Draw(
		    newImg,
	            newImg.Bounds().Add(image.Point{x, y}),
		    catImage,
		    image.Point{X: 0, Y: 0},
		    draw.Over,
                )
	    })
        }(catImg, randX, randY)
    }
   
    wg.Wait()

    return newImg, nil
}
