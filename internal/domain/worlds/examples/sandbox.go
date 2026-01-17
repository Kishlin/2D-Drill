package examples

import (
	"github.com/Kishlin/drill-game/internal/domain/worlds"
)

// init registers the sandbox world on package initialization
func init() {
	worlds.RegisterWorld("sandbox", NewSandboxWorld())
}

// NewSandboxWorld creates an easy mode sandbox for testing and creative play
// Features: 10x ore values, free upgrades, no hazards, fully equipped player
func NewSandboxWorld() *worlds.WorldGameConfig {
	config := worlds.NewDefaultWorld()

	// 10x ore values for faster progression
	for i := range config.Ores.Ores {
		config.Ores.Ores[i].Value *= 10
	}

	// Free upgrades - all tiers cost $0
	for i := range config.Upgrades.Engine {
		config.Upgrades.Engine[i].Price = 0
	}
	for i := range config.Upgrades.Hull {
		config.Upgrades.Hull[i].Price = 0
	}
	for i := range config.Upgrades.FuelTank {
		config.Upgrades.FuelTank[i].Price = 0
	}
	for i := range config.Upgrades.CargoHold {
		config.Upgrades.CargoHold[i].Price = 0
	}
	for i := range config.Upgrades.HeatShield {
		config.Upgrades.HeatShield[i].Price = 0
	}
	for i := range config.Upgrades.Drill {
		config.Upgrades.Drill[i].Price = 0
	}

	// Free items
	for i := range config.Items.Items {
		config.Items.Items[i].Price = 0
	}

	// Remove hazards (no dangers, no progression gates)
	config.Hazards.Hazards = []worlds.HazardDefinition{}

	// Start with max upgrades, lots of money and items
	config.PlayerInitial = worlds.PlayerInitialStatus{
		Money:          100000,
		EngineTier:     5, // Mk5 (highest tier)
		HullTier:       5, // Mk5
		FuelTankTier:   5, // Mk5
		CargoHoldTier:  5, // Mk5
		HeatShieldTier: 5, // Mk5
		DrillTier:      5, // Mk5
		ItemInventory:  [5]int{99, 99, 99, 99, 99}, // Max of each item
	}

	return config
}
