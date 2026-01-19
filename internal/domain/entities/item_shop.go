package entities

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

const (
	ItemShopWidth  = 320.0
	ItemShopHeight = 192.0
)

type ItemCatalogEntry struct {
	ItemType ItemType
	Price    int
}

type ItemShop struct {
	AABB    types.AABB
	Catalog [5]ItemCatalogEntry // Fixed array: Teleport, Repair, Refuel, Bomb, BigBomb
}

func NewItemShopFromConfig(x, y float32, itemCfg config.ItemConfig) *ItemShop {
	return &ItemShop{
		AABB: types.NewAABB(x, y, ItemShopWidth, ItemShopHeight),
		Catalog: [5]ItemCatalogEntry{
			{ItemType: ItemTeleport, Price: itemCfg.Teleport.Price},
			{ItemType: ItemRepair, Price: itemCfg.Repair.Price},
			{ItemType: ItemRefuel, Price: itemCfg.Refuel.Price},
			{ItemType: ItemBomb, Price: itemCfg.Bomb.Price},
			{ItemType: ItemBigBomb, Price: itemCfg.BigBomb.Price},
		},
	}
}

func (s *ItemShop) IsPlayerInRange(player *Player) bool {
	return s.AABB.Intersects(player.AABB)
}

func (s *ItemShop) GetItem(index int) *ItemCatalogEntry {
	if index < 0 || index >= len(s.Catalog) {
		return nil
	}
	return &s.Catalog[index]
}
