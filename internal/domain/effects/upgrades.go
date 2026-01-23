package effects

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/upgrades"
)

type SetUpgrade struct {
	Upgrade upgrades.Upgrade
}

func (e SetUpgrade) Apply(player *entities.Player) {
	player.SetUpgrade(e.Upgrade)
}
