package gamemath

import "math"

// Vec2 is a 2D vector.
type Vec2 struct {
	X, Y float64
}

// Zero is the zero vector.
var Zero = Vec2{0, 0}

// Add returns the sum of two vectors.
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

// Sub returns the difference of two vectors.
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{v.X - other.X, v.Y - other.Y}
}

// Scale returns the vector scaled by a scalar.
func (v Vec2) Scale(s float64) Vec2 {
	return Vec2{v.X * s, v.Y * s}
}

// Length returns the magnitude of the vector.
func (v Vec2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
