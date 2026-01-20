package attacks

import (
	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// Attack represents a boss attack pattern
type Attack interface {
	// Update updates the attack state and spawns projectiles if ready
	// Returns new projectiles spawned this frame
	Update(bossAABB, playerAABB types.AABB, dt float32) []*bosses.Projectile

	// IsReady returns true if the attack can be executed
	IsReady() bool

	// GetCooldown returns the current cooldown remaining
	GetCooldown() float32

	// Reset resets the attack cooldown
	Reset()
}
