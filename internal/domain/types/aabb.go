package types

// AABB represents an Axis-Aligned Bounding Box for collision detection
type AABB struct {
	X, Y          float32 // Top-left corner position
	Width, Height float32 // Dimensions
}

// NewAABB creates a new AABB
func NewAABB(x, y, width, height float32) AABB {
	return AABB{X: x, Y: y, Width: width, Height: height}
}

// Min returns the top-left corner (minimum point)
func (a AABB) Min() Vec2 {
	return Vec2{X: a.X, Y: a.Y}
}

// Max returns the bottom-right corner (maximum point)
func (a AABB) Max() Vec2 {
	return Vec2{X: a.X + a.Width, Y: a.Y + a.Height}
}

// Intersects checks if this AABB overlaps with another
func (a AABB) Intersects(b AABB) bool {
	return a.X < b.X+b.Width &&
		a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height &&
		a.Y+a.Height > b.Y
}

// Penetration calculates overlap for collision resolution
// Returns (dx, dy) where values indicate how to push 'a' out of 'b'
// Subtract from position to resolve: position.X -= dx
func (a AABB) Penetration(b AABB) (float32, float32) {
	overlapX := min(a.X+a.Width, b.X+b.Width) - max(a.X, b.X)
	overlapY := min(a.Y+a.Height, b.Y+b.Height) - max(a.Y, b.Y)

	if overlapX <= 0 || overlapY <= 0 {
		return 0, 0 // No intersection
	}

	var dx, dy float32

	// Determine push direction based on relative positions
	if a.X < b.X {
		dx = overlapX // Push left (positive = move right to separate)
	} else {
		dx = -overlapX // Push right (negative = move left to separate)
	}

	if a.Y < b.Y {
		dy = overlapY // Push up (positive = move down to separate)
	} else {
		dy = -overlapY // Push down (negative = move up to separate)
	}

	return dx, dy
}
