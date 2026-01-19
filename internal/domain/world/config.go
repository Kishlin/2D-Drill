package world

import "github.com/Kishlin/drill-game/internal/domain/config"

// Type aliases for backward compatibility
// These allow existing code to continue using world.WorldConfig, world.PlayerSpawn, etc.
type (
	WorldConfig    = config.WorldConfig
	PlayerSpawn    = config.PlayerSpawn
	BuildingLayout = config.BuildingLayout
)

// NewWorldConfigForTesting creates a config suitable for tests with minimal building positions
// This avoids needing to specify all building positions in test code
func NewWorldConfigForTesting(width, height, groundLevel float32, seed int64) *WorldConfig {
	return config.NewWorldConfigForTesting(width, height, groundLevel, seed)
}
