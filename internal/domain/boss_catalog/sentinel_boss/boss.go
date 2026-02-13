package sentinel_boss

import (
	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/bosses/attacks"
	"github.com/Kishlin/drill-game/internal/domain/bosses/movement"
	"github.com/Kishlin/drill-game/internal/domain/bosses/phases"
	"github.com/Kishlin/drill-game/internal/domain/bosses/statemachine"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// phaseConfig holds SentinelBoss-specific parameters that vary per phase.
type phaseConfig struct {
	MovementSpeed      float32
	ProjectileCooldown float32
	ChargeCooldown     float32 // 0 = no charge in this phase
	LaserCooldown      float32 // 0 = no laser in this phase
	StunDuration       float32
	UseHoming          bool // false = Sinusoidal, true = Homing
}

func init() {
	bosses.Register("sentinel_boss", func(roomStartY, worldWidth float32) bosses.Boss {
		return New(roomStartY, worldWidth)
	})
}

const (
	MaxHP         = 150.0
	Width         = 80.0
	Height        = 120.0
	ContactDamage = 15.0 // Damage per second on contact

	// Hover offset: boss floats above the floor
	HoverOffset = 200.0

	// Animation timings
	ChargeWindupDuration = 0.8
	ChargeMaxDuration    = 1.5
	ChargeSpeed          = 400.0
	LaserAimDuration     = 1.0
	LaserFireDuration    = 0.5
	LaserDamagePerSec    = 25.0
	LaserBeamWidth       = 30.0

	// Projectile parameters
	ProjectileCount     = 1
	ProjectileSpeed     = 180.0
	ProjectileSize      = 14.0
	ProjectileDamage    = 8.0
	SinusoidalAmplitude = 20.0
	SinusoidalFrequency = 5.0

	// Bobbing parameters
	BobAmplitude = 8.0
	BobFrequency = 3.0
)

// HP thresholds for the phase manager
var phaseThresholds = []phases.Config{
	{HPThreshold: 0.60}, // Phase 1: 100% - 60% HP
	{HPThreshold: 0.30}, // Phase 2: 60% - 30% HP
	{HPThreshold: 0.0},  // Phase 3: 30% - 0% HP
}

// Boss-specific parameters per phase
var phaseConfigs = []phaseConfig{
	// Phase 1: Sinusoidal projectiles, no charge, no laser, always vulnerable
	{MovementSpeed: 60, ProjectileCooldown: 2.5, ChargeCooldown: 0, LaserCooldown: 0, StunDuration: 0, UseHoming: false},
	// Phase 2: Sinusoidal projectiles, charge every 8s, no laser, vulnerable after charge stun (3s)
	{MovementSpeed: 80, ProjectileCooldown: 2.0, ChargeCooldown: 8.0, LaserCooldown: 0, StunDuration: 3.0, UseHoming: false},
	// Phase 3: Homing projectiles, charge every 6s, laser every 10s, vulnerable after charge stun (2s)
	{MovementSpeed: 100, ProjectileCooldown: 1.5, ChargeCooldown: 6.0, LaserCooldown: 10.0, StunDuration: 2.0, UseHoming: true},
}

// LaserInfo contains information about an active laser for rendering
type LaserInfo struct {
	StartX, StartY float32 // Laser origin (boss center)
	EndX, EndY     float32 // Laser target direction endpoint
	Width          float32 // Beam width
	IsAiming       bool    // Telegraph phase
	IsFiring       bool    // Damage phase
}

type SentinelBoss struct {
	*bosses.BaseBoss

	// Boss-specific components
	movement         *movement.Hovering
	projectileAttack *attacks.ProjectileAttack
	worldWidth       float32
	floorY           float32
	hoverY           float32

	// Charge state
	chargeCooldown float32
	chargeTarget   types.Vec2 // Target position locked during windup

	// Laser state
	laserCooldown float32
	laserStartX   float32
	laserStartY   float32
	laserEndX     float32
	laserEndY     float32

	// Player position tracking (for homing projectiles)
	playerPos *types.Vec2
}

func New(roomStartY, worldWidth float32) *SentinelBoss {
	floorY := roomStartY + 680 - Height
	hoverY := floorY - HoverOffset
	centerX := (worldWidth - Width) / 2

	// Create hovering movement behavior
	moveCfg := movement.HoveringConfig{
		Speed:        phaseConfigs[0].MovementSpeed,
		MinX:         0,
		MaxX:         worldWidth,
		HoverY:       hoverY,
		BossWidth:    Width,
		BobAmplitude: BobAmplitude,
		BobFrequency: BobFrequency,
	}
	hoveringMovement := movement.NewHovering(moveCfg)

	// Create projectile attack with Sinusoidal movement
	projAttack := newProjectileAttack(phaseConfigs[0], nil)

	// Create base boss
	baseBoss := bosses.NewBaseBoss(bosses.BaseBossConfig{
		Position: types.NewVec2(centerX, hoverY),
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

	b := &SentinelBoss{
		BaseBoss:         baseBoss,
		movement:         hoveringMovement,
		projectileAttack: projAttack,
		worldWidth:       worldWidth,
		floorY:           floorY,
		hoverY:           hoverY,
		chargeCooldown:   phaseConfigs[0].ChargeCooldown,
		laserCooldown:    phaseConfigs[0].LaserCooldown,
		playerPos:        &types.Vec2{},
	}

	// Configure self reference for virtual dispatch and handlers
	b.Self = b
	b.PhaseChangeHandler = b
	b.DamageReactionHandler = b

	// Build state machine
	b.SetStateMachine(statemachine.NewStateMachine(b.buildStates(), StateHover))

	return b
}

func newProjectileAttack(cfg phaseConfig, playerPos *types.Vec2) *attacks.ProjectileAttack {
	var factory attacks.MovementFactory

	if cfg.UseHoming {
		factory = func(velocity types.Vec2) projectiles.Movement {
			return projectiles.Homing{
				Speed:  velocity.Magnitude(),
				Target: playerPos,
			}
		}
	} else {
		factory = func(velocity types.Vec2) projectiles.Movement {
			return projectiles.NewSinusoidal(velocity, SinusoidalAmplitude, SinusoidalFrequency)
		}
	}

	return attacks.NewProjectileAttack(attacks.ProjectileAttackConfig{
		Cooldown:        cfg.ProjectileCooldown,
		ProjectileCount: ProjectileCount,
		ProjectileSpeed: ProjectileSpeed,
		ProjectileSize:  ProjectileSize,
		Damage:          ProjectileDamage,
		MovementFactory: factory,
	})
}

// OnPhaseChange implements PhaseChangeHandler
func (b *SentinelBoss) OnPhaseChange(phaseIndex int) {
	cfg := phaseConfigs[phaseIndex]

	// Update movement speed
	b.movement.SetSpeed(cfg.MovementSpeed)

	// Recreate projectile attack with new cooldown and movement type
	b.projectileAttack = newProjectileAttack(cfg, b.playerPos)

	// Reset charge/laser cooldowns for new phase
	if b.StateMachine.CurrentState() == StateHover {
		b.chargeCooldown = cfg.ChargeCooldown
		b.laserCooldown = cfg.LaserCooldown
	}
}

// OnDamageReceived implements DamageReactionHandler
func (b *SentinelBoss) OnDamageReceived(_ string, _ float32) {
	// End stun on damage if in stunned state
	if b.StateMachine.CurrentState() == StateStunned {
		b.StateMachine.TransitionTo(StateHover, &statemachine.StateContext{})
	}
}

// buildStates creates the state machine states
func (b *SentinelBoss) buildStates() map[statemachine.StateID]*statemachine.State {
	return map[statemachine.StateID]*statemachine.State{
		StateHover: {
			ID: StateHover,
			OnEnter: func(ctx *statemachine.StateContext) {
				b.movement.Resume()
			},
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				b.Position = b.movement.Update(b.Position, ctx.Dt)
				spawnRequests := b.updateProjectileAttack(ctx.Dt)

				// Check if it should start charge
				if b.hasCharge() {
					b.chargeCooldown -= ctx.Dt
					if b.chargeCooldown <= 0 {
						return statemachine.StateResult{
							NextState:     StateChargeWindup,
							SpawnRequests: spawnRequests,
						}
					}
				}

				// Check if it should start laser
				if b.hasLaser() {
					b.laserCooldown -= ctx.Dt
					if b.laserCooldown <= 0 {
						return statemachine.StateResult{
							NextState:     StateLaserAim,
							SpawnRequests: spawnRequests,
						}
					}
				}

				return statemachine.StateResult{NextState: statemachine.StateIDNone, SpawnRequests: spawnRequests}
			},
		},

		StateChargeWindup: {
			ID: StateChargeWindup,
			OnEnter: func(ctx *statemachine.StateContext) {
				b.movement.Pause()
				// Lock target position (player's current position)
				if b.CurrentPlayer != nil {
					b.chargeTarget = types.NewVec2(
						b.CurrentPlayer.AABB.X+b.CurrentPlayer.AABB.Width/2-Width/2,
						b.CurrentPlayer.AABB.Y+b.CurrentPlayer.AABB.Height/2-Height/2,
					)
				}
			},
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				// Still bob during windup
				b.Position = b.movement.Update(b.Position, ctx.Dt)

				if ctx.Elapsed >= ChargeWindupDuration {
					return statemachine.StateResult{NextState: StateCharge}
				}
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			},
		},

		StateCharge: {
			ID: StateCharge,
			OnEnter: func(ctx *statemachine.StateContext) {
				// Movement is already paused from windup
			},
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				// Move toward charge target
				dx := b.chargeTarget.X - b.Position.X
				dy := b.chargeTarget.Y - b.Position.Y
				dir := types.NewVec2(dx, dy)
				dist := dir.Magnitude()

				if dist < 10 || ctx.Elapsed >= ChargeMaxDuration {
					return statemachine.StateResult{NextState: StateStunned}
				}

				normalized := dir.Normalize()
				moveAmount := ChargeSpeed * ctx.Dt
				if moveAmount > dist {
					moveAmount = dist
				}
				b.Position.X += normalized.X * moveAmount
				b.Position.Y += normalized.Y * moveAmount

				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			},
		},

		StateStunned: {
			ID: StateStunned,
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				if ctx.Elapsed >= b.currentPhaseConfig().StunDuration {
					return statemachine.StateResult{NextState: StateHover}
				}
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			},
			OnExit: func(ctx *statemachine.StateContext) {
				// Only reset charge cooldown. Laser cooldown accumulates
				// across charge cycles and resets only when laser fires.
				b.chargeCooldown = b.currentPhaseConfig().ChargeCooldown
			},
		},

		StateLaserAim: {
			ID: StateLaserAim,
			OnEnter: func(ctx *statemachine.StateContext) {
				b.movement.Pause()
				b.updateLaserTarget()
			},
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				// Still bob during aim
				b.Position = b.movement.Update(b.Position, ctx.Dt)
				// Update laser origin to track boss position
				b.laserStartX = b.Position.X + Width/2
				b.laserStartY = b.Position.Y + Height/2

				if ctx.Elapsed >= LaserAimDuration {
					return statemachine.StateResult{NextState: StateLaser}
				}
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			},
		},

		StateLaser: {
			ID: StateLaser,
			OnEnter: func(ctx *statemachine.StateContext) {
				// Lock the laser direction from aim phase
				b.laserStartX = b.Position.X + Width/2
				b.laserStartY = b.Position.Y + Height/2
			},
			OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				// Still bob during laser
				b.Position = b.movement.Update(b.Position, ctx.Dt)
				b.laserStartX = b.Position.X + Width/2
				b.laserStartY = b.Position.Y + Height/2

				// Deal laser damage
				b.dealLaserDamage(ctx.Dt)

				if ctx.Elapsed >= LaserFireDuration {
					return statemachine.StateResult{NextState: StateHover}
				}
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			},
			OnExit: func(ctx *statemachine.StateContext) {
				b.laserCooldown = b.currentPhaseConfig().LaserCooldown
			},
		},
	}
}

