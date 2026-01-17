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

	// Compact world dimensions (500px padding on each side of buildings)
	worldWidth  = 3072     // 500px + buildings + 500px
	worldHeight = 64 * 800 // 51200 pixels (800 tiles × 64px)
	groundLevel = 640.0    // Aligned to tile boundary (10 * TileSize)

	worldSeed = int64(42) // Seed for procedural world generation

	// Player spawn position (centered in world)
	playerSpawnX = worldWidth / 2
	playerSpawnY = 570.0 // Just above ground

	// 320 per building

	// Building positions (X coordinates, Y calculated from ground level)
	// Layout: 480px pad | Hospital | 50px | FuelStation | 230px | Market | 130px | UpgradeShop | 50px | ItemShop | 532px pad
	hospitalX    = 480.0
	fuelStationX = 850.0
	marketX      = 1400
	upgradeShopX = 1850.0
	itemShopX    = 2220.0
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Drill Game")

	renderer := rendering.NewRaylibRenderer(screenWidth, screenHeight)
	inputAdapter := input.NewRaylibInputAdapter()

	// Initialize window
	renderer.InitWindow(screenWidth, screenHeight, "Drill Game")
	defer renderer.CloseWindow()

	renderer.SetTargetFPS(targetFPS)

	slog.Info("Initializing Game")

	config := &world.WorldConfig{
		Width:       worldWidth,
		Height:      worldHeight,
		GroundLevel: groundLevel,
		Seed:        worldSeed,
		PlayerSpawn: world.PlayerSpawn{
			X: playerSpawnX,
			Y: playerSpawnY,
		},
		BuildingLayout: world.BuildingLayout{
			MarketX:      marketX,
			FuelStationX: fuelStationX,
			HospitalX:    hospitalX,
			UpgradeShopX: upgradeShopX,
			ItemShopX:    itemShopX,
		},
	}

	if err := config.Validate(); err != nil {
		slog.Error("Invalid world configuration", "error", err)
		return
	}

	gameWorld := world.NewWorld(config)
	game := engine.NewGame(gameWorld, config)

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
