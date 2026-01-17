package worlds

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// HazardConfigLookup provides a backward-compatible interface to hazard configuration data
// It enables systems to query hazard properties without direct dependency on global maps
// It implements the world.HazardLookup interface
type HazardConfigLookup struct {
	config       *HazardConfig
	byHazardType map[entities.HazardType]*HazardDefinition
}

// NewHazardConfigLookup creates a lookup interface for hazard configuration
func NewHazardConfigLookup(config *HazardConfig) *HazardConfigLookup {
	lookup := &HazardConfigLookup{
		config:       config,
		byHazardType: make(map[entities.HazardType]*HazardDefinition),
	}

	// Build index for O(1) lookup
	for i := range config.Hazards {
		lookup.byHazardType[config.Hazards[i].HazardType] = &config.Hazards[i]
	}

	return lookup
}

// GetHardness returns the drilling difficulty multiplier for a hazard type
// For rock (impenetrable), this returns 0
// For lava (drillable), this returns the fixed drilling duration
func (l *HazardConfigLookup) GetHardness(hazardType entities.HazardType) float32 {
	if def, exists := l.byHazardType[hazardType]; exists {
		return def.Hardness
	}
	return 0.0 // Default if hazard not in config
}

// GetName returns the display name for a hazard type
func (l *HazardConfigLookup) GetName(hazardType entities.HazardType) string {
	if def, exists := l.byHazardType[hazardType]; exists {
		return def.Name
	}
	return "Unknown"
}

// GetDistribution returns the Gaussian distribution parameters for a hazard type
func (l *HazardConfigLookup) GetDistribution(hazardType entities.HazardType) *DistributionParams {
	if def, exists := l.byHazardType[hazardType]; exists {
		return &def.Distribution
	}
	return nil
}

// GetAllDistributions implements world.HazardLookup interface
// Returns a map of all hazard types to their distribution parameters
func (l *HazardConfigLookup) GetAllDistributions() map[entities.HazardType]world.HazardDistributionParams {
	distributions := make(map[entities.HazardType]world.HazardDistributionParams)
	for hazardType, def := range l.byHazardType {
		distributions[hazardType] = world.HazardDistributionParams{
			PeakDepth: def.Distribution.PeakDepth,
			Sigma:     def.Distribution.Sigma,
			MaxWeight: def.Distribution.MaxWeight,
		}
	}
	return distributions
}

// GetAllHazards returns all hazard definitions from the config
func (l *HazardConfigLookup) GetAllHazards() []HazardDefinition {
	return l.config.Hazards
}
