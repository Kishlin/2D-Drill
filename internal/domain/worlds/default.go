package worlds

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// NewDefaultWorld returns the default world configuration with the current game balance
func NewDefaultWorld() *WorldGameConfig {
	return &WorldGameConfig{
		World: &world.WorldConfig{
			Width:       3072.0,
			Height:      64 * 800,          // 51200 pixels
			GroundLevel: 640.0,
			Seed:        42,
			PlayerSpawn: world.PlayerSpawn{
				X: 1536.0, // Center of world
				Y: 570.0,  // Just above ground
			},
			BuildingLayout: world.BuildingLayout{
				HospitalX:    480.0,
				FuelStationX: 850.0,
				MarketX:      1400.0,
				UpgradeShopX: 1850.0,
				ItemShopX:    2220.0,
			},
		},
		Ores: OreConfig{
			Ores: []OreDefinition{
				{
					OreType:  entities.OreCopper,
					Name:     "Copper",
					Value:    25,
					Hardness: 1.2,
					Distribution: DistributionParams{
						PeakDepth: -75,
						Sigma:     120,
						MaxWeight: 8.0,
					},
				},
				{
					OreType:  entities.OreIron,
					Name:     "Iron",
					Value:    75,
					Hardness: 1.5,
					Distribution: DistributionParams{
						PeakDepth: 70,
						Sigma:     90,
						MaxWeight: 5.0,
					},
				},
				{
					OreType:  entities.OreGold,
					Name:     "Gold",
					Value:    300,
					Hardness: 1.8,
					Distribution: DistributionParams{
						PeakDepth: 230,
						Sigma:     80,
						MaxWeight: 3.0,
					},
				},
				{
					OreType:  entities.OreMythril,
					Name:     "Mythril",
					Value:    1500,
					Hardness: 2.1,
					Distribution: DistributionParams{
						PeakDepth: 360,
						Sigma:     70,
						MaxWeight: 2.2,
					},
				},
				{
					OreType:  entities.OrePlatinum,
					Name:     "Platinum",
					Value:    10000,
					Hardness: 2.5,
					Distribution: DistributionParams{
						PeakDepth: 500,
						Sigma:     80,
						MaxWeight: 1.8,
					},
				},
				{
					OreType:  entities.OreDiamond,
					Name:     "Diamond",
					Value:    30000,
					Hardness: 3.0,
					Distribution: DistributionParams{
						PeakDepth: 600,
						Sigma:     180,
						MaxWeight: 0.15,
					},
				},
			},
		},
		Hazards: HazardConfig{
			Hazards: []HazardDefinition{
				{
					HazardType: entities.HazardRock,
					Name:       "Rock",
					Hardness:   0, // Rock is not drillable
					Distribution: DistributionParams{
						PeakDepth: 650.0, // ~80% depth
						Sigma:     200.0, // Wide spread to start at ~40%
						MaxWeight: 15.0,  // High to dominate at depth
					},
				},
				{
					HazardType: entities.HazardLava,
					Name:       "Lava",
					Hardness:   0.3,
					Distribution: DistributionParams{
						PeakDepth: 750.0, // ~85% depth
						Sigma:     150.0, // Narrower, starts ~60-65%
						MaxWeight: 12.0,  // High to dominate at depth
					},
				},
			},
		},
		BaseTiles: BaseTileConfig{
			EmptyBaseWeight:    8.0,
			EmptyDepthReduction: 7.5,
			DirtBaseWeight:     20.0,
			DirtDepthReduction:  18.0,
			DirtHardness:       1.0,
		},
		Upgrades: UpgradeConfig{
			Engine: []EngineCatalogDefinition{
				{
					Price:           0,
					MaxSpeed:        450.0,
					Acceleration:    2500.0,
					FlyAcceleration: 2500.0,
					MaxUpwardSpeed:  -600.0,
				},
				{
					Price:           100,
					MaxSpeed:        475.0,
					Acceleration:    2667.0,
					FlyAcceleration: 2667.0,
					MaxUpwardSpeed:  -635.0,
				},
				{
					Price:           300,
					MaxSpeed:        500.0,
					Acceleration:    2833.0,
					FlyAcceleration: 2833.0,
					MaxUpwardSpeed:  -670.0,
				},
				{
					Price:           750,
					MaxSpeed:        525.0,
					Acceleration:    3000.0,
					FlyAcceleration: 3000.0,
					MaxUpwardSpeed:  -705.0,
				},
				{
					Price:           1500,
					MaxSpeed:        562.0,
					Acceleration:    3250.0,
					FlyAcceleration: 3250.0,
					MaxUpwardSpeed:  -740.0,
				},
				{
					Price:           5000,
					MaxSpeed:        600.0,
					Acceleration:    3500.0,
					FlyAcceleration: 3500.0,
					MaxUpwardSpeed:  -775.0,
				},
			},
			Hull: []HullCatalogDefinition{
				{Price: 0, MaxHP: 10.0},
				{Price: 150, MaxHP: 15.0},
				{Price: 400, MaxHP: 20.0},
				{Price: 1000, MaxHP: 30.0},
				{Price: 2500, MaxHP: 45.0},
				{Price: 8000, MaxHP: 75.0},
			},
			FuelTank: []FuelTankCatalogDefinition{
				{Price: 0, Capacity: 10.0},
				{Price: 100, Capacity: 15.0},
				{Price: 250, Capacity: 22.0},
				{Price: 600, Capacity: 32.0},
				{Price: 1500, Capacity: 45.0},
				{Price: 4000, Capacity: 65.0},
			},
			CargoHold: []CargoHoldCatalogDefinition{
				{Price: 0, Capacity: 10},
				{Price: 125, Capacity: 14},
				{Price: 350, Capacity: 18},
				{Price: 800, Capacity: 24},
				{Price: 2000, Capacity: 31},
				{Price: 6000, Capacity: 40},
			},
			HeatShield: []HeatShieldCatalogDefinition{
				{Price: 0, HeatResistance: 50.0},
				{Price: 200, HeatResistance: 90.0},
				{Price: 500, HeatResistance: 140.0},
				{Price: 1200, HeatResistance: 190.0},
				{Price: 3000, HeatResistance: 250.0},
				{Price: 7500, HeatResistance: 320.0},
			},
			Drill: []DrillCatalogDefinition{
				{Price: 0, DrillSpeed: 1.0},
				{Price: 125, DrillSpeed: 2.0},
				{Price: 350, DrillSpeed: 3.0},
				{Price: 875, DrillSpeed: 4.0},
				{Price: 2000, DrillSpeed: 5.0},
				{Price: 6500, DrillSpeed: 6.0},
			},
		},
		Items: ItemConfig{
			Items: []ItemDefinition{
				{ItemType: entities.ItemTeleport, Name: "Teleport", Price: 500},
				{ItemType: entities.ItemRepair, Name: "Repair Kit", Price: 200},
				{ItemType: entities.ItemRefuel, Name: "Fuel Can", Price: 100},
				{ItemType: entities.ItemBomb, Name: "Bomb", Price: 300},
				{ItemType: entities.ItemBigBomb, Name: "Big Bomb", Price: 800},
			},
		},
		PlayerInitial: PlayerInitialStatus{
			Money:          0,
			EngineTier:     0,     // Base
			HullTier:       0,     // Base
			FuelTankTier:   0,     // Base
			CargoHoldTier:  0,     // Base
			HeatShieldTier: 0,     // Base
			DrillTier:      0,     // Base
			ItemInventory:  [5]int{0, 0, 0, 0, 0}, // No starting items
		},
	}
}
