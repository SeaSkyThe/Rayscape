package hittable

import (
	"github.com/seaskythe/rayscape/ray"
)

// HittableList is a list of hittable objects
type HittableList []Hittable

func (h *HittableList) Clear() {
	for i := range *h {
		(*h)[i] = nil
	}
}

func (h *HittableList) Add(obj Hittable) {
	*h = append(*h, obj)
}

func (h HittableList) Hit(r ray.Ray, t_min float64, t_max float64, rec *HitRecord) bool {
    var temp_rec HitRecord
    var hit_anything bool = false
    closest_so_far := t_max

    for _, obj := range h {
        if obj.Hit(r, t_min, closest_so_far, &temp_rec) {
            hit_anything = true
            closest_so_far = temp_rec.T
            *rec = temp_rec
        }
    }

    return hit_anything

}
