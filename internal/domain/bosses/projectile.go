package bosses

import "github.com/Kishlin/drill-game/internal/domain/types"

// Projectile represents a projectile entity that can be spawned by bosses
type Projectile struct {
	AABB     types.AABB
	Velocity types.Vec2
	Damage   float32
	Active   bool
}

// NewProjectile creates a new projectile at the given position with the given velocity and damage
func NewProjectile(x, y, width, height float32, velocity types.Vec2, damage float32) *Projectile {
	return &Projectile{
		AABB:     types.AABB{X: x, Y: y, Width: width, Height: height},
		Velocity: velocity,
		Damage:   damage,
		Active:   true,
	}
}

// Update moves the projectile based on velocity and delta time
func (p *Projectile) Update(dt float32) {
	if !p.Active {
		return
	}
	p.AABB.X += p.Velocity.X * dt
	p.AABB.Y += p.Velocity.Y * dt
}

// Intersects checks if this projectile overlaps with the given AABB
func (p *Projectile) Intersects(aabb types.AABB) bool {
	return p.AABB.Intersects(aabb)
}

// Deactivate marks the projectile as inactive
func (p *Projectile) Deactivate() {
	p.Active = false
}
