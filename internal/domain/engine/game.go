package engine

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/systems"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type Game struct {
	world         *world.World
	player        *entities.Player
	physicsSystem *systems.PhysicsSystem
}

func NewGame(w *world.World) *Game {
	return &Game{
		world:         w,
		player:        entities.NewPlayer(400, w.GetGroundLevel()-entities.PlayerHeight),
		physicsSystem: systems.NewPhysicsSystem(w),
	}
}

func (g *Game) Update(dt float32, inputState input.InputState) error {
	g.physicsSystem.UpdatePhysics(g.player, inputState, dt)

	return nil
}

func (g *Game) GetWorld() *world.World {
	return g.world
}

func (g *Game) GetPlayer() *entities.Player {
	return g.player
}
