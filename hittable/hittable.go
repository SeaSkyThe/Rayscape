package hittable

import (
	"github.com/seaskythe/rayscape/interval"
	"github.com/seaskythe/rayscape/material"
	"github.com/seaskythe/rayscape/ray"
)

type Hittable interface {
	Hit(r ray.Ray, ray_t interval.Interval, rec *material.HitRecord) bool
}
