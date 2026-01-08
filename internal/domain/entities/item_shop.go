package entities

import "github.com/Kishlin/drill-game/internal/domain/types"

const (
	ItemShopWidth  = 160.0
	ItemShopHeight = 192.0
)

type ItemShop struct {
	AABB     types.AABB
	ItemType ItemType
	Price    int
	Name     string
}

func NewItemShop(x, y float32, itemType ItemType, price int, name string) *ItemShop {
	return &ItemShop{
		AABB:     types.NewAABB(x, y, ItemShopWidth, ItemShopHeight),
		ItemType: itemType,
		Price:    price,
		Name:     name,
	}
}

func (s *ItemShop) IsPlayerInRange(player *Player) bool {
	return s.AABB.Intersects(player.AABB)
}
