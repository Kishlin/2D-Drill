package effects

import "github.com/Kishlin/drill-game/internal/domain/entities"

type ClearOreInventory struct{}

func (e ClearOreInventory) Apply(player *entities.Player) {
	player.OreInventory = make(map[string]int)
}

type AddItem struct {
	ItemType entities.ItemType
}

func (e AddItem) Apply(player *entities.Player) {
	player.AddItem(e.ItemType)
}
