package main

import (
	"log/slog"
	"os"

	"github.com/Kishlin/drill-game/internal/engine"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	targetFPS    = 60
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Drill Game")

	// Initialize window
	rl.InitWindow(screenWidth, screenHeight, "Drill Game")
	defer rl.CloseWindow()

	rl.SetTargetFPS(targetFPS)

	// Initialize game
	game := engine.NewGame()

	slog.Info("Game initialized")

	// Main game loop
	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime() // Delta time in seconds

		// Update
		if err := game.Update(dt); err != nil {
			slog.Error("Error during update", "error", err)
			break
		}

		// Render
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		game.Render()

		rl.EndDrawing()
	}

	slog.Info("Shutting down Drill Game")
}
