package input

import (
	"github.com/Kishlin/drill-game/internal/domain/input"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// RaylibInputAdapter converts Raylib keyboard state to domain InputState
type RaylibInputAdapter struct{}

func NewRaylibInputAdapter() *RaylibInputAdapter {
	return &RaylibInputAdapter{}
}

// ReadInput reads Raylib keys and returns platform-agnostic InputState
func (a *RaylibInputAdapter) ReadInput() input.InputState {
	return input.InputState{
		Left:  rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA),
		Right: rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD),
		Up:    rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW),
		Dig:   rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS),
	}
}
