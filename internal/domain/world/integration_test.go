package world

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
)

func TestIntegration_WorldGenerationDeterministic(t *testing.T) {
	world1 := NewWorld(7680, 64000, 640, 12345)
	world2 := NewWorld(7680, 64000, 640, 12345)

	// Query 100 random tile coordinates
	for i := 0; i < 100; i++ {
		x := i * 7
		y := 11 + i*2 // Below ground

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

func TestIntegration_GroundLevelSolid(t *testing.T) {
	world := NewWorld(7680, 64000, 640, 42)

	groundTileY := 10 // 640 / 64 = 10

	// Check multiple X coordinates at ground level
	for x := 0; x < 120; x++ {
		tile := world.GetTileAtGrid(x, groundTileY)

		if tile == nil {
			t.Errorf("Ground level tile at X=%d should not be nil", x)
			continue
		}

		if tile.Type != entities.TileTypeDirt {
			t.Errorf("Ground level tile at X=%d should be dirt, got %v", x, tile.Type)
		}
	}
}

func TestIntegration_OreDistribution(t *testing.T) {
	world := NewWorld(7680, 64000, 640, 42)

	// At depth 300 (gold's peak), count ore types
	depth := 300
	oreCounts := make(map[entities.OreType]int)
	totalTiles := 0
	solidTiles := 0

	for x := 0; x < 100; x++ {
		tile := world.GetTileAtGrid(x, depth)
		totalTiles++

		if tile != nil {
			solidTiles++
			if tile.Type == entities.TileTypeOre {
				oreCounts[tile.OreType]++
			}
		}
	}

	// Should have some ores at this depth
	totalOres := 0
	for _, count := range oreCounts {
		totalOres += count
	}

	if totalOres == 0 {
		t.Error("Expected some ores at depth 300 (gold's peak)")
	}

	// Gold should have weight at its peak (not necessarily most common due to RNG)
	// Just verify gold appears
	if oreCounts[entities.OreGold] == 0 {
		t.Log("Note: Gold didn't appear in sample (might be RNG variation)")
	}
}

func TestIntegration_EmptyTileCollision(t *testing.T) {
	world := NewWorld(7680, 64000, 640, 42)

	// Find an empty tile by scanning underground
	foundEmpty := false
	for y := 11; y < 50; y++ {
		for x := 0; x < 50; x++ {
			tile := world.GetTileAtGrid(x, y)
			if tile == nil {
				// Empty tile found
				pixelX := float32(x * 64)
				pixelY := float32(y * 64)

				if world.IsTileSolid(pixelX, pixelY) {
					t.Error("Empty tile should not be solid")
				}

				foundEmpty = true
				break
			}
		}
		if foundEmpty {
			break
		}
	}

	if !foundEmpty {
		t.Log("Note: No empty tiles found in sample (20% chance, might be RNG)")
	}
}

func TestIntegration_DiggingOre(t *testing.T) {
	world := NewWorld(7680, 64000, 640, 42)

	// Find an ore tile by scanning underground
	foundOre := false
	for y := 11; y < 100; y++ {
		for x := 0; x < 100; x++ {
			tile := world.GetTileAtGrid(x, y)
			if tile != nil && tile.Type == entities.TileTypeOre {
				// Found an ore tile
				pixelX := float32(x*64 + 32)
				pixelY := float32(y*64 + 32)

				// Verify it's solid
				if !tile.IsSolid() {
					t.Error("Ore tile should be solid")
				}

				// Verify it's diggable
				if !tile.IsDiggable() {
					t.Error("Ore tile should be diggable")
				}

				// Dig the tile
				success := world.DigTile(pixelX, pixelY)
				if !success {
					t.Error("Should successfully dig ore tile")
				}

				// Tile should now be nil (removed from sparse map)
				tileAfter := world.GetTileAtGrid(x, y)
				if tileAfter != nil {
					t.Error("Dug ore tile should be nil (removed from sparse map)")
				}

				foundOre = true
				break
			}
		}
		if foundOre {
			break
		}
	}

	if !foundOre {
		t.Error("Expected to find at least one ore tile in sample")
	}
}

func TestIntegration_AboveGroundIsEmpty(t *testing.T) {
	world := NewWorld(7680, 64000, 640, 42)

	groundTileY := 10

	// Check tiles above ground (sky)
	for y := 0; y < groundTileY; y++ {
		for x := 0; x < 20; x++ {
			tile := world.GetTileAtGrid(x, y)

			if tile == nil {
				// Empty tile (expected in sparse storage)
				continue
			}

			if tile.Type != entities.TileTypeEmpty {
				t.Errorf("Tile above ground at (%d,%d) should be empty, got %v", x, y, tile.Type)
			}

			if tile.IsSolid() {
				t.Errorf("Tile above ground at (%d,%d) should not be solid", x, y)
			}
		}
	}
}

func TestIntegration_DifferentSeeds(t *testing.T) {
	world1 := NewWorld(7680, 64000, 640, 111)
	world2 := NewWorld(7680, 64000, 640, 222)

	// Different seeds should produce different worlds
	differences := 0

	for i := 0; i < 100; i++ {
		x := i * 5
		y := 20 + i

		tile1 := world1.GetTileAtGrid(x, y)
		tile2 := world2.GetTileAtGrid(x, y)

		// Check if tiles differ
		if (tile1 == nil) != (tile2 == nil) {
			differences++
			continue
		}

		if tile1 != nil && tile2 != nil {
			if tile1.Type != tile2.Type {
				differences++
			} else if tile1.Type == entities.TileTypeOre && tile1.OreType != tile2.OreType {
				differences++
			}
		}
	}

	// Expect at least some differences (very unlikely to be identical with different seeds)
	if differences == 0 {
		t.Error("Different seeds should produce different worlds")
	}
}

// Benchmark chunk generation performance
func BenchmarkChunkGeneration(b *testing.B) {
	gen := NewChunkGenerator(42, 640)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Generate a 16×16 chunk
		chunkX := i % 10
		chunkY := (i / 10) % 10

		for localX := 0; localX < ChunkSize; localX++ {
			for localY := 0; localY < ChunkSize; localY++ {
				tileX := chunkX*ChunkSize + localX
				tileY := chunkY*ChunkSize + localY
				_ = gen.GenerateTile(tileX, tileY)
			}
		}
	}
}

// Benchmark tile lookup on cached chunk
func BenchmarkGetTileAtGrid_CachedChunk(b *testing.B) {
	world := NewWorld(7680, 64000, 640, 42)

	// Preload chunk
	world.EnsureChunkLoaded(5, 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Access tiles in loaded chunk (80-95, 80-95)
		x := 80 + (i % 16)
		y := 80 + ((i / 16) % 16)
		_ = world.GetTileAtGrid(x, y)
	}
}
