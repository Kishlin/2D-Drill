package config

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

type HazardConfig struct {
	ID            string           // Unique identifier (e.g., "rock", "lava")
	Name          string           // Display name (e.g., "Rock", "Lava")
	Distribution  TileDistribution // Gaussian distribution for generation
	Drillable     bool             // false = impenetrable (rock), true = can drill (lava)
	FixedDuration float32          // If drillable: fixed drill time (0 = use depth formula)
	OnDrillDamage float32          // Damage dealt when drilling completes (lava = 100)
	OnHitDamage   float32          // Future: damage on collision (spikes)
	Color         [4]uint8         // RGBA for rendering
}

type GenerationConfig struct {
	Empty        TileDistribution // Gaussian for empty tiles
	Dirt         TileDistribution // Gaussian for dirt tiles
	DirtHardness float32          // Baseline hardness for dirt
	Ores         []OreConfig      // Dynamic list of ores (varies per level)
	Hazards      []HazardConfig   // Dynamic list of hazards (varies per level)
}

func (g *GenerationConfig) GetOreByID(id string) *OreConfig {
	for i := range g.Ores {
		if g.Ores[i].ID == id {
			return &g.Ores[i]
		}
	}
	return nil
}

func (g *GenerationConfig) GetHazardByID(id string) *HazardConfig {
	for i := range g.Hazards {
		if g.Hazards[i].ID == id {
			return &g.Hazards[i]
		}
	}
	return nil
}
