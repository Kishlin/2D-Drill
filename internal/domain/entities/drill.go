package entities

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

func NewDrillBase() Drill {
	return Drill{
		tier:       0,
		name:       "Base Drill",
		drillSpeed: 1.0,
	}
}

func NewDrillMk1() Drill {
	return Drill{
		tier:       1,
		name:       "Drill Mk1",
		drillSpeed: 2.0,
	}
}

func NewDrillMk2() Drill {
	return Drill{
		tier:       2,
		name:       "Drill Mk2",
		drillSpeed: 3.0,
	}
}

func NewDrillMk3() Drill {
	return Drill{
		tier:       3,
		name:       "Drill Mk3",
		drillSpeed: 4.0,
	}
}

func NewDrillMk4() Drill {
	return Drill{
		tier:       4,
		name:       "Drill Mk4",
		drillSpeed: 5.0,
	}
}

func NewDrillMk5() Drill {
	return Drill{
		tier:       5,
		name:       "Drill Mk5",
		drillSpeed: 6.0,
	}
}
