package systems

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// CalculateTemperature returns the temperature in °C at the given Y position
func CalculateTemperature(playerY, groundLevel, worldHeight float32, heatCfg config.HeatConfig) float32 {
	depthBelowGround := playerY - groundLevel
	if depthBelowGround <= 0 {
		return heatCfg.BaseTemperature
	}

	maxDepth := worldHeight - groundLevel
	normalizedDepth := depthBelowGround / maxDepth
	temperature := heatCfg.BaseTemperature +
		normalizedDepth*(heatCfg.MaxTemperature-heatCfg.BaseTemperature)

	return temperature
}

// UpdateHeat calculates and applies damage based on depth-based temperature
func UpdateHeat(player *entities.Player, w *world.World, dt float32, heatCfg config.HeatConfig) {
	temperature := CalculateTemperature(player.AABB.Y, w.GroundLevel, w.Height, heatCfg)

	excessHeat := temperature - player.HeatResistance()
	if excessHeat <= 0 {
		return
	}

	damagePerSecond := heatCfg.DamageBaseDPS *
		float32(math.Pow(float64(excessHeat/heatCfg.DamageDivisor), float64(heatCfg.DamageExponent)))

	damage := damagePerSecond * dt

	player.DealDamage(damage)
}
