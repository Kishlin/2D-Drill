package ui

import (
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type ItemShopUI struct {
	catalog *entities.ItemCatalog
	state   *ItemShopState
}

func NewItemShopUI(catalog *entities.ItemCatalog) *ItemShopUI {
	return &ItemShopUI{
		catalog: catalog,
		state:   NewItemShopState(),
	}
}

func (u *ItemShopUI) Process(player *entities.Player, inputState input.InputState) Result {
	if inputState.CloseShop {
		return Close()
	}

	// Handle navigation
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

func (u *ItemShopUI) tryPurchase(player *entities.Player) []effects.Effect {
	catalogEntry := u.catalog.GetItem(u.state.Selected)
	if catalogEntry == nil {
		return nil
	}

	price := catalogEntry.Price
	if player.CanAfford(price) == false {
		return nil
	}

	return []effects.Effect{
		effects.TakeMoney{Amount: price},
		effects.AddItem{ItemType: catalogEntry.ItemType},
	}
}

func (u *ItemShopUI) GetRenderState() interface{} {
	return u.state
}

func (u *ItemShopUI) GetCatalog() *entities.ItemCatalog {
	return u.catalog
}

func (u *ItemShopUI) ResetState() {
	u.state.Reset()
}
