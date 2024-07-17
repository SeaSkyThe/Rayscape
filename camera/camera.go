package camera

import (
	"fmt"
	"os"
	"time"

	"github.com/seaskythe/rayscape/color"
	"github.com/seaskythe/rayscape/hittable"
	"github.com/seaskythe/rayscape/interval"
	"github.com/seaskythe/rayscape/ray"
	"github.com/seaskythe/rayscape/rtweekend"
	"github.com/seaskythe/rayscape/vector"
)

type Camera struct {
	AspectRatio       float64
	ImageWidth        int
	ImageHeight       int
	Center            vector.Point3
	Pixel00Loc        vector.Point3
	PixelDeltaU       vector.Vec3
	PixelDeltaV       vector.Vec3
	SamplesPerPixel   int
	PixelSamplesScale float64
	MaxDepth          int
}

func (c *Camera) Render(world hittable.Hittable) {
    start_time := time.Now()
    fmt.Fprintln(os.Stderr, "\nRender Started.")
	c.Initialize()

	file := CreatePPMImage(c.ImageWidth, c.ImageHeight)
	defer file.Close()

	for j := 0; j < c.ImageHeight; j++ {
		fmt.Fprintf(os.Stderr, "\033[2K\rScanlines remaining: %d", c.ImageHeight-j)
		os.Stderr.Sync()
		for i := 0; i < c.ImageWidth; i++ {
			pixel_color := color.Color3{X: 0, Y: 0, Z: 0}

			for sample := 0; sample < c.SamplesPerPixel; sample++ {
				r := c.GetRay(i, j)
				pixel_color = vector.Add(pixel_color, c.RayColor(r, c.MaxDepth, world))
			}
			color.WriteColor(file, vector.Scale(pixel_color, c.PixelSamplesScale))
		}
	}
	fmt.Fprintf(os.Stderr, "\nRender Finished in %f s!\n", time.Since(start_time).Seconds())

}

func (c *Camera) Initialize() {
	// Calculate image height and make sure that its at least 1
	c.ImageHeight = int(float64(c.ImageWidth) / c.AspectRatio)
	if c.ImageHeight < 1 {
		c.ImageHeight = 1
	}
	c.PixelSamplesScale = 1.0 / float64(c.SamplesPerPixel)

	// Camera
	var focal_length float64 = 1.0
	var viewport_height float64 = 2.0
	var viewport_width float64 = viewport_height * float64(float64(c.ImageWidth)/float64(c.ImageHeight))
	c.Center = vector.Point3{X: 0, Y: 0, Z: 0}

	// Calculate the vectors across the horizontal and down the vertical viewport edges
	var viewport_u vector.Vec3 = vector.Vec3{X: viewport_width, Y: 0, Z: 0}
	var viewport_v vector.Vec3 = vector.Vec3{X: 0, Y: -viewport_height, Z: 0}

	// Calculate the horizontal and vertical delta vectors from pixel to pixel
	c.PixelDeltaU = vector.Divide(viewport_u, float64(c.ImageWidth))
	c.PixelDeltaV = vector.Divide(viewport_v, float64(c.ImageHeight))

	// Calculate the location of the upper left pixel
	var viewport_upper_left = vector.Vec3{
		X: c.Center.X - viewport_u.X/2 - viewport_v.X/2,
		Y: c.Center.Y - viewport_u.Y/2 - viewport_v.Y/2,
		Z: c.Center.Z - viewport_u.Z/2 - viewport_v.Z/2 - focal_length,
	}

	var pixel_delta_u_plus_v = vector.Add(c.PixelDeltaU, c.PixelDeltaV)
	c.Pixel00Loc = vector.Add(
		viewport_upper_left,
		vector.Scale(pixel_delta_u_plus_v, 0.5),
	)
}

func (c Camera) GetRay(i, j int) ray.Ray {
	// Construct a camera Ray originating from the origin and
	// directed at randomly sampled point around the pixel location i, j
	offset := c.SampleSquare()

	pixel_delta_U_sample := vector.Scale(c.PixelDeltaU, float64(i)+offset.X)
	pixel_delta_V_sample := vector.Scale(c.PixelDeltaV, float64(j)+offset.Y)
	pixel_delta_UV_sample := vector.Add(pixel_delta_U_sample, pixel_delta_V_sample)
	pixel_sample := vector.Add(c.Pixel00Loc, pixel_delta_UV_sample)

	ray_origin := c.Center
	ray_direction := vector.Subtract(pixel_sample, ray_origin)

	return ray.Ray{Origin: ray_origin, Direction: ray_direction}
}

func (c Camera) SampleSquare() vector.Vec3 {
	// Returns the vector to a random point in the [-.5,-.5]-[+.5,+.5] unit square.
	return vector.Vec3{X: rtweekend.RandomDouble() - 0.5, Y: rtweekend.RandomDouble() - 0.5, Z: 0}
}

// Shader
func (c Camera) RayColor(r ray.Ray, depth int, world hittable.Hittable) color.Color3 {
	// If we exceed the ray bounce limit, no more light is processed
	if depth <= 0 {
		return color.Color3{X: 0, Y: 0, Z: 0}
	}

	var rec hittable.HitRecord
	if (world.Hit(r, interval.Interval{Min: 0.001, Max: rtweekend.INFINITY}, &rec)) {
		// var direction vector.Vec3 = vector.RandomOnHemisphere(rec.Normal) // Reflect Rays Evenly
        var direction = vector.Add(rec.Normal, vector.RandomUnitVector()) // Lambertian Reflection (reflect randomly towards normal)
        reflectance := 0.2
		return vector.Scale(c.RayColor(ray.Ray{Origin: rec.P, Direction: direction}, depth-1, world), reflectance)
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
