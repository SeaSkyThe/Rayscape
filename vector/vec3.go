package vector

import (
	"fmt"
	"math"

	"github.com/seaskythe/rayscape/rtweekend"
)

type Vec3 struct {
	X, Y, Z float64
}

// Create New Vec3
func NewVec3(x, y, z float64) Vec3 {
	return Vec3{X: x, Y: y, Z: z}
}

// Aliases
type (
	Point3 = Vec3
)

// Negate returns the negation of the vector
func (v Vec3) Negate() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

// Length returns the length of the vector
func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// LengthSquared returns the square of the length of the vector
func (v Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3) NearZero() bool {
	s := 1e-8
	return math.Abs(v.X) < s && math.Abs(v.Y) < s && math.Abs(v.Z) < s
}

// Vector Utility Functions
func Add(u, v Vec3) Vec3 {
	return Vec3{u.X + v.X, u.Y + v.Y, u.Z + v.Z}
}

func Subtract(u, v Vec3) Vec3 {
	return Vec3{u.X - v.X, u.Y - v.Y, u.Z - v.Z}
}

func Multiply(u, v Vec3) Vec3 {
	return Vec3{u.X * v.X, u.Y * v.Y, u.Z * v.Z}
}

func Scale(u Vec3, t float64) Vec3 {
	return Vec3{u.X * t, u.Y * t, u.Z * t}
}

func Divide(v Vec3, t float64) Vec3 {
	return Scale(v, 1/t)
}

func Dot(u, v Vec3) float64 {
	return u.X*v.X + u.Y*v.Y + u.Z*v.Z
}

func Cross(u, v Vec3) Vec3 {
	return Vec3{
		u.Y*v.Z - u.Z*v.Y,
		u.Z*v.X - u.X*v.Z,
		u.X*v.Y - u.Y*v.X,
	}
}

func UnitVector(v Vec3) Vec3 {
	length := v.Length()
	if length == 0 {
		return v // Return original vector if length is zero to avoid division by zero
	}
	return Divide(v, length)
}

func Random() Vec3 {
	return Vec3{
		X: rtweekend.RandomDouble(),
		Y: rtweekend.RandomDouble(),
		Z: rtweekend.RandomDouble(),
	}
}

func RandomInInterval(min, max float64) Vec3 {
	return Vec3{
		X: rtweekend.RandomDoubleInInterval(min, max),
		Y: rtweekend.RandomDoubleInInterval(min, max),
		Z: rtweekend.RandomDoubleInInterval(min, max),
	}
}

func RandomInUnitSphere() Vec3 {
	p := RandomInInterval(-1, 1)
	for true {
		p = RandomInInterval(-1, 1)
		if p.LengthSquared() < 1 {
			return p
		}
	}
	return p
}

func RandomUnitVector() Vec3 {
	return UnitVector(RandomInUnitSphere())
}

func RandomOnHemisphere(normal Vec3) Vec3 {
	var on_unit_sphere = RandomUnitVector()
	if Dot(on_unit_sphere, normal) > 0.0 { // In the same hemisphere of the normal (ray reflection correct)
		return on_unit_sphere
	} else {
		return Scale(on_unit_sphere, -1)
	}
}

func Reflect(v Vec3, normal Vec3) Vec3 {
	return Subtract(v, Scale(normal, Dot(v, normal)*2))
}

func PrintVec3(v Vec3) {
	fmt.Printf("%f %f %f\n", v.X, v.Y, v.Z)
}
