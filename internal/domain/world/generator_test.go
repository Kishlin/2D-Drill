package world

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
)

func TestGaussianWeight_AtPeak(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())

	// Test each ore at its peak depth
	for _, oreCfg := range genCfg.Ores {
		t.Run(oreCfg.ID+" at peak", func(t *testing.T) {
			weight := gen.gaussianWeight(oreCfg.Distribution.PeakDepth, oreCfg.Distribution)

			// Weight at peak should be close to MaxWeight
			expectedMin := oreCfg.Distribution.MaxWeight * 0.9
			if weight < expectedMin {
				t.Errorf("Weight at peak for %s = %f, expected >= %f", oreCfg.ID, weight, expectedMin)
			}
		})
	}
}

func TestGaussianWeight_Symmetry(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())

	// Test with gold ore
	goldCfg := genCfg.GetOreByID("gold")

	// Weight should be equal at equal distance from peak (230)
	// Test at ±50 tiles from peak: 180 and 280
	weightAbove := gen.gaussianWeight(goldCfg.Distribution.PeakDepth-50, goldCfg.Distribution)
	weightBelow := gen.gaussianWeight(goldCfg.Distribution.PeakDepth+50, goldCfg.Distribution)

	diff := weightAbove - weightBelow
	if diff < -0.01 || diff > 0.01 {
		t.Errorf("Gaussian should be symmetric: weight above=%f, weight below=%f", weightAbove, weightBelow)
	}
}

func TestGaussianWeight_FarFromPeak(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())

	// Test with diamond ore (peaks deep)
	diamondCfg := genCfg.GetOreByID("diamond")

	// Diamond peaks at 600, should have very low weight at 100 (500px away)
	weight := gen.gaussianWeight(100, diamondCfg.Distribution)

	if weight > 0.005 {
		t.Errorf("Weight far from peak should be near zero, got %f", weight)
	}
}

func TestCalculateAllTileWeights_MultipleOres(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())

	// Get gold's peak depth
	goldCfg := genCfg.GetOreByID("gold")

	// At depth near gold's peak, multiple ores should have weights
	tileY := int(goldCfg.Distribution.PeakDepth)
	weights := gen.calculateAllTileWeights(tileY)

	if len(weights.Ores) == 0 {
		t.Error("Expected multiple ores at gold's peak depth, got none")
	}

	// Gold should have weight at its peak
	if _, hasGold := weights.Ores["gold"]; !hasGold {
		t.Error("Expected gold to have weight at its peak depth")
	}
}

func TestGenerateTile_Deterministic(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen1 := NewChunkGeneratorFromConfig(12345, 640, 64000, genCfg, testBossRoomConfig())
	gen2 := NewChunkGeneratorFromConfig(12345, 640, 64000, genCfg, testBossRoomConfig())

	// Same seed + coords = same tile
	for i := 0; i < 10; i++ {
		x, y := i*10, i*20
		tile1 := gen1.GenerateTile(x, y)
		tile2 := gen2.GenerateTile(x, y)

		if tile1.Type != tile2.Type {
			t.Errorf("Tile type mismatch at (%d,%d): %v vs %v", x, y, tile1.Type, tile2.Type)
		}

		if tile1.Type == entities.TileTypeOre && tile1.OreID != tile2.OreID {
			t.Errorf("Ore ID mismatch at (%d,%d): %v vs %v", x, y, tile1.OreID, tile2.OreID)
		}
	}
}

func TestGenerateTile_GroundLevel(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())
	groundTileY := 10 // 640 / 64 = 10

	// Test multiple X coordinates at ground level
	for x := 0; x < 100; x++ {
		tile := gen.GenerateTile(x, groundTileY)

		if tile.Type != entities.TileTypeDirt {
			t.Errorf("Ground level tile at X=%d should be dirt, got %v", x, tile.Type)
		}
	}
}

func TestGenerateTile_NoOreAtGroundLevel(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())
	groundTileY := 10

	// Ground level should never generate ore
	for x := 0; x < 100; x++ {
		tile := gen.GenerateTile(x, groundTileY)

		if tile.Type == entities.TileTypeOre {
			t.Errorf("Ground level should never be ore, got ore at X=%d", x)
		}
	}
}

