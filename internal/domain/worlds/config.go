package worlds

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/world" // Must be imported for BaseTileConfig
)

// DistributionParams holds Gaussian distribution parameters for ore and hazard generation
type DistributionParams struct {
	PeakDepth float32 // Tile Y coordinate where spawn is most common
	Sigma     float32 // Standard deviation (spread of distribution)
	MaxWeight float32 // Weight at peak depth (relative spawn chance)
}

// OreDefinition defines a single ore type within a world's configuration
type OreDefinition struct {
	OreType     entities.OreType
	Name        string
	Value       int // Market sell price
	Hardness    float32
	Distribution DistributionParams
}

// HazardDefinition defines a single hazard type within a world's configuration
type HazardDefinition struct {
	HazardType   entities.HazardType
	Name         string
	Hardness     float32
	Distribution DistributionParams
}

// OreConfig contains all ore definitions for a world
type OreConfig struct {
	Ores []OreDefinition
}

// HazardConfig contains all hazard definitions for a world
type HazardConfig struct {
	Hazards []HazardDefinition
}

// BaseTileConfig is an alias to world.BaseTileConfig for convenience
// The actual definition is in the world package
type BaseTileConfig = world.BaseTileConfig

// EngineCatalogDefinition defines a single engine tier for config
type EngineCatalogDefinition struct {
	Price           int
	MaxSpeed        float32
	Acceleration    float32
	FlyAcceleration float32
	MaxUpwardSpeed  float32
}

// HullCatalogDefinition defines a single hull tier for config
type HullCatalogDefinition struct {
	Price int
	MaxHP float32
}

// FuelTankCatalogDefinition defines a single fuel tank tier for config
type FuelTankCatalogDefinition struct {
	Price    int
	Capacity float32
}

// CargoHoldCatalogDefinition defines a single cargo hold tier for config
type CargoHoldCatalogDefinition struct {
	Price    int
	Capacity int
}

// HeatShieldCatalogDefinition defines a single heat shield tier for config
type HeatShieldCatalogDefinition struct {
	Price           int
	HeatResistance  float32
}

// DrillCatalogDefinition defines a single drill tier for config
type DrillCatalogDefinition struct {
	Price      int
	DrillSpeed float32
}

// UpgradeConfig contains all upgrade catalogs for a world
type UpgradeConfig struct {
	Engine    []EngineCatalogDefinition
	Hull      []HullCatalogDefinition
	FuelTank  []FuelTankCatalogDefinition
	CargoHold []CargoHoldCatalogDefinition
	HeatShield []HeatShieldCatalogDefinition
	Drill     []DrillCatalogDefinition
}

// ItemDefinition defines a single item for purchase in the item shop
type ItemDefinition struct {
	ItemType entities.ItemType
	Name     string
	Price    int
}

// ItemConfig contains all item definitions for a world
type ItemConfig struct {
	Items []ItemDefinition
}

// PlayerInitialStatus defines the initial state of the player for a world
// This enables different scenarios: basic start ($0, Base upgrades), endgame ($100k, Mk5), etc.
type PlayerInitialStatus struct {
	Money          int    // Starting money (default: 0)
	EngineTier     int    // Starting engine tier (0-5: Base-Mk5, default: 0)
	HullTier       int    // Starting hull tier (0-5: Base-Mk5, default: 0)
	FuelTankTier   int    // Starting fuel tank tier (0-5: Base-Mk5, default: 0)
	CargoHoldTier  int    // Starting cargo hold tier (0-5: Base-Mk5, default: 0)
	HeatShieldTier int    // Starting heat shield tier (0-5: Base-Mk5, default: 0)
	DrillTier      int    // Starting drill tier (0-5: Base-Mk5, default: 0)
	ItemInventory  [5]int // Starting item counts: [Teleport, Repair, Refuel, Bomb, BigBomb] (default: all 0)
}

// WorldGameConfig represents a complete world configuration with all game content
type WorldGameConfig struct {
	World         *world.WorldConfig    // World dimensions, seed, spawn, buildings
	Ores          OreConfig             // Ore types and their properties
	Hazards       HazardConfig          // Hazard types and their properties
	BaseTiles     BaseTileConfig        // Empty/dirt weight formulas
	Upgrades      UpgradeConfig         // All upgrade tiers for 6 upgrade types
	Items         ItemConfig            // All purchasable items
	PlayerInitial PlayerInitialStatus   // Initial player status (money, upgrades, items)
}
