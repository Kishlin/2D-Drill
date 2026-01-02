package world

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
)

func TestEnsureChunkLoaded_OnlyOnce(t *testing.T) {
	world := NewWorld(7680, 64000, 640, 42)

	// First load
	world.EnsureChunkLoaded(0, 0)
	if !world.loadedChunks[[2]int{0, 0}] {
		t.Error("Chunk should be marked as loaded")
	}

	// Second load should not regenerate
	tilesBefore := len(world.tiles)
	world.EnsureChunkLoaded(0, 0)
	tilesAfter := len(world.tiles)

	if tilesBefore != tilesAfter {
		t.Error("Chunk should not regenerate tiles on second load")
	}
}

func TestEnsureChunkLoaded_StoresOnlySolid(t *testing.T) {
	world := NewWorld(7680, 64000, 640, 42)

	// Load a chunk
	world.EnsureChunkLoaded(0, 1) // Below ground

	// Count tiles in chunk
	tilesInChunk := 0
	for x := 0; x < ChunkSize; x++ {
		for y := ChunkSize; y < ChunkSize*2; y++ {
			if world.tiles[[2]int{x, y}] != nil {
				tilesInChunk++
			}
		}
	}

	// Should have tiles, but not all 256 (due to 20% empty rate)
	if tilesInChunk == 0 {
		t.Error("Chunk should have some solid tiles")
	}

	if tilesInChunk == ChunkSize*ChunkSize {
		t.Error("Chunk should not store all tiles (some should be empty)")
	}
}

func TestGetTileAtGrid_TriggersLoad(t *testing.T) {
	world := NewWorld(7680, 64000, 640, 42)

	// Chunk should not be loaded initially
	if world.loadedChunks[[2]int{5, 5}] {
		t.Error("Chunk should not be loaded initially")
	}

	// Access tile in chunk (5*16, 5*16) = (80, 80)
	_ = world.GetTileAtGrid(80, 80)

	// Chunk should now be loaded
	if !world.loadedChunks[[2]int{5, 5}] {
		t.Error("GetTileAtGrid should trigger chunk load")
	}
}

func TestUpdateChunksAroundPlayer_Loads3x3(t *testing.T) {
	world := NewWorld(7680, 64000, 640, 42)

	// Player at center of chunk (5, 5): pixel (80*64, 80*64) = (5120, 5120)
	playerX := float32(5120)
	playerY := float32(5120)

	world.UpdateChunksAroundPlayer(playerX, playerY)

	// Should load 3×3 grid of chunks around chunk (5, 5)
	expectedChunks := [][2]int{
		{4, 4}, {4, 5}, {4, 6},
		{5, 4}, {5, 5}, {5, 6},
		{6, 4}, {6, 5}, {6, 6},
	}

	for _, chunk := range expectedChunks {
		if !world.loadedChunks[chunk] {
			t.Errorf("Chunk %v should be loaded", chunk)
		}
	}
}

func TestWorld_Deterministic(t *testing.T) {
	world1 := NewWorld(7680, 64000, 640, 12345)
	world2 := NewWorld(7680, 64000, 640, 12345)

	// Query 100 random tile coordinates
	for i := 0; i < 100; i++ {
		x := i * 3
		y := i * 5

		tile1 := world1.GetTileAtGrid(x, y)
		tile2 := world2.GetTileAtGrid(x, y)

		// Both should be nil or both should be non-nil
		if (tile1 == nil) != (tile2 == nil) {
			t.Errorf("Tile existence mismatch at (%d,%d)", x, y)
		}

		// If both exist, types should match
		if tile1 != nil && tile2 != nil {
			if tile1.Type != tile2.Type {
				t.Errorf("Tile type mismatch at (%d,%d): %v vs %v", x, y, tile1.Type, tile2.Type)
			}

			if tile1.Type == entities.TileTypeOre && tile1.OreType != tile2.OreType {
				t.Errorf("Ore type mismatch at (%d,%d): %v vs %v", x, y, tile1.OreType, tile2.OreType)
			}
		}
	}
}

func TestGetTileAt_PixelToGrid(t *testing.T) {
	world := NewWorld(7680, 64000, 640, 42)

	// Pixel (128, 192) should map to grid (2, 3)
	// Grid (2, 3) is in chunk (0, 0)

	tile := world.GetTileAt(128, 192)

	// Should trigger chunk load and return tile
	if !world.loadedChunks[[2]int{0, 0}] {
		t.Error("GetTileAt should trigger chunk load")
	}

	// Tile may be nil (empty) or non-nil (solid) - both are valid
	// Just verify no panic/crash
	_ = tile
}

func TestDigTile_RemovesFromSparseMap(t *testing.T) {
	world := NewWorld(7680, 64000, 640, 42)

	// Ensure ground level chunk is loaded
	world.EnsureChunkLoaded(0, 0)

	// Get tile at ground level (should be dirt)
	pixelX := float32(100)
	pixelY := float32(640) // Ground level

	tileBefore := world.GetTileAt(pixelX, pixelY)
	if tileBefore == nil || !tileBefore.IsDiggable() {
		t.Skip("Ground tile not diggable, skipping test")
	}

	// Dig the tile
	success := world.DigTile(pixelX, pixelY)
	if !success {
		t.Error("Should successfully dig diggable tile")
	}

	// Tile should now be nil (empty)
	tileAfter := world.GetTileAt(pixelX, pixelY)
	if tileAfter != nil {
		t.Error("Dug tile should be nil (removed from sparse map)")
	}
}
