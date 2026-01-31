package test_boss

import "github.com/Kishlin/drill-game/internal/domain/bosses/statemachine"

// State IDs
const (
	StatePatrol statemachine.StateID = iota
	StateWindup
	StateWindupBetween
	StateSlam
	StateVulnerable
)
