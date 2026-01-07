package entities

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/types"
)

const (
	PlayerWidth  = 54.0
	PlayerHeight = 54.0
)

type Player struct {
	AABB         types.AABB // Position and dimensions
	Velocity     types.Vec2 // Pixels per second
	OnGround     bool       // Collision state
	OreInventory [7]int     // Ore counts indexed by OreType
	Money        int        // Player's currency from selling ores
	Fuel         float32    // Current fuel in liters
	HP           float32    // Current hit points
	Engine       Engine     // Engine component (exported)
	Hull         Hull       // Hull component (exported)
	FuelTank     FuelTank   // FuelTank component (exported)
}

func NewPlayer(startX, startY float32) *Player {
	engine := NewEngineBase()
	hull := NewHullBase()
	fuelTank := NewFuelTankBase()

	return &Player{
		AABB:         types.NewAABB(startX, startY, PlayerWidth, PlayerHeight),
		Velocity:     types.Zero(),
		OnGround:     false,
		OreInventory: [7]int{},
		Fuel:         fuelTank.Capacity(),
		HP:           hull.MaxHP(),
		Engine:       engine,
		Hull:         hull,
		FuelTank:     fuelTank,
		Money:        1000000,
	}
}

// Purchase methods

func (p *Player) CanAfford(cost int) bool {
	return p.Money >= cost
}

func (p *Player) BuyEngine(e Engine, cost int) {
	p.Money -= cost
	p.Engine = e
}

func (p *Player) BuyHull(h Hull, cost int) {
	p.Money -= cost
	p.Hull = h
}

func (p *Player) BuyFuelTank(ft FuelTank, cost int) {
	p.Money -= cost
	p.FuelTank = ft
}

// Refuel fills the tank if player can afford it, returns success
func (p *Player) Refuel() bool {
	fuelCapacity := p.FuelTank.Capacity()
	litersNeeded := fuelCapacity - p.Fuel
	cost := int(math.Ceil(float64(litersNeeded)))

	if !p.CanAfford(cost) {
		return false
	}

	p.Money -= cost
	p.Fuel = fuelCapacity
	return true
}

// Heal restores HP to max if player can afford it, returns success
func (p *Player) Heal() bool {
	maxHP := p.Hull.MaxHP()
	hpNeeded := maxHP - p.HP

	if hpNeeded <= 0 {
		return true // Already full
	}

	cost := int(math.Ceil(float64(hpNeeded) * 2.0))

	if !p.CanAfford(cost) {
		return false
	}

	p.Money -= cost
	p.HP = maxHP
	return true
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
