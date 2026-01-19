package entities

import "github.com/Kishlin/drill-game/internal/domain/config"

type HeatShield struct {
	tier           int
	name           string
	heatResistance float32
}

func (hs HeatShield) Tier() int {
	return hs.tier
}

func (hs HeatShield) Name() string {
	return hs.name
}

func (hs HeatShield) HeatResistance() float32 {
	return hs.heatResistance
}

func NewHeatShieldFromConfig(tier int, name string, stats config.HeatShieldStats) HeatShield {
	return HeatShield{
		tier:           tier,
		name:           name,
		heatResistance: stats.HeatResistance,
	}
}
