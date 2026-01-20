package rendering

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	bossrenderers "github.com/Kishlin/drill-game/internal/adapters/rendering/bosses"
	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/engine"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/physics"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

var (
	PlayerColor         = rl.Red
	GroundColor         = rl.Brown
	SkyColor            = rl.SkyBlue
	DirtColor           = rl.NewColor(139, 90, 43, 255)   // Brown dirt
	GridColor           = rl.NewColor(100, 65, 30, 128)   // Semi-transparent grid lines
	MarketColor         = rl.NewColor(34, 139, 34, 255)   // Forest Green
	FuelStationColor    = rl.NewColor(255, 165, 0, 255)   // Orange
	HospitalColor       = rl.NewColor(220, 20, 60, 255)   // Crimson
	EngineShopColor     = rl.NewColor(70, 130, 180, 255)  // Steel Blue
	HullShopColor       = rl.NewColor(105, 105, 105, 255) // Dim Gray
	FuelTankShopColor   = rl.NewColor(255, 99, 71, 255)   // Tomato
	CargoHoldShopColor  = rl.NewColor(148, 0, 211, 255)   // Dark Violet
	HeatShieldShopColor = rl.NewColor(255, 69, 0, 255)    // Orange Red
	DrillShopColor      = rl.NewColor(184, 134, 11, 255)  // Dark Goldenrod
	TeleportShopColor   = rl.NewColor(138, 43, 226, 255)  // Blue Violet
	RepairShopColor     = rl.NewColor(34, 139, 34, 255)   // Forest Green
	RefuelShopColor     = rl.NewColor(255, 165, 0, 255)   // Orange
	BombShopColor       = rl.NewColor(255, 20, 147, 255)  // Deep Pink
	BigBombShopColor    = rl.NewColor(220, 20, 60, 255)   // Crimson
	FloorConcreteColor = rl.NewColor(100, 100, 100, 255) // Dark Gray
	FloorLavaColor     = rl.NewColor(255, 100, 0, 255)   // Orange
)

type RaylibRenderer struct {
	camera       rl.Camera2D
	screenWidth  float32
	screenHeight float32
	worldWidth   float32                  // Cached for boundary clamping
	genCfg       *config.GenerationConfig // Ore/hazard colors from config
	oreColors    map[string]rl.Color      // Cached ore colors by ID
	hazardColors map[string]rl.Color      // Cached hazard colors by ID
}

func NewRaylibRendererWithConfig(screenWidth, screenHeight int32, genCfg *config.GenerationConfig) *RaylibRenderer {
	r := &RaylibRenderer{
		camera: rl.Camera2D{
			Offset:   rl.Vector2{X: float32(screenWidth) / 2, Y: float32(screenHeight) / 2},
			Target:   rl.Vector2{X: 0, Y: 0},
			Rotation: 0.0,
			Zoom:     1.0,
		},
		screenWidth:  float32(screenWidth),
		screenHeight: float32(screenHeight),
		worldWidth:   0, // Set on first render
		genCfg:       genCfg,
		oreColors:    make(map[string]rl.Color),
		hazardColors: make(map[string]rl.Color),
	}

	for _, oreCfg := range genCfg.Ores {
		r.oreColors[oreCfg.ID] = rl.NewColor(oreCfg.Color[0], oreCfg.Color[1], oreCfg.Color[2], oreCfg.Color[3])
	}
	for _, hazardCfg := range genCfg.Hazards {
		r.hazardColors[hazardCfg.ID] = rl.NewColor(hazardCfg.Color[0], hazardCfg.Color[1], hazardCfg.Color[2], hazardCfg.Color[3])
	}

	return r
}

// updateCamera sets camera target to player position with boundary clamping
func (r *RaylibRenderer) updateCamera(player *entities.Player, w *world.World) {
	// Cache world width on first call
	if r.worldWidth == 0 {
		r.worldWidth = w.Width
	}

	// Camera targets player center (AABB is top-left corner)
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerCenterY := player.AABB.Y + player.AABB.Height/2

	// Clamp camera to prevent viewing outside world bounds
	halfScreenW := r.screenWidth / 2
	halfScreenH := r.screenHeight / 2

	minX := halfScreenW
	maxX := r.worldWidth - halfScreenW
	minY := w.GetGroundLevel() - halfScreenH // Can't view above sky
	maxY := w.Height - r.screenHeight + 128  // Stop at world bottom with 128px margin

	// Horizontal clamping
	targetX := playerCenterX
	if targetX < minX {
		targetX = minX
	} else if targetX > maxX {
		targetX = maxX
	}

	// Vertical clamping (top and bottom)
	targetY := playerCenterY
	if targetY < minY {
		targetY = minY
	}
	if targetY > maxY {
		targetY = maxY
	}

	r.camera.Target = rl.Vector2{X: targetX, Y: targetY}
}

