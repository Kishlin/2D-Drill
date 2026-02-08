package test_boss

import (
	"math/rand"

	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/bosses/attacks"
	"github.com/Kishlin/drill-game/internal/domain/bosses/movement"
	"github.com/Kishlin/drill-game/internal/domain/bosses/phases"
	"github.com/Kishlin/drill-game/internal/domain/bosses/statemachine"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// phaseConfig holds TestBoss-specific parameters that vary per phase.
type phaseConfig struct {
	MovementSpeed      float32
	ProjectileCooldown float32
	AOECooldown        float32
}

func init() {
	bosses.Register("test_boss", func(roomStartY, worldWidth float32) bosses.Boss {
		return New(roomStartY, worldWidth)
	})
}

const (
	MaxHP         = 100.0
	Width         = 100.0
	Height        = 100.0
	BaseSpeed     = 80.0
	ContactDamage = 20.0 // Damage per second on contact

	// Animation timings
	WindupDuration  = 1.0 // Vibration warning before slam
	SlamDuration    = 0.3 // Actual slam damage window
	DoubleSlamPause = 0.4 // Pause between slams in phase 3

	// Vulnerability durations per phase (after slam)
	Phase2VulnerableDuration = 3.0
	Phase3VulnerableDuration = 2.0

	// Projectile parameters
	ProjectileCount = 3
	ProjectileSpeed = 200.0
	ProjectileSize  = 16.0
	ProjectileDamage = 5.0
)

// HP thresholds for the phase manager (generic infrastructure)
var phaseThresholds = []phases.Config{
	{HPThreshold: 0.66}, // Phase 1: 100% - 66% HP
	{HPThreshold: 0.33}, // Phase 2: 66% - 33% HP
	{HPThreshold: 0.0},  // Phase 3: 33% - 0% HP
}

// Boss-specific parameters per phase
var phaseConfigs = []phaseConfig{
	// Phase 1
	{MovementSpeed: BaseSpeed, ProjectileCooldown: 3.0, AOECooldown: 0},
	// Phase 2
	{MovementSpeed: BaseSpeed * 1.25, ProjectileCooldown: 2.0, AOECooldown: 6.0},
	// Phase 3
	{MovementSpeed: BaseSpeed * 1.5, ProjectileCooldown: 1.0, AOECooldown: 4.0},
}

type TestBoss struct {
	*bosses.BaseBoss

	// Boss-specific components
	movement         *movement.Grounded
	projectileAttack *attacks.ProjectileAttack
	worldWidth       float32
	floorY           float32

	// Boss-specific data (accessed by state behaviors)
	aoeCooldown float32
	slamCount   int
	maxSlams    int
	aoeRadius   float32
	aoeDamage   float32
	aoePosition types.Vec2
}

func New(roomStartY, worldWidth float32) *TestBoss {
	floorY := roomStartY + 680 - Height
	centerX := (worldWidth - Width) / 2

	// Create movement behavior
	moveCfg := movement.GroundedConfig{
		Speed:     phaseConfigs[0].MovementSpeed,
		MinX:      0,
		MaxX:      worldWidth,
		FloorY:    floorY,
		BossWidth: Width,
	}
	groundedMovement := movement.NewGrounded(moveCfg)

	// Create projectile attack
	projCfg := attacks.ProjectileAttackConfig{
		Cooldown:        phaseConfigs[0].ProjectileCooldown,
		ProjectileCount: ProjectileCount,
		ProjectileSpeed: ProjectileSpeed,
		ProjectileSize:  ProjectileSize,
		Damage:          ProjectileDamage,
	}
	projAttack := attacks.NewProjectileAttack(projCfg)

	// Create base boss
	baseBoss := bosses.NewBaseBoss(bosses.BaseBossConfig{
		Position: types.NewVec2(centerX, floorY),
		MaxHP:    MaxHP,
		BoxSet: bosses.NewBodyBoxSet(bosses.BodyBoxConfig{
			ID:               "body",
			Width:            Width,
			Height:           Height,
			DamagePerSec:     ContactDamage,
			DamageMultiplier: 1.0,
		}),
		Phases: phaseThresholds,
	})

	b := &TestBoss{
		BaseBoss:         baseBoss,
		movement:         groundedMovement,
		projectileAttack: projAttack,
		worldWidth:       worldWidth,
		floorY:           floorY,
		aoeCooldown:      phaseConfigs[0].AOECooldown,
		slamCount:        0,
		maxSlams:         1,
		aoeRadius:        150.0,
		aoeDamage:        15.0,
	}

	// Configure self reference for virtual dispatch and handlers
	b.Self = b
	b.PhaseChangeHandler = b
	b.DamageReactionHandler = b

	// Build state machine
	b.SetStateMachine(statemachine.NewStateMachine(b.buildStates(), StatePatrol))

	return b
}

// vulnerableDuration returns the vulnerability window duration for the current phase
func (b *TestBoss) vulnerableDuration() float32 {
	switch b.PhaseManager.GetCurrentPhase() {
	case 1:
		return Phase2VulnerableDuration
	case 2:
		return Phase3VulnerableDuration
	default:
		return 0
	}
}

// OnPhaseChange implements PhaseChangeHandler
func (b *TestBoss) OnPhaseChange(phaseIndex int) {
	cfg := phaseConfigs[phaseIndex]

	// Update movement speed
	b.movement.SetSpeed(cfg.MovementSpeed)

	// Update projectile cooldown
	b.projectileAttack = attacks.NewProjectileAttack(attacks.ProjectileAttackConfig{
		Cooldown:        cfg.ProjectileCooldown,
		ProjectileCount: ProjectileCount,
		ProjectileSpeed: ProjectileSpeed,
		ProjectileSize:  ProjectileSize,
		Damage:          ProjectileDamage,
	})

	// Reset slam cooldown for new phase
	if b.StateMachine.CurrentState() == StatePatrol {
		b.aoeCooldown = cfg.AOECooldown
	}
}

// OnDamageReceived implements DamageReactionHandler
func (b *TestBoss) OnDamageReceived(hurtboxID string, damage float32) {
	// End vulnerability on damage if in vulnerable state
	if b.StateMachine.CurrentState() == StateVulnerable {
		b.StateMachine.TransitionTo(StatePatrol, &statemachine.StateContext{})
	}
}

// buildStates creates the state machine states with direct access to boss fields
func (b *TestBoss) buildStates() map[statemachine.StateID]*statemachine.State {
	return map[statemachine.StateID]*statemachine.State{
		StatePatrol: {
			ID:      StatePatrol,
			CanMove: true,
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				b.Position = b.movement.Update(b.Position, ctx.Dt)
				spawnRequests := b.updateProjectileAttack(ctx.Dt)

				// Check if should start slam (only in phases with AOE)
				if b.hasAOEAttack() {
					b.aoeCooldown -= ctx.Dt
					if b.aoeCooldown <= 0 {
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
				b.slamCount = 0
				b.determineMaxSlams()
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
				b.dealAOEDamage(ctx.Dt)

				if ctx.Elapsed >= SlamDuration {
					b.slamCount++

					// Check if more slams to do
					if b.slamCount < b.maxSlams {
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
				if ctx.Elapsed >= b.vulnerableDuration() {
					return statemachine.StateResult{NextState: StatePatrol}
				}
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			},
			OnExit: func(ctx *statemachine.StateContext) {
				b.aoeCooldown = b.currentPhaseConfig().AOECooldown
			},
		},
	}
}

// updateProjectileAttack handles projectile spawning during patrol
func (b *TestBoss) updateProjectileAttack(dt float32) []projectiles.SpawnRequest {
	if b.CurrentPlayer == nil {
		return nil
	}
	bossAABB := types.AABB{X: b.Position.X, Y: b.Position.Y, Width: Width, Height: Height}
	return b.projectileAttack.Update(bossAABB, b.CurrentPlayer.AABB, dt)
}

// currentPhaseConfig returns the boss-specific config for the current phase
func (b *TestBoss) currentPhaseConfig() phaseConfig {
	return phaseConfigs[b.PhaseManager.GetCurrentPhase()]
}

// hasAOEAttack returns true if the current phase has AOE attacks
func (b *TestBoss) hasAOEAttack() bool {
	return b.currentPhaseConfig().AOECooldown > 0
}

// determineMaxSlams sets up the slam sequence for this attack
func (b *TestBoss) determineMaxSlams() {
	phase := b.PhaseManager.GetCurrentPhase()
	if phase >= 2 && rand.Float32() < 0.5 {
		b.maxSlams = 2 // 50% chance of double slam in phase 3
	} else {
		b.maxSlams = 1
	}
	// Store AOE position at boss's feet
	b.aoePosition = types.NewVec2(
		b.Position.X+Width/2,
		b.Position.Y+Height,
	)
}

// dealAOEDamage applies damage to player if within AOE radius
func (b *TestBoss) dealAOEDamage(dt float32) {
	if b.CurrentPlayer == nil {
		return
	}
	playerCenterX := b.CurrentPlayer.AABB.X + b.CurrentPlayer.AABB.Width/2
	playerCenterY := b.CurrentPlayer.AABB.Y + b.CurrentPlayer.AABB.Height/2
	dx := playerCenterX - b.aoePosition.X
	dy := playerCenterY - b.aoePosition.Y
	distSq := dx*dx + dy*dy
	radiusSq := b.aoeRadius * b.aoeRadius

	if distSq <= radiusSq {
		b.CurrentPlayer.DealDamage(b.aoeDamage * dt / SlamDuration)
	}
}

func (b *TestBoss) Update(player *entities.Player, dt float32) []projectiles.SpawnRequest {
	return b.BaseUpdate(player, dt)
}

// GetHurtboxes returns hurtboxes based on phase and state.
// Phase 1: always has hurtboxes
// Phase 2+: only has hurtboxes during StateVulnerable
func (b *TestBoss) GetHurtboxes() []bosses.Hurtbox {
	if b.PhaseManager.GetCurrentPhase() == 0 || b.StateMachine.CurrentState() == StateVulnerable {
		return b.BoxSet.Hurtboxes
	}
	return []bosses.Hurtbox{}
}

// IsVulnerable returns true if the boss can currently take damage (for rendering)
func (b *TestBoss) IsVulnerable() bool {
	return len(b.GetHurtboxes()) > 0
}

// GetVulnerableTimer returns remaining vulnerability duration (for UI rendering)
// Returns -1 in phase 1 (always vulnerable), 0 when not vulnerable, or remaining time
func (b *TestBoss) GetVulnerableTimer() float32 {
	if b.PhaseManager.GetCurrentPhase() == 0 {
		return -1
	}
	if b.StateMachine.CurrentState() == StateVulnerable {
		return b.vulnerableDuration() - b.StateMachine.Elapsed()
	}
	return 0
}

func (b *TestBoss) GetAOEInfo() *bosses.AOEInfo {
	currentState := b.StateMachine.CurrentState()
	elapsed := b.StateMachine.Elapsed()

	switch currentState {
	case StateWindup:
		return &bosses.AOEInfo{
			Position:    b.aoePosition,
			Radius:      b.aoeRadius,
			IsTelegraph: true,
			IsDamaging:  false,
			StateTimer:  WindupDuration - elapsed,
		}
	case StateWindupBetween:
		return &bosses.AOEInfo{
			Position:    b.aoePosition,
			Radius:      b.aoeRadius,
			IsTelegraph: true,
			IsDamaging:  false,
			StateTimer:  DoubleSlamPause - elapsed,
		}
	case StateSlam:
		return &bosses.AOEInfo{
			Position:    b.aoePosition,
			Radius:      b.aoeRadius,
			IsTelegraph: false,
			IsDamaging:  true,
			StateTimer:  SlamDuration - elapsed,
		}
	default:
		return nil
	}
}

// GetState returns the current state ID for rendering
func (b *TestBoss) GetState() statemachine.StateID {
	return b.StateMachine.CurrentState()
}

// GetStateTimer returns the remaining time in current state
func (b *TestBoss) GetStateTimer() float32 {
	currentState := b.StateMachine.CurrentState()
	elapsed := b.StateMachine.Elapsed()

	switch currentState {
	case StateWindup:
		return WindupDuration - elapsed
	case StateWindupBetween:
		return DoubleSlamPause - elapsed
	case StateSlam:
		return SlamDuration - elapsed
	case StateVulnerable:
		return b.vulnerableDuration() - elapsed
	default:
		return 0
	}
}
