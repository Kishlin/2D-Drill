package world

import "github.com/Kishlin/drill-game/internal/domain/entities"

const TileSize = 64 // pixels

type World struct {
	GroundLevel float32
	Width       float32
	Height      float32
	tiles       map[[2]int]*entities.Tile // Sparse map: [x, y] -> Tile

	generator    *ChunkGenerator
	loadedChunks map[[2]int]bool
	seed         int64
}

func NewWorld(width, height, groundLevel float32, seed int64) *World {
	return &World{
		Width:        width,
		Height:       height,
		GroundLevel:  groundLevel,
		tiles:        make(map[[2]int]*entities.Tile),
		generator:    NewChunkGenerator(seed, groundLevel),
		loadedChunks: make(map[[2]int]bool),
		seed:         seed,
	}
}

// EnsureChunkLoaded generates a chunk if not already loaded
func (w *World) EnsureChunkLoaded(chunkX, chunkY int) {
	key := [2]int{chunkX, chunkY}
	if w.loadedChunks[key] {
		return
	}

	for localX := 0; localX < ChunkSize; localX++ {
		for localY := 0; localY < ChunkSize; localY++ {
			tileX := chunkX*ChunkSize + localX
			tileY := chunkY*ChunkSize + localY

			if !w.isGridInBounds(tileX, tileY) {
				continue
			}

			tile := w.generator.GenerateTile(tileX, tileY)

			// Only store solid tiles (Dirt, Ore) - sparse storage
			if tile.Type != entities.TileTypeEmpty {
				w.tiles[[2]int{tileX, tileY}] = tile
			}
		}
	}

	w.loadedChunks[key] = true
}

// UpdateChunksAroundPlayer proactively loads a 3×3 grid of chunks around player
func (w *World) UpdateChunksAroundPlayer(playerX, playerY float32) {
	playerChunkX := int(playerX/TileSize) / ChunkSize
	playerChunkY := int(playerY/TileSize) / ChunkSize

	// Load 3×3 grid around player
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			w.EnsureChunkLoaded(playerChunkX+dx, playerChunkY+dy)
		}
	}
}

// isGridInBounds checks if tile coordinates are within world bounds
func (w *World) isGridInBounds(gridX, gridY int) bool {
	pixelX := float32(gridX) * TileSize
	pixelY := float32(gridY) * TileSize
	return w.IsInBounds(pixelX, pixelY)
}

func (w *World) GetGroundLevel() float32 {
	return w.GroundLevel
}

// GetTileAt returns tile at pixel coordinates (returns nil if empty/air)
func (w *World) GetTileAt(pixelX, pixelY float32) *entities.Tile {
	tileX := int(pixelX / TileSize)
	tileY := int(pixelY / TileSize)
	return w.GetTileAtGrid(tileX, tileY)
}

// GetTileAtGrid returns tile at grid coordinates (triggers chunk load if needed)
func (w *World) GetTileAtGrid(gridX, gridY int) *entities.Tile {
	chunkX := gridX / ChunkSize
	chunkY := gridY / ChunkSize
	w.EnsureChunkLoaded(chunkX, chunkY)

	return w.tiles[[2]int{gridX, gridY}]
}

// DigTile removes tile at pixel coordinates
// Returns the removed tile (if any) and success status
func (w *World) DigTile(pixelX, pixelY float32) (*entities.Tile, bool) {
	tileX := int(pixelX / TileSize)
	tileY := int(pixelY / TileSize)

	tile := w.tiles[[2]int{tileX, tileY}]
	if tile != nil && tile.IsDiggable() {
		delete(w.tiles, [2]int{tileX, tileY})
		return tile, true
	}
	return nil, false
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

// SetTile sets a tile at the given grid coordinates (for testing)
func (w *World) SetTile(gridX, gridY int, tile *entities.Tile) {
	if tile == nil || tile.Type == entities.TileTypeEmpty {
		delete(w.tiles, [2]int{gridX, gridY})
	} else {
		w.tiles[[2]int{gridX, gridY}] = tile
	}
}

func (w *World) IsInBounds(x, y float32) bool {
	return x >= 0 && x <= w.Width && y >= 0 && y <= w.Height
}
