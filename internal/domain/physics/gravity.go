package physics

import "github.com/Kishlin/drill-game/internal/domain/types"

// ApplyGravity applies gravitational acceleration to velocity
func ApplyGravity(velocity types.Vec2, dt float32) types.Vec2 {
	return types.Vec2{
		X: velocity.X,
		Y: velocity.Y + Gravity*dt,
	}
}

// IntegrateVelocity updates position based on current velocity
func IntegrateVelocity(position, velocity types.Vec2, dt float32) types.Vec2 {
	return types.Vec2{
		X: position.X + velocity.X*dt,
		Y: position.Y + velocity.Y*dt,
	}
}
