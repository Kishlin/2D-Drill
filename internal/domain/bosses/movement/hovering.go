package movement

import "github.com/Kishlin/drill-game/internal/domain/types"

// HoveringConfig holds configuration for hovering movement
type HoveringConfig struct {
	Speed        float32 // Horizontal movement speed in pixels per second
	MinX         float32 // Left boundary
	MaxX         float32 // Right boundary
	HoverY       float32 // Base Y position to hover at
	BossWidth    float32 // Width of the boss (for boundary calculation)
	BobAmplitude float32 // Vertical bobbing amplitude in pixels
	BobFrequency float32 // Vertical bobbing frequency in radians per second
}

// Hovering implements horizontal patrol with vertical sine-wave bobbing
type Hovering struct {
	config    HoveringConfig
	direction float32 // 1 for right, -1 for left
	elapsed   float32
	paused    bool
}

// NewHovering creates a new hovering movement behavior
func NewHovering(cfg HoveringConfig) *Hovering {
	return &Hovering{
		config:    cfg,
		direction: 1, // Start moving right
		elapsed:   0,
		paused:    false,
	}
}

// Update updates position with horizontal patrol and vertical bobbing
func (h *Hovering) Update(currentPos types.Vec2, dt float32) types.Vec2 {
	h.elapsed += dt

	if h.paused {
		// When paused, still bob vertically but don't move horizontally
		bobOffset := h.config.BobAmplitude * sin(h.config.BobFrequency*h.elapsed)
		return types.NewVec2(currentPos.X, h.config.HoverY+bobOffset)
	}

	// Calculate new X position
	newX := currentPos.X + h.config.Speed*h.direction*dt

	// Check boundaries and reverse direction if needed
	if newX < h.config.MinX {
		newX = h.config.MinX
		h.reverseDirection()
	}
	if newX+h.config.BossWidth > h.config.MaxX {
		newX = h.config.MaxX - h.config.BossWidth
		h.reverseDirection()
	}

	// Apply vertical bobbing
	bobOffset := h.config.BobAmplitude * sin(h.config.BobFrequency*h.elapsed)
	newY := h.config.HoverY + bobOffset

	return types.NewVec2(newX, newY)
}

// reverseDirection flips the movement direction
func (h *Hovering) reverseDirection() {
	h.direction *= -1
}

// SetSpeed updates the movement speed (preserves direction)
func (h *Hovering) SetSpeed(speed float32) {
	h.config.Speed = speed
}

// GetSpeed returns the current speed
func (h *Hovering) GetSpeed() float32 {
	return h.config.Speed
}

// GetDirection returns the current direction (1 = right, -1 = left)
func (h *Hovering) GetDirection() float32 {
	return h.direction
}

// Pause stops horizontal movement but continues bobbing
func (h *Hovering) Pause() {
	h.paused = true
}

// Resume restores horizontal movement
func (h *Hovering) Resume() {
	h.paused = false
}

// IsPaused returns whether movement is paused
func (h *Hovering) IsPaused() bool {
	return h.paused
}

// sin uses Taylor series approximation (matching the projectiles package pattern)
func sin(x float32) float32 {
	for x > 3.14159 {
		x -= 6.28318
	}
	for x < -3.14159 {
		x += 6.28318
	}
	x3 := x * x * x
	x5 := x3 * x * x
	return x - x3/6 + x5/120
}
