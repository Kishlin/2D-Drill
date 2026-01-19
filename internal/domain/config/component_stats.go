package config

type EngineStats struct {
	MaxSpeed        float32
	Acceleration    float32
	FlyAcceleration float32
	MaxUpwardSpeed  float32
}

type HullStats struct {
	MaxHP float32
}

type FuelTankStats struct {
	Capacity float32
}

type CargoHoldStats struct {
	Capacity int
}

type HeatShieldStats struct {
	HeatResistance float32
}

type DrillStats struct {
	DrillSpeed float32
}
