package rendering

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/Kishlin/drill-game/internal/domain/engine"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/physics"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

var (
	PlayerColor       = rl.Red
	GroundColor       = rl.Brown
	SkyColor          = rl.SkyBlue
	DirtColor         = rl.NewColor(139, 90, 43, 255)   // Brown dirt
	GridColor         = rl.NewColor(100, 65, 30, 128)   // Semi-transparent grid lines
	ShopColor         = rl.NewColor(34, 139, 34, 255)   // Forest Green
	FuelStationColor  = rl.NewColor(255, 165, 0, 255)   // Orange
	HospitalColor     = rl.NewColor(220, 20, 60, 255)   // Crimson
	EngineShopColor    = rl.NewColor(70, 130, 180, 255)  // Steel Blue
	HullShopColor      = rl.NewColor(105, 105, 105, 255) // Dim Gray
	FuelTankShopColor  = rl.NewColor(255, 99, 71, 255)   // Tomato
	CargoHoldShopColor = rl.NewColor(148, 0, 211, 255)   // Dark Violet
	HeatShieldShopColor = rl.NewColor(255, 69, 0, 255)   // Orange Red

	// Ore colors for different ore types
	OreColors = map[entities.OreType]rl.Color{
		entities.OreCopper:   rl.NewColor(255, 140, 0, 255),   // Orange
		entities.OreIron:     rl.NewColor(128, 128, 128, 255), // Gray
		entities.OreSilver:   rl.NewColor(192, 192, 192, 255), // Light Gray
		entities.OreGold:     rl.NewColor(255, 215, 0, 255),   // Gold
		entities.OreMythril:  rl.NewColor(0, 255, 255, 255),   // Cyan
		entities.OrePlatinum: rl.NewColor(230, 230, 250, 255), // White-ish
		entities.OreDiamond:  rl.NewColor(0, 191, 255, 255),   // Blue
	}
)

type RaylibRenderer struct {
	camera       rl.Camera2D
	screenWidth  float32
	screenHeight float32
	worldWidth   float32 // Cached for boundary clamping
}

func NewRaylibRenderer(screenWidth, screenHeight int32) *RaylibRenderer {
	return &RaylibRenderer{
		camera: rl.Camera2D{
			Offset:   rl.Vector2{X: float32(screenWidth) / 2, Y: float32(screenHeight) / 2},
			Target:   rl.Vector2{X: 0, Y: 0},
			Rotation: 0.0,
			Zoom:     1.0,
		},
		screenWidth:  float32(screenWidth),
		screenHeight: float32(screenHeight),
		worldWidth:   0, // Set on first render
	}
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

	// Horizontal clamping
	targetX := playerCenterX
	if targetX < minX {
		targetX = minX
	} else if targetX > maxX {
		targetX = maxX
	}

	// Vertical clamping (top only, no bottom limit per user requirement)
	targetY := playerCenterY
	if targetY < minY {
		targetY = minY
	}
	// No maxY check - camera follows player infinitely downward

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
	r.renderShop(game.GetShop())
	r.renderFuelStation(game.GetFuelStation())
	r.renderHospital(game.GetHospital())
	r.renderUpgradeShop(game.GetEngineShop().AABB, EngineShopColor, rl.DarkBlue)
	r.renderUpgradeShop(game.GetHullShop().AABB, HullShopColor, rl.DarkGray)
	r.renderUpgradeShop(game.GetFuelTankShop().AABB, FuelTankShopColor, rl.Maroon)
	r.renderUpgradeShop(game.GetCargoHoldShop().AABB, CargoHoldShopColor, rl.NewColor(75, 0, 130, 255))
	r.renderUpgradeShop(game.GetHeatShieldShop().AABB, HeatShieldShopColor, rl.Red)
	r.renderPlayer(game.GetPlayer())

	rl.EndMode2D()

	// === SCREEN SPACE (no camera, always visible) ===
	r.renderDebugInfo(game.GetPlayer(), inputState)

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

func (r *RaylibRenderer) renderShop(shop *entities.Shop) {
	aabb := shop.AABB
	rlPos := rl.Vector2{X: aabb.X, Y: aabb.Y}
	rlSize := rl.Vector2{X: aabb.Width, Y: aabb.Height}

	rl.DrawRectangleV(rlPos, rlSize, ShopColor)

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
			var ok bool
			color, ok = OreColors[tile.OreType]
			if !ok {
				color = rl.Magenta // Error color for unknown ore
			}
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

	// Draw player position and velocity
	posVelText := fmt.Sprintf("Pos: X=%.1f, Y=%.1f | Vel: X=%.1f, Y=%.1f", player.AABB.X, player.AABB.Y, player.Velocity.X, player.Velocity.Y)
	rl.DrawText(posVelText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw on ground status
	onGroundText := fmt.Sprintf("OnGround: %v", player.OnGround)
	rl.DrawText(onGroundText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw input state
	inputText := fmt.Sprintf("Input: L=%v R=%v U=%v D=%v", inputState.Left, inputState.Right, inputState.Up, inputState.Dig)
	rl.DrawText(inputText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw ore inventory
	inventoryText := fmt.Sprintf("Ore: Cu=%d Fe=%d Ag=%d Au=%d My=%d Pt=%d Di=%d",
		player.OreInventory[entities.OreCopper],
		player.OreInventory[entities.OreIron],
		player.OreInventory[entities.OreSilver],
		player.OreInventory[entities.OreGold],
		player.OreInventory[entities.OreMythril],
		player.OreInventory[entities.OrePlatinum],
		player.OreInventory[entities.OreDiamond])
	rl.DrawText(inventoryText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw player money, fuel, HP, and cargo
	totalOre := player.GetTotalOreCount()
	moneyFuelHPText := fmt.Sprintf("Money: $%d | Fuel: %.2fL | HP: %.1f | Cargo: %d/%d",
		player.Money, player.Fuel, player.HP, totalOre, player.CargoHold.Capacity())
	rl.DrawText(moneyFuelHPText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw upgrade levels
	upgradeText := fmt.Sprintf("Upgrades: Engine=%d Hull=%d Tank=%d Cargo=%d Heat=%d",
		player.Engine.Tier(), player.Hull.Tier(), player.FuelTank.Tier(), player.CargoHold.Tier(), player.HeatShield.Tier())
	rl.DrawText(upgradeText, posX, posY, fontSize, textColor)
	posY += lineHeight

	// Draw temperature
	temperature := physics.CalculateTemperature(player.AABB.Y)
	tempText := fmt.Sprintf("Temperature: %.1f°C (Resistance: %.1f°C)",
		temperature, player.HeatShield.HeatResistance())
	rl.DrawText(tempText, posX, posY, fontSize, textColor)
}
