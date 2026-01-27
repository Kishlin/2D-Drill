package ui

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type FuelStationUI struct {
	state *FuelStationState
}

func NewFuelStationUI() *FuelStationUI {
	return &FuelStationUI{
		state: NewFuelStationState(),
	}
}

func (u *FuelStationUI) Process(player *entities.Player, inputState input.InputState) Result {
	// Handle close shop
	if inputState.CloseShop {
		return Close()
	}

	// Skip first frame to prevent input from previous context
	if u.state.IsFirstFrame() {
		u.state.ClearFirstFrame()
		return NoChange()
	}

	// Handle navigation
	if inputState.NavUp {
		u.state.NavigateUp()
		return NoChange()
	}
	if inputState.NavDown {
		u.state.NavigateDown()
		return NoChange()
	}

	// Handle interaction (refuel)
	if inputState.Interact {
		fuelAmount := u.GetFuelAmount(u.state.SelectedIndex, player)
		if fuelAmount <= 0 {
			return Close()
		}

		cost := u.GetCost(fuelAmount)
		if player.CanAfford(cost) == false {
			return NoChange()
		}

		newFuel := player.Fuel + fuelAmount
		fuelEffects := []effects.Effect{
			effects.TakeMoney{Amount: cost},
			effects.SetFuel{Amount: newFuel},
		}

		// Options 0 and 1 (fixed amounts) stay open for repeated purchases
		// Options 2 and 3 (full/max) close after purchase
		if u.state.SelectedIndex <= 1 {
			return WithEffects(fuelEffects...)
		}
		return CloseWithEffects(fuelEffects...)
	}

	return NoChange()
}

func (u *FuelStationUI) GetRenderState() interface{} {
	return u.state
}

func (u *FuelStationUI) ResetState() {
	u.state.Reset()
}

func (u *FuelStationUI) GetFuelAmount(index int, player *entities.Player) float32 {
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
