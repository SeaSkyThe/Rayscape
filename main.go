package main

import (
	"github.com/seaskythe/rayscape/camera"
	"github.com/seaskythe/rayscape/color"
	"github.com/seaskythe/rayscape/hittable"
	"github.com/seaskythe/rayscape/material"
	"github.com/seaskythe/rayscape/vector"
)

func main() {
	// Materials
	material_ground := material.Lambertian{Albedo: color.Color3{X: 0.8, Y: 0.8, Z: 0.0}}
	material_center := material.Lambertian{Albedo: color.Color3{X: 0.1, Y: 0.2, Z: 0.5}}
	material_left := material.Metal{Albedo: color.Color3{X: 0.8, Y: 0.8, Z: 0.8}, Fuzz: 0.3}
    material_right := material.Metal{Albedo: color.Color3{X: 0.8, Y: 0.6, Z: 0.2}, Fuzz: 0.8}

	// World
	var world hittable.HittableList
	world.Add(hittable.Sphere{Center: vector.Point3{X: 0.0, Y: -100.5, Z: -1.0}, Radius: 100.0, Mat: material_ground})
	world.Add(hittable.Sphere{Center: vector.Point3{X: 0.0, Y: 0.0, Z: -1.2}, Radius: 0.5, Mat: material_center})
	world.Add(hittable.Sphere{Center: vector.Point3{X: -1.0, Y: 0.0, Z: -1.0}, Radius: 0.5, Mat: material_left})
	world.Add(hittable.Sphere{Center: vector.Point3{X: 1.5, Y: 0.0, Z: -2.0}, Radius: 0.5, Mat: material_right})

	// Camera
	var cam camera.Camera
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 400
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 50

	cam.Render(world)
}
