package entities

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

func NewFuelTankBase() FuelTank {
	return FuelTank{
		tier:     0,
		name:     "Base Tank",
		capacity: 10.0,
	}
}

func NewFuelTankMk1() FuelTank {
	return FuelTank{
		tier:     1,
		name:     "Tank Mk1",
		capacity: 15.0,
	}
}

func NewFuelTankMk2() FuelTank {
	return FuelTank{
		tier:     2,
		name:     "Tank Mk2",
		capacity: 22.0,
	}
}

func NewFuelTankMk3() FuelTank {
	return FuelTank{
		tier:     3,
		name:     "Tank Mk3",
		capacity: 32.0,
	}
}

func NewFuelTankMk4() FuelTank {
	return FuelTank{
		tier:     4,
		name:     "Tank Mk4",
		capacity: 45.0,
	}
}

func NewFuelTankMk5() FuelTank {
	return FuelTank{
		tier:     5,
		name:     "Tank Mk5",
		capacity: 65.0,
	}
}
