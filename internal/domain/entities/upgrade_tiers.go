package entities

type EngineTier struct {
	Name            string
	MaxSpeed        float32 // horizontal max speed (pixels per second)
	Acceleration    float32 // horizontal acceleration (pixels per second squared)
	FlyAcceleration float32 // upward acceleration (pixels per second squared)
	MaxUpwardSpeed  float32 // max upward velocity (negative = upward, pixels per second)
	Cost            int
}

var EngineTiers = []EngineTier{
	{Name: "Base", MaxSpeed: 450.0, Acceleration: 2500.0, FlyAcceleration: 2500.0, MaxUpwardSpeed: -600.0, Cost: 0},
	{Name: "Engine Mk1", MaxSpeed: 500.0, Acceleration: 2750.0, FlyAcceleration: 2750.0, MaxUpwardSpeed: -650.0, Cost: 100},
	{Name: "Engine Mk2", MaxSpeed: 550.0, Acceleration: 3000.0, FlyAcceleration: 3000.0, MaxUpwardSpeed: -700.0, Cost: 300},
	{Name: "Engine Mk3", MaxSpeed: 600.0, Acceleration: 3500.0, FlyAcceleration: 3500.0, MaxUpwardSpeed: -775.0, Cost: 750},
	{Name: "Engine Mk4", MaxSpeed: 675.0, Acceleration: 4000.0, FlyAcceleration: 4000.0, MaxUpwardSpeed: -850.0, Cost: 1500},
	{Name: "Engine Mk5", MaxSpeed: 750.0, Acceleration: 5000.0, FlyAcceleration: 5000.0, MaxUpwardSpeed: -1000.0, Cost: 5000},
}

type HullTier struct {
	Name  string
	MaxHP float32
	Cost  int
}

var HullTiers = []HullTier{
	{Name: "Base", MaxHP: 10.0, Cost: 0},
	{Name: "Hull Mk1", MaxHP: 15.0, Cost: 150},
	{Name: "Hull Mk2", MaxHP: 20.0, Cost: 400},
	{Name: "Hull Mk3", MaxHP: 30.0, Cost: 1000},
	{Name: "Hull Mk4", MaxHP: 45.0, Cost: 2500},
	{Name: "Hull Mk5", MaxHP: 75.0, Cost: 8000},
}

type FuelTankTier struct {
	Name     string
	Capacity float32 // liters
	Cost     int
}

var FuelTankTiers = []FuelTankTier{
	{Name: "Base", Capacity: 10.0, Cost: 0},
	{Name: "Tank Mk1", Capacity: 15.0, Cost: 100},
	{Name: "Tank Mk2", Capacity: 22.0, Cost: 250},
	{Name: "Tank Mk3", Capacity: 32.0, Cost: 600},
	{Name: "Tank Mk4", Capacity: 45.0, Cost: 1500},
	{Name: "Tank Mk5", Capacity: 65.0, Cost: 4000},
}

func GetEngineMaxSpeed(level UpgradeLevel) float32 {
	if int(level) >= len(EngineTiers) {
		return EngineTiers[len(EngineTiers)-1].MaxSpeed
	}
	return EngineTiers[level].MaxSpeed
}

func GetEngineAcceleration(level UpgradeLevel) float32 {
	if int(level) >= len(EngineTiers) {
		return EngineTiers[len(EngineTiers)-1].Acceleration
	}
	return EngineTiers[level].Acceleration
}

func GetEngineFlyAcceleration(level UpgradeLevel) float32 {
	if int(level) >= len(EngineTiers) {
		return EngineTiers[len(EngineTiers)-1].FlyAcceleration
	}
	return EngineTiers[level].FlyAcceleration
}

func GetEngineMaxUpwardSpeed(level UpgradeLevel) float32 {
	if int(level) >= len(EngineTiers) {
		return EngineTiers[len(EngineTiers)-1].MaxUpwardSpeed
	}
	return EngineTiers[level].MaxUpwardSpeed
}

// GetEngineNextCost Returns 0 if already at max level
func GetEngineNextCost(level UpgradeLevel) int {
	nextLevel := int(level) + 1
	if nextLevel >= len(EngineTiers) {
		return 0 // Max level reached
	}
	return EngineTiers[nextLevel].Cost
}

func GetHullMaxHP(level UpgradeLevel) float32 {
	if int(level) >= len(HullTiers) {
		return HullTiers[len(HullTiers)-1].MaxHP
	}
	return HullTiers[level].MaxHP
}

// GetHullNextCost Returns 0 if already at max level
func GetHullNextCost(level UpgradeLevel) int {
	nextLevel := int(level) + 1
	if nextLevel >= len(HullTiers) {
		return 0 // Max level reached
	}
	return HullTiers[nextLevel].Cost
}

func GetFuelTankCapacity(level UpgradeLevel) float32 {
	if int(level) >= len(FuelTankTiers) {
		return FuelTankTiers[len(FuelTankTiers)-1].Capacity
	}
	return FuelTankTiers[level].Capacity
}

// GetFuelTankNextCost Returns 0 if already at max level
func GetFuelTankNextCost(level UpgradeLevel) int {
	nextLevel := int(level) + 1
	if nextLevel >= len(FuelTankTiers) {
		return 0 // Max level reached
	}
	return FuelTankTiers[nextLevel].Cost
}
