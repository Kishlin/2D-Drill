package physics

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// Test helper - creates components from config for fall damage tests
func testFallDamagePlayerComponents() (entities.Hull, entities.Engine) {
	hullStats := config.HullStats{MaxHP: 10}
	engineStats := config.EngineStats{MaxSpeed: 450, Acceleration: 2500, FlyAcceleration: 2500, MaxUpwardSpeed: -600}

	hull := entities.NewHullFromConfig(0, "Base", hullStats)
	engine := entities.NewEngineFromConfig(0, "Base", engineStats)

	return hull, engine
}

func TestApplyFallDamage_BelowThreshold(t *testing.T) {
	hull, engine := testFallDamagePlayerComponents()
	player := &entities.Player{
		AABB:   types.NewAABB(0, 0, 64, 64),
		HP:     10.0,
		Hull:   hull,
		Engine: engine,
	}

	// Fall at 400 px/sec (below 500 threshold)
	ApplyFallDamage(player, 400.0)

	if player.HP != 10.0 {
		t.Errorf("Expected no damage below threshold, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_AtThreshold(t *testing.T) {
	hull, engine := testFallDamagePlayerComponents()
	player := &entities.Player{
		AABB:   types.NewAABB(0, 0, 64, 64),
		HP:     10.0,
		Hull:   hull,
		Engine: engine,
	}

	// Fall at exactly 500 px/sec (threshold)
	ApplyFallDamage(player, 500.0)

	if player.HP != 10.0 {
		t.Errorf("Expected no damage at threshold, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_SlightlyAboveThreshold(t *testing.T) {
	hull, engine := testFallDamagePlayerComponents()
	player := &entities.Player{
		AABB:   types.NewAABB(0, 0, 64, 64),
		HP:     10.0,
		Hull:   hull,
		Engine: engine,
	}

	// Fall at 520 px/sec: damage = (520 - 500) / 20 = 1.0
	ApplyFallDamage(player, 520.0)

	if player.HP != 9.0 {
		t.Errorf("Expected 1.0 damage, got HP: %f (damage: %f)", player.HP, 10.0-player.HP)
	}
}

func TestApplyFallDamage_ModerateFall(t *testing.T) {
	hull, engine := testFallDamagePlayerComponents()
	player := &entities.Player{
		AABB:   types.NewAABB(0, 0, 64, 64),
		HP:     10.0,
		Hull:   hull,
		Engine: engine,
	}

	// Fall at 600 px/sec: damage = (600 - 500) / 20 = 5.0
	ApplyFallDamage(player, 600.0)

	if player.HP != 5.0 {
		t.Errorf("Expected 5.0 damage, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_LethalFall(t *testing.T) {
	hull, engine := testFallDamagePlayerComponents()
	player := &entities.Player{
		AABB:   types.NewAABB(0, 0, 64, 64),
		HP:     10.0,
		Hull:   hull,
		Engine: engine,
	}

	// Fall at 700 px/sec: damage = (700 - 500) / 20 = 10.0 (lethal)
	ApplyFallDamage(player, 700.0)

	if player.HP != 0.0 {
		t.Errorf("Expected 0.0 HP (clamped), got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_ExtremeVelocity(t *testing.T) {
	hull, engine := testFallDamagePlayerComponents()
	player := &entities.Player{
		AABB:   types.NewAABB(0, 0, 64, 64),
		HP:     10.0,
		Hull:   hull,
		Engine: engine,
	}

	// Fall at 1500 px/sec: damage = (1500 - 500) / 20 = 50.0
	// But HP should clamp at 0
	ApplyFallDamage(player, 1500.0)

	if player.HP != 0.0 {
		t.Errorf("Expected HP clamped at 0, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_PreservesPartialHealth(t *testing.T) {
	hull, engine := testFallDamagePlayerComponents()
	player := &entities.Player{
		AABB:   types.NewAABB(0, 0, 64, 64),
		HP:     8.0, // Damaged player
		Hull:   hull,
		Engine: engine,
	}

	// Fall at 600 px/sec: damage = (600 - 500) / 20 = 5.0
	// Should reduce to 3.0, not clamp to 0
	ApplyFallDamage(player, 600.0)

	if player.HP != 3.0 {
		t.Errorf("Expected 3.0 HP, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_AlreadyDead(t *testing.T) {
	hull, engine := testFallDamagePlayerComponents()
	player := &entities.Player{
		AABB:   types.NewAABB(0, 0, 64, 64),
		HP:     0.0, // Already dead
		Hull:   hull,
		Engine: engine,
	}

	// Fall at 600 px/sec
	ApplyFallDamage(player, 600.0)

	// Should remain at 0, not go negative
	if player.HP != 0.0 {
		t.Errorf("Expected 0.0 HP for dead player, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_NegativeVelocity(t *testing.T) {
	hull, engine := testFallDamagePlayerComponents()
	player := &entities.Player{
		AABB:   types.NewAABB(0, 0, 64, 64),
		HP:     10.0,
		Hull:   hull,
		Engine: engine,
	}

	// Negative velocity (moving upward) - should not apply damage
	ApplyFallDamage(player, -600.0)

	if player.HP != 10.0 {
		t.Errorf("Expected no damage for upward movement, got HP: %f", player.HP)
	}
}

func TestApplyFallDamage_ZeroVelocity(t *testing.T) {
	hull, engine := testFallDamagePlayerComponents()
	player := &entities.Player{
		AABB:   types.NewAABB(0, 0, 64, 64),
		HP:     10.0,
		Hull:   hull,
		Engine: engine,
	}

	// Zero velocity - no damage
	ApplyFallDamage(player, 0.0)

	if player.HP != 10.0 {
		t.Errorf("Expected no damage at zero velocity, got HP: %f", player.HP)
	}
}
