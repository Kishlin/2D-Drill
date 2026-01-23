package entities

import "github.com/Kishlin/drill-game/internal/domain/config"

// Item catalog entry type

type ItemCatalogEntry struct {
	ItemType ItemType
	Price    int
}

// ItemCatalog contains all purchasable items
type ItemCatalog struct {
	Items [5]ItemCatalogEntry
}

func NewItemCatalogFromConfig(itemCfg config.ItemConfig) *ItemCatalog {
	return &ItemCatalog{
		Items: [5]ItemCatalogEntry{
			{ItemType: ItemTeleport, Price: itemCfg.Teleport.Price},
			{ItemType: ItemRepair, Price: itemCfg.Repair.Price},
			{ItemType: ItemRefuel, Price: itemCfg.Refuel.Price},
			{ItemType: ItemBomb, Price: itemCfg.Bomb.Price},
			{ItemType: ItemBigBomb, Price: itemCfg.BigBomb.Price},
		},
	}
}

func (c *ItemCatalog) GetItem(index int) *ItemCatalogEntry {
	if index < 0 || index >= len(c.Items) {
		return nil
	}
	return &c.Items[index]
}
