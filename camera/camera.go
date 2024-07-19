package camera

import (
	"fmt"
	"math"
	"os"
	// "sync"
	"time"

	"github.com/seaskythe/rayscape/color"
	"github.com/seaskythe/rayscape/hittable"
	"github.com/seaskythe/rayscape/interval"
	"github.com/seaskythe/rayscape/material"
	"github.com/seaskythe/rayscape/ray"
	"github.com/seaskythe/rayscape/rtweekend"
	"github.com/seaskythe/rayscape/vector"
)

var ray_bounces_count int = 0

type Camera struct {
	AspectRatio     float64 // Ratio of image width over height
	ImageWidth      int     // Image width in pixels
	SamplesPerPixel int     // Count of random ray samples to generate each pixel color
	MaxDepth        int     // Maximum number of ray bounces into scene

	Vfov     float64       // Vertical view angle (field of view)
	Lookfrom vector.Point3 // Point camera is looking from
	Lookat   vector.Point3 // Point camera is looking at
	Vup      vector.Vec3   // Camera up vector

	DefocusAngle  float64 // Variation angle of rays through each pixel
	FocusDistance float64 // Distance from camera lookfrom point to plane of perfect focus

	// Should be private
	ImageHeight       int           // Image height in pixels
	PixelSamplesScale float64       // Color scale factor for a sum of pixel samples
	Center            vector.Point3 // Camera Center
	Pixel00Loc        vector.Point3 // Location of pixel (0, 0)
	PixelDeltaU       vector.Vec3   // Offset to pixel to the right
	PixelDeltaV       vector.Vec3   // Offset to pixel down
	U, V, W           vector.Vec3   // Camera frame basis vectors
	defocusDiskU      vector.Vec3   // Defocus disk horizontal radius
	defocusDiskV      vector.Vec3   // Defocus disk vertical radius
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
			// fmt.Println(pixel_color, vector.Scale(pixel_color, c.PixelSamplesScale))
			color.WriteColor(file, vector.Scale(pixel_color, c.PixelSamplesScale))
		}
	}
	fmt.Fprintf(os.Stderr, "\nRender Finished in %fs | Ray Bounces: %d\n", time.Since(start_time).Seconds(), ray_bounces_count)
}


// func (c *Camera) Render(world hittable.Hittable) {
// 	start_time := time.Now()
// 	fmt.Fprintln(os.Stderr, "\nRender Started.")
// 	c.Initialize()
//
// 	// Create an in-memory buffer for pixel data
// 	pixelData := make([][]color.Color3, c.ImageHeight)
// 	for i := range pixelData {
// 		pixelData[i] = make([]color.Color3, c.ImageWidth)
// 	}
//
// 	var wg sync.WaitGroup
//     semaphore := make(chan struct{}, 10)
//
//     count_goroutines := 0
// 	for j := 0; j < c.ImageHeight; j++ {
// 		wg.Add(1)
// 		go func(j int) {
//             // fmt.Fprintf(os.Stderr, "\033[2K\rScanlines remaining: %d", c.ImageHeight-j)
//             //
// 			defer wg.Done()
//             semaphore <- struct{}{}
//             defer func() {<-semaphore}()
//
// 			os.Stderr.Sync()
// 			for i := 0; i < c.ImageWidth; i++ {
// 				pixel_color := color.Color3{X: 0, Y: 0, Z: 0}
//
// 				for sample := 0; sample < c.SamplesPerPixel; sample++ {
// 					r := c.GetRay(i, j)
// 					pixel_color = vector.Add(pixel_color, c.RayColor(r, c.MaxDepth, world))
// 				}
// 				pixelData[j][i] = vector.Scale(pixel_color, c.PixelSamplesScale)
// 			}
//             count_goroutines += 1
//             fmt.Println("Finished: ", j, " | ", count_goroutines)
// 		}(j)
// 	}
//
//     fmt.Println("Waiting for all goroutines to finish...")
// 	wg.Wait()
//
// 	// Write the in-memory buffer to the file
// 	file := CreatePPMImage(c.ImageWidth, c.ImageHeight)
// 	defer file.Close()
//
// 	for j := 0; j < c.ImageHeight; j++ {
// 		for i := 0; i < c.ImageWidth; i++ {
// 			color.WriteColor(file, pixelData[j][i])
// 		}
// 	}
//
// 	fmt.Fprintf(os.Stderr, "\nRender Finished in %fs\n", time.Since(start_time).Seconds())
// }

