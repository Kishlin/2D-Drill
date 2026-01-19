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

func DefaultUpgradeConfig() UpgradeConfig {
	return UpgradeConfig{
		Engines: []UpgradeTier[EngineStats]{
			{Name: "Base Engine", Price: 0, Stats: EngineStats{MaxSpeed: 450.0, Acceleration: 2500.0, FlyAcceleration: 2500.0, MaxUpwardSpeed: -600.0}},
			{Name: "Engine Mk1", Price: 100, Stats: EngineStats{MaxSpeed: 475.0, Acceleration: 2667.0, FlyAcceleration: 2667.0, MaxUpwardSpeed: -635.0}},
			{Name: "Engine Mk2", Price: 300, Stats: EngineStats{MaxSpeed: 500.0, Acceleration: 2833.0, FlyAcceleration: 2833.0, MaxUpwardSpeed: -670.0}},
			{Name: "Engine Mk3", Price: 750, Stats: EngineStats{MaxSpeed: 525.0, Acceleration: 3000.0, FlyAcceleration: 3000.0, MaxUpwardSpeed: -705.0}},
			{Name: "Engine Mk4", Price: 1500, Stats: EngineStats{MaxSpeed: 562.0, Acceleration: 3250.0, FlyAcceleration: 3250.0, MaxUpwardSpeed: -740.0}},
			{Name: "Engine Mk5", Price: 5000, Stats: EngineStats{MaxSpeed: 600.0, Acceleration: 3500.0, FlyAcceleration: 3500.0, MaxUpwardSpeed: -775.0}},
		},
		Hulls: []UpgradeTier[HullStats]{
			{Name: "Base Hull", Price: 0, Stats: HullStats{MaxHP: 10.0}},
			{Name: "Hull Mk1", Price: 150, Stats: HullStats{MaxHP: 15.0}},
			{Name: "Hull Mk2", Price: 400, Stats: HullStats{MaxHP: 20.0}},
			{Name: "Hull Mk3", Price: 1000, Stats: HullStats{MaxHP: 30.0}},
			{Name: "Hull Mk4", Price: 2500, Stats: HullStats{MaxHP: 45.0}},
			{Name: "Hull Mk5", Price: 8000, Stats: HullStats{MaxHP: 75.0}},
		},
		FuelTanks: []UpgradeTier[FuelTankStats]{
			{Name: "Base Tank", Price: 0, Stats: FuelTankStats{Capacity: 10.0}},
			{Name: "Tank Mk1", Price: 100, Stats: FuelTankStats{Capacity: 15.0}},
			{Name: "Tank Mk2", Price: 250, Stats: FuelTankStats{Capacity: 22.0}},
			{Name: "Tank Mk3", Price: 600, Stats: FuelTankStats{Capacity: 32.0}},
			{Name: "Tank Mk4", Price: 1500, Stats: FuelTankStats{Capacity: 45.0}},
			{Name: "Tank Mk5", Price: 4000, Stats: FuelTankStats{Capacity: 65.0}},
		},
		CargoHolds: []UpgradeTier[CargoHoldStats]{
			{Name: "Base Cargo Hold", Price: 0, Stats: CargoHoldStats{Capacity: 10}},
			{Name: "Cargo Hold Mk1", Price: 125, Stats: CargoHoldStats{Capacity: 14}},
			{Name: "Cargo Hold Mk2", Price: 350, Stats: CargoHoldStats{Capacity: 18}},
			{Name: "Cargo Hold Mk3", Price: 800, Stats: CargoHoldStats{Capacity: 24}},
			{Name: "Cargo Hold Mk4", Price: 2000, Stats: CargoHoldStats{Capacity: 31}},
			{Name: "Cargo Hold Mk5", Price: 6000, Stats: CargoHoldStats{Capacity: 40}},
		},
		HeatShields: []UpgradeTier[HeatShieldStats]{
			{Name: "Base Heat Shield", Price: 0, Stats: HeatShieldStats{HeatResistance: 50.0}},
			{Name: "Heat Shield Mk1", Price: 200, Stats: HeatShieldStats{HeatResistance: 90.0}},
			{Name: "Heat Shield Mk2", Price: 500, Stats: HeatShieldStats{HeatResistance: 140.0}},
			{Name: "Heat Shield Mk3", Price: 1200, Stats: HeatShieldStats{HeatResistance: 190.0}},
			{Name: "Heat Shield Mk4", Price: 3000, Stats: HeatShieldStats{HeatResistance: 250.0}},
			{Name: "Heat Shield Mk5", Price: 7500, Stats: HeatShieldStats{HeatResistance: 320.0}},
		},
		Drills: []UpgradeTier[DrillStats]{
			{Name: "Base Drill", Price: 0, Stats: DrillStats{DrillSpeed: 1.0}},
			{Name: "Drill Mk1", Price: 125, Stats: DrillStats{DrillSpeed: 2.0}},
			{Name: "Drill Mk2", Price: 350, Stats: DrillStats{DrillSpeed: 3.0}},
			{Name: "Drill Mk3", Price: 875, Stats: DrillStats{DrillSpeed: 4.0}},
			{Name: "Drill Mk4", Price: 2000, Stats: DrillStats{DrillSpeed: 5.0}},
			{Name: "Drill Mk5", Price: 6500, Stats: DrillStats{DrillSpeed: 6.0}},
		},
	}
}
