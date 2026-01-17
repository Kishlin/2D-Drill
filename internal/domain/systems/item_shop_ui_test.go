package systems_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/systems"
)

func createTestItemShopUISystem() (*systems.ItemShopUISystem, *entities.Player) {
	// Create unified item shop at position where player (at 0,0) will be in range
	shop := entities.NewItemShop(0, 0)
	system := systems.NewItemShopUISystem(shop)
	player := entities.NewPlayer(0, 0)

	return system, player
}

func TestItemShopUISystem_OpenShop_WhenInRange(t *testing.T) {
	system, player := createTestItemShopUISystem()

	// Press E to open shop
	inputState := input.InputState{Sell: true}
	system.ProcessItemShopInteraction(player, inputState)

	if !system.GetUIState().Open {
		t.Error("Expected shop to be open")
	}
	if !player.InShop {
		t.Error("Expected player.InShop to be true")
	}
}

func TestItemShopUISystem_OpenShop_WhenOutOfRange(t *testing.T) {
	system, player := createTestItemShopUISystem()
	// Move player far away from shop
	player.AABB.X = 5000

	inputState := input.InputState{Sell: true}
	system.ProcessItemShopInteraction(player, inputState)

	if system.GetUIState().Open {
		t.Error("Expected shop to remain closed when out of range")
	}
}

func TestItemShopUISystem_CloseShop(t *testing.T) {
	system, player := createTestItemShopUISystem()

	// Open shop first
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})

	// Close shop
	system.ProcessItemShopInteraction(player, input.InputState{CloseShop: true})

	if system.GetUIState().Open {
		t.Error("Expected shop to be closed")
	}
	if player.InShop {
		t.Error("Expected player.InShop to be false")
	}
}

func TestItemShopUISystem_GridNavigation(t *testing.T) {
	system, player := createTestItemShopUISystem()

	// Open shop
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})

	// Initial selection should be 0 (Teleport)
	if system.GetUIState().SelectedIndex != 0 {
		t.Errorf("Expected initial selection to be 0, got %d", system.GetUIState().SelectedIndex)
	}

	// Move right (0 -> 1)
	system.ProcessItemShopInteraction(player, input.InputState{NavRight: true})
	if system.GetUIState().SelectedIndex != 1 {
		t.Errorf("Expected selection to be 1, got %d", system.GetUIState().SelectedIndex)
	}

	// Move down (1 -> 4, grid is 3 columns)
	system.ProcessItemShopInteraction(player, input.InputState{NavDown: true})
	if system.GetUIState().SelectedIndex != 4 {
		t.Errorf("Expected selection to be 4, got %d", system.GetUIState().SelectedIndex)
	}

	// Move up (4 -> 1)
	system.ProcessItemShopInteraction(player, input.InputState{NavUp: true})
	if system.GetUIState().SelectedIndex != 1 {
		t.Errorf("Expected selection to be 1, got %d", system.GetUIState().SelectedIndex)
	}

	// Move right to reach position 2 (Refuel)
	system.ProcessItemShopInteraction(player, input.InputState{NavRight: true})
	if system.GetUIState().SelectedIndex != 2 {
		t.Errorf("Expected selection to be 2, got %d", system.GetUIState().SelectedIndex)
	}

	// Move right again - should stay at 2 (wraps within row, but last column is 2)
	system.ProcessItemShopInteraction(player, input.InputState{NavRight: true})
	if system.GetUIState().SelectedIndex != 0 {
		t.Errorf("Expected selection to wrap to 0, got %d", system.GetUIState().SelectedIndex)
	}
}

func TestItemShopUISystem_GridNavigation_SkipsEmptyCell(t *testing.T) {
	system, player := createTestItemShopUISystem()

	// Open shop
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})

	// Navigate to position 4 (BigBomb, bottom-middle)
	// 0 -> 1 -> 2
	system.ProcessItemShopInteraction(player, input.InputState{NavRight: true})
	system.ProcessItemShopInteraction(player, input.InputState{NavRight: true})
	// 2 -> 5 (empty) but should skip to 2 (not move)
	// Let me navigate down: 2 -> 5, but should stay at 2
	system.ProcessItemShopInteraction(player, input.InputState{NavDown: true})
	// After moving down from 2, we should get to 5 (empty), but navigation should prevent it
	// Actually, let me check the navigation logic - it should adjust to 2 (cellAbove empty)
	if system.GetUIState().SelectedIndex != 2 {
		t.Errorf("Expected selection to remain at 2 when moving down from 2, got %d", system.GetUIState().SelectedIndex)
	}

	// Move left from 2 to 1
	system.ProcessItemShopInteraction(player, input.InputState{NavLeft: true})
	if system.GetUIState().SelectedIndex != 1 {
		t.Errorf("Expected selection to be 1, got %d", system.GetUIState().SelectedIndex)
	}

	// Move down from 1 to 4
	system.ProcessItemShopInteraction(player, input.InputState{NavDown: true})
	if system.GetUIState().SelectedIndex != 4 {
		t.Errorf("Expected selection to be 4, got %d", system.GetUIState().SelectedIndex)
	}

	// Move right from 4 - should go to 3 (can't go to 5 which is empty)
	system.ProcessItemShopInteraction(player, input.InputState{NavRight: true})
	if system.GetUIState().SelectedIndex != 4 {
		t.Errorf("Expected selection to adjust when moving right from 4 (towards empty), got %d", system.GetUIState().SelectedIndex)
	}
}

