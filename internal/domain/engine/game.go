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
	upgradeShopUISystem *systems.UpgradeShopUISystem
	itemSystem        *systems.ItemSystem
	itemShopUISystem  *systems.ItemShopUISystem
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

	// Create unified upgrade shop to the right of the ore market (with void spacing)
	upgradeShopY := w.GetGroundLevel() - entities.UpgradeShopHeight
	upgradeShopX := marketX + 450.0
	upgradeShop := entities.NewUpgradeShop(upgradeShopX, upgradeShopY)

	// Create unified item shop to the right of upgrade shop (with 40px gap, matching fuel/hospital spacing)
	itemShopY := w.GetGroundLevel() - entities.ItemShopHeight
	itemShopX := upgradeShopX + 360.0
	itemShop := entities.NewItemShop(itemShopX, itemShopY)

	return &Game{
		world:             w,
		player:            entities.NewPlayer(spawnX, spawnY),
		physicsSystem:     systems.NewPhysicsSystem(w),
		drillingSystem:    systems.NewDrillingSystem(w),
		marketSystem:      systems.NewMarketSystem(market),
		fuelSystem:        systems.NewFuelSystem(),
		fuelStationSystem: systems.NewFuelStationSystem(fuelStation),
		hospitalSystem:    systems.NewHospitalSystem(hospital),
		upgradeShopUISystem: systems.NewUpgradeShopUISystem(upgradeShop),
		itemSystem:        systems.NewItemSystem(w, spawnX, spawnY),
		itemShopUISystem:  systems.NewItemShopUISystem(itemShop),
	}
}

func (g *Game) Update(dt float32, inputState input.InputState) error {
	// 0. Update chunks around player (proactive loading)
	playerX := g.player.AABB.X + g.player.AABB.Width/2
	playerY := g.player.AABB.Y + g.player.AABB.Height/2
	g.world.UpdateChunksAroundPlayer(playerX, playerY)

	// 1. Modal pause: if upgrade shop is open, pause all gameplay
	g.upgradeShopUISystem.ProcessShopInteraction(g.player, inputState)
	if g.upgradeShopUISystem.GetUIState().Open {
		return nil
	}

	// 2. Modal pause: if item shop is open, pause all gameplay
	g.itemShopUISystem.ProcessItemShopInteraction(g.player, inputState)
	if g.itemShopUISystem.GetUIState().Open {
		return nil
	}

	// 3. Physics - handles landing/fall damage before drilling, heat damage, prevents movement during drilling
	g.physicsSystem.UpdatePhysics(g.player, inputState, dt)

	// 4. Fuel consumption (runs even during drilling animation to maintain resource pressure)
	g.fuelSystem.ConsumeFuel(g.player, inputState, dt)

	// 5. Drilling animation (vertical + horizontal)
	g.drillingSystem.ProcessDrilling(g.player, inputState, dt)

	// Skip interactions if drilling animation is active
	if g.player.IsDrilling {
		return nil
	}

	// 6. Handle item usage
	g.itemSystem.ProcessItemUsage(g.player, inputState)

	// 7. Handle market selling
	g.marketSystem.ProcessSelling(g.player, inputState)

	// 8. Handle fuel station refueling
	g.fuelStationSystem.ProcessRefueling(g.player, inputState)

	// 9. Handle hospital healing
	g.hospitalSystem.ProcessHealing(g.player, inputState)

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
	return g.upgradeShopUISystem.GetShop()
}

func (g *Game) GetUpgradeShopUIState() *entities.UpgradeShopUIState {
	return g.upgradeShopUISystem.GetUIState()
}

func (g *Game) GetItemShop() *entities.ItemShop {
	return g.itemShopUISystem.GetShop()
}

func (g *Game) GetItemShopUIState() *entities.ItemShopUIState {
	return g.itemShopUISystem.GetUIState()
}
