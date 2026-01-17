package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type UpgradeShopUISystem struct {
	shop    *entities.UpgradeShop
	uiState *entities.UpgradeShopUIState
}

func NewUpgradeShopUISystem(shop *entities.UpgradeShop) *UpgradeShopUISystem {
	return &UpgradeShopUISystem{
		shop:    shop,
		uiState: entities.NewUpgradeShopUIState(),
	}
}

func (s *UpgradeShopUISystem) ProcessShopInteraction(
	player *entities.Player,
	inputState input.InputState,
) {
	if s.uiState.Open {
		s.processOpenShop(player, inputState)
	} else {
		s.processClosedShop(player, inputState)
	}
}

func (s *UpgradeShopUISystem) processClosedShop(
	player *entities.Player,
	inputState input.InputState,
) {
	// Check if player is in range and presses E to open
	if inputState.Sell && s.shop.IsPlayerInRange(player) {
		s.openShop(player)
	}
}

func (s *UpgradeShopUISystem) processOpenShop(
	player *entities.Player,
	inputState input.InputState,
) {
	// Close shop
	if inputState.CloseShop {
		s.closeShop(player)
		return
	}

	// Tab navigation
	if inputState.PrevTab {
		s.uiState.PrevTab()
	}
	if inputState.NextTab {
		s.uiState.NextTab()
	}

	// Grid navigation
	if inputState.NavLeft {
		s.uiState.NavigateLeft()
	}
	if inputState.NavRight {
		s.uiState.NavigateRight()
	}
	if inputState.NavUp {
		s.uiState.NavigateUp()
	}
	if inputState.NavDown {
		s.uiState.NavigateDown()
	}

	// Purchase
	if inputState.Sell {
		s.tryPurchase(player)
	}
}

func (s *UpgradeShopUISystem) openShop(player *entities.Player) {
	s.uiState.OpenShop()
	player.InShop = true
}

func (s *UpgradeShopUISystem) closeShop(player *entities.Player) {
	s.uiState.CloseShop()
	player.InShop = false
}

func (s *UpgradeShopUISystem) tryPurchase(player *entities.Player) {
	selectedTier := s.uiState.SelectedTier
	currentTier := entities.GetPlayerCurrentTier(player, s.uiState.ActiveTab)

	// Cannot purchase already-owned tiers
	if selectedTier <= currentTier {
		return
	}

	// Check price and affordability
	price := s.shop.GetUpgradePrice(s.uiState.ActiveTab, selectedTier)
	if !player.CanAfford(price) {
		return
	}

	// Apply the upgrade
	s.applyUpgrade(player, s.uiState.ActiveTab, selectedTier, price)
}

func (s *UpgradeShopUISystem) applyUpgrade(player *entities.Player, upgradeType entities.UpgradeType, tier int, price int) {
	switch upgradeType {
	case entities.UpgradeEngine:
		if entry := s.shop.GetEngineCatalogEntry(tier); entry != nil {
			player.BuyEngine(entry.Engine, price)
		}
	case entities.UpgradeHull:
		if entry := s.shop.GetHullCatalogEntry(tier); entry != nil {
			player.BuyHull(entry.Hull, price)
		}
	case entities.UpgradeFuelTank:
		if entry := s.shop.GetFuelTankCatalogEntry(tier); entry != nil {
			player.BuyFuelTank(entry.FuelTank, price)
		}
	case entities.UpgradeCargoHold:
		if entry := s.shop.GetCargoCatalogEntry(tier); entry != nil {
			player.BuyCargoHold(entry.CargoHold, price)
		}
	case entities.UpgradeHeatShield:
		if entry := s.shop.GetHeatCatalogEntry(tier); entry != nil {
			player.BuyHeatShield(entry.HeatShield, price)
		}
	case entities.UpgradeDrill:
		if entry := s.shop.GetDrillCatalogEntry(tier); entry != nil {
			player.BuyDrill(entry.Drill, price)
		}
	}
}

func (s *UpgradeShopUISystem) GetShop() *entities.UpgradeShop {
	return s.shop
}

func (s *UpgradeShopUISystem) GetUIState() *entities.UpgradeShopUIState {
	return s.uiState
}
