package world

import (
	"math"
	"math/rand"

	"github.com/Kishlin/drill-game/internal/domain/entities"
)

const ChunkSize = 16 // 16x16 tiles per chunk

// ChunkGenerator handles procedural tile generation using depth-dependent Gaussian distributions
type ChunkGenerator struct {
	seed           int64
	groundTileY    int
	oreLookup      OreLookup
	hazardLookup   HazardLookup
	baseTileConfig *BaseTileConfig
}

// NewChunkGenerator creates a generator with the given world seed and ground level
// Deprecated: Use NewChunkGeneratorWithConfig instead
func NewChunkGenerator(seed int64, groundLevel float32) *ChunkGenerator {
	return &ChunkGenerator{
		seed:        seed,
		groundTileY: int(groundLevel / TileSize),
	}
}

// NewChunkGeneratorWithConfig creates a generator with configuration lookups
func NewChunkGeneratorWithConfig(
	seed int64,
	groundLevel float32,
	oreLookup OreLookup,
	hazardLookup HazardLookup,
	baseTileConfig *BaseTileConfig,
) *ChunkGenerator {
	return &ChunkGenerator{
		seed:           seed,
		groundTileY:    int(groundLevel / TileSize),
		oreLookup:      oreLookup,
		hazardLookup:   hazardLookup,
		baseTileConfig: baseTileConfig,
	}
}

// TileWeights holds spawn weights for all tile types at a given depth
type TileWeights struct {
	Empty   float32
	Dirt    float32
	Ores    map[entities.OreType]float32
	Hazards map[entities.HazardType]float32
}

// GenerateTile creates a single tile at the given tile coordinates
// Returns a tile based on depth-dependent Gaussian distributions for all tile types
func (cg *ChunkGenerator) GenerateTile(tileX, tileY int) *entities.Tile {
	// Above ground: always empty (sky)
	if tileY < cg.groundTileY {
		return entities.NewTile(entities.TileTypeEmpty)
	}

	// Ground level is always solid dirt
	if tileY == cg.groundTileY {
		return entities.NewTile(entities.TileTypeDirt)
	}

	// Underground: use depth-dependent weighted selection
	rng := cg.seedRNG(tileX, tileY)
	weights := cg.calculateAllTileWeights(tileY)
	totalWeight := cg.sumAllWeights(weights)
	return cg.selectTileByWeight(rng, weights, totalWeight)
}

// gaussianWeight calculates the weight of a tile type at a given depth using Gaussian distribution
// Formula: weight = maxWeight × e^(-(depth - peak)² / (2σ²))
func (cg *ChunkGenerator) gaussianWeight(tileY float32, peak, sigma, maxWeight float32) float32 {
	exponent := -math.Pow(float64(tileY-peak), 2) / (2 * math.Pow(float64(sigma), 2))
	return maxWeight * float32(math.Exp(exponent))
}

// calculateAllTileWeights computes spawn weights for all tile types at the given depth
func (cg *ChunkGenerator) calculateAllTileWeights(tileY int) TileWeights {
	weights := TileWeights{
		Ores:    make(map[entities.OreType]float32),
		Hazards: make(map[entities.HazardType]float32),
	}

	// Depth factor: 0.0 at ground, 1.0 at max depth
	depthFactor := float32(tileY-cg.groundTileY) / (1000.0 - float32(cg.groundTileY))
	depthFactor = clamp(depthFactor, 0, 1)

	// Empty and Dirt: use config or hardcoded defaults if no config
	if cg.baseTileConfig != nil {
		weights.Empty = cg.baseTileConfig.EmptyBaseWeight - (depthFactor * cg.baseTileConfig.EmptyDepthReduction)
		weights.Dirt = cg.baseTileConfig.DirtBaseWeight - (depthFactor * cg.baseTileConfig.DirtDepthReduction)
	} else {
		// Fallback to hardcoded values for backward compatibility
		weights.Empty = 8.0 - 7.5*depthFactor
		weights.Dirt = 20.0 - 18.0*depthFactor
	}

	// Ore weights using config or global map as fallback
	if cg.oreLookup != nil {
		for oreType, dist := range cg.oreLookup.GetAllDistributions() {
			weight := cg.gaussianWeight(float32(tileY), dist.PeakDepth, dist.Sigma, dist.MaxWeight)
			if weight >= 0.01 {
				weights.Ores[oreType] = weight
			}
		}
	} else {
		// Fallback to global map for backward compatibility
		for oreType, meta := range entities.OreDistributions {
			weight := cg.gaussianWeight(float32(tileY), meta.PeakDepth, meta.Sigma, meta.MaxWeight)
			if weight >= 0.01 {
				weights.Ores[oreType] = weight
			}
		}
	}

	// Hazard weights using config or global map as fallback
	if cg.hazardLookup != nil {
		for hazardType, dist := range cg.hazardLookup.GetAllDistributions() {
			weight := cg.gaussianWeight(float32(tileY), dist.PeakDepth, dist.Sigma, dist.MaxWeight)
			if weight >= 0.01 {
				weights.Hazards[hazardType] = weight
			}
		}
	} else {
		// Fallback to global map for backward compatibility
		for hazardType, meta := range entities.HazardDistributions {
			weight := cg.gaussianWeight(float32(tileY), meta.PeakDepth, meta.Sigma, meta.MaxWeight)
			if weight >= 0.01 {
				weights.Hazards[hazardType] = weight
			}
		}
	}

	return weights
}

// calculateOreWeights computes the spawn weight for each ore type at the given depth (legacy, kept for compatibility)
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

// sumAllWeights calculates the total weight of all tile types
func (cg *ChunkGenerator) sumAllWeights(weights TileWeights) float32 {
	total := weights.Empty + weights.Dirt

	for _, w := range weights.Ores {
		total += w
	}

	for _, w := range weights.Hazards {
		total += w
	}

	return total
}

// selectTileByWeight performs weighted random selection from all tile types
// Maintains deterministic order for consistent world generation
func (cg *ChunkGenerator) selectTileByWeight(
	rng *rand.Rand,
	weights TileWeights,
	totalWeight float32,
) *entities.Tile {
	if totalWeight < 0.01 {
		return entities.NewTile(entities.TileTypeDirt) // Fallback
	}

	r := rng.Float32() * totalWeight

	// Check Empty
	r -= weights.Empty
	if r <= 0 {
		return entities.NewTile(entities.TileTypeEmpty)
	}

	// Check Dirt
	r -= weights.Dirt
	if r <= 0 {
		return entities.NewTile(entities.TileTypeDirt)
	}

	// Check Ores (deterministic order)
	for _, oreType := range entities.GetAllOreTypes() {
		if weight, exists := weights.Ores[oreType]; exists {
			r -= weight
			if r <= 0 {
				return entities.NewOreTile(oreType)
			}
		}
	}

	// Check Hazards (deterministic order)
	for _, hazardType := range entities.GetAllHazardTypes() {
		if weight, exists := weights.Hazards[hazardType]; exists {
			r -= weight
			if r <= 0 {
				return entities.NewHazardTile(hazardType)
			}
		}
	}

	return entities.NewTile(entities.TileTypeDirt) // Fallback
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

// sumWeights calculates the total weight of all ores (legacy, kept for compatibility)
func sumWeights(weights map[entities.OreType]float32) float32 {
	total := float32(0)
	for _, w := range weights {
		total += w
	}
	return total
}

func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
