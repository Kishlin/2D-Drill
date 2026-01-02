package engine

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/systems"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type Game struct {
	world          *world.World
	player         *entities.Player
	physicsSystem  *systems.PhysicsSystem
	diggingSystem  *systems.DiggingSystem
}

func NewGame(w *world.World) *Game {
	// Spawn player at center of world horizontally, just above ground
	spawnX := (w.Width / 2) - (entities.PlayerWidth / 2)
	spawnY := w.GetGroundLevel() - entities.PlayerHeight - 10

	return &Game{
		world:          w,
		player:         entities.NewPlayer(spawnX, spawnY),
		physicsSystem:  systems.NewPhysicsSystem(w),
		diggingSystem:  systems.NewDiggingSystem(w),
	}
}

func (g *Game) Update(dt float32, inputState input.InputState) error {
	// 1. Handle digging (before physics, so alignment happens first)
	g.diggingSystem.ProcessDigging(g.player, inputState)

	// 2. Update physics
	g.physicsSystem.UpdatePhysics(g.player, inputState, dt)

	return nil
}

func (g *Game) GetWorld() *world.World {
	return g.world
}

func (g *Game) GetPlayer() *entities.Player {
	return g.player
}
