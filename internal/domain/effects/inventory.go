package effects

import "github.com/Kishlin/drill-game/internal/domain/entities"

type ClearOreInventory struct{}

func (e ClearOreInventory) Apply(ctx *EffectContext) {
	ctx.Player.OreInventory = make(map[string]int)
}

type AddItem struct {
	ItemType entities.ItemType
}

func (e AddItem) Apply(ctx *EffectContext) {
	ctx.Player.AddItem(e.ItemType)
}
