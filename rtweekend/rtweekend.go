package rtweekend

import (
	"math"
	"math/rand"
)

// Constants
var INFINITY float64 = math.Inf(1)
var PI float64 = 3.1415926535897932385

// Utility functions

func DegreesToRadians(degrees float64) float64 {
	return degrees * PI / 180.0
}

func RandomDouble() float64 {
	// Returns random real in [0, 1)
	return rand.Float64()
}

func RandomDoubleInInterval(min float64, max float64) float64 {
	return min + (max-min)*RandomDouble()
}
