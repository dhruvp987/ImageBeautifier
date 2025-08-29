/*
 * Authors: Dhruv Patel and Ayush Sharma
 * File: pipeline.go
 * Description:
 *   The pipeline of transformations to run images through.
 */

package main

import "image"

/*
 * Pipe an image through a sequence of transformations.
 *
 * Parameters:
 *   img: The image to pipe
 *   tfms: The transformations to pipe the image through
 *
 * Returns: The new image after going through all of the transformations.
 */
func Pipe(img image.Image, tfms []func(image.Image) (image.Image, error)) (image.Image, error) {
    res := img
    var err error = nil
    for _, transform := range tfms {
	res, err = transform(res)
	if err != nil {
            return nil, err
	}
    }
    return res, nil
}
