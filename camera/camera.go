package camera

import (
	"fmt"
	"os"

	"github.com/seaskythe/rayscape/color"
	"github.com/seaskythe/rayscape/hittable"
	"github.com/seaskythe/rayscape/interval"
	"github.com/seaskythe/rayscape/ray"
	"github.com/seaskythe/rayscape/rtweekend"
	"github.com/seaskythe/rayscape/vector"
)

type Camera struct {
	AspectRatio float64
	ImageWidth  int
	ImageHeight int
	Center      vector.Point3
	Pixel00Loc  vector.Point3
	PixelDeltaU vector.Vec3
	PixelDeltaV vector.Vec3
}

func (c *Camera) Render(world hittable.Hittable) {
	c.Initialize()

	file := CreatePPMImage(c.ImageWidth, c.ImageHeight)
	defer file.Close()

	for j := 0; j < c.ImageHeight; j++ {
		fmt.Fprintf(os.Stderr, "\033[2K\rScanlines remaining: %d", c.ImageHeight-j)
		os.Stderr.Sync()
		for i := 0; i < c.ImageWidth; i++ {
			pixel_delta_u_i := vector.Scale(c.PixelDeltaU, float64(i))
			pixel_delta_v_j := vector.Scale(c.PixelDeltaV, float64(j))
			pixel_deltas := vector.Add(pixel_delta_u_i, pixel_delta_v_j)
			var pixel_center = vector.Add(c.Pixel00Loc, pixel_deltas)
			var ray_direction = vector.Subtract(pixel_center, c.Center)
			var r ray.Ray = ray.Ray{Origin: c.Center, Direction: ray_direction}

			pixel_color := c.RayColor(r, world)
			color.WriteColor(file, pixel_color)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone!")

}

func (c *Camera) Initialize() {
	// Calculate image height and make sure that its at least 1
	c.ImageHeight = int(float64(c.ImageWidth) / c.AspectRatio)
	if c.ImageHeight < 1 {
		c.ImageHeight = 1
	}

	// Camera
	var focal_length float64 = 1.0
	var viewport_height float64 = 2.0
	var viewport_width float64 = viewport_height * float64(float64(c.ImageWidth)/float64(c.ImageHeight))
	var camera_center vector.Point3 = vector.Point3{X: 0, Y: 0, Z: 0}

	// Calculate the vectors across the horizontal and down the vertical viewport edges
	var viewport_u vector.Vec3 = vector.Vec3{X: viewport_width, Y: 0, Z: 0}
	var viewport_v vector.Vec3 = vector.Vec3{X: 0, Y: -viewport_height, Z: 0}

	// Calculate the horizontal and vertical delta vectors from pixel to pixel
	c.PixelDeltaU = vector.Divide(viewport_u, float64(c.ImageWidth))
	c.PixelDeltaV = vector.Divide(viewport_v, float64(c.ImageHeight))

	// Calculate the location of the upper left pixel
	var viewport_upper_left = vector.Vec3{
		X: camera_center.X - viewport_u.X/2 - viewport_v.X/2,
		Y: camera_center.Y - viewport_u.Y/2 - viewport_v.Y/2,
		Z: camera_center.Z - viewport_u.Z/2 - viewport_v.Z/2 - focal_length,
	}

	var pixel_delta_u_plus_v = vector.Add(c.PixelDeltaU, c.PixelDeltaV)
	c.Pixel00Loc = vector.Add(
		viewport_upper_left,
		vector.Scale(pixel_delta_u_plus_v, 0.5),
	)
}

func (c Camera) RayColor(r ray.Ray, world hittable.Hittable) color.Color3 {
	var rec hittable.HitRecord
	if (world.Hit(r, interval.Interval{Min: 0, Max: rtweekend.INFINITY}, &rec)) {
		return vector.Scale(vector.Add(rec.Normal, color.Color3{X: 1, Y: 1, Z: 1}), 0.5)
	}

	unit_direction := vector.UnitVector(r.Direction)
	a := 0.5 * (unit_direction.Y + 1.0)

	white := color.Color3{X: 1.0, Y: 1.0, Z: 1.0}
	blue := color.Color3{X: 0.5, Y: 0.7, Z: 1.0}

	return vector.Add(vector.Scale(white, 1.0-a), vector.Scale(blue, a))
}

func CreatePPMImage(width, height int) *os.File {
	file, err := os.OpenFile("image.ppm", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(file, "P3\n%d %d\n255\n", width, height)

	return file
}
