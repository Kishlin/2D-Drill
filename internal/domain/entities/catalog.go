package entities

import "github.com/Kishlin/drill-game/internal/domain/config"

// Upgrade catalog entry types

type EngineCatalogEntry struct {
	Price  int
	Engine Engine
}

type HullCatalogEntry struct {
	Price int
	Hull  Hull
}

type FuelTankCatalogEntry struct {
	Price    int
	FuelTank FuelTank
}

type CargoHoldCatalogEntry struct {
	Price     int
	CargoHold CargoHold
}

type HeatShieldCatalogEntry struct {
	Price      int
	HeatShield HeatShield
}

type DrillCatalogEntry struct {
	Price int
	Drill Drill
}

// UpgradeCatalog contains all upgrade types and tiers
type UpgradeCatalog struct {
	Engines    []EngineCatalogEntry
	Hulls      []HullCatalogEntry
	FuelTanks  []FuelTankCatalogEntry
	CargoHolds []CargoHoldCatalogEntry
	HeatShields []HeatShieldCatalogEntry
	Drills     []DrillCatalogEntry
}

func NewUpgradeCatalogFromConfig(upgradeCfg config.UpgradeConfig) *UpgradeCatalog {
	catalog := &UpgradeCatalog{
		Engines:    make([]EngineCatalogEntry, len(upgradeCfg.Engines)),
		Hulls:      make([]HullCatalogEntry, len(upgradeCfg.Hulls)),
		FuelTanks:  make([]FuelTankCatalogEntry, len(upgradeCfg.FuelTanks)),
		CargoHolds: make([]CargoHoldCatalogEntry, len(upgradeCfg.CargoHolds)),
		HeatShields: make([]HeatShieldCatalogEntry, len(upgradeCfg.HeatShields)),
		Drills:     make([]DrillCatalogEntry, len(upgradeCfg.Drills)),
	}

	for i, tier := range upgradeCfg.Engines {
		catalog.Engines[i] = EngineCatalogEntry{
			Price:  tier.Price,
			Engine: NewEngineFromConfig(i, tier.Name, tier.Stats),
		}
	}

	for i, tier := range upgradeCfg.Hulls {
		catalog.Hulls[i] = HullCatalogEntry{
			Price: tier.Price,
			Hull:  NewHullFromConfig(i, tier.Name, tier.Stats),
		}
	}

	for i, tier := range upgradeCfg.FuelTanks {
		catalog.FuelTanks[i] = FuelTankCatalogEntry{
			Price:    tier.Price,
			FuelTank: NewFuelTankFromConfig(i, tier.Name, tier.Stats),
		}
	}

	for i, tier := range upgradeCfg.CargoHolds {
		catalog.CargoHolds[i] = CargoHoldCatalogEntry{
			Price:     tier.Price,
			CargoHold: NewCargoHoldFromConfig(i, tier.Name, tier.Stats),
		}
	}

	for i, tier := range upgradeCfg.HeatShields {
		catalog.HeatShields[i] = HeatShieldCatalogEntry{
			Price:      tier.Price,
			HeatShield: NewHeatShieldFromConfig(i, tier.Name, tier.Stats),
		}
	}

	for i, tier := range upgradeCfg.Drills {
		catalog.Drills[i] = DrillCatalogEntry{
			Price: tier.Price,
			Drill: NewDrillFromConfig(i, tier.Name, tier.Stats),
		}
	}

	return catalog
}

func (c *UpgradeCatalog) GetEngineCatalogEntry(tier int) *EngineCatalogEntry {
	if tier < 0 || tier >= len(c.Engines) {
		return nil
	}
	return &c.Engines[tier]
}

func (c *UpgradeCatalog) GetHullCatalogEntry(tier int) *HullCatalogEntry {
	if tier < 0 || tier >= len(c.Hulls) {
		return nil
	}
	return &c.Hulls[tier]
}

func (c *UpgradeCatalog) GetFuelTankCatalogEntry(tier int) *FuelTankCatalogEntry {
	if tier < 0 || tier >= len(c.FuelTanks) {
		return nil
	}
	return &c.FuelTanks[tier]
}

