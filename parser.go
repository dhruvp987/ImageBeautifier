/*
 * Authors: Dhruv Patel and Ayush Sharma
 * File: parser.go
 * Description:
 *   Source for a command line argument parser.
 */

package main

import (
    "errors"
    "fmt"
    "image"
    "strconv"
)

/*
 * Convert tokens into transformation functions.
 *
 * Parameter:
 *   tokens: The tokens to convert
 *
 * Returns: An array of transformation functions or an error.
 */
func TokensToTfms(tokens []string) ([]func(img image.Image) (image.Image, error), error) {
    lenTkns := len(tokens)
    tfms := make([]func(image.Image) (image.Image, error), 0)

    for i := 0; i < lenTkns; i++ {
        switch tokens[i] {
	    case "blur":
                tfms = append(tfms, BlurParallel)
	    case "grayscale":
	        tfms = append(tfms, GrayscaleParallel)
            case "upsidedown":
		tfms = append(tfms, FlipUpsideDownParallel)
            case "resize":
		scalar, err := strconv.ParseFloat(tokens[i + 1], 64)
		if err != nil {
                    return nil, err
		}
		i++
		tfms = append(tfms, ResizeT(scalar))
	    case "cats":
	        tfms = append(tfms, CatImages)
	    default:
		return nil, errors.New(fmt.Sprintf("%s is not a valid option", tokens[i]))
	}
    }

    return tfms, nil
}
