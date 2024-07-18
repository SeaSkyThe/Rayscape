package material

import (
	"math"

	"github.com/seaskythe/rayscape/color"
	"github.com/seaskythe/rayscape/ray"
	"github.com/seaskythe/rayscape/rtweekend"
	"github.com/seaskythe/rayscape/vector"
)

// Hit record is here to avoid circular dependency - removed from hittable/hittable.go
type HitRecord struct {
	P         vector.Point3
	Normal    vector.Vec3
	T         float64
	FrontFace bool
	Mat       Material
}

func (h *HitRecord) SetFaceNormal(r ray.Ray, outward_normal vector.Vec3) {
	h.FrontFace = vector.Dot(r.Direction, outward_normal) < 0
	if h.FrontFace {
		h.Normal = outward_normal
	} else {
		h.Normal = outward_normal.Negate()
	}
}

// Material is the interface for all material types
type Material interface {
	Scatter(r_in ray.Ray, rec HitRecord, attenuation *color.Color3, scattered *ray.Ray) bool
}

// Lambertian is Diffuse material
type Lambertian struct {
	Albedo color.Color3 // Latin for "whiteness"
}

func (l Lambertian) Scatter(r_in ray.Ray, rec HitRecord, attenuation *color.Color3, scattered *ray.Ray) bool {
	scatter_direction := vector.Add(rec.Normal, vector.RandomUnitVector())
	if scatter_direction.NearZero() {
		scatter_direction = rec.Normal
	}
	*scattered = ray.Ray{Origin: rec.P, Direction: scatter_direction}
	*attenuation = l.Albedo
	return true
}

// Metal Material

type Metal struct {
	Albedo color.Color3
	Fuzz   float64
}

func (m Metal) Scatter(r_in ray.Ray, rec HitRecord, attenuation *color.Color3, scattered *ray.Ray) bool {
	var reflected vector.Vec3 = vector.Reflect(r_in.Direction, rec.Normal)
	reflected = vector.Add(vector.UnitVector(reflected), vector.Scale(vector.RandomUnitVector(), m.Fuzz))
	*scattered = ray.Ray{Origin: rec.P, Direction: reflected}
	*attenuation = m.Albedo

	return vector.Dot(scattered.Direction, rec.Normal) > 0
}

// Dieletric
type Dieletric struct {
	RefractionIndex float64
}

func (d Dieletric) Scatter(r_in ray.Ray, rec HitRecord, attenuation *color.Color3, scattered *ray.Ray) bool {
	*attenuation = color.Color3{X: 1.0, Y: 1.0, Z: 1.0}
	// Review ri and rec.FontFace and rec.Normal
	// how its coming to this function
	var ri float64 = d.RefractionIndex
	if rec.FrontFace {
		ri = 1.0 / d.RefractionIndex
	}

	var unit_direction vector.Vec3 = vector.UnitVector(r_in.Direction)

    cos_theta := math.Min(vector.Dot(unit_direction.Negate(), rec.Normal), 1.0)
    sin_theta := math.Sqrt(1.0 - cos_theta*cos_theta)

    cannot_refract := ri * sin_theta > 1.0
    var direction vector.Vec3
    if cannot_refract || d.Reflectance(cos_theta, ri) > rtweekend.RandomDouble(){
        direction = vector.Reflect(unit_direction, rec.Normal)
    } else {
        direction = vector.Refract(unit_direction, rec.Normal, ri)
    }

	*scattered = ray.Ray{Origin: rec.P, Direction: direction}
	return true
}

func (d Dieletric) Reflectance(cosine float64, refraction_index float64) float64 {
    var r0 float64 = (1.0 - d.RefractionIndex) / (1.0 + d.RefractionIndex )
    r0 = r0 * r0
    return r0 + (1.0 - r0) * math.Pow(1.0 - cosine, 5)
}
