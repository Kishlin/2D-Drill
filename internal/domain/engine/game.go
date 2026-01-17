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
	shopUISystem      *systems.ShopUISystem
	itemSystem        *systems.ItemSystem
	itemShopSystem    *systems.ItemShopSystem
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

	// Create unified upgrade shop to the right of the ore market
	upgradeShopY := w.GetGroundLevel() - entities.UpgradeShopHeight
	upgradeShopX := marketX + 360.0
	upgradeShop := entities.NewUpgradeShop(upgradeShopX, upgradeShopY)

	// Create item shops to the right of upgrade shop
	itemShopY := w.GetGroundLevel() - entities.ItemShopHeight
	teleportShopX := upgradeShopX + 360.0
	teleportShop := entities.NewItemShop(teleportShopX, itemShopY, entities.ItemTeleport, 500, "Teleport")

	repairShopX := teleportShopX + 200.0
	repairShop := entities.NewItemShop(repairShopX, itemShopY, entities.ItemRepair, 200, "Repair Kit")

	refuelShopX := repairShopX + 200.0
	refuelShop := entities.NewItemShop(refuelShopX, itemShopY, entities.ItemRefuel, 100, "Fuel Can")

	bombShopX := refuelShopX + 200.0
	bombShop := entities.NewItemShop(bombShopX, itemShopY, entities.ItemBomb, 300, "Bomb")

	bigBombShopX := bombShopX + 200.0
	bigBombShop := entities.NewItemShop(bigBombShopX, itemShopY, entities.ItemBigBomb, 800, "Big Bomb")

	return &Game{
		world:             w,
		player:            entities.NewPlayer(spawnX, spawnY),
		physicsSystem:     systems.NewPhysicsSystem(w),
		drillingSystem:    systems.NewDrillingSystem(w),
		marketSystem:      systems.NewMarketSystem(market),
		fuelSystem:        systems.NewFuelSystem(),
		fuelStationSystem: systems.NewFuelStationSystem(fuelStation),
		hospitalSystem:    systems.NewHospitalSystem(hospital),
		shopUISystem:      systems.NewShopUISystem(upgradeShop),
		itemSystem:        systems.NewItemSystem(w, spawnX, spawnY),
		itemShopSystem:    systems.NewItemShopSystem(teleportShop, repairShop, refuelShop, bombShop, bigBombShop),
	}
}

func (g *Game) Update(dt float32, inputState input.InputState) error {
	// 0. Update chunks around player (proactive loading)
	playerX := g.player.AABB.X + g.player.AABB.Width/2
	playerY := g.player.AABB.Y + g.player.AABB.Height/2
	g.world.UpdateChunksAroundPlayer(playerX, playerY)

	// 1. Modal pause: if shop is open, only process shop UI and pause all gameplay
	g.shopUISystem.ProcessShopInteraction(g.player, inputState)
	if g.player.InShop {
		return nil
	}

	// 2. Physics - handles landing/fall damage before drilling, heat damage, prevents movement during drilling
	g.physicsSystem.UpdatePhysics(g.player, inputState, dt)

	// 3. Fuel consumption (runs even during drilling animation to maintain resource pressure)
	g.fuelSystem.ConsumeFuel(g.player, inputState, dt)

	// 4. Drilling animation (vertical + horizontal)
	g.drillingSystem.ProcessDrilling(g.player, inputState, dt)

	// Skip interactions if drilling animation is active
	if g.player.IsDrilling {
		return nil
	}

	// 5. Handle item usage
	g.itemSystem.ProcessItemUsage(g.player, inputState)

	// 6. Handle market selling
	g.marketSystem.ProcessSelling(g.player, inputState)

	// 7. Handle fuel station refueling
	g.fuelStationSystem.ProcessRefueling(g.player, inputState)

	// 8. Handle hospital healing
	g.hospitalSystem.ProcessHealing(g.player, inputState)

	// 9. Handle item shop purchases
	g.itemShopSystem.ProcessPurchase(g.player, inputState)

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

func (g *Game) GetUpgradeShop() *entities.UpgradeShop {
	return g.shopUISystem.GetShop()
}

func (g *Game) GetShopUIState() *entities.ShopUIState {
	return g.shopUISystem.GetUIState()
}

func (g *Game) GetItemShops() []*entities.ItemShop {
	return g.itemShopSystem.GetShops()
}
