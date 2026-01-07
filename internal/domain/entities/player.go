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
	IsDrilling   bool       // Drilling animation state
	OreInventory [6]int     // Ore counts indexed by OreType
	Money        int        // Player's currency from selling ores
	Fuel         float32    // Current fuel in liters
	HP           float32    // Current hit points
	Engine       Engine     // Engine component (exported)
	Hull         Hull       // Hull component (exported)
	FuelTank     FuelTank   // FuelTank component (exported)
	CargoHold    CargoHold  // CargoHold component (exported)
	HeatShield   HeatShield // HeatShield component (exported)
}

func NewPlayer(startX, startY float32) *Player {
	engine := NewEngineMk3()
	hull := NewHullMk2()
	fuelTank := NewFuelTankMk4()
	cargoHold := NewCargoHoldMk2()
	heatShield := NewHeatShieldMk1()

	return &Player{
		AABB:         types.NewAABB(startX, startY, PlayerWidth, PlayerHeight),
		Velocity:     types.Zero(),
		OnGround:     false,
		OreInventory: [6]int{2, 5, 3, 0, 1, 0},
		Fuel:         fuelTank.Capacity(),
		HP:           hull.MaxHP(),
		Engine:       engine,
		Hull:         hull,
		FuelTank:     fuelTank,
		CargoHold:    cargoHold,
		HeatShield:   heatShield,
		Money:        1854,
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

func (p *Player) BuyCargoHold(ch CargoHold, cost int) {
	p.Money -= cost
	p.CargoHold = ch
}

func (p *Player) BuyHeatShield(hs HeatShield, cost int) {
	p.Money -= cost
	p.HeatShield = hs
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

// DealDamage applies damage to player HP, clamping at zero
func (p *Player) DealDamage(damage float32) {
	p.HP -= damage
	if p.HP < 0 {
		p.HP = 0
	}
}

func (p *Player) GetTotalOreCount() int {
	total := 0
	for _, count := range p.OreInventory {
		total += count
	}
	return total
}

// AddOre increments ore count for given type if capacity allows
// Returns true if ore was added, false if cargo is full
func (p *Player) AddOre(oreType OreType) bool {
	if oreType < 0 || oreType >= 6 {
		return false
	}
	if p.GetTotalOreCount() >= p.CargoHold.Capacity() {
		return false // Cargo full
	}
	p.OreInventory[oreType]++
	return true
}

// SellInventory sells all ore in inventory and adds value to player's money
func (p *Player) SellInventory() {
	totalValue := CalculateInventoryValue(p.OreInventory)
	p.Money += totalValue
	p.OreInventory = [6]int{}
}
