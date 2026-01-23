package bosses

import (
	"github.com/Kishlin/drill-game/internal/domain/components"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

type Boss interface {
	// Update runs boss AI and updates state
	Update(player *entities.Player, dt float32)

	// GetHP returns current health points
	GetHP() float32

	// GetMaxHP returns maximum health points
	GetMaxHP() float32

	// IsDefeated returns true if boss is dead
	IsDefeated() bool

	// IsActive returns true if boss is currently active in the fight
	IsActive() bool

	// Activate starts the boss fight
	Activate()

	// Deactivate pauses the boss fight
	Deactivate()

	// GetProjectiles returns all projectiles the boss has spawned
	// Returns empty slice if boss doesn't use projectiles
	GetProjectiles() []*Projectile
}

// PhysicalBoss is a boss that has a physical AABB and can be hit by bombs
type PhysicalBoss interface {
	Boss

	// GetAABB returns the boss's axis-aligned bounding box for collision detection
	GetAABB() types.AABB

	// GetDamageable returns the boss's damageable component
	GetDamageable() *components.Damageable

	// TakeDamage applies damage to the boss (delegates to Damageable)
	TakeDamage(damage float32)

	// IsVulnerable returns true if the boss can currently take damage
	IsVulnerable() bool

	// GetVulnerableTimer returns remaining vulnerability duration (for UI)
	GetVulnerableTimer() float32

	// GetContactDamage returns damage dealt per second on player contact
	// Returns 0 if boss doesn't deal contact damage
	GetContactDamage() float32
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

