package material

import (
	"github.com/seaskythe/rayscape/color"
	"github.com/seaskythe/rayscape/ray"
	"github.com/seaskythe/rayscape/vector"
)

// Hit record is here to avoid circular dependency - removed from hittable/hittable.go
type HitRecord struct {
	P        vector.Point3
	Normal   vector.Vec3
	T        float64
	FontFace bool
	Mat      Material
}

func (h *HitRecord) SetFaceNormal(r ray.Ray, outward_normal vector.Vec3) {
	front_face := vector.Dot(r.Direction, outward_normal) < 0
	if front_face {
		(*h).Normal = outward_normal
	} else {
		(*h).Normal = vector.Scale(outward_normal, -1)
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
