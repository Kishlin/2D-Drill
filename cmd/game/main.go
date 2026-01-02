package main

import (
	"log/slog"
	"os"

	"github.com/Kishlin/drill-game/internal/adapters/input"
	"github.com/Kishlin/drill-game/internal/adapters/rendering"
	"github.com/Kishlin/drill-game/internal/domain/engine"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	targetFPS    = 60
	groundLevel  = 256.0 // Aligned to tile boundary (4 * TileSize)
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Drill Game")

	renderer := rendering.NewRaylibRenderer()
	inputAdapter := input.NewRaylibInputAdapter()

	// Initialize window
	renderer.InitWindow(screenWidth, screenHeight, "Drill Game")
	defer renderer.CloseWindow()

	renderer.SetTargetFPS(targetFPS)

	slog.Info("Initializing Game")

	gameWorld := world.NewWorld(screenWidth, screenHeight, groundLevel)
	game := engine.NewGame(gameWorld)

	for renderer.WindowShouldClose() == false {
		dt := renderer.GetFrameTime() // Delta time in seconds

		inputState := inputAdapter.ReadInput()

		err := game.Update(dt, inputState)
		if err != nil {
			slog.Error("Error during update", "error", err)
			break
		}

		renderer.Render(game, inputState)
	}

	slog.Info("Shutting down Drill Game")
}
