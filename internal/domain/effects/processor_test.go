package effects

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
)

func TestProcessor_Apply_EmptyEffects(t *testing.T) {
	processor := NewProcessor()
	player := testPlayer()
	player.Money = 100

	processor.Apply(player, []Effect{})

	if player.Money != 100 {
		t.Errorf("expected money to remain 100, got %d", player.Money)
	}
}

func TestProcessor_Apply_NilEffects(t *testing.T) {
	processor := NewProcessor()
	player := testPlayer()
	player.Money = 100

	processor.Apply(player, nil)

	if player.Money != 100 {
		t.Errorf("expected money to remain 100, got %d", player.Money)
	}
}

func TestProcessor_Apply_SingleEffect(t *testing.T) {
	processor := NewProcessor()
	player := testPlayer()
	player.Money = 100

	effects := []Effect{
		AddMoney{Amount: 50},
	}

	processor.Apply(player, effects)

	if player.Money != 150 {
		t.Errorf("expected money to be 150, got %d", player.Money)
	}
}

func TestProcessor_Apply_MultipleEffects(t *testing.T) {
	processor := NewProcessor()
	player := testPlayer()
	player.Money = 100
	player.Fuel = 50.0
	player.HP = 75.0

	effects := []Effect{
		TakeMoney{Amount: 25},
		SetFuel{Amount: 100.0},
		SetHP{Amount: 100.0},
	}

	processor.Apply(player, effects)

	if player.Money != 75 {
		t.Errorf("expected money to be 75, got %d", player.Money)
	}
	if player.Fuel != 100.0 {
		t.Errorf("expected fuel to be 100.0, got %.2f", player.Fuel)
	}
	if player.HP != 100.0 {
		t.Errorf("expected HP to be 100.0, got %.2f", player.HP)
	}
}

func TestProcessor_Apply_EffectsAppliedInOrder(t *testing.T) {
	processor := NewProcessor()
	player := testPlayer()
	player.Money = 100

	// Add 50, then take 30: should result in 120
	effects := []Effect{
		AddMoney{Amount: 50},
		TakeMoney{Amount: 30},
	}

	processor.Apply(player, effects)

	if player.Money != 120 {
		t.Errorf("expected money to be 120 (100 + 50 - 30), got %d", player.Money)
	}
}

func TestProcessor_Apply_PurchaseScenario(t *testing.T) {
	processor := NewProcessor()
	player := testPlayer()
	player.Money = 500
	player.ItemInventory = [5]int{0, 0, 0, 0, 0}

	// Simulate buying an item: take money, add item
	effects := []Effect{
		TakeMoney{Amount: 100},
		AddItem{ItemType: entities.ItemTeleport},
	}

	processor.Apply(player, effects)

	if player.Money != 400 {
		t.Errorf("expected money to be 400, got %d", player.Money)
	}
	if player.ItemInventory[entities.ItemTeleport] != 1 {
		t.Errorf("expected teleport count to be 1, got %d", player.ItemInventory[entities.ItemTeleport])
	}
}

func TestProcessor_Apply_MarketSellScenario(t *testing.T) {
	processor := NewProcessor()
	player := testPlayer()
	player.Money = 100
	player.OreInventory = map[string]int{
		"copper": 5,
		"gold":   3,
	}

	// Simulate selling at market: add money, clear inventory
	effects := []Effect{
		AddMoney{Amount: 250},
		ClearOreInventory{},
	}

	processor.Apply(player, effects)

	if player.Money != 350 {
		t.Errorf("expected money to be 350, got %d", player.Money)
	}
	if len(player.OreInventory) != 0 {
		t.Errorf("expected ore inventory to be empty, got %d items", len(player.OreInventory))
	}
}

func TestProcessor_Apply_HospitalHealScenario(t *testing.T) {
	processor := NewProcessor()
	player := testPlayer()
	player.Money = 200
	player.HP = 50.0

	// Simulate healing at hospital: take money, set HP
	effects := []Effect{
		TakeMoney{Amount: 100},
		SetHP{Amount: 100.0},
	}

	processor.Apply(player, effects)

	if player.Money != 100 {
		t.Errorf("expected money to be 100, got %d", player.Money)
	}
	if player.HP != 100.0 {
		t.Errorf("expected HP to be 100.0, got %.2f", player.HP)
	}
}

func TestProcessor_Apply_FuelStationRefuelScenario(t *testing.T) {
	processor := NewProcessor()
	player := testPlayer()
	player.Money = 150
	player.Fuel = 25.0

	// Simulate refueling: take money, set fuel
	effects := []Effect{
		TakeMoney{Amount: 75},
		SetFuel{Amount: 100.0},
	}

	processor.Apply(player, effects)

	if player.Money != 75 {
		t.Errorf("expected money to be 75, got %d", player.Money)
	}
	if player.Fuel != 100.0 {
		t.Errorf("expected fuel to be 100.0, got %.2f", player.Fuel)
	}
}
