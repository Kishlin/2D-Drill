package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

const (
	// Fuel consumption rates in liters per second
	FuelConsumptionMoving float32 = 10.0 / 30.0  // 0.33333 L/s when actively moving/digging
	FuelConsumptionIdle   float32 = 10.0 / 120.0 // 0.08333 L/s when idle (no inputs)
)

type FuelSystem struct {
	// Empty for now - could add config or state later
}

func NewFuelSystem() *FuelSystem {
	return &FuelSystem{}
}

// ConsumeFuel drains fuel based on player input state
// Movement inputs (Left, Right, Up, Dig) consume fuel faster than idle state
func (fs *FuelSystem) ConsumeFuel(
	player *entities.Player,
	inputState input.InputState,
	dt float32,
) {
	// Determine consumption rate based on input
	var rate float32
	if inputState.HasMovementInput() {
		rate = FuelConsumptionMoving
	} else {
		rate = FuelConsumptionIdle
	}

	// Calculate fuel consumed this frame
	fuelConsumed := rate * dt

	// Drain fuel (clamp at zero, never go negative)
	player.Fuel -= fuelConsumed
	if player.Fuel < 0 {
		player.Fuel = 0
	}
}
