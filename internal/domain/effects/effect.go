package effects

import "github.com/Kishlin/drill-game/internal/domain/entities"

type Effect interface {
	Apply(player *entities.Player)
}
