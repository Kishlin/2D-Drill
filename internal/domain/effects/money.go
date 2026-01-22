package effects

import "github.com/Kishlin/drill-game/internal/domain/entities"

type TakeMoney struct {
	Amount int
}

func (e TakeMoney) Apply(player *entities.Player) {
	player.Money -= e.Amount
}

type AddMoney struct {
	Amount int
}

func (e AddMoney) Apply(player *entities.Player) {
	player.Money += e.Amount
}
