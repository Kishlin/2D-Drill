package effects

import "github.com/Kishlin/drill-game/internal/domain/entities"

type SetEngine struct {
	Engine entities.Engine
}

func (e SetEngine) Apply(player *entities.Player) {
	player.Engine = e.Engine
}

type SetHull struct {
	Hull entities.Hull
}

func (e SetHull) Apply(player *entities.Player) {
	player.Hull = e.Hull
}

type SetFuelTank struct {
	FuelTank entities.FuelTank
}

func (e SetFuelTank) Apply(player *entities.Player) {
	player.FuelTank = e.FuelTank
}

type SetCargoHold struct {
	CargoHold entities.CargoHold
}

func (e SetCargoHold) Apply(player *entities.Player) {
	player.CargoHold = e.CargoHold
}

type SetHeatShield struct {
	HeatShield entities.HeatShield
}

func (e SetHeatShield) Apply(player *entities.Player) {
	player.HeatShield = e.HeatShield
}

type SetDrill struct {
	Drill entities.Drill
}

func (e SetDrill) Apply(player *entities.Player) {
	player.Drill = e.Drill
}
