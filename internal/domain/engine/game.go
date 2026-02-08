package engine

import (
	"fmt"

	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/components"
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/systems"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/ui"
	"github.com/Kishlin/drill-game/internal/domain/upgrades"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type Game struct {
	// Core state
	World     *world.World
	Player    *entities.Player
	Buildings []*entities.Building
	Boss      bosses.Boss
	GameState entities.GameState

	// UI
	UIManager   *ui.Manager
	InventoryUI *ui.InventoryUI

	// Render data
	Projectiles []types.AABB

	// Systems
	drillingSystem  *systems.DrillingSystem
	bossFightSystem *systems.BossFightSystem

	// Projectile internals
	projectilePool   []systems.Projectile
	projectileBounds systems.ProjectileBounds

	// Effects
	effectProcessor *effects.Processor
	effectContext   *effects.EffectContext

	// Config
	config *config.GameConfig
}

func NewGame(gameCfg *config.GameConfig) (*Game, error) {
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
	boss, err := bosses.Create(
		gameCfg.Level.BossRoom.BossType,
		worldCfg.Height-gameCfg.Level.BossRoom.RoomHeight-gameCfg.Level.BossRoom.FloorHeight*world.TileSize,
		worldCfg.Width,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create boss: %w", err)
	}
	bossFightSystem := systems.NewBossFightSystem(boss, gameCfg.Level.BossRoom, worldCfg.Height)

	// Build damageables list for effect context
	var damageables []effects.DamageableEntity
	if physicalBoss, ok := boss.(effects.DamageableEntity); ok {
		damageables = append(damageables, physicalBoss)
	}

	effectContext := &effects.EffectContext{
		Player:      player,
		World:       w,
		Damageables: damageables,
	}

	return &Game{
		World:           w,
		Player:          player,
		Buildings:       buildings,
		Boss:            boss,
		GameState:       entities.GameStatePlaying,
		UIManager:       uiManager,
		InventoryUI:     inventoryUI,
		drillingSystem:  systems.NewDrillingSystemWithConfig(w, gameCfg.Generation, gameCfg.Drilling),
		bossFightSystem: bossFightSystem,
		projectilePool:  systems.NewProjectilePool(),
		projectileBounds: systems.ProjectileBounds{
			MinX: -100,
			MaxX: worldCfg.Width + 100,
			MinY: -100,
			MaxY: worldCfg.Height + 100,
		},
		effectProcessor: effects.NewProcessor(),
		effectContext:   effectContext,
		config:          gameCfg,
	}, nil
}

func (g *Game) Update(dt float32, inputState input.InputState) error {
	// 0. Update chunks around player (proactive loading)
	playerX := g.Player.AABB.X + g.Player.AABB.Width/2
	playerY := g.Player.AABB.Y + g.Player.AABB.Height/2
	g.World.UpdateChunksAroundPlayer(playerX, playerY)

	// 1a. If inventory UI is active, process it
	if g.InventoryUI.IsActive() {
		closed := g.InventoryUI.Process(inputState)
		if closed {
			g.Player.InUI = false
			return nil
		}
		return nil // Still open (modal) - pause gameplay
	}

	// 1b. If a building UI is active, process it
	if g.UIManager.HasActiveUI() {
		result := g.UIManager.Process(g.Player, inputState)
		g.effectProcessor.Apply(g.effectContext, result.Effects)

		// If UI closed, resume gameplay but skip interaction check this frame
		// to prevent the same key press from reopening the UI
		if g.UIManager.HasActiveUI() == false {
			g.Player.InUI = false
			return nil
		}
		return nil // Still open (modal) - pause gameplay
	}

	// 2a. Check for inventory key press (when no UI is active)
	if inputState.Inventory {
		g.InventoryUI.Open()
		g.Player.InUI = true
		return nil
	}

	// 2b. Check for new building interactions
	if interactionType := systems.DetectInteraction(g.Player, g.Buildings, inputState); interactionType != nil {
		// Reset UI state for modal UIs before opening
		g.resetUIState(*interactionType)

		if g.UIManager.OpenUI(*interactionType) {
			// Process immediately (handles both instant and modal first frame)
			result := g.UIManager.Process(g.Player, inputState)
			g.effectProcessor.Apply(g.effectContext, result.Effects)

			// If still open after first process, it's modal - pause
			if g.UIManager.HasActiveUI() {
				g.Player.InUI = true
				return nil
			}
			// Otherwise it was instant, continue with gameplay
		}
	}

	// 3. Heat damage - applies damage based on depth-based temperature
	systems.UpdateHeat(g.Player, g.World, dt, g.config.Heat)

	// 4. Physics - handles landing/fall damage before drilling, prevents movement during drilling
	systems.UpdatePhysics(g.Player, g.World, inputState, dt)

	// 5. Fuel consumption (runs even during drilling animation to maintain resource pressure)
	systems.ConsumeFuel(g.Player, inputState, dt, g.config.Fuel)

	// 6. Drilling animation (vertical + horizontal)
	drillEffects := g.drillingSystem.ProcessDrilling(g.Player, inputState, dt)
	if len(drillEffects) > 0 {
		g.effectProcessor.Apply(g.effectContext, drillEffects)
	}

	// Skip interactions if drilling animation is active
	if g.Player.IsDrilling {
		return nil
	}

	// 7. Handle item usage
	itemEffects := systems.DetectItemUsage(g.Player, inputState, g.config.Items)
	if len(itemEffects) > 0 {
		g.effectProcessor.Apply(g.effectContext, itemEffects)
	}

	// 8. Update boss fight system - returns spawn requests
	if g.bossFightSystem != nil {
		result := g.bossFightSystem.Update(g.Player, dt)
		g.GameState = result.GameState
		systems.SpawnProjectiles(g.projectilePool, result.SpawnRequests)
	}

	// 9. Update projectile system (moves, culls, detects collisions)
	projectileEffects := systems.UpdateProjectiles(g.projectilePool, g.projectileBounds, dt, []systems.CollisionTarget{g.Player})
	if len(projectileEffects) > 0 {
		g.effectProcessor.Apply(g.effectContext, projectileEffects)
	}

	// 10. Collect active projectile positions for rendering
	g.Projectiles = g.Projectiles[:0]
	for i := range g.projectilePool {
		if g.projectilePool[i].IsActive() {
			g.Projectiles = append(g.Projectiles, g.projectilePool[i].AABB())
		}
	}

	return nil
}

func (g *Game) resetUIState(interactionType components.InteractableType) {
	registeredUI := g.UIManager.GetRegisteredUI(interactionType)
	if registeredUI == nil {
		return
	}

	registeredUI.ResetState()
}

func (g *Game) IsBossFightActive() bool {
	if g.bossFightSystem == nil {
		return false
	}
	return g.bossFightSystem.IsBossFightActive()
}
