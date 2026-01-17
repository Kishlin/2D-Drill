package worlds

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// OreConfigLookup provides a backward-compatible interface to ore configuration data
// It enables systems to query ore properties without direct dependency on global maps
// It implements the world.OreLookup interface
type OreConfigLookup struct {
	config    *OreConfig
	byOreType map[entities.OreType]*OreDefinition
}

// NewOreConfigLookup creates a lookup interface for ore configuration
func NewOreConfigLookup(config *OreConfig) *OreConfigLookup {
	lookup := &OreConfigLookup{
		config:    config,
		byOreType: make(map[entities.OreType]*OreDefinition),
	}

	// Build index for O(1) lookup
	for i := range config.Ores {
		lookup.byOreType[config.Ores[i].OreType] = &config.Ores[i]
	}

	return lookup
}

// GetValue returns the market sell price for an ore type
func (l *OreConfigLookup) GetValue(oreType entities.OreType) int {
	if def, exists := l.byOreType[oreType]; exists {
		return def.Value
	}
	return 0 // Ore not in config
}

// GetHardness returns the drilling difficulty multiplier for an ore type
func (l *OreConfigLookup) GetHardness(oreType entities.OreType) float32 {
	if def, exists := l.byOreType[oreType]; exists {
		return def.Hardness
	}
	return 1.0 // Default hardness if ore not in config
}

// GetName returns the display name for an ore type
func (l *OreConfigLookup) GetName(oreType entities.OreType) string {
	if def, exists := l.byOreType[oreType]; exists {
		return def.Name
	}
	return "Unknown"
}

// GetDistribution returns the Gaussian distribution parameters for an ore type
func (l *OreConfigLookup) GetDistribution(oreType entities.OreType) *DistributionParams {
	if def, exists := l.byOreType[oreType]; exists {
		return &def.Distribution
	}
	return nil
}

// GetAllDistributions implements world.OreLookup interface
// Returns a map of all ore types to their distribution parameters
func (l *OreConfigLookup) GetAllDistributions() map[entities.OreType]world.OreDistributionParams {
	distributions := make(map[entities.OreType]world.OreDistributionParams)
	for oreType, def := range l.byOreType {
		distributions[oreType] = world.OreDistributionParams{
			PeakDepth: def.Distribution.PeakDepth,
			Sigma:     def.Distribution.Sigma,
			MaxWeight: def.Distribution.MaxWeight,
		}
	}
	return distributions
}

// GetAllOres returns all ore definitions from the config
func (l *OreConfigLookup) GetAllOres() []OreDefinition {
	return l.config.Ores
}
