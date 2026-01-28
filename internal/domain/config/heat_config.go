package config

// HeatConfig defines temperature and heat damage parameters
type HeatConfig struct {
	BaseTemperature float32 // Temperature at ground level (°C)
	MaxTemperature  float32 // Temperature at max depth (°C)
	DamageBaseDPS   float32 // Base damage per second when over resistance
	DamageDivisor   float32 // Scaling factor for excess heat
	DamageExponent  float32 // Exponential scaling factor
}
