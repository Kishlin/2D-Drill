package config

import "fmt"

type LevelConfig struct {
	Number   int
	Name     string
	BossRoom BossRoomConfig
}

type GameConfig struct {
	World      WorldConfig
	Player     PlayerConfig
	Generation GenerationConfig
	Upgrades   UpgradeConfig
	Items      ItemConfig
	Level      LevelConfig
}

func (c *GameConfig) Validate() error {
	if err := c.World.Validate(); err != nil {
		return fmt.Errorf("world config: %w", err)
	}

	// Validate player starting upgrades don't exceed available tiers
	if c.Player.StartingUpgrades.Engine >= len(c.Upgrades.Engines) {
		return fmt.Errorf("starting engine tier %d exceeds available tiers %d",
			c.Player.StartingUpgrades.Engine, len(c.Upgrades.Engines))
	}
	if c.Player.StartingUpgrades.Hull >= len(c.Upgrades.Hulls) {
		return fmt.Errorf("starting hull tier %d exceeds available tiers %d",
			c.Player.StartingUpgrades.Hull, len(c.Upgrades.Hulls))
	}
	if c.Player.StartingUpgrades.FuelTank >= len(c.Upgrades.FuelTanks) {
		return fmt.Errorf("starting fuel tank tier %d exceeds available tiers %d",
			c.Player.StartingUpgrades.FuelTank, len(c.Upgrades.FuelTanks))
	}
	if c.Player.StartingUpgrades.CargoHold >= len(c.Upgrades.CargoHolds) {
		return fmt.Errorf("starting cargo hold tier %d exceeds available tiers %d",
			c.Player.StartingUpgrades.CargoHold, len(c.Upgrades.CargoHolds))
	}
	if c.Player.StartingUpgrades.HeatShield >= len(c.Upgrades.HeatShields) {
		return fmt.Errorf("starting heat shield tier %d exceeds available tiers %d",
			c.Player.StartingUpgrades.HeatShield, len(c.Upgrades.HeatShields))
	}
	if c.Player.StartingUpgrades.Drill >= len(c.Upgrades.Drills) {
		return fmt.Errorf("starting drill tier %d exceeds available tiers %d",
			c.Player.StartingUpgrades.Drill, len(c.Upgrades.Drills))
	}

	// Validate generation config has at least one ore and hazard
	if len(c.Generation.Ores) == 0 {
		return fmt.Errorf("generation config must have at least one ore type")
	}

	// Validate ore IDs are unique
	oreIDs := make(map[string]bool)
	for _, ore := range c.Generation.Ores {
		if oreIDs[ore.ID] {
			return fmt.Errorf("duplicate ore ID: %s", ore.ID)
		}
		oreIDs[ore.ID] = true
	}

	// Validate hazard IDs are unique
	hazardIDs := make(map[string]bool)
	for _, hazard := range c.Generation.Hazards {
		if hazardIDs[hazard.ID] {
			return fmt.Errorf("duplicate hazard ID: %s", hazard.ID)
		}
		hazardIDs[hazard.ID] = true
	}

	// Validate boss room config (required for all levels)
	if c.Level.BossRoom.BossType == "" {
		return fmt.Errorf("boss room must have a boss type")
	}
	if c.Level.BossRoom.RoomHeight <= 0 {
		return fmt.Errorf("boss room height must be positive")
	}
	if c.Level.BossRoom.FloorHeight < 1 {
		return fmt.Errorf("boss room floor height must be at least 1 tile")
	}

	return nil
}
