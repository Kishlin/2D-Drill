package statemachine

// StateMachine manages state transitions and lifecycle
type StateMachine struct {
	states  map[StateID]*State
	current StateID
	elapsed float32
}

// NewStateMachine creates a state machine with the given states and initial state
func NewStateMachine(states map[StateID]*State, initial StateID) *StateMachine {
	sm := &StateMachine{
		states:  states,
		current: initial,
		elapsed: 0,
	}

	// Call OnEnter for initial state
	if state, ok := sm.states[initial]; ok && state.OnEnter != nil {
		state.OnEnter(&StateContext{Elapsed: 0})
	}

	return sm
}

// Update runs the current state's OnUpdate and handles transitions
func (sm *StateMachine) Update(ctx *StateContext) StateResult {
	ctx.Elapsed = sm.elapsed

	state, ok := sm.states[sm.current]
	if ok == false || state.OnUpdate == nil {
		sm.elapsed += ctx.Dt
		return StateResult{}
	}

	result := state.OnUpdate(ctx)

	// Handle transition if NextState is set
	if result.NextState != "" && result.NextState != sm.current {
		sm.TransitionTo(result.NextState, ctx)
	} else {
		sm.elapsed += ctx.Dt
	}

	return result
}

// TransitionTo changes to a new state, calling OnExit and OnEnter
func (sm *StateMachine) TransitionTo(next StateID, ctx *StateContext) {
	// Call OnExit for current state
	if current, ok := sm.states[sm.current]; ok && current.OnExit != nil {
		current.OnExit(ctx)
	}

	sm.current = next
	sm.elapsed = 0

	// Call OnEnter for new state
	if nextState, ok := sm.states[next]; ok && nextState.OnEnter != nil {
		ctx.Elapsed = 0
		nextState.OnEnter(ctx)
	}
}

// CurrentState returns the current state ID
func (sm *StateMachine) CurrentState() StateID {
	return sm.current
}

// Elapsed returns time since entering the current state
func (sm *StateMachine) Elapsed() float32 {
	return sm.elapsed
}

// CanMove returns whether movement is allowed in the current state
func (sm *StateMachine) CanMove() bool {
	if state, ok := sm.states[sm.current]; ok {
		return state.CanMove
	}
	return false
}
