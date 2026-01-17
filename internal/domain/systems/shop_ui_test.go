package systems_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/systems"
)

func createTestShopUISystem() (*systems.ShopUISystem, *entities.Player) {
	// Create unified shop at position where player (at 0,0) will be in range
	shop := entities.NewUpgradeShop(0, 0)
	system := systems.NewShopUISystem(shop)
	player := entities.NewPlayer(0, 0)

	return system, player
}

func TestShopUISystem_OpenShop_WhenInRange(t *testing.T) {
	system, player := createTestShopUISystem()

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

func TestShopUISystem_OpenShop_WhenOutOfRange(t *testing.T) {
	system, player := createTestShopUISystem()
	// Move player far away from shop
	player.AABB.X = 5000

	inputState := input.InputState{Sell: true}
	system.ProcessShopInteraction(player, inputState)

	if system.GetUIState().Open {
		t.Error("Expected shop to remain closed when out of range")
	}
}

func TestShopUISystem_CloseShop(t *testing.T) {
	system, player := createTestShopUISystem()

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

func TestShopUISystem_TabNavigation(t *testing.T) {
	system, player := createTestShopUISystem()

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

func TestShopUISystem_GridNavigation(t *testing.T) {
	system, player := createTestShopUISystem()

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

func TestShopUISystem_PurchaseEngine_Success(t *testing.T) {
	system, player := createTestShopUISystem()
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

func TestShopUISystem_PurchaseEngine_InsufficientFunds(t *testing.T) {
	system, player := createTestShopUISystem()
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

func TestShopUISystem_PurchaseOwnedUpgrade_Fails(t *testing.T) {
	system, player := createTestShopUISystem()
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

func TestShopUISystem_PurchaseSkipTier_Success(t *testing.T) {
	system, player := createTestShopUISystem()
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

func TestShopUISystem_ProgressivePurchases(t *testing.T) {
	system, player := createTestShopUISystem()
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

func TestShopUISystem_PurchaseHull_Success(t *testing.T) {
	system, player := createTestShopUISystem()
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
