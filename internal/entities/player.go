package entities

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	PlayerWidth  = 64.0
	PlayerHeight = 64.0
)

var (
	PlayerColor = rl.Red
)

type Player struct {
	Position rl.Vector2 // World position (pixels)
	Velocity rl.Vector2 // Pixels per second
	OnGround bool       // Collision state (used by physics system for ground collision)
}

func NewPlayer(startX, startY float32) *Player {
	return &Player{
		Position: rl.Vector2{X: startX, Y: startY},
		Velocity: rl.Vector2{X: 0, Y: 0},
		OnGround: false,
	}
}

func (p *Player) Render() {
	rl.DrawRectangleV(p.Position, rl.Vector2{X: PlayerWidth, Y: PlayerHeight}, PlayerColor)
}

func (p *Player) GetVelocity() *rl.Vector2 {
	return &p.Velocity
}

func (p *Player) IsOnGround() bool {
	return p.OnGround
}

func (p *Player) GetPosition() *rl.Vector2 {
	return &p.Position
}

func (p *Player) SetOnGround(onGround bool) {
	p.OnGround = onGround
}

func (p *Player) GetHeight() float32 {
	return PlayerHeight
}
