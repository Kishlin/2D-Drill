package physics

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
)

func testFallDamagePlayer() *entities.Player {
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
	return entities.NewPlayerFromConfig(0, 0, playerCfg, upgradeCfg)
}

func TestApplyFallDamage_BelowThreshold(t *testing.T) {
	player := testFallDamagePlayer()
	player.HP = 10.0

	// Fall at 400 px/sec (below 500 threshold)
	ApplyFallDamage(player, 400.0)

	if player.HP != 10.0 {
		t.Errorf("Expected no damage below threshold, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_AtThreshold(t *testing.T) {
	player := testFallDamagePlayer()
	player.HP = 10.0

	// Fall at exactly 500 px/sec (threshold)
	ApplyFallDamage(player, 500.0)

	if player.HP != 10.0 {
		t.Errorf("Expected no damage at threshold, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_SlightlyAboveThreshold(t *testing.T) {
	player := testFallDamagePlayer()
	player.HP = 10.0

	// Fall at 520 px/sec: damage = (520 - 500) / 20 = 1.0
	ApplyFallDamage(player, 520.0)

	if player.HP != 9.0 {
		t.Errorf("Expected 1.0 damage, got HP: %f (damage: %f)", player.HP, 10.0-player.HP)
	}
}

func TestApplyFallDamage_ModerateFall(t *testing.T) {
	player := testFallDamagePlayer()
	player.HP = 10.0

	// Fall at 600 px/sec: damage = (600 - 500) / 20 = 5.0
	ApplyFallDamage(player, 600.0)

	if player.HP != 5.0 {
		t.Errorf("Expected 5.0 damage, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_LethalFall(t *testing.T) {
	player := testFallDamagePlayer()
	player.HP = 10.0

	// Fall at 700 px/sec: damage = (700 - 500) / 20 = 10.0 (lethal)
	ApplyFallDamage(player, 700.0)

	if player.HP != 0.0 {
		t.Errorf("Expected 0.0 HP (clamped), got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_ExtremeVelocity(t *testing.T) {
	player := testFallDamagePlayer()
	player.HP = 10.0

	// Fall at 1500 px/sec: damage = (1500 - 500) / 20 = 50.0
	// But HP should clamp at 0
	ApplyFallDamage(player, 1500.0)

	if player.HP != 0.0 {
		t.Errorf("Expected HP clamped at 0, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_PreservesPartialHealth(t *testing.T) {
	player := testFallDamagePlayer()
	player.HP = 8.0 // Damaged player

	// Fall at 600 px/sec: damage = (600 - 500) / 20 = 5.0
	// Should reduce to 3.0, not clamp to 0
	ApplyFallDamage(player, 600.0)

	if player.HP != 3.0 {
		t.Errorf("Expected 3.0 HP, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_AlreadyDead(t *testing.T) {
	player := testFallDamagePlayer()
	player.HP = 0.0 // Already dead

	// Fall at 600 px/sec
	ApplyFallDamage(player, 600.0)

	// Should remain at 0, not go negative
	if player.HP != 0.0 {
		t.Errorf("Expected 0.0 HP for dead player, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_NegativeVelocity(t *testing.T) {
	player := testFallDamagePlayer()
	player.HP = 10.0

	// Negative velocity (moving upward) - should not apply damage
	ApplyFallDamage(player, -600.0)

	if player.HP != 10.0 {
		t.Errorf("Expected no damage for upward movement, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_ZeroVelocity(t *testing.T) {
	player := testFallDamagePlayer()
	player.HP = 10.0

	// Zero velocity - no damage
	ApplyFallDamage(player, 0.0)

	if player.HP != 10.0 {
		t.Errorf("Expected no damage at zero velocity, got HP: %f", player.HP)
	}
}
