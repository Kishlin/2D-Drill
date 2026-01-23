package upgrades

import "github.com/Kishlin/drill-game/internal/domain/config"

type Engine struct {
	tier            int
	name            string
	maxSpeed        float32
	acceleration    float32
	flyAcceleration float32
	maxUpwardSpeed  float32
}

func (e Engine) Tier() int {
	return e.tier
}

func (e Engine) Name() string {
	return e.name
}

func (e Engine) Type() UpgradeType {
	return TypeEngine
}

func (e Engine) MaxSpeed() float32 {
	return e.maxSpeed
}

func (e Engine) Acceleration() float32 {
	return e.acceleration
}

func (e Engine) FlyAcceleration() float32 {
	return e.flyAcceleration
}

func (e Engine) MaxUpwardSpeed() float32 {
	return e.maxUpwardSpeed
}

func NewEngineFromConfig(tier int, name string, stats config.EngineStats) Engine {
	return Engine{
		tier:            tier,
		name:            name,
		maxSpeed:        stats.MaxSpeed,
		acceleration:    stats.Acceleration,
		flyAcceleration: stats.FlyAcceleration,
		maxUpwardSpeed:  stats.MaxUpwardSpeed,
	}
}
