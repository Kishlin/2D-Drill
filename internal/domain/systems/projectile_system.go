package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// Projectile represents an active projectile in the pool
type Projectile struct {
	aabb     types.AABB
	movement projectiles.Movement
	damage   float32
	active   bool
}

func (p *Projectile) IsActive() bool {
	return p.active
}

func (p *Projectile) AABB() types.AABB {
	return p.aabb
}

// ProjectileBounds defines world bounds for culling
type ProjectileBounds struct {
	MinX, MaxX, MinY, MaxY float32
}

// CollisionTarget is anything a projectile can hit
type CollisionTarget interface {
	GetAABB() types.AABB
}

const DefaultPoolSize = 64

func NewProjectilePool() []Projectile {
	return make([]Projectile, DefaultPoolSize)
}

// SpawnProjectiles spawns projectiles from a slice of requests into the pool
func SpawnProjectiles(pool []Projectile, requests []projectiles.SpawnRequest) {
	for _, req := range requests {
		spawnProjectile(pool, req)
	}
}

func spawnProjectile(pool []Projectile, req projectiles.SpawnRequest) {
	for i := range pool {
		if pool[i].active == false {
			pool[i] = Projectile{
				aabb: types.AABB{
					X:      req.Position.X - req.Size/2,
					Y:      req.Position.Y - req.Size/2,
					Width:  req.Size,
					Height: req.Size,
				},
				movement: req.Movement,
				damage:   req.Damage,
				active:   true,
			}
			return
		}
	}
}

// UpdateProjectiles moves, culls, and checks collisions for all active projectiles
func UpdateProjectiles(pool []Projectile, bounds ProjectileBounds, dt float32, targets []CollisionTarget) []effects.Effect {
	var result []effects.Effect

	for i := range pool {
		p := &pool[i]
		if p.active == false {
			continue
		}

		// Move via movement behavior
		newPos := p.movement.Update(types.Vec2{X: p.aabb.X, Y: p.aabb.Y}, dt)
		p.aabb.X = newPos.X
		p.aabb.Y = newPos.Y

		// Cull out-of-bounds
		if isProjectileOutOfBounds(p.aabb, bounds) {
			p.active = false
			continue
		}

		// Check collisions
		for _, target := range targets {
			if p.aabb.Intersects(target.GetAABB()) {
				result = append(result, effects.ProjectileDamage{Damage: p.damage})
				p.active = false
				break
			}
		}
	}

	return result
}

func isProjectileOutOfBounds(aabb types.AABB, bounds ProjectileBounds) bool {
	return aabb.X < bounds.MinX || aabb.X > bounds.MaxX ||
		aabb.Y < bounds.MinY || aabb.Y > bounds.MaxY
}
