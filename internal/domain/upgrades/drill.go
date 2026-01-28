package upgrades

import "github.com/Kishlin/drill-game/internal/domain/config"

type Drill struct {
	tier            int
	name            string
	speedAtSurface  float32
	speedAtMaxDepth float32
}

func (d Drill) Tier() int {
	return d.tier
}

func (d Drill) Name() string {
	return d.name
}

func (d Drill) Type() UpgradeType {
	return TypeDrill
}

func (d Drill) SpeedAtSurface() float32 {
	return d.speedAtSurface
}

func (d Drill) SpeedAtMaxDepth() float32 {
	return d.speedAtMaxDepth
}

func NewDrillFromConfig(tier int, name string, stats config.DrillStats) Drill {
	return Drill{
		tier:            tier,
		name:            name,
		speedAtSurface:  stats.SpeedAtSurface,
		speedAtMaxDepth: stats.SpeedAtMaxDepth,
	}
}
