package levels

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
)

// GetTestLevelConfig returns a test configuration with an advanced player
// for easier testing. Use levelNum -1 to load this config.
func GetTestLevelConfig() *config.GameConfig {
	// Start with Level 1 as base
	cfg := GetLevel1Config()

	// Override player config for testing
	cfg.Player = config.PlayerConfig{
		StartingMoney: 100000,
		StartingItems: [5]int{
			3,    // Teleport
			5,    // Repair
			5,    // Refuel
			10,   // Bomb
			2000, // Big Bomb
		},
		StartingUpgrades: config.StartingUpgrades{
			Engine:     5, // Max tier
			Hull:       3, // Mid-tier
			FuelTank:   3, // Mid-tier
			CargoHold:  3, // Mid-tier
			HeatShield: 3, // Mid-tier
			Drill:      5, // Max tier
		},
	}

	// Update level metadata
	cfg.Level = config.LevelConfig{
		Number: -1,
		Name:   "Test Level",
	}

	return cfg
}
