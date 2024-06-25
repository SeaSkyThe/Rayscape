package main

import (
	"fmt"
	"os"
)

func main() {
	// Image
	var image_width int = 256
	var image_height int = 256

    file := CreatePPMImage(image_width, image_height)
    defer file.Close()

	for j := 0; j < image_height; j++ {
		fmt.Fprintf(os.Stderr, "Scanlines remaining: %d \n", image_height-j)
		os.Stderr.Sync()
		for i := 0; i < image_width; i++ {
			pixel_color := Color3{float64(i) / float64(image_width-1), float64(j) / float64(image_height-1), 0.0}
			WriteColor(file, pixel_color)
		}
	}
}
