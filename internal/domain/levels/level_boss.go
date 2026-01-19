package levels

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
)

// GetBossTestLevelConfig returns a configuration for testing the boss system
// Uses level number -2. Minimal world with quick access to boss room.
func GetBossTestLevelConfig() *config.GameConfig {
	cfg := GetTestLevelConfig() // Start with test level config for advanced player stats

	// Create minimal world: ~10 tiles to dig + boss room + floor
	// Each tile is 64 pixels, so:
	// - Surface to boss room: ~10 tiles = 640 pixels
	// - Boss room: ~720 pixels
	// - Floor: 2 tiles = 128 pixels
	// Total: ~1488 pixels
	cfg.World = config.WorldConfig{
		Width:       1280,
		Height:      1488,
		GroundLevel: 128,
		Seed:        42,
		PlayerSpawn: config.PlayerSpawn{X: 100, Y: 58},
		BuildingLayout: config.BuildingLayout{
			MarketX:      100,
			FuelStationX: 300,
			HospitalX:    500,
			UpgradeShopX: 700,
			ItemShopX:    900,
		},
	}

	// Simplified generation for quick testing
	cfg.Generation = config.GenerationConfig{
		Empty: config.TileDistribution{
			PeakDepth: 5,
			Sigma:     100,
			MaxWeight: 0.3,
		},
		Dirt: config.TileDistribution{
			PeakDepth: 5,
			Sigma:     100,
			MaxWeight: 0.7,
		},
		Ores: []config.OreConfig{
			{
				ID: "copper",
				Distribution: config.TileDistribution{
					PeakDepth: 5,
					Sigma:     100,
					MaxWeight: 0.1,
				},
				Value: 10,
			},
		},
		Hazards: []config.HazardConfig{
			{
				ID:        "rock",
				Drillable: false,
				Distribution: config.TileDistribution{
					PeakDepth: 10,
					Sigma:     100,
					MaxWeight: 0.2,
				},
				OnDrillDamage: 0,
				FixedDuration: 0,
			},
		},
	}

	// Update level metadata with boss room
	cfg.Level = config.LevelConfig{
		Number: -2,
		Name:   "Boss Test Level",
		BossRoom: &config.BossRoomConfig{
			BossType:    "test_boss",
			FloorType:   config.FloorConcrete,
			RoomHeight:  720.0,
			FloorHeight: 2.0,
		},
	}

	return cfg
}
