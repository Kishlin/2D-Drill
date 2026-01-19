package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type MarketSystem struct {
	market     *entities.Market
	oreConfigs []config.OreConfig
}

func NewMarketSystemWithConfig(market *entities.Market, oreConfigs []config.OreConfig) *MarketSystem {
	return &MarketSystem{
		market:     market,
		oreConfigs: oreConfigs,
	}
}

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

	player.SellInventory(ms.oreConfigs)
}

func (ms *MarketSystem) GetMarket() *entities.Market {
	return ms.market
}

func (ms *MarketSystem) GetOreConfigs() []config.OreConfig {
	return ms.oreConfigs
}