func (r *RaylibRenderer) Render(game *engine.Game, inputState input.InputState) {
	// Update camera position before rendering
	r.updateCamera(game.GetPlayer(), game.GetWorld())

	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	// === WORLD SPACE (camera transform applied) ===
	rl.BeginMode2D(r.camera)

	r.renderWorld(game.GetWorld())
	r.renderTiles(game.GetWorld())
	r.renderMarket(game.GetMarket())
	r.renderFuelStation(game.GetFuelStation())
	r.renderHospital(game.GetHospital())
	r.renderUpgradeShop(game.GetUpgradeShop().AABB, EngineShopColor, rl.DarkBlue)
	r.renderUpgradeShop(game.GetItemShop().AABB, TeleportShopColor, rl.Purple)
	r.renderPlayer(game.GetPlayer())

	// Boss rendering
	if game.GetBoss() != nil {
		r.renderBoss(game.GetBoss())
		r.renderProjectiles(game.GetBoss())
	}

	rl.EndMode2D()

	// === SCREEN SPACE (no camera, always visible) ===
	r.renderDebugInfo(game.GetPlayer(), inputState)

	// Boss HP bar
	if game.GetBoss() != nil {
		r.renderBossHPBar(game.GetBoss())
	}

	// Game state overlay (victory/defeat screens)
	r.renderGameStateOverlay(game)

	// Render shop modals if open
	if game.GetUpgradeShopUIState().Open {
		r.renderUpgradeShopModal(game)
	}
	if game.GetItemShopUIState().Open {
		r.renderItemShopModal(game)
	}

	rl.EndDrawing()
}

func (r *RaylibRenderer) InitWindow(width, height int32, title string) {
	rl.InitWindow(width, height, title)
}

func (r *RaylibRenderer) CloseWindow() {
	rl.CloseWindow()
}

func (r *RaylibRenderer) WindowShouldClose() bool {
	return rl.WindowShouldClose()
}

func (r *RaylibRenderer) SetTargetFPS(fps int32) {
	rl.SetTargetFPS(fps)
}

func (r *RaylibRenderer) GetFrameTime() float32 {
	return rl.GetFrameTime()
}

func (r *RaylibRenderer) renderPlayer(player *entities.Player) {
	// Convert domain AABB to Raylib rendering
	aabb := player.AABB
	rlPos := rl.Vector2{X: aabb.X, Y: aabb.Y}
	rlSize := rl.Vector2{X: aabb.Width, Y: aabb.Height}
	rl.DrawRectangleV(rlPos, rlSize, PlayerColor)
}

func (r *RaylibRenderer) renderMarket(market *entities.Market) {
	aabb := market.AABB
	rlPos := rl.Vector2{X: aabb.X, Y: aabb.Y}
	rlSize := rl.Vector2{X: aabb.Width, Y: aabb.Height}

	rl.DrawRectangleV(rlPos, rlSize, MarketColor)

	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: aabb.X, Y: aabb.Y, Width: aabb.Width, Height: aabb.Height},
		2.0,
		rl.DarkGreen,
	)
}

func (r *RaylibRenderer) renderFuelStation(fuelStation *entities.FuelStation) {
	aabb := fuelStation.AABB
	rlPos := rl.Vector2{X: aabb.X, Y: aabb.Y}
	rlSize := rl.Vector2{X: aabb.Width, Y: aabb.Height}

	rl.DrawRectangleV(rlPos, rlSize, FuelStationColor)

	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: aabb.X, Y: aabb.Y, Width: aabb.Width, Height: aabb.Height},
		2.0,
		rl.Orange,
	)
}

func (r *RaylibRenderer) renderHospital(hospital *entities.Hospital) {
	aabb := hospital.AABB
	rlPos := rl.Vector2{X: aabb.X, Y: aabb.Y}
	rlSize := rl.Vector2{X: aabb.Width, Y: aabb.Height}

	rl.DrawRectangleV(rlPos, rlSize, HospitalColor)

	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: aabb.X, Y: aabb.Y, Width: aabb.Width, Height: aabb.Height},
		2.0,
		rl.White,
	)
}

func (r *RaylibRenderer) renderUpgradeShop(aabb types.AABB, fillColor, borderColor rl.Color) {
	rlPos := rl.Vector2{X: aabb.X, Y: aabb.Y}
	rlSize := rl.Vector2{X: aabb.Width, Y: aabb.Height}

	rl.DrawRectangleV(rlPos, rlSize, fillColor)

	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: aabb.X, Y: aabb.Y, Width: aabb.Width, Height: aabb.Height},
		2.0,
		borderColor,
	)
}

