package config

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
