package statemachine

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/projectiles"
)

func TestStateMachine_StartsInInitialState(t *testing.T) {
	states := map[StateID]*State{
		"idle": {ID: "idle", CanMove: false},
	}
	sm := NewStateMachine(states, "idle")

	if sm.CurrentState() != "idle" {
		t.Errorf("expected state 'idle', got '%s'", sm.CurrentState())
	}
}

func TestStateMachine_CallsOnEnterOnInit(t *testing.T) {
	enterCalled := false
	states := map[StateID]*State{
		"idle": {
			ID:      "idle",
			CanMove: false,
			OnEnter: func(ctx *StateContext) { enterCalled = true },
		},
	}
	NewStateMachine(states, "idle")

	if enterCalled == false {
		t.Error("expected OnEnter to be called on init")
	}
}

func TestStateMachine_ElapsedIncrements(t *testing.T) {
	states := map[StateID]*State{
		"idle": {ID: "idle", CanMove: false},
	}
	sm := NewStateMachine(states, "idle")

	ctx := &StateContext{Dt: 0.5}
	sm.Update(ctx)

	if sm.Elapsed() != 0.5 {
		t.Errorf("expected elapsed 0.5, got %f", sm.Elapsed())
	}

	ctx.Dt = 0.3
	sm.Update(ctx)

	if sm.Elapsed() != 0.8 {
		t.Errorf("expected elapsed 0.8, got %f", sm.Elapsed())
	}
}

func TestStateMachine_TransitionResetsElapsed(t *testing.T) {
	states := map[StateID]*State{
		"idle":   {ID: "idle", CanMove: false},
		"active": {ID: "active", CanMove: true},
	}
	sm := NewStateMachine(states, "idle")

	// Accumulate time
	ctx := &StateContext{Dt: 1.0}
	sm.Update(ctx)

	// Transition
	sm.TransitionTo("active", ctx)

	if sm.Elapsed() != 0 {
		t.Errorf("expected elapsed 0 after transition, got %f", sm.Elapsed())
	}
}

func TestStateMachine_TransitionCallsOnExitAndOnEnter(t *testing.T) {
	exitCalled := false
	enterCalled := false

	states := map[StateID]*State{
		"idle": {
			ID:     "idle",
			OnExit: func(ctx *StateContext) { exitCalled = true },
		},
		"active": {
			ID:      "active",
			OnEnter: func(ctx *StateContext) { enterCalled = true },
		},
	}
	sm := NewStateMachine(states, "idle")

	ctx := &StateContext{}
	sm.TransitionTo("active", ctx)

	if exitCalled == false {
		t.Error("expected OnExit to be called")
	}
	if enterCalled == false {
		t.Error("expected OnEnter to be called")
	}
}

func TestStateMachine_OnUpdateCanTransition(t *testing.T) {
	states := map[StateID]*State{
		"idle": {
			ID: "idle",
			OnUpdate: func(ctx *StateContext) StateResult {
				if ctx.Elapsed >= 1.0 {
					return StateResult{NextState: "active"}
				}
				return StateResult{}
			},
		},
		"active": {ID: "active"},
	}
	sm := NewStateMachine(states, "idle")

	// Update but don't hit threshold (elapsed=0 during check, becomes 0.5 after)
	ctx := &StateContext{Dt: 0.5}
	sm.Update(ctx)
	if sm.CurrentState() != "idle" {
		t.Error("should still be in idle")
	}

	// Update again (elapsed=0.5 during check, becomes 1.0 after)
	ctx.Dt = 0.5
	sm.Update(ctx)
	if sm.CurrentState() != "idle" {
		t.Error("should still be in idle at elapsed=0.5")
	}

	// Update past threshold (elapsed=1.0 during check, triggers transition)
	ctx.Dt = 0.1
	sm.Update(ctx)
	if sm.CurrentState() != "active" {
		t.Errorf("expected transition to active, got %s", sm.CurrentState())
	}
}

func TestStateMachine_OnUpdateReturnsSpawnRequests(t *testing.T) {
	expectedRequest := projectiles.SpawnRequest{Damage: 10}
	states := map[StateID]*State{
		"shooting": {
			ID: "shooting",
			OnUpdate: func(ctx *StateContext) StateResult {
				return StateResult{
					SpawnRequests: []projectiles.SpawnRequest{expectedRequest},
				}
			},
		},
	}
	sm := NewStateMachine(states, "shooting")

	ctx := &StateContext{Dt: 0.1}
	result := sm.Update(ctx)

	if len(result.SpawnRequests) != 1 {
		t.Fatal("expected 1 spawn request")
	}
	if result.SpawnRequests[0].Damage != 10 {
		t.Error("spawn request damage mismatch")
	}
}

func TestStateMachine_CanMove(t *testing.T) {
	states := map[StateID]*State{
		"patrol": {ID: "patrol", CanMove: true},
		"stunned": {ID: "stunned", CanMove: false},
	}

	sm := NewStateMachine(states, "patrol")
	if sm.CanMove() == false {
		t.Error("patrol should allow movement")
	}

	sm.TransitionTo("stunned", &StateContext{})
	if sm.CanMove() == true {
		t.Error("stunned should not allow movement")
	}
}

func TestStateMachine_ContextReceivesElapsed(t *testing.T) {
	var receivedElapsed float32
	states := map[StateID]*State{
		"idle": {
			ID: "idle",
			OnUpdate: func(ctx *StateContext) StateResult {
				receivedElapsed = ctx.Elapsed
				return StateResult{}
			},
		},
	}
	sm := NewStateMachine(states, "idle")

	ctx := &StateContext{Dt: 0.5}
	sm.Update(ctx)
	sm.Update(ctx)

	// After second update, elapsed should be 0.5 (from first update)
	// because elapsed is set before OnUpdate runs
	if receivedElapsed != 0.5 {
		t.Errorf("expected received elapsed 0.5, got %f", receivedElapsed)
	}
}
