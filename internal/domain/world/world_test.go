package world

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
)

// Test helper - creates minimal world config for tests
func testWorldConfig(width, height, groundLevel float32, seed int64) config.WorldConfig {
	return config.WorldConfig{
		Width:       width,
		Height:      height,
		GroundLevel: groundLevel,
		Seed:        seed,
		PlayerSpawn: config.PlayerSpawn{X: width / 2, Y: groundLevel - 10},
		BuildingLayout: config.BuildingLayout{
			HospitalX: 0, FuelStationX: 0, MarketX: 0, UpgradeShopX: 0, ItemShopX: 0,
		},
	}
}

func testGenConfig() config.GenerationConfig {
	return config.GenerationConfig{
		Empty:        config.TileDistribution{PeakDepth: 0, Sigma: 1000, MaxWeight: 20},
		Dirt:         config.TileDistribution{PeakDepth: 0, Sigma: 500, MaxWeight: 100},
		DirtHardness: 1.0,
		Ores: []config.OreConfig{
			{ID: "copper", Name: "Copper", Value: 25, Hardness: 1.2, Distribution: config.TileDistribution{PeakDepth: -75, Sigma: 120, MaxWeight: 8}, Color: [4]uint8{184, 115, 51, 255}},
			{ID: "gold", Name: "Gold", Value: 300, Hardness: 1.8, Distribution: config.TileDistribution{PeakDepth: 230, Sigma: 80, MaxWeight: 3}, Color: [4]uint8{255, 215, 0, 255}},
		},
		Hazards: []config.HazardConfig{
			{ID: "rock", Name: "Rock", Drillable: false, Distribution: config.TileDistribution{PeakDepth: 650, Sigma: 200, MaxWeight: 15}, Color: [4]uint8{80, 80, 80, 255}},
			{ID: "lava", Name: "Lava", Drillable: true, FixedDuration: 0.3, OnDrillDamage: 100, Distribution: config.TileDistribution{PeakDepth: 750, Sigma: 150, MaxWeight: 12}, Color: [4]uint8{255, 100, 0, 255}},
		},
	}
}

func testWorld(width, height, groundLevel float32, seed int64) *World {
	return NewWorldFromConfig(testWorldConfig(width, height, groundLevel, seed), testGenConfig(), testBossRoomConfig())
}

func TestEnsureChunkLoaded_OnlyOnce(t *testing.T) {
	world := testWorld(7680, 64000, 640, 42)

	// First load
	world.EnsureChunkLoaded(0, 0)
	if world.loadedChunks[[2]int{0, 0}] == false {
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
	world := testWorld(7680, 64000, 640, 42)

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
	world := testWorld(7680, 64000, 640, 42)

	// Chunk should not be loaded initially
	if world.loadedChunks[[2]int{5, 5}] {
		t.Error("Chunk should not be loaded initially")
	}

	// Access tile in chunk (5*16, 5*16) = (80, 80)
	_ = world.GetTileAtGrid(80, 80)

	// Chunk should now be loaded
	if world.loadedChunks[[2]int{5, 5}] == false {
		t.Error("GetTileAtGrid should trigger chunk load")
	}
}

func TestUpdateChunksAroundPlayer_Loads3x3(t *testing.T) {
	world := testWorld(7680, 64000, 640, 42)

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
		if world.loadedChunks[chunk] == false {
			t.Errorf("Chunk %v should be loaded", chunk)
		}
	}
}

func TestWorld_Deterministic(t *testing.T) {
	world1 := testWorld(7680, 64000, 640, 12345)
	world2 := testWorld(7680, 64000, 640, 12345)

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

			if tile1.Type == entities.TileTypeOre && tile1.OreID != tile2.OreID {
				t.Errorf("Ore ID mismatch at (%d,%d): %v vs %v", x, y, tile1.OreID, tile2.OreID)
			}
		}
	}
}

func TestGetTileAt_PixelToGrid(t *testing.T) {
	world := testWorld(7680, 64000, 640, 42)

	// Pixel (128, 192) should map to grid (2, 3)
	// Grid (2, 3) is in chunk (0, 0)

	tile := world.GetTileAt(128, 192)

	// Should trigger chunk load and return tile
	if world.loadedChunks[[2]int{0, 0}] == false {
		t.Error("GetTileAt should trigger chunk load")
	}

	// Tile may be nil (empty) or non-nil (solid) - both are valid
	// Just verify no panic/crash
	_ = tile
}

