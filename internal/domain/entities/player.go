package entities

import (
	"github.com/Kishlin/drill-game/internal/domain/types"
)

const (
	PlayerWidth  = 54.0
	PlayerHeight = 54.0
)

type Player struct {
	AABB         types.AABB // Position and dimensions - direct access
	Velocity     types.Vec2 // Pixels per second - direct access
	OnGround     bool       // Collision state - direct access
	OreInventory [7]int     // Ore counts indexed by OreType
	Money        int        // Player's currency from selling ores
}

func NewPlayer(startX, startY float32) *Player {
	return &Player{
		AABB:         types.NewAABB(startX, startY, PlayerWidth, PlayerHeight),
		Velocity:     types.Zero(),
		OnGround:     false,
		OreInventory: [7]int{},
	}
}

// AddOre increments ore count for given type
func (p *Player) AddOre(oreType OreType, amount int) {
	if oreType >= 0 && oreType < 7 {
		p.OreInventory[oreType] += amount
	}
}

// SellInventory sells all ore in inventory and adds value to player's money
func (p *Player) SellInventory() {
	totalValue := CalculateInventoryValue(p.OreInventory)
	p.Money += totalValue
	p.OreInventory = [7]int{}
}
