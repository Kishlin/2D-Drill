package systems

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// Test constants matching the original hardcoded values
const (
	testGroundLevel  = 640.0
	testWorldHeight  = 64000.0
)

// Test helper - creates player with specified heat resistance
func testHeatPlayer(heatResistance float32) *entities.Player {
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
		HeatShields: []config.UpgradeTier[config.HeatShieldStats]{{Name: "Base", Price: 0, Stats: config.HeatShieldStats{HeatResistance: heatResistance}}},
		Drills:      []config.UpgradeTier[config.DrillStats]{{Name: "Base", Price: 0, Stats: config.DrillStats{SpeedAtSurface: 1.0, SpeedAtMaxDepth: 1.0}}},
	}
	return entities.NewPlayerFromConfig(0, 0, playerCfg, upgradeCfg)
}

// Test helper - creates a test world with standard dimensions for heat tests
func testHeatWorld() *world.World {
	worldCfg := config.WorldConfig{
		Width:       1280,
		Height:      testWorldHeight,
		GroundLevel: testGroundLevel,
		Seed:        12345,
	}
	genCfg := config.GenerationConfig{}
	bossRoomCfg := config.BossRoomConfig{
		BossType:    "test_boss",
		FloorType:   config.FloorConcrete,
		RoomHeight:  680.0,
		FloorHeight: 6.0,
	}
	return world.NewWorldFromConfig(worldCfg, genCfg, bossRoomCfg)
}

// TemperatureTests

func TestCalculateTemperature_AtGroundLevel(t *testing.T) {
	// At ground level (Y=640), temperature should be 15°C
	temp := CalculateTemperature(640.0, testGroundLevel, testWorldHeight)

	if temp != 15.0 {
		t.Errorf("Expected 15°C at ground level, got %f°C", temp)
	}
}

func TestCalculateTemperature_AboveGround(t *testing.T) {
	// Above ground (Y < 640), temperature should be clamped to base (15°C)
	temp := CalculateTemperature(500.0, testGroundLevel, testWorldHeight)

	if temp != 15.0 {
		t.Errorf("Expected 15°C above ground, got %f°C", temp)
	}
}

func TestCalculateTemperature_AtMaxDepth(t *testing.T) {
	// At max depth (Y=64000), temperature should be 350°C
	temp := CalculateTemperature(64000.0, testGroundLevel, testWorldHeight)

	if temp != 350.0 {
		t.Errorf("Expected 350°C at max depth, got %f°C", temp)
	}
}

func TestCalculateTemperature_Midpoint(t *testing.T) {
	// Midpoint between ground and max depth
	// Y = (640 + 64000) / 2 = 32320
	// temp = 15 + 0.5 * (350 - 15) = 15 + 167.5 = 182.5°C
	temp := CalculateTemperature(32320.0, testGroundLevel, testWorldHeight)

	const expected = 182.5
	const tolerance = 0.1
	if temp < expected-tolerance || temp > expected+tolerance {
		t.Errorf("Expected ~182.5°C at midpoint, got %f°C", temp)
	}
}

func TestCalculateTemperature_OneQuarter(t *testing.T) {
	// One quarter down: 640 + 0.25 * (64000 - 640) = 640 + 15840 = 16480
	// temp = 15 + 0.25 * 335 = 15 + 83.75 = 98.75°C
	temp := CalculateTemperature(16480.0, testGroundLevel, testWorldHeight)

	const expected = 98.75
	const tolerance = 0.1
	if temp < expected-tolerance || temp > expected+tolerance {
		t.Errorf("Expected ~98.75°C at 1/4 depth, got %f°C", temp)
	}
}

func TestCalculateTemperature_ThreeQuarters(t *testing.T) {
	// Three quarters down: 640 + 0.75 * (64000 - 640) = 640 + 47520 = 48160
	// temp = 15 + 0.75 * 335 = 15 + 251.25 = 266.25°C
	temp := CalculateTemperature(48160.0, testGroundLevel, testWorldHeight)

	const expected = 266.25
	const tolerance = 0.1
	if temp < expected-tolerance || temp > expected+tolerance {
		t.Errorf("Expected ~266.25°C at 3/4 depth, got %f°C", temp)
	}
}

// Heat Damage Tests

func TestUpdateHeat_NoExcessHeat(t *testing.T) {
	player := testHeatPlayer(50) // 50°C resistance
	player.AABB.Y = 640          // At ground level (15°C)
	player.HP = 10.0
	w := testHeatWorld()

	// Temperature 15°C < resistance 50°C, no damage
	UpdateHeat(player, w, 0.016) // ~60 FPS

	if player.HP != 10.0 {
		t.Errorf("Expected no damage when below resistance, got HP: %f", player.HP)
	}
}

func TestUpdateHeat_AtResistanceLimit(t *testing.T) {
	// At 640 + 800 = 1440px, temp = 15 + (800/63360) * 335 ≈ 19.24°C
	// Resistance 50°C > temp, no damage
	player := testHeatPlayer(50)
	player.AABB.Y = 1440
	player.HP = 10.0
	w := testHeatWorld()

	UpdateHeat(player, w, 0.016)

	if player.HP != 10.0 {
		t.Errorf("Expected no damage within resistance margin, got HP: %f", player.HP)
	}
}

