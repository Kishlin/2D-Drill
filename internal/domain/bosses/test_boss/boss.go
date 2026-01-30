package test_boss

import (
	"math/rand"

	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/bosses/attacks"
	"github.com/Kishlin/drill-game/internal/domain/bosses/movement"
	"github.com/Kishlin/drill-game/internal/domain/bosses/statemachine"
	"github.com/Kishlin/drill-game/internal/domain/components"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

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
)

// Phase configurations
var phases = []bosses.PhaseConfig{
	// Phase 1: 100% - 66% HP - Easy mode, always vulnerable
	{
		HPThreshold:        0.66,
		MovementSpeed:      BaseSpeed,
		ProjectileCooldown: 3.0,
		AOECooldown:        0, // No AOE in phase 1
		AlwaysVulnerable:   true,
		VulnerableDuration: 0,
	},
	// Phase 2: 66% - 33% HP - Medium difficulty
	{
		HPThreshold:        0.33,
		MovementSpeed:      BaseSpeed * 1.25,
		ProjectileCooldown: 2.0,
		AOECooldown:        6.0,
		AlwaysVulnerable:   false,
		VulnerableDuration: 3.0,
	},
	// Phase 3: 33% - 0% HP - Hard mode
	{
		HPThreshold:        0.0,
		MovementSpeed:      BaseSpeed * 1.5,
		ProjectileCooldown: 1.0,
		AOECooldown:        4.0,
		AlwaysVulnerable:   false,
		VulnerableDuration: 2.0,
	},
}

type TestBoss struct {
	aabb             types.AABB
	damageable       components.Damageable
	active           bool
	movement         *movement.Grounded
	projectileAttack *attacks.ProjectileAttack
	phaseManager     *bosses.PhaseManager
	worldWidth       float32
	floorY           float32

	// State machine
	stateMachine *statemachine.StateMachine

	// Boss data (accessed by state behaviors)
	aoeCooldown float32
	slamCount   int
	maxSlams    int
	aoeRadius   float32
	aoeDamage   float32
	aoePosition types.Vec2

	// Reference to player for AOE damage (set during Update)
	currentPlayer *entities.Player
}

func New(roomStartY, worldWidth float32) *TestBoss {
	floorY := roomStartY + 680 - Height
	centerX := (worldWidth - Width) / 2

	// Create movement behavior
	moveCfg := movement.GroundedConfig{
		Speed:     phases[0].MovementSpeed,
		MinX:      0,
		MaxX:      worldWidth,
		FloorY:    floorY,
		BossWidth: Width,
	}
	groundedMovement := movement.NewGrounded(moveCfg)

	// Create projectile attack
	projCfg := attacks.ProjectileAttackConfig{
		Cooldown:        phases[0].ProjectileCooldown,
		ProjectileCount: 3,
		ProjectileSpeed: 200.0,
		ProjectileSize:  16.0,
		Damage:          5.0,
	}
	projAttack := attacks.NewProjectileAttack(projCfg)

	// Create phase manager
	phaseManager := bosses.NewPhaseManager(MaxHP, phases)

	b := &TestBoss{
		aabb: types.AABB{
			X:      centerX,
			Y:      floorY,
			Width:  Width,
			Height: Height,
		},
		damageable:       components.NewDamageable(MaxHP, MaxHP),
		active:           false,
		movement:         groundedMovement,
		projectileAttack: projAttack,
		phaseManager:     phaseManager,
		worldWidth:       worldWidth,
		floorY:           floorY,
		aoeCooldown:      phases[0].AOECooldown,
		slamCount:        0,
		maxSlams:         1,
		aoeRadius:        150.0,
		aoeDamage:        15.0,
	}

	// Build state machine with behaviors
	behaviors := b.buildStateBehaviors()
	states := BuildStates(behaviors)
	b.stateMachine = statemachine.NewStateMachine(states, StatePatrol)

	return b
}

func (b *TestBoss) buildStateBehaviors() *StateBehaviors {
	return &StateBehaviors{
		GetAOECooldown:    func() float32 { return b.aoeCooldown },
		SetAOECooldown:    func(cd float32) { b.aoeCooldown = cd },
		DecrementCooldown: func(dt float32) { b.aoeCooldown -= dt },

		GetSlamCount:   func() int { return b.slamCount },
		IncrementSlam:  func() { b.slamCount++ },
		ResetSlamCount: func() { b.slamCount = 0 },
		GetMaxSlams:    func() int { return b.maxSlams },
		SetMaxSlams:    func(n int) { b.maxSlams = n },
		DetermineMaxSlams: func() {
			phase := b.phaseManager.GetCurrentPhase()
			if phase >= 2 && rand.Float32() < 0.5 {
				b.maxSlams = 2 // 50% chance of double slam in phase 3
			} else {
				b.maxSlams = 1
			}
			// Store AOE position at boss's feet
			b.aoePosition = types.NewVec2(
				b.aabb.X+b.aabb.Width/2,
				b.aabb.Y+b.aabb.Height,
			)
		},

		SetAOEPosition: func(pos types.Vec2) { b.aoePosition = pos },

		UpdateMovement: func(dt float32) {
			newPos := b.movement.Update(types.NewVec2(b.aabb.X, b.aabb.Y), dt)
			b.aabb.X = newPos.X
			b.aabb.Y = newPos.Y
		},

		UpdateProjectileAttack: func(dt float32) []projectiles.SpawnRequest {
			if b.currentPlayer == nil {
				return nil
			}
			return b.projectileAttack.Update(b.aabb, b.currentPlayer.AABB, dt)
		},

		GetVulnerableDuration: func() float32 {
			return b.phaseManager.GetCurrentConfig().VulnerableDuration
		},

		HasAOEAttack: func() bool {
			return b.phaseManager.GetCurrentConfig().AOECooldown > 0
		},

		DealAOEDamage: func(dt float32) {
			if b.currentPlayer == nil {
				return
			}
			playerCenterX := b.currentPlayer.AABB.X + b.currentPlayer.AABB.Width/2
			playerCenterY := b.currentPlayer.AABB.Y + b.currentPlayer.AABB.Height/2
			dx := playerCenterX - b.aoePosition.X
			dy := playerCenterY - b.aoePosition.Y
			distSq := dx*dx + dy*dy
			radiusSq := b.aoeRadius * b.aoeRadius

			if distSq <= radiusSq {
				b.currentPlayer.DealDamage(b.aoeDamage * dt / SlamDuration)
			}
		},

		EndVulnerability: func() {
			b.aoeCooldown = b.phaseManager.GetCurrentConfig().AOECooldown
		},
	}
}

