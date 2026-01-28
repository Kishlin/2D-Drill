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
	"github.com/Kishlin/drill-game/internal/domain/upgrades"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type Game struct {
	world           *world.World
	player          *entities.Player
	buildings       []*entities.Building
	upgradeCatalog  *upgrades.Catalog
	itemCatalog     *entities.ItemCatalog
	drillingSystem  *systems.DrillingSystem
	uiManager       *ui.Manager
	inventoryUI     *ui.InventoryUI
	effectProcessor *effects.Processor
	effectContext   *effects.EffectContext
	damageables     []effects.DamageableEntity
	boss            bosses.Boss
	bossFightSystem *systems.BossFightSystem
	gameState       entities.GameState
	config          *config.GameConfig
}

func NewGame(gameCfg *config.GameConfig) *Game {
	worldCfg := gameCfg.World

	// Create the world from config
	w := world.NewWorldFromConfig(worldCfg, gameCfg.Generation, gameCfg.Level.BossRoom)

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
	upgradeCatalog := upgrades.NewCatalogFromConfig(gameCfg.Upgrades)
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

	// Create inventory UI (separate from building UIs)
	inventoryUI := ui.NewInventoryUI(gameCfg.Generation.Ores)

	// Create boss and boss fight system
	var damageables []effects.DamageableEntity
	boss, err := createBossByType(
		gameCfg.Level.BossRoom.BossType,
		worldCfg.Height-gameCfg.Level.BossRoom.RoomHeight-gameCfg.Level.BossRoom.FloorHeight*world.TileSize,
		worldCfg.Width,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create boss: %v", err))
	}
	bossFightSystem := systems.NewBossFightSystem(boss, gameCfg.Level.BossRoom, worldCfg.Height)
	// Add physical boss to damageables list
	if physicalBoss, ok := boss.(effects.DamageableEntity); ok {
		damageables = append(damageables, physicalBoss)
	}

	// Create effect context
	effectContext := &effects.EffectContext{
		Player:      player,
		World:       w,
		Damageables: damageables,
	}

	return &Game{
		world:           w,
		player:          player,
		buildings:       buildings,
		upgradeCatalog:  upgradeCatalog,
		itemCatalog:     itemCatalog,
		drillingSystem:  systems.NewDrillingSystemWithConfig(w, gameCfg.Generation, gameCfg.Drilling),
		uiManager:       uiManager,
		inventoryUI:     inventoryUI,
		effectProcessor: effects.NewProcessor(),
		effectContext:   effectContext,
		damageables:     damageables,
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

	// 1a. If inventory UI is active, process it
	if g.inventoryUI.IsActive() {
		closed := g.inventoryUI.Process(inputState)
		if closed {
			g.player.InUI = false
			return nil
		}
		return nil // Still open (modal) - pause gameplay
	}

	// 1b. If a building UI is active, process it
	if g.uiManager.HasActiveUI() {
		result := g.uiManager.Process(g.player, inputState)
		g.effectProcessor.Apply(g.effectContext, result.Effects)

		// If UI closed, resume gameplay but skip interaction check this frame
		// to prevent the same key press from reopening the UI
		if g.uiManager.HasActiveUI() == false {
			g.player.InUI = false
			return nil
		}
		return nil // Still open (modal) - pause gameplay
	}

	// 2a. Check for inventory key press (when no UI is active)
	if inputState.Inventory {
		g.inventoryUI.Open()
		g.player.InUI = true
		return nil
	}

	// 2b. Check for new building interactions
	if interactionType := systems.DetectInteraction(g.player, g.buildings, inputState); interactionType != nil {
		// Reset UI state for modal UIs before opening
		g.resetUIState(*interactionType)

		if g.uiManager.OpenUI(*interactionType) {
			// Process immediately (handles both instant and modal first frame)
			result := g.uiManager.Process(g.player, inputState)
			g.effectProcessor.Apply(g.effectContext, result.Effects)

			// If still open after first process, it's modal - pause
			if g.uiManager.HasActiveUI() {
				g.player.InUI = true
				return nil
			}
			// Otherwise it was instant, continue with gameplay
		}
	}

	// 3. Heat damage - applies damage based on depth-based temperature
	systems.UpdateHeat(g.player, g.world, dt)

	// 4. Physics - handles landing/fall damage before drilling, prevents movement during drilling
	systems.UpdatePhysics(g.player, g.world, inputState, dt)

	// 5. Fuel consumption (runs even during drilling animation to maintain resource pressure)
	systems.ConsumeFuel(g.player, inputState, dt)

	// 6. Drilling animation (vertical + horizontal)
	drillEffects := g.drillingSystem.ProcessDrilling(g.player, inputState, dt)
	if len(drillEffects) > 0 {
		g.effectProcessor.Apply(g.effectContext, drillEffects)
	}

	// Skip interactions if drilling animation is active
	if g.player.IsDrilling {
		return nil
	}

	// 7. Handle item usage
	itemEffects := systems.DetectItemUsage(g.player, inputState, g.config.Items)
	if len(itemEffects) > 0 {
		g.effectProcessor.Apply(g.effectContext, itemEffects)
	}

	// 8. Update boss fight system (if active)
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
	case components.InteractableMarket:
		if marketUI, ok := registeredUI.(*ui.MarketUI); ok {
			marketUI.ResetState()
		}
	case components.InteractableHospital:
		if hospitalUI, ok := registeredUI.(*ui.HospitalUI); ok {
			hospitalUI.ResetState()
		}
	case components.InteractableFuelStation:
		if fuelStationUI, ok := registeredUI.(*ui.FuelStationUI); ok {
			fuelStationUI.ResetState()
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

func (g *Game) GetUpgradeCatalog() *upgrades.Catalog {
	return g.upgradeCatalog
}

func (g *Game) GetItemCatalog() *entities.ItemCatalog {
	return g.itemCatalog
}

func (g *Game) GetUIManager() *ui.Manager {
	return g.uiManager
}

func (g *Game) GetInventoryUI() *ui.InventoryUI {
	return g.inventoryUI
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
