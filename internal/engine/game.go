package engine

import (
	"github.com/Kishlin/drill-game/internal/entities"
	"github.com/Kishlin/drill-game/internal/systems"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	GroundColor = rl.Brown
	SkyColor    = rl.SkyBlue
)

type Game struct {
	player        *entities.Player
	inputSystem   *systems.InputSystem
	physicsSystem *systems.PhysicsSystem
}

func NewGame() *Game {
	return &Game{
		player:        entities.NewPlayer(400, systems.GroundLevel-entities.PlayerHeight), // Start position on ground (GroundLevel - PlayerHeight)
		inputSystem:   systems.NewInputSystem(),
		physicsSystem: systems.NewPhysicsSystem(),
	}
}

func (g *Game) Update(dt float32) error {
	// Process input
	g.inputSystem.UpdatePlayerInput(g.player, dt)

	// Update physics
	g.physicsSystem.UpdatePhysics(g.player, dt)

	return nil
}

func (g *Game) Render() {
	// Draw sky (upper portion)
	rl.DrawRectangle(0, 0, 1280, int32(systems.GroundLevel), SkyColor)

	// Draw ground (lower portion)
	rl.DrawRectangle(0, int32(systems.GroundLevel), 1280, 720, GroundColor)

	// Draw player
	g.player.Render()
}
