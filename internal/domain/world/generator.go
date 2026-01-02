package world

import (
	"math"
	"math/rand"

	"github.com/Kishlin/drill-game/internal/domain/entities"
)

const ChunkSize = 16 // 16x16 tiles per chunk

// ChunkGenerator handles procedural tile generation using Gaussian ore distribution
type ChunkGenerator struct {
	seed                int64
	emptyRate, dirtRate float32
	groundTileY         int
}

// NewChunkGenerator creates a generator with the given world seed and ground level
func NewChunkGenerator(seed int64, groundLevel float32) *ChunkGenerator {
	return &ChunkGenerator{
		seed:        seed,
		emptyRate:   0.20, // 20% of underground tiles are empty
		dirtRate:    0.65, // 65% of underground tiles are dirt
		groundTileY: int(groundLevel / TileSize),
	}
}

// GenerateTile creates a single tile at the given tile coordinates
// Returns a tile (Dirt, Ore, or Empty) based on Gaussian distribution
func (cg *ChunkGenerator) GenerateTile(tileX, tileY int) *entities.Tile {
	// Above ground: always empty (sky)
	if tileY < cg.groundTileY {
		return entities.NewTile(entities.TileTypeEmpty)
	}

	// Ground level is always solid dirt
	if tileY == cg.groundTileY {
		return entities.NewTile(entities.TileTypeDirt)
	}

	// Seed RNG deterministically for this tile
	rng := cg.seedRNG(tileX, tileY)

	// Roll the tile type (empty, dirt, ore) using cumulative probability ranges
	random := rng.Float32()
	if random < cg.emptyRate {
		// Range [0.0, 0.20) → Empty
		return entities.NewTile(entities.TileTypeEmpty)
	}
	// Range [0.20, 0.85) → Dirt
	// Note: Must include emptyRate in threshold, otherwise dirtRate would only apply to the
	// remaining (1-emptyRate) portion, giving ~45% dirt instead of the intended 65%
	if random < cg.emptyRate+cg.dirtRate {
		return entities.NewTile(entities.TileTypeDirt)
	}
	// Range [0.85, 1.0) → Ore (distributed by Gaussian weight)

	// Calculate ore weights at this depth
	weights := cg.calculateOreWeights(tileY)
	totalWeight := sumWeights(weights)

	// If no ore weight, return dirt
	if totalWeight < 0.01 {
		return entities.NewTile(entities.TileTypeDirt)
	}

	// Select ore type by weighted random
	oreType := cg.selectOreByWeight(rng, weights, totalWeight)
	if oreType == nil {
		return entities.NewTile(entities.TileTypeDirt) // Fallback
	}

	return entities.NewOreTile(*oreType)
}

// gaussianWeight calculates the weight of an ore at a given depth using Gaussian distribution
// Formula: weight = maxWeight × e^(-(depth - peak)² / (2σ²))
func (cg *ChunkGenerator) gaussianWeight(tileY float32, peak, sigma, maxWeight float32) float32 {
	exponent := -math.Pow(float64(tileY-peak), 2) / (2 * math.Pow(float64(sigma), 2))
	return maxWeight * float32(math.Exp(exponent))
}

// calculateOreWeights computes the spawn weight for each ore type at the given depth
func (cg *ChunkGenerator) calculateOreWeights(tileY int) map[entities.OreType]float32 {
	weights := make(map[entities.OreType]float32)

	for oreType, meta := range entities.OreDistributions {
		weight := cg.gaussianWeight(float32(tileY), meta.PeakDepth, meta.Sigma, meta.MaxWeight)
		if weight >= 0.01 {
			weights[oreType] = weight
		}
	}

	return weights
}

// selectOreByWeight performs weighted random selection from available ores
// Iterates in deterministic order (all ore types) to ensure consistent results
func (cg *ChunkGenerator) selectOreByWeight(
	rng *rand.Rand,
	weights map[entities.OreType]float32,
	totalWeight float32,
) *entities.OreType {
	r := rng.Float32() * totalWeight

	// Iterate in fixed order for determinism (map iteration is non-deterministic)
	for _, oreType := range entities.GetAllOreTypes() {
		weight, exists := weights[oreType]
		if !exists {
			continue
		}

		r -= weight
		if r <= 0 {
			return &oreType
		}
	}

	return nil // Shouldn't happen if totalWeight > 0
}

// seedRNG creates a deterministic RNG for this tile based on world seed and coordinates
func (cg *ChunkGenerator) seedRNG(tileX, tileY int) *rand.Rand {
	chunkX := tileX / ChunkSize
	chunkY := tileY / ChunkSize
	localX := tileX % ChunkSize
	localY := tileY % ChunkSize

	// Handle negative coordinates properly
	if localX < 0 {
		localX += ChunkSize
	}
	if localY < 0 {
		localY += ChunkSize
	}

	seed := hashCoordinates(cg.seed, chunkX, chunkY, localX, localY)
	return rand.New(rand.NewSource(seed))
}

// sumWeights calculates the total weight of all ores
func sumWeights(weights map[entities.OreType]float32) float32 {
	total := float32(0)
	for _, w := range weights {
		total += w
	}
	return total
}
