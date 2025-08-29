/*
Author: Ayush Sharma and Dhruv Patel
main.go
*/
package main

import (
        "flag"
        "fmt"
        "image"
        "image/draw"
        "image/jpeg"
        "image/png"
        "log"
        "os"
        "path/filepath"
	"math/rand"
	"runtime"
        "strings"
	"time"
)

// The global Worker Pool which all transformations will send
// concurrent tasks to if they have such tasks. This centralizes
// Goroutine management and simplifies transformations.
var pool WorkerPool

/*
 * Get the program's global Worker Pool to send concurrent tasks to.
 */
func GetGlobalWorkers() WorkerPool {
    return pool
}

func main() {
        // Seed rand for transformations that may use it.
	rand.Seed(time.Now().UnixNano())

        // Parse flags
        inputPath := flag.String("i", "", "input image path")
        outputPath := flag.String("o", "output.jpg", "output image path")
        commands := flag.String("c", "", "transformation commands (e.g. blur, grayscale, upsidedown, cats)")
        flag.Parse()

        // Validate input
        if *inputPath == "" {
                log.Fatal("Input path is required")
		os.Exit(1)
        }

        // Parse commands into parallel transformations
        tfms, err := TokensToTfms(strings.Split(*commands, ","))
        if err != nil {
                log.Fatalf("Command parsing failed: %v", err)
		os.Exit(1)
        }

        // Decode image
        img, err := decodeImage(*inputPath)
        if err != nil {
                log.Fatalf("Decode failed: %v", err)
		os.Exit(1)
        }

        // Convert to RGBA if needed (for parallel transforms)
        if _, ok := img.(*image.RGBA); !ok {
                rgba := image.NewRGBA(img.Bounds())
                draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
                img = rgba
        }

	// Start a global worker pool
	numCpu := runtime.NumCPU()
	pool = WpNew(numCpu, numCpu)
	pool.Start()

        // Run pipeline
        result, err := Pipe(img, tfms)
        if err != nil {
                log.Fatalf("Pipeline failed: %v", err)
		os.Exit(1)
        }

	// Wait for all tasks to complete, and then stop workers.
	pool.WaitAndStop()

        fmt.Println("Saving output image")

        // Save output
        if err := saveImage(result, *outputPath); err != nil {
                log.Fatalf("Save failed: %v", err)
		os.Exit(1)
        }

        fmt.Printf("Success! Saved to %s\n", *outputPath)
}

// decodeImage handles JPEG/PNG decoding
func decodeImage(path string) (image.Image, error) {
        file, err := os.Open(path)
        if err != nil {
                return nil, err
        }
        defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

// saveImage encodes based on file extension
func saveImage(img image.Image, path string) error {
        file, err := os.Create(path)
        if err != nil {
                return err
        }
        defer file.Close()

        switch strings.ToLower(filepath.Ext(path)) {
        case ".png":
                return png.Encode(file, img)
        default: // Default to JPEG
                return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
        }
}