func (r *RaylibRenderer) renderWorld(w *world.World) {
	groundLevel := w.GetGroundLevel()

	// Draw sky from off-screen above to ground level
	// Extended upward to cover viewport when camera is near top
	skyTop := int32(-r.screenHeight)
	skyHeight := int32(groundLevel) - skyTop

	rl.DrawRectangle(0, skyTop, int32(w.Width), skyHeight, SkyColor)

	// Draw ground (lower portion from groundLevel to world bottom)
	rl.DrawRectangle(0, int32(groundLevel), int32(w.Width), int32(w.Height), GroundColor)
}

func (r *RaylibRenderer) renderTiles(w *world.World) {
	tiles := w.GetAllTiles()

	// Calculate visible tile range based on camera viewport
	// Add 1-tile margin to prevent pop-in at edges
	minVisibleX := int((r.camera.Target.X-r.screenWidth/2)/world.TileSize) - 1
	maxVisibleX := int((r.camera.Target.X+r.screenWidth/2)/world.TileSize) + 1
	minVisibleY := int((r.camera.Target.Y-r.screenHeight/2)/world.TileSize) - 1
	maxVisibleY := int((r.camera.Target.Y+r.screenHeight/2)/world.TileSize) + 1

	for coord, tile := range tiles {
		gridX, gridY := coord[0], coord[1]

		// Skip tiles outside viewport (culling optimization)
		if gridX < minVisibleX || gridX > maxVisibleX ||
			gridY < minVisibleY || gridY > maxVisibleY {
			continue
		}

		pixelX := float32(gridX * world.TileSize)
		pixelY := float32(gridY * world.TileSize)

		// Render tile based on type
		var color rl.Color
		switch tile.Type {
		case entities.TileTypeEmpty:
			continue // Skip empty tiles
		case entities.TileTypeDirt:
			color = DirtColor
		case entities.TileTypeOre:
			if oreColor, ok := r.oreColors[tile.OreID]; ok {
				color = oreColor
			} else {
				color = rl.Magenta // Error color for unknown ore
			}
		case entities.TileTypeRock, entities.TileTypeLava:
			if hazardColor, ok := r.hazardColors[tile.HazardID]; ok {
				color = hazardColor
			} else {
				color = rl.Magenta // Error color for unknown hazard
			}
		case entities.TileTypeFloor:
			color = FloorConcreteColor
		default:
			color = rl.Magenta // Error color for unknown tile type
		}

		// Draw filled tile
		rl.DrawRectangle(
			int32(pixelX),
			int32(pixelY),
			world.TileSize,
			world.TileSize,
			color,
		)

		// Draw grid lines for visual clarity
		rl.DrawRectangleLines(
			int32(pixelX),
			int32(pixelY),
			world.TileSize,
			world.TileSize,
			GridColor,
		)
	}
}

