package config_test

import (
	"strings"
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
)

func TestGameConfig_Validate_ValidConfig(t *testing.T) {
	cfg := validGameConfig()
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Valid config should not return error, got: %v", err)
	}
}

func TestGameConfig_Validate_InvalidWorldConfig(t *testing.T) {
	cfg := validGameConfig()
	cfg.World.Width = -100 // Invalid

	err := cfg.Validate()
	if err == nil {
		t.Error("Invalid world config should return error")
	}
	if strings.Contains(err.Error(), "world config") == false {
		t.Errorf("Error should mention world config, got: %v", err)
	}
}

func TestGameConfig_Validate_StartingEngineTierExceedsAvailable(t *testing.T) {
	cfg := validGameConfig()
	cfg.Player.StartingUpgrades.Engine = 10 // Only 1 tier available

	err := cfg.Validate()
	if err == nil {
		t.Error("Starting engine tier exceeding available should return error")
	}
	if strings.Contains(err.Error(), "engine tier") == false {
		t.Errorf("Error should mention engine tier, got: %v", err)
	}
}

func TestGameConfig_Validate_StartingHullTierExceedsAvailable(t *testing.T) {
	cfg := validGameConfig()
	cfg.Player.StartingUpgrades.Hull = 10

	err := cfg.Validate()
	if err == nil {
		t.Error("Starting hull tier exceeding available should return error")
	}
	if strings.Contains(err.Error(), "hull tier") == false {
		t.Errorf("Error should mention hull tier, got: %v", err)
	}
}

func TestGameConfig_Validate_StartingFuelTankTierExceedsAvailable(t *testing.T) {
	cfg := validGameConfig()
	cfg.Player.StartingUpgrades.FuelTank = 10

	err := cfg.Validate()
	if err == nil {
		t.Error("Starting fuel tank tier exceeding available should return error")
	}
	if strings.Contains(err.Error(), "fuel tank tier") == false {
		t.Errorf("Error should mention fuel tank tier, got: %v", err)
	}
}

func TestGameConfig_Validate_StartingCargoHoldTierExceedsAvailable(t *testing.T) {
	cfg := validGameConfig()
	cfg.Player.StartingUpgrades.CargoHold = 10

	err := cfg.Validate()
	if err == nil {
		t.Error("Starting cargo hold tier exceeding available should return error")
	}
	if strings.Contains(err.Error(), "cargo hold tier") == false {
		t.Errorf("Error should mention cargo hold tier, got: %v", err)
	}
}

func TestGameConfig_Validate_StartingHeatShieldTierExceedsAvailable(t *testing.T) {
	cfg := validGameConfig()
	cfg.Player.StartingUpgrades.HeatShield = 10

	err := cfg.Validate()
	if err == nil {
		t.Error("Starting heat shield tier exceeding available should return error")
	}
	if strings.Contains(err.Error(), "heat shield tier") == false {
		t.Errorf("Error should mention heat shield tier, got: %v", err)
	}
}

func TestGameConfig_Validate_StartingDrillTierExceedsAvailable(t *testing.T) {
	cfg := validGameConfig()
	cfg.Player.StartingUpgrades.Drill = 10

	err := cfg.Validate()
	if err == nil {
		t.Error("Starting drill tier exceeding available should return error")
	}
	if strings.Contains(err.Error(), "drill tier") == false {
		t.Errorf("Error should mention drill tier, got: %v", err)
	}
}

func TestGameConfig_Validate_NoOres(t *testing.T) {
	cfg := validGameConfig()
	cfg.Generation.Ores = []config.OreConfig{}

	err := cfg.Validate()
	if err == nil {
		t.Error("Config with no ores should return error")
	}
	if strings.Contains(err.Error(), "at least one ore") == false {
		t.Errorf("Error should mention ore requirement, got: %v", err)
	}
}

