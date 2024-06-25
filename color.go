package main

import (
	"fmt"
	"io"
)

type (
    Color3 = Vec3
)

func WriteColor(out io.Writer, c Color3) {
    // Transform from [0, 1] to [0,255]
    r:= int(255.999 * c.X)
    g:= int(255.999 * c.Y)
    b:= int(255.999 * c.Z)

    _, err := fmt.Fprintf(out, "%d %d %d\n", r, g, b)
    if err != nil {
        panic(err)
    }
}
