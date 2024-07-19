package main

import (
	"math"

	"github.com/seaskythe/rayscape/camera"
	"github.com/seaskythe/rayscape/color"
	"github.com/seaskythe/rayscape/hittable"
	"github.com/seaskythe/rayscape/material"
	"github.com/seaskythe/rayscape/rtweekend"
	"github.com/seaskythe/rayscape/vector"
)

func build_world_scene_1() hittable.HittableList {
	// Materials
	material_ground := material.Lambertian{Albedo: color.Color3{X: 0.8, Y: 0.8, Z: 0.0}}
	material_center := material.Lambertian{Albedo: color.Color3{X: 0.1, Y: 0.2, Z: 0.5}}
	material_left := material.Dieletric{RefractionIndex: 1.50}
	material_bubble := material.Dieletric{RefractionIndex: 1.0 / 1.50}
	material_right := material.Metal{Albedo: color.Color3{X: 0.8, Y: 0.6, Z: 0.2}, Fuzz: 1}

	// World
	var world hittable.HittableList
	world.Add(hittable.Sphere{Center: vector.Point3{X: 0.0, Y: -100.5, Z: -1.0}, Radius: 100.0, Mat: material_ground})
	world.Add(hittable.Sphere{Center: vector.Point3{X: 0.0, Y: 0.0, Z: -1.2}, Radius: 0.5, Mat: material_center})
	world.Add(hittable.Sphere{Center: vector.Point3{X: -1.0, Y: 0.0, Z: -1.0}, Radius: 0.5, Mat: material_left})
	world.Add(hittable.Sphere{Center: vector.Point3{X: -1.0, Y: 0.0, Z: -1.0}, Radius: 0.4, Mat: material_bubble})
	world.Add(hittable.Sphere{Center: vector.Point3{X: 1.0, Y: 0.0, Z: -1.0}, Radius: 0.5, Mat: material_right})

	return world
}

func build_world_scene_test_fov() hittable.HittableList {
	R := math.Cos(rtweekend.PI / 4)
	// Materials
	material_left := material.Lambertian{Albedo: color.Color3{X: 0.0, Y: 0.0, Z: 1.0}}
	material_right := material.Lambertian{Albedo: color.Color3{X: 1.0, Y: 0.0, Z: 0.0}}

	// World
	var world hittable.HittableList
	world.Add(hittable.Sphere{Center: vector.Point3{X: -R, Y: 0.0, Z: -1.0}, Radius: R, Mat: material_left})
	world.Add(hittable.Sphere{Center: vector.Point3{X: R, Y: 0.0, Z: -1.0}, Radius: R, Mat: material_right})

	return world
}

func generateCam() camera.Camera {
	var cam camera.Camera
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 400
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 50

	cam.Vfov = 20
	cam.Lookfrom = vector.Point3{X: -2, Y: 2, Z: 1}
	cam.Lookat = vector.Point3{X: 0, Y: 0, Z: -1}
    cam.Vup = vector.Vec3{X: 0, Y: 1, Z: 0}

	return cam
}

func main() {
	world := build_world_scene_1()
	// world := build_world_scene_test_fov()

	cam := generateCam()

	cam.Render(world)
}
