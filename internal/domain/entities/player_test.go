package entities

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
)

// Test helper - creates player with base stats for tests
func testPlayer() *Player {
	playerCfg := config.PlayerConfig{
		StartingMoney:    0,
		StartingItems:    [5]int{0, 0, 0, 0, 0},
		StartingUpgrades: config.StartingUpgrades{},
	}
	upgradeCfg := config.UpgradeConfig{
		Engines:     []config.UpgradeTier[config.EngineStats]{{Price: 0, Stats: config.EngineStats{MaxSpeed: 450, Acceleration: 2500, FlyAcceleration: 2500, MaxUpwardSpeed: -600}}},
		Hulls:       []config.UpgradeTier[config.HullStats]{{Price: 0, Stats: config.HullStats{MaxHP: 10}}},
		FuelTanks:   []config.UpgradeTier[config.FuelTankStats]{{Price: 0, Stats: config.FuelTankStats{Capacity: 10}}},
		CargoHolds:  []config.UpgradeTier[config.CargoHoldStats]{{Price: 0, Stats: config.CargoHoldStats{Capacity: 10}}},
		HeatShields: []config.UpgradeTier[config.HeatShieldStats]{{Price: 0, Stats: config.HeatShieldStats{HeatResistance: 50}}},
		Drills:      []config.UpgradeTier[config.DrillStats]{{Price: 0, Stats: config.DrillStats{SpeedAtSurface: 1.0, SpeedAtMaxDepth: 1.0}}},
	}
	return NewPlayerFromConfig(0, 0, playerCfg, upgradeCfg)
}

func TestPlayer_AddOreByID_SingleType(t *testing.T) {
	player := testPlayer()

	success := player.AddOreByID("copper")

	if success == false {
		t.Errorf("Expected AddOreByID to succeed")
	}
	if player.OreInventory["copper"] != 1 {
		t.Errorf("Expected 1 copper, got %d", player.OreInventory["copper"])
	}
}

func TestPlayer_AddOreByID_MultipleTypes(t *testing.T) {
	player := testPlayer()

	player.AddOreByID("copper")
	player.AddOreByID("copper")
	player.AddOreByID("copper")
	player.AddOreByID("gold")
	player.AddOreByID("gold")
	player.AddOreByID("gold")
	player.AddOreByID("gold")
	player.AddOreByID("diamond")

	if player.OreInventory["copper"] != 3 {
		t.Errorf("Expected 3 copper, got %d", player.OreInventory["copper"])
	}
	if player.OreInventory["gold"] != 4 {
		t.Errorf("Expected 4 gold, got %d", player.OreInventory["gold"])
	}
	if player.OreInventory["diamond"] != 1 {
		t.Errorf("Expected 1 diamond, got %d", player.OreInventory["diamond"])
	}
}

func TestPlayer_AddOreByID_Accumulates(t *testing.T) {
	player := testPlayer()

	for i := 0; i < 10; i++ {
		player.AddOreByID("iron")
	}

	if player.OreInventory["iron"] != 10 {
		t.Errorf("Expected 10 iron, got %d", player.OreInventory["iron"])
	}
}

func TestPlayer_NewPlayer_StartsWithZeroOres(t *testing.T) {
	player := testPlayer()

	// Check some ore types
	oreIDs := []string{"copper", "iron", "gold", "diamond"}
	for _, oreID := range oreIDs {
		if player.OreInventory[oreID] != 0 {
			t.Errorf("New player should have 0 of ore %s, got %d", oreID, player.OreInventory[oreID])
		}
	}
}

func TestPlayer_AddOreByID_BoundsCheck(t *testing.T) {
	player := testPlayer()

	// Should return false on empty ore ID
	if player.AddOreByID("") {
		t.Errorf("Should return false for empty ore ID")
	}

	// Empty string should not be in inventory
	if player.OreInventory[""] != 0 {
		t.Errorf("Empty ore ID should not affect inventory")
	}
}

func TestPlayer_AddOreByID_CargoCapacity(t *testing.T) {
	player := testPlayer()
	// Player starts with Base CargoHold, capacity 10

	// Fill cargo to capacity
	for i := 0; i < 10; i++ {
		if player.AddOreByID("copper") == false {
			t.Errorf("AddOreByID should succeed at position %d", i)
		}
	}

	if player.GetTotalOreCount() != 10 {
		t.Errorf("Expected 10 total ore, got %d", player.GetTotalOreCount())
	}

	// Next ore should fail
	if player.AddOreByID("copper") {
		t.Errorf("AddOreByID should fail when cargo is full")
	}

	if player.GetTotalOreCount() != 10 {
		t.Errorf("Failed AddOreByID should not change inventory")
	}
}

// DealDamage tests

func TestPlayer_DealDamage_ReducesHP(t *testing.T) {
	player := testPlayer()
	initialHP := player.HP

	player.DealDamage(2.0)

	if player.HP != initialHP-2.0 {
		t.Errorf("Expected HP %f, got %f", initialHP-2.0, player.HP)
	}
}

func TestPlayer_DealDamage_SmallDamage(t *testing.T) {
	player := testPlayer()
	// Player starts with 10 HP

	player.DealDamage(1.5)

	if player.HP != 8.5 {
		t.Errorf("Expected 8.5 HP, got %f", player.HP)
	}
}

func TestPlayer_DealDamage_LethalDamage(t *testing.T) {
	player := testPlayer()

	player.DealDamage(10.0)

	if player.HP != 0.0 {
		t.Errorf("Expected 0 HP, got %f", player.HP)
	}
}

func TestPlayer_DealDamage_OverDamage(t *testing.T) {
	player := testPlayer()

	// Deal more damage than current HP
	player.DealDamage(100.0)

	// Should clamp at 0, not go negative
	if player.HP != 0.0 {
		t.Errorf("Expected HP clamped at 0, got %f", player.HP)
	}
}

func TestPlayer_DealDamage_MultipleDamageInstances(t *testing.T) {
	player := testPlayer()
	// Player starts with 10 HP

	player.DealDamage(2.0)
	player.DealDamage(3.0)
	player.DealDamage(1.0)

	if player.HP != 4.0 {
		t.Errorf("Expected 4 HP after 3 damage instances, got %f", player.HP)
	}
}

func TestPlayer_DealDamage_AlreadyDead(t *testing.T) {
	player := testPlayer()

	// Kill player
	player.DealDamage(10.0)

	// Deal additional damage
	player.DealDamage(5.0)

	// Should remain at 0, not go negative
	if player.HP != 0.0 {
		t.Errorf("Expected HP to stay at 0 for dead player, got %f", player.HP)
	}
}

func TestPlayer_DealDamage_PartialDamage(t *testing.T) {
	player := testPlayer()
	// Player starts with 10 HP

	player.DealDamage(3.7)

	expectedHP := float32(10.0 - 3.7)
	if player.HP != expectedHP {
		t.Errorf("Expected %f HP, got %f", expectedHP, player.HP)
	}
}

func TestPlayer_DealDamage_ZeroDamage(t *testing.T) {
	player := testPlayer()
	initialHP := player.HP

	player.DealDamage(0.0)

	if player.HP != initialHP {
		t.Errorf("Expected no change with zero damage, got HP %f", player.HP)
	}
}
