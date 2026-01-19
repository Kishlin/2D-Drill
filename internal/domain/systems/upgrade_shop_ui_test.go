package systems_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/systems"
)

func TestUpgradeShopUISystem_OpenShop_WhenInRange(t *testing.T) {
	system, player := createTestUpgradeShopUISystem()

	// Press E to open shop
	inputState := input.InputState{Sell: true}
	system.ProcessShopInteraction(player, inputState)

	if !system.GetUIState().Open {
		t.Error("Expected shop to be open")
	}
	if !player.InShop {
		t.Error("Expected player.InShop to be true")
	}
}

func TestUpgradeShopUISystem_OpenShop_WhenOutOfRange(t *testing.T) {
	system, player := createTestUpgradeShopUISystem()
	// Move player far away from shop
	player.AABB.X = 5000

	inputState := input.InputState{Sell: true}
	system.ProcessShopInteraction(player, inputState)

	if system.GetUIState().Open {
		t.Error("Expected shop to remain closed when out of range")
	}
}

func TestUpgradeShopUISystem_CloseShop(t *testing.T) {
	system, player := createTestUpgradeShopUISystem()

	// Open shop first
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	// Close shop
	system.ProcessShopInteraction(player, input.InputState{CloseShop: true})

	if system.GetUIState().Open {
		t.Error("Expected shop to be closed")
	}
	if player.InShop {
		t.Error("Expected player.InShop to be false")
	}
}

func TestUpgradeShopUISystem_TabNavigation(t *testing.T) {
	system, player := createTestUpgradeShopUISystem()

	// Open shop
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	// Initial tab should be Engine (0)
	if system.GetUIState().ActiveTab != entities.UpgradeEngine {
		t.Errorf("Expected initial tab to be Engine, got %v", system.GetUIState().ActiveTab)
	}

	// Press X to go to next tab (Hull)
	system.ProcessShopInteraction(player, input.InputState{NextTab: true})
	if system.GetUIState().ActiveTab != entities.UpgradeHull {
		t.Errorf("Expected tab to be Hull, got %v", system.GetUIState().ActiveTab)
	}

	// Press Z to go back to Engine
	system.ProcessShopInteraction(player, input.InputState{PrevTab: true})
	if system.GetUIState().ActiveTab != entities.UpgradeEngine {
		t.Errorf("Expected tab to be Engine, got %v", system.GetUIState().ActiveTab)
	}

	// Press Z again to wrap to Drill (last tab)
	system.ProcessShopInteraction(player, input.InputState{PrevTab: true})
	if system.GetUIState().ActiveTab != entities.UpgradeDrill {
		t.Errorf("Expected tab to wrap to Drill, got %v", system.GetUIState().ActiveTab)
	}
}

func TestUpgradeShopUISystem_GridNavigation(t *testing.T) {
	system, player := createTestUpgradeShopUISystem()

	// Open shop
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	// Initial selection should be tier 0 (Base)
	if system.GetUIState().SelectedTier != 0 {
		t.Errorf("Expected initial tier to be 0, got %d", system.GetUIState().SelectedTier)
	}

	// Move right (0 -> 1)
	system.ProcessShopInteraction(player, input.InputState{NavRight: true})
	if system.GetUIState().SelectedTier != 1 {
		t.Errorf("Expected tier to be 1, got %d", system.GetUIState().SelectedTier)
	}

	// Move down (1 -> 4, grid is 3 columns)
	system.ProcessShopInteraction(player, input.InputState{NavDown: true})
	if system.GetUIState().SelectedTier != 4 {
		t.Errorf("Expected tier to be 4, got %d", system.GetUIState().SelectedTier)
	}

	// Move up (4 -> 1)
	system.ProcessShopInteraction(player, input.InputState{NavUp: true})
	if system.GetUIState().SelectedTier != 1 {
		t.Errorf("Expected tier to be 1, got %d", system.GetUIState().SelectedTier)
	}
}

func TestUpgradeShopUISystem_PurchaseEngine_Success(t *testing.T) {
	system, player := createTestUpgradeShopUISystem()
	player.Money = 200 // More than enough for Engine Mk1 ($100)

	// Open shop
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	// Navigate to Mk1 (tier 1)
	system.ProcessShopInteraction(player, input.InputState{NavRight: true})

	// Purchase
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	if player.Engine.Tier() != 1 {
		t.Errorf("Expected engine tier 1, got %d", player.Engine.Tier())
	}
	if player.Money != 100 {
		t.Errorf("Expected money to be 100 after purchase, got %d", player.Money)
	}
}

func TestUpgradeShopUISystem_PurchaseEngine_InsufficientFunds(t *testing.T) {
	system, player := createTestUpgradeShopUISystem()
	player.Money = 50 // Not enough for Engine Mk1 ($100)

	// Open shop
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	// Navigate to Mk1
	system.ProcessShopInteraction(player, input.InputState{NavRight: true})

	// Attempt purchase
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	if player.Engine.Tier() != 0 {
		t.Errorf("Expected engine tier to remain 0, got %d", player.Engine.Tier())
	}
	if player.Money != 50 {
		t.Errorf("Expected money to remain 50, got %d", player.Money)
	}
}

