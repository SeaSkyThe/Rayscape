package main

import "math"

// Constants
var infinity float64 = math.Inf(1)
var pi float64 = 3.1415926535897932385

// Utility functions

func DegreesToRadians(degrees float64) float64{
    return degrees * pi / 180.0
}

