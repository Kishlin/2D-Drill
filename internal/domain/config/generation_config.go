package config

import "fmt"

type TileDistribution struct {
	PeakDepth float32 // Tile Y coordinate where this type is most common
	Sigma     float32 // Standard deviation (spread of distribution)
	MaxWeight float32 // Weight at peak depth (relative spawn chance)
}

type OreConfig struct {
	ID           string           // Unique identifier (e.g., "copper", "uranium")
	Name         string           // Display name (e.g., "Copper", "Uranium")
	Distribution TileDistribution // Gaussian distribution for generation
	Value        int              // Sell price (varies per level)
	Hardness     float32          // Drilling multiplier (varies per level)
	Color        [4]uint8         // RGBA for rendering
}

type HazardEffectType string

const (
	HazardEffectNone       HazardEffectType = "none"
	HazardEffectDamage     HazardEffectType = "damage"
	HazardEffectHeatDamage HazardEffectType = "heat_damage"
	HazardEffectMoney      HazardEffectType = "money"
)

type HazardEffectConfig struct {
	Type               HazardEffectType
	Damage             float32 // For "damage" effect
	BaseDamage         float32 // For "heat_damage" effect
	MaxHeatResistance  float32 // For "heat_damage" (e.g., 320.0)
	MaxDamageReduction float32 // For "heat_damage" (e.g., 0.5 = 50%)
	MoneyAmount        int     // For "money" effect
}

type HazardConfig struct {
	ID            string             // Unique identifier (e.g., "rock", "lava")
	Name          string             // Display name (e.g., "Rock", "Lava")
	Distribution  TileDistribution   // Gaussian distribution for generation
	Drillable     bool               // false = impenetrable (rock), true = can drill (lava)
	FixedDuration float32            // If drillable: fixed drill time (0 = use depth formula)
	OnDrillEffect HazardEffectConfig // Effect applied when drilling completes
	OnHitDamage   float32            // Future: damage on collision (spikes)
	Color         [4]uint8           // RGBA for rendering
}

type GenerationConfig struct {
	Empty        TileDistribution // Gaussian for empty tiles
	Dirt         TileDistribution // Gaussian for dirt tiles
	DirtHardness float32          // Baseline hardness for dirt
	Ores         []OreConfig      // Dynamic list of ores (varies per level)
	Hazards      []HazardConfig   // Dynamic list of hazards (varies per level)
}

func (g *GenerationConfig) GetOreByID(id string) OreConfig {
	for _, ore := range g.Ores {
		if ore.ID == id {
			return ore
		}
	}
	panic(fmt.Sprintf("unknown ore ID: %s", id))
}

func (g *GenerationConfig) GetHazardByID(id string) HazardConfig {
	for _, hazard := range g.Hazards {
		if hazard.ID == id {
			return hazard
		}
	}
	panic(fmt.Sprintf("unknown hazard ID: %s", id))
}
