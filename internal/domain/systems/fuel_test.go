package systems

import (
	"math"
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

func TestFuelSystem_ConsumesMovingRateWhenMovingLeft(t *testing.T) {
	player := testFuelPlayer()
	fuelCapacity := player.FuelCapacity()

	if player.Fuel != fuelCapacity {
		t.Fatalf("expected full tank (%.2f), got %.2f", fuelCapacity, player.Fuel)
	}

	// Simulate 1 second of leftward movement
	inputState := input.InputState{Left: true}
	ConsumeFuel(player, inputState, 1.0, testFuelConfig())

	expectedFuel := fuelCapacity - testFuelConfig().ConsumptionMoving
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s moving left, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesMovingRateWhenMovingRight(t *testing.T) {
	player := testFuelPlayer()
	fuelCapacity := player.FuelCapacity()

	inputState := input.InputState{Right: true}
	ConsumeFuel(player, inputState, 1.0, testFuelConfig())

	expectedFuel := fuelCapacity - testFuelConfig().ConsumptionMoving
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s moving right, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesMovingRateWhenMovingUp(t *testing.T) {
	player := testFuelPlayer()
	fuelCapacity := player.FuelCapacity()

	inputState := input.InputState{Up: true}
	ConsumeFuel(player, inputState, 1.0, testFuelConfig())

	expectedFuel := fuelCapacity - testFuelConfig().ConsumptionMoving
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s flying up, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesMovingRateWhenDrilling(t *testing.T) {
	player := testFuelPlayer()
	fuelCapacity := player.FuelCapacity()

	// Drilling should use movement rate (active work)
	inputState := input.InputState{Drill: true}
	ConsumeFuel(player, inputState, 1.0, testFuelConfig())

	expectedFuel := fuelCapacity - testFuelConfig().ConsumptionMoving
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s drilling, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesIdleRateWhenNoInput(t *testing.T) {
	player := testFuelPlayer()
	fuelCapacity := player.FuelCapacity()

	// No input = idle state
	inputState := input.InputState{}
	ConsumeFuel(player, inputState, 1.0, testFuelConfig())

	expectedFuel := fuelCapacity - testFuelConfig().ConsumptionIdle
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s idle, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesIdleRateWhenOnlyInteractingInput(t *testing.T) {
	player := testFuelPlayer()
	fuelCapacity := player.FuelCapacity()

	// Interact input alone should use idle rate (not active movement)
	inputState := input.InputState{Interact: true}
	ConsumeFuel(player, inputState, 1.0, testFuelConfig())

	expectedFuel := fuelCapacity - testFuelConfig().ConsumptionIdle
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s with sell input, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesMovingRateWhenMovingAndInteracting(t *testing.T) {
	player := testFuelPlayer()
	fuelCapacity := player.FuelCapacity()

	// Moving + selling = use movement rate (movement takes priority)
	inputState := input.InputState{Up: true, Interact: true}
	ConsumeFuel(player, inputState, 1.0, testFuelConfig())

	expectedFuel := fuelCapacity - testFuelConfig().ConsumptionMoving
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel with movement + sell, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_FuelDoesNotGoBelowZero(t *testing.T) {
	player := testFuelPlayer()

	// Consume all fuel in one very large frame
	inputState := input.InputState{Left: true}
	ConsumeFuel(player, inputState, 1000.0, testFuelConfig()) // Way more than 10 liters

	if player.Fuel < 0 {
		t.Errorf("expected fuel >= 0, got %.4f", player.Fuel)
	}

	if player.Fuel != 0 {
		t.Errorf("expected fuel to clamp at 0, got %.4f", player.Fuel)
	}
}

func TestFuelSystem_FrameRateIndependence(t *testing.T) {
	// 60 seconds at two different frame rates should consume the same fuel

	// Test at 60 FPS
	player60 := testFuelPlayer()
	inputState := input.InputState{Up: true}
	frameTime60 := float32(1.0 / 60.0)
	for i := 0; i < 3600; i++ { // 60 frames/sec * 60 seconds
		ConsumeFuel(player60, inputState, frameTime60, testFuelConfig())
	}

	// Test at 30 FPS
	player30 := testFuelPlayer()
	frameTime30 := float32(1.0 / 30.0)
	for i := 0; i < 1800; i++ { // 30 frames/sec * 60 seconds
		ConsumeFuel(player30, inputState, frameTime30, testFuelConfig())
	}

	// Both should have consumed approximately the same amount of fuel
	const tolerance = 0.01 // Within 0.01 liters due to floating point precision
	if math.Abs(float64(player60.Fuel-player30.Fuel)) > tolerance {
		t.Errorf("frame rate dependent consumption: 60fps=%.4f, 30fps=%.4f (diff=%.4f)",
			player60.Fuel, player30.Fuel, math.Abs(float64(player60.Fuel-player30.Fuel)))
	}
}

func TestFuelSystem_FullTankDurationMoving(t *testing.T) {
	// 10 liters in 30 seconds = 0.333 L/s
	// Starting with 10L, moving continuously should last 30 seconds
	player := testFuelPlayer()
	fuelCapacity := player.FuelCapacity()

	inputState := input.InputState{Up: true}
	dt := float32(0.1) // Simulate 0.1 second frames
	remainingFuel := fuelCapacity
	frames := 0
	maxFrames := 500 // Safety limit

	for remainingFuel > 0 && frames < maxFrames {
		ConsumeFuel(player, inputState, dt, testFuelConfig())
		remainingFuel = player.Fuel
		frames++
	}

	elapsedSeconds := float32(frames) * dt
	expectedSeconds := float32(30.0)

	// Should last approximately 30 seconds (within 0.5 seconds due to discrete frame steps)
	if math.Abs(float64(elapsedSeconds-expectedSeconds)) > 0.5 {
		t.Errorf("expected full tank to last ~%.1f seconds when moving, lasted %.1f seconds",
			expectedSeconds, elapsedSeconds)
	}
}

func TestFuelSystem_FullTankDurationIdle(t *testing.T) {
	// 10 liters in 120 seconds = 0.08333 L/s
	// Starting with 10L, idle should last 120 seconds
	player := testFuelPlayer()
	fuelCapacity := player.FuelCapacity()

	inputState := input.InputState{} // No input = idle
	dt := float32(1.0)               // Simulate 1 second frames for speed
	remainingFuel := fuelCapacity
	frames := 0
	maxFrames := 150 // Safety limit

	for remainingFuel > 0 && frames < maxFrames {
		ConsumeFuel(player, inputState, dt, testFuelConfig())
		remainingFuel = player.Fuel
		frames++
	}

	expectedSeconds := float32(120.0)

	// Should last approximately 120 seconds (within 2 seconds due to discrete frame steps)
	if math.Abs(float64(frames)-float64(expectedSeconds)) > 2.0 {
		t.Errorf("expected full tank to last ~%.0f seconds when idle, lasted %.0f seconds",
			expectedSeconds, float32(frames))
	}
}

func TestFuelSystem_MultipleConsumptionsAccumulate(t *testing.T) {
	player := testFuelPlayer()
	fuelCapacity := player.FuelCapacity()
	fuelCfg := testFuelConfig()

	// Consume fuel multiple times
	inputState1 := input.InputState{Left: true}
	ConsumeFuel(player, inputState1, 1.0, fuelCfg) // -0.0833L

	inputState2 := input.InputState{}             // Idle
	ConsumeFuel(player, inputState2, 1.0, fuelCfg) // -0.0167L

	inputState3 := input.InputState{Up: true}
	ConsumeFuel(player, inputState3, 2.0, fuelCfg) // -0.1667L (2 seconds moving)

	// Total should be: 10 - 0.0833 - 0.0167 - 0.1667 = 9.7333
	expectedFuel := fuelCapacity - fuelCfg.ConsumptionMoving - fuelCfg.ConsumptionIdle - (fuelCfg.ConsumptionMoving * 2.0)
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f after multiple consumptions, got %.4f", expectedFuel, player.Fuel)
	}
}

// Test helpers

func testFuelConfig() config.FuelSystemConfig {
	return config.FuelSystemConfig{
		ConsumptionMoving: 10.0 / 30.0,
		ConsumptionIdle:   10.0 / 120.0,
	}
}

func testFuelPlayer() *entities.Player {
	playerCfg := config.PlayerConfig{
		StartingMoney:    0,
		StartingItems:    [5]int{0, 0, 0, 0, 0},
		StartingUpgrades: config.StartingUpgrades{},
	}
	upgradeCfg := config.UpgradeConfig{
		Engines:     []config.UpgradeTier[config.EngineStats]{{Name: "Base", Price: 0, Stats: config.EngineStats{MaxSpeed: 450, Acceleration: 2500, FlyAcceleration: 2500, MaxUpwardSpeed: -600}}},
		Hulls:       []config.UpgradeTier[config.HullStats]{{Name: "Base", Price: 0, Stats: config.HullStats{MaxHP: 10}}},
		FuelTanks:   []config.UpgradeTier[config.FuelTankStats]{{Name: "Base", Price: 0, Stats: config.FuelTankStats{Capacity: 10}}},
		CargoHolds:  []config.UpgradeTier[config.CargoHoldStats]{{Name: "Base", Price: 0, Stats: config.CargoHoldStats{Capacity: 10}}},
		HeatShields: []config.UpgradeTier[config.HeatShieldStats]{{Name: "Base", Price: 0, Stats: config.HeatShieldStats{HeatResistance: 50}}},
		Drills:      []config.UpgradeTier[config.DrillStats]{{Name: "Base", Price: 0, Stats: config.DrillStats{SpeedAtSurface: 1.0, SpeedAtMaxDepth: 1.0}}},
	}
	return entities.NewPlayerFromConfig(0, 0, playerCfg, upgradeCfg)
}
