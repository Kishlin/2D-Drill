package bosses

import (
	"github.com/Kishlin/drill-game/internal/domain/bosses/statemachine"
	"github.com/Kishlin/drill-game/internal/domain/components"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// PhaseChangeHandler is called when the boss transitions to a new phase
type PhaseChangeHandler interface {
	OnPhaseChange(phaseIndex int, config PhaseConfig)
}

// DamageReactionHandler is called when the boss receives damage
type DamageReactionHandler interface {
	OnDamageReceived(hurtboxID string, damage float32)
}

// BaseBossConfig contains initialization parameters for BaseBoss
type BaseBossConfig struct {
	Position types.Vec2
	MaxHP    float32
	BoxSet   *BoxSet
	Phases   []PhaseConfig
}

// BaseBoss provides default implementations for common Boss interface methods.
// Embed this struct in concrete boss implementations to reduce boilerplate.
type BaseBoss struct {
	Position      types.Vec2
	Damageable    components.Damageable
	Active        bool
	BoxSet        *BoxSet
	StateMachine  *statemachine.StateMachine
	PhaseManager  *PhaseManager
	CurrentPlayer *entities.Player

	// Optional handlers (nil = skip)
	PhaseChangeHandler    PhaseChangeHandler
	DamageReactionHandler DamageReactionHandler
}

// NewBaseBoss creates a BaseBoss with the given configuration
func NewBaseBoss(cfg BaseBossConfig) *BaseBoss {
	return &BaseBoss{
		Position:     cfg.Position,
		Damageable:   components.NewDamageable(cfg.MaxHP, cfg.MaxHP),
		Active:       false,
		BoxSet:       cfg.BoxSet,
		PhaseManager: NewPhaseManager(cfg.MaxHP, cfg.Phases),
	}
}

// SetStateMachine sets the state machine for this boss.
// Call this after creating the boss and building the state machine.
func (b *BaseBoss) SetStateMachine(sm *statemachine.StateMachine) {
	b.StateMachine = sm
}

// Activate sets the boss as active
func (b *BaseBoss) Activate() {
	b.Active = true
}

// Deactivate sets the boss as inactive
func (b *BaseBoss) Deactivate() {
	b.Active = false
}

// IsActive returns whether the boss is currently active
func (b *BaseBoss) IsActive() bool {
	return b.Active
}

// IsDefeated returns whether the boss has been defeated
func (b *BaseBoss) IsDefeated() bool {
	return b.Damageable.IsDefeated()
}

// GetHP returns the boss's current HP
func (b *BaseBoss) GetHP() float32 {
	return b.Damageable.HP
}

// GetMaxHP returns the boss's maximum HP
func (b *BaseBoss) GetMaxHP() float32 {
	return b.Damageable.MaxHP
}

// GetDamageable returns a pointer to the boss's Damageable component
func (b *BaseBoss) GetDamageable() *components.Damageable {
	return &b.Damageable
}

// GetPosition returns the boss's current position
func (b *BaseBoss) GetPosition() types.Vec2 {
	return b.Position
}

// GetCollisionBoxes returns the boss's collision boxes
func (b *BaseBoss) GetCollisionBoxes() []CollisionBox {
	return b.BoxSet.CollisionBoxes
}

// GetHitboxes returns the boss's hitboxes
func (b *BaseBoss) GetHitboxes() []Hitbox {
	return b.BoxSet.Hitboxes
}

// GetHurtboxes returns the boss's hurtboxes.
// Override this method in concrete bosses to implement conditional vulnerability.
func (b *BaseBoss) GetHurtboxes() []Hurtbox {
	return b.BoxSet.Hurtboxes
}

// TakeDamageAt applies damage to the boss at a specific hurtbox.
// Calls DamageReactionHandler if set.
func (b *BaseBoss) TakeDamageAt(hurtboxID string, baseDamage float32) float32 {
	for _, hb := range b.BoxSet.Hurtboxes {
		if hb.ID == hurtboxID {
			actual := baseDamage * hb.DamageMultiplier
			b.Damageable.TakeDamage(actual)

			if b.DamageReactionHandler != nil {
				b.DamageReactionHandler.OnDamageReceived(hurtboxID, actual)
			}
			return actual
		}
	}
	return 0
}

// BaseUpdate performs common update logic for bosses.
// Call this from concrete boss Update() methods.
// Returns spawn requests from the state machine.
func (b *BaseBoss) BaseUpdate(player *entities.Player, dt float32) []projectiles.SpawnRequest {
	if b.Active == false || b.Damageable.IsDefeated() {
		return nil
	}

	b.CurrentPlayer = player

	// Check for phase transitions
	if b.PhaseManager.Update(b.Damageable.HP) {
		if b.PhaseChangeHandler != nil {
			phaseIndex := b.PhaseManager.GetCurrentPhase()
			config := b.PhaseManager.GetCurrentConfig()
			b.PhaseChangeHandler.OnPhaseChange(phaseIndex, config)
		}
	}

	// Update state machine
	ctx := &statemachine.StateContext{
		Player: player,
		Dt:     dt,
	}
	result := b.StateMachine.Update(ctx)

	// Update box positions
	b.BoxSet.UpdatePositions(b.Position.X, b.Position.Y)

	return result.SpawnRequests
}

// GetCurrentPhase returns the current phase number (1-indexed for display)
func (b *BaseBoss) GetCurrentPhase() int {
	return b.PhaseManager.GetCurrentPhase() + 1
}

// GetState returns the current state ID
func (b *BaseBoss) GetState() statemachine.StateID {
	return b.StateMachine.CurrentState()
}

// GetStateElapsed returns the time elapsed in the current state
func (b *BaseBoss) GetStateElapsed() float32 {
	return b.StateMachine.Elapsed()
}
