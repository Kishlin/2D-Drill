package entities

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

func NewEngineBase() Engine {
	return Engine{
		tier:            0,
		name:            "Base Engine",
		maxSpeed:        450.0,
		acceleration:    2500.0,
		flyAcceleration: 2500.0,
		maxUpwardSpeed:  -600.0,
	}
}

func NewEngineMk1() Engine {
	return Engine{
		tier:            1,
		name:            "Engine Mk1",
		maxSpeed:        475.0,
		acceleration:    2667.0,
		flyAcceleration: 2667.0,
		maxUpwardSpeed:  -635.0,
	}
}

func NewEngineMk2() Engine {
	return Engine{
		tier:            2,
		name:            "Engine Mk2",
		maxSpeed:        500.0,
		acceleration:    2833.0,
		flyAcceleration: 2833.0,
		maxUpwardSpeed:  -670.0,
	}
}

func NewEngineMk3() Engine {
	return Engine{
		tier:            3,
		name:            "Engine Mk3",
		maxSpeed:        525.0,
		acceleration:    3000.0,
		flyAcceleration: 3000.0,
		maxUpwardSpeed:  -705.0,
	}
}

func NewEngineMk4() Engine {
	return Engine{
		tier:            4,
		name:            "Engine Mk4",
		maxSpeed:        562.0,
		acceleration:    3250.0,
		flyAcceleration: 3250.0,
		maxUpwardSpeed:  -740.0,
	}
}

func NewEngineMk5() Engine {
	return Engine{
		tier:            5,
		name:            "Engine Mk5",
		maxSpeed:        600.0,
		acceleration:    3500.0,
		flyAcceleration: 3500.0,
		maxUpwardSpeed:  -775.0,
	}
}
