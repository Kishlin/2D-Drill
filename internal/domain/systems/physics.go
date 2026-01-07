package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/physics"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type PhysicsSystem struct {
	world *world.World
}

func NewPhysicsSystem(w *world.World) *PhysicsSystem {
	return &PhysicsSystem{world: w}
}

// UpdatePhysics applies physics using axis-separated AABB collision
func (ps *PhysicsSystem) UpdatePhysics(
	player *entities.Player,
	inputState input.InputState,
	dt float32,
) {
	physics.ApplyHeatDamage(player, dt)

	if player.IsDrilling {
		return
	}

	// 1. Apply movement and gravity to velocity
	player.Velocity = physics.ApplyHorizontalMovement(
		player.Velocity, inputState, dt,
		player.Engine.MaxSpeed(), player.Engine.Acceleration(),
	)
	player.Velocity = physics.ApplyVerticalMovement(
		player.Velocity, inputState, dt,
		player.Engine.FlyAcceleration(), player.Engine.MaxUpwardSpeed(),
	)
	player.Velocity = physics.ApplyGravity(player.Velocity, dt)

	// 2. AXIS-SEPARATED COLLISION RESOLUTION

	// X-axis: integrate position → check → resolve
	player.AABB.X += player.Velocity.X * dt
	collisionsX := physics.CheckCollisions(player.AABB, ps.world)
	player.AABB, player.Velocity = physics.ResolveCollisionsX(player.AABB, player.Velocity, collisionsX)

	// Y-axis: integrate position → check → resolve
	player.AABB.Y += player.Velocity.Y * dt
	collisionsY := physics.CheckCollisions(player.AABB, ps.world)

	// Capture state before Y-resolution for fall damage calculation
	wasAirborne := !player.OnGround
	ySpeedBeforeLanding := player.Velocity.Y

	player.AABB, player.Velocity, player.OnGround = physics.ResolveCollisionsY(player.AABB, player.Velocity, collisionsY)

	// Apply fall damage on landing transition
	if wasAirborne && player.OnGround {
		physics.ApplyFallDamage(player, ySpeedBeforeLanding)
	}

	// 3. Enforce world boundary constraints (prevent player from leaving game area)
	ps.constrainPlayerToWorldBounds(player)
}

// constrainPlayerToWorldBounds clamps player position to world boundaries
func (ps *PhysicsSystem) constrainPlayerToWorldBounds(player *entities.Player) {
	// Horizontal bounds: player cannot go left of 0 or right of (worldWidth - playerWidth)
	minX := float32(0.0)
	maxX := ps.world.Width - float32(entities.PlayerWidth)

	if player.AABB.X < minX {
		player.AABB.X = minX
		player.Velocity.X = 0 // Stop horizontal velocity at boundary
	} else if player.AABB.X > maxX {
		player.AABB.X = maxX
		player.Velocity.X = 0 // Stop horizontal velocity at boundary
	}

	// Vertical bounds: player cannot go above y=0
	minY := float32(0.0)

	if player.AABB.Y < minY {
		player.AABB.Y = minY
		player.Velocity.Y = 0 // Stop vertical velocity at top boundary
	}

	// No maximum Y - player can drill infinitely deep
}
