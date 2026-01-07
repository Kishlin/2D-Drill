package entities

import "github.com/Kishlin/drill-game/internal/domain/types"

const (
	UpgradeShopWidth  = 320.0
	UpgradeShopHeight = 192.0
)

// Engine Upgrade Shop

type EngineCatalogEntry struct {
	Price  int
	Engine Engine
}

type EngineUpgradeShop struct {
	AABB    types.AABB
	Catalog []EngineCatalogEntry
}

func NewEngineUpgradeShop(x, y float32) *EngineUpgradeShop {
	return &EngineUpgradeShop{
		AABB: types.NewAABB(x, y, UpgradeShopWidth, UpgradeShopHeight),
		Catalog: []EngineCatalogEntry{
			{Price: 100, Engine: NewEngineMk1()},
			{Price: 300, Engine: NewEngineMk2()},
			{Price: 750, Engine: NewEngineMk3()},
			{Price: 1500, Engine: NewEngineMk4()},
			{Price: 5000, Engine: NewEngineMk5()},
		},
	}
}

func (s *EngineUpgradeShop) IsPlayerInRange(player *Player) bool {
	return s.AABB.Intersects(player.AABB)
}

func (s *EngineUpgradeShop) GetNextEngine(currentTier int) *EngineCatalogEntry {
	nextTier := currentTier + 1
	for i := range s.Catalog {
		if s.Catalog[i].Engine.Tier() == nextTier {
			return &s.Catalog[i]
		}
	}
	return nil // Max level reached
}

// Hull Upgrade Shop

type HullCatalogEntry struct {
	Price int
	Hull  Hull
}

type HullUpgradeShop struct {
	AABB    types.AABB
	Catalog []HullCatalogEntry
}

func NewHullUpgradeShop(x, y float32) *HullUpgradeShop {
	return &HullUpgradeShop{
		AABB: types.NewAABB(x, y, UpgradeShopWidth, UpgradeShopHeight),
		Catalog: []HullCatalogEntry{
			{Price: 150, Hull: NewHullMk1()},
			{Price: 400, Hull: NewHullMk2()},
			{Price: 1000, Hull: NewHullMk3()},
			{Price: 2500, Hull: NewHullMk4()},
			{Price: 8000, Hull: NewHullMk5()},
		},
	}
}

func (s *HullUpgradeShop) IsPlayerInRange(player *Player) bool {
	return s.AABB.Intersects(player.AABB)
}

func (s *HullUpgradeShop) GetNextHull(currentTier int) *HullCatalogEntry {
	nextTier := currentTier + 1
	for i := range s.Catalog {
		if s.Catalog[i].Hull.Tier() == nextTier {
			return &s.Catalog[i]
		}
	}
	return nil // Max level reached
}

// FuelTank Upgrade Shop

type FuelTankCatalogEntry struct {
	Price    int
	FuelTank FuelTank
}

type FuelTankUpgradeShop struct {
	AABB    types.AABB
	Catalog []FuelTankCatalogEntry
}

func NewFuelTankUpgradeShop(x, y float32) *FuelTankUpgradeShop {
	return &FuelTankUpgradeShop{
		AABB: types.NewAABB(x, y, UpgradeShopWidth, UpgradeShopHeight),
		Catalog: []FuelTankCatalogEntry{
			{Price: 100, FuelTank: NewFuelTankMk1()},
			{Price: 250, FuelTank: NewFuelTankMk2()},
			{Price: 600, FuelTank: NewFuelTankMk3()},
			{Price: 1500, FuelTank: NewFuelTankMk4()},
			{Price: 4000, FuelTank: NewFuelTankMk5()},
		},
	}
}

func (s *FuelTankUpgradeShop) IsPlayerInRange(player *Player) bool {
	return s.AABB.Intersects(player.AABB)
}

func (s *FuelTankUpgradeShop) GetNextFuelTank(currentTier int) *FuelTankCatalogEntry {
	nextTier := currentTier + 1
	for i := range s.Catalog {
		if s.Catalog[i].FuelTank.Tier() == nextTier {
			return &s.Catalog[i]
		}
	}
	return nil // Max level reached
}

// CargoHold Upgrade Shop

type CargoHoldCatalogEntry struct {
	Price     int
	CargoHold CargoHold
}

type CargoHoldUpgradeShop struct {
	AABB    types.AABB
	Catalog []CargoHoldCatalogEntry
}

func NewCargoHoldUpgradeShop(x, y float32) *CargoHoldUpgradeShop {
	return &CargoHoldUpgradeShop{
		AABB: types.NewAABB(x, y, UpgradeShopWidth, UpgradeShopHeight),
		Catalog: []CargoHoldCatalogEntry{
			{Price: 125, CargoHold: NewCargoHoldMk1()},
			{Price: 350, CargoHold: NewCargoHoldMk2()},
			{Price: 800, CargoHold: NewCargoHoldMk3()},
			{Price: 2000, CargoHold: NewCargoHoldMk4()},
			{Price: 6000, CargoHold: NewCargoHoldMk5()},
		},
	}
}

func (s *CargoHoldUpgradeShop) IsPlayerInRange(player *Player) bool {
	return s.AABB.Intersects(player.AABB)
}

func (s *CargoHoldUpgradeShop) GetNextCargoHold(currentTier int) *CargoHoldCatalogEntry {
	nextTier := currentTier + 1
	for i := range s.Catalog {
		if s.Catalog[i].CargoHold.Tier() == nextTier {
			return &s.Catalog[i]
		}
	}
	return nil // Max level reached
}
