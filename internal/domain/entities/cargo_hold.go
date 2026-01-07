package entities

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

func (ch CargoHold) Capacity() int {
	return ch.capacity
}

func NewCargoHoldBase() CargoHold {
	return CargoHold{
		tier:     0,
		name:     "Base Cargo Hold",
		capacity: 10,
	}
}

func NewCargoHoldMk1() CargoHold {
	return CargoHold{
		tier:     1,
		name:     "Cargo Hold Mk1",
		capacity: 14,
	}
}

func NewCargoHoldMk2() CargoHold {
	return CargoHold{
		tier:     2,
		name:     "Cargo Hold Mk2",
		capacity: 18,
	}
}

func NewCargoHoldMk3() CargoHold {
	return CargoHold{
		tier:     3,
		name:     "Cargo Hold Mk3",
		capacity: 24,
	}
}

func NewCargoHoldMk4() CargoHold {
	return CargoHold{
		tier:     4,
		name:     "Cargo Hold Mk4",
		capacity: 31,
	}
}

func NewCargoHoldMk5() CargoHold {
	return CargoHold{
		tier:     5,
		name:     "Cargo Hold Mk5",
		capacity: 40,
	}
}