func (b *TestBoss) Update(player *entities.Player, dt float32) []projectiles.SpawnRequest {
	if b.active == false || b.damageable.IsDefeated() {
		return nil
	}

	// Store player reference for state behaviors
	b.currentPlayer = player

	// Check for phase transitions
	if b.phaseManager.Update(b.damageable.HP) {
		b.onPhaseChange()
	}

	// Update state machine
	ctx := &statemachine.StateContext{
		Player: player,
		Dt:     dt,
	}
	result := b.stateMachine.Update(ctx)

	return result.SpawnRequests
}

func (b *TestBoss) onPhaseChange() {
	phaseCfg := b.phaseManager.GetCurrentConfig()

	// Update movement speed
	b.movement.SetSpeed(phaseCfg.MovementSpeed)

	// Update projectile cooldown
	b.projectileAttack = attacks.NewProjectileAttack(attacks.ProjectileAttackConfig{
		Cooldown:        phaseCfg.ProjectileCooldown,
		ProjectileCount: 3,
		ProjectileSpeed: 200.0,
		ProjectileSize:  16.0,
		Damage:          5.0,
	})

	// Reset slam cooldown for new phase
	if b.stateMachine.CurrentState() == StatePatrol {
		b.aoeCooldown = phaseCfg.AOECooldown
	}
}

func (b *TestBoss) GetHP() float32 {
	return b.damageable.HP
}

func (b *TestBoss) GetMaxHP() float32 {
	return b.damageable.MaxHP
}

func (b *TestBoss) IsDefeated() bool {
	return b.damageable.IsDefeated()
}

func (b *TestBoss) IsActive() bool {
	return b.active
}

func (b *TestBoss) Activate() {
	b.active = true
}

func (b *TestBoss) Deactivate() {
	b.active = false
}

func (b *TestBoss) GetAABB() types.AABB {
	return b.aabb
}

func (b *TestBoss) GetContactDamage() float32 {
	return ContactDamage
}

func (b *TestBoss) TakeDamage(damage float32) {
	if b.IsVulnerable() == false {
		return
	}

	b.damageable.TakeDamage(damage)

	// Close vulnerability window after taking damage
	if b.stateMachine.CurrentState() == StateVulnerable {
		ctx := &statemachine.StateContext{}
		b.stateMachine.TransitionTo(StatePatrol, ctx)
	}
}

func (b *TestBoss) IsVulnerable() bool {
	return b.phaseManager.IsAlwaysVulnerable() ||
		b.stateMachine.CurrentState() == StateVulnerable
}

func (b *TestBoss) GetDamageable() *components.Damageable {
	return &b.damageable
}

func (b *TestBoss) GetVulnerableTimer() float32 {
	if b.phaseManager.GetCurrentConfig().AlwaysVulnerable {
		return -1
	}
	if b.stateMachine.CurrentState() == StateVulnerable {
		return b.phaseManager.GetVulnerableDuration() - b.stateMachine.Elapsed()
	}
	return 0
}

func (b *TestBoss) GetAOEInfo() *bosses.AOEInfo {
	currentState := b.stateMachine.CurrentState()
	elapsed := b.stateMachine.Elapsed()

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

func (b *TestBoss) GetCurrentPhase() int {
	return b.phaseManager.GetCurrentPhase() + 1
}

// GetState returns the current state ID for rendering
func (b *TestBoss) GetState() statemachine.StateID {
	return b.stateMachine.CurrentState()
}

// GetStateTimer returns the remaining time in current state
func (b *TestBoss) GetStateTimer() float32 {
	currentState := b.stateMachine.CurrentState()
	elapsed := b.stateMachine.Elapsed()

	switch currentState {
	case StateWindup:
		return WindupDuration - elapsed
	case StateWindupBetween:
		return DoubleSlamPause - elapsed
	case StateSlam:
		return SlamDuration - elapsed
	case StateVulnerable:
		return b.phaseManager.GetVulnerableDuration() - elapsed
	default:
		return 0
	}
}
