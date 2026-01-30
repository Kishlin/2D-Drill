package effects

import (
	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/components"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// DamageableEntity is what effects need to damage entities in area
type DamageableEntity interface {
	GetHurtboxes() []bosses.Hurtbox
	GetDamageable() *components.Damageable
	TakeDamageAt(hurtboxID string, baseDamage float32) float32
}

// EffectContext provides access to game state for effects
type EffectContext struct {
	Player      *entities.Player
	World       *world.World
	Damageables []DamageableEntity // Boss + future enemies
}

type Effect interface {
	Apply(ctx *EffectContext)
}
