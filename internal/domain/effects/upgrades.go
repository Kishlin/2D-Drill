package effects

import (
	"github.com/Kishlin/drill-game/internal/domain/upgrades"
)

type SetUpgrade struct {
	Upgrade upgrades.Upgrade
}

func (e SetUpgrade) Apply(ctx *EffectContext) {
	ctx.Player.SetUpgrade(e.Upgrade)
}
