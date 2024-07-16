package main

import (
	"github.com/seaskythe/rayscape/camera"
	"github.com/seaskythe/rayscape/hittable"
	"github.com/seaskythe/rayscape/vector"
)

func main() {
	// World
	var world hittable.HittableList
	world.Add(hittable.Sphere{Center: vector.Point3{X: 0, Y: 0, Z: -1}, Radius: 0.5})
	world.Add(hittable.Sphere{Center: vector.Point3{X: 0, Y: -100.5, Z: -1}, Radius: 100})

	// Camera
	var cam camera.Camera
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 400

	cam.Render(world)
}
