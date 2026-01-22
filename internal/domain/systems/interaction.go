package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/components"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

// DetectInteraction checks if the player is pressing interact and overlapping with a building
// Returns the interactable type if an interaction should occur, nil otherwise
func DetectInteraction(player *entities.Player, buildings []*entities.Building, inputState input.InputState) *components.InteractableType {
	if !inputState.Interact {
		return nil
	}

	playerPos := components.Position{AABB: player.AABB}

	for _, b := range buildings {
		if b.Position.Intersects(playerPos) {
			interactableType := b.Interactable.Type
			return &interactableType
		}
	}
	return nil
}
