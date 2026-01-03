package physics

const (
	// Gravity constant (pixels per second squared)
	Gravity = 800.0

	// Horizontal movement constants
	MaxMoveSpeed     = 450.0  // Maximum horizontal velocity (pixels per second)
	MoveAcceleration = 2500.0 // Horizontal acceleration when key is held (pixels per second squared)
	MoveDamping      = 1000.0 // Horizontal deceleration when key is released (pixels per second squared)

	// Vertical movement (flying) constants
	FlyAcceleration   = 2500.0 // Upward acceleration when key is held (pixels per second squared)
	MaxUpwardVelocity = -600.0 // Maximum upward speed (negative Y is up)
	FlyDamping        = 300.0  // Deceleration when fly key is released (pixels per second squared)

	// Fall damage constants
	FallDamageThreshold = 500.0 // Minimum downward speed (px/sec) to deal damage
	FallDamageDivisor   = 20.0  // Damage scaling: (speed - threshold) / divisor
)
