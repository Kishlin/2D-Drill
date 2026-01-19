package systems

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

func TestHospitalSystem_ProcessHealing_FullHP(t *testing.T) {
	// Setup: Player with full HP
	player := testHospitalPlayer(100, 100)
	player.Money = 100
	// player.HP is already at max from NewPlayerFromConfig

	hospital := entities.NewHospital(80, 80)
	system := NewHospitalSystem(hospital)

	inputState := input.InputState{Sell: true}

	initialMoney := player.Money
	maxHP := player.Hull.MaxHP()

	// Execute
	system.ProcessHealing(player, inputState)

	// Verify: No money deducted, HP stays full
	if player.Money != initialMoney {
		t.Errorf("Expected money %d, got %d", initialMoney, player.Money)
	}
	if player.HP != maxHP {
		t.Errorf("Expected HP %.2f, got %.2f", maxHP, player.HP)
	}
}

func TestHospitalSystem_ProcessHealing_ZeroHP(t *testing.T) {
	// Setup: Player with zero HP
	player := testHospitalPlayer(100, 100)
	player.Money = 100
	player.HP = 0.0

	hospital := entities.NewHospital(80, 80)
	system := NewHospitalSystem(hospital)

	inputState := input.InputState{Sell: true}

	maxHP := player.Hull.MaxHP()

	// Execute
	system.ProcessHealing(player, inputState)

	// Verify: 20 money deducted (10 HP * $2), HP restored to max
	expectedMoney := 100 - int(maxHP)*2
	if player.Money != expectedMoney {
		t.Errorf("Expected money %d, got %d", expectedMoney, player.Money)
	}
	if player.HP != maxHP {
		t.Errorf("Expected HP %.2f, got %.2f", maxHP, player.HP)
	}
}

func TestHospitalSystem_ProcessHealing_PartialHPRoundedUp(t *testing.T) {
	// Setup: Player at 7.2 HP (need 2.8 HP = ceil(5.6) = $6)
	player := testHospitalPlayer(100, 100)
	player.Money = 100
	player.HP = 7.2

	hospital := entities.NewHospital(80, 80)
	system := NewHospitalSystem(hospital)

	inputState := input.InputState{Sell: true}

	maxHP := player.Hull.MaxHP()

	// Execute
	system.ProcessHealing(player, inputState)

	// Verify: 6 money deducted (ceil(2.8 * 2) = ceil(5.6) = 6), HP full
	expectedMoney := 100 - 6
	if player.Money != expectedMoney {
		t.Errorf("Expected money %d, got %d", expectedMoney, player.Money)
	}
	if player.HP != maxHP {
		t.Errorf("Expected HP %.2f, got %.2f", maxHP, player.HP)
	}
}

func TestHospitalSystem_ProcessHealing_InsufficientMoney(t *testing.T) {
	// Setup: Player with zero HP but insufficient money
	player := testHospitalPlayer(100, 100)
	player.Money = 5 // Need 20, only have 5
	player.HP = 0.0

	hospital := entities.NewHospital(80, 80)
	system := NewHospitalSystem(hospital)

	inputState := input.InputState{Sell: true}

	// Execute
	system.ProcessHealing(player, inputState)

	// Verify: No transaction (money and HP unchanged)
	if player.Money != 5 {
		t.Errorf("Expected money 5, got %d", player.Money)
	}
	if player.HP != 0.0 {
		t.Errorf("Expected HP 0.0, got %.2f", player.HP)
	}
}

func TestHospitalSystem_ProcessHealing_NoInput(t *testing.T) {
	// Setup: Player in range but no Sell input
	player := testHospitalPlayer(100, 100)
	player.Money = 100
	player.HP = 0.0

	hospital := entities.NewHospital(80, 80)
	system := NewHospitalSystem(hospital)

	inputState := input.InputState{Sell: false}

	// Execute
	system.ProcessHealing(player, inputState)

	// Verify: No transaction
	if player.Money != 100 {
		t.Errorf("Expected money 100, got %d", player.Money)
	}
	if player.HP != 0.0 {
		t.Errorf("Expected HP 0.0, got %.2f", player.HP)
	}
}

func TestHospitalSystem_ProcessHealing_OutOfRange(t *testing.T) {
	// Setup: Player far from hospital
	player := testHospitalPlayer(500, 500)
	player.Money = 100
	player.HP = 0.0

	hospital := entities.NewHospital(80, 80)
	system := NewHospitalSystem(hospital)

	inputState := input.InputState{Sell: true}

	// Execute
	system.ProcessHealing(player, inputState)

	// Verify: No transaction
	if player.Money != 100 {
		t.Errorf("Expected money 100, got %d", player.Money)
	}
	if player.HP != 0.0 {
		t.Errorf("Expected HP 0.0, got %.2f", player.HP)
	}
}

func TestHospitalSystem_ProcessHealing_SmallFractionalHP(t *testing.T) {
	// Setup: Player at 9.9 HP (need 0.1 HP = ceil(0.2) = $1)
	player := testHospitalPlayer(100, 100)
	player.Money = 100
	player.HP = 9.9

	hospital := entities.NewHospital(80, 80)
	system := NewHospitalSystem(hospital)

	inputState := input.InputState{Sell: true}

	maxHP := player.Hull.MaxHP()

	// Execute
	system.ProcessHealing(player, inputState)

	// Verify: 1 money deducted (ceil(0.1 * 2) = ceil(0.2) = 1), HP full
	expectedMoney := 100 - 1
	if player.Money != expectedMoney {
		t.Errorf("Expected money %d, got %d", expectedMoney, player.Money)
	}
	if player.HP != maxHP {
		t.Errorf("Expected HP %.2f, got %.2f", maxHP, player.HP)
	}
}

// Test helpers

func testHospitalPlayer(x, y float32) *entities.Player {
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