func TestGameConfig_Validate_DuplicateOreIDs(t *testing.T) {
	cfg := validGameConfig()
	cfg.Generation.Ores = []config.OreConfig{
		{ID: "copper", Name: "Copper", Value: 25, Hardness: 1.2, Distribution: config.TileDistribution{PeakDepth: 0, Sigma: 100, MaxWeight: 5}},
		{ID: "copper", Name: "Another Copper", Value: 50, Hardness: 1.5, Distribution: config.TileDistribution{PeakDepth: 100, Sigma: 100, MaxWeight: 5}},
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Config with duplicate ore IDs should return error")
	}
	if strings.Contains(err.Error(), "duplicate ore ID") == false {
		t.Errorf("Error should mention duplicate ore ID, got: %v", err)
	}
}

func TestGameConfig_Validate_DuplicateHazardIDs(t *testing.T) {
	cfg := validGameConfig()
	cfg.Generation.Hazards = []config.HazardConfig{
		{ID: "rock", Name: "Rock", Drillable: false, Distribution: config.TileDistribution{PeakDepth: 500, Sigma: 100, MaxWeight: 10}},
		{ID: "rock", Name: "Another Rock", Drillable: false, Distribution: config.TileDistribution{PeakDepth: 600, Sigma: 100, MaxWeight: 10}},
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Config with duplicate hazard IDs should return error")
	}
	if strings.Contains(err.Error(), "duplicate hazard ID") == false {
		t.Errorf("Error should mention duplicate hazard ID, got: %v", err)
	}
}

func TestGameConfig_Validate_NoHazardsIsValid(t *testing.T) {
	cfg := validGameConfig()
	cfg.Generation.Hazards = []config.HazardConfig{}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Config with no hazards should be valid, got: %v", err)
	}
}

func TestGameConfig_Validate_StartingTierAtMaxValid(t *testing.T) {
	cfg := validGameConfig()
	// Add multiple tiers and set starting to the last one
	cfg.Upgrades.Engines = []config.UpgradeTier[config.EngineStats]{
		{Name: "Base", Price: 0, Stats: config.EngineStats{MaxSpeed: 400}},
		{Name: "Mk1", Price: 100, Stats: config.EngineStats{MaxSpeed: 500}},
		{Name: "Mk2", Price: 200, Stats: config.EngineStats{MaxSpeed: 600}},
	}
	cfg.Player.StartingUpgrades.Engine = 2 // Last tier (index 2)

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Starting at max tier should be valid, got: %v", err)
	}
}

func TestGameConfig_Validate_MultipleOresUniqueIDs(t *testing.T) {
	cfg := validGameConfig()
	cfg.Generation.Ores = []config.OreConfig{
		{ID: "copper", Name: "Copper", Value: 25, Hardness: 1.2, Distribution: config.TileDistribution{PeakDepth: 0, Sigma: 100, MaxWeight: 5}},
		{ID: "iron", Name: "Iron", Value: 50, Hardness: 1.5, Distribution: config.TileDistribution{PeakDepth: 100, Sigma: 100, MaxWeight: 5}},
		{ID: "gold", Name: "Gold", Value: 200, Hardness: 2.0, Distribution: config.TileDistribution{PeakDepth: 200, Sigma: 100, MaxWeight: 3}},
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Multiple ores with unique IDs should be valid, got: %v", err)
	}
}

func TestGameConfig_Validate_MultipleHazardsUniqueIDs(t *testing.T) {
	cfg := validGameConfig()
	cfg.Generation.Hazards = []config.HazardConfig{
		{ID: "rock", Name: "Rock", Drillable: false, Distribution: config.TileDistribution{PeakDepth: 500, Sigma: 100, MaxWeight: 10}},
		{ID: "lava", Name: "Lava", Drillable: true, FixedDuration: 0.3, OnDrillEffect: config.HazardEffectConfig{Type: config.HazardEffectHeatDamage, BaseDamage: 100, MaxHeatResistance: 320, MaxDamageReduction: 0.5}, Distribution: config.TileDistribution{PeakDepth: 600, Sigma: 100, MaxWeight: 8}},
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Multiple hazards with unique IDs should be valid, got: %v", err)
	}
}

// Test helpers

