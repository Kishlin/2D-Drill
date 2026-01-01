package entities

import "github.com/Kishlin/drill-game/internal/domain/types"

// PhysicsEntity defines the interface that objects needing physics simulation must implement
// This allows the physics system to work with any entity that supports movement and collision
type PhysicsEntity interface {
	// Position and velocity accessors
	GetPositionVec() *types.Vec2
	GetVelocityVec() *types.Vec2

	// Position and velocity setters
	SetPosition(types.Vec2)
	SetVelocity(types.Vec2)

	// Collision state
	SetOnGround(bool)

	// Dimensions
	GetHeight() float32
}
