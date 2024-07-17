package color

import (
	"fmt"
	"io"
	"math"

	"github.com/seaskythe/rayscape/interval"
	"github.com/seaskythe/rayscape/vector"
)

type (
	Color3 = vector.Vec3
)

func LinearToGamma(linear_component float64) float64 {
    if linear_component > 0 {
        return math.Sqrt(linear_component)
    }

    return 0
}
func WriteColor(out io.Writer, pixel_color Color3) {
	// Transform from [0, 1] to [0,255]
	r := pixel_color.X
	g := pixel_color.Y
	b := pixel_color.Z

    // Apply a linear to gamma transform for gamma 2
    r = LinearToGamma(r)
    g = LinearToGamma(g)
    b = LinearToGamma(b)

	// Translate the [0, 1] component values to the range [0, 255]
	var intensity interval.Interval = interval.Interval{Min: 0.000, Max: 0.999}
    rbyte := int(256 * intensity.Clamp(r))
    gbyte := int(256 * intensity.Clamp(g))
    bbyte := int(256 * intensity.Clamp(b))

	_, err := fmt.Fprintf(out, "%d %d %d\n", rbyte, gbyte, bbyte)
	if err != nil {
		panic(err)
	}
}
