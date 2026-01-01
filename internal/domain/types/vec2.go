package types

import "math"

// Vec2 represents a 2D vector with float32 precision
type Vec2 struct {
	X float32
	Y float32
}

// NewVec2 creates a new Vec2 with the given coordinates
func NewVec2(x, y float32) Vec2 {
	return Vec2{X: x, Y: y}
}

// Zero returns a zero vector
func Zero() Vec2 {
	return Vec2{X: 0, Y: 0}
}

// Add returns the sum of two vectors
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub returns the difference of two vectors
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{X: v.X - other.X, Y: v.Y - other.Y}
}

// Scale returns the vector scaled by a scalar
func (v Vec2) Scale(scalar float32) Vec2 {
	return Vec2{X: v.X * scalar, Y: v.Y * scalar}
}

// Magnitude returns the length of the vector
func (v Vec2) Magnitude() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// Normalize returns a unit vector in the same direction
func (v Vec2) Normalize() Vec2 {
	mag := v.Magnitude()
	if mag == 0 {
		return Zero()
	}
	return Vec2{X: v.X / mag, Y: v.Y / mag}
}

// AddMut adds another vector to this vector in-place
func (v *Vec2) AddMut(other Vec2) {
	v.X += other.X
	v.Y += other.Y
}

// ScaleMut scales this vector by a scalar in-place
func (v *Vec2) ScaleMut(scalar float32) {
	v.X *= scalar
	v.Y *= scalar
}
