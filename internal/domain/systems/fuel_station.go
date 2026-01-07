package systems

import (
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

	player.Refuel()
}

func (fss *FuelStationSystem) GetFuelStation() *entities.FuelStation {
	return fss.fuelStation
}