func (c *Camera) Initialize() {
	// Calculate image height and make sure that its at least 1
	c.ImageHeight = int(float64(c.ImageWidth) / c.AspectRatio)
	if c.ImageHeight < 1 {
		c.ImageHeight = 1
	}
	c.PixelSamplesScale = 1.0 / float64(c.SamplesPerPixel)

	c.Center = c.Lookfrom

	// Determine viewport dimensions
	var theta = rtweekend.DegreesToRadians(c.Vfov)
	var h = math.Tan(theta / 2)
	var viewport_height float64 = 2.0 * h * c.FocusDistance
	var viewport_width float64 = viewport_height * float64(float64(c.ImageWidth)/float64(c.ImageHeight))

	// Calculate the u,v,w unit basis vector for the camera coordinate frame
	c.W = vector.UnitVector(vector.Subtract(c.Lookfrom, c.Lookat))
	c.U = vector.UnitVector(vector.Cross(c.Vup, c.W))
	c.V = vector.Cross(c.W, c.U)

	// Calculate the vectors across the horizontal and down the vertical viewport edges
	var viewport_u vector.Vec3 = vector.Scale(c.U, viewport_width)
	var viewport_v vector.Vec3 = vector.Scale(c.V.Negate(), viewport_height)

	// Calculate the horizontal and vertical delta vectors from pixel to pixel
	c.PixelDeltaU = vector.Divide(viewport_u, float64(c.ImageWidth))
	c.PixelDeltaV = vector.Divide(viewport_v, float64(c.ImageHeight))

	// Calculate the location of the upper left pixel
	var viewport_upper_left = vector.Vec3{
		X: c.Center.X - (c.W.X * c.FocusDistance) - viewport_u.X/2 - viewport_v.X/2,
		Y: c.Center.Y - (c.W.Y * c.FocusDistance) - viewport_u.Y/2 - viewport_v.Y/2,
		Z: c.Center.Z - (c.W.Z * c.FocusDistance) - viewport_u.Z/2 - viewport_v.Z/2,
	}

	var pixel_delta_u_plus_v = vector.Add(c.PixelDeltaU, c.PixelDeltaV)
	c.Pixel00Loc = vector.Add(
		viewport_upper_left,
		vector.Scale(pixel_delta_u_plus_v, 0.5),
	)

	// Calculate the camera defocus disk basis vectors
	defocus_radius := c.FocusDistance * math.Tan(rtweekend.DegreesToRadians(c.DefocusAngle/2))
	c.defocusDiskU = vector.Scale(c.U, defocus_radius)
	c.defocusDiskV = vector.Scale(c.V, defocus_radius)
}

func (c Camera) GetRay(i, j int) ray.Ray {
	// Construct a camera Ray originating from the defocus disk and
	// directed at randomly sampled point around the pixel location i, j
	offset := c.SampleSquare()

	pixel_delta_U_sample := vector.Scale(c.PixelDeltaU, float64(i)+offset.X)
	pixel_delta_V_sample := vector.Scale(c.PixelDeltaV, float64(j)+offset.Y)
	pixel_delta_UV_sample := vector.Add(pixel_delta_U_sample, pixel_delta_V_sample)
	pixel_sample := vector.Add(c.Pixel00Loc, pixel_delta_UV_sample)

    ray_origin := c.Center
	if c.DefocusAngle > 0 {
        ray_origin = c.DefocusDiskSample()
	}

	ray_direction := vector.Subtract(pixel_sample, ray_origin)

	return ray.Ray{Origin: ray_origin, Direction: ray_direction}
}

func (c Camera) SampleSquare() vector.Vec3 {
	// Returns the vector to a random point in the [-.5,-.5]-[+.5,+.5] unit square.
	return vector.Vec3{X: rtweekend.RandomDouble() - 0.5, Y: rtweekend.RandomDouble() - 0.5, Z: 0}
}

func (c Camera) DefocusDiskSample() vector.Vec3 {
	// Returns a random point in the camera defocus disk
	p := vector.RandomInUnitDisk()
	return vector.Add(c.Center, vector.Add(vector.Scale(c.defocusDiskU, p.X), vector.Scale(c.defocusDiskV, p.Y)))
}

// Shader
func (c Camera) RayColor(r ray.Ray, depth int, world hittable.Hittable) color.Color3 {
	ray_bounces_count += 1
	// If we exceed the ray bounce limit, no more light is processed
	if depth <= 0 {
		return color.Color3{X: 0, Y: 0, Z: 0}
	}
	var rec material.HitRecord
	if (world.Hit(r, interval.Interval{Min: 0.001, Max: rtweekend.INFINITY}, &rec)) {
		var scattered ray.Ray
		var attenuation color.Color3
		if rec.Mat.Scatter(r, rec, &attenuation, &scattered) {
			// fmt.Println("attenuation = ", attenuation, "nextColor = ", c.RayColor(scattered, depth - 1, world), "resultColor = ", vector.Multiply(attenuation, c.RayColor(scattered, depth-1, world)))
			return vector.Multiply(attenuation, c.RayColor(scattered, depth-1, world))
		}
		return color.Color3{X: 0, Y: 0, Z: 0}
	}

	// If its not any object, just render a background
	// fmt.Println("here")
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
