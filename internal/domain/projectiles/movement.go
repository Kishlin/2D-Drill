package projectiles

import "github.com/Kishlin/drill-game/internal/domain/types"

// Movement defines how a projectile moves each frame
type Movement interface {
	Update(position types.Vec2, dt float32) types.Vec2
}

// Linear moves in a straight line at constant velocity
type Linear struct {
	Velocity types.Vec2
}

func (m Linear) Update(position types.Vec2, dt float32) types.Vec2 {
	return types.Vec2{
		X: position.X + m.Velocity.X*dt,
		Y: position.Y + m.Velocity.Y*dt,
	}
}

// Sinusoidal moves with a wave pattern perpendicular to base velocity
type Sinusoidal struct {
	BaseVelocity types.Vec2
	Amplitude    float32
	Frequency    float32
	elapsed      float32
}

func NewSinusoidal(velocity types.Vec2, amplitude, frequency float32) *Sinusoidal {
	return &Sinusoidal{
		BaseVelocity: velocity,
		Amplitude:    amplitude,
		Frequency:    frequency,
		elapsed:      0,
	}
}

func (m *Sinusoidal) Update(position types.Vec2, dt float32) types.Vec2 {
	m.elapsed += dt

	// Calculate perpendicular direction for wave
	// Perpendicular to (vx, vy) is (-vy, vx)
	mag := m.BaseVelocity.Magnitude()
	if mag == 0 {
		return position
	}

	perpX := -m.BaseVelocity.Y / mag
	perpY := m.BaseVelocity.X / mag

	// Sine wave offset
	offset := m.Amplitude * sin(m.Frequency*m.elapsed)

	return types.Vec2{
		X: position.X + m.BaseVelocity.X*dt + perpX*offset*dt,
		Y: position.Y + m.BaseVelocity.Y*dt + perpY*offset*dt,
	}
}

// Homing tracks toward a target position
type Homing struct {
	Speed  float32
	Target *types.Vec2 // Pointer so it can track moving target
}

func (m Homing) Update(position types.Vec2, dt float32) types.Vec2 {
	if m.Target == nil {
		return position
	}

	dx := m.Target.X - position.X
	dy := m.Target.Y - position.Y
	dir := types.NewVec2(dx, dy).Normalize()

	return types.Vec2{
		X: position.X + dir.X*m.Speed*dt,
		Y: position.Y + dir.Y*m.Speed*dt,
	}
}

// Orbital orbits around a center point
type Orbital struct {
	Center *types.Vec2 // Pointer so it can track moving center
	Radius float32
	Speed  float32 // Radians per second
	angle  float32
}

func NewOrbital(center *types.Vec2, radius, speed float32) *Orbital {
	return &Orbital{
		Center: center,
		Radius: radius,
		Speed:  speed,
		angle:  0,
	}
}

func (m *Orbital) Update(position types.Vec2, dt float32) types.Vec2 {
	if m.Center == nil {
		return position
	}

	m.angle += m.Speed * dt

	return types.Vec2{
		X: m.Center.X + m.Radius*cos(m.angle),
		Y: m.Center.Y + m.Radius*sin(m.angle),
	}
}

// Simple sin/cos using Taylor series approximation (avoid math import)
func sin(x float32) float32 {
	// Normalize to [-π, π]
	for x > 3.14159 {
		x -= 6.28318
	}
	for x < -3.14159 {
		x += 6.28318
	}
	// Taylor series: sin(x) ≈ x - x³/6 + x⁵/120
	x3 := x * x * x
	x5 := x3 * x * x
	return x - x3/6 + x5/120
}

func cos(x float32) float32 {
	return sin(x + 1.5708) // cos(x) = sin(x + π/2)
}
