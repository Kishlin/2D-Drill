package effects

import "github.com/Kishlin/drill-game/internal/domain/entities"

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Apply(player *entities.Player, effects []Effect) {
	for _, effect := range effects {
		effect.Apply(player)
	}
}
