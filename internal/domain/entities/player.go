package entities

import (
	"github.com/Kishlin/drill-game/internal/domain/types"
)

const (
	PlayerWidth  = 54.0
	PlayerHeight = 54.0
)

type Player struct {
	AABB     types.AABB // Position and dimensions - direct access
	Velocity types.Vec2 // Pixels per second - direct access
	OnGround bool       // Collision state - direct access
}

func NewPlayer(startX, startY float32) *Player {
	return &Player{
		AABB:     types.NewAABB(startX, startY, PlayerWidth, PlayerHeight),
		Velocity: types.Zero(),
		OnGround: false,
	}
}
