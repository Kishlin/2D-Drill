package levels

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
)

// GetTestLevelConfig returns a test configuration with an advanced player
// for easier testing. Use levelNum -1 to load this config.
func GetTestLevelConfig() *config.GameConfig {
	return &config.GameConfig{
		World: config.WorldConfig{
			Width:       3072,
			Height:      64 * 800,
			GroundLevel: 640.0,
			Seed:        42,
			PlayerSpawn: config.PlayerSpawn{
				X: 1536.0,
				Y: 570.0,
			},
			BuildingLayout: config.BuildingLayout{
				HospitalX:    480.0,
				FuelStationX: 850.0,
				MarketX:      1400.0,
				UpgradeShopX: 1850.0,
				ItemShopX:    2220.0,
			},
		},

		Player: config.PlayerConfig{
			StartingMoney: 100000,
			StartingItems: [5]int{
				3,    // Teleport
				5,    // Repair
				5,    // Refuel
				10,   // Bomb
				2000, // Big Bomb
			},
			StartingUpgrades: config.StartingUpgrades{
				Engine:     5,
				Hull:       3,
				FuelTank:   3,
				CargoHold:  3,
				HeatShield: 3,
				Drill:      5,
			},
		},

		Generation: config.GenerationConfig{
			Empty: config.TileDistribution{
				PeakDepth: 0,
				Sigma:     1000,
				MaxWeight: 20.0,
			},
			Dirt: config.TileDistribution{
				PeakDepth: 0,
				Sigma:     500,
				MaxWeight: 100.0,
			},
			DirtHardness: 1.0,
			Ores: []config.OreConfig{
				{
					ID:   "copper",
					Name: "Copper",
					Distribution: config.TileDistribution{
						PeakDepth: -75,
						Sigma:     120,
						MaxWeight: 8.0,
					},
					Value:    25,
					Hardness: 1.2,
					Color:    [4]uint8{184, 115, 51, 255},
				},
				{
					ID:   "iron",
					Name: "Iron",
					Distribution: config.TileDistribution{
						PeakDepth: 70,
						Sigma:     90,
						MaxWeight: 5.0,
					},
					Value:    75,
					Hardness: 1.5,
					Color:    [4]uint8{112, 128, 144, 255},
				},
				{
					ID:   "gold",
					Name: "Gold",
					Distribution: config.TileDistribution{
						PeakDepth: 230,
						Sigma:     80,
						MaxWeight: 3.0,
					},
					Value:    300,
					Hardness: 1.8,
					Color:    [4]uint8{255, 215, 0, 255},
				},
				{
					ID:   "mythril",
					Name: "Mythril",
					Distribution: config.TileDistribution{
						PeakDepth: 360,
						Sigma:     70,
						MaxWeight: 2.2,
					},
					Value:    1500,
					Hardness: 2.1,
					Color:    [4]uint8{0, 191, 255, 255},
				},
				{
					ID:   "platinum",
					Name: "Platinum",
					Distribution: config.TileDistribution{
						PeakDepth: 500,
						Sigma:     80,
						MaxWeight: 1.8,
					},
					Value:    10000,
					Hardness: 2.5,
					Color:    [4]uint8{229, 228, 226, 255},
				},
				{
					ID:   "diamond",
					Name: "Diamond",
					Distribution: config.TileDistribution{
						PeakDepth: 600,
						Sigma:     180,
						MaxWeight: 0.15,
					},
					Value:    30000,
					Hardness: 3.0,
					Color:    [4]uint8{185, 242, 255, 255},
				},
			},
			Hazards: []config.HazardConfig{
				{
					ID:   "rock",
					Name: "Rock",
					Distribution: config.TileDistribution{
						PeakDepth: 650.0,
						Sigma:     200.0,
						MaxWeight: 15.0,
					},
					Drillable:     false,
					FixedDuration: 0,
					OnDrillEffect: config.HazardEffectConfig{
						Type: config.HazardEffectNone,
					},
					OnHitDamage: 0,
					Color:       [4]uint8{80, 80, 80, 255},
				},
				{
					ID:   "lava",
					Name: "Lava",
					Distribution: config.TileDistribution{
						PeakDepth: 750.0,
						Sigma:     150.0,
						MaxWeight: 12.0,
					},
					Drillable:     true,
					FixedDuration: 0.3,
					OnDrillEffect: config.HazardEffectConfig{
						Type:               config.HazardEffectHeatDamage,
						BaseDamage:         100.0,
						MaxHeatResistance:  320.0,
						MaxDamageReduction: 0.5,
					},
					OnHitDamage: 0,
					Color:       [4]uint8{255, 100, 0, 255},
				},
			},
		},

		Upgrades: config.UpgradeConfig{
			Engines: []config.UpgradeTier[config.EngineStats]{
				{Price: 0, Stats: config.EngineStats{MaxSpeed: 450.0, Acceleration: 2500.0, FlyAcceleration: 2500.0, MaxUpwardSpeed: -600.0}},
				{Price: 500, Stats: config.EngineStats{MaxSpeed: 475.0, Acceleration: 2667.0, FlyAcceleration: 2667.0, MaxUpwardSpeed: -635.0}},
				{Price: 2000, Stats: config.EngineStats{MaxSpeed: 500.0, Acceleration: 2833.0, FlyAcceleration: 2833.0, MaxUpwardSpeed: -670.0}},
				{Price: 8000, Stats: config.EngineStats{MaxSpeed: 525.0, Acceleration: 3000.0, FlyAcceleration: 3000.0, MaxUpwardSpeed: -705.0}},
				{Price: 32000, Stats: config.EngineStats{MaxSpeed: 562.0, Acceleration: 3250.0, FlyAcceleration: 3250.0, MaxUpwardSpeed: -740.0}},
				{Price: 128000, Stats: config.EngineStats{MaxSpeed: 600.0, Acceleration: 3500.0, FlyAcceleration: 3500.0, MaxUpwardSpeed: -775.0}},
			},
			Hulls: []config.UpgradeTier[config.HullStats]{
				{Price: 0, Stats: config.HullStats{MaxHP: 10.0}},
				{Price: 750, Stats: config.HullStats{MaxHP: 15.0}},
				{Price: 3000, Stats: config.HullStats{MaxHP: 20.0}},
				{Price: 12000, Stats: config.HullStats{MaxHP: 30.0}},
				{Price: 48000, Stats: config.HullStats{MaxHP: 45.0}},
				{Price: 192000, Stats: config.HullStats{MaxHP: 75.0}},
			},
			FuelTanks: []config.UpgradeTier[config.FuelTankStats]{
				{Price: 0, Stats: config.FuelTankStats{Capacity: 10.0}},
				{Price: 400, Stats: config.FuelTankStats{Capacity: 15.0}},
				{Price: 1600, Stats: config.FuelTankStats{Capacity: 22.0}},
				{Price: 6400, Stats: config.FuelTankStats{Capacity: 32.0}},
				{Price: 25600, Stats: config.FuelTankStats{Capacity: 45.0}},
				{Price: 102400, Stats: config.FuelTankStats{Capacity: 65.0}},
			},
			CargoHolds: []config.UpgradeTier[config.CargoHoldStats]{
				{Price: 0, Stats: config.CargoHoldStats{Capacity: 10}},
				{Price: 600, Stats: config.CargoHoldStats{Capacity: 14}},
				{Price: 2400, Stats: config.CargoHoldStats{Capacity: 18}},
				{Price: 9600, Stats: config.CargoHoldStats{Capacity: 24}},
				{Price: 38400, Stats: config.CargoHoldStats{Capacity: 31}},
				{Price: 153600, Stats: config.CargoHoldStats{Capacity: 40}},
			},
			HeatShields: []config.UpgradeTier[config.HeatShieldStats]{
				{Price: 0, Stats: config.HeatShieldStats{HeatResistance: 50.0}},
				{Price: 800, Stats: config.HeatShieldStats{HeatResistance: 90.0}},
				{Price: 3200, Stats: config.HeatShieldStats{HeatResistance: 140.0}},
				{Price: 12800, Stats: config.HeatShieldStats{HeatResistance: 190.0}},
				{Price: 51200, Stats: config.HeatShieldStats{HeatResistance: 250.0}},
				{Price: 204800, Stats: config.HeatShieldStats{HeatResistance: 320.0}},
			},
			Drills: []config.UpgradeTier[config.DrillStats]{
				{Price: 0, Stats: config.DrillStats{SpeedAtSurface: 1.0, SpeedAtMaxDepth: 1.0}},
				{Price: 1000, Stats: config.DrillStats{SpeedAtSurface: 1.1, SpeedAtMaxDepth: 2.0}},
				{Price: 4000, Stats: config.DrillStats{SpeedAtSurface: 1.2, SpeedAtMaxDepth: 3.0}},
				{Price: 16000, Stats: config.DrillStats{SpeedAtSurface: 1.3, SpeedAtMaxDepth: 4.0}},
				{Price: 64000, Stats: config.DrillStats{SpeedAtSurface: 1.4, SpeedAtMaxDepth: 5.0}},
				{Price: 256000, Stats: config.DrillStats{SpeedAtSurface: 1.5, SpeedAtMaxDepth: 6.0}},
			},
		},

		Items: config.ItemConfig{
			Teleport: config.ItemEntry{Price: 1000, Radius: 0, Damage: 0},
			Repair:   config.ItemEntry{Price: 2500, Radius: 0, Damage: 0},
			Refuel:   config.ItemEntry{Price: 500, Radius: 0, Damage: 0},
			Bomb:     config.ItemEntry{Price: 3000, Radius: 2, Damage: 10.0},
			BigBomb:  config.ItemEntry{Price: 10000, Radius: 4, Damage: 25.0},
		},

		Level: config.LevelConfig{
			Number: -1,
			Name:   "Test Level",
			BossRoom: config.BossRoomConfig{
				BossType:    "test_boss",
				FloorType:   config.FloorConcrete,
				RoomHeight:  680.0,
				FloorHeight: 6.0,
			},
		},

		Drilling: config.DrillingConfig{
			MinDrillingDuration:   1.0,
			MaxDrillingDuration:   24.0,
			FloorDrillingDuration: 0.5,
		},

		Heat: config.HeatConfig{
			BaseTemperature: 15.0,
			MaxTemperature:  350.0,
			DamageBaseDPS:   0.5,
			DamageDivisor:   10.0,
			DamageExponent:  1.5,
		},

		Fuel: config.FuelSystemConfig{
			ConsumptionMoving: 10.0 / 30.0,  // 0.33333 L/s
			ConsumptionIdle:   10.0 / 120.0, // 0.08333 L/s
		},
	}
}
