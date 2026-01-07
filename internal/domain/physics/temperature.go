package physics

// CalculateTemperature returns the temperature in °C at the given Y position
func CalculateTemperature(playerY float32) float32 {
	depthBelowGround := playerY - float32(GroundLevelY)
	if depthBelowGround <= 0 {
		return float32(BaseTemperature) // At or above ground level
	}

	maxDepth := float32(MaxUndergroundY) - float32(GroundLevelY)
	normalizedDepth := depthBelowGround / maxDepth
	temperature := float32(BaseTemperature) +
		normalizedDepth * (float32(MaxTemperature) - float32(BaseTemperature))

	return temperature
}
