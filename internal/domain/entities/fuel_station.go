package entities

import "github.com/Kishlin/drill-game/internal/domain/types"

const (
	FuelStationWidth  = 320.0 // 5 tiles * 64px
	FuelStationHeight = 192.0 // 3 tiles * 64px
)

type FuelStation struct {
	AABB types.AABB
}

func NewFuelStation(x, y float32) *FuelStation {
	return &FuelStation{
		AABB: types.NewAABB(x, y, FuelStationWidth, FuelStationHeight),
	}
}

// IsPlayerInRange checks if player AABB overlaps with fuel station AABB
func (fs *FuelStation) IsPlayerInRange(player *Player) bool {
	return fs.AABB.Intersects(player.AABB)
}
