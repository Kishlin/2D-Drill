package statemachine

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/projectiles"
)

// Test state IDs
const (
	testStateIdle StateID = iota
	testStateActive
	testStateShooting
)

func TestStateMachine_StartsInInitialState(t *testing.T) {
	states := map[StateID]*State{
		testStateIdle: {ID: testStateIdle},
	}
	sm := NewStateMachine(states, testStateIdle)

	if sm.CurrentState() != testStateIdle {
		t.Errorf("expected state idle, got %d", sm.CurrentState())
	}
}

func TestStateMachine_CallsOnEnterOnInit(t *testing.T) {
	enterCalled := false
	states := map[StateID]*State{
		testStateIdle: {
			ID:      testStateIdle,
			OnEnter: func(ctx *StateContext) { enterCalled = true },
		},
	}
	NewStateMachine(states, testStateIdle)

	if enterCalled == false {
		t.Error("expected OnEnter to be called on init")
	}
}

func TestStateMachine_ElapsedIncrements(t *testing.T) {
	states := map[StateID]*State{
		testStateIdle: {ID: testStateIdle},
	}
	sm := NewStateMachine(states, testStateIdle)

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
		testStateIdle:   {ID: testStateIdle},
		testStateActive: {ID: testStateActive},
	}
	sm := NewStateMachine(states, testStateIdle)

	// Accumulate time
	ctx := &StateContext{Dt: 1.0}
	sm.Update(ctx)

	// Transition
	sm.TransitionTo(testStateActive, ctx)

	if sm.Elapsed() != 0 {
		t.Errorf("expected elapsed 0 after transition, got %f", sm.Elapsed())
	}
}

func TestStateMachine_TransitionCallsOnExitAndOnEnter(t *testing.T) {
	exitCalled := false
	enterCalled := false

	states := map[StateID]*State{
		testStateIdle: {
			ID:     testStateIdle,
			OnExit: func(ctx *StateContext) { exitCalled = true },
		},
		testStateActive: {
			ID:      testStateActive,
			OnEnter: func(ctx *StateContext) { enterCalled = true },
		},
	}
	sm := NewStateMachine(states, testStateIdle)

	ctx := &StateContext{}
	sm.TransitionTo(testStateActive, ctx)

	if exitCalled == false {
		t.Error("expected OnExit to be called")
	}
	if enterCalled == false {
		t.Error("expected OnEnter to be called")
	}
}

func TestStateMachine_OnUpdateCanTransition(t *testing.T) {
	states := map[StateID]*State{
		testStateIdle: {
			ID: testStateIdle,
			OnUpdate: func(ctx *StateContext) StateResult {
				if ctx.Elapsed >= 1.0 {
					return StateResult{NextState: testStateActive}
				}
				return StateResult{NextState: StateIDNone}
			},
		},
		testStateActive: {ID: testStateActive},
	}
	sm := NewStateMachine(states, testStateIdle)

	// Update but don't hit threshold (elapsed=0 during check, becomes 0.5 after)
	ctx := &StateContext{Dt: 0.5}
	sm.Update(ctx)
	if sm.CurrentState() != testStateIdle {
		t.Error("should still be in idle")
	}

	// Update again (elapsed=0.5 during check, becomes 1.0 after)
	ctx.Dt = 0.5
	sm.Update(ctx)
	if sm.CurrentState() != testStateIdle {
		t.Error("should still be in idle at elapsed=0.5")
	}

	// Update past threshold (elapsed=1.0 during check, triggers transition)
	ctx.Dt = 0.1
	sm.Update(ctx)
	if sm.CurrentState() != testStateActive {
		t.Errorf("expected transition to active, got %d", sm.CurrentState())
	}
}

func TestStateMachine_OnUpdateReturnsSpawnRequests(t *testing.T) {
	expectedRequest := projectiles.SpawnRequest{Damage: 10}
	states := map[StateID]*State{
		testStateShooting: {
			ID: testStateShooting,
			OnUpdate: func(ctx *StateContext) StateResult {
				return StateResult{
					NextState:     StateIDNone,
					SpawnRequests: []projectiles.SpawnRequest{expectedRequest},
				}
			},
		},
	}
	sm := NewStateMachine(states, testStateShooting)

	ctx := &StateContext{Dt: 0.1}
	result := sm.Update(ctx)

	if len(result.SpawnRequests) != 1 {
		t.Fatal("expected 1 spawn request")
	}
	if result.SpawnRequests[0].Damage != 10 {
		t.Error("spawn request damage mismatch")
	}
}

func TestStateMachine_ContextReceivesElapsed(t *testing.T) {
	var receivedElapsed float32
	states := map[StateID]*State{
		testStateIdle: {
			ID: testStateIdle,
			OnUpdate: func(ctx *StateContext) StateResult {
				receivedElapsed = ctx.Elapsed
				return StateResult{NextState: StateIDNone}
			},
		},
	}
	sm := NewStateMachine(states, testStateIdle)

	ctx := &StateContext{Dt: 0.5}
	sm.Update(ctx)
	sm.Update(ctx)

	// After second update, elapsed should be 0.5 (from first update)
	// because elapsed is set before OnUpdate runs
	if receivedElapsed != 0.5 {
		t.Errorf("expected received elapsed 0.5, got %f", receivedElapsed)
	}
}
