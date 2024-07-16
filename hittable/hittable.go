package hittable

import (
	"github.com/seaskythe/rayscape/interval"
	"github.com/seaskythe/rayscape/ray"
	"github.com/seaskythe/rayscape/vector"
)

type HitRecord struct {
	P        vector.Point3
	Normal   vector.Vec3
	T        float64
	FontFace bool
}

type Hittable interface {
	Hit(r ray.Ray, ray_t interval.Interval, rec *HitRecord) bool
}

func (h *HitRecord) SetFaceNormal(r ray.Ray, outward_normal vector.Vec3) {
	front_face := vector.Dot(r.Direction, outward_normal) < 0
	if front_face {
		(*h).Normal = outward_normal
	} else {
		(*h).Normal = vector.Scale(outward_normal, -1)
	}
}

