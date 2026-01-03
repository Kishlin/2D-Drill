package systems

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

func TestHospitalSystem_ProcessHealing_FullHP(t *testing.T) {
	// Setup: Player with full HP
	player := entities.NewPlayer(100, 100)
	player.Money = 100
	player.HP = entities.MaxHP

	hospital := entities.NewHospital(80, 80)
	system := NewHospitalSystem(hospital)

	inputState := input.InputState{Sell: true}

	initialMoney := player.Money

	// Execute
	system.ProcessHealing(player, inputState)

	// Verify: No money deducted, HP stays full
	if player.Money != initialMoney {
		t.Errorf("Expected money %d, got %d", initialMoney, player.Money)
	}
	if player.HP != entities.MaxHP {
		t.Errorf("Expected HP %.2f, got %.2f", entities.MaxHP, player.HP)
	}
}

func TestHospitalSystem_ProcessHealing_ZeroHP(t *testing.T) {
	// Setup: Player with zero HP
	player := entities.NewPlayer(100, 100)
	player.Money = 100
	player.HP = 0.0

	hospital := entities.NewHospital(80, 80)
	system := NewHospitalSystem(hospital)

	inputState := input.InputState{Sell: true}

	// Execute
	system.ProcessHealing(player, inputState)

	// Verify: 20 money deducted (10 HP * $2), HP restored to max
	expectedMoney := 100 - 20
	if player.Money != expectedMoney {
		t.Errorf("Expected money %d, got %d", expectedMoney, player.Money)
	}
	if player.HP != entities.MaxHP {
		t.Errorf("Expected HP %.2f, got %.2f", entities.MaxHP, player.HP)
	}
}

func TestHospitalSystem_ProcessHealing_PartialHPRoundedUp(t *testing.T) {
	// Setup: Player at 7.2 HP (need 2.8 HP = ceil(5.6) = $6)
	player := entities.NewPlayer(100, 100)
	player.Money = 100
	player.HP = 7.2

	hospital := entities.NewHospital(80, 80)
	system := NewHospitalSystem(hospital)

	inputState := input.InputState{Sell: true}

	// Execute
	system.ProcessHealing(player, inputState)

	// Verify: 6 money deducted (ceil(2.8 * 2) = ceil(5.6) = 6), HP full
	expectedMoney := 100 - 6
	if player.Money != expectedMoney {
		t.Errorf("Expected money %d, got %d", expectedMoney, player.Money)
	}
	if player.HP != entities.MaxHP {
		t.Errorf("Expected HP %.2f, got %.2f", entities.MaxHP, player.HP)
	}
}

func TestHospitalSystem_ProcessHealing_InsufficientMoney(t *testing.T) {
	// Setup: Player with zero HP but insufficient money
	player := entities.NewPlayer(100, 100)
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
	player := entities.NewPlayer(100, 100)
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
	player := entities.NewPlayer(500, 500)
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
	player := entities.NewPlayer(100, 100)
	player.Money = 100
	player.HP = 9.9

	hospital := entities.NewHospital(80, 80)
	system := NewHospitalSystem(hospital)

	inputState := input.InputState{Sell: true}

	// Execute
	system.ProcessHealing(player, inputState)

	// Verify: 1 money deducted (ceil(0.1 * 2) = ceil(0.2) = 1), HP full
	expectedMoney := 100 - 1
	if player.Money != expectedMoney {
		t.Errorf("Expected money %d, got %d", expectedMoney, player.Money)
	}
	if player.HP != entities.MaxHP {
		t.Errorf("Expected HP %.2f, got %.2f", entities.MaxHP, player.HP)
	}
}
