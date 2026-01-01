package physics

import (
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// CollisionResult contains the results of collision detection and resolution
type CollisionResult struct {
	Position types.Vec2
	Velocity types.Vec2
	OnGround bool
}

// ResolveGroundCollision checks for and resolves collision with the ground
// Returns updated position, velocity, and ground contact state
func ResolveGroundCollision(position, velocity types.Vec2, height float32, w *world.World) CollisionResult {
	playerBottom := position.Y + height
	groundLevel := w.GetGroundLevel()

	if playerBottom >= groundLevel {
		// Snap to ground
		newPos := types.Vec2{
			X: position.X,
			Y: groundLevel - height,
		}

		// Stop downward velocity only (allow takeoff from ground)
		newVel := velocity
		if newVel.Y >= 0 {
			newVel.Y = 0
		}

		return CollisionResult{
			Position: newPos,
			Velocity: newVel,
			OnGround: true,
		}
	}

	return CollisionResult{
		Position: position,
		Velocity: velocity,
		OnGround: false,
	}
}
