package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

// ConsumeFuel drains fuel based on player input state.
// Movement inputs (Left, Right, Up, Drill) consume fuel faster than idle state.
func ConsumeFuel(player *entities.Player, inputState input.InputState, dt float32, fuelCfg config.FuelSystemConfig) {
	var rate float32
	if inputState.HasMovementInput() {
		rate = fuelCfg.ConsumptionMoving
	} else {
		rate = fuelCfg.ConsumptionIdle
	}

	fuelConsumed := rate * dt

	player.Fuel -= fuelConsumed
	if player.Fuel < 0 {
		player.Fuel = 0
	}
}
