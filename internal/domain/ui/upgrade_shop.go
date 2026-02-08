package ui

import (
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/upgrades"
)

type UpgradeShopUI struct {
	catalog *upgrades.Catalog
	state   *UpgradeShopState
}

func NewUpgradeShopUI(catalog *upgrades.Catalog) *UpgradeShopUI {
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
	currentTier := player.GetUpgradeTier(u.state.ActiveTab)
	if u.state.Selected <= currentTier {
		return nil
	}

	entry := u.catalog.GetEntry(u.state.ActiveTab, u.state.Selected)
	if entry == nil || player.CanAfford(entry.Price) == false {
		return nil
	}

	return []effects.Effect{
		effects.TakeMoney{Amount: entry.Price},
		effects.SetUpgrade{Upgrade: entry.Upgrade},
	}
}

func (u *UpgradeShopUI) GetRenderState() interface{} {
	return u.state
}

func (u *UpgradeShopUI) GetCatalog() *upgrades.Catalog {
	return u.catalog
}

func (u *UpgradeShopUI) ResetState() {
	u.state.Reset()
}
