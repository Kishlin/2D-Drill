package systems

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

func TestFuelStationSystem_ProcessRefueling_FullTank(t *testing.T) {
	// Setup: Player with full tank
	player := testFuelStationPlayer(100, 100)
	player.Money = 100
	// player.Fuel is already at capacity from NewPlayerFromConfig

	fuelStation := entities.NewFuelStation(80, 80)
	system := NewFuelStationSystem(fuelStation)

	inputState := input.InputState{Sell: true}

	initialMoney := player.Money
	fuelCapacity := player.FuelTank.Capacity()

	// Execute
	system.ProcessRefueling(player, inputState)

	// Verify: No money deducted, fuel stays full
	if player.Money != initialMoney {
		t.Errorf("Expected money %d, got %d", initialMoney, player.Money)
	}
	if player.Fuel != fuelCapacity {
		t.Errorf("Expected fuel %.2f, got %.2f", fuelCapacity, player.Fuel)
	}
}

func TestFuelStationSystem_ProcessRefueling_EmptyTank(t *testing.T) {
	// Setup: Player with empty tank
	player := testFuelStationPlayer(100, 100)
	player.Money = 100
	player.Fuel = 0.0

	fuelStation := entities.NewFuelStation(80, 80)
	system := NewFuelStationSystem(fuelStation)

	inputState := input.InputState{Sell: true}

	fuelCapacity := player.FuelTank.Capacity()

	// Execute
	system.ProcessRefueling(player, inputState)

	// Verify: 10 money deducted (10 liters * $1), fuel full
	expectedMoney := 100 - int(fuelCapacity)
	if player.Money != expectedMoney {
		t.Errorf("Expected money %d, got %d", expectedMoney, player.Money)
	}
	if player.Fuel != fuelCapacity {
		t.Errorf("Expected fuel %.2f, got %.2f", fuelCapacity, player.Fuel)
	}
}

func TestFuelStationSystem_ProcessRefueling_PartialTankRoundedUp(t *testing.T) {
	// Setup: Player with partial tank (3.2 liters needed = $4 cost)
	player := testFuelStationPlayer(100, 100)
	player.Money = 100
	player.Fuel = 6.8 // Need 3.2 liters

	fuelStation := entities.NewFuelStation(80, 80)
	system := NewFuelStationSystem(fuelStation)

	inputState := input.InputState{Sell: true}

	fuelCapacity := player.FuelTank.Capacity()

	// Execute
	system.ProcessRefueling(player, inputState)

	// Verify: 4 money deducted (ceil(3.2) = 4), fuel full
	expectedMoney := 100 - 4
	if player.Money != expectedMoney {
		t.Errorf("Expected money %d, got %d", expectedMoney, player.Money)
	}
	if player.Fuel != fuelCapacity {
		t.Errorf("Expected fuel %.2f, got %.2f", fuelCapacity, player.Fuel)
	}
}

func TestFuelStationSystem_ProcessRefueling_InsufficientMoney(t *testing.T) {
	// Setup: Player with empty tank but insufficient money
	player := testFuelStationPlayer(100, 100)
	player.Money = 5 // Need 10, only have 5
	player.Fuel = 0.0

	fuelStation := entities.NewFuelStation(80, 80)
	system := NewFuelStationSystem(fuelStation)

	inputState := input.InputState{Sell: true}

	// Execute
	system.ProcessRefueling(player, inputState)

	// Verify: No transaction (money and fuel unchanged)
	if player.Money != 5 {
		t.Errorf("Expected money 5, got %d", player.Money)
	}
	if player.Fuel != 0.0 {
		t.Errorf("Expected fuel 0.0, got %.2f", player.Fuel)
	}
}

func TestFuelStationSystem_ProcessRefueling_NoInput(t *testing.T) {
	// Setup: Player in range but no Sell input
	player := testFuelStationPlayer(100, 100)
	player.Money = 100
	player.Fuel = 0.0

	fuelStation := entities.NewFuelStation(80, 80)
	system := NewFuelStationSystem(fuelStation)

	inputState := input.InputState{Sell: false}

	// Execute
	system.ProcessRefueling(player, inputState)

	// Verify: No transaction
	if player.Money != 100 {
		t.Errorf("Expected money 100, got %d", player.Money)
	}
	if player.Fuel != 0.0 {
		t.Errorf("Expected fuel 0.0, got %.2f", player.Fuel)
	}
}

func TestFuelStationSystem_ProcessRefueling_OutOfRange(t *testing.T) {
	// Setup: Player far from fuel station
	player := testFuelStationPlayer(500, 500)
	player.Money = 100
	player.Fuel = 0.0

	fuelStation := entities.NewFuelStation(80, 80)
	system := NewFuelStationSystem(fuelStation)

	inputState := input.InputState{Sell: true}

	// Execute
	system.ProcessRefueling(player, inputState)

	// Verify: No transaction
	if player.Money != 100 {
		t.Errorf("Expected money 100, got %d", player.Money)
	}
	if player.Fuel != 0.0 {
		t.Errorf("Expected fuel 0.0, got %.2f", player.Fuel)
	}
}

// Test helpers

func testFuelStationPlayer(x, y float32) *entities.Player {
	playerCfg := config.PlayerConfig{
		StartingMoney:    0,
		StartingItems:    [5]int{0, 0, 0, 0, 0},
		StartingUpgrades: config.StartingUpgrades{},
	}
	upgradeCfg := config.UpgradeConfig{
		Engines:     []config.UpgradeTier[config.EngineStats]{{Name: "Base", Price: 0, Stats: config.EngineStats{MaxSpeed: 450, Acceleration: 2500, FlyAcceleration: 2500, MaxUpwardSpeed: -600}}},
		Hulls:       []config.UpgradeTier[config.HullStats]{{Name: "Base", Price: 0, Stats: config.HullStats{MaxHP: 10}}},
		FuelTanks:   []config.UpgradeTier[config.FuelTankStats]{{Name: "Base", Price: 0, Stats: config.FuelTankStats{Capacity: 10}}},
		CargoHolds:  []config.UpgradeTier[config.CargoHoldStats]{{Name: "Base", Price: 0, Stats: config.CargoHoldStats{Capacity: 10}}},
		HeatShields: []config.UpgradeTier[config.HeatShieldStats]{{Name: "Base", Price: 0, Stats: config.HeatShieldStats{HeatResistance: 50}}},
		Drills:      []config.UpgradeTier[config.DrillStats]{{Name: "Base", Price: 0, Stats: config.DrillStats{DrillSpeed: 1.0}}},
	}
	return entities.NewPlayerFromConfig(x, y, playerCfg, upgradeCfg)
}