func TestHashCoordinates_Deterministic(t *testing.T) {
	seed := int64(12345)

	hash1 := hashCoordinates(seed, 5, 10, 3, 7)
	hash2 := hashCoordinates(seed, 5, 10, 3, 7)

	if hash1 != hash2 {
		t.Error("Hash should be deterministic for same inputs")
	}
}

func TestHashCoordinates_Unique(t *testing.T) {
	seed := int64(12345)

	hash1 := hashCoordinates(seed, 0, 0, 0, 0)
	hash2 := hashCoordinates(seed, 0, 0, 0, 1)
	hash3 := hashCoordinates(seed, 0, 0, 1, 0)
	hash4 := hashCoordinates(seed, 1, 0, 0, 0)

	// Different coordinates should produce different hashes
	if hash1 == hash2 || hash1 == hash3 || hash1 == hash4 {
		t.Error("Different coordinates should produce different hashes")
	}
}

// === Hazard Tile Tests ===

func TestCalculateAllTileWeights_RockAtShallowDepth(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())

	// At shallow depth (~tile 50, ~40% depth), rock should have minimal weight
	weights := gen.calculateAllTileWeights(50)

	if rockWeight, exists := weights.Hazards["rock"]; exists && rockWeight > 0.5 {
		t.Errorf("Rock should be rare at shallow depth, got weight %f", rockWeight)
	}
}

func TestCalculateAllTileWeights_HazardsAtDeepDepth(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())

	// At very deep depth (~tile 800, ~80%+ depth), hazards should have significant weight
	weights := gen.calculateAllTileWeights(800)

	totalHazardWeight := float32(0)
	for _, w := range weights.Hazards {
		totalHazardWeight += w
	}

	if totalHazardWeight < 10.0 {
		t.Errorf("Hazards should have significant weight at deep depth, got total %f", totalHazardWeight)
	}
}

func TestGenerateTile_NoHazardsAtSurface(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())

	// Generate 100 tiles at shallow depth (near ground)
	hazardCount := 0
	for i := 0; i < 100; i++ {
		tile := gen.GenerateTile(i, 15) // Just 5 tiles below ground level (tile 10)
		if tile.Type == entities.TileTypeHazard {
			hazardCount++
		}
	}

	// Should have very few or no hazards at surface
	if hazardCount > 5 {
		t.Errorf("Should have very few hazards at shallow depth, got %d/100", hazardCount)
	}
}

func TestGenerateTile_HazardsAtDeepDepth(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())

	// Generate 100 tiles at very deep depth (80%+)
	hazardCount := 0
	for i := 0; i < 100; i++ {
		tile := gen.GenerateTile(i, 800) // ~80% depth
		if tile.Type == entities.TileTypeHazard {
			hazardCount++
		}
	}

	// Should have significant number of hazards at depth
	if hazardCount < 20 {
		t.Errorf("Should have many hazards at deep depth (80+), got %d/100", hazardCount)
	}
}

func TestGenerateTile_RockAndLavaAreDistinct(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())

	rockCount := 0
	lavaCount := 0

	// At deep depth, collect both rock and lava statistics
	for i := 0; i < 500; i++ {
		tile := gen.GenerateTile(i, 800)
		if tile.Type == entities.TileTypeHazard && tile.Drillable == false {
			rockCount++
		} else if tile.Type == entities.TileTypeHazard {
			lavaCount++
		}
	}

	// Both should appear, with rock more common than lava (per default config)
	if rockCount == 0 {
		t.Error("Rock tiles should appear at deep depth")
	}
	if lavaCount == 0 {
		t.Error("Lava tiles should appear at deep depth")
	}
	if rockCount < lavaCount {
		t.Errorf("Rock should be more common than lava at deep depth: rock=%d, lava=%d", rockCount, lavaCount)
	}
}

