package entities

import "github.com/Kishlin/drill-game/internal/domain/types"

const (
	ShopWidth  = 320.0 // 5 tiles * 64px
	ShopHeight = 192.0 // 3 tiles * 64px
)

type Shop struct {
	AABB types.AABB
}

func NewShop(x, y float32) *Shop {
	return &Shop{
		AABB: types.NewAABB(x, y, ShopWidth, ShopHeight),
	}
}

// IsPlayerInRange checks if player AABB overlaps with shop AABB
func (s *Shop) IsPlayerInRange(player *Player) bool {
	return s.AABB.Intersects(player.AABB)
}
