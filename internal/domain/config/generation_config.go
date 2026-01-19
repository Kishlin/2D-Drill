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

func DefaultGenerationConfig() GenerationConfig {
	return GenerationConfig{
		Empty: TileDistribution{
			PeakDepth: -500, // Peaks well above ground (negative = above ground level)
			Sigma:     800,  // Wide spread
			MaxWeight: 8.0,  // High weight at surface
		},
		Dirt: TileDistribution{
			PeakDepth: -200, // Peaks slightly above ground
			Sigma:     600,  // Wide spread
			MaxWeight: 20.0, // Dominant at surface
		},
		DirtHardness: 1.0,
		Ores: []OreConfig{
			{
				ID:           "copper",
				Name:         "Copper",
				Distribution: TileDistribution{PeakDepth: -75, Sigma: 120, MaxWeight: 8.0},
				Value:        25,
				Hardness:     1.2,
				Color:        [4]uint8{184, 115, 51, 255}, // Copper brown
			},
			{
				ID:           "iron",
				Name:         "Iron",
				Distribution: TileDistribution{PeakDepth: 70, Sigma: 90, MaxWeight: 5.0},
				Value:        75,
				Hardness:     1.5,
				Color:        [4]uint8{165, 165, 165, 255}, // Iron gray
			},
			{
				ID:           "gold",
				Name:         "Gold",
				Distribution: TileDistribution{PeakDepth: 230, Sigma: 80, MaxWeight: 3.0},
				Value:        300,
				Hardness:     1.8,
				Color:        [4]uint8{255, 215, 0, 255}, // Gold
			},
			{
				ID:           "mythril",
				Name:         "Mythril",
				Distribution: TileDistribution{PeakDepth: 360, Sigma: 70, MaxWeight: 2.2},
				Value:        1500,
				Hardness:     2.1,
				Color:        [4]uint8{64, 224, 208, 255}, // Turquoise
			},
			{
				ID:           "platinum",
				Name:         "Platinum",
				Distribution: TileDistribution{PeakDepth: 500, Sigma: 80, MaxWeight: 1.8},
				Value:        10000,
				Hardness:     2.5,
				Color:        [4]uint8{229, 228, 226, 255}, // Platinum silver
			},
			{
				ID:           "diamond",
				Name:         "Diamond",
				Distribution: TileDistribution{PeakDepth: 600, Sigma: 180, MaxWeight: 0.15},
				Value:        30000,
				Hardness:     3.0,
				Color:        [4]uint8{185, 242, 255, 255}, // Light blue
			},
		},
		Hazards: []HazardConfig{
			{
				ID:            "rock",
				Name:          "Rock",
				Distribution:  TileDistribution{PeakDepth: 650, Sigma: 200, MaxWeight: 15.0},
				Drillable:     false,
				FixedDuration: 0,
				OnDrillDamage: 0,
				OnHitDamage:   0,
				Color:         [4]uint8{80, 80, 80, 255}, // Dark gray
			},
			{
				ID:            "lava",
				Name:          "Lava",
				Distribution:  TileDistribution{PeakDepth: 750, Sigma: 150, MaxWeight: 12.0},
				Drillable:     true,
				FixedDuration: 0.3, // Fixed 0.3s, depth-independent
				OnDrillDamage: 100, // Damage when drilling completes
				OnHitDamage:   0,
				Color:         [4]uint8{255, 69, 0, 255}, // Orange-red
			},
		},
	}
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
