package config

// FuelSystemConfig defines fuel consumption rates
type FuelSystemConfig struct {
	ConsumptionMoving float32 // Fuel consumption rate when moving/drilling (L/s)
	ConsumptionIdle   float32 // Fuel consumption rate when idle (L/s)
}
