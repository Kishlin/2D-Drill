package physics

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// CollisionResult contains the results of collision detection and resolution
type CollisionResult struct {
	Position types.Vec2
	Velocity types.Vec2
	OnGround bool
}

// ResolveGroundCollision checks for and resolves collision with tiles
// Checks the full width of the player to support standing on tiles even with partial overlap
// Returns updated position, velocity, and ground contact state
func ResolveGroundCollision(position, velocity types.Vec2, height float32, w *world.World) CollisionResult {
	playerBottom := position.Y + height

	// Check for tile collision across the full width of the player (left, center, right)
	// This allows the player to stand on a tile even if only a small portion is on top
	checkPoints := []float32{
		position.X,                          // left edge
		position.X + entities.PlayerWidth/2, // center
		position.X + entities.PlayerWidth,   // right edge
	}

	for _, checkX := range checkPoints {
		tile := w.GetTileAt(checkX, playerBottom)
		if tile != nil && tile.IsSolid() {
			// Found a solid tile, snap to its top
			tileGridY := int(playerBottom / world.TileSize)
			tileTopY := float32(tileGridY * world.TileSize)

			newPos := types.Vec2{
				X: position.X,
				Y: tileTopY - height,
			}

			newVel := velocity
			if newVel.Y >= 0 {
				newVel.Y = 0
			}

			return CollisionResult{
				Position: newPos,
				Velocity: newVel,
				OnGround: true,
			}
		}
	}

	// No tile collision found, so no collision at all
	// (Removed legacy ground level - now purely tile-based)
	return CollisionResult{
		Position: position,
		Velocity: velocity,
		OnGround: false,
	}
}
