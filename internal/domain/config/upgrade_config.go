package config

type UpgradeTier[T any] struct {
	Name  string
	Price int
	Stats T
}

type UpgradeConfig struct {
	Engines     []UpgradeTier[EngineStats]
	Hulls       []UpgradeTier[HullStats]
	FuelTanks   []UpgradeTier[FuelTankStats]
	CargoHolds  []UpgradeTier[CargoHoldStats]
	HeatShields []UpgradeTier[HeatShieldStats]
	Drills      []UpgradeTier[DrillStats]
}
