package world

import "fmt"

type PlayerSpawn struct {
	X float32
	Y float32
}

// BuildingLayout defines the X positions of all buildings on the surface
// Y positions are always calculated based on GroundLevel - building height
type BuildingLayout struct {
	MarketX      float32
	FuelStationX float32
	HospitalX    float32
	UpgradeShopX float32
	ItemShopX    float32
}

type WorldConfig struct {
	Width          float32
	Height         float32
	GroundLevel    float32
	Seed           int64
	PlayerSpawn    PlayerSpawn
	BuildingLayout BuildingLayout
}

// NewWorldConfigForTesting creates a config suitable for tests with minimal building positions
// This avoids needing to specify all building positions in test code
func NewWorldConfigForTesting(width, height, groundLevel float32, seed int64) *WorldConfig {
	return &WorldConfig{
		Width:       width,
		Height:      height,
		GroundLevel: groundLevel,
		Seed:        seed,
		PlayerSpawn: PlayerSpawn{
			X: width / 2,
			Y: groundLevel - 10,
		},
		BuildingLayout: BuildingLayout{
			MarketX:      width / 2,
			FuelStationX: width / 2,
			HospitalX:    width / 2,
			UpgradeShopX: width / 2,
			ItemShopX:    width / 2,
		},
	}
}

func (c *WorldConfig) Validate() error {
	if c.Width <= 0 {
		return fmt.Errorf("world width must be positive, got %f", c.Width)
	}
	if c.Height <= 0 {
		return fmt.Errorf("world height must be positive, got %f", c.Height)
	}
	if c.GroundLevel < 0 || c.GroundLevel > c.Height {
		return fmt.Errorf("ground level %f must be between 0 and height %f", c.GroundLevel, c.Height)
	}

	// Validate player spawn is within bounds
	if c.PlayerSpawn.X < 0 || c.PlayerSpawn.X > c.Width {
		return fmt.Errorf("player spawn X %f must be between 0 and width %f", c.PlayerSpawn.X, c.Width)
	}
	if c.PlayerSpawn.Y < 0 || c.PlayerSpawn.Y > c.Height {
		return fmt.Errorf("player spawn Y %f must be between 0 and height %f", c.PlayerSpawn.Y, c.Height)
	}

	// Validate building X positions are within bounds
	buildingXs := map[string]float32{
		"Market":      c.BuildingLayout.MarketX,
		"FuelStation": c.BuildingLayout.FuelStationX,
		"Hospital":    c.BuildingLayout.HospitalX,
		"UpgradeShop": c.BuildingLayout.UpgradeShopX,
		"ItemShop":    c.BuildingLayout.ItemShopX,
	}

	for name, x := range buildingXs {
		// Building extends from x to x+320
		// Reject only if completely off-screen (not partially)
		if x+320 <= 0 || x >= c.Width {
			return fmt.Errorf("%s X %f completely off-screen (world width %f)", name, x, c.Width)
		}
	}

	return nil
}
