package entities

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/types"
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
	Engine        Engine         // Engine component (exported)
	Hull          Hull           // Hull component (exported)
	FuelTank      FuelTank       // FuelTank component (exported)
	CargoHold     CargoHold      // CargoHold component (exported)
	HeatShield    HeatShield     // HeatShield component (exported)
	Drill         Drill          // Drill component (exported)
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
	engine := NewEngineFromConfig(playerCfg.StartingUpgrades.Engine, engineTier.Name, engineTier.Stats)
	hull := NewHullFromConfig(playerCfg.StartingUpgrades.Hull, hullTier.Name, hullTier.Stats)
	fuelTank := NewFuelTankFromConfig(playerCfg.StartingUpgrades.FuelTank, fuelTankTier.Name, fuelTankTier.Stats)
	cargoHold := NewCargoHoldFromConfig(playerCfg.StartingUpgrades.CargoHold, cargoTier.Name, cargoTier.Stats)
	heatShield := NewHeatShieldFromConfig(playerCfg.StartingUpgrades.HeatShield, heatTier.Name, heatTier.Stats)
	drill := NewDrillFromConfig(playerCfg.StartingUpgrades.Drill, drillTier.Name, drillTier.Stats)

	return &Player{
		AABB:          types.NewAABB(startX, startY, PlayerWidth, PlayerHeight),
		Velocity:      types.Zero(),
		OnGround:      false,
		OreInventory:  make(map[string]int),
		ItemInventory: playerCfg.StartingItems,
		Fuel:          fuelTank.Capacity(),
		HP:            hull.MaxHP(),
		Engine:        engine,
		Hull:          hull,
		FuelTank:      fuelTank,
		CargoHold:     cargoHold,
		HeatShield:    heatShield,
		Drill:         drill,
		Money:         playerCfg.StartingMoney,
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

func (p *Player) BuyDrill(d Drill, cost int) {
	p.Money -= cost
	p.Drill = d
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

// AddOreByID increments ore count for given ore ID if capacity allows
// Returns true if ore was added, false if cargo is full
func (p *Player) AddOreByID(oreID string) bool {
	if oreID == "" {
		return false
	}
	if p.GetTotalOreCount() >= p.CargoHold.Capacity() {
		return false // Cargo full
	}
	p.OreInventory[oreID]++
	return true
}

func (p *Player) SellInventory(oreConfigs []config.OreConfig) {
	totalValue := 0
	for oreID, count := range p.OreInventory {
		if count > 0 {
			for _, oreCfg := range oreConfigs {
				if oreCfg.ID == oreID {
					totalValue += oreCfg.Value * count
					break
				}
			}
		}
	}
	p.Money += totalValue
	p.OreInventory = make(map[string]int)
}

func (p *Player) CalculateInventoryValue(oreConfigs []config.OreConfig) int {
	totalValue := 0
	for oreID, count := range p.OreInventory {
		if count > 0 {
			for _, oreCfg := range oreConfigs {
				if oreCfg.ID == oreID {
					totalValue += oreCfg.Value * count
					break
				}
			}
		}
	}
	return totalValue
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
