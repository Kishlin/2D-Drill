package attacks

import (
	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// AOEState represents the current state of the AOE attack
type AOEState int

const (
	AOEStateIdle AOEState = iota
	AOEStateTelegraph
	AOEStateDamage
	AOEStateVulnerable
)

// AOEAttackConfig holds configuration for an AOE attack
type AOEAttackConfig struct {
	Cooldown           float32 // Time between attacks
	TelegraphDuration  float32 // Warning duration before damage
	DamageDuration     float32 // Duration of damage zone
	VulnerableDuration float32 // Duration of vulnerability window after attack
	Radius             float32 // Radius of the AOE
	Damage             float32 // Damage dealt
}

// AOEAttack represents a ground slam attack with telegraph and vulnerability
type AOEAttack struct {
	config        AOEAttackConfig
	cooldownTimer float32
	stateTimer    float32
	state         AOEState
	position      types.Vec2 // Center of the AOE effect
}

// NewAOEAttack creates a new AOE attack
func NewAOEAttack(cfg AOEAttackConfig) *AOEAttack {
	return &AOEAttack{
		config:        cfg,
		cooldownTimer: cfg.Cooldown, // Start on cooldown
		stateTimer:    0,
		state:         AOEStateIdle,
	}
}

// Update updates the AOE attack state
// Returns damage to apply to player (0 if no damage this frame)
func (a *AOEAttack) Update(bossAABB, playerAABB types.AABB, dt float32) []*bosses.Projectile {
	switch a.state {
	case AOEStateIdle:
		if a.cooldownTimer > 0 {
			a.cooldownTimer -= dt
		}
	case AOEStateTelegraph:
		a.stateTimer -= dt
		if a.stateTimer <= 0 {
			a.state = AOEStateDamage
			a.stateTimer = a.config.DamageDuration
		}
	case AOEStateDamage:
		a.stateTimer -= dt
		if a.stateTimer <= 0 {
			a.state = AOEStateVulnerable
			a.stateTimer = a.config.VulnerableDuration
		}
	case AOEStateVulnerable:
		a.stateTimer -= dt
		if a.stateTimer <= 0 {
			a.state = AOEStateIdle
			a.cooldownTimer = a.config.Cooldown
		}
	}

	// AOE attacks don't spawn projectiles - damage is handled separately
	return nil
}

// StartAttack begins the attack sequence (called by boss when ready)
func (a *AOEAttack) StartAttack(bossAABB types.AABB) {
	if a.state != AOEStateIdle || a.cooldownTimer > 0 {
		return
	}

	// Position AOE at boss's feet
	a.position = types.NewVec2(
		bossAABB.X+bossAABB.Width/2,
		bossAABB.Y+bossAABB.Height,
	)
	a.state = AOEStateTelegraph
	a.stateTimer = a.config.TelegraphDuration
}

// GetDamageToPlayer returns damage to apply if player is in AOE during damage phase
func (a *AOEAttack) GetDamageToPlayer(playerAABB types.AABB) float32 {
	if a.state != AOEStateDamage {
		return 0
	}

	// Check if player center is within radius
	playerCenterX := playerAABB.X + playerAABB.Width/2
	playerCenterY := playerAABB.Y + playerAABB.Height/2

	dx := playerCenterX - a.position.X
	dy := playerCenterY - a.position.Y
	distSq := dx*dx + dy*dy
	radiusSq := a.config.Radius * a.config.Radius

	if distSq <= radiusSq {
		return a.config.Damage
	}
	return 0
}

// IsReady returns true if the attack can be started
func (a *AOEAttack) IsReady() bool {
	return a.state == AOEStateIdle && a.cooldownTimer <= 0
}

// GetCooldown returns remaining cooldown
func (a *AOEAttack) GetCooldown() float32 {
	return a.cooldownTimer
}

// Reset resets the attack to idle state
func (a *AOEAttack) Reset() {
	a.state = AOEStateIdle
	a.cooldownTimer = a.config.Cooldown
	a.stateTimer = 0
}

// GetState returns the current attack state
func (a *AOEAttack) GetState() AOEState {
	return a.state
}

// IsVulnerableWindow returns true if in vulnerability window
func (a *AOEAttack) IsVulnerableWindow() bool {
	return a.state == AOEStateVulnerable
}

// IsTelegraphing returns true if showing warning
func (a *AOEAttack) IsTelegraphing() bool {
	return a.state == AOEStateTelegraph
}

// IsDamaging returns true if dealing damage
func (a *AOEAttack) IsDamaging() bool {
	return a.state == AOEStateDamage
}

// GetPosition returns the center of the AOE effect
func (a *AOEAttack) GetPosition() types.Vec2 {
	return a.position
}

// GetRadius returns the AOE radius
func (a *AOEAttack) GetRadius() float32 {
	return a.config.Radius
}

// GetStateTimer returns the remaining time in current state
func (a *AOEAttack) GetStateTimer() float32 {
	return a.stateTimer
}
