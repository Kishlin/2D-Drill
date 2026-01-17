package entities

import "github.com/Kishlin/drill-game/internal/domain/types"

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

func NewUpgradeShop(x, y float32) *UpgradeShop {
	return &UpgradeShop{
		AABB: types.NewAABB(x, y, UpgradeShopWidth, UpgradeShopHeight),
		EngineCatalog: []EngineCatalogEntry{
			{Price: 0, Engine: NewEngineBase()},
			{Price: 100, Engine: NewEngineMk1()},
			{Price: 300, Engine: NewEngineMk2()},
			{Price: 750, Engine: NewEngineMk3()},
			{Price: 1500, Engine: NewEngineMk4()},
			{Price: 5000, Engine: NewEngineMk5()},
		},
		HullCatalog: []HullCatalogEntry{
			{Price: 0, Hull: NewHullBase()},
			{Price: 150, Hull: NewHullMk1()},
			{Price: 400, Hull: NewHullMk2()},
			{Price: 1000, Hull: NewHullMk3()},
			{Price: 2500, Hull: NewHullMk4()},
			{Price: 8000, Hull: NewHullMk5()},
		},
		FuelTankCatalog: []FuelTankCatalogEntry{
			{Price: 0, FuelTank: NewFuelTankBase()},
			{Price: 100, FuelTank: NewFuelTankMk1()},
			{Price: 250, FuelTank: NewFuelTankMk2()},
			{Price: 600, FuelTank: NewFuelTankMk3()},
			{Price: 1500, FuelTank: NewFuelTankMk4()},
			{Price: 4000, FuelTank: NewFuelTankMk5()},
		},
		CargoCatalog: []CargoHoldCatalogEntry{
			{Price: 0, CargoHold: NewCargoHoldBase()},
			{Price: 125, CargoHold: NewCargoHoldMk1()},
			{Price: 350, CargoHold: NewCargoHoldMk2()},
			{Price: 800, CargoHold: NewCargoHoldMk3()},
			{Price: 2000, CargoHold: NewCargoHoldMk4()},
			{Price: 6000, CargoHold: NewCargoHoldMk5()},
		},
		HeatCatalog: []HeatShieldCatalogEntry{
			{Price: 0, HeatShield: NewHeatShieldBase()},
			{Price: 200, HeatShield: NewHeatShieldMk1()},
			{Price: 500, HeatShield: NewHeatShieldMk2()},
			{Price: 1200, HeatShield: NewHeatShieldMk3()},
			{Price: 3000, HeatShield: NewHeatShieldMk4()},
			{Price: 7500, HeatShield: NewHeatShieldMk5()},
		},
		DrillCatalog: []DrillCatalogEntry{
			{Price: 0, Drill: NewDrillBase()},
			{Price: 125, Drill: NewDrillMk1()},
			{Price: 350, Drill: NewDrillMk2()},
			{Price: 875, Drill: NewDrillMk3()},
			{Price: 2000, Drill: NewDrillMk4()},
			{Price: 6500, Drill: NewDrillMk5()},
		},
	}
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
