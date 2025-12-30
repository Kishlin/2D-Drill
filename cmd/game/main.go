package main

import (
	"log/slog"
	"os"

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

	slog.Info("Window initialized", "width", screenWidth, "height", screenHeight, "fps", targetFPS)

	// TODO: Initialize game state, load assets, etc.

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		if err := update(); err != nil {
			slog.Error("Error during update", "error", err)
			break
		}

		// Render
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		render()

		rl.EndDrawing()
	}

	slog.Info("Shutting down Drill Game")
}

func update() error {
	// TODO: Implement game update logic
	// - Handle input
	// - Update physics
	// - Update entities
	// - Check collisions
	// - etc.

	return nil
}

func render() {
	// TODO: Implement game rendering
	// - Draw world
	// - Draw entities
	// - Draw UI/HUD
	// - etc.

	// Placeholder: Draw "Hello, Drill Game!" text
	rl.DrawText("Hello, Drill Game!", 190, 200, 20, rl.LightGray)
	rl.DrawText("Press ESC to exit", 190, 240, 20, rl.LightGray)
}
