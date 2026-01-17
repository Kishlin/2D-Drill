package examples

import (
	"github.com/Kishlin/drill-game/internal/domain/worlds"
)

// init registers the endgame test world on package initialization
func init() {
	worlds.RegisterWorld("endgame_test", NewEndgameTestWorld())
}

// NewEndgameTestWorld creates a mid-tier starting configuration for testing deep content
// Player starts with Mk3 upgrades and some items, allowing safe exploration to ~30k depth
func NewEndgameTestWorld() *worlds.WorldGameConfig {
	config := worlds.NewDefaultWorld()

	// Keep default balance for testing
	// (ore values, prices, hazard distribution unchanged)

	// Start with mid-tier upgrades (Mk3) to test endgame content
	// Mk3 Heat Shield provides enough resistance to reach ~30k depth safely
	config.PlayerInitial = worlds.PlayerInitialStatus{
		Money:          5000,
		EngineTier:     3, // Mk3
		HullTier:       3, // Mk3 (30 HP max)
		FuelTankTier:   3, // Mk3 (32 L capacity)
		CargoHoldTier:  3, // Mk3 (24 units capacity)
		HeatShieldTier: 3, // Mk3 (190 heat resistance - safe to ~30k depth)
		DrillTier:      3, // Mk3 (4x speed)
		ItemInventory:  [5]int{5, 5, 5, 3, 1}, // Varied item counts for testing
	}

	return config
}
