package projectiles

import "github.com/Kishlin/drill-game/internal/domain/types"

// SpawnRequest contains data for spawning a projectile
type SpawnRequest struct {
	Position types.Vec2
	Size     float32
	Damage   float32
	Movement Movement
}
