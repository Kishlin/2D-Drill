package entities

import (
	"github.com/Kishlin/drill-game/internal/domain/types"
)

const (
	PlayerWidth  = 64.0
	PlayerHeight = 64.0
)

type Player struct {
	Position types.Vec2 // World position (pixels)
	Velocity types.Vec2 // Pixels per second
	OnGround bool       // Collision state (used by physics system for ground collision)
}

func NewPlayer(startX, startY float32) *Player {
	return &Player{
		Position: types.NewVec2(startX, startY),
		Velocity: types.Zero(),
		OnGround: false,
	}
}

func (p *Player) GetPositionVec() *types.Vec2 {
	return &p.Position
}

func (p *Player) GetVelocityVec() *types.Vec2 {
	return &p.Velocity
}

func (p *Player) IsOnGround() bool {
	return p.OnGround
}

func (p *Player) SetPosition(pos types.Vec2) {
	p.Position = pos
}

func (p *Player) SetVelocity(vel types.Vec2) {
	p.Velocity = vel
}

func (p *Player) SetOnGround(onGround bool) {
	p.OnGround = onGround
}

func (p *Player) GetHeight() float32 {
	return PlayerHeight
}