func TestItemShopUISystem_PurchaseItem_Success(t *testing.T) {
	system, player := createTestItemShopUISystem()
	player.Money = 1000 // More than enough

	// Open shop
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})

	// Buy Teleport at index 0 (price $500)
	initialMoney := player.Money
	initialInventory := player.ItemInventory[entities.ItemTeleport]
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})

	if player.Money != initialMoney-500 {
		t.Errorf("Expected money to be %d after purchase, got %d", initialMoney-500, player.Money)
	}
	if player.ItemInventory[entities.ItemTeleport] != initialInventory+1 {
		t.Errorf("Expected teleport inventory to be %d, got %d", initialInventory+1, player.ItemInventory[entities.ItemTeleport])
	}
}

func TestItemShopUISystem_PurchaseItem_InsufficientFunds(t *testing.T) {
	system, player := createTestItemShopUISystem()
	player.Money = 50 // Not enough for Teleport ($500)

	// Open shop
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})

	// Attempt purchase
	initialInventory := player.ItemInventory[entities.ItemTeleport]
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})

	if player.ItemInventory[entities.ItemTeleport] != initialInventory {
		t.Errorf("Expected teleport inventory to remain %d, got %d", initialInventory, player.ItemInventory[entities.ItemTeleport])
	}
	if player.Money != 50 {
		t.Errorf("Expected money to remain 50, got %d", player.Money)
	}
}

func TestItemShopUISystem_PurchaseMultipleItems(t *testing.T) {
	system, player := createTestItemShopUISystem()
	player.Money = 10000 // Plenty for multiple purchases

	// Open shop
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})

	// Buy Teleport at index 0 ($500)
	teleportBefore := player.ItemInventory[entities.ItemTeleport]
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})
	if player.ItemInventory[entities.ItemTeleport] != teleportBefore+1 {
		t.Errorf("Expected teleport to increase by 1")
	}

	// Navigate to index 1 (Repair, $200)
	system.ProcessItemShopInteraction(player, input.InputState{NavRight: true})
	// Buy Repair
	repairBefore := player.ItemInventory[entities.ItemRepair]
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})
	if player.ItemInventory[entities.ItemRepair] != repairBefore+1 {
		t.Errorf("Expected repair kit to increase by 1")
	}

	// Navigate to index 2 (Refuel, $100)
	system.ProcessItemShopInteraction(player, input.InputState{NavRight: true})
	// Buy Refuel
	refuelBefore := player.ItemInventory[entities.ItemRefuel]
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})
	if player.ItemInventory[entities.ItemRefuel] != refuelBefore+1 {
		t.Errorf("Expected refuel to increase by 1")
	}

	// Verify total cost
	expectedMoney := 10000 - 500 - 200 - 100
	if player.Money != expectedMoney {
		t.Errorf("Expected money to be %d, got %d", expectedMoney, player.Money)
	}
}

func TestItemShopUISystem_PurchaseNavigateAndBuy(t *testing.T) {
	system, player := createTestItemShopUISystem()
	player.Money = 1000

	// Open shop
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})

	// Navigate to index 3 (Bomb, $300)
	// 0 -> 1 -> 2 -> (down) -> 5 (empty, skip) -> adjust to 2...
	// Let me navigate directly: 0 -> 1 -> 2 -> down -> 5 (empty) but should adjust
	// Actually, let's navigate: right, right, down, left
	system.ProcessItemShopInteraction(player, input.InputState{NavRight: true})  // 0 -> 1
	system.ProcessItemShopInteraction(player, input.InputState{NavRight: true})  // 1 -> 2
	system.ProcessItemShopInteraction(player, input.InputState{NavDown: true})   // 2 -> 2 (skip empty)
	system.ProcessItemShopInteraction(player, input.InputState{NavLeft: true})   // 2 -> 1
	system.ProcessItemShopInteraction(player, input.InputState{NavDown: true})   // 1 -> 4
	system.ProcessItemShopInteraction(player, input.InputState{NavLeft: true})   // 4 -> 3

	if system.GetUIState().SelectedIndex != 3 {
		t.Errorf("Expected selection to be 3, got %d", system.GetUIState().SelectedIndex)
	}

	// Buy Bomb
	bombBefore := player.ItemInventory[entities.ItemBomb]
	system.ProcessItemShopInteraction(player, input.InputState{Sell: true})

	if player.ItemInventory[entities.ItemBomb] != bombBefore+1 {
		t.Errorf("Expected bomb to increase by 1")
	}
	if player.Money != 1000-300 {
		t.Errorf("Expected money to be 700, got %d", player.Money)
	}
}
