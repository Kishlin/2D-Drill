package statemachine

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
)

// StateID identifies a state by a type-safe integer constant.
// Each boss package defines its own iota-based constants.
type StateID int

// StateIDNone represents no state / stay in current state
const StateIDNone StateID = -1

// StateContext provides data to state handlers
type StateContext struct {
	Player  *entities.Player
	Dt      float32
	Elapsed float32 // Time since state entered
}

// StateResult is returned by OnUpdate to signal transitions and spawn requests
type StateResult struct {
	NextState     StateID // StateIDNone = stay in current state
	SpawnRequests []projectiles.SpawnRequest
}

// State defines a declarative state with lifecycle hooks
type State struct {
	ID StateID

	OnEnter  func(ctx *StateContext)
	OnUpdate func(ctx *StateContext) StateResult
	OnExit   func(ctx *StateContext)
}
