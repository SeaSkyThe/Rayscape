package main

import (
	"fmt"
	"math"
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

func PrintVec3(v Vec3) {
	fmt.Printf("%f %f %f\n", v.X, v.Y, v.Z)
}

