package upgrades

import "github.com/Kishlin/drill-game/internal/domain/config"

type CargoHold struct {
	tier     int
	name     string
	capacity int
}

func (ch CargoHold) Tier() int {
	return ch.tier
}

func (ch CargoHold) Name() string {
	return ch.name
}

func (ch CargoHold) Type() UpgradeType {
	return TypeCargoHold
}

func (ch CargoHold) Capacity() int {
	return ch.capacity
}

func NewCargoHoldFromConfig(tier int, name string, stats config.CargoHoldStats) CargoHold {
	return CargoHold{
		tier:     tier,
		name:     name,
		capacity: stats.Capacity,
	}
}