func TestUpdateHeat_SlightExcess(t *testing.T) {
	// At 640 + 6650 = 7290px, temp ≈ 50.35°C
	// Base resistance = 50°C
	// Excess = 0.35°C (minimal)
	// damage = 0.5 * (0.35/10)^1.5 * dt
	player := testHeatPlayer(50)
	player.AABB.Y = 7290
	player.HP = 10.0
	w := testHeatWorld()

	UpdateHeat(player, w, 1.0) // 1 second

	if player.HP >= 10.0 {
		t.Errorf("Expected some damage with excess heat, got HP: %f", player.HP)
	}
	if player.HP < 9.95 {
		t.Errorf("Expected minimal damage, got HP: %f", player.HP)
	}
}

func TestUpdateHeat_SignificantExcess(t *testing.T) {
	// At 640 + 6650*3 = 20590px, temp ≈ 116.04°C
	// Base resistance = 50°C
	// Excess = 66.04°C
	// damage/sec ≈ 0.5 * (66.04/10)^1.5 ≈ 0.5 * ~17 ≈ 8.5 HP/sec
	player := testHeatPlayer(50)
	player.AABB.Y = 20590
	player.HP = 10.0
	w := testHeatWorld()

	UpdateHeat(player, w, 1.0) // 1 second

	// Should take significant damage
	if player.HP >= 5.0 {
		t.Errorf("Expected ~8+ damage per second at this depth, got HP: %f", player.HP)
	}
	if player.HP < 0.0 {
		t.Errorf("HP should be clamped at 0, got HP: %f", player.HP)
	}
}

func TestUpdateHeat_ClampsAtZero(t *testing.T) {
	// Very deep: temperature far exceeds resistance
	player := testHeatPlayer(50) // 50°C resistance
	player.AABB.Y = 64000        // Max depth (350°C)
	player.HP = 10.0
	w := testHeatWorld()

	// Apply 10 seconds of heat damage
	UpdateHeat(player, w, 10.0)

	if player.HP != 0.0 {
		t.Errorf("Expected HP clamped at 0, got HP: %f", player.HP)
	}
}

func TestUpdateHeat_UpgradedShield(t *testing.T) {
	// At depth with 140°C temp, Mk2 shield (140°C resistance) should take minimal/no damage
	// Y = 640 + (140-15)/335 * 63360 = 640 + 23647 ≈ 24287
	player := testHeatPlayer(140) // 140°C resistance (Mk2)
	player.AABB.Y = 24287
	player.HP = 10.0
	w := testHeatWorld()

	UpdateHeat(player, w, 0.016) // One frame at 60 FPS

	// At this depth, temp ≈ 140°C, resistance = 140°C
	// Allow for floating-point tolerance (tiny rounding errors)
	if player.HP < 9.999 {
		t.Errorf("Expected minimal/no damage at resistance limit with upgraded shield, got HP: %f", player.HP)
	}
}

func TestUpdateHeat_ScalesWithDeltaTime(t *testing.T) {
	// Same depth, test that damage scales with dt
	depth := float32(20590.0) // ~116°C temp
	w := testHeatWorld()

	player1 := testHeatPlayer(50)
	player1.AABB.Y = depth
	player1.HP = 10.0

	player2 := testHeatPlayer(50)
	player2.AABB.Y = depth
	player2.HP = 10.0

	UpdateHeat(player1, w, 0.5) // Half second
	UpdateHeat(player2, w, 1.0) // Full second

	// Damage should roughly double with 2x delta time
	damage1 := 10.0 - player1.HP
	damage2 := 10.0 - player2.HP

	if damage2 < damage1*1.8 || damage2 > damage1*2.2 {
		t.Errorf("Damage should scale linearly with dt. 0.5s: %f, 1.0s: %f", damage1, damage2)
	}
}

func TestUpdateHeat_AlreadyDead(t *testing.T) {
	player := testHeatPlayer(50)
	player.AABB.Y = 64000
	player.HP = 0.0 // Already dead
	w := testHeatWorld()

	UpdateHeat(player, w, 10.0)

	// Should remain at 0, not go negative
	if player.HP != 0.0 {
		t.Errorf("Expected 0.0 HP for dead player, got HP: %f", player.HP)
	}
}

func TestUpdateHeat_PreservesPartialHealth(t *testing.T) {
	// Player at 8 HP with excess heat
	player := testHeatPlayer(50)
	player.AABB.Y = 20590 // ~116°C
	player.HP = 8.0       // Damaged
	w := testHeatWorld()

	UpdateHeat(player, w, 0.5)

	// Should reduce proportionally but not clamp to 0 if still above 0
	if player.HP < 0.0 || player.HP >= 8.0 {
		t.Errorf("Expected partial health reduction, got HP: %f", player.HP)
	}
}

func TestUpdateHeat_MoreDamageThanExposure(t *testing.T) {
	// Verify that deeper (hotter) locations take more damage
	// Shallow location: temp ≈ 75°C, excess = 25°C
	// Deep location: temp ≈ 200°C, excess = 150°C
	w := testHeatWorld()

	shallowPlayer := testHeatPlayer(50)
	shallowPlayer.AABB.Y = 6650 // Shallow depth
	shallowPlayer.HP = 10.0

	deepPlayer := testHeatPlayer(50)
	deepPlayer.AABB.Y = 30000 // Deeper depth
	deepPlayer.HP = 10.0

	UpdateHeat(shallowPlayer, w, 1.0)
	UpdateHeat(deepPlayer, w, 1.0)

	shallowDamage := 10.0 - shallowPlayer.HP
	deepDamage := 10.0 - deepPlayer.HP

	// Deeper location should take significantly more damage
	if deepDamage <= shallowDamage {
		t.Errorf("Expected more damage at deeper location. Shallow: %f, Deep: %f", shallowDamage, deepDamage)
	}
}
