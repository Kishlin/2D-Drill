package examples

import (
	"github.com/Kishlin/drill-game/internal/domain/worlds"
)

// init registers the hard mode world on package initialization
func init() {
	worlds.RegisterWorld("hard_mode", NewHardModeWorld())
}

// NewHardModeWorld creates a harder difficulty configuration
// Reduces ore values by 50%, doubles upgrade prices, increases hazard spawn rates
func NewHardModeWorld() *worlds.WorldGameConfig {
	config := worlds.NewDefaultWorld()

	// Reduce ore values by 50%
	for i := range config.Ores.Ores {
		config.Ores.Ores[i].Value /= 2
	}

	// Double all upgrade prices
	for i := range config.Upgrades.Engine {
		config.Upgrades.Engine[i].Price *= 2
	}
	for i := range config.Upgrades.Hull {
		config.Upgrades.Hull[i].Price *= 2
	}
	for i := range config.Upgrades.FuelTank {
		config.Upgrades.FuelTank[i].Price *= 2
	}
	for i := range config.Upgrades.CargoHold {
		config.Upgrades.CargoHold[i].Price *= 2
	}
	for i := range config.Upgrades.HeatShield {
		config.Upgrades.HeatShield[i].Price *= 2
	}
	for i := range config.Upgrades.Drill {
		config.Upgrades.Drill[i].Price *= 2
	}

	// Double item prices
	for i := range config.Items.Items {
		config.Items.Items[i].Price *= 2
	}

	// Increase hazard spawn rates
	for i := range config.Hazards.Hazards {
		config.Hazards.Hazards[i].Distribution.MaxWeight *= 1.5
	}

	// Player starts with no money, basic upgrades
	config.PlayerInitial = worlds.PlayerInitialStatus{
		Money:          0,
		EngineTier:     0,
		HullTier:       0,
		FuelTankTier:   0,
		CargoHoldTier:  0,
		HeatShieldTier: 0,
		DrillTier:      0,
		ItemInventory:  [5]int{0, 0, 0, 0, 0},
	}

	return config
}
