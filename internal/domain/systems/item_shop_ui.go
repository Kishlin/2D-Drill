package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type ItemShopUISystem struct {
	shop    *entities.ItemShop
	uiState *entities.ItemShopUIState
}

func NewItemShopUISystem(shop *entities.ItemShop) *ItemShopUISystem {
	return &ItemShopUISystem{
		shop:    shop,
		uiState: entities.NewItemShopUIState(),
	}
}

func (iss *ItemShopUISystem) ProcessItemShopInteraction(player *entities.Player, inputState input.InputState) {
	if !iss.uiState.Open {
		iss.processClosedShop(player, inputState)
	} else {
		iss.processOpenShop(player, inputState)
	}
}

func (iss *ItemShopUISystem) processClosedShop(player *entities.Player, inputState input.InputState) {
	if !inputState.Sell {
		return
	}

	if !iss.shop.IsPlayerInRange(player) {
		return
	}

	iss.uiState.OpenShop()
	player.InShop = true
}

func (iss *ItemShopUISystem) processOpenShop(player *entities.Player, inputState input.InputState) {
	// Handle close
	if inputState.CloseShop {
		iss.uiState.CloseShop()
		player.InShop = false
		return
	}

	// Handle navigation
	if inputState.NavLeft {
		iss.uiState.NavigateLeft()
	}
	if inputState.NavRight {
		iss.uiState.NavigateRight()
	}
	if inputState.NavUp {
		iss.uiState.NavigateUp()
	}
	if inputState.NavDown {
		iss.uiState.NavigateDown()
	}

	// Handle purchase
	if inputState.Sell {
		iss.tryPurchase(player, iss.uiState.SelectedIndex)
	}
}

func (iss *ItemShopUISystem) tryPurchase(player *entities.Player, index int) {
	catalogEntry := iss.shop.GetItem(index)
	if catalogEntry == nil {
		return
	}

	price := catalogEntry.Price
	if !player.CanAfford(price) {
		return
	}

	player.Money -= price
	player.AddItem(catalogEntry.ItemType)
}

func (iss *ItemShopUISystem) GetUIState() *entities.ItemShopUIState {
	return iss.uiState
}

func (iss *ItemShopUISystem) GetShop() *entities.ItemShop {
	return iss.shop
}
