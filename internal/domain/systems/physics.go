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

// UpdatePhysics applies all physics: movement, gravity, integration, collision
// This method orchestrates pure physics functions from the physics package
func (ps *PhysicsSystem) UpdatePhysics(
	entity entities.PhysicsEntity,
	inputState input.InputState,
	dt float32,
) {
	// Get current position and velocity
	position := *entity.GetPositionVec()
	velocity := *entity.GetVelocityVec()

	// Apply physics using pure functions from physics package
	velocity = physics.ApplyHorizontalMovement(velocity, inputState, dt)
	velocity = physics.ApplyVerticalMovement(velocity, inputState, dt)
	velocity = physics.ApplyGravity(velocity, dt)
	position = physics.IntegrateVelocity(position, velocity, dt)

	// Resolve collisions
	result := physics.ResolveGroundCollision(position, velocity, entity.GetHeight(), ps.world)

	// Update entity with new state
	entity.SetPosition(result.Position)
	entity.SetVelocity(result.Velocity)
	entity.SetOnGround(result.OnGround)
}
