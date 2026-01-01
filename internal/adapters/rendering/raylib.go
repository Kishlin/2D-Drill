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
)

type RaylibRenderer struct{}

func NewRaylibRenderer() *RaylibRenderer {
	return &RaylibRenderer{}
}

func (r *RaylibRenderer) Render(game *engine.Game) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	r.renderWorld(game.GetWorld())
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
	// Convert domain Vec2 to Raylib Vector2
	pos := player.GetPositionVec()
	rlPos := rl.Vector2{X: pos.X, Y: pos.Y}
	rlSize := rl.Vector2{X: entities.PlayerWidth, Y: entities.PlayerHeight}
	rl.DrawRectangleV(rlPos, rlSize, PlayerColor)
}

func (r *RaylibRenderer) renderWorld(w *world.World) {
	groundLevel := w.GetGroundLevel()

	// Draw sky (upper portion)
	rl.DrawRectangle(0, 0, int32(w.Width), int32(groundLevel), SkyColor)

	// Draw ground (lower portion)
	rl.DrawRectangle(0, int32(groundLevel), int32(w.Width), int32(w.Height), GroundColor)
}
