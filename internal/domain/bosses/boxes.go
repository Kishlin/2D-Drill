package bosses

import "github.com/Kishlin/drill-game/internal/domain/types"

// CollisionBox represents a physical presence that blocks player movement
type CollisionBox struct {
	ID     string
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// AABB returns the bounding box for this collision box
func (c CollisionBox) AABB() types.AABB {
	return types.AABB{
		X:      c.X,
		Y:      c.Y,
		Width:  c.Width,
		Height: c.Height,
	}
}

// Hitbox represents an attack zone that damages the player on intersection
type Hitbox struct {
	ID           string
	X            float32
	Y            float32
	Width        float32
	Height       float32
	DamagePerSec float32
}

// AABB returns the bounding box for this hitbox
func (h Hitbox) AABB() types.AABB {
	return types.AABB{
		X:      h.X,
		Y:      h.Y,
		Width:  h.Width,
		Height: h.Height,
	}
}

// Hurtbox represents a vulnerable zone where the boss receives damage
type Hurtbox struct {
	ID               string
	X                float32
	Y                float32
	Width            float32
	Height           float32
	DamageMultiplier float32 // 1.0 = normal, 2.0 = weak point
}

// AABB returns the bounding box for this hurtbox
func (h Hurtbox) AABB() types.AABB {
	return types.AABB{
		X:      h.X,
		Y:      h.Y,
		Width:  h.Width,
		Height: h.Height,
	}
}
