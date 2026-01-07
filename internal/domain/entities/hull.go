package entities

type Hull struct {
	tier  int
	name  string
	maxHP float32
}

func (h Hull) Tier() int {
	return h.tier
}

func (h Hull) Name() string {
	return h.name
}

func (h Hull) MaxHP() float32 {
	return h.maxHP
}

func NewHullBase() Hull {
	return Hull{
		tier:  0,
		name:  "Base Hull",
		maxHP: 10.0,
	}
}

func NewHullMk1() Hull {
	return Hull{
		tier:  1,
		name:  "Hull Mk1",
		maxHP: 15.0,
	}
}

func NewHullMk2() Hull {
	return Hull{
		tier:  2,
		name:  "Hull Mk2",
		maxHP: 20.0,
	}
}

func NewHullMk3() Hull {
	return Hull{
		tier:  3,
		name:  "Hull Mk3",
		maxHP: 30.0,
	}
}

func NewHullMk4() Hull {
	return Hull{
		tier:  4,
		name:  "Hull Mk4",
		maxHP: 45.0,
	}
}

func NewHullMk5() Hull {
	return Hull{
		tier:  5,
		name:  "Hull Mk5",
		maxHP: 75.0,
	}
}
