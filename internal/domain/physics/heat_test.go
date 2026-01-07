package physics

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// TemperatureTests

func TestCalculateTemperature_AtGroundLevel(t *testing.T) {
	// At ground level (Y=640), temperature should be 15°C
	temp := CalculateTemperature(640.0)

	if temp != 15.0 {
		t.Errorf("Expected 15°C at ground level, got %f°C", temp)
	}
}

func TestCalculateTemperature_AboveGround(t *testing.T) {
	// Above ground (Y < 640), temperature should be clamped to base (15°C)
	temp := CalculateTemperature(500.0)

	if temp != 15.0 {
		t.Errorf("Expected 15°C above ground, got %f°C", temp)
	}
}

func TestCalculateTemperature_AtMaxDepth(t *testing.T) {
	// At max depth (Y=64000), temperature should be 350°C
	temp := CalculateTemperature(64000.0)

	if temp != 350.0 {
		t.Errorf("Expected 350°C at max depth, got %f°C", temp)
	}
}

func TestCalculateTemperature_Midpoint(t *testing.T) {
	// Midpoint between ground and max depth
	// Y = (640 + 64000) / 2 = 32320
	// temp = 15 + 0.5 * (350 - 15) = 15 + 167.5 = 182.5°C
	temp := CalculateTemperature(32320.0)

	const expected = 182.5
	const tolerance = 0.1
	if temp < expected-tolerance || temp > expected+tolerance {
		t.Errorf("Expected ~182.5°C at midpoint, got %f°C", temp)
	}
}

func TestCalculateTemperature_OneQuarter(t *testing.T) {
	// One quarter down: 640 + 0.25 * (64000 - 640) = 640 + 15840 = 16480
	// temp = 15 + 0.25 * 335 = 15 + 83.75 = 98.75°C
	temp := CalculateTemperature(16480.0)

	const expected = 98.75
	const tolerance = 0.1
	if temp < expected-tolerance || temp > expected+tolerance {
		t.Errorf("Expected ~98.75°C at 1/4 depth, got %f°C", temp)
	}
}

func TestCalculateTemperature_ThreeQuarters(t *testing.T) {
	// Three quarters down: 640 + 0.75 * (64000 - 640) = 640 + 47520 = 48160
	// temp = 15 + 0.75 * 335 = 15 + 251.25 = 266.25°C
	temp := CalculateTemperature(48160.0)

	const expected = 266.25
	const tolerance = 0.1
	if temp < expected-tolerance || temp > expected+tolerance {
		t.Errorf("Expected ~266.25°C at 3/4 depth, got %f°C", temp)
	}
}

// Heat Damage Tests

func TestApplyHeatDamage_NoExcessHeat(t *testing.T) {
	player := &entities.Player{
		AABB:      types.NewAABB(0, 640, 64, 64), // At ground level (15°C)
		HP:        10.0,
		Hull:      entities.NewHullBase(),
		Engine:    entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldBase(), // 50°C resistance
	}

	// Temperature 15°C < resistance 50°C, no damage
	ApplyHeatDamage(player, 0.016) // ~60 FPS

	if player.HP != 10.0 {
		t.Errorf("Expected no damage when below resistance, got HP: %f", player.HP)
	}
}

func TestApplyHeatDamage_AtResistanceLimit(t *testing.T) {
	// At 640 + 800 = 1440px, temp = 15 + (800/63360) * 335 ≈ 19.24°C
	// Resistance 50°C > temp, no damage
	player := &entities.Player{
		AABB:       types.NewAABB(0, 1440, 64, 64),
		HP:         10.0,
		Hull:       entities.NewHullBase(),
		Engine:     entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldBase(),
	}

	ApplyHeatDamage(player, 0.016)

	if player.HP != 10.0 {
		t.Errorf("Expected no damage within resistance margin, got HP: %f", player.HP)
	}
}

func TestApplyHeatDamage_SlightExcess(t *testing.T) {
	// At 640 + 6650 = 7290px, temp ≈ 50.35°C
	// Base resistance = 50°C
	// Excess = 0.35°C (minimal)
	// damage = 0.5 * (0.35/10)^1.5 * dt
	player := &entities.Player{
		AABB:       types.NewAABB(0, 7290, 64, 64),
		HP:         10.0,
		Hull:       entities.NewHullBase(),
		Engine:     entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldBase(),
	}

	ApplyHeatDamage(player, 1.0) // 1 second

	if player.HP >= 10.0 {
		t.Errorf("Expected some damage with excess heat, got HP: %f", player.HP)
	}
	if player.HP < 9.95 {
		t.Errorf("Expected minimal damage, got HP: %f", player.HP)
	}
}

func TestApplyHeatDamage_SignificantExcess(t *testing.T) {
	// At 640 + 6650*3 = 20590px, temp ≈ 116.04°C
	// Base resistance = 50°C
	// Excess = 66.04°C
	// damage/sec ≈ 0.5 * (66.04/10)^1.5 ≈ 0.5 * ~17 ≈ 8.5 HP/sec
	player := &entities.Player{
		AABB:       types.NewAABB(0, 20590, 64, 64),
		HP:         10.0,
		Hull:       entities.NewHullBase(),
		Engine:     entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldBase(),
	}

	ApplyHeatDamage(player, 1.0) // 1 second

	// Should take significant damage
	if player.HP >= 5.0 {
		t.Errorf("Expected ~8+ damage per second at this depth, got HP: %f", player.HP)
	}
	if player.HP < 0.0 {
		t.Errorf("HP should be clamped at 0, got HP: %f", player.HP)
	}
}

