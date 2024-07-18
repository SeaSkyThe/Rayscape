package hittable

import (
	"math"

	"github.com/seaskythe/rayscape/interval"
	"github.com/seaskythe/rayscape/material"
	"github.com/seaskythe/rayscape/ray"
	"github.com/seaskythe/rayscape/vector"
)

type Sphere struct {
	Center vector.Point3
	Radius float64
    Mat material.Material
}

func (s Sphere) Hit(r ray.Ray, ray_t interval.Interval, rec *material.HitRecord) bool {
	oc := vector.Subtract(s.Center, r.Origin) // Vector from the ray origin to the sphere center
	a := r.Direction.LengthSquared()
	half_b := vector.Dot(r.Direction, oc)
	c := oc.LengthSquared() - s.Radius*s.Radius

	discriminant := half_b*half_b - a*c

	if discriminant < 0 { // No intersection
		return false
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range
	root := (half_b - sqrtd) / a  // Check if root is in the interval
	if !ray_t.Surrounds(root) {
		root = (half_b + sqrtd) / a // Check other root
		if !ray_t.Surrounds(root) {
			return false
		}
	}

	rec.T = root
	rec.P = r.At(rec.T)
	rec.Normal = vector.Divide(vector.Subtract(rec.P, s.Center), s.Radius)

	var outward_normal vector.Vec3 = vector.Divide(vector.Subtract(rec.P, s.Center), s.Radius)
	rec.SetFaceNormal(r, outward_normal)
    rec.Mat = s.Mat

	return true

}
