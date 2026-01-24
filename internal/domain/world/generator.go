package world

import (
	"math"
	"math/rand"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
)

const ChunkSize = 16 // 16x16 tiles per chunk

type ChunkGenerator struct {
	seed           int64
	groundTileY    int
	genCfg         config.GenerationConfig
	bossRoomConfig config.BossRoomConfig
	bossRoomStartY int // Starting tile Y of boss room (in pixels / TileSize)
	floorStartY    int // Starting tile Y of floor area
	floorEndY      int // Ending tile Y of floor area (world bottom)
}

func NewChunkGeneratorFromConfig(seed int64, groundLevel, worldHeight float32, genCfg config.GenerationConfig, bossRoomCfg config.BossRoomConfig) *ChunkGenerator {
	worldHeightTiles := int(worldHeight / TileSize)
	floorHeightTiles := int(bossRoomCfg.FloorHeight)
	roomHeightTiles := int(bossRoomCfg.RoomHeight / TileSize)

	return &ChunkGenerator{
		seed:           seed,
		groundTileY:    int(groundLevel / TileSize),
		genCfg:         genCfg,
		bossRoomConfig: bossRoomCfg,
		floorEndY:      worldHeightTiles,
		floorStartY:    worldHeightTiles - floorHeightTiles,
		bossRoomStartY: worldHeightTiles - floorHeightTiles - roomHeightTiles,
	}
}

// TileWeights holds spawn weights for all tile types at a given depth
type TileWeights struct {
	Empty   float32
	Dirt    float32
	Ores    map[string]float32 // Keyed by ore ID
	Hazards map[string]float32 // Keyed by hazard ID
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

	// Boss room (empty air for the room itself)
	if cg.isBossRoomTile(tileY) {
		return entities.NewTile(entities.TileTypeEmpty)
	}

	// Floor (solid, not drillable, not nukeable)
	if cg.isFloorTile(tileY) {
		return entities.NewTile(entities.TileTypeFloor)
	}

	// Underground: use depth-dependent weighted selection
	rng := cg.seedRNG(tileX, tileY)
	weights := cg.calculateAllTileWeights(tileY)
	totalWeight := cg.sumAllWeights(weights)
	return cg.selectTileByWeight(rng, weights, totalWeight)
}

func (cg *ChunkGenerator) isBossRoomTile(tileY int) bool {
	return tileY >= cg.bossRoomStartY && tileY < cg.floorStartY
}

func (cg *ChunkGenerator) isFloorTile(tileY int) bool {
	return tileY >= cg.floorStartY && tileY < cg.floorEndY
}

// gaussianWeight calculates the weight of a tile type at a given depth using Gaussian distribution
// Formula: weight = maxWeight × e^(-(depth - peak)² / (2σ²))
func (cg *ChunkGenerator) gaussianWeight(tileY float32, dist config.TileDistribution) float32 {
	exponent := -math.Pow(float64(tileY-dist.PeakDepth), 2) / (2 * math.Pow(float64(dist.Sigma), 2))
	return dist.MaxWeight * float32(math.Exp(exponent))
}

// calculateAllTileWeights computes spawn weights for all tile types at the given depth
func (cg *ChunkGenerator) calculateAllTileWeights(tileY int) TileWeights {
	weights := TileWeights{
		Ores:    make(map[string]float32),
		Hazards: make(map[string]float32),
	}

	// Empty: use Gaussian distribution from config
	weights.Empty = cg.gaussianWeight(float32(tileY), cg.genCfg.Empty)

	// Dirt: use Gaussian distribution from config
	weights.Dirt = cg.gaussianWeight(float32(tileY), cg.genCfg.Dirt)

	// Ore weights (Gaussian from config)
	for _, oreCfg := range cg.genCfg.Ores {
		weight := cg.gaussianWeight(float32(tileY), oreCfg.Distribution)
		if weight >= 0.01 {
			weights.Ores[oreCfg.ID] = weight
		}
	}

	// Hazard weights (Gaussian from config)
	for _, hazardCfg := range cg.genCfg.Hazards {
		weight := cg.gaussianWeight(float32(tileY), hazardCfg.Distribution)
		if weight >= 0.01 {
			weights.Hazards[hazardCfg.ID] = weight
		}
	}

	return weights
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

	// Check Ores (deterministic order based on config order)
	for _, oreCfg := range cg.genCfg.Ores {
		if weight, exists := weights.Ores[oreCfg.ID]; exists {
			r -= weight
			if r <= 0 {
				return entities.NewOreTileByID(oreCfg.ID)
			}
		}
	}

	// Check Hazards (deterministic order based on config order)
	for i := range cg.genCfg.Hazards {
		hazardCfg := &cg.genCfg.Hazards[i]
		if weight, exists := weights.Hazards[hazardCfg.ID]; exists {
			r -= weight
			if r <= 0 {
				return entities.NewHazardTileByID(hazardCfg.ID, hazardCfg)
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

func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// GetGenerationConfig returns the generator's configuration (for systems that need ore/hazard info)
func (cg *ChunkGenerator) GetGenerationConfig() config.GenerationConfig {
	return cg.genCfg
}
