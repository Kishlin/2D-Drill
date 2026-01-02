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
// Now works with *Player directly instead of interface
func (ps *PhysicsSystem) UpdatePhysics(
	player *entities.Player,
	inputState input.InputState,
	dt float32,
) {
	// 1. Apply movement and gravity to velocity
	player.Velocity = physics.ApplyHorizontalMovement(player.Velocity, inputState, dt)
	player.Velocity = physics.ApplyVerticalMovement(player.Velocity, inputState, dt)
	player.Velocity = physics.ApplyGravity(player.Velocity, dt)

	// 2. AXIS-SEPARATED COLLISION RESOLUTION

	// X-axis: integrate position → check → resolve
	player.AABB.X += player.Velocity.X * dt
	collisionsX := physics.CheckCollisions(player.AABB, ps.world)
	player.AABB, player.Velocity = physics.ResolveCollisionsX(player.AABB, player.Velocity, collisionsX)

	// Y-axis: integrate position → check → resolve
	player.AABB.Y += player.Velocity.Y * dt
	collisionsY := physics.CheckCollisions(player.AABB, ps.world)
	player.AABB, player.Velocity, player.OnGround = physics.ResolveCollisionsY(player.AABB, player.Velocity, collisionsY)
}
