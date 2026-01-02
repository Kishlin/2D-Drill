package rendering

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/Kishlin/drill-game/internal/domain/engine"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

var (
	PlayerColor = rl.Red
	GroundColor = rl.Brown
	SkyColor    = rl.SkyBlue
	DirtColor   = rl.NewColor(139, 90, 43, 255) // Brown dirt
	GridColor   = rl.NewColor(100, 65, 30, 128) // Semi-transparent grid lines
)

type RaylibRenderer struct{}

func NewRaylibRenderer() *RaylibRenderer {
	return &RaylibRenderer{}
}

func (r *RaylibRenderer) Render(game *engine.Game) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	r.renderWorld(game.GetWorld())
	r.renderTiles(game.GetWorld())
	r.renderPlayer(game.GetPlayer())

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

func (r *RaylibRenderer) renderWorld(w *world.World) {
	groundLevel := w.GetGroundLevel()

	// Draw sky (upper portion)
	rl.DrawRectangle(0, 0, int32(w.Width), int32(groundLevel), SkyColor)

	// Draw ground (lower portion)
	rl.DrawRectangle(0, int32(groundLevel), int32(w.Width), int32(w.Height), GroundColor)
}

func (r *RaylibRenderer) renderTiles(w *world.World) {
	tiles := w.GetAllTiles()

	for coord, tile := range tiles {
		gridX, gridY := coord[0], coord[1]
		pixelX := float32(gridX * world.TileSize)
		pixelY := float32(gridY * world.TileSize)

		// Render tile based on type
		var color rl.Color
		switch tile.Type {
		case entities.TileTypeDirt:
			color = DirtColor
		default:
			continue // Skip empty tiles
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
