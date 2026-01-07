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
	diggingSystem     *systems.DiggingSystem
	shopSystem        *systems.ShopSystem
	fuelSystem        *systems.FuelSystem
	fuelStationSystem *systems.FuelStationSystem
	hospitalSystem    *systems.HospitalSystem
	upgradeSystem     *systems.UpgradeSystem
}

func NewGame(w *world.World) *Game {
	// Spawn player at center of world horizontally, just above ground
	spawnX := (w.Width / 2) - (entities.PlayerWidth / 2)
	spawnY := w.GetGroundLevel() - entities.PlayerHeight - 10

	// Create shop to the right of player spawn
	shopX := spawnX + 200.0 // ~3 tiles to the right
	shopY := w.GetGroundLevel() - entities.ShopHeight
	shop := entities.NewShop(shopX, shopY)

	// Create fuel station to the left of player spawn
	fuelStationX := spawnX - 520.0 // ~8 tiles to the left
	fuelStationY := w.GetGroundLevel() - entities.FuelStationHeight
	fuelStation := entities.NewFuelStation(fuelStationX, fuelStationY)

	// Create hospital to the left of fuel station
	hospitalX := fuelStationX - 360.0 // ~5 tiles + gap to the left
	hospitalY := w.GetGroundLevel() - entities.HospitalHeight
	hospital := entities.NewHospital(hospitalX, hospitalY)

	// Create upgrade shops to the right of the ore shop
	upgradeShopY := w.GetGroundLevel() - entities.UpgradeShopHeight
	engineShopX := shopX + 360.0
	engineShop := entities.NewEngineUpgradeShop(engineShopX, upgradeShopY)

	hullShopX := engineShopX + 360.0
	hullShop := entities.NewHullUpgradeShop(hullShopX, upgradeShopY)

	fuelTankShopX := hullShopX + 360.0
	fuelTankShop := entities.NewFuelTankUpgradeShop(fuelTankShopX, upgradeShopY)

	cargoHoldShopX := fuelTankShopX + 360.0
	cargoHoldShop := entities.NewCargoHoldUpgradeShop(cargoHoldShopX, upgradeShopY)

	heatShieldShopX := cargoHoldShopX + 360.0
	heatShieldShop := entities.NewHeatShieldUpgradeShop(heatShieldShopX, upgradeShopY)

	return &Game{
		world:             w,
		player:            entities.NewPlayer(spawnX, spawnY),
		physicsSystem:     systems.NewPhysicsSystem(w),
		diggingSystem:     systems.NewDiggingSystem(w),
		shopSystem:        systems.NewShopSystem(shop),
		fuelSystem:        systems.NewFuelSystem(),
		fuelStationSystem: systems.NewFuelStationSystem(fuelStation),
		hospitalSystem:    systems.NewHospitalSystem(hospital),
		upgradeSystem:     systems.NewUpgradeSystem(engineShop, hullShop, fuelTankShop, cargoHoldShop, heatShieldShop),
	}
}

func (g *Game) Update(dt float32, inputState input.InputState) error {
	// 0. Update chunks around player (proactive loading)
	playerX := g.player.AABB.X + g.player.AABB.Width/2
	playerY := g.player.AABB.Y + g.player.AABB.Height/2
	g.world.UpdateChunksAroundPlayer(playerX, playerY)

	// 1. Handle downward digging (before physics, so alignment happens first)
	g.diggingSystem.ProcessDigging(g.player, inputState)

	// 2. Handle horizontal digging (before physics, so blocked tiles can be dug)
	g.diggingSystem.ProcessHorizontalDigging(g.player, inputState)

	// 3. Handle shop selling (before physics, so player position is stable)
	g.shopSystem.ProcessSelling(g.player, inputState)

	// 4. Handle fuel station refueling (before physics, so player position is stable)
	g.fuelStationSystem.ProcessRefueling(g.player, inputState)

	// 5. Handle hospital healing (before physics, so player position is stable)
	g.hospitalSystem.ProcessHealing(g.player, inputState)

	// 6. Handle upgrade purchases (before physics, so player position is stable)
	g.upgradeSystem.ProcessUpgrade(g.player, inputState)

	// 7. Update physics (includes heat damage)
	g.physicsSystem.UpdatePhysics(g.player, inputState, dt)

	// 8. Consume fuel based on activity
	g.fuelSystem.ConsumeFuel(g.player, inputState, dt)

	return nil
}

func (g *Game) GetWorld() *world.World {
	return g.world
}

func (g *Game) GetPlayer() *entities.Player {
	return g.player
}

func (g *Game) GetShop() *entities.Shop {
	return g.shopSystem.GetShop()
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
