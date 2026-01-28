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
	Drilling   DrillingConfig
	Heat       HeatConfig
	Fuel       FuelSystemConfig
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

	// Validate hazard IDs are unique and drillable hazards have FixedDuration
	hazardIDs := make(map[string]bool)
	for _, hazard := range c.Generation.Hazards {
		if hazardIDs[hazard.ID] {
			return fmt.Errorf("duplicate hazard ID: %s", hazard.ID)
		}
		hazardIDs[hazard.ID] = true

		if hazard.Drillable && hazard.FixedDuration <= 0 {
			return fmt.Errorf("drillable hazard %s must have positive FixedDuration", hazard.ID)
		}
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

	// Validate drilling config
	if c.Drilling.MinDrillingDuration <= 0 {
		return fmt.Errorf("drilling min duration must be positive")
	}
	if c.Drilling.MaxDrillingDuration <= 0 {
		return fmt.Errorf("drilling max duration must be positive")
	}
	if c.Drilling.FloorDrillingDuration <= 0 {
		return fmt.Errorf("drilling floor duration must be positive")
	}

	// Validate heat config
	if c.Heat.BaseTemperature < 0 {
		return fmt.Errorf("heat base temperature must be non-negative")
	}
	if c.Heat.MaxTemperature <= c.Heat.BaseTemperature {
		return fmt.Errorf("heat max temperature must exceed base temperature")
	}
	if c.Heat.DamageBaseDPS <= 0 || c.Heat.DamageDivisor <= 0 || c.Heat.DamageExponent <= 0 {
		return fmt.Errorf("heat damage parameters must be positive")
	}

	// Validate fuel config
	if c.Fuel.ConsumptionMoving <= 0 {
		return fmt.Errorf("fuel consumption moving rate must be positive")
	}
	if c.Fuel.ConsumptionIdle <= 0 {
		return fmt.Errorf("fuel consumption idle rate must be positive")
	}

	return nil
}
