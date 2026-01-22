package engine

import (
	"fmt"

	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/bosses/test_boss"
	"github.com/Kishlin/drill-game/internal/domain/components"
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/systems"
	"github.com/Kishlin/drill-game/internal/domain/ui"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type Game struct {
	world           *world.World
	player          *entities.Player
	buildings       []*entities.Building
	upgradeCatalog  *entities.UpgradeCatalog
	itemCatalog     *entities.ItemCatalog
	physicsSystem   *systems.PhysicsSystem
	drillingSystem  *systems.DrillingSystem
	fuelSystem      *systems.FuelSystem
	itemSystem      *systems.ItemSystem
	uiManager       *ui.Manager
	effectProcessor *effects.Processor
	boss            bosses.Boss
	bossFightSystem *systems.BossFightSystem
	gameState       entities.GameState
	config          *config.GameConfig
}

func NewGame(gameCfg *config.GameConfig) *Game {
	worldCfg := gameCfg.World

	// Create the world from config
	w := world.NewWorldFromConfigWithBoss(&worldCfg, gameCfg.Generation, gameCfg.Level.BossRoom)

	// Use configured player spawn position
	spawnX := worldCfg.PlayerSpawn.X
	spawnY := worldCfg.PlayerSpawn.Y

	// Calculate building Y positions (always ground level - building height)
	buildingY := w.GetGroundLevel() - entities.BuildingHeight

	// Create buildings
	buildings := []*entities.Building{
		entities.NewMarketBuilding(worldCfg.BuildingLayout.MarketX, buildingY),
		entities.NewFuelStationBuilding(worldCfg.BuildingLayout.FuelStationX, buildingY),
		entities.NewHospitalBuilding(worldCfg.BuildingLayout.HospitalX, buildingY),
		entities.NewUpgradeShopBuilding(worldCfg.BuildingLayout.UpgradeShopX, buildingY),
		entities.NewItemShopBuilding(worldCfg.BuildingLayout.ItemShopX, buildingY),
	}

	// Create catalogs
	upgradeCatalog := entities.NewUpgradeCatalogFromConfig(gameCfg.Upgrades)
	itemCatalog := entities.NewItemCatalogFromConfig(gameCfg.Items)

	// Create player from config
	player := entities.NewPlayerFromConfig(
		spawnX,
		spawnY,
		gameCfg.Player,
		gameCfg.Upgrades,
	)

	// Create UI manager and register UIs
	uiManager := ui.NewManager()
	uiManager.Register(components.InteractableMarket, ui.NewMarketUI(gameCfg.Generation.Ores))
	uiManager.Register(components.InteractableFuelStation, ui.NewFuelStationUI())
	uiManager.Register(components.InteractableHospital, ui.NewHospitalUI())
	uiManager.Register(components.InteractableUpgradeShop, ui.NewUpgradeShopUI(upgradeCatalog))
	uiManager.Register(components.InteractableItemShop, ui.NewItemShopUI(itemCatalog))

	// Create item system
	itemSystem := systems.NewItemSystemWithConfig(w, spawnX, spawnY, gameCfg.Items)

	// Create boss and boss fight system if configured
	var boss bosses.Boss
	var bossFightSystem *systems.BossFightSystem
	if gameCfg.Level.BossRoom != nil {
		var err error
		boss, err = createBossByType(
			gameCfg.Level.BossRoom.BossType,
			worldCfg.Height-gameCfg.Level.BossRoom.RoomHeight-gameCfg.Level.BossRoom.FloorHeight*world.TileSize,
			worldCfg.Width,
		)
		if err == nil && boss != nil {
			bossFightSystem = systems.NewBossFightSystem(boss, gameCfg.Level.BossRoom, worldCfg.Height)
			itemSystem.SetBossFightSystem(bossFightSystem)
		}
	}

	return &Game{
		world:           w,
		player:          player,
		buildings:       buildings,
		upgradeCatalog:  upgradeCatalog,
		itemCatalog:     itemCatalog,
		physicsSystem:   systems.NewPhysicsSystem(w),
		drillingSystem:  systems.NewDrillingSystem(w),
		fuelSystem:      systems.NewFuelSystem(),
		itemSystem:      itemSystem,
		uiManager:       uiManager,
		effectProcessor: effects.NewProcessor(),
		boss:            boss,
		bossFightSystem: bossFightSystem,
		gameState:       entities.GameStatePlaying,
		config:          gameCfg,
	}
}

func (g *Game) Update(dt float32, inputState input.InputState) error {
	// 0. Update chunks around player (proactive loading)
	playerX := g.player.AABB.X + g.player.AABB.Width/2
	playerY := g.player.AABB.Y + g.player.AABB.Height/2
	g.world.UpdateChunksAroundPlayer(playerX, playerY)

	// 1. If a UI is active, process it
	if g.uiManager.HasActiveUI() {
		result := g.uiManager.Process(g.player, inputState)
		g.effectProcessor.Apply(g.player, result.Effects)

		// If UI closed, resume gameplay
		if !g.uiManager.HasActiveUI() {
			g.player.InShop = false
		} else {
			return nil // Still open (modal) - pause gameplay
		}
	}

	// 2. Check for new interactions
	if interactionType := systems.DetectInteraction(g.player, g.buildings, inputState); interactionType != nil {
		// Reset UI state for modal UIs before opening
		g.resetUIState(*interactionType)

		if g.uiManager.OpenUI(*interactionType) {
			// Process immediately (handles both instant and modal first frame)
			result := g.uiManager.Process(g.player, inputState)
			g.effectProcessor.Apply(g.player, result.Effects)

			// If still open after first process, it's modal - pause
			if g.uiManager.HasActiveUI() {
				g.player.InShop = true
				return nil
			}
			// Otherwise it was instant, continue with gameplay
		}
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

	// 7. Update boss fight system (if active)
	if g.bossFightSystem != nil {
		g.gameState = g.bossFightSystem.Update(g.player, dt)
	}

	return nil
}

func (g *Game) resetUIState(interactionType components.InteractableType) {
	registeredUI := g.uiManager.GetRegisteredUI(interactionType)
	if registeredUI == nil {
		return
	}

	switch interactionType {
	case components.InteractableUpgradeShop:
		if upgradeUI, ok := registeredUI.(*ui.UpgradeShopUI); ok {
			upgradeUI.ResetState()
		}
	case components.InteractableItemShop:
		if itemUI, ok := registeredUI.(*ui.ItemShopUI); ok {
			itemUI.ResetState()
		}
	}
}

func (g *Game) GetWorld() *world.World {
	return g.world
}

func (g *Game) GetPlayer() *entities.Player {
	return g.player
}

func (g *Game) GetBuildings() []*entities.Building {
	return g.buildings
}

func (g *Game) GetUpgradeCatalog() *entities.UpgradeCatalog {
	return g.upgradeCatalog
}

func (g *Game) GetItemCatalog() *entities.ItemCatalog {
	return g.itemCatalog
}

func (g *Game) GetUIManager() *ui.Manager {
	return g.uiManager
}

func (g *Game) GetConfig() *config.GameConfig {
	return g.config
}

func (g *Game) GetBoss() bosses.Boss {
	return g.boss
}

func (g *Game) GetBossFightSystem() *systems.BossFightSystem {
	return g.bossFightSystem
}

func (g *Game) GetGameState() entities.GameState {
	return g.gameState
}

func (g *Game) IsBossFightActive() bool {
	if g.bossFightSystem == nil {
		return false
	}
	return g.bossFightSystem.IsBossFightActive()
}

func createBossByType(bossType string, roomStartY, worldWidth float32) (bosses.Boss, error) {
	switch bossType {
	case "test_boss":
		return test_boss.New(roomStartY, worldWidth), nil
	default:
		return nil, fmt.Errorf("unknown boss type: %s", bossType)
	}
}
