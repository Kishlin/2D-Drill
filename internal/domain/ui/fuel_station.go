package ui

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type FuelStationUI struct{}

func NewFuelStationUI() *FuelStationUI {
	return &FuelStationUI{}
}

func (u *FuelStationUI) Process(player *entities.Player, inputState input.InputState) Result {
	fuelNeeded := player.FuelCapacity() - player.Fuel
	if fuelNeeded <= 0 {
		return Close()
	}

	cost := int(math.Ceil(float64(fuelNeeded)))
	if !player.CanAfford(cost) {
		return Close()
	}

	return CloseWithEffects(
		effects.TakeMoney{Amount: cost},
		effects.SetFuel{Amount: player.FuelCapacity()},
	)
}

func (u *FuelStationUI) GetRenderState() interface{} {
	return nil
}
