package ui

import (
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type UpgradeShopUI struct {
	catalog *entities.UpgradeCatalog
	state   *UpgradeShopState
}

func NewUpgradeShopUI(catalog *entities.UpgradeCatalog) *UpgradeShopUI {
	return &UpgradeShopUI{
		catalog: catalog,
		state:   NewUpgradeShopState(),
	}
}

func (u *UpgradeShopUI) Process(player *entities.Player, inputState input.InputState) Result {
	if inputState.CloseShop {
		return Close()
	}

	// Handle navigation
	if inputState.PrevTab {
		u.state.PrevTab()
	}
	if inputState.NextTab {
		u.state.NextTab()
	}
	if inputState.NavLeft {
		u.state.NavigateLeft()
	}
	if inputState.NavRight {
		u.state.NavigateRight()
	}
	if inputState.NavUp {
		u.state.NavigateUp()
	}
	if inputState.NavDown {
		u.state.NavigateDown()
	}

	// Handle purchase
	if inputState.Interact {
		if effs := u.tryPurchase(player); effs != nil {
			return WithEffects(effs...)
		}
	}

	return NoChange()
}

func (u *UpgradeShopUI) tryPurchase(player *entities.Player) []effects.Effect {
	selectedTier := u.state.SelectedTier
	currentTier := entities.GetPlayerCurrentTier(player, u.state.ActiveTab)

	// Cannot purchase already-owned tiers
	if selectedTier <= currentTier {
		return nil
	}

	// Check price and affordability
	price := u.catalog.GetUpgradePrice(u.state.ActiveTab, selectedTier)
	if !player.CanAfford(price) {
		return nil
	}

	// Create effects for the purchase
	effs := []effects.Effect{effects.TakeMoney{Amount: price}}

	// Add the upgrade effect
	switch u.state.ActiveTab {
	case entities.UpgradeEngine:
		if entry := u.catalog.GetEngineCatalogEntry(selectedTier); entry != nil {
			effs = append(effs, effects.SetEngine{Engine: entry.Engine})
		}
	case entities.UpgradeHull:
		if entry := u.catalog.GetHullCatalogEntry(selectedTier); entry != nil {
			effs = append(effs, effects.SetHull{Hull: entry.Hull})
		}
	case entities.UpgradeFuelTank:
		if entry := u.catalog.GetFuelTankCatalogEntry(selectedTier); entry != nil {
			effs = append(effs, effects.SetFuelTank{FuelTank: entry.FuelTank})
		}
	case entities.UpgradeCargoHold:
		if entry := u.catalog.GetCargoHoldCatalogEntry(selectedTier); entry != nil {
			effs = append(effs, effects.SetCargoHold{CargoHold: entry.CargoHold})
		}
	case entities.UpgradeHeatShield:
		if entry := u.catalog.GetHeatShieldCatalogEntry(selectedTier); entry != nil {
			effs = append(effs, effects.SetHeatShield{HeatShield: entry.HeatShield})
		}
	case entities.UpgradeDrill:
		if entry := u.catalog.GetDrillCatalogEntry(selectedTier); entry != nil {
			effs = append(effs, effects.SetDrill{Drill: entry.Drill})
		}
	}

	return effs
}

func (u *UpgradeShopUI) GetRenderState() interface{} {
	return u.state
}

func (u *UpgradeShopUI) GetCatalog() *entities.UpgradeCatalog {
	return u.catalog
}

func (u *UpgradeShopUI) ResetState() {
	u.state.Reset()
}