func TestUpgradeShopUISystem_PurchaseOwnedUpgrade_Fails(t *testing.T) {
	system, player := createTestUpgradeShopUISystem()
	player.Money = 1000

	// Open shop
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	// Try to purchase Base (tier 0) which is already owned
	initialMoney := player.Money
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	if player.Money != initialMoney {
		t.Errorf("Expected money to remain unchanged when trying to buy owned upgrade")
	}
}

func TestUpgradeShopUISystem_PurchaseSkipTier_Success(t *testing.T) {
	system, player := createTestUpgradeShopUISystem()
	player.Money = 10000 // Plenty of money

	// Open shop
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	// Navigate to Mk2 (tier 2), skipping Mk1
	system.ProcessShopInteraction(player, input.InputState{NavRight: true}) // to tier 1
	system.ProcessShopInteraction(player, input.InputState{NavRight: true}) // to tier 2

	// Purchase Mk2 (should succeed even without Mk1 owned)
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	if player.Engine.Tier() != 2 {
		t.Errorf("Expected engine tier to be 2 (can skip tiers), got %d", player.Engine.Tier())
	}
}

func TestUpgradeShopUISystem_ProgressivePurchases(t *testing.T) {
	system, player := createTestUpgradeShopUISystem()
	player.Money = 10000 // Plenty for multiple upgrades

	// Open shop
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	// Grid layout:
	// [0] [1] [2]  <- Base, Mk1, Mk2
	// [3] [4] [5]  <- Mk3, Mk4, Mk5

	// Purchase Mk1 (navigate to tier 1, then buy)
	system.ProcessShopInteraction(player, input.InputState{NavRight: true}) // 0 -> 1
	system.ProcessShopInteraction(player, input.InputState{Sell: true})
	if player.Engine.Tier() != 1 {
		t.Errorf("Expected engine tier 1, got %d", player.Engine.Tier())
	}

	// Purchase Mk2 (navigate to tier 2, then buy)
	system.ProcessShopInteraction(player, input.InputState{NavRight: true}) // 1 -> 2
	system.ProcessShopInteraction(player, input.InputState{Sell: true})
	if player.Engine.Tier() != 2 {
		t.Errorf("Expected engine tier 2, got %d", player.Engine.Tier())
	}

	// Purchase Mk3 (navigate to tier 3, which is down-left from tier 2)
	system.ProcessShopInteraction(player, input.InputState{NavDown: true}) // 2 -> 5
	system.ProcessShopInteraction(player, input.InputState{NavLeft: true}) // 5 -> 4
	system.ProcessShopInteraction(player, input.InputState{NavLeft: true}) // 4 -> 3
	system.ProcessShopInteraction(player, input.InputState{Sell: true})
	if player.Engine.Tier() != 3 {
		t.Errorf("Expected engine tier 3, got %d", player.Engine.Tier())
	}
}

func TestUpgradeShopUISystem_PurchaseHull_Success(t *testing.T) {
	system, player := createTestUpgradeShopUISystem()
	player.Money = 200

	// Open shop
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	// Switch to Hull tab
	system.ProcessShopInteraction(player, input.InputState{NextTab: true})

	// Navigate to Mk1
	system.ProcessShopInteraction(player, input.InputState{NavRight: true})

	// Purchase
	system.ProcessShopInteraction(player, input.InputState{Sell: true})

	if player.Hull.Tier() != 1 {
		t.Errorf("Expected hull tier 1, got %d", player.Hull.Tier())
	}
	if player.Money != 50 { // Hull Mk1 costs $150
		t.Errorf("Expected money to be 50 after purchase, got %d", player.Money)
	}
}

// Test helpers

