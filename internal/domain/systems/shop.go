package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type ShopSystem struct {
	shop *entities.Shop
}

func NewShopSystem(shop *entities.Shop) *ShopSystem {
	return &ShopSystem{shop: shop}
}

// ProcessSelling handles player selling inventory at shop
func (ss *ShopSystem) ProcessSelling(
	player *entities.Player,
	inputState input.InputState,
) {
	if !inputState.Sell {
		return
	}

	if !ss.shop.IsPlayerInRange(player) {
		return
	}

	player.SellInventory()
}

// GetShop returns the shop entity for rendering
func (ss *ShopSystem) GetShop() *entities.Shop {
	return ss.shop
}
