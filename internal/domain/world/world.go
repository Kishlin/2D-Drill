package world

import "github.com/Kishlin/drill-game/internal/domain/entities"

const TileSize = 64 // pixels

type World struct {
	GroundLevel float32
	Width       float32
	Height      float32
	tiles       map[[2]int]*entities.Tile // Sparse map: [x, y] -> Tile
}

func NewWorld(width, height, groundLevel float32) *World {
	w := &World{
		Width:       width,
		Height:      height,
		GroundLevel: groundLevel,
		tiles:       make(map[[2]int]*entities.Tile),
	}
	w.generateInitialTiles()
	return w
}

// generateInitialTiles initializes world with solid dirt below ground level
func (w *World) generateInitialTiles() {
	tilesWide := int(w.Width / TileSize)
	tilesHigh := int(w.Height / TileSize)
	groundTileY := int(w.GroundLevel / TileSize)

	for x := 0; x < tilesWide; x++ {
		for y := groundTileY; y < tilesHigh; y++ {
			w.tiles[[2]int{x, y}] = entities.NewTile(entities.TileTypeDirt)
		}
	}
}

func (w *World) GetGroundLevel() float32 {
	return w.GroundLevel
}

// GetTileAt returns tile at pixel coordinates (returns nil if empty/air)
func (w *World) GetTileAt(pixelX, pixelY float32) *entities.Tile {
	tileX := int(pixelX / TileSize)
	tileY := int(pixelY / TileSize)
	return w.tiles[[2]int{tileX, tileY}]
}

// GetTileAtGrid returns tile at grid coordinates
func (w *World) GetTileAtGrid(gridX, gridY int) *entities.Tile {
	return w.tiles[[2]int{gridX, gridY}]
}

// DigTile removes tile at pixel coordinates, returns true if successful
func (w *World) DigTile(pixelX, pixelY float32) bool {
	tileX := int(pixelX / TileSize)
	tileY := int(pixelY / TileSize)

	tile := w.tiles[[2]int{tileX, tileY}]
	if tile != nil && tile.IsDiggable() {
		delete(w.tiles, [2]int{tileX, tileY})
		return true
	}
	return false
}

// IsTileSolid checks if there's a solid tile at pixel coordinates
func (w *World) IsTileSolid(pixelX, pixelY float32) bool {
	tile := w.GetTileAt(pixelX, pixelY)
	return tile != nil && tile.IsSolid()
}

// GetAllTiles returns all tiles (for rendering)
func (w *World) GetAllTiles() map[[2]int]*entities.Tile {
	return w.tiles
}

func (w *World) IsInBounds(x, y float32) bool {
	return x >= 0 && x <= w.Width && y >= 0 && y <= w.Height
}
