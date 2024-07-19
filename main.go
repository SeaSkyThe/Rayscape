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

func buildWorldScene1() hittable.HittableList {
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

func buildWorldSceneBookCover() hittable.HittableList {
	var world hittable.HittableList

	ground_material := material.Lambertian{Albedo: color.Color3{X: 0.5, Y: 0.5, Z: 0.5}}
	world.Add(hittable.Sphere{Center: vector.Point3{X: 0, Y: -1000, Z: 0}, Radius: 1000, Mat: ground_material})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			choose_material := rtweekend.RandomDouble()
			center := vector.Point3{X: float64(a) + 0.9*rtweekend.RandomDouble(), Y: 0.2, Z: float64(b) + 0.9*rtweekend.RandomDouble()}

			if (vector.Subtract(center, vector.Point3{X: 4, Y: 0.2, Z: 0}).Length() > 0.9) {
				var sphere_material material.Material
				if choose_material < 0.8 {
					//Diffuse
					albedo := vector.Multiply(color.Random(0, 1), color.Random(0, 1))
					sphere_material = material.Lambertian{Albedo: albedo}
					world.Add(hittable.Sphere{Center: center, Radius: 0.2, Mat: sphere_material})
				} else if choose_material < 0.95 {
					//Metal
					albedo := vector.Multiply(color.Random(0.5, 1), color.Random(0.5, 1))
					fuzz := rtweekend.RandomDoubleInInterval(0, 0.5)
					sphere_material = material.Metal{Albedo: albedo, Fuzz: fuzz}
					world.Add(hittable.Sphere{Center: center, Radius: 0.2, Mat: sphere_material})
				} else {
					// Glass
					sphere_material = material.Dieletric{RefractionIndex: 1.5}
					world.Add(hittable.Sphere{Center: center, Radius: 0.2, Mat: sphere_material})
				}
			}
		}
	}
	material1 := material.Dieletric{RefractionIndex: 1.5}
	world.Add(hittable.Sphere{Center: vector.Point3{X: 0, Y: 1, Z: 0}, Radius: 1.0, Mat: material1})

	material2 := material.Lambertian{Albedo: color.Color3{X: 0.4, Y: 0.2, Z: 0.1}}
	world.Add(hittable.Sphere{Center: vector.Point3{X: -4, Y: 1, Z: 0}, Radius: 1.0, Mat: material2})

	material3 := material.Metal{Albedo: color.Color3{X: 0.7, Y: 0.6, Z: 0.5}, Fuzz: 0.0}
	world.Add(hittable.Sphere{Center: vector.Point3{X: 4, Y: 1, Z: 0}, Radius: 1.0, Mat: material3})

	return world
}

func buildWorldSceneTestFov() hittable.HittableList {
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

	cam.DefocusAngle = 10.0
	cam.FocusDistance = 3.4

	return cam
}

func generateCamBookCover() camera.Camera {
    var cam camera.Camera
    cam.AspectRatio = 16.0 / 9.0
    cam.ImageWidth = 1200
    cam.SamplesPerPixel = 500
    cam.MaxDepth = 500

    cam.Vfov = 20
    cam.Lookfrom = vector.Point3{X: 13, Y: 2, Z: 3}
    cam.Lookat = vector.Point3{X: 0, Y: 0, Z: 0}
    cam.Vup = vector.Vec3{X: 0, Y: 1, Z: 0}

    cam.DefocusAngle = 0.6
    cam.FocusDistance = 10.0

    return cam
}

func main() {
	world := buildWorldSceneBookCover()

	cam := generateCamBookCover()

	cam.Render(world)
}
