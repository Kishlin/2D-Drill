package systems

import (
	"math"

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

	// Use player's upgraded max HP
	maxHP := player.GetMaxHP()
	hpNeeded := maxHP - player.HP

	if hpNeeded <= 0 {
		return
	}

	// Calculate cost: $2 per HP, rounded up
	cost := int(math.Ceil(float64(hpNeeded) * 2.0))

	if player.Money < cost {
		return // Cannot afford healing
	}

	player.Money -= cost
	player.HP = maxHP
}

func (hs *HospitalSystem) GetHospital() *entities.Hospital {
	return hs.hospital
}
