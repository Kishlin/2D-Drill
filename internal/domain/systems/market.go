package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type MarketSystem struct {
	market *entities.Market
}

func NewMarketSystem(market *entities.Market) *MarketSystem {
	return &MarketSystem{market: market}
}

// ProcessSelling handles player selling inventory at market
func (ms *MarketSystem) ProcessSelling(
	player *entities.Player,
	inputState input.InputState,
) {
	if !inputState.Sell {
		return
	}

	if !ms.market.IsPlayerInRange(player) {
		return
	}

	player.SellInventory()
}

// GetMarket returns the market entity for rendering
func (ms *MarketSystem) GetMarket() *entities.Market {
	return ms.market
}
