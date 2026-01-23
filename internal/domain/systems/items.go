package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

// DetectItemUsage checks for item inputs and returns effects to apply
func DetectItemUsage(player *entities.Player, inputState input.InputState, itemCfg config.ItemConfig) []effects.Effect {
	var result []effects.Effect

	if inputState.UseTeleport && player.UseItem(entities.ItemTeleport) {
		result = append(result, effects.Teleport{})
	}
	if inputState.UseRepair && player.UseItem(entities.ItemRepair) {
		result = append(result, effects.Repair{})
	}
	if inputState.UseRefuel && player.UseItem(entities.ItemRefuel) {
		result = append(result, effects.Refuel{})
	}
	if inputState.UseBomb && player.UseItem(entities.ItemBomb) {
		result = append(result, effects.Bomb{Radius: itemCfg.Bomb.Radius, Damage: 10.0})
	}
	if inputState.UseBigBomb && player.UseItem(entities.ItemBigBomb) {
		result = append(result, effects.Bomb{Radius: itemCfg.BigBomb.Radius, Damage: 25.0})
	}

	return result
}