func validGameConfig() *config.GameConfig {
	return &config.GameConfig{
		World: config.WorldConfig{
			Width:       2360,
			Height:      64000,
			GroundLevel: 640,
			Seed:        42,
			PlayerSpawn: config.PlayerSpawn{X: 1180, Y: 570},
			BuildingLayout: config.BuildingLayout{
				HospitalX:    100,
				FuelStationX: 460,
				MarketX:      820,
				UpgradeShopX: 1180,
				ItemShopX:    1540,
			},
		},
		Player: config.PlayerConfig{
			StartingMoney:    0,
			StartingItems:    [5]int{0, 0, 0, 0, 0},
			StartingUpgrades: config.StartingUpgrades{},
		},
		Generation: config.GenerationConfig{
			Empty:        config.TileDistribution{PeakDepth: 0, Sigma: 1000, MaxWeight: 20},
			Dirt:         config.TileDistribution{PeakDepth: 0, Sigma: 500, MaxWeight: 100},
			DirtHardness: 1.0,
			Ores: []config.OreConfig{
				{ID: "copper", Name: "Copper", Value: 25, Hardness: 1.2, Distribution: config.TileDistribution{PeakDepth: 0, Sigma: 100, MaxWeight: 5}, Color: [4]uint8{184, 115, 51, 255}},
			},
			Hazards: []config.HazardConfig{},
		},
		Upgrades: config.UpgradeConfig{
			Engines:     []config.UpgradeTier[config.EngineStats]{{Name: "Base", Price: 0, Stats: config.EngineStats{MaxSpeed: 450, Acceleration: 2500, FlyAcceleration: 2500, MaxUpwardSpeed: -600}}},
			Hulls:       []config.UpgradeTier[config.HullStats]{{Name: "Base", Price: 0, Stats: config.HullStats{MaxHP: 10}}},
			FuelTanks:   []config.UpgradeTier[config.FuelTankStats]{{Name: "Base", Price: 0, Stats: config.FuelTankStats{Capacity: 10}}},
			CargoHolds:  []config.UpgradeTier[config.CargoHoldStats]{{Name: "Base", Price: 0, Stats: config.CargoHoldStats{Capacity: 10}}},
			HeatShields: []config.UpgradeTier[config.HeatShieldStats]{{Name: "Base", Price: 0, Stats: config.HeatShieldStats{HeatResistance: 50}}},
			Drills:      []config.UpgradeTier[config.DrillStats]{{Name: "Base", Price: 0, Stats: config.DrillStats{SpeedAtSurface: 1.0, SpeedAtMaxDepth: 1.0}}},
		},
		Items: config.ItemConfig{
			Teleport: config.ItemEntry{Price: 500, Radius: 0, Damage: 0},
			Repair:   config.ItemEntry{Price: 200, Radius: 0, Damage: 0},
			Refuel:   config.ItemEntry{Price: 100, Radius: 0, Damage: 0},
			Bomb:     config.ItemEntry{Price: 300, Radius: 1, Damage: 10.0},
			BigBomb:  config.ItemEntry{Price: 800, Radius: 2, Damage: 25.0},
		},
		Level: config.LevelConfig{
			Number: 1,
			Name:   "Test Level",
			BossRoom: config.BossRoomConfig{
				BossType:    "test_boss",
				FloorType:   config.FloorConcrete,
				RoomHeight:  680.0,
				FloorHeight: 6.0,
			},
		},
		Drilling: config.DrillingConfig{
			MinDrillingDuration:   1.0,
			MaxDrillingDuration:   24.0,
			FloorDrillingDuration: 0.5,
		},
		Heat: config.HeatConfig{
			BaseTemperature: 15.0,
			MaxTemperature:  350.0,
			DamageBaseDPS:   0.5,
			DamageDivisor:   10.0,
			DamageExponent:  1.5,
		},
	}
}

func TestGameConfig_Validate_MissingBossType(t *testing.T) {
	cfg := validGameConfig()
	cfg.Level.BossRoom.BossType = ""

	err := cfg.Validate()
	if err == nil {
		t.Error("Config with empty boss type should return error")
	}
	if strings.Contains(err.Error(), "boss type") == false {
		t.Errorf("Error should mention boss type, got: %v", err)
	}
}
