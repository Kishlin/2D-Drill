package entities

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

const (
	UpgradeShopWidth  = 320.0
	UpgradeShopHeight = 192.0
)

// Catalog entry types

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

type UpgradeShop struct {
	AABB            types.AABB
	EngineCatalog   []EngineCatalogEntry
	HullCatalog     []HullCatalogEntry
	FuelTankCatalog []FuelTankCatalogEntry
	CargoCatalog    []CargoHoldCatalogEntry
	HeatCatalog     []HeatShieldCatalogEntry
	DrillCatalog    []DrillCatalogEntry
}

func NewUpgradeShopFromConfig(x, y float32, upgradeCfg config.UpgradeConfig) *UpgradeShop {
	shop := &UpgradeShop{
		AABB:            types.NewAABB(x, y, UpgradeShopWidth, UpgradeShopHeight),
		EngineCatalog:   make([]EngineCatalogEntry, len(upgradeCfg.Engines)),
		HullCatalog:     make([]HullCatalogEntry, len(upgradeCfg.Hulls)),
		FuelTankCatalog: make([]FuelTankCatalogEntry, len(upgradeCfg.FuelTanks)),
		CargoCatalog:    make([]CargoHoldCatalogEntry, len(upgradeCfg.CargoHolds)),
		HeatCatalog:     make([]HeatShieldCatalogEntry, len(upgradeCfg.HeatShields)),
		DrillCatalog:    make([]DrillCatalogEntry, len(upgradeCfg.Drills)),
	}

	for i, tier := range upgradeCfg.Engines {
		shop.EngineCatalog[i] = EngineCatalogEntry{
			Price:  tier.Price,
			Engine: NewEngineFromConfig(i, tier.Name, tier.Stats),
		}
	}

	for i, tier := range upgradeCfg.Hulls {
		shop.HullCatalog[i] = HullCatalogEntry{
			Price: tier.Price,
			Hull:  NewHullFromConfig(i, tier.Name, tier.Stats),
		}
	}

	for i, tier := range upgradeCfg.FuelTanks {
		shop.FuelTankCatalog[i] = FuelTankCatalogEntry{
			Price:    tier.Price,
			FuelTank: NewFuelTankFromConfig(i, tier.Name, tier.Stats),
		}
	}

	for i, tier := range upgradeCfg.CargoHolds {
		shop.CargoCatalog[i] = CargoHoldCatalogEntry{
			Price:     tier.Price,
			CargoHold: NewCargoHoldFromConfig(i, tier.Name, tier.Stats),
		}
	}

	for i, tier := range upgradeCfg.HeatShields {
		shop.HeatCatalog[i] = HeatShieldCatalogEntry{
			Price:      tier.Price,
			HeatShield: NewHeatShieldFromConfig(i, tier.Name, tier.Stats),
		}
	}

	for i, tier := range upgradeCfg.Drills {
		shop.DrillCatalog[i] = DrillCatalogEntry{
			Price: tier.Price,
			Drill: NewDrillFromConfig(i, tier.Name, tier.Stats),
		}
	}

	return shop
}

func (s *UpgradeShop) IsPlayerInRange(player *Player) bool {
	return s.AABB.Intersects(player.AABB)
}

func (s *UpgradeShop) GetEngineCatalogEntry(tier int) *EngineCatalogEntry {
	if tier < 0 || tier >= len(s.EngineCatalog) {
		return nil
	}
	return &s.EngineCatalog[tier]
}

func (s *UpgradeShop) GetHullCatalogEntry(tier int) *HullCatalogEntry {
	if tier < 0 || tier >= len(s.HullCatalog) {
		return nil
	}
	return &s.HullCatalog[tier]
}

func (s *UpgradeShop) GetFuelTankCatalogEntry(tier int) *FuelTankCatalogEntry {
	if tier < 0 || tier >= len(s.FuelTankCatalog) {
		return nil
	}
	return &s.FuelTankCatalog[tier]
}

func (s *UpgradeShop) GetCargoCatalogEntry(tier int) *CargoHoldCatalogEntry {
	if tier < 0 || tier >= len(s.CargoCatalog) {
		return nil
	}
	return &s.CargoCatalog[tier]
}

func (s *UpgradeShop) GetHeatCatalogEntry(tier int) *HeatShieldCatalogEntry {
	if tier < 0 || tier >= len(s.HeatCatalog) {
		return nil
	}
	return &s.HeatCatalog[tier]
}

func (s *UpgradeShop) GetDrillCatalogEntry(tier int) *DrillCatalogEntry {
	if tier < 0 || tier >= len(s.DrillCatalog) {
		return nil
	}
	return &s.DrillCatalog[tier]
}

func (s *UpgradeShop) GetUpgradePrice(upgradeType UpgradeType, tier int) int {
	switch upgradeType {
	case UpgradeEngine:
		if entry := s.GetEngineCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	case UpgradeHull:
		if entry := s.GetHullCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	case UpgradeFuelTank:
		if entry := s.GetFuelTankCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	case UpgradeCargoHold:
		if entry := s.GetCargoCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	case UpgradeHeatShield:
		if entry := s.GetHeatCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	case UpgradeDrill:
		if entry := s.GetDrillCatalogEntry(tier); entry != nil {
			return entry.Price
		}
	}
	return 0
}

func (s *UpgradeShop) GetUpgradeName(upgradeType UpgradeType, tier int) string {
	switch upgradeType {
	case UpgradeEngine:
		if entry := s.GetEngineCatalogEntry(tier); entry != nil {
			return entry.Engine.Name()
		}
	case UpgradeHull:
		if entry := s.GetHullCatalogEntry(tier); entry != nil {
			return entry.Hull.Name()
		}
	case UpgradeFuelTank:
		if entry := s.GetFuelTankCatalogEntry(tier); entry != nil {
			return entry.FuelTank.Name()
		}
	case UpgradeCargoHold:
		if entry := s.GetCargoCatalogEntry(tier); entry != nil {
			return entry.CargoHold.Name()
		}
	case UpgradeHeatShield:
		if entry := s.GetHeatCatalogEntry(tier); entry != nil {
			return entry.HeatShield.Name()
		}
	case UpgradeDrill:
		if entry := s.GetDrillCatalogEntry(tier); entry != nil {
			return entry.Drill.Name()
		}
	}
	return "Unknown"
}

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
