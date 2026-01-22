package effects

import "github.com/Kishlin/drill-game/internal/domain/entities"

type SetFuel struct {
	Amount float32
}

func (e SetFuel) Apply(player *entities.Player) {
	player.Fuel = e.Amount
}

type SetHP struct {
	Amount float32
}

func (e SetHP) Apply(player *entities.Player) {
	player.HP = e.Amount
}
