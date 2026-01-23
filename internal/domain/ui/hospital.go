package ui

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type HospitalUI struct{}

func NewHospitalUI() *HospitalUI {
	return &HospitalUI{}
}

func (u *HospitalUI) Process(player *entities.Player, inputState input.InputState) Result {
	hpNeeded := player.MaxHP() - player.HP
	if hpNeeded <= 0 {
		return Close()
	}

	cost := int(math.Ceil(float64(hpNeeded) * 2.0))
	if !player.CanAfford(cost) {
		return Close()
	}

	return CloseWithEffects(
		effects.TakeMoney{Amount: cost},
		effects.SetHP{Amount: player.MaxHP()},
	)
}

func (u *HospitalUI) GetRenderState() interface{} {
	return nil
}
