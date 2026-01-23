package upgrades

import "github.com/Kishlin/drill-game/internal/domain/config"

type Hull struct {
	tier  int
	name  string
	maxHP float32
}

func (h Hull) Tier() int {
	return h.tier
}

func (h Hull) Name() string {
	return h.name
}

func (h Hull) Type() UpgradeType {
	return TypeHull
}

func (h Hull) MaxHP() float32 {
	return h.maxHP
}

func NewHullFromConfig(tier int, name string, stats config.HullStats) Hull {
	return Hull{
		tier:  tier,
		name:  name,
		maxHP: stats.MaxHP,
	}
}
