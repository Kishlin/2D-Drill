package upgrades

import "github.com/Kishlin/drill-game/internal/domain/config"

type CatalogEntry struct {
	Price   int
	Upgrade Upgrade
}

type Catalog struct {
	entries [TypeCount][]CatalogEntry
}

func NewCatalogFromConfig(cfg config.UpgradeConfig) *Catalog {
	c := &Catalog{}

	c.entries[TypeEngine] = make([]CatalogEntry, len(cfg.Engines))
	for i, tier := range cfg.Engines {
		c.entries[TypeEngine][i] = CatalogEntry{
			Price:   tier.Price,
			Upgrade: NewEngineFromConfig(i, tier.Name, tier.Stats),
		}
	}

	c.entries[TypeHull] = make([]CatalogEntry, len(cfg.Hulls))
	for i, tier := range cfg.Hulls {
		c.entries[TypeHull][i] = CatalogEntry{
			Price:   tier.Price,
			Upgrade: NewHullFromConfig(i, tier.Name, tier.Stats),
		}
	}

	c.entries[TypeFuelTank] = make([]CatalogEntry, len(cfg.FuelTanks))
	for i, tier := range cfg.FuelTanks {
		c.entries[TypeFuelTank][i] = CatalogEntry{
			Price:   tier.Price,
			Upgrade: NewFuelTankFromConfig(i, tier.Name, tier.Stats),
		}
	}

	c.entries[TypeCargoHold] = make([]CatalogEntry, len(cfg.CargoHolds))
	for i, tier := range cfg.CargoHolds {
		c.entries[TypeCargoHold][i] = CatalogEntry{
			Price:   tier.Price,
			Upgrade: NewCargoHoldFromConfig(i, tier.Name, tier.Stats),
		}
	}

	c.entries[TypeHeatShield] = make([]CatalogEntry, len(cfg.HeatShields))
	for i, tier := range cfg.HeatShields {
		c.entries[TypeHeatShield][i] = CatalogEntry{
			Price:   tier.Price,
			Upgrade: NewHeatShieldFromConfig(i, tier.Name, tier.Stats),
		}
	}

	c.entries[TypeDrill] = make([]CatalogEntry, len(cfg.Drills))
	for i, tier := range cfg.Drills {
		c.entries[TypeDrill][i] = CatalogEntry{
			Price:   tier.Price,
			Upgrade: NewDrillFromConfig(i, tier.Name, tier.Stats),
		}
	}

	return c
}

func (c *Catalog) GetEntry(t UpgradeType, tier int) *CatalogEntry {
	if t < 0 || t >= TypeCount {
		return nil
	}
	entries := c.entries[t]
	if tier < 0 || tier >= len(entries) {
		return nil
	}
	return &entries[tier]
}

func (c *Catalog) GetPrice(t UpgradeType, tier int) int {
	if entry := c.GetEntry(t, tier); entry != nil {
		return entry.Price
	}
	return 0
}

func (c *Catalog) GetName(t UpgradeType, tier int) string {
	if entry := c.GetEntry(t, tier); entry != nil {
		return entry.Upgrade.Name()
	}
	return "Unknown"
}

func (c *Catalog) TierCount(t UpgradeType) int {
	if t < 0 || t >= TypeCount {
		return 0
	}
	return len(c.entries[t])
}
