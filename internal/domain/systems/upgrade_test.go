package systems_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/systems"
)

func createTestUpgradeSystem() (*systems.UpgradeSystem, *entities.Player) {
	// Create shops at positions where player (at 0,0) will be in range
	engineShop := entities.NewEngineUpgradeShop(0, 0)
	hullShop := entities.NewHullUpgradeShop(400, 0)
	fuelTankShop := entities.NewFuelTankUpgradeShop(800, 0)

	system := systems.NewUpgradeSystem(engineShop, hullShop, fuelTankShop)
	player := entities.NewPlayer(0, 0)

	return system, player
}

func TestUpgradeSystem_BuyEngineMk1_Success(t *testing.T) {
	system, player := createTestUpgradeSystem()
	player.Money = 200 // More than enough for Engine Mk1 ($100)

	inputState := input.InputState{Sell: true}
	system.ProcessUpgrade(player, inputState)

	if player.Engine.Tier() != 1 {
		t.Errorf("Expected engine tier 1, got %d", player.Engine.Tier())
	}
	if player.Money != 100 {
		t.Errorf("Expected money to be 100 after purchase, got %d", player.Money)
	}
}

func TestUpgradeSystem_BuyEngine_InsufficientFunds(t *testing.T) {
	system, player := createTestUpgradeSystem()
	player.Money = 50 // Not enough for Engine Mk1 ($100)

	inputState := input.InputState{Sell: true}
	system.ProcessUpgrade(player, inputState)

	if player.Engine.Tier() != 0 {
		t.Errorf("Expected engine tier to remain 0, got %d", player.Engine.Tier())
	}
	if player.Money != 50 {
		t.Errorf("Expected money to remain 50, got %d", player.Money)
	}
}

func TestUpgradeSystem_BuyEngine_NoInput(t *testing.T) {
	system, player := createTestUpgradeSystem()
	player.Money = 200

	inputState := input.InputState{Sell: false}
	system.ProcessUpgrade(player, inputState)

	if player.Engine.Tier() != 0 {
		t.Errorf("Expected engine tier to remain 0, got %d", player.Engine.Tier())
	}
}

func TestUpgradeSystem_BuyEngine_OutOfRange(t *testing.T) {
	system, player := createTestUpgradeSystem()
	player.Money = 200
	// Move player far away from all shops
	player.AABB.X = 5000

	inputState := input.InputState{Sell: true}
	system.ProcessUpgrade(player, inputState)

	if player.Engine.Tier() != 0 {
		t.Errorf("Expected engine tier to remain 0, got %d", player.Engine.Tier())
	}
}

func TestUpgradeSystem_BuyEngine_MaxLevel(t *testing.T) {
	system, player := createTestUpgradeSystem()
	player.Money = 100000 // Plenty of money

	inputState := input.InputState{Sell: true}

	// Buy all engine upgrades (Mk1 through Mk5)
	for i := 0; i < 5; i++ {
		system.ProcessUpgrade(player, inputState)
	}

	if player.Engine.Tier() != 5 {
		t.Errorf("Expected engine tier 5, got %d", player.Engine.Tier())
	}

	initialMoney := player.Money
	// Try to buy again at max level
	system.ProcessUpgrade(player, inputState)

	if player.Engine.Tier() != 5 {
		t.Errorf("Expected engine tier to remain 5, got %d", player.Engine.Tier())
	}
	if player.Money != initialMoney {
		t.Errorf("Expected money to remain unchanged at max level")
	}
}

func TestUpgradeSystem_BuyHull_Success(t *testing.T) {
	system, player := createTestUpgradeSystem()
	player.Money = 200
	// Move player to hull shop position
	player.AABB.X = 400

	inputState := input.InputState{Sell: true}
	system.ProcessUpgrade(player, inputState)

	if player.Hull.Tier() != 1 {
		t.Errorf("Expected hull tier 1, got %d", player.Hull.Tier())
	}
	if player.Money != 50 { // Hull Mk1 costs $150
		t.Errorf("Expected money to be 50 after purchase, got %d", player.Money)
	}
}

func TestUpgradeSystem_BuyFuelTank_Success(t *testing.T) {
	system, player := createTestUpgradeSystem()
	player.Money = 200
	// Move player to fuel tank shop position
	player.AABB.X = 800

	inputState := input.InputState{Sell: true}
	system.ProcessUpgrade(player, inputState)

	if player.FuelTank.Tier() != 1 {
		t.Errorf("Expected fuel tank tier 1, got %d", player.FuelTank.Tier())
	}
	if player.Money != 100 { // Tank Mk1 costs $100
		t.Errorf("Expected money to be 100 after purchase, got %d", player.Money)
	}
}

func TestUpgradeSystem_ProgressiveUpgrades(t *testing.T) {
	system, player := createTestUpgradeSystem()
	player.Money = 10000 // Plenty for multiple upgrades

	inputState := input.InputState{Sell: true}

	// Buy Mk1
	system.ProcessUpgrade(player, inputState)
	if player.Engine.Tier() != 1 {
		t.Errorf("Expected engine tier 1 after first purchase, got %d", player.Engine.Tier())
	}

	// Buy Mk2
	system.ProcessUpgrade(player, inputState)
	if player.Engine.Tier() != 2 {
		t.Errorf("Expected engine tier 2 after second purchase, got %d", player.Engine.Tier())
	}

	// Buy Mk3
	system.ProcessUpgrade(player, inputState)
	if player.Engine.Tier() != 3 {
		t.Errorf("Expected engine tier 3 after third purchase, got %d", player.Engine.Tier())
	}
}
