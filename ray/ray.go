package ray

import "github.com/seaskythe/rayscape/vector"

type Ray struct {
	Origin    vector.Point3
	Direction vector.Vec3
}

func (ray *Ray) At(t float64) vector.Vec3 {
	return vector.Add(ray.Origin, vector.Scale(ray.Direction, t))
}
