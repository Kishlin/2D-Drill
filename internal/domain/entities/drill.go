package entities

import "github.com/Kishlin/drill-game/internal/domain/config"

type Drill struct {
	tier       int
	name       string
	drillSpeed float32
}

func (d Drill) Tier() int {
	return d.tier
}

func (d Drill) Name() string {
	return d.name
}

func (d Drill) DrillSpeed() float32 {
	return d.drillSpeed
}

func NewDrillFromConfig(tier int, name string, stats config.DrillStats) Drill {
	return Drill{
		tier:       tier,
		name:       name,
		drillSpeed: stats.DrillSpeed,
	}
}
