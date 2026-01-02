package physics

import (
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// TileCollision represents a collision with a tile
type TileCollision struct {
	GridX, GridY int
	TileAABB     types.AABB
}

// GetOccupiedTileRange calculates which tiles an AABB overlaps
// Returns (minX, maxX, minY, maxY) in grid coordinates
func GetOccupiedTileRange(aabb types.AABB, tileSize float32) (minX, maxX, minY, maxY int) {
	minX = int(aabb.X / tileSize)
	maxX = int((aabb.X + aabb.Width - 0.001) / tileSize)
	minY = int(aabb.Y / tileSize)
	maxY = int((aabb.Y + aabb.Height - 0.001) / tileSize)
	return
}

// CheckCollisions finds all solid tiles intersecting the AABB
func CheckCollisions(aabb types.AABB, w *world.World) []TileCollision {
	minX, maxX, minY, maxY := GetOccupiedTileRange(aabb, world.TileSize)

	var collisions []TileCollision

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			tile := w.GetTileAtGrid(x, y)

			if tile == nil || !tile.IsSolid() {
				continue
			}

			tileAABB := tile.GetAABB(x, y, world.TileSize)

			if aabb.Intersects(tileAABB) {
				collisions = append(collisions, TileCollision{
					GridX:    x,
					GridY:    y,
					TileAABB: tileAABB,
				})
			}
		}
	}

	return collisions
}

// ResolveCollisionsX resolves X-axis collisions only
// Takes AABB and velocity by value, returns updated AABB and velocity
func ResolveCollisionsX(aabb types.AABB, velocity types.Vec2, collisions []TileCollision) (types.AABB, types.Vec2) {
	if len(collisions) == 0 {
		return aabb, velocity
	}

	newAABB := aabb
	newVel := velocity

	for _, col := range collisions {
		dx, _ := newAABB.Penetration(col.TileAABB)

		if dx != 0 {
			newAABB.X -= dx // Adjust position
			newVel.X = 0    // Stop horizontal movement
		}
	}

	return newAABB, newVel
}

// ResolveCollisionsY resolves Y-axis collisions and detects ground
// Takes AABB and velocity by value, returns updated AABB, velocity, and ground state
func ResolveCollisionsY(aabb types.AABB, velocity types.Vec2, collisions []TileCollision) (types.AABB, types.Vec2, bool) {
	if len(collisions) == 0 {
		return aabb, velocity, false
	}

	newAABB := aabb
	newVel := velocity
	onGround := false

	for _, col := range collisions {
		_, dy := newAABB.Penetration(col.TileAABB)

		if dy != 0 {
			newAABB.Y -= dy // Adjust position

			if dy > 0 {
				// Positive dy with subtraction = push up = ground collision
				onGround = true
				newVel.Y = 0
			} else {
				// Negative dy with subtraction = push down = ceiling collision
				newVel.Y = 0
			}
		}
	}

	return newAABB, newVel, onGround
}
