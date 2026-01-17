package world

import "github.com/Kishlin/drill-game/internal/domain/entities"

// OreLookup interface defines methods to query ore configuration
// Implemented by worlds.OreConfigLookup
type OreLookup interface {
	GetValue(oreType entities.OreType) int
	GetHardness(oreType entities.OreType) float32
	GetAllDistributions() map[entities.OreType]OreDistributionParams
}

// OreDistributionParams represents Gaussian distribution parameters for ore generation
type OreDistributionParams struct {
	PeakDepth float32
	Sigma     float32
	MaxWeight float32
}

// HazardLookup interface defines methods to query hazard configuration
// Implemented by worlds.HazardConfigLookup
type HazardLookup interface {
	GetHardness(hazardType entities.HazardType) float32
	GetAllDistributions() map[entities.HazardType]HazardDistributionParams
}

// HazardDistributionParams represents Gaussian distribution parameters for hazard generation
type HazardDistributionParams struct {
	PeakDepth float32
	Sigma     float32
	MaxWeight float32
}

// BaseTileConfig defines the weight formulas for empty and dirt tiles
type BaseTileConfig struct {
	EmptyBaseWeight       float32
	EmptyDepthReduction   float32
	DirtBaseWeight        float32
	DirtDepthReduction    float32
	DirtHardness          float32
}
