package main

import (
	"log/slog"
	"os"

	"github.com/Kishlin/drill-game/internal/adapters/input"
	"github.com/Kishlin/drill-game/internal/adapters/rendering"
	"github.com/Kishlin/drill-game/internal/domain/engine"
	"github.com/Kishlin/drill-game/internal/domain/levels"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	targetFPS    = 60
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Drill Game")
	slog.Info("Initializing Game")

	gameCfg, err := levels.GetLevelConfig(-2)
	if err != nil {
		slog.Error("Failed to load level config", "error", err)
		return
	}

	if err := gameCfg.Validate(); err != nil {
		slog.Error("Invalid game configuration", "error", err)
		return
	}

	renderer := rendering.NewRaylibRendererWithConfig(screenWidth, screenHeight, &gameCfg.Generation)
	inputAdapter := input.NewRaylibInputAdapter()

	renderer.InitWindow(screenWidth, screenHeight, "Drill Game")
	defer renderer.CloseWindow()

	renderer.SetTargetFPS(targetFPS)

	gameWorld := world.NewWorldFromConfigWithBoss(&gameCfg.World, gameCfg.Generation, gameCfg.Level.BossRoom)

	game := engine.NewGame(gameWorld, gameCfg)

	for {
		dt := renderer.GetFrameTime() // Delta time in seconds

		inputState := inputAdapter.ReadInput()

		err := game.Update(dt, inputState)
		if err != nil {
			slog.Error("Error during update", "error", err)
			break
		}

		renderer.Render(game, inputState)

		// Only check for window close if not in shop (Escape key should close shop, not game)
		if !game.GetPlayer().InShop && renderer.WindowShouldClose() {
			break
		}
	}

	slog.Info("Shutting down Drill Game")
}
