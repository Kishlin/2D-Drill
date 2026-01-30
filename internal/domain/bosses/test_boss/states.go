package test_boss

import (
	"github.com/Kishlin/drill-game/internal/domain/bosses/statemachine"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// State IDs
const (
	StatePatrol statemachine.StateID = iota
	StateWindup
	StateWindupBetween
	StateSlam
	StateVulnerable
)

// StateBehaviors provides callbacks for state handlers to interact with boss data
type StateBehaviors struct {
	// Cooldown management
	GetAOECooldown    func() float32
	SetAOECooldown    func(float32)
	DecrementCooldown func(dt float32)

	// Slam management
	GetSlamCount     func() int
	IncrementSlam    func()
	ResetSlamCount   func()
	GetMaxSlams      func() int
	SetMaxSlams      func(int)
	DetermineMaxSlams func()

	// AOE position
	SetAOEPosition func(pos types.Vec2)

	// Movement and attacks
	UpdateMovement         func(dt float32)
	UpdateProjectileAttack func(dt float32) []projectiles.SpawnRequest

	// Phase info
	GetVulnerableDuration func() float32
	HasAOEAttack          func() bool

	// Damage
	DealAOEDamage func(dt float32)

	// Vulnerability
	EndVulnerability func()
}

// BuildStates creates the state definitions for TestBoss
func BuildStates(behaviors *StateBehaviors) map[statemachine.StateID]*statemachine.State {
	return map[statemachine.StateID]*statemachine.State{
		StatePatrol: {
			ID:      StatePatrol,
			CanMove: true,
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				behaviors.UpdateMovement(ctx.Dt)
				spawnRequests := behaviors.UpdateProjectileAttack(ctx.Dt)

				// Check if should start slam (only in phases with AOE)
				if behaviors.HasAOEAttack() {
					behaviors.DecrementCooldown(ctx.Dt)
					if behaviors.GetAOECooldown() <= 0 {
						return statemachine.StateResult{
							NextState:     StateWindup,
							SpawnRequests: spawnRequests,
						}
					}
				}

				return statemachine.StateResult{NextState: statemachine.StateIDNone, SpawnRequests: spawnRequests}
			},
		},

		StateWindup: {
			ID:      StateWindup,
			CanMove: false,
			OnEnter: func(ctx *statemachine.StateContext) {
				behaviors.ResetSlamCount()
				behaviors.DetermineMaxSlams()
			},
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				if ctx.Elapsed >= WindupDuration {
					return statemachine.StateResult{NextState: StateSlam}
				}
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			},
		},

		StateWindupBetween: {
			ID:      StateWindupBetween,
			CanMove: false,
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				if ctx.Elapsed >= DoubleSlamPause {
					return statemachine.StateResult{NextState: StateSlam}
				}
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			},
		},

		StateSlam: {
			ID:      StateSlam,
			CanMove: false,
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				behaviors.DealAOEDamage(ctx.Dt)

				if ctx.Elapsed >= SlamDuration {
					behaviors.IncrementSlam()

					// Check if more slams to do
					if behaviors.GetSlamCount() < behaviors.GetMaxSlams() {
						return statemachine.StateResult{NextState: StateWindupBetween}
					}

					// Done slamming, enter vulnerable state
					return statemachine.StateResult{NextState: StateVulnerable}
				}
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			},
		},

		StateVulnerable: {
			ID:      StateVulnerable,
			CanMove: false,
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				if ctx.Elapsed >= behaviors.GetVulnerableDuration() {
					return statemachine.StateResult{NextState: StatePatrol}
				}
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			},
			OnExit: func(ctx *statemachine.StateContext) {
				behaviors.EndVulnerability()
			},
		},
	}
}
