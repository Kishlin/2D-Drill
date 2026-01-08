package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type ItemShopSystem struct {
	shops []*entities.ItemShop
}

func NewItemShopSystem(shops ...*entities.ItemShop) *ItemShopSystem {
	return &ItemShopSystem{
		shops: shops,
	}
}

// ProcessPurchase checks if player is in an item shop and wants to buy
func (iss *ItemShopSystem) ProcessPurchase(player *entities.Player, inputState input.InputState) {
	if !inputState.Sell {
		return
	}

	for _, shop := range iss.shops {
		if shop.IsPlayerInRange(player) {
			iss.tryPurchase(player, shop)
			return
		}
	}
}

func (iss *ItemShopSystem) tryPurchase(player *entities.Player, shop *entities.ItemShop) {
	if !player.CanAfford(shop.Price) {
		return
	}

	player.Money -= shop.Price
	player.AddItem(shop.ItemType)
}

// GetShops returns all item shops for rendering
func (iss *ItemShopSystem) GetShops() []*entities.ItemShop {
	return iss.shops
}
