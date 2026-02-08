package ui

import (
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

// ModalServiceProvider defines the service-specific logic for a modal service UI.
type ModalServiceProvider interface {
	GetAmount(index int, player *entities.Player) float32
	GetCost(amount float32) int
	BuildEffect(amount float32, player *entities.Player) effects.Effect
}

func processModalService(state *ModalServiceState, provider ModalServiceProvider, player *entities.Player, inputState input.InputState) Result {
	if inputState.CloseShop {
		return Close()
	}

	if state.IsFirstFrame() {
		state.ClearFirstFrame()
		return NoChange()
	}

	if inputState.NavUp {
		state.NavigateUp()
		return NoChange()
	}
	if inputState.NavDown {
		state.NavigateDown()
		return NoChange()
	}

	if inputState.Interact {
		amount := provider.GetAmount(state.SelectedIndex, player)
		if amount <= 0 {
			return Close()
		}

		cost := provider.GetCost(amount)
		if player.CanAfford(cost) == false {
			return NoChange()
		}

		serviceEffects := []effects.Effect{
			effects.TakeMoney{Amount: cost},
			provider.BuildEffect(amount, player),
		}

		if state.SelectedIndex <= 1 {
			return WithEffects(serviceEffects...)
		}
		return CloseWithEffects(serviceEffects...)
	}

	return NoChange()
}
