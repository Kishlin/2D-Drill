package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type HospitalSystem struct {
	hospital *entities.Hospital
}

func NewHospitalSystem(hospital *entities.Hospital) *HospitalSystem {
	return &HospitalSystem{hospital: hospital}
}

func (hs *HospitalSystem) ProcessHealing(
	player *entities.Player,
	inputState input.InputState,
) {
	if !inputState.Sell {
		return
	}

	if !hs.hospital.IsPlayerInRange(player) {
		return
	}

	player.Heal()
}

func (hs *HospitalSystem) GetHospital() *entities.Hospital {
	return hs.hospital
}
