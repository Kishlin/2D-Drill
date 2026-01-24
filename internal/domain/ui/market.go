package ui

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type MarketUI struct {
	oreConfigs []config.OreConfig
	state      *MarketState
}

func NewMarketUI(oreConfigs []config.OreConfig) *MarketUI {
	return &MarketUI{
		oreConfigs: oreConfigs,
		state:      NewMarketState(),
	}
}

func (u *MarketUI) Process(player *entities.Player, inputState input.InputState) Result {
	// Close on Q/Escape without selling
	if inputState.CloseShop {
		return Close()
	}

	// Skip the first frame to avoid processing the E that opened the shop
	if u.state.IsFirstFrame() {
		u.state.ClearFirstFrame()
		return NoChange()
	}

	// Sell on E press
	if inputState.Interact {
		totalValue := u.calculateValue(player)
		if totalValue == 0 {
			return Close()
		}

		return CloseWithEffects(
			effects.AddMoney{Amount: totalValue},
			effects.ClearOreInventory{},
		)
	}

	// Stay open (modal behavior)
	return NoChange()
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
	return u.state
}

func (u *MarketUI) GetOreConfigs() []config.OreConfig {
	return u.oreConfigs
}

func (u *MarketUI) ResetState() {
	u.state.Reset()
}
