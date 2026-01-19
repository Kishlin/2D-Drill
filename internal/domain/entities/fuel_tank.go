package entities

import "github.com/Kishlin/drill-game/internal/domain/config"

type FuelTank struct {
	tier     int
	name     string
	capacity float32
}

func (ft FuelTank) Tier() int {
	return ft.tier
}

func (ft FuelTank) Name() string {
	return ft.name
}

func (ft FuelTank) Capacity() float32 {
	return ft.capacity
}

func NewFuelTankFromConfig(tier int, name string, stats config.FuelTankStats) FuelTank {
	return FuelTank{
		tier:     tier,
		name:     name,
		capacity: stats.Capacity,
	}
}
