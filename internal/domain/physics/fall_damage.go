package physics

import "github.com/Kishlin/drill-game/internal/domain/entities"

// ApplyFallDamage calculates and applies damage based on fall velocity.
// ySpeed is positive when falling downward (screen coordinates).
func ApplyFallDamage(player *entities.Player, ySpeed float32) {
	if ySpeed < FallDamageThreshold {
		return
	}

	damage := (ySpeed - FallDamageThreshold) / FallDamageDivisor

	player.DealDamage(damage)
}