func (r *RaylibRenderer) renderDebugInfo(player *entities.Player, inputState input.InputState) {
	fontSize := int32(20)
	textColor := rl.Black
	lineHeight := int32(25)
	posX := int32(10)
	posY := int32(10)

	// Draw FPS
	fps := rl.GetFPS()
	fpsText := fmt.Sprintf("FPS: %d", fps)
	rl.DrawText(fpsText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw player position, velocity, onground status, and drilling state
	posVelText := fmt.Sprintf("Pos: X=%.1f, Y=%.1f | Vel: X=%.1f, Y=%.1f | OnGround: %v | IsDrilling: %v",
		player.AABB.X, player.AABB.Y, player.Velocity.X, player.Velocity.Y, player.OnGround, player.IsDrilling)
	rl.DrawText(posVelText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw input state
	inputText := fmt.Sprintf("Input: L=%v R=%v U=%v Drill=%v", inputState.Left, inputState.Right, inputState.Up, inputState.Drill)
	rl.DrawText(inputText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw ore inventory (dynamically from config)
	inventoryText := "Ore:"
	for _, oreCfg := range r.genCfg.Ores {
		count := player.OreInventory[oreCfg.ID]
		inventoryText += fmt.Sprintf(" %s=%d", oreCfg.ID[:2], count)
	}
	rl.DrawText(inventoryText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw player money, fuel, HP, and cargo
	totalOre := player.GetTotalOreCount()
	moneyFuelHPText := fmt.Sprintf("Money: $%d | Fuel: %.2fL | HP: %.1f | Cargo: %d/%d",
		player.Money, player.Fuel, player.HP, totalOre, player.CargoHold.Capacity())
	rl.DrawText(moneyFuelHPText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw upgrade levels
	upgradeText := fmt.Sprintf("Upgrades: Engine=%d Hull=%d Tank=%d Cargo=%d Heat=%d Drill=%d",
		player.Engine.Tier(), player.Hull.Tier(), player.FuelTank.Tier(), player.CargoHold.Tier(), player.HeatShield.Tier(), player.Drill.Tier())
	rl.DrawText(upgradeText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw temperature
	temperature := physics.CalculateTemperature(player.AABB.Y)
	tempText := fmt.Sprintf("Temperature: %.1f°C (Resistance: %.1f°C)",
		temperature, player.HeatShield.HeatResistance())
	rl.DrawText(tempText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw item inventory
	itemText := fmt.Sprintf("Items: Teleport=%d Repair=%d Fuel=%d Bomb=%d BigBomb=%d",
		player.ItemInventory[entities.ItemTeleport],
		player.ItemInventory[entities.ItemRepair],
		player.ItemInventory[entities.ItemRefuel],
		player.ItemInventory[entities.ItemBomb],
		player.ItemInventory[entities.ItemBigBomb])
	rl.DrawText(itemText, posX, posY, fontSize, textColor)
}

// renderUpgradeShopModal draws the upgrade shop modal UI
func (r *RaylibRenderer) renderUpgradeShopModal(game *engine.Game) {
	uiState := game.GetUpgradeShopUIState()
	shop := game.GetUpgradeShop()
	player := game.GetPlayer()

	// Modal dimensions
	modalWidth := float32(900)
	modalHeight := float32(550)
	modalX := (r.screenWidth - modalWidth) / 2
	modalY := (r.screenHeight - modalHeight) / 2

	// Draw semi-transparent overlay
	rl.DrawRectangle(0, 0, int32(r.screenWidth), int32(r.screenHeight), rl.NewColor(0, 0, 0, 150))

	// Draw modal background
	rl.DrawRectangle(int32(modalX), int32(modalY), int32(modalWidth), int32(modalHeight), rl.NewColor(40, 40, 50, 255))
	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: modalX, Y: modalY, Width: modalWidth, Height: modalHeight},
		3.0,
		rl.NewColor(100, 100, 120, 255),
	)

	// Title
	titleText := "UPGRADE SHOP"
	titleFontSize := int32(30)
	titleWidth := rl.MeasureText(titleText, titleFontSize)
	rl.DrawText(titleText, int32(modalX)+(int32(modalWidth)-titleWidth)/2, int32(modalY)+10, titleFontSize, rl.White)

	// Tab bar
	tabY := modalY + 50
	tabWidth := float32(140)
	tabHeight := float32(35)
	tabGap := float32(5)
	tabStartX := modalX + (modalWidth-6*(tabWidth+tabGap))/2

	tabNames := []string{"Engine", "Hull", "Fuel", "Cargo", "Heat", "Drill"}
	tabColors := []rl.Color{EngineShopColor, HullShopColor, FuelTankShopColor, CargoHoldShopColor, HeatShieldShopColor, DrillShopColor}

	for i := 0; i < 6; i++ {
		tabX := tabStartX + float32(i)*(tabWidth+tabGap)
		isActive := int(uiState.ActiveTab) == i

		// Tab background
		bgColor := rl.NewColor(60, 60, 70, 255)
		if isActive {
			bgColor = tabColors[i]
		}
		rl.DrawRectangle(int32(tabX), int32(tabY), int32(tabWidth), int32(tabHeight), bgColor)

		// Tab border
		borderColor := rl.NewColor(100, 100, 110, 255)
		if isActive {
			borderColor = rl.White
		}
		rl.DrawRectangleLinesEx(
			rl.Rectangle{X: tabX, Y: tabY, Width: tabWidth, Height: tabHeight},
			2.0,
			borderColor,
		)

		// Tab text
		textWidth := rl.MeasureText(tabNames[i], 18)
		textX := int32(tabX) + (int32(tabWidth)-textWidth)/2
		textY := int32(tabY) + 8
		textColor := rl.LightGray
		if isActive {
			textColor = rl.White
		}
		rl.DrawText(tabNames[i], textX, textY, 18, textColor)
	}

	// Grid area
	gridStartX := modalX + 50
	gridStartY := tabY + tabHeight + 30
	cellSize := float32(100)
	cellGap := float32(15)

	// Get current tier for this upgrade type
	currentTier := entities.GetPlayerCurrentTier(player, uiState.ActiveTab)

	// Draw 2x3 grid of upgrades
	tierNames := []string{"Base", "Mk1", "Mk2", "Mk3", "Mk4", "Mk5"}
	for tier := 0; tier < 6; tier++ {
		row := tier / 3
		col := tier % 3

		cellX := gridStartX + float32(col)*(cellSize+cellGap)
		cellY := gridStartY + float32(row)*(cellSize+cellGap)

		isSelected := uiState.SelectedTier == tier
		isOwned := tier <= currentTier

		// Cell background: owned (green) or check affordability for others
		var bgColor rl.Color
		if isOwned {
			bgColor = rl.NewColor(40, 60, 40, 255) // Green tint for owned
		} else {
			price := shop.GetUpgradePrice(uiState.ActiveTab, tier)
			if player.CanAfford(price) {
				bgColor = rl.NewColor(60, 60, 80, 255) // Light for affordable
			} else {
				bgColor = rl.NewColor(50, 50, 55, 255) // Dark for expensive
			}
		}
		rl.DrawRectangle(int32(cellX), int32(cellY), int32(cellSize), int32(cellSize), bgColor)

		// Selection border
		if isSelected {
			rl.DrawRectangleLinesEx(
				rl.Rectangle{X: cellX, Y: cellY, Width: cellSize, Height: cellSize},
				3.0,
				rl.Yellow,
			)
		} else {
			rl.DrawRectangleLinesEx(
				rl.Rectangle{X: cellX, Y: cellY, Width: cellSize, Height: cellSize},
				1.0,
				rl.NewColor(80, 80, 90, 255),
			)
		}

		// Tier name
		tierTextWidth := rl.MeasureText(tierNames[tier], 20)
		tierTextX := int32(cellX) + (int32(cellSize)-tierTextWidth)/2
		tierTextY := int32(cellY) + 15
		tierTextColor := rl.White
		if isOwned && tier < currentTier {
			tierTextColor = rl.Gray
		}
		rl.DrawText(tierNames[tier], tierTextX, tierTextY, 20, tierTextColor)

		// Status indicator: always show price for non-owned, or status for owned
		var statusText string
		var statusColor rl.Color
		if tier == currentTier {
			statusText = "EQUIPPED"
			statusColor = rl.Green
		} else if isOwned {
			statusText = "OWNED"
			statusColor = rl.NewColor(100, 150, 100, 255)
		} else {
			// Always show price for non-owned upgrades
			price := shop.GetUpgradePrice(uiState.ActiveTab, tier)
			statusText = fmt.Sprintf("$%d", price)
			if player.CanAfford(price) {
				statusColor = rl.Yellow
			} else {
				statusColor = rl.Red
			}
		}
		statusWidth := rl.MeasureText(statusText, 14)
		statusX := int32(cellX) + (int32(cellSize)-statusWidth)/2
		statusY := int32(cellY) + int32(cellSize) - 25
		rl.DrawText(statusText, statusX, statusY, 14, statusColor)
	}

	// Details panel (right side)
	detailsX := gridStartX + 3*(cellSize+cellGap) + 30
	detailsY := gridStartY
	detailsWidth := modalWidth - (detailsX - modalX) - 30

	// Draw details box
	rl.DrawRectangle(int32(detailsX), int32(detailsY), int32(detailsWidth), 200, rl.NewColor(30, 30, 40, 255))
	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: detailsX, Y: detailsY, Width: detailsWidth, Height: 200},
		1.0,
		rl.NewColor(70, 70, 80, 255),
	)

	// Selected upgrade details
	selectedName := shop.GetUpgradeName(uiState.ActiveTab, uiState.SelectedTier)
	rl.DrawText(selectedName, int32(detailsX)+10, int32(detailsY)+10, 24, rl.White)

	// Stats based on upgrade type
	statsY := int32(detailsY) + 45
	r.renderUpgradeStats(shop, uiState.ActiveTab, uiState.SelectedTier, int32(detailsX)+10, statsY)

	// Price
	selectedPrice := shop.GetUpgradePrice(uiState.ActiveTab, uiState.SelectedTier)
	priceText := fmt.Sprintf("Price: $%d", selectedPrice)
	if uiState.SelectedTier <= currentTier {
		priceText = "Already owned"
	}
	rl.DrawText(priceText, int32(detailsX)+10, int32(detailsY)+160, 20, rl.Yellow)

	// Player info
	playerInfoY := detailsY + 220
	rl.DrawText(fmt.Sprintf("Your Money: $%d", player.Money), int32(detailsX)+10, int32(playerInfoY), 18, rl.White)
	rl.DrawText(fmt.Sprintf("Current: %s", shop.GetUpgradeName(uiState.ActiveTab, currentTier)), int32(detailsX)+10, int32(playerInfoY)+25, 18, rl.LightGray)

	// Controls hint at bottom
	controlsY := modalY + modalHeight - 40
	controlsText := "[Z] Prev Tab   [X] Next Tab   [Arrows] Select   [E] Buy   [Q] Close"
	controlsWidth := rl.MeasureText(controlsText, 16)
	rl.DrawText(controlsText, int32(modalX)+(int32(modalWidth)-controlsWidth)/2, int32(controlsY), 16, rl.LightGray)
}

// renderUpgradeStats renders the stats for a specific upgrade
func (r *RaylibRenderer) renderUpgradeStats(shop *entities.UpgradeShop, upgradeType entities.UpgradeType, tier int, x, y int32) {
	lineHeight := int32(22)
	fontSize := int32(16)

	switch upgradeType {
	case entities.UpgradeEngine:
		if entry := shop.GetEngineCatalogEntry(tier); entry != nil {
			rl.DrawText(fmt.Sprintf("Max Speed: %.0f px/s", entry.Engine.MaxSpeed()), x, y, fontSize, rl.LightGray)
			rl.DrawText(fmt.Sprintf("Acceleration: %.0f px/s²", entry.Engine.Acceleration()), x, y+lineHeight, fontSize, rl.LightGray)
			rl.DrawText(fmt.Sprintf("Fly Accel: %.0f px/s²", entry.Engine.FlyAcceleration()), x, y+lineHeight*2, fontSize, rl.LightGray)
			rl.DrawText(fmt.Sprintf("Max Upward: %.0f px/s", -entry.Engine.MaxUpwardSpeed()), x, y+lineHeight*3, fontSize, rl.LightGray)
		}
	case entities.UpgradeHull:
		if entry := shop.GetHullCatalogEntry(tier); entry != nil {
			rl.DrawText(fmt.Sprintf("Max HP: %.0f", entry.Hull.MaxHP()), x, y, fontSize, rl.LightGray)
		}
	case entities.UpgradeFuelTank:
		if entry := shop.GetFuelTankCatalogEntry(tier); entry != nil {
			rl.DrawText(fmt.Sprintf("Capacity: %.0f L", entry.FuelTank.Capacity()), x, y, fontSize, rl.LightGray)
		}
	case entities.UpgradeCargoHold:
		if entry := shop.GetCargoCatalogEntry(tier); entry != nil {
			rl.DrawText(fmt.Sprintf("Capacity: %d ore", entry.CargoHold.Capacity()), x, y, fontSize, rl.LightGray)
		}
	case entities.UpgradeHeatShield:
		if entry := shop.GetHeatCatalogEntry(tier); entry != nil {
			rl.DrawText(fmt.Sprintf("Heat Resistance: %.0f°C", entry.HeatShield.HeatResistance()), x, y, fontSize, rl.LightGray)
		}
	case entities.UpgradeDrill:
		if entry := shop.GetDrillCatalogEntry(tier); entry != nil {
			rl.DrawText(fmt.Sprintf("Drill Speed: %.1fx", entry.Drill.DrillSpeed()), x, y, fontSize, rl.LightGray)
		}
	}
}

// renderItemShopModal draws the item shop modal UI
func (r *RaylibRenderer) renderItemShopModal(game *engine.Game) {
	uiState := game.GetItemShopUIState()
	shop := game.GetItemShop()
	player := game.GetPlayer()

	// Modal dimensions
	modalWidth := float32(750)
	modalHeight := float32(500)
	modalX := (r.screenWidth - modalWidth) / 2
	modalY := (r.screenHeight - modalHeight) / 2

	// Draw semi-transparent overlay
	rl.DrawRectangle(0, 0, int32(r.screenWidth), int32(r.screenHeight), rl.NewColor(0, 0, 0, 150))

	// Draw modal background
	rl.DrawRectangle(int32(modalX), int32(modalY), int32(modalWidth), int32(modalHeight), rl.NewColor(40, 40, 50, 255))
	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: modalX, Y: modalY, Width: modalWidth, Height: modalHeight},
		3.0,
		rl.NewColor(100, 100, 120, 255),
	)

	// Title
	titleText := "ITEM SHOP"
	titleFontSize := int32(30)
	titleWidth := rl.MeasureText(titleText, titleFontSize)
	rl.DrawText(titleText, int32(modalX)+(int32(modalWidth)-titleWidth)/2, int32(modalY)+10, titleFontSize, rl.White)

	// Grid area (2x3 grid)
	gridStartX := modalX + 50
	gridStartY := modalY + 70
	cellSize := float32(100)
	cellGap := float32(20)

	// Item names
	itemNames := []string{"Teleport", "Repair", "Refuel", "Bomb", "Big Bomb"}

	// Draw 2x3 grid (5 items, last cell empty)
	for index := 0; index < 6; index++ {
		row := index / 3
		col := index % 3

		cellX := gridStartX + float32(col)*(cellSize+cellGap)
		cellY := gridStartY + float32(row)*(cellSize+cellGap)

		// Skip the empty cell
		if index == 5 {
			continue
		}

		isSelected := uiState.SelectedIndex == index

		// Get item info
		catalogEntry := shop.GetItem(index)
		if catalogEntry == nil {
			continue
		}

		price := catalogEntry.Price
		owned := player.ItemInventory[catalogEntry.ItemType]

		// Cell background: check affordability
		var bgColor rl.Color
		if player.CanAfford(price) {
			bgColor = rl.NewColor(60, 60, 80, 255) // Light for affordable
		} else {
			bgColor = rl.NewColor(50, 50, 55, 255) // Dark for expensive
		}
		rl.DrawRectangle(int32(cellX), int32(cellY), int32(cellSize), int32(cellSize), bgColor)

		// Selection border
		if isSelected {
			rl.DrawRectangleLinesEx(
				rl.Rectangle{X: cellX, Y: cellY, Width: cellSize, Height: cellSize},
				3.0,
				rl.Yellow,
			)
		} else {
			rl.DrawRectangleLinesEx(
				rl.Rectangle{X: cellX, Y: cellY, Width: cellSize, Height: cellSize},
				1.0,
				rl.NewColor(80, 80, 90, 255),
			)
		}

		// Item name
		itemNameWidth := rl.MeasureText(itemNames[index], 18)
		itemNameX := int32(cellX) + (int32(cellSize)-itemNameWidth)/2
		itemNameY := int32(cellY) + 10
		rl.DrawText(itemNames[index], itemNameX, itemNameY, 18, rl.White)

		// Price
		priceText := fmt.Sprintf("$%d", price)
		var priceColor rl.Color
		if player.CanAfford(price) {
			priceColor = rl.Yellow
		} else {
			priceColor = rl.Red
		}
		priceWidth := rl.MeasureText(priceText, 16)
		priceX := int32(cellX) + (int32(cellSize)-priceWidth)/2
		priceY := int32(cellY) + 40
		rl.DrawText(priceText, priceX, priceY, 16, priceColor)

		// Owned count
		ownedText := fmt.Sprintf("Own: %d", owned)
		ownedWidth := rl.MeasureText(ownedText, 14)
		ownedX := int32(cellX) + (int32(cellSize)-ownedWidth)/2
		ownedY := int32(cellY) + 65
		rl.DrawText(ownedText, ownedX, ownedY, 14, rl.LightGray)
	}

	// Details panel (right side)
	detailsX := gridStartX + 3*(cellSize+cellGap) + 30
	detailsY := gridStartY
	detailsWidth := modalWidth - (detailsX - modalX) - 30

	// Draw details box
	rl.DrawRectangle(int32(detailsX), int32(detailsY), int32(detailsWidth), 200, rl.NewColor(30, 30, 40, 255))
	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: detailsX, Y: detailsY, Width: detailsWidth, Height: 200},
		1.0,
		rl.NewColor(70, 70, 80, 255),
	)

	// Selected item details
	selectedEntry := shop.GetItem(uiState.SelectedIndex)
	if selectedEntry != nil {
		selectedName := itemNames[uiState.SelectedIndex]
		rl.DrawText(selectedName, int32(detailsX)+10, int32(detailsY)+10, 24, rl.White)

		// Item effect text
		effectX := int32(detailsX) + 10
		effectY := int32(detailsY) + 50
		effectText := ""
		switch uiState.SelectedIndex {
		case 0:
			effectText = "Teleport to surface"
		case 1:
			effectText = "Restore full HP"
		case 2:
			effectText = "Refill fuel tank"
		case 3:
			effectText = "Destroy nearby tiles"
		case 4:
			effectText = "Destroy large area"
		}
		rl.DrawText(effectText, effectX, effectY, 16, rl.LightGray)

		// Price
		price := selectedEntry.Price
		priceText := fmt.Sprintf("Price: $%d", price)
		priceColor := rl.Yellow
		if !player.CanAfford(price) {
			priceColor = rl.Red
		}
		rl.DrawText(priceText, int32(detailsX)+10, int32(detailsY)+160, 20, priceColor)
	}

	// Player info
	playerInfoY := detailsY + 220
	rl.DrawText(fmt.Sprintf("Your Money: $%d", player.Money), int32(detailsX)+10, int32(playerInfoY), 18, rl.White)
	if selectedEntry != nil {
		rl.DrawText(fmt.Sprintf("Owned: %d", player.ItemInventory[selectedEntry.ItemType]), int32(detailsX)+10, int32(playerInfoY)+25, 18, rl.LightGray)
	}

	// Controls hint at bottom
	controlsY := modalY + modalHeight - 40
	controlsText := "[Arrows] Navigate   [E] Buy   [Q] Close"
	controlsWidth := rl.MeasureText(controlsText, 16)
	rl.DrawText(controlsText, int32(modalX)+(int32(modalWidth)-controlsWidth)/2, int32(controlsY), 16, rl.LightGray)
}

