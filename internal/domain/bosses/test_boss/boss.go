package test_boss

import (
	"math/rand"

	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/bosses/attacks"
	"github.com/Kishlin/drill-game/internal/domain/bosses/movement"
	"github.com/Kishlin/drill-game/internal/domain/components"
	"github.com/Kishlin/drill-game/internal/domain/entities"
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

// BossState represents the current animation state
type BossState int

const (
	StatePatrol     BossState = iota // Moving and shooting
	StateWindup                      // Stopped, vibrating, preparing to slam
	StateSlam                        // Dealing AOE damage
	StateVulnerable                  // Immobile, can be damaged
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
	projectiles      []*bosses.Projectile
	movement         *movement.Grounded
	projectileAttack *attacks.ProjectileAttack
	phaseManager     *bosses.PhaseManager
	worldWidth       float32
	floorY           float32

	// State machine
	state        BossState
	stateTimer   float32
	slamCooldown float32
	slamCount    int // For phase 3 double slams
	maxSlams     int // 1 for phases 1-2, 1-2 for phase 3
	aoeRadius    float32
	aoeDamage    float32
	aoePosition  types.Vec2
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

	return &TestBoss{
		aabb: types.AABB{
			X:      centerX,
			Y:      floorY,
			Width:  Width,
			Height: Height,
		},
		damageable:       components.NewDamageable(MaxHP, MaxHP),
		active:           false,
		projectiles:      make([]*bosses.Projectile, 0),
		movement:         groundedMovement,
		projectileAttack: projAttack,
		phaseManager:     phaseManager,
		worldWidth:       worldWidth,
		floorY:           floorY,
		state:            StatePatrol,
		stateTimer:       0,
		slamCooldown:     phases[0].AOECooldown,
		slamCount:        0,
		maxSlams:         1,
		aoeRadius:        150.0,
		aoeDamage:        15.0,
	}
}

func (b *TestBoss) Update(player *entities.Player, dt float32) {
	if !b.active || b.damageable.IsDefeated() {
		return
	}

	// Check for phase transitions
	if b.phaseManager.Update(b.damageable.HP) {
		b.onPhaseChange()
	}

	phaseCfg := b.phaseManager.GetCurrentConfig()

	switch b.state {
	case StatePatrol:
		b.updatePatrol(player, dt, phaseCfg)
	case StateWindup:
		b.updateWindup(dt)
	case StateSlam:
		b.updateSlam(player, dt)
	case StateVulnerable:
		b.updateVulnerable(dt)
	}

	// Update existing projectiles (always, regardless of state)
	b.updateProjectiles(dt, player)
}

func (b *TestBoss) updatePatrol(player *entities.Player, dt float32, phaseCfg bosses.PhaseConfig) {
	// Move
	newPos := b.movement.Update(types.NewVec2(b.aabb.X, b.aabb.Y), dt)
	b.aabb.X = newPos.X
	b.aabb.Y = newPos.Y

	// Shoot projectiles
	newProjectiles := b.projectileAttack.Update(b.aabb, player.AABB, dt)
	if len(newProjectiles) > 0 {
		b.projectiles = append(b.projectiles, newProjectiles...)
	}

	// Check if should start slam (only in phases with AOE)
	if phaseCfg.AOECooldown > 0 {
		b.slamCooldown -= dt
		if b.slamCooldown <= 0 {
			b.startWindup()
		}
	}
}

func (b *TestBoss) startWindup() {
	b.state = StateWindup
	b.stateTimer = WindupDuration
	b.slamCount = 0

	// Determine max slams for this attack
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
}

func (b *TestBoss) updateWindup(dt float32) {
	b.stateTimer -= dt
	if b.stateTimer <= 0 {
		b.state = StateSlam
		b.stateTimer = SlamDuration
	}
}

func (b *TestBoss) updateSlam(player *entities.Player, dt float32) {
	// Deal AOE damage to player if in range
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerCenterY := player.AABB.Y + player.AABB.Height/2
	dx := playerCenterX - b.aoePosition.X
	dy := playerCenterY - b.aoePosition.Y
	distSq := dx*dx + dy*dy
	radiusSq := b.aoeRadius * b.aoeRadius

	if distSq <= radiusSq {
		player.DealDamage(b.aoeDamage * dt / SlamDuration) // Spread damage over slam duration
	}

	b.stateTimer -= dt
	if b.stateTimer <= 0 {
		b.slamCount++

		// Check if more slams to do
		if b.slamCount < b.maxSlams {
			// Brief pause before next slam
			b.state = StateWindup
			b.stateTimer = DoubleSlamPause
		} else {
			// Done slamming, enter vulnerable state
			b.state = StateVulnerable
			b.stateTimer = b.phaseManager.GetCurrentConfig().VulnerableDuration
		}
	}
}

func (b *TestBoss) updateVulnerable(dt float32) {
	b.stateTimer -= dt
	if b.stateTimer <= 0 {
		b.endVulnerability()
	}
}

func (b *TestBoss) endVulnerability() {
	b.state = StatePatrol
	b.slamCooldown = b.phaseManager.GetCurrentConfig().AOECooldown
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
	if b.state == StatePatrol {
		b.slamCooldown = phaseCfg.AOECooldown
	}
}

func (b *TestBoss) updateProjectiles(dt float32, player *entities.Player) {
	activeProjectiles := make([]*bosses.Projectile, 0, len(b.projectiles))
	for _, proj := range b.projectiles {
		if !proj.Active {
			continue
		}
		proj.Update(dt)

		if proj.AABB.X < -100 || proj.AABB.X > b.worldWidth+100 ||
			proj.AABB.Y < -100 || proj.AABB.Y > b.floorY+200 {
			proj.Deactivate()
			continue
		}

		activeProjectiles = append(activeProjectiles, proj)
	}
	b.projectiles = activeProjectiles
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
	b.projectiles = make([]*bosses.Projectile, 0)
}

func (b *TestBoss) GetProjectiles() []*bosses.Projectile {
	return b.projectiles
}

func (b *TestBoss) GetAABB() types.AABB {
	return b.aabb
}

func (b *TestBoss) GetContactDamage() float32 {
	return ContactDamage
}

func (b *TestBoss) TakeDamage(damage float32) {
	if !b.IsVulnerable() {
		return
	}

	b.damageable.TakeDamage(damage)

	// Close vulnerability window after taking damage
	if b.state == StateVulnerable {
		b.endVulnerability()
	}
}

func (b *TestBoss) IsVulnerable() bool {
	phaseCfg := b.phaseManager.GetCurrentConfig()
	return phaseCfg.AlwaysVulnerable || b.state == StateVulnerable
}

func (b *TestBoss) GetDamageable() *components.Damageable {
	return &b.damageable
}

func (b *TestBoss) GetVulnerableTimer() float32 {
	if b.phaseManager.GetCurrentConfig().AlwaysVulnerable {
		return -1
	}
	if b.state == StateVulnerable {
		return b.stateTimer
	}
	return 0
}

func (b *TestBoss) GetAOEInfo() *bosses.AOEInfo {
	switch b.state {
	case StateWindup:
		return &bosses.AOEInfo{
			Position:    b.aoePosition,
			Radius:      b.aoeRadius,
			IsTelegraph: true,
			IsDamaging:  false,
			StateTimer:  b.stateTimer,
		}
	case StateSlam:
		return &bosses.AOEInfo{
			Position:    b.aoePosition,
			Radius:      b.aoeRadius,
			IsTelegraph: false,
			IsDamaging:  true,
			StateTimer:  b.stateTimer,
		}
	default:
		return nil
	}
}

func (b *TestBoss) GetCurrentPhase() int {
	return b.phaseManager.GetCurrentPhase() + 1
}

// GetState returns the current boss state for rendering
func (b *TestBoss) GetState() BossState {
	return b.state
}

// GetStateTimer returns the remaining time in current state
func (b *TestBoss) GetStateTimer() float32 {
	return b.stateTimer
}
