package main

import (
	"fmt"
	"os"
)

func CreatePPMImage(width, height int) *os.File {
	file, err := os.OpenFile("image.ppm", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(file, "P3\n%d %d\n255\n", width, height)

	return file
}
