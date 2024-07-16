package hittable

import (
	"math"

	"github.com/seaskythe/rayscape/interval"
	"github.com/seaskythe/rayscape/ray"
	"github.com/seaskythe/rayscape/vector"
)

type Sphere struct {
	Center vector.Point3
	Radius float64
}

func (s Sphere) Hit(r ray.Ray, ray_t interval.Interval, rec *HitRecord) bool {
	oc := vector.Subtract(s.Center, r.Origin)
	a := r.Direction.LengthSquared()
	half_b := vector.Dot(r.Direction, oc)
	c := oc.LengthSquared() - s.Radius*s.Radius

	discriminant := half_b*half_b - a*c

	if discriminant < 0 {
		return false
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range
	root := (half_b - sqrtd) / a
	if !ray_t.Sorrounds(root) {
		root = (half_b + sqrtd) / a
		if !ray_t.Sorrounds(root) {
			return false
		}
	}

	rec.T = root
	rec.P = r.At(rec.T)
	rec.Normal = vector.Divide(vector.Subtract(rec.P, s.Center), s.Radius)

	var outward_normal vector.Vec3 = vector.Divide(vector.Subtract(rec.P, s.Center), s.Radius)
	rec.SetFaceNormal(r, outward_normal)

	return true

}
