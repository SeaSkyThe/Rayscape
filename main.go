package main

import (
	"fmt"
	"os"
)

func main() {
	// Image
	var image_width int = 256
	var image_height int = 256

	fmt.Printf("P3\n%d %d\n255\n", image_width, image_height)

	for j := 0; j < image_height; j++ {
        fmt.Fprintf(os.Stderr, "Scanlines remaining: %d \n", image_height - j)
        os.Stderr.Sync()
		for i := 0; i < image_width; i++ {
			r := float64(i) / float64(image_width-1)
			g := float64(j) / float64(image_height-1)
			b := 0.0

			// fmt.Printf("%d %d %d\n", r, g, b)

			base := 255.999
			ir := int64(base * r)
			ig := int64(base * g)
			ib := int64(base * b)

			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}
