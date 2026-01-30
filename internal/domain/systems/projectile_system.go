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

// ProjectileBounds defines world bounds for culling
type ProjectileBounds struct {
	MinX, MaxX, MinY, MaxY float32
}

// CollisionTarget is anything a projectile can hit
type CollisionTarget interface {
	GetAABB() types.AABB
}

// ProjectileRenderData is read-only data for rendering
type ProjectileRenderData struct {
	AABB types.AABB
}

// ProjectileSystem manages all projectiles centrally
type ProjectileSystem struct {
	pool         []Projectile
	bounds       ProjectileBounds
	renderBuffer []ProjectileRenderData
}

const DefaultPoolSize = 64

func NewProjectileSystem(bounds ProjectileBounds) *ProjectileSystem {
	return &ProjectileSystem{
		pool:         make([]Projectile, DefaultPoolSize),
		bounds:       bounds,
		renderBuffer: make([]ProjectileRenderData, 0, DefaultPoolSize),
	}
}

// SpawnAll spawns projectiles from a slice of requests
func (ps *ProjectileSystem) SpawnAll(requests []projectiles.SpawnRequest) {
	for _, req := range requests {
		ps.spawn(req)
	}
}

func (ps *ProjectileSystem) spawn(req projectiles.SpawnRequest) bool {
	for i := range ps.pool {
		if ps.pool[i].active == false {
			ps.pool[i] = Projectile{
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
			return true
		}
	}
	return false // Pool exhausted
}

func (ps *ProjectileSystem) Update(dt float32, targets []CollisionTarget) []effects.Effect {
	var result []effects.Effect

	for i := range ps.pool {
		p := &ps.pool[i]
		if p.active == false {
			continue
		}

		// Move via movement behavior
		newPos := p.movement.Update(types.Vec2{X: p.aabb.X, Y: p.aabb.Y}, dt)
		p.aabb.X = newPos.X
		p.aabb.Y = newPos.Y

		// Cull out-of-bounds
		if ps.isOutOfBounds(p.aabb) {
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

func (ps *ProjectileSystem) isOutOfBounds(aabb types.AABB) bool {
	return aabb.X < ps.bounds.MinX || aabb.X > ps.bounds.MaxX ||
		aabb.Y < ps.bounds.MinY || aabb.Y > ps.bounds.MaxY
}

func (ps *ProjectileSystem) GetActiveProjectiles() []ProjectileRenderData {
	ps.renderBuffer = ps.renderBuffer[:0] // Reuse backing array
	for i := range ps.pool {
		if ps.pool[i].active {
			ps.renderBuffer = append(ps.renderBuffer, ProjectileRenderData{
				AABB: ps.pool[i].aabb,
			})
		}
	}
	return ps.renderBuffer
}

func (ps *ProjectileSystem) Clear() {
	for i := range ps.pool {
		ps.pool[i].active = false
	}
}
