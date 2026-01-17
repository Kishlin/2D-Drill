package main

import (
	"log/slog"
	"os"

	"github.com/Kishlin/drill-game/internal/adapters/input"
	"github.com/Kishlin/drill-game/internal/adapters/rendering"
	"github.com/Kishlin/drill-game/internal/domain/engine"
	"github.com/Kishlin/drill-game/internal/domain/world"
	"github.com/Kishlin/drill-game/internal/domain/worlds"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	targetFPS    = 60

	// World selection: change this to switch between different world configurations
	// Available worlds:
	//   - "default":      Normal game balance, player starts with $0 and Base upgrades
	//   - "hard_mode":    Harder difficulty - 50% ore values, 2x upgrade prices, more hazards
	//   - "sandbox":      Easy testing - 10x ore values, free upgrades, no hazards, Mk5 start
	//   - "endgame_test": Mid-tier start with Mk3 upgrades for testing deep content
	selectedWorld = "default"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Drill Game")

	// Load world configuration
	slog.Info("Loading world configuration", "world", selectedWorld)
	worldConfig := worlds.MustGetWorld(selectedWorld)

	renderer := rendering.NewRaylibRenderer(screenWidth, screenHeight)
	inputAdapter := input.NewRaylibInputAdapter()

	// Initialize window
	renderer.InitWindow(screenWidth, screenHeight, "Drill Game")
	defer renderer.CloseWindow()

	renderer.SetTargetFPS(targetFPS)

	slog.Info("Initializing Game")

	// Validate world configuration
	if err := worldConfig.World.Validate(); err != nil {
		slog.Error("Invalid world configuration", "error", err)
		return
	}

	// Create configuration lookups for procedural generation
	oreLookup := worlds.NewOreConfigLookup(&worldConfig.Ores)
	hazardLookup := worlds.NewHazardConfigLookup(&worldConfig.Hazards)

	// Create game world with config lookups
	gameWorld := world.NewWorldWithConfig(worldConfig.World, oreLookup, hazardLookup, &worldConfig.BaseTiles)

	// Create game (to be updated in Phase 4 for full config support)
	game := engine.NewGame(gameWorld, worldConfig.World)

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
