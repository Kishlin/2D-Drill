package entities

import "github.com/Kishlin/drill-game/internal/domain/types"

const (
	HospitalWidth  = 320.0 // 5 tiles * 64px
	HospitalHeight = 192.0 // 3 tiles * 64px
)

type Hospital struct {
	AABB types.AABB
}

func NewHospital(x, y float32) *Hospital {
	return &Hospital{
		AABB: types.NewAABB(x, y, HospitalWidth, HospitalHeight),
	}
}

func (h *Hospital) IsPlayerInRange(player *Player) bool {
	return h.AABB.Intersects(player.AABB)
}
