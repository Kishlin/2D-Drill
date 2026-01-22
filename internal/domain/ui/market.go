package ui

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type MarketUI struct {
	oreConfigs []config.OreConfig
}

func NewMarketUI(oreConfigs []config.OreConfig) *MarketUI {
	return &MarketUI{oreConfigs: oreConfigs}
}

func (u *MarketUI) Process(player *entities.Player, inputState input.InputState) Result {
	totalValue := u.calculateValue(player)
	if totalValue == 0 {
		return Close()
	}

	return CloseWithEffects(
		effects.AddMoney{Amount: totalValue},
		effects.ClearOreInventory{},
	)
}

func (u *MarketUI) calculateValue(player *entities.Player) int {
	totalValue := 0
	for oreID, count := range player.OreInventory {
		if count > 0 {
			for _, oreCfg := range u.oreConfigs {
				if oreCfg.ID == oreID {
					totalValue += oreCfg.Value * count
					break
				}
			}
		}
	}
	return totalValue
}

func (u *MarketUI) GetRenderState() interface{} {
	return nil
}