func TestDrillTile_RemovesFromSparseMap(t *testing.T) {
	world := testWorld(7680, 64000, 640, 42)

	// Ensure ground level chunk is loaded
	world.EnsureChunkLoaded(0, 0)

	// Get tile at ground level (should be dirt)
	pixelX := float32(100)
	pixelY := float32(640) // Ground level

	tileBefore := world.GetTileAt(pixelX, pixelY)
	if tileBefore == nil || tileBefore.IsDrillable() == false {
		t.Skip("Ground tile not drillable, skipping test")
	}

	// Drill the tile
	_, success := world.DrillTile(pixelX, pixelY)
	if success == false {
		t.Error("Should successfully drill drillable tile")
	}

	// Tile should now be nil (empty)
	tileAfter := world.GetTileAt(pixelX, pixelY)
	if tileAfter != nil {
		t.Error("Drilled tile should be nil (removed from sparse map)")
	}
}

// === Hazard Tile and Bomb Tests ===

func TestNukeTileAtGrid_RemovesRock(t *testing.T) {
	world := testWorld(7680, 64000, 640, 42)

	// Place rock tile at grid (100, 500) directly in the sparse map
	gridX, gridY := 100, 500
	rockCfg := &config.HazardConfig{Drillable: false}
	rockTile := entities.NewHazardTileByID("rock", rockCfg)
	world.tiles[[2]int{gridX, gridY}] = rockTile

	// Nuke the tile (use bomb)
	tile, success := world.NukeTileAtGrid(gridX, gridY)
	if success == false {
		t.Error("Should successfully nuke rock tile")
	}

	// Should return the removed rock tile
	if tile == nil || tile.Type != entities.TileTypeRock {
		t.Error("Should return rock tile")
	}

	// Tile should now be removed
	if world.tiles[[2]int{gridX, gridY}] != nil {
		t.Error("Rock tile should be removed after nuke")
	}
}

func TestNukeTileAtGrid_RemovesLava(t *testing.T) {
	world := testWorld(7680, 64000, 640, 42)

	// Place lava tile directly in sparse map
	gridX, gridY := 100, 500
	lavaTile := entities.NewHazardTileByID("lava", nil)
	world.tiles[[2]int{gridX, gridY}] = lavaTile

	// Nuke the tile
	tile, success := world.NukeTileAtGrid(gridX, gridY)
	if success == false {
		t.Error("Should successfully nuke lava tile")
	}

	if tile == nil || tile.Type != entities.TileTypeLava {
		t.Error("Should return lava tile")
	}

	// Tile should now be removed
	if world.tiles[[2]int{gridX, gridY}] != nil {
		t.Error("Lava tile should be removed after nuke")
	}
}

func TestNukeTileAtGrid_BypassesDrillability(t *testing.T) {
	world := testWorld(7680, 64000, 640, 42)

	// Place rock tile (not drillable)
	gridX, gridY := 100, 500
	rockCfg := &config.HazardConfig{Drillable: false}
	rockTile := entities.NewHazardTileByID("rock", rockCfg)
	world.SetTile(gridX, gridY, rockTile)

	// DrillTileAtGrid should fail (rock not drillable)
	_, drillSuccess := world.DrillTileAtGrid(gridX, gridY)
	if drillSuccess {
		t.Error("Should NOT be able to drill rock tile")
	}

	// But rock should still exist
	if world.GetTileAtGrid(gridX, gridY) == nil {
		t.Error("Rock should still exist after failed drill")
	}

	// NukeTileAtGrid should succeed (bypasses drillability)
	_, nukeSuccess := world.NukeTileAtGrid(gridX, gridY)
	if nukeSuccess == false {
		t.Error("Should be able to nuke rock tile (bypasses drillability)")
	}

	// Rock should now be removed
	if world.GetTileAtGrid(gridX, gridY) != nil {
		t.Error("Rock should be removed after nuke")
	}
}

func TestNukeTileAtGrid_DoesNotAffectEmpty(t *testing.T) {
	world := testWorld(7680, 64000, 640, 42)

	// Try to nuke at a location with no tile (empty)
	gridX, gridY := 100, 500
	tile, success := world.NukeTileAtGrid(gridX, gridY)

	if success {
		t.Error("Should not be able to nuke empty tile")
	}

	if tile != nil {
		t.Error("Should return nil when nuking empty space")
	}
}

func TestNukeTileAtGrid_RemovesDirt(t *testing.T) {
	world := testWorld(7680, 64000, 640, 42)

	// Place dirt tile directly in sparse map (should also be removable with nuke)
	gridX, gridY := 100, 500
	dirtTile := entities.NewTile(entities.TileTypeDirt)
	world.tiles[[2]int{gridX, gridY}] = dirtTile

	// Nuke the dirt tile
	tile, success := world.NukeTileAtGrid(gridX, gridY)
	if success == false {
		t.Error("Should be able to nuke dirt tile")
	}

	if tile == nil || tile.Type != entities.TileTypeDirt {
		t.Error("Should return dirt tile")
	}

	// Dirt should now be removed
	if world.tiles[[2]int{gridX, gridY}] != nil {
		t.Error("Dirt tile should be removed after nuke")
	}
}
