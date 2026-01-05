package systems_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/systems"
)

func createTestUpgradeSystem() (*systems.UpgradeSystem, *entities.Player) {
	// Create shops at positions where player (at 0,0) will be in range
	engineShop := entities.NewUpgradeShop(0, 0, entities.UpgradeTypeEngine)
	hullShop := entities.NewUpgradeShop(400, 0, entities.UpgradeTypeHull)
	fuelTankShop := entities.NewUpgradeShop(800, 0, entities.UpgradeTypeFuelTank)

	system := systems.NewUpgradeSystem(engineShop, hullShop, fuelTankShop)
	player := entities.NewPlayer(0, 0)

	return system, player
}

func TestUpgradeSystem_BuyEngineMk1_Success(t *testing.T) {
	system, player := createTestUpgradeSystem()
	player.Money = 200 // More than enough for Engine Mk1 ($100)

	inputState := input.InputState{Sell: true}
	system.ProcessUpgrade(player, inputState)

	if player.Upgrades.Engine != entities.UpgradeLevelMk1 {
		t.Errorf("Expected engine level Mk1 (1), got %d", player.Upgrades.Engine)
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

	if player.Upgrades.Engine != entities.UpgradeLevelBase {
		t.Errorf("Expected engine level to remain Base (0), got %d", player.Upgrades.Engine)
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

	if player.Upgrades.Engine != entities.UpgradeLevelBase {
		t.Errorf("Expected engine level to remain Base (0), got %d", player.Upgrades.Engine)
	}
}

func TestUpgradeSystem_BuyEngine_OutOfRange(t *testing.T) {
	system, player := createTestUpgradeSystem()
	player.Money = 200
	// Move player far away from all shops
	player.AABB.X = 5000

	inputState := input.InputState{Sell: true}
	system.ProcessUpgrade(player, inputState)

	if player.Upgrades.Engine != entities.UpgradeLevelBase {
		t.Errorf("Expected engine level to remain Base (0), got %d", player.Upgrades.Engine)
	}
}

func TestUpgradeSystem_BuyEngine_MaxLevel(t *testing.T) {
	system, player := createTestUpgradeSystem()
	player.Money = 100000 // Plenty of money
	player.Upgrades.Engine = entities.UpgradeLevelMk5 // Already at max

	initialMoney := player.Money
	inputState := input.InputState{Sell: true}
	system.ProcessUpgrade(player, inputState)

	if player.Upgrades.Engine != entities.UpgradeLevelMk5 {
		t.Errorf("Expected engine level to remain Mk5 (5), got %d", player.Upgrades.Engine)
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

	if player.Upgrades.Hull != entities.UpgradeLevelMk1 {
		t.Errorf("Expected hull level Mk1 (1), got %d", player.Upgrades.Hull)
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

	if player.Upgrades.FuelTank != entities.UpgradeLevelMk1 {
		t.Errorf("Expected fuel tank level Mk1 (1), got %d", player.Upgrades.FuelTank)
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
	if player.Upgrades.Engine != entities.UpgradeLevelMk1 {
		t.Errorf("Expected engine Mk1 after first purchase, got %d", player.Upgrades.Engine)
	}

	// Buy Mk2
	system.ProcessUpgrade(player, inputState)
	if player.Upgrades.Engine != entities.UpgradeLevelMk2 {
		t.Errorf("Expected engine Mk2 after second purchase, got %d", player.Upgrades.Engine)
	}

	// Buy Mk3
	system.ProcessUpgrade(player, inputState)
	if player.Upgrades.Engine != entities.UpgradeLevelMk3 {
		t.Errorf("Expected engine Mk3 after third purchase, got %d", player.Upgrades.Engine)
	}
}