func (r *RaylibRenderer) renderBoss(boss bosses.Boss) {
	if boss == nil {
		return
	}

	// Try boss-specific renderers first
	if bossrenderers.RenderBoss(boss) {
		return
	}

	// Fallback to generic rendering
	bossrenderers.RenderGeneric(boss)
}

func (r *RaylibRenderer) renderBossHPBar(boss bosses.Boss) {
	if boss == nil || !boss.IsActive() {
		return
	}

	barWidth := float32(200)
	barHeight := float32(20)
	barX := r.screenWidth/2 - barWidth/2 // Center horizontally
	barY := float32(20)                  // Top of screen

	// Background (empty bar)
	rl.DrawRectangle(
		int32(barX),
		int32(barY),
		int32(barWidth),
		int32(barHeight),
		rl.DarkGray,
	)

	// Health fill
	healthPercent := boss.GetHP() / boss.GetMaxHP()
	fillWidth := barWidth * healthPercent

	// Color changes based on health
	var color rl.Color
	if healthPercent > 0.66 {
		color = rl.Green
	} else if healthPercent > 0.33 {
		color = rl.Orange
	} else {
		color = rl.Red
	}

	rl.DrawRectangle(
		int32(barX),
		int32(barY),
		int32(fillWidth),
		int32(barHeight),
		color,
	)

	// Border
	rl.DrawRectangleLines(
		int32(barX),
		int32(barY),
		int32(barWidth),
		int32(barHeight),
		rl.Black,
	)

	// Text: "HP: X / Y"
	hpText := fmt.Sprintf("HP: %.0f / %.0f", boss.GetHP(), boss.GetMaxHP())
	rl.DrawText(
		hpText,
		int32(barX+5),
		int32(barY+2),
		12,
		rl.White,
	)
}

