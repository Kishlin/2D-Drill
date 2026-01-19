package bosses

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// Boss is the interface all bosses must implement
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

	// TakeDamage applies damage to the boss
	TakeDamage(damage float32)
}
