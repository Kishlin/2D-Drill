package entities

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

func NewHeatShieldBase() HeatShield {
	return HeatShield{
		tier:           0,
		name:           "Base Heat Shield",
		heatResistance: 50.0,
	}
}

func NewHeatShieldMk1() HeatShield {
	return HeatShield{
		tier:           1,
		name:           "Heat Shield Mk1",
		heatResistance: 90.0,
	}
}

func NewHeatShieldMk2() HeatShield {
	return HeatShield{
		tier:           2,
		name:           "Heat Shield Mk2",
		heatResistance: 140.0,
	}
}

func NewHeatShieldMk3() HeatShield {
	return HeatShield{
		tier:           3,
		name:           "Heat Shield Mk3",
		heatResistance: 190.0,
	}
}

func NewHeatShieldMk4() HeatShield {
	return HeatShield{
		tier:           4,
		name:           "Heat Shield Mk4",
		heatResistance: 250.0,
	}
}

func NewHeatShieldMk5() HeatShield {
	return HeatShield{
		tier:           5,
		name:           "Heat Shield Mk5",
		heatResistance: 320.0,
	}
}
