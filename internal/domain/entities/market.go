package entities

import "github.com/Kishlin/drill-game/internal/domain/types"

const (
	MarketWidth  = 320.0 // 5 tiles * 64px
	MarketHeight = 192.0 // 3 tiles * 64px
)

type Market struct {
	AABB types.AABB
}

func NewMarket(x, y float32) *Market {
	return &Market{
		AABB: types.NewAABB(x, y, MarketWidth, MarketHeight),
	}
}

func (m *Market) IsPlayerInRange(player *Player) bool {
	return m.AABB.Intersects(player.AABB)
}
