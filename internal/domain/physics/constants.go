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

	// Heat system constants
	GroundLevelY      = 640.0   // Ground level position (pixels)
	MaxUndergroundY   = 64000.0 // Maximum underground depth (pixels)
	BaseTemperature   = 15.0    // Temperature at ground level (°C)
	MaxTemperature    = 350.0   // Temperature at max depth (°C)

	// Heat damage constants
	HeatDamageBaseDPS  = 0.5  // Base damage per second
	HeatDamageDivisor  = 10.0 // Scaling factor for excess heat
	HeatDamageExponent = 1.5  // Exponential scaling factor
)
