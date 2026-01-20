package movement

import "github.com/Kishlin/drill-game/internal/domain/types"

// GroundedConfig holds configuration for grounded movement
type GroundedConfig struct {
	Speed     float32 // Movement speed in pixels per second
	MinX      float32 // Left boundary
	MaxX      float32 // Right boundary
	FloorY    float32 // Y position of the floor (boss bottom will be here)
	BossWidth float32 // Width of the boss (for boundary calculation)
}

// Grounded implements left-right patrol movement along the floor
type Grounded struct {
	config    GroundedConfig
	velocity  types.Vec2
	direction float32 // 1 for right, -1 for left
}

// NewGrounded creates a new grounded movement behavior
func NewGrounded(cfg GroundedConfig) *Grounded {
	return &Grounded{
		config:    cfg,
		velocity:  types.NewVec2(cfg.Speed, 0),
		direction: 1, // Start moving right
	}
}

// Update updates position and handles boundary reversal
func (g *Grounded) Update(currentPos types.Vec2, dt float32) types.Vec2 {
	// Calculate new X position
	newX := currentPos.X + g.velocity.X*dt

	// Check boundaries and reverse direction if needed
	// Left boundary: boss left edge hits MinX
	if newX < g.config.MinX {
		newX = g.config.MinX
		g.reverseDirection()
	}
	// Right boundary: boss right edge hits MaxX
	if newX+g.config.BossWidth > g.config.MaxX {
		newX = g.config.MaxX - g.config.BossWidth
		g.reverseDirection()
	}

	return types.NewVec2(newX, currentPos.Y)
}

// reverseDirection flips the movement direction
func (g *Grounded) reverseDirection() {
	g.direction *= -1
	g.velocity.X = g.config.Speed * g.direction
}

// GetVelocity returns current velocity
func (g *Grounded) GetVelocity() types.Vec2 {
	return g.velocity
}

// SetSpeed updates the movement speed (preserves direction)
func (g *Grounded) SetSpeed(speed float32) {
	g.config.Speed = speed
	g.velocity.X = speed * g.direction
}

// GetSpeed returns the current speed
func (g *Grounded) GetSpeed() float32 {
	return g.config.Speed
}

// GetDirection returns the current direction (1 = right, -1 = left)
func (g *Grounded) GetDirection() float32 {
	return g.direction
}
