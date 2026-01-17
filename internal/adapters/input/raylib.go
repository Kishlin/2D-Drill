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
		// Continuous inputs (for movement)
		Left:  rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA),
		Right: rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD),
		Up:    rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW),
		Drill: rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS),

		// Discrete inputs (single press actions)
		Sell:        rl.IsKeyPressed(rl.KeyE),
		UseTeleport: rl.IsKeyPressed(rl.KeyT),
		UseRepair:   rl.IsKeyPressed(rl.KeyR),
		UseRefuel:   rl.IsKeyPressed(rl.KeyF),
		UseBomb:     rl.IsKeyPressed(rl.KeyB),
		UseBigBomb:  rl.IsKeyPressed(rl.KeyG),
		PrevTab:     rl.IsKeyPressed(rl.KeyZ),
		NextTab:     rl.IsKeyPressed(rl.KeyX),
		CloseShop:   rl.IsKeyPressed(rl.KeyQ) || rl.IsKeyPressed(rl.KeyEscape),

		// Discrete navigation (for UI)
		NavLeft:  rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA),
		NavRight: rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD),
		NavUp:    rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW),
		NavDown:  rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS),
	}
}
