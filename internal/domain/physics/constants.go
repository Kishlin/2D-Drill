package physics

const (
	// Gravity constant (pixels per second squared)
	Gravity = 800.0

	// Damping constants (deceleration when input is released)
	MoveDamping = 1000.0 // Horizontal deceleration (pixels per second squared)
	FlyDamping  = 300.0  // Vertical deceleration when fly key is released (pixels per second squared)

	// Fall damage constants
	FallDamageThreshold = 500.0 // Minimum downward speed (px/sec) to deal damage
	FallDamageDivisor   = 20.0  // Damage scaling: (speed - threshold) / divisor
)