func TestApplyHeatDamage_ClampsAtZero(t *testing.T) {
	// Very deep: temperature far exceeds resistance
	player := &entities.Player{
		AABB:       types.NewAABB(0, 64000, 64, 64), // Max depth (350°C)
		HP:         10.0,
		Hull:       entities.NewHullBase(),
		Engine:     entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldBase(), // 50°C resistance
	}

	// Apply 10 seconds of heat damage
	ApplyHeatDamage(player, 10.0)

	if player.HP != 0.0 {
		t.Errorf("Expected HP clamped at 0, got HP: %f", player.HP)
	}
}

func TestApplyHeatDamage_UpgradedShield(t *testing.T) {
	// At depth with 140°C temp, Mk2 shield (140°C resistance) should take minimal/no damage
	// Y = 640 + (140-15)/335 * 63360 = 640 + 23647 ≈ 24287
	player := &entities.Player{
		AABB:       types.NewAABB(0, 24287, 64, 64),
		HP:         10.0,
		Hull:       entities.NewHullBase(),
		Engine:     entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldMk2(), // 140°C resistance
	}

	ApplyHeatDamage(player, 0.016) // One frame at 60 FPS

	// At this depth, temp ≈ 140°C, resistance = 140°C
	// Allow for floating-point tolerance (tiny rounding errors)
	if player.HP < 9.999 {
		t.Errorf("Expected minimal/no damage at resistance limit with upgraded shield, got HP: %f", player.HP)
	}
}

func TestApplyHeatDamage_ScalesWithDeltaTime(t *testing.T) {
	// Same depth, test that damage scales with dt
	depth := float32(20590.0) // ~116°C temp

	player1 := &entities.Player{
		AABB:       types.NewAABB(0, depth, 64, 64),
		HP:         10.0,
		Hull:       entities.NewHullBase(),
		Engine:     entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldBase(),
	}

	player2 := &entities.Player{
		AABB:       types.NewAABB(0, depth, 64, 64),
		HP:         10.0,
		Hull:       entities.NewHullBase(),
		Engine:     entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldBase(),
	}

	ApplyHeatDamage(player1, 0.5)  // Half second
	ApplyHeatDamage(player2, 1.0)  // Full second

	// Damage should roughly double with 2x delta time
	damage1 := 10.0 - player1.HP
	damage2 := 10.0 - player2.HP

	if damage2 < damage1*1.8 || damage2 > damage1*2.2 {
		t.Errorf("Damage should scale linearly with dt. 0.5s: %f, 1.0s: %f", damage1, damage2)
	}
}

func TestApplyHeatDamage_AlreadyDead(t *testing.T) {
	player := &entities.Player{
		AABB:       types.NewAABB(0, 64000, 64, 64),
		HP:         0.0, // Already dead
		Hull:       entities.NewHullBase(),
		Engine:     entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldBase(),
	}

	ApplyHeatDamage(player, 10.0)

	// Should remain at 0, not go negative
	if player.HP != 0.0 {
		t.Errorf("Expected 0.0 HP for dead player, got HP: %f", player.HP)
	}
}

func TestApplyHeatDamage_PreservesPartialHealth(t *testing.T) {
	// Player at 8 HP with excess heat
	player := &entities.Player{
		AABB:       types.NewAABB(0, 20590, 64, 64), // ~116°C
		HP:         8.0, // Damaged
		Hull:       entities.NewHullBase(),
		Engine:     entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldBase(),
	}

	ApplyHeatDamage(player, 0.5)

	// Should reduce proportionally but not clamp to 0 if still above 0
	if player.HP < 0.0 || player.HP >= 8.0 {
		t.Errorf("Expected partial health reduction, got HP: %f", player.HP)
	}
}

func TestApplyHeatDamage_MoreDamageThanExposure(t *testing.T) {
	// Verify that deeper (hotter) locations take more damage
	// Shallow location: temp ≈ 75°C, excess = 25°C
	// Deep location: temp ≈ 200°C, excess = 150°C

	shallowPlayer := &entities.Player{
		AABB:       types.NewAABB(0, float32(6650), 64, 64), // Shallow depth
		HP:         10.0,
		Hull:       entities.NewHullBase(),
		Engine:     entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldBase(),
	}

	deepPlayer := &entities.Player{
		AABB:       types.NewAABB(0, float32(30000), 64, 64), // Deeper depth
		HP:         10.0,
		Hull:       entities.NewHullBase(),
		Engine:     entities.NewEngineBase(),
		HeatShield: entities.NewHeatShieldBase(),
	}

	ApplyHeatDamage(shallowPlayer, 1.0)
	ApplyHeatDamage(deepPlayer, 1.0)

	shallowDamage := 10.0 - shallowPlayer.HP
	deepDamage := 10.0 - deepPlayer.HP

	// Deeper location should take significantly more damage
	if deepDamage <= shallowDamage {
		t.Errorf("Expected more damage at deeper location. Shallow: %f, Deep: %f", shallowDamage, deepDamage)
	}
}
