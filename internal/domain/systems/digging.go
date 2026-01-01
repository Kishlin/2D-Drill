package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type DiggingSystem struct {
	world *world.World
}

func NewDiggingSystem(w *world.World) *DiggingSystem {
	return &DiggingSystem{world: w}
}

// ProcessDigging handles digging logic and player alignment
func (ds *DiggingSystem) ProcessDigging(
	player *entities.Player,
	inputState input.InputState,
) {
	if !inputState.Dig {
		return
	}

	// Calculate tile beneath player's center-bottom
	playerCenterX := player.Position.X + entities.PlayerWidth/2
	playerBottomY := player.Position.Y + entities.PlayerHeight

	// Check tile directly below player
	tile := ds.world.GetTileAt(playerCenterX, playerBottomY)
	if tile != nil && tile.IsDiggable() {
		// Snap player to tile grid (Motherload behavior)
		tileGridX := int(playerCenterX / world.TileSize)
		tileCenterX := float32(tileGridX)*world.TileSize + world.TileSize/2

		// Align player center with tile center
		player.Position.X = tileCenterX - entities.PlayerWidth/2

		// Dig the tile
		ds.world.DigTile(playerCenterX, playerBottomY)
	}
}