func (c *UpgradeCatalog) GetCargoHoldCatalogEntry(tier int) *CargoHoldCatalogEntry {
	if tier < 0 || tier >= len(c.CargoHolds) {
		return nil
	}
	return &c.CargoHolds[tier]
}

func (c *UpgradeCatalog) GetHeatShieldCatalogEntry(tier int) *HeatShieldCatalogEntry {
	if tier < 0 || tier >= len(c.HeatShields) {
		return nil
	}
	return &c.HeatShields[tier]
}

func (c *UpgradeCatalog) GetDrillCatalogEntry(tier int) *DrillCatalogEntry {
	if tier < 0 || tier >= len(c.Drills) {
		return nil
	}
	return &c.Drills[tier]
}

func (c *UpgradeCatalog) GetUpgradePrice(upgradeType UpgradeType, tier int) int {
	switch upgradeType {
	case UpgradeEngine:
		if entry := c.GetEngineCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	case UpgradeHull:
		if entry := c.GetHullCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	case UpgradeFuelTank:
		if entry := c.GetFuelTankCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	case UpgradeCargoHold:
		if entry := c.GetCargoHoldCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	case UpgradeHeatShield:
		if entry := c.GetHeatShieldCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	case UpgradeDrill:
		if entry := c.GetDrillCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	}
	return 0
}

func (c *UpgradeCatalog) GetUpgradeName(upgradeType UpgradeType, tier int) string {
	switch upgradeType {
	case UpgradeEngine:
		if entry := c.GetEngineCatalogEntry(tier); entry != nil {
			return entry.Engine.Name()
		}
	case UpgradeHull:
		if entry := c.GetHullCatalogEntry(tier); entry != nil {
			return entry.Hull.Name()
		}
	case UpgradeFuelTank:
		if entry := c.GetFuelTankCatalogEntry(tier); entry != nil {
			return entry.FuelTank.Name()
		}
	case UpgradeCargoHold:
		if entry := c.GetCargoHoldCatalogEntry(tier); entry != nil {
			return entry.CargoHold.Name()
		}
	case UpgradeHeatShield:
		if entry := c.GetHeatShieldCatalogEntry(tier); entry != nil {
			return entry.HeatShield.Name()
		}
	case UpgradeDrill:
		if entry := c.GetDrillCatalogEntry(tier); entry != nil {
			return entry.Drill.Name()
		}
	}
	return "Unknown"
}

// Item catalog entry type

type ItemCatalogEntry struct {
	ItemType ItemType
	Price    int
}

// ItemCatalog contains all purchasable items
type ItemCatalog struct {
	Items [5]ItemCatalogEntry
}

func NewItemCatalogFromConfig(itemCfg config.ItemConfig) *ItemCatalog {
	return &ItemCatalog{
		Items: [5]ItemCatalogEntry{
			{ItemType: ItemTeleport, Price: itemCfg.Teleport.Price},
			{ItemType: ItemRepair, Price: itemCfg.Repair.Price},
			{ItemType: ItemRefuel, Price: itemCfg.Refuel.Price},
			{ItemType: ItemBomb, Price: itemCfg.Bomb.Price},
			{ItemType: ItemBigBomb, Price: itemCfg.BigBomb.Price},
		},
	}
}

func (c *ItemCatalog) GetItem(index int) *ItemCatalogEntry {
	if index < 0 || index >= len(c.Items) {
		return nil
	}
	return &c.Items[index]
}

// GetPlayerCurrentTier returns the player's current tier for a given upgrade type
func GetPlayerCurrentTier(player *Player, upgradeType UpgradeType) int {
	switch upgradeType {
	case UpgradeEngine:
		return player.Engine.Tier()
	case UpgradeHull:
		return player.Hull.Tier()
	case UpgradeFuelTank:
		return player.FuelTank.Tier()
	case UpgradeCargoHold:
		return player.CargoHold.Tier()
	case UpgradeHeatShield:
		return player.HeatShield.Tier()
	case UpgradeDrill:
		return player.Drill.Tier()
	}
	return 0
}
