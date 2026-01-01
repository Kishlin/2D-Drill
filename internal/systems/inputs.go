package systems

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	MoveSpeed    = 200.0 // Pixels per second (horizontal)
	JumpStrength = 400.0 // Initial upward velocity (pixels per second)
)

type InputSystem struct{}

func NewInputSystem() *InputSystem {
	return &InputSystem{}
}

// UpdatePlayerInput modifies player velocity based on keyboard input
func (is *InputSystem) UpdatePlayerInput(player interface {
	GetVelocity() *rl.Vector2
	IsOnGround() bool
}, dt float32) {
	vel := player.GetVelocity()

	// Horizontal movement (instant, not acceleration-based for now)
	vel.X = 0
	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		vel.X = MoveSpeed
	}
	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		vel.X = -MoveSpeed
	}

	// Jump (only when on ground)
	if (rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW)) && player.IsOnGround() {
		vel.Y = -JumpStrength // Negative Y is up in Raylib
	}
}