// renderProjectiles renders boss projectiles
func (r *RaylibRenderer) renderProjectiles(boss bosses.Boss) {
	if boss == nil {
		return
	}

	projectiles := boss.GetProjectiles()
	for _, proj := range projectiles {
		if !proj.Active {
			continue
		}

		// Draw projectile as small circle
		centerX := proj.AABB.X + proj.AABB.Width/2
		centerY := proj.AABB.Y + proj.AABB.Height/2

		rl.DrawCircle(
			int32(centerX),
			int32(centerY),
			proj.AABB.Width/2,
			rl.Yellow,
		)
	}
}

// renderGameStateOverlay renders victory/defeat screens
func (r *RaylibRenderer) renderGameStateOverlay(game *engine.Game) {
	state := game.GetGameState()

	switch state {
	case entities.GameStateVictory:
		r.renderVictoryScreen()
	case entities.GameStateDefeat:
		r.renderDefeatScreen()
	}
}

func (r *RaylibRenderer) renderVictoryScreen() {
	// Semi-transparent overlay
	rl.DrawRectangle(0, 0, int32(r.screenWidth), int32(r.screenHeight), rl.NewColor(0, 0, 0, 200))

	// Victory text
	victoryText := "VICTORY!"
	victoryWidth := rl.MeasureText(victoryText, 40)
	rl.DrawText(
		victoryText,
		int32(r.screenWidth)/2-victoryWidth/2,
		int32(r.screenHeight)/2-40,
		40,
		rl.Green,
	)

	// Instructions
	instructText := "Press ESC to return to main menu"
	instructWidth := rl.MeasureText(instructText, 20)
	rl.DrawText(
		instructText,
		int32(r.screenWidth)/2-instructWidth/2,
		int32(r.screenHeight)/2+40,
		20,
		rl.White,
	)
}

func (r *RaylibRenderer) renderDefeatScreen() {
	// Semi-transparent overlay
	rl.DrawRectangle(0, 0, int32(r.screenWidth), int32(r.screenHeight), rl.NewColor(0, 0, 0, 200))

	// Defeat text
	defeatText := "DEFEATED!"
	defeatWidth := rl.MeasureText(defeatText, 40)
	rl.DrawText(
		defeatText,
		int32(r.screenWidth)/2-defeatWidth/2,
		int32(r.screenHeight)/2-40,
		40,
		rl.Red,
	)

	// Instructions
	instructText := "Press ESC to try again"
	instructWidth := rl.MeasureText(instructText, 20)
	rl.DrawText(
		instructText,
		int32(r.screenWidth)/2-instructWidth/2,
		int32(r.screenHeight)/2+40,
		20,
		rl.White,
	)
}
