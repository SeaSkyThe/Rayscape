package rtweekend

import "math"

// Constants
var INFINITY float64 = math.Inf(1)
var PI float64 = 3.1415926535897932385

// Utility functions

func DegreesToRadians(degrees float64) float64{
    return degrees * PI / 180.0
}

