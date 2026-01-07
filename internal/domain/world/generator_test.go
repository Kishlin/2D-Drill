package world

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
)

func TestGaussianWeight_AtPeak(t *testing.T) {
	gen := NewChunkGenerator(42, 640)

	// Test each ore at its peak depth
	tests := []struct {
		name      string
		oreType   entities.OreType
		depth     float32
		minWeight float32
	}{
		{"Copper at peak", entities.OreCopper, -75, 7.5},   // Should be ~8.0
		{"Gold at peak", entities.OreGold, 230, 2.5},       // Should be ~3.0
		{"Diamond at peak", entities.OreDiamond, 600, 0.1}, // Should be ~0.15 (extremely rare)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := entities.OreDistributions[tt.oreType]
			weight := gen.gaussianWeight(tt.depth, meta.PeakDepth, meta.Sigma, meta.MaxWeight)

			if weight < tt.minWeight {
				t.Errorf("Weight at peak for %v = %f, expected >= %f", tt.name, weight, tt.minWeight)
			}
		})
	}
}

func TestGaussianWeight_Symmetry(t *testing.T) {
	gen := NewChunkGenerator(42, 640)
	meta := entities.OreDistributions[entities.OreGold]

	// Weight should be equal at equal distance from peak (230)
	// Test at ±50 tiles from peak: 180 and 280
	weightAbove := gen.gaussianWeight(180, meta.PeakDepth, meta.Sigma, meta.MaxWeight)
	weightBelow := gen.gaussianWeight(280, meta.PeakDepth, meta.Sigma, meta.MaxWeight)

	diff := weightAbove - weightBelow
	if diff < -0.01 || diff > 0.01 {
		t.Errorf("Gaussian should be symmetric: weight at 180=%f, weight at 280=%f", weightAbove, weightBelow)
	}
}

func TestGaussianWeight_FarFromPeak(t *testing.T) {
	gen := NewChunkGenerator(42, 640)
	meta := entities.OreDistributions[entities.OreDiamond]

	// Diamond peaks at 600, should have very low weight at 100 (500px away)
	weight := gen.gaussianWeight(100, meta.PeakDepth, meta.Sigma, meta.MaxWeight)

	if weight > 0.005 {
		t.Errorf("Weight far from peak should be near zero, got %f", weight)
	}
}

func TestCalculateOreWeights_MultipleOres(t *testing.T) {
	gen := NewChunkGenerator(42, 640)

	// At depth 230 (gold's peak), multiple ores should have weights
	weights := gen.calculateOreWeights(int(230 / 64)) // Convert pixels to tile coordinates

	if len(weights) == 0 {
		t.Error("Expected multiple ores at gold's peak depth, got none")
	}

	// Gold should have weight at its peak
	if _, hasGold := weights[entities.OreGold]; !hasGold {
		t.Error("Expected gold to have weight at depth 230")
	}
}

func TestGenerateTile_Deterministic(t *testing.T) {
	gen1 := NewChunkGenerator(12345, 640)
	gen2 := NewChunkGenerator(12345, 640)

	// Same seed + coords = same tile
	for i := 0; i < 10; i++ {
		x, y := i*10, i*20
		tile1 := gen1.GenerateTile(x, y)
		tile2 := gen2.GenerateTile(x, y)

		if tile1.Type != tile2.Type {
			t.Errorf("Tile type mismatch at (%d,%d): %v vs %v", x, y, tile1.Type, tile2.Type)
		}

		if tile1.Type == entities.TileTypeOre && tile1.OreType != tile2.OreType {
			t.Errorf("Ore type mismatch at (%d,%d): %v vs %v", x, y, tile1.OreType, tile2.OreType)
		}
	}
}

func TestGenerateTile_GroundLevel(t *testing.T) {
	gen := NewChunkGenerator(42, 640)
	groundTileY := 10 // 640 / 64 = 10

	// Test multiple X coordinates at ground level
	for x := 0; x < 100; x++ {
		tile := gen.GenerateTile(x, groundTileY)

		if tile.Type != entities.TileTypeDirt {
			t.Errorf("Ground level tile at X=%d should be dirt, got %v", x, tile.Type)
		}
	}
}

func TestGenerateTile_EmptyRate(t *testing.T) {
	gen := NewChunkGenerator(42, 640)

	emptyCount := 0
	totalTiles := 1000

	// Test underground tiles (not ground level)
	for i := 0; i < totalTiles; i++ {
		tile := gen.GenerateTile(i, 50) // Below ground
		if tile.Type == entities.TileTypeEmpty {
			emptyCount++
		}
	}

	// Should be ~23% (allow 18-28% margin for sampling variance)
	emptyRate := float32(emptyCount) / float32(totalTiles)
	if emptyRate < 0.18 || emptyRate > 0.28 {
		t.Errorf("Empty rate = %f, expected ~0.23 (18-28%% range)", emptyRate)
	}
}

func TestGenerateTile_NoOreAtGroundLevel(t *testing.T) {
	gen := NewChunkGenerator(42, 640)
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

func TestSelectOreByWeight_Distribution(t *testing.T) {
	gen := NewChunkGenerator(42, 640)

	// Create simple weight distribution
	weights := map[entities.OreType]float32{
		entities.OreCopper: 10.0,
		entities.OreIron:   5.0,
	}
	totalWeight := sumWeights(weights)

	// Run many selections to verify weighted distribution
	copperCount := 0
	ironCount := 0
	iterations := 1000

	for i := 0; i < iterations; i++ {
		rng := gen.seedRNG(i, 100)
		oreType := gen.selectOreByWeight(rng, weights, totalWeight)

		if oreType == nil {
			t.Error("Should always select an ore when totalWeight > 0")
			continue
		}

		if *oreType == entities.OreCopper {
			copperCount++
		} else if *oreType == entities.OreIron {
			ironCount++
		}
	}

	// Copper should appear ~2x more than iron (10.0 vs 5.0 weight)
	ratio := float32(copperCount) / float32(ironCount)
	if ratio < 1.5 || ratio > 2.5 {
		t.Errorf("Copper/Iron ratio = %f, expected ~2.0", ratio)
	}
}

func TestSumWeights(t *testing.T) {
	weights := map[entities.OreType]float32{
		entities.OreCopper: 10.0,
		entities.OreIron:   5.0,
		entities.OreGold:   2.5,
	}

	total := sumWeights(weights)
	expected := float32(17.5)

	if total != expected {
		t.Errorf("sumWeights = %f, expected %f", total, expected)
	}
}
