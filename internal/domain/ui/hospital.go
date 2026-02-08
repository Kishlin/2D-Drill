package ui

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type HospitalUI struct {
	state *ModalServiceState
}

func NewHospitalUI() *HospitalUI {
	return &HospitalUI{
		state: NewModalServiceState(),
	}
}

func (u *HospitalUI) Process(player *entities.Player, inputState input.InputState) Result {
	return processModalService(u.state, u, player, inputState)
}

func (u *HospitalUI) GetRenderState() interface{} {
	return u.state
}

func (u *HospitalUI) ResetState() {
	u.state.Reset()
}

func (u *HospitalUI) GetAmount(index int, player *entities.Player) float32 {
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

func (u *HospitalUI) GetCost(amount float32) int {
	return int(math.Ceil(float64(amount) * 2.0))
}

func (u *HospitalUI) BuildEffect(amount float32, player *entities.Player) effects.Effect {
	return effects.SetHP{Amount: player.HP + amount}
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
