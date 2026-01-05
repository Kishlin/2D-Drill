package systems

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type FuelStationSystem struct {
	fuelStation *entities.FuelStation
}

func NewFuelStationSystem(fuelStation *entities.FuelStation) *FuelStationSystem {
	return &FuelStationSystem{fuelStation: fuelStation}
}

func (fss *FuelStationSystem) ProcessRefueling(
	player *entities.Player,
	inputState input.InputState,
) {
	if !inputState.Sell {
		return
	}

	if !fss.fuelStation.IsPlayerInRange(player) {
		return
	}

	// Calculate liters needed to fill tank (use player's upgraded capacity)
	fuelCapacity := player.GetFuelCapacity()
	litersNeeded := fuelCapacity - player.Fuel

	// Calculate cost (1 money per liter, rounded up)
	cost := int(math.Ceil(float64(litersNeeded)))

	// Check if player has enough money
	if player.Money < cost {
		return // Cannot afford refueling
	}

	// Deduct money and refuel
	player.Money -= cost
	player.Fuel = fuelCapacity
}

func (fss *FuelStationSystem) GetFuelStation() *entities.FuelStation {
	return fss.fuelStation
}
