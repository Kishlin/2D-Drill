package systems

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// Heat system constants
const (
	BaseTemperature = 15.0  // Temperature at ground level (°C)
	MaxTemperature  = 350.0 // Temperature at max depth (°C)

	// Heat damage constants
	HeatDamageBaseDPS  = 0.5  // Base damage per second
	HeatDamageDivisor  = 10.0 // Scaling factor for excess heat
	HeatDamageExponent = 1.5  // Exponential scaling factor
)

// CalculateTemperature returns the temperature in °C at the given Y position
func CalculateTemperature(playerY, groundLevel, worldHeight float32) float32 {
	depthBelowGround := playerY - groundLevel
	if depthBelowGround <= 0 {
		return float32(BaseTemperature) // At or above ground level
	}

	maxDepth := worldHeight - groundLevel
	normalizedDepth := depthBelowGround / maxDepth
	temperature := float32(BaseTemperature) +
		normalizedDepth*(float32(MaxTemperature)-float32(BaseTemperature))

	return temperature
}

// UpdateHeat calculates and applies damage based on depth-based temperature
func UpdateHeat(player *entities.Player, w *world.World, dt float32) {
	temperature := CalculateTemperature(player.AABB.Y, w.GroundLevel, w.Height)

	excessHeat := temperature - player.HeatResistance()
	if excessHeat <= 0 {
		return // Player is within safe temperature range
	}

	// damage = baseDPS * (excessHeat / divisor)^exponent * dt
	damagePerSecond := float32(HeatDamageBaseDPS) *
		float32(math.Pow(float64(excessHeat/float32(HeatDamageDivisor)), float64(HeatDamageExponent)))

	damage := damagePerSecond * dt

	player.DealDamage(damage)
}
