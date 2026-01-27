package ui

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type HospitalUI struct {
	state *HospitalState
}

func NewHospitalUI() *HospitalUI {
	return &HospitalUI{
		state: NewHospitalState(),
	}
}

func (u *HospitalUI) Process(player *entities.Player, inputState input.InputState) Result {
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

	// Handle interaction (heal)
	if inputState.Interact {
		healAmount := u.GetHealAmount(u.state.SelectedIndex, player)
		if healAmount <= 0 {
			return Close()
		}

		cost := u.GetCost(healAmount)
		if player.CanAfford(cost) == false {
			return NoChange()
		}

		newHP := player.HP + healAmount
		healEffects := []effects.Effect{
			effects.TakeMoney{Amount: cost},
			effects.SetHP{Amount: newHP},
		}

		// Options 0 and 1 (fixed amounts) stay open for repeated purchases
		// Options 2 and 3 (full/max) close after purchase
		if u.state.SelectedIndex <= 1 {
			return WithEffects(healEffects...)
		}
		return CloseWithEffects(healEffects...)
	}

	return NoChange()
}

func (u *HospitalUI) GetRenderState() interface{} {
	return u.state
}

func (u *HospitalUI) ResetState() {
	u.state.Reset()
}

func (u *HospitalUI) GetHealAmount(index int, player *entities.Player) float32 {
	hpNeeded := player.MaxHP() - player.HP
	if hpNeeded <= 0 {
		return 0
	}

	switch index {
	case 0: // Restore 1 HP
		return min(1, hpNeeded)
	case 1: // Restore 10 HP
		return min(10, hpNeeded)
	case 2: // Restore All HP
		return hpNeeded
	case 3: // Max Affordable
		maxAffordable := float32(player.Money) / 2.0
		return min(hpNeeded, float32(math.Floor(float64(maxAffordable))))
	}
	return 0
}

func (u *HospitalUI) GetCost(healAmount float32) int {
	return int(math.Ceil(float64(healAmount) * 2.0))
}

func (u *HospitalUI) GetOptionLabel(index int) string {
	switch index {
	case 0:
		return "Restore 1 HP"
	case 1:
		return "Restore 10 HP"
	case 2:
		return "Restore All HP"
	case 3:
		return "Max Affordable"
	}
	return ""
}
