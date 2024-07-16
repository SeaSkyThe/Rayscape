package main

import (
	"fmt"
	"os"

	"github.com/seaskythe/rayscape/color"
	"github.com/seaskythe/rayscape/hittable"
	"github.com/seaskythe/rayscape/interval"
	"github.com/seaskythe/rayscape/ray"
	"github.com/seaskythe/rayscape/vector"
)

func ray_color(ray ray.Ray, world hittable.Hittable) color.Color3 {
	var rec hittable.HitRecord
	if (world.Hit(ray, interval.Interval{Min: 0, Max: infinity}, &rec)) {
		return vector.Scale(vector.Add(rec.Normal, color.Color3{X: 1, Y: 1, Z: 1}), 0.5)
	}

	unit_direction := vector.UnitVector(ray.Direction)
	a := 0.5 * (unit_direction.Y + 1.0)

	white := color.Color3{X: 1.0, Y: 1.0, Z: 1.0}
	blue := color.Color3{X: 0.5, Y: 0.7, Z: 1.0}

	return vector.Add(vector.Scale(white, 1.0-a), vector.Scale(blue, a))
}

func main() {
	// Image
	var aspect_ratio float64 = 16.0 / 9.0
	var image_width int = 400

	// Calculate image height and make sure that its at least 1
	var image_height int = int(float64(image_width) / aspect_ratio)
	if image_height < 1 {
		image_height = 1
	}

	// World
	var world hittable.HittableList
	world.Add(hittable.Sphere{Center: vector.Point3{X: 0, Y: 0, Z: -1}, Radius: 0.5})
	world.Add(hittable.Sphere{Center: vector.Point3{X: 0, Y: -100.5, Z: -1}, Radius: 100})

	// Camera
	var focal_length float64 = 1.0
	var viewport_height float64 = 2.0
	var viewport_width float64 = viewport_height * float64(float64(image_width)/float64(image_height))
	var camera_center vector.Point3 = vector.Point3{X: 0, Y: 0, Z: 0}

	// Calculate the vectors across the horizontal and down the vertical viewport edges
	var viewport_u vector.Vec3 = vector.Vec3{X: viewport_width, Y: 0, Z: 0}
	var viewport_v vector.Vec3 = vector.Vec3{X: 0, Y: -viewport_height, Z: 0}

	// Calculate the horizontal and vertical delta vectors from pixel to pixel
	var pixel_delta_u = vector.Divide(viewport_u, float64(image_width))
	var pixel_delta_v = vector.Divide(viewport_v, float64(image_height))

	// Calculate the location of the upper left pixel
	var viewport_upper_left = vector.Vec3{
		X: camera_center.X - viewport_u.X/2 - viewport_v.X/2,
		Y: camera_center.Y - viewport_u.Y/2 - viewport_v.Y/2,
		Z: camera_center.Z - viewport_u.Z/2 - viewport_v.Z/2 - focal_length,
	}

	var pixel_delta_u_plus_v = vector.Add(pixel_delta_u, pixel_delta_v)
	var pixel00_loc = vector.Add(
		viewport_upper_left,
		vector.Scale(pixel_delta_u_plus_v, 0.5),
	)

	file := CreatePPMImage(image_width, image_height)
	defer file.Close()

	for j := 0; j < image_height; j++ {
		fmt.Fprintf(os.Stderr, "\033[2K\rScanlines remaining: %d", image_height-j)
		os.Stderr.Sync()
		for i := 0; i < image_width; i++ {
			pixel_delta_u_i := vector.Scale(pixel_delta_u, float64(i))
			pixel_delta_v_j := vector.Scale(pixel_delta_v, float64(j))
			pixel_deltas := vector.Add(pixel_delta_u_i, pixel_delta_v_j)
			var pixel_center = vector.Add(pixel00_loc, pixel_deltas)
			var ray_direction = vector.Subtract(pixel_center, camera_center)
			var r ray.Ray = ray.Ray{Origin: camera_center, Direction: ray_direction}

			pixel_color := ray_color(r, world)
			color.WriteColor(file, pixel_color)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone!")
}