func TestGenerateTile_HazardDeterministic(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen1 := NewChunkGeneratorFromConfig(12345, 640, 64000, genCfg, testBossRoomConfig())
	gen2 := NewChunkGeneratorFromConfig(12345, 640, 64000, genCfg, testBossRoomConfig())

	// Same seed + coords = same hazard type
	for i := 0; i < 10; i++ {
		x, y := i*10, i*20+800 // Deep depth
		tile1 := gen1.GenerateTile(x, y)
		tile2 := gen2.GenerateTile(x, y)

		if tile1.Type == entities.TileTypeHazard && tile2.Type == entities.TileTypeHazard {
			if tile1.HazardID != tile2.HazardID {
				t.Errorf("Hazard ID should be deterministic at (%d,%d): %v vs %v",
					x, y, tile1.HazardID, tile2.HazardID)
			}
		}
	}
}

func TestNewChunkGeneratorFromConfig(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())

	// Should be able to generate tiles
	tile := gen.GenerateTile(10, 15)
	if tile == nil {
		t.Error("Generator should create tiles")
	}

	// Should have access to config - verify by checking a field
	cfg := gen.GetGenerationConfig()
	if cfg.DirtHardness == 0 && len(cfg.Ores) == 0 {
		t.Error("Generator should have config accessible")
	}
}

func TestSumAllWeights(t *testing.T) {
	genCfg := testGeneratorConfig()
	gen := NewChunkGeneratorFromConfig(42, 640, 64000, genCfg, testBossRoomConfig())
	weights := gen.calculateAllTileWeights(50)

	total := gen.sumAllWeights(weights)

	// Should be sum of all tile types
	expected := weights.Empty + weights.Dirt
	for _, w := range weights.Ores {
		expected += w
	}
	for _, w := range weights.Hazards {
		expected += w
	}

	if total != expected {
		t.Errorf("sumAllWeights = %f, expected %f", total, expected)
	}
}

// Test helpers

func testBossRoomConfig() config.BossRoomConfig {
	return config.BossRoomConfig{
		BossType:    "test_boss",
		FloorType:   config.FloorConcrete,
		RoomHeight:  680.0,
		FloorHeight: 6.0,
	}
}

func testGeneratorConfig() config.GenerationConfig {
	return config.GenerationConfig{
		Empty:        config.TileDistribution{PeakDepth: 0, Sigma: 1000, MaxWeight: 20},
		Dirt:         config.TileDistribution{PeakDepth: 0, Sigma: 500, MaxWeight: 100},
		DirtHardness: 1.0,
		Ores: []config.OreConfig{
			{ID: "copper", Name: "Copper", Value: 25, Hardness: 1.2, Distribution: config.TileDistribution{PeakDepth: -75, Sigma: 120, MaxWeight: 8}, Color: [4]uint8{184, 115, 51, 255}},
			{ID: "gold", Name: "Gold", Value: 300, Hardness: 1.8, Distribution: config.TileDistribution{PeakDepth: 230, Sigma: 80, MaxWeight: 3}, Color: [4]uint8{255, 215, 0, 255}},
			{ID: "diamond", Name: "Diamond", Value: 500, Hardness: 2.5, Distribution: config.TileDistribution{PeakDepth: 600, Sigma: 100, MaxWeight: 2}, Color: [4]uint8{185, 242, 255, 255}},
		},
		Hazards: []config.HazardConfig{
			{ID: "rock", Name: "Rock", Drillable: false, Distribution: config.TileDistribution{PeakDepth: 650, Sigma: 200, MaxWeight: 15}, OnDrillEffect: config.HazardEffectConfig{Type: config.HazardEffectNone}, Color: [4]uint8{80, 80, 80, 255}},
			{ID: "lava", Name: "Lava", Drillable: true, FixedDuration: 0.3, OnDrillEffect: config.HazardEffectConfig{Type: config.HazardEffectHeatDamage, BaseDamage: 100, MaxHeatResistance: 320, MaxDamageReduction: 0.5}, Distribution: config.TileDistribution{PeakDepth: 750, Sigma: 150, MaxWeight: 12}, Color: [4]uint8{255, 100, 0, 255}},
		},
	}
}