// updateLaserTarget calculates laser endpoint toward the player
func (b *SentinelBoss) updateLaserTarget() {
	b.laserStartX = b.Position.X + Width/2
	b.laserStartY = b.Position.Y + Height/2

	if b.CurrentPlayer == nil {
		b.laserEndX = b.laserStartX
		b.laserEndY = b.laserStartY + 500
		return
	}

	playerCenterX := b.CurrentPlayer.AABB.X + b.CurrentPlayer.AABB.Width/2
	playerCenterY := b.CurrentPlayer.AABB.Y + b.CurrentPlayer.AABB.Height/2

	dx := playerCenterX - b.laserStartX
	dy := playerCenterY - b.laserStartY
	dir := types.NewVec2(dx, dy).Normalize()

	// Extend the laser far beyond the player
	b.laserEndX = b.laserStartX + dir.X*1500
	b.laserEndY = b.laserStartY + dir.Y*1500
}

// dealLaserDamage applies damage to player if within laser beam
func (b *SentinelBoss) dealLaserDamage(dt float32) {
	if b.CurrentPlayer == nil {
		return
	}

	playerCenterX := b.CurrentPlayer.AABB.X + b.CurrentPlayer.AABB.Width/2
	playerCenterY := b.CurrentPlayer.AABB.Y + b.CurrentPlayer.AABB.Height/2

	// Calculate distance from player center to laser line segment
	dist := distanceToLineSegment(
		playerCenterX, playerCenterY,
		b.laserStartX, b.laserStartY,
		b.laserEndX, b.laserEndY,
	)

	if dist < LaserBeamWidth/2 {
		b.CurrentPlayer.DealDamage(LaserDamagePerSec * dt)
	}
}

