package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	Gravity     = 800.0 // Pixels per second squared
	GroundLevel = 600.0 // Y position of ground surface
)

type PhysicsSystem struct{}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

// UpdatePhysics applies gravity and updates position based on velocity
func (ps *PhysicsSystem) UpdatePhysics(player interface {
	GetPosition() *rl.Vector2
	GetVelocity() *rl.Vector2
	SetOnGround(bool)
	GetHeight() float32
}, dt float32) {
	pos := player.GetPosition()
	vel := player.GetVelocity()

	// Apply gravity (always)
	vel.Y += Gravity * dt

	// Update position based on velocity
	pos.X += vel.X * dt
	pos.Y += vel.Y * dt

	// Ground collision (simple: just check if below ground level)
	playerBottom := pos.Y + player.GetHeight()
	if playerBottom >= GroundLevel {
		pos.Y = GroundLevel - player.GetHeight() // Snap to ground
		vel.Y = 0                                // Stop vertical movement
		player.SetOnGround(true)
	} else {
		player.SetOnGround(false)
	}
}
