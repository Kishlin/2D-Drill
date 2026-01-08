package engine

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/systems"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type Game struct {
	world             *world.World
	player            *entities.Player
	physicsSystem     *systems.PhysicsSystem
	drillingSystem    *systems.DrillingSystem
	marketSystem      *systems.MarketSystem
	fuelSystem        *systems.FuelSystem
	fuelStationSystem *systems.FuelStationSystem
	hospitalSystem    *systems.HospitalSystem
	upgradeSystem     *systems.UpgradeSystem
	itemSystem        *systems.ItemSystem
}

func NewGame(w *world.World) *Game {
	// Spawn player at center of world horizontally, just above ground
	spawnX := (w.Width / 2) - (entities.PlayerWidth / 2)
	spawnY := w.GetGroundLevel() - entities.PlayerHeight - 10

	// Create market to the right of player spawn
	marketX := spawnX + 200.0 // ~3 tiles to the right
	marketY := w.GetGroundLevel() - entities.MarketHeight
	market := entities.NewMarket(marketX, marketY)

	// Create fuel station to the left of player spawn
	fuelStationX := spawnX - 520.0 // ~8 tiles to the left
	fuelStationY := w.GetGroundLevel() - entities.FuelStationHeight
	fuelStation := entities.NewFuelStation(fuelStationX, fuelStationY)

	// Create hospital to the left of fuel station
	hospitalX := fuelStationX - 360.0 // ~5 tiles + gap to the left
	hospitalY := w.GetGroundLevel() - entities.HospitalHeight
	hospital := entities.NewHospital(hospitalX, hospitalY)

	// Create upgrade shops to the right of the ore market
	upgradeShopY := w.GetGroundLevel() - entities.UpgradeShopHeight
	engineShopX := marketX + 360.0
	engineShop := entities.NewEngineUpgradeShop(engineShopX, upgradeShopY)

	hullShopX := engineShopX + 360.0
	hullShop := entities.NewHullUpgradeShop(hullShopX, upgradeShopY)

	fuelTankShopX := hullShopX + 360.0
	fuelTankShop := entities.NewFuelTankUpgradeShop(fuelTankShopX, upgradeShopY)

	cargoHoldShopX := fuelTankShopX + 360.0
	cargoHoldShop := entities.NewCargoHoldUpgradeShop(cargoHoldShopX, upgradeShopY)

	heatShieldShopX := cargoHoldShopX + 360.0
	heatShieldShop := entities.NewHeatShieldUpgradeShop(heatShieldShopX, upgradeShopY)

	drillShopX := heatShieldShopX + 360.0
	drillShop := entities.NewDrillUpgradeShop(drillShopX, upgradeShopY)

	return &Game{
		world:             w,
		player:            entities.NewPlayer(spawnX, spawnY),
		physicsSystem:     systems.NewPhysicsSystem(w),
		drillingSystem:    systems.NewDrillingSystem(w),
		marketSystem:      systems.NewMarketSystem(market),
		fuelSystem:        systems.NewFuelSystem(),
		fuelStationSystem: systems.NewFuelStationSystem(fuelStation),
		hospitalSystem:    systems.NewHospitalSystem(hospital),
		upgradeSystem:     systems.NewUpgradeSystem(engineShop, hullShop, fuelTankShop, cargoHoldShop, heatShieldShop, drillShop),
		itemSystem:        systems.NewItemSystem(w, spawnX, spawnY),
	}
}

func (g *Game) Update(dt float32, inputState input.InputState) error {
	// 0. Update chunks around player (proactive loading)
	playerX := g.player.AABB.X + g.player.AABB.Width/2
	playerY := g.player.AABB.Y + g.player.AABB.Height/2
	g.world.UpdateChunksAroundPlayer(playerX, playerY)

	// 1. Physics FIRST - handles landing/fall damage before drilling can start
	//    Also applies heat damage and skips movement during drilling animation
	g.physicsSystem.UpdatePhysics(g.player, inputState, dt)

	// 2. Always: fuel consumption (runs even during drilling animation)
	g.fuelSystem.ConsumeFuel(g.player, inputState, dt)

	// 3. Handle drilling (vertical + horizontal, with animation)
	g.drillingSystem.ProcessDrilling(g.player, inputState, dt)

	// Skip interactions during drilling animation
	if g.player.IsDrilling {
		return nil
	}

	// 4. Handle item usage
	g.itemSystem.ProcessItemUsage(g.player, inputState)

	// 5. Handle market selling
	g.marketSystem.ProcessSelling(g.player, inputState)

	// 6. Handle fuel station refueling
	g.fuelStationSystem.ProcessRefueling(g.player, inputState)

	// 7. Handle hospital healing
	g.hospitalSystem.ProcessHealing(g.player, inputState)

	// 8. Handle upgrade purchases
	g.upgradeSystem.ProcessUpgrade(g.player, inputState)

	return nil
}

func (g *Game) GetWorld() *world.World {
	return g.world
}

func (g *Game) GetPlayer() *entities.Player {
	return g.player
}

func (g *Game) GetMarket() *entities.Market {
	return g.marketSystem.GetMarket()
}

func (g *Game) GetFuelStation() *entities.FuelStation {
	return g.fuelStationSystem.GetFuelStation()
}

func (g *Game) GetHospital() *entities.Hospital {
	return g.hospitalSystem.GetHospital()
}

func (g *Game) GetEngineShop() *entities.EngineUpgradeShop {
	return g.upgradeSystem.GetEngineShop()
}

func (g *Game) GetHullShop() *entities.HullUpgradeShop {
	return g.upgradeSystem.GetHullShop()
}

func (g *Game) GetFuelTankShop() *entities.FuelTankUpgradeShop {
	return g.upgradeSystem.GetFuelTankShop()
}

func (g *Game) GetCargoHoldShop() *entities.CargoHoldUpgradeShop {
	return g.upgradeSystem.GetCargoHoldShop()
}

func (g *Game) GetHeatShieldShop() *entities.HeatShieldUpgradeShop {
	return g.upgradeSystem.GetHeatShieldShop()
}

func (g *Game) GetDrillShop() *entities.DrillUpgradeShop {
	return g.upgradeSystem.GetDrillShop()
}