// distanceToLineSegment calculates the distance from point (px, py) to line segment (x1,y1)-(x2,y2)
func distanceToLineSegment(px, py, x1, y1, x2, y2 float32) float32 {
	dx := x2 - x1
	dy := y2 - y1
	lengthSq := dx*dx + dy*dy

	if lengthSq == 0 {
		// Degenerate segment (start == end)
		ddx := px - x1
		ddy := py - y1
		return sqrt(ddx*ddx + ddy*ddy)
	}

	// Project point onto line, clamped to [0, 1]
	t := ((px-x1)*dx + (py-y1)*dy) / lengthSq
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	// Closest point on segment
	closestX := x1 + t*dx
	closestY := y1 + t*dy

	ddx := px - closestX
	ddy := py - closestY
	return sqrt(ddx*ddx + ddy*ddy)
}

// sqrt approximation using Newton's method (avoid math import in domain)
func sqrt(x float32) float32 {
	if x <= 0 {
		return 0
	}
	guess := x
	for i := 0; i < 10; i++ {
		guess = (guess + x/guess) / 2
	}
	return guess
}

// updateProjectileAttack handles projectile spawning during hover
func (b *SentinelBoss) updateProjectileAttack(dt float32) []projectiles.SpawnRequest {
	if b.CurrentPlayer == nil {
		return nil
	}
	bossAABB := types.AABB{X: b.Position.X, Y: b.Position.Y, Width: Width, Height: Height}
	return b.projectileAttack.Update(bossAABB, b.CurrentPlayer.AABB, dt)
}

