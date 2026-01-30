package statemachine

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
)

// StateID identifies a state by name
type StateID string

// StateContext provides data to state handlers
type StateContext struct {
	Player  *entities.Player
	Dt      float32
	Elapsed float32 // Time since state entered
}

// StateResult is returned by OnUpdate to signal transitions and spawn requests
type StateResult struct {
	NextState     StateID // Empty = stay in current state
	SpawnRequests []projectiles.SpawnRequest
}

// State defines a declarative state with lifecycle hooks
type State struct {
	ID      StateID
	CanMove bool // Movement behavior active in this state

	OnEnter  func(ctx *StateContext)
	OnUpdate func(ctx *StateContext) StateResult
	OnExit   func(ctx *StateContext)
}
