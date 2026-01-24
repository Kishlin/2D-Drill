package bosses

import "github.com/Kishlin/drill-game/internal/domain/types"

// Projectile represents a projectile entity that can be spawned by bosses
type Projectile struct {
	AABB     types.AABB
	Velocity types.Vec2
	Damage   float32
	Active   bool
}

func NewProjectile(x, y, width, height float32, velocity types.Vec2, damage float32) *Projectile {
	return &Projectile{
		AABB:     types.AABB{X: x, Y: y, Width: width, Height: height},
		Velocity: velocity,
		Damage:   damage,
		Active:   true,
	}
}

func (p *Projectile) Update(dt float32) {
	if p.Active == false {
		return
	}
	p.AABB.X += p.Velocity.X * dt
	p.AABB.Y += p.Velocity.Y * dt
}

func (p *Projectile) Intersects(aabb types.AABB) bool {
	return p.AABB.Intersects(aabb)
}

func (p *Projectile) Deactivate() {
	p.Active = false
}
