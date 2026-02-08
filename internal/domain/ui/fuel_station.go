package ui

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type FuelStationUI struct {
	state *ModalServiceState
}

func NewFuelStationUI() *FuelStationUI {
	return &FuelStationUI{
		state: NewModalServiceState(),
	}
}

func (u *FuelStationUI) Process(player *entities.Player, inputState input.InputState) Result {
	return processModalService(u.state, u, player, inputState)
}

func (u *FuelStationUI) GetRenderState() interface{} {
	return u.state
}

func (u *FuelStationUI) ResetState() {
	u.state.Reset()
}

func (u *FuelStationUI) GetAmount(index int, player *entities.Player) float32 {
	fuelNeeded := player.FuelCapacity() - player.Fuel
	if fuelNeeded <= 0 {
		return 0
	}

	switch index {
	case 0: // Refuel 1L
		return min(1, fuelNeeded)
	case 1: // Refuel 10L
		return min(10, fuelNeeded)
	case 2: // Full Tank
		return fuelNeeded
	case 3: // Max Affordable
		maxAffordable := float32(player.Money)
		return min(fuelNeeded, float32(math.Floor(float64(maxAffordable))))
	}
	return 0
}

func (u *FuelStationUI) GetCost(fuelAmount float32) int {
	return int(math.Ceil(float64(fuelAmount)))
}

func (u *FuelStationUI) BuildEffect(amount float32, player *entities.Player) effects.Effect {
	return effects.SetFuel{Amount: player.Fuel + amount}
}

func (u *FuelStationUI) GetOptionLabel(index int) string {
	switch index {
	case 0:
		return "Refuel 1L"
	case 1:
		return "Refuel 10L"
	case 2:
		return "Full Tank"
	case 3:
		return "Max Affordable"
	}
	return ""
}
