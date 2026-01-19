package engine

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/systems"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type Game struct {
	world               *world.World
	player              *entities.Player
	physicsSystem       *systems.PhysicsSystem
	drillingSystem      *systems.DrillingSystem
	marketSystem        *systems.MarketSystem
	fuelSystem          *systems.FuelSystem
	fuelStationSystem   *systems.FuelStationSystem
	hospitalSystem      *systems.HospitalSystem
	upgradeShopUISystem *systems.UpgradeShopUISystem
	itemSystem          *systems.ItemSystem
	itemShopUISystem    *systems.ItemShopUISystem
	config              *config.GameConfig
}

func NewGame(w *world.World, gameCfg *config.GameConfig) *Game {
	worldCfg := gameCfg.World

	// Use configured player spawn position
	spawnX := worldCfg.PlayerSpawn.X
	spawnY := worldCfg.PlayerSpawn.Y

	// Calculate building Y positions (always ground level - building height)
	marketY := w.GetGroundLevel() - entities.MarketHeight
	fuelStationY := w.GetGroundLevel() - entities.FuelStationHeight
	hospitalY := w.GetGroundLevel() - entities.HospitalHeight
	upgradeShopY := w.GetGroundLevel() - entities.UpgradeShopHeight
	itemShopY := w.GetGroundLevel() - entities.ItemShopHeight

	// Create buildings at configured positions
	market := entities.NewMarket(worldCfg.BuildingLayout.MarketX, marketY)
	fuelStation := entities.NewFuelStation(worldCfg.BuildingLayout.FuelStationX, fuelStationY)
	hospital := entities.NewHospital(worldCfg.BuildingLayout.HospitalX, hospitalY)

	// Create shops from config (allows per-level pricing and tiers)
	upgradeShop := entities.NewUpgradeShopFromConfig(
		worldCfg.BuildingLayout.UpgradeShopX,
		upgradeShopY,
		gameCfg.Upgrades,
	)
	itemShop := entities.NewItemShopFromConfig(
		worldCfg.BuildingLayout.ItemShopX,
		itemShopY,
		gameCfg.Items,
	)

	// Create player from config (allows per-level starting upgrades, money, items)
	player := entities.NewPlayerFromConfig(
		spawnX,
		spawnY,
		gameCfg.Player,
		gameCfg.Upgrades,
	)

	// Create market system with ore configs for value lookup
	marketSystem := systems.NewMarketSystemWithConfig(market, gameCfg.Generation.Ores)

	return &Game{
		world:               w,
		player:              player,
		physicsSystem:       systems.NewPhysicsSystem(w),
		drillingSystem:      systems.NewDrillingSystem(w),
		marketSystem:        marketSystem,
		fuelSystem:          systems.NewFuelSystem(),
		fuelStationSystem:   systems.NewFuelStationSystem(fuelStation),
		hospitalSystem:      systems.NewHospitalSystem(hospital),
		upgradeShopUISystem: systems.NewUpgradeShopUISystem(upgradeShop),
		itemSystem:          systems.NewItemSystemWithConfig(w, spawnX, spawnY, gameCfg.Items),
		itemShopUISystem:    systems.NewItemShopUISystem(itemShop),
		config:              gameCfg,
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

func (g *Game) GetConfig() *config.GameConfig {
	return g.config
}
