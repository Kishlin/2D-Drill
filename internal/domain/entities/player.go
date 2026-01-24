package entities

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/upgrades"
)

const (
	PlayerWidth  = 54.0
	PlayerHeight = 54.0
)

type Player struct {
	AABB          types.AABB     // Position and dimensions
	Velocity      types.Vec2     // Pixels per second
	OnGround      bool           // Collision state
	IsDrilling    bool           // Drilling animation state
	InShop        bool           // Shop UI is open, pauses gameplay
	OreInventory  map[string]int // Ore counts keyed by ore ID (e.g., "copper", "gold")
	ItemInventory [5]int         // Item counts indexed by ItemType
	Money         int            // Player's currency from selling ores
	Fuel          float32        // Current fuel in liters
	HP            float32        // Current hit points
	SpawnX        float32        // Spawn position X for teleport
	SpawnY        float32        // Spawn position Y for teleport
	engine        upgrades.Engine
	hull          upgrades.Hull
	fuelTank      upgrades.FuelTank
	cargoHold     upgrades.CargoHold
	heatShield    upgrades.HeatShield
	drill         upgrades.Drill
}

func NewPlayerFromConfig(startX, startY float32, playerCfg config.PlayerConfig, upgradeCfg config.UpgradeConfig) *Player {
	// Get starting upgrade tiers
	engineTier := upgradeCfg.Engines[playerCfg.StartingUpgrades.Engine]
	hullTier := upgradeCfg.Hulls[playerCfg.StartingUpgrades.Hull]
	fuelTankTier := upgradeCfg.FuelTanks[playerCfg.StartingUpgrades.FuelTank]
	cargoTier := upgradeCfg.CargoHolds[playerCfg.StartingUpgrades.CargoHold]
	heatTier := upgradeCfg.HeatShields[playerCfg.StartingUpgrades.HeatShield]
	drillTier := upgradeCfg.Drills[playerCfg.StartingUpgrades.Drill]

	// Create components from config
	engine := upgrades.NewEngineFromConfig(playerCfg.StartingUpgrades.Engine, engineTier.Name, engineTier.Stats)
	hull := upgrades.NewHullFromConfig(playerCfg.StartingUpgrades.Hull, hullTier.Name, hullTier.Stats)
	fuelTank := upgrades.NewFuelTankFromConfig(playerCfg.StartingUpgrades.FuelTank, fuelTankTier.Name, fuelTankTier.Stats)
	cargoHold := upgrades.NewCargoHoldFromConfig(playerCfg.StartingUpgrades.CargoHold, cargoTier.Name, cargoTier.Stats)
	heatShield := upgrades.NewHeatShieldFromConfig(playerCfg.StartingUpgrades.HeatShield, heatTier.Name, heatTier.Stats)
	drill := upgrades.NewDrillFromConfig(playerCfg.StartingUpgrades.Drill, drillTier.Name, drillTier.Stats)

	return &Player{
		AABB:          types.NewAABB(startX, startY, PlayerWidth, PlayerHeight),
		Velocity:      types.Zero(),
		OnGround:      false,
		OreInventory:  make(map[string]int),
		ItemInventory: playerCfg.StartingItems,
		Fuel:          fuelTank.Capacity(),
		HP:            hull.MaxHP(),
		SpawnX:        startX,
		SpawnY:        startY,
		engine:        engine,
		hull:          hull,
		fuelTank:      fuelTank,
		cargoHold:     cargoHold,
		heatShield:    heatShield,
		drill:         drill,
		Money:         playerCfg.StartingMoney,
	}
}

// Movement stats (from Engine)
func (p *Player) MaxSpeed() float32        { return p.engine.MaxSpeed() }
func (p *Player) Acceleration() float32    { return p.engine.Acceleration() }
func (p *Player) FlyAcceleration() float32 { return p.engine.FlyAcceleration() }
func (p *Player) MaxUpwardSpeed() float32  { return p.engine.MaxUpwardSpeed() }

// Defense stats (from Hull)
func (p *Player) MaxHP() float32 { return p.hull.MaxHP() }

// Resource stats
func (p *Player) FuelCapacity() float32 { return p.fuelTank.Capacity() }
func (p *Player) CargoCapacity() int    { return p.cargoHold.Capacity() }

// Heat/Environment stats
func (p *Player) HeatResistance() float32 { return p.heatShield.HeatResistance() }

// Drilling stats
func (p *Player) DrillSpeed() float32 { return p.drill.DrillSpeed() }

// Upgrade management
func (p *Player) GetUpgrade(t upgrades.UpgradeType) upgrades.Upgrade {
	switch t {
	case upgrades.TypeEngine:
		return p.engine
	case upgrades.TypeHull:
		return p.hull
	case upgrades.TypeFuelTank:
		return p.fuelTank
	case upgrades.TypeCargoHold:
		return p.cargoHold
	case upgrades.TypeHeatShield:
		return p.heatShield
	case upgrades.TypeDrill:
		return p.drill
	}
	return nil
}

func (p *Player) SetUpgrade(u upgrades.Upgrade) {
	switch u.Type() {
	case upgrades.TypeEngine:
		p.engine = u.(upgrades.Engine)
	case upgrades.TypeHull:
		p.hull = u.(upgrades.Hull)
	case upgrades.TypeFuelTank:
		p.fuelTank = u.(upgrades.FuelTank)
	case upgrades.TypeCargoHold:
		p.cargoHold = u.(upgrades.CargoHold)
	case upgrades.TypeHeatShield:
		p.heatShield = u.(upgrades.HeatShield)
	case upgrades.TypeDrill:
		p.drill = u.(upgrades.Drill)
	}
}

func (p *Player) GetUpgradeTier(t upgrades.UpgradeType) int {
	if u := p.GetUpgrade(t); u != nil {
		return u.Tier()
	}
	return 0
}

// Purchase methods

func (p *Player) CanAfford(cost int) bool {
	return p.Money >= cost
}

// Refuel fills the tank if player can afford it, returns success
func (p *Player) Refuel() bool {
	fuelCapacity := p.fuelTank.Capacity()
	litersNeeded := fuelCapacity - p.Fuel
	cost := int(math.Ceil(float64(litersNeeded)))

	if p.CanAfford(cost) == false {
		return false
	}

	p.Money -= cost
	p.Fuel = fuelCapacity
	return true
}

// Heal restores HP to max if player can afford it, returns success
func (p *Player) Heal() bool {
	maxHP := p.hull.MaxHP()
	hpNeeded := maxHP - p.HP

	if hpNeeded <= 0 {
		return true // Already full
	}

	cost := int(math.Ceil(float64(hpNeeded) * 2.0))

	if p.CanAfford(cost) == false {
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

// AddOreByID increments ore count for given ore ID if capacity allows
// Returns true if ore was added, false if cargo is full
func (p *Player) AddOreByID(oreID string) bool {
	if oreID == "" {
		return false
	}
	if p.GetTotalOreCount() >= p.cargoHold.Capacity() {
		return false // Cargo full
	}
	p.OreInventory[oreID]++
	return true
}

func (p *Player) AddItem(itemType ItemType) bool {
	if itemType < 0 || itemType >= 5 {
		return false
	}
	p.ItemInventory[itemType]++
	return true
}

func (p *Player) UseItem(itemType ItemType) bool {
	if itemType < 0 || itemType >= 5 {
		return false
	}
	if p.ItemInventory[itemType] <= 0 {
		return false
	}
	p.ItemInventory[itemType]--
	return true
}