func testUpgradeConfig() config.UpgradeConfig {
	return config.UpgradeConfig{
		Engines: []config.UpgradeTier[config.EngineStats]{
			{Name: "Base", Price: 0, Stats: config.EngineStats{MaxSpeed: 450, Acceleration: 2500, FlyAcceleration: 2500, MaxUpwardSpeed: -600}},
			{Name: "Mk1", Price: 100, Stats: config.EngineStats{MaxSpeed: 500, Acceleration: 2750, FlyAcceleration: 2750, MaxUpwardSpeed: -650}},
			{Name: "Mk2", Price: 300, Stats: config.EngineStats{MaxSpeed: 550, Acceleration: 3000, FlyAcceleration: 3000, MaxUpwardSpeed: -700}},
			{Name: "Mk3", Price: 750, Stats: config.EngineStats{MaxSpeed: 600, Acceleration: 3250, FlyAcceleration: 3250, MaxUpwardSpeed: -750}},
			{Name: "Mk4", Price: 1500, Stats: config.EngineStats{MaxSpeed: 650, Acceleration: 3500, FlyAcceleration: 3500, MaxUpwardSpeed: -800}},
			{Name: "Mk5", Price: 5000, Stats: config.EngineStats{MaxSpeed: 700, Acceleration: 4000, FlyAcceleration: 4000, MaxUpwardSpeed: -900}},
		},
		Hulls: []config.UpgradeTier[config.HullStats]{
			{Name: "Base", Price: 0, Stats: config.HullStats{MaxHP: 10}},
			{Name: "Mk1", Price: 150, Stats: config.HullStats{MaxHP: 15}},
			{Name: "Mk2", Price: 400, Stats: config.HullStats{MaxHP: 20}},
			{Name: "Mk3", Price: 1000, Stats: config.HullStats{MaxHP: 30}},
			{Name: "Mk4", Price: 2500, Stats: config.HullStats{MaxHP: 45}},
			{Name: "Mk5", Price: 8000, Stats: config.HullStats{MaxHP: 70}},
		},
		FuelTanks: []config.UpgradeTier[config.FuelTankStats]{
			{Name: "Base", Price: 0, Stats: config.FuelTankStats{Capacity: 10}},
			{Name: "Mk1", Price: 100, Stats: config.FuelTankStats{Capacity: 15}},
			{Name: "Mk2", Price: 250, Stats: config.FuelTankStats{Capacity: 22}},
			{Name: "Mk3", Price: 600, Stats: config.FuelTankStats{Capacity: 32}},
			{Name: "Mk4", Price: 1500, Stats: config.FuelTankStats{Capacity: 45}},
			{Name: "Mk5", Price: 4000, Stats: config.FuelTankStats{Capacity: 65}},
		},
		CargoHolds: []config.UpgradeTier[config.CargoHoldStats]{
			{Name: "Base", Price: 0, Stats: config.CargoHoldStats{Capacity: 10}},
			{Name: "Mk1", Price: 75, Stats: config.CargoHoldStats{Capacity: 15}},
			{Name: "Mk2", Price: 200, Stats: config.CargoHoldStats{Capacity: 25}},
			{Name: "Mk3", Price: 500, Stats: config.CargoHoldStats{Capacity: 40}},
			{Name: "Mk4", Price: 1250, Stats: config.CargoHoldStats{Capacity: 60}},
			{Name: "Mk5", Price: 3500, Stats: config.CargoHoldStats{Capacity: 100}},
		},
		HeatShields: []config.UpgradeTier[config.HeatShieldStats]{
			{Name: "Base", Price: 0, Stats: config.HeatShieldStats{HeatResistance: 50}},
			{Name: "Mk1", Price: 200, Stats: config.HeatShieldStats{HeatResistance: 90}},
			{Name: "Mk2", Price: 500, Stats: config.HeatShieldStats{HeatResistance: 140}},
			{Name: "Mk3", Price: 1200, Stats: config.HeatShieldStats{HeatResistance: 200}},
			{Name: "Mk4", Price: 3000, Stats: config.HeatShieldStats{HeatResistance: 270}},
			{Name: "Mk5", Price: 10000, Stats: config.HeatShieldStats{HeatResistance: 350}},
		},
		Drills: []config.UpgradeTier[config.DrillStats]{
			{Name: "Base", Price: 0, Stats: config.DrillStats{DrillSpeed: 1.0}},
			{Name: "Mk1", Price: 125, Stats: config.DrillStats{DrillSpeed: 1.5}},
			{Name: "Mk2", Price: 350, Stats: config.DrillStats{DrillSpeed: 2.2}},
			{Name: "Mk3", Price: 875, Stats: config.DrillStats{DrillSpeed: 3.2}},
			{Name: "Mk4", Price: 2000, Stats: config.DrillStats{DrillSpeed: 4.5}},
			{Name: "Mk5", Price: 6500, Stats: config.DrillStats{DrillSpeed: 6.5}},
		},
	}
}

func testPlayerConfig() config.PlayerConfig {
	return config.PlayerConfig{
		StartingMoney:    0,
		StartingItems:    [5]int{0, 0, 0, 0, 0},
		StartingUpgrades: config.StartingUpgrades{},
	}
}

func testUpgradeShopPlayer() *entities.Player {
	return entities.NewPlayerFromConfig(0, 0, testPlayerConfig(), testUpgradeConfig())
}

func testUpgradeShop() *entities.UpgradeShop {
	return entities.NewUpgradeShopFromConfig(0, 0, testUpgradeConfig())
}

func createTestUpgradeShopUISystem() (*systems.UpgradeShopUISystem, *entities.Player) {
	shop := testUpgradeShop()
	system := systems.NewUpgradeShopUISystem(shop)
	player := testUpgradeShopPlayer()

	return system, player
}
