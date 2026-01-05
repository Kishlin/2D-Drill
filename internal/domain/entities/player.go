package entities

import (
	"github.com/Kishlin/drill-game/internal/domain/types"
)

const (
	PlayerWidth  = 54.0
	PlayerHeight = 54.0
	FuelCapacity = 10.0 // Liters
	MaxHP        = 10.0 // Maximum hit points
)

type Player struct {
	AABB         types.AABB // Position and dimensions - direct access
	Velocity     types.Vec2 // Pixels per second - direct access
	OnGround     bool       // Collision state - direct access
	OreInventory [7]int     // Ore counts indexed by OreType
	Money        int        // Player's currency from selling ores
	Fuel         float32    // Current fuel in liters (0.0 to FuelCapacity)
	HP           float32    // Current hit points (0.0 to MaxHP)
	Upgrades     Upgrades   // Current upgrade levels
}

func NewPlayer(startX, startY float32) *Player {
	return &Player{
		AABB:         types.NewAABB(startX, startY, PlayerWidth, PlayerHeight),
		Velocity:     types.Zero(),
		OnGround:     false,
		OreInventory: [7]int{},
		Fuel:         FuelCapacity,
		HP:           MaxHP,
		Upgrades:     NewUpgrades(),
	}
}

func (p *Player) GetMaxMoveSpeed() float32 {
	return GetEngineMaxSpeed(p.Upgrades.Engine)
}

func (p *Player) GetMoveAcceleration() float32 {
	return GetEngineAcceleration(p.Upgrades.Engine)
}

func (p *Player) GetFlyAcceleration() float32 {
	return GetEngineFlyAcceleration(p.Upgrades.Engine)
}

func (p *Player) GetMaxUpwardSpeed() float32 {
	return GetEngineMaxUpwardSpeed(p.Upgrades.Engine)
}

func (p *Player) GetMaxHP() float32 {
	return GetHullMaxHP(p.Upgrades.Hull)
}

func (p *Player) GetFuelCapacity() float32 {
	return GetFuelTankCapacity(p.Upgrades.FuelTank)
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
