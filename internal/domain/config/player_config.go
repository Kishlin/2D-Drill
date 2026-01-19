package config

// StartingUpgrades defines which tier each component starts at (0 = Base)
type StartingUpgrades struct {
	Engine     int
	Hull       int
	FuelTank   int
	CargoHold  int
	HeatShield int
	Drill      int
}

type PlayerConfig struct {
	StartingMoney    int
	StartingItems    [5]int           // Teleport, Repair, Refuel, Bomb, BigBomb
	StartingUpgrades StartingUpgrades // Tier indices (0=Base, 1=Mk1, etc.)
}

func DefaultPlayerConfig() PlayerConfig {
	return PlayerConfig{
		StartingMoney: 100000,
		StartingItems: [5]int{5, 5, 5, 5, 5}, // 5 of each item for testing
		StartingUpgrades: StartingUpgrades{
			Engine:     0,
			Hull:       0,
			FuelTank:   0,
			CargoHold:  0,
			HeatShield: 0,
			Drill:      0,
		},
	}
}
