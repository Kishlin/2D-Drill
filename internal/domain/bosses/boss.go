package bosses

import (
	"github.com/Kishlin/drill-game/internal/domain/components"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

type Boss interface {
	// Lifecycle
	Update(player *entities.Player, dt float32) []projectiles.SpawnRequest
	Activate()
	Deactivate()
	IsActive() bool
	IsDefeated() bool

	// Health
	GetHP() float32
	GetMaxHP() float32
	GetDamageable() *components.Damageable

	// Position (origin for box offsets)
	GetPosition() types.Vec2

	// Three box types (state-dependent)
	GetCollisionBoxes() []CollisionBox // Blocks player movement
	GetHitboxes() []Hitbox             // Damages player
	GetHurtboxes() []Hurtbox           // Receives damage (empty = invulnerable)

	// Damage (only works if hurtbox exists)
	TakeDamageAt(hurtboxID string, baseDamage float32) float32
}

// AOEInfo contains information about an active AOE effect for rendering
// Used by boss-specific renderers that type-assert to concrete boss types
type AOEInfo struct {
	Position    types.Vec2
	Radius      float32
	IsTelegraph bool // Warning phase
	IsDamaging  bool // Damage phase
	StateTimer  float32
}

