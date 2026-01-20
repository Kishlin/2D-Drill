package movement

import "github.com/Kishlin/drill-game/internal/domain/types"

// MovementBehavior defines how a boss moves
type MovementBehavior interface {
	// Update updates the movement and returns the new position
	Update(currentPos types.Vec2, dt float32) types.Vec2

	// GetVelocity returns the current movement velocity
	GetVelocity() types.Vec2

	// SetSpeed changes the movement speed
	SetSpeed(speed float32)

	// GetSpeed returns the current speed
	GetSpeed() float32
}
