package main

import "math"

// Vector : 3d vector
type Vector struct {
	x, y, z float64
}

// Length : L2 norm of vector
func (v *Vector) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

// LengthSquared : squared length of vector
func (v *Vector) LengthSquared() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

// Normalize : normalize vector
func (v *Vector) Normalize() {
	s := 1.0 / v.Length()
	v.x *= s
	v.y *= s
	v.z *= s
}

// Normalized : return normalized vector
func Normalized(v Vector) Vector {
	v.Normalize()
	return v
}

// Add : add two vectors
func Add(u Vector, v Vector) Vector {
	return Vector{u.x + v.x, u.y + v.y, u.z + v.z}
}

// Sub : sub two vectors
func Sub(u Vector, v Vector) Vector {
	return Vector{u.x - v.x, u.y - v.y, u.z - v.z}
}

// Mul : Scale vector by scalar
func Mul(u Vector, s float64) Vector {
	return Vector{u.x * s, u.y * s, u.z * s}
}

// Dot : dot product
func Dot(u, v Vector) float64 {
	return u.x*v.x + u.y*v.y + u.z*v.z
}
