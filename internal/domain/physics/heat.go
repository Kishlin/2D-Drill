package physics

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/entities"
)

// CalculateTemperature returns the temperature in °C at the given Y position
func CalculateTemperature(playerY float32) float32 {
	depthBelowGround := playerY - float32(GroundLevelY)
	if depthBelowGround <= 0 {
		return float32(BaseTemperature) // At or above ground level
	}

	maxDepth := float32(MaxUndergroundY) - float32(GroundLevelY)
	normalizedDepth := depthBelowGround / maxDepth
	temperature := float32(BaseTemperature) +
		normalizedDepth*(float32(MaxTemperature)-float32(BaseTemperature))

	return temperature
}

// ApplyHeatDamage calculates and applies damage based on depth-based temperature
func ApplyHeatDamage(player *entities.Player, dt float32) {
	temperature := CalculateTemperature(player.AABB.Y)

	excessHeat := temperature - player.HeatShield.HeatResistance()
	if excessHeat <= 0 {
		return // Player is within safe temperature range
	}

	// damage = baseDPS * (excessHeat / divisor)^exponent * dt
	damagePerSecond := float32(HeatDamageBaseDPS) *
		float32(math.Pow(float64(excessHeat/float32(HeatDamageDivisor)), float64(HeatDamageExponent)))

	damage := damagePerSecond * dt

	player.DealDamage(damage)
}
