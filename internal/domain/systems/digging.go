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

// ProcessDigging handles downward digging logic and player alignment
func (ds *DiggingSystem) ProcessDigging(
	player *entities.Player,
	inputState input.InputState,
) {
	if !inputState.Dig {
		return
	}

	// Calculate tile beneath player's center-bottom
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height

	// Check tile directly below player
	tile := ds.world.GetTileAt(playerCenterX, playerBottomY)
	if tile != nil && tile.IsDiggable() {
		// Snap player to tile grid (Motherload behavior)
		tileGridX := int(playerCenterX / world.TileSize)
		tileCenterX := float32(tileGridX)*world.TileSize + world.TileSize/2

		// Align player center with tile center
		player.AABB.X = tileCenterX - player.AABB.Width/2

		// Dig the tile
		ds.world.DigTile(playerCenterX, playerBottomY)
	}
}

// ProcessHorizontalDigging handles auto-digging when moving left/right against walls
// Only works when player is grounded
func (ds *DiggingSystem) ProcessHorizontalDigging(
	player *entities.Player,
	inputState input.InputState,
) {
	// Only dig horizontally when grounded
	if !player.OnGround {
		return
	}

	playerCenterY := player.AABB.Y + player.AABB.Height/2

	if inputState.Left {
		// Check tile just left of player's left edge
		tileX := player.AABB.X - 1
		tile := ds.world.GetTileAt(tileX, playerCenterY)
		if tile != nil && tile.IsDiggable() {
			ds.world.DigTile(tileX, playerCenterY)
			return
		}
	}

	if inputState.Right {
		// Check tile just right of player's right edge
		tileX := player.AABB.X + player.AABB.Width + 1
		tile := ds.world.GetTileAt(tileX, playerCenterY)
		if tile != nil && tile.IsDiggable() {
			ds.world.DigTile(tileX, playerCenterY)
			return
		}
	}
}
