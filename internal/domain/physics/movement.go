package physics

import (
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// ApplyHorizontalMovement updates velocity based on horizontal input (left/right)
// Applies acceleration when input is active, damping when no input
// maxSpeed and acceleration are passed as parameters to support upgrades
func ApplyHorizontalMovement(velocity types.Vec2, inputState input.InputState, dt float32, maxSpeed, acceleration float32) types.Vec2 {
	vel := velocity

	if inputState.Right {
		// Accelerate to the right
		vel.X += acceleration * dt
		// Cap at max speed
		if vel.X > maxSpeed {
			vel.X = maxSpeed
		}
	} else if inputState.Left {
		// Accelerate to the left
		vel.X -= acceleration * dt
		// Cap at max speed (in negative direction)
		if vel.X < -maxSpeed {
			vel.X = -maxSpeed
		}
	} else {
		// No input: apply damping to gradually slow down
		if vel.X > 0 {
			// Moving right, apply damping to slow down
			vel.X -= MoveDamping * dt
			// Don't reverse direction; stop at 0
			if vel.X < 0 {
				vel.X = 0
			}
		} else if vel.X < 0 {
			// Moving left, apply damping to slow down
			vel.X += MoveDamping * dt
			// Don't reverse direction; stop at 0
			if vel.X > 0 {
				vel.X = 0
			}
		}
	}

	return vel
}

// ApplyVerticalMovement updates velocity based on vertical input (flying)
// Applies upward acceleration when input is active, damping when no input
// flyAcceleration and maxUpwardVelocity are passed as parameters to support upgrades
func ApplyVerticalMovement(velocity types.Vec2, inputState input.InputState, dt float32, flyAcceleration, maxUpwardVelocity float32) types.Vec2 {
	vel := velocity

	if inputState.Up {
		// Increase upward velocity (decrease Y velocity, which is negative)
		vel.Y -= flyAcceleration * dt
		// Cap the maximum upward speed
		if vel.Y < maxUpwardVelocity {
			vel.Y = maxUpwardVelocity
		}
	} else {
		// When key is released, apply damping to gradually slow down upward movement
		// and eventually fall. This creates the inertia effect.
		if vel.Y < 0 {
			// Player is moving upward, apply damping to slow the ascent
			vel.Y += FlyDamping * dt
			// Don't let damping reverse the direction; stop at 0
			if vel.Y > 0 {
				vel.Y = 0
			}
		}
		// If vel.Y >= 0, gravity will handle the falling
	}

	return vel
}