// currentPhaseConfig returns the boss-specific config for the current phase
func (b *SentinelBoss) currentPhaseConfig() phaseConfig {
	return phaseConfigs[b.PhaseManager.GetCurrentPhase()]
}

// hasCharge returns true if the current phase has charge attacks
func (b *SentinelBoss) hasCharge() bool {
	return b.currentPhaseConfig().ChargeCooldown > 0
}

// hasLaser returns true if the current phase has laser attacks
func (b *SentinelBoss) hasLaser() bool {
	return b.currentPhaseConfig().LaserCooldown > 0
}

func (b *SentinelBoss) Update(player *entities.Player, dt float32) []projectiles.SpawnRequest {
	// Update player position tracking for homing projectiles
	if player != nil {
		b.playerPos.X = player.AABB.X + player.AABB.Width/2
		b.playerPos.Y = player.AABB.Y + player.AABB.Height/2
	}

	return b.BaseUpdate(player, dt)
}

// GetHurtboxes returns hurtboxes based on phase and state.
// Phase 1 (index 0): always has hurtboxes
// Phase 2+ (index 1+): only has hurtboxes during StateStunned
func (b *SentinelBoss) GetHurtboxes() []bosses.Hurtbox {
	if b.PhaseManager.GetCurrentPhase() == 0 || b.StateMachine.CurrentState() == StateStunned {
		return b.BoxSet.Hurtboxes
	}
	return []bosses.Hurtbox{}
}

// IsVulnerable returns true if the boss can currently take damage (for rendering)
func (b *SentinelBoss) IsVulnerable() bool {
	return len(b.GetHurtboxes()) > 0
}

// GetLaserInfo returns laser rendering information, or nil if no laser active
func (b *SentinelBoss) GetLaserInfo() *LaserInfo {
	currentState := b.StateMachine.CurrentState()

	switch currentState {
	case StateLaserAim:
		return &LaserInfo{
			StartX:   b.laserStartX,
			StartY:   b.laserStartY,
			EndX:     b.laserEndX,
			EndY:     b.laserEndY,
			Width:    LaserBeamWidth,
			IsAiming: true,
			IsFiring: false,
		}
	case StateLaser:
		return &LaserInfo{
			StartX:   b.laserStartX,
			StartY:   b.laserStartY,
			EndX:     b.laserEndX,
			EndY:     b.laserEndY,
			Width:    LaserBeamWidth,
			IsAiming: false,
			IsFiring: true,
		}
	default:
		return nil
	}
}

// GetChargeTarget returns the charge target position (for rendering)
func (b *SentinelBoss) GetChargeTarget() types.Vec2 {
	return b.chargeTarget
}

// GetVulnerableTimer returns remaining vulnerability duration (for UI rendering)
// Returns -1 in phase 1 (always vulnerable), 0 when not vulnerable, or remaining time
func (b *SentinelBoss) GetVulnerableTimer() float32 {
	if b.PhaseManager.GetCurrentPhase() == 0 {
		return -1
	}
	if b.StateMachine.CurrentState() == StateStunned {
		return b.currentPhaseConfig().StunDuration - b.StateMachine.Elapsed()
	}
	return 0
}

// GetStateTimer returns the remaining time in current state
func (b *SentinelBoss) GetStateTimer() float32 {
	currentState := b.StateMachine.CurrentState()
	elapsed := b.StateMachine.Elapsed()

	switch currentState {
	case StateChargeWindup:
		return ChargeWindupDuration - elapsed
	case StateCharge:
		return ChargeMaxDuration - elapsed
	case StateStunned:
		return b.currentPhaseConfig().StunDuration - elapsed
	case StateLaserAim:
		return LaserAimDuration - elapsed
	case StateLaser:
		return LaserFireDuration - elapsed
	default:
		return 0
	}
}
