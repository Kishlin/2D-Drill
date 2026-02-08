package ui

import (
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type Result struct {
	ShouldClose bool
	Effects     []effects.Effect
}

func NoChange() Result {
	return Result{}
}

func Close() Result {
	return Result{ShouldClose: true}
}

func WithEffects(effs ...effects.Effect) Result {
	return Result{Effects: effs}
}

func CloseWithEffects(effs ...effects.Effect) Result {
	return Result{ShouldClose: true, Effects: effs}
}

// UI interface - all UIs implement this
type UI interface {
	// Process handles input, returns result
	// Modal UIs return ShouldClose=false to stay open
	// Instant UIs return ShouldClose=true immediately (with effects)
	Process(player *entities.Player, inputState input.InputState) Result

	// GetRenderState returns UI-specific state for rendering (nil for instant UIs)
	GetRenderState() interface{}

	// ResetState resets the UI to its initial state
	ResetState()
}
