package systems

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	MaxMoveSpeed      = 300.0  // Pixels per second (horizontal max velocity)
	MoveAcceleration  = 2500.0 // Pixels per second squared (horizontal acceleration when key is held)
	MoveDamping       = 1000.0 // Pixels per second squared (deceleration when key is released)
	FlyAcceleration   = 2000.0 // Pixels per second squared (upward acceleration when key is held)
	MaxUpwardVelocity = -300.0 // Maximum upward speed (negative Y is up)
	FlyDamping        = 300.0  // Pixels per second squared (deceleration when key is released)
)

type InputSystem struct{}

func NewInputSystem() *InputSystem {
	return &InputSystem{}
}

// UpdatePlayerInput modifies player velocity based on keyboard input
func (is *InputSystem) UpdatePlayerInput(player interface {
	GetVelocity() *rl.Vector2
}, dt float32) {
	vel := player.GetVelocity()

	// Horizontal movement with inertia: acceleration and damping
	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		// Accelerate to the right
		vel.X += MoveAcceleration * dt
		// Cap at max speed
		if vel.X > MaxMoveSpeed {
			vel.X = MaxMoveSpeed
		}
	} else if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		// Accelerate to the left
		vel.X -= MoveAcceleration * dt
		// Cap at max speed (in negative direction)
		if vel.X < -MaxMoveSpeed {
			vel.X = -MaxMoveSpeed
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

	// Flying mechanics: gradually accelerate upward when key is held
	if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		// Increase upward velocity (decrease Y velocity, which is negative)
		vel.Y -= FlyAcceleration * dt
		// Cap the maximum upward speed
		if vel.Y < MaxUpwardVelocity {
			vel.Y = MaxUpwardVelocity
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
		// If vel.Y >= 0, gravity will handle the falling (in physics system)
	}
}
