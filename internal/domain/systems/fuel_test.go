package systems

import (
	"math"
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

func TestFuelSystem_ConsumesMovingRateWhenMovingLeft(t *testing.T) {
	fs := NewFuelSystem()
	player := entities.NewPlayer(0, 0)
	fuelCapacity := player.FuelTank.Capacity()

	if player.Fuel != fuelCapacity {
		t.Fatalf("expected full tank (%.2f), got %.2f", fuelCapacity, player.Fuel)
	}

	// Simulate 1 second of leftward movement
	inputState := input.InputState{Left: true}
	fs.ConsumeFuel(player, inputState, 1.0)

	expectedFuel := fuelCapacity - FuelConsumptionMoving
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s moving left, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesMovingRateWhenMovingRight(t *testing.T) {
	fs := NewFuelSystem()
	player := entities.NewPlayer(0, 0)
	fuelCapacity := player.FuelTank.Capacity()

	inputState := input.InputState{Right: true}
	fs.ConsumeFuel(player, inputState, 1.0)

	expectedFuel := fuelCapacity - FuelConsumptionMoving
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s moving right, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesMovingRateWhenMovingUp(t *testing.T) {
	fs := NewFuelSystem()
	player := entities.NewPlayer(0, 0)
	fuelCapacity := player.FuelTank.Capacity()

	inputState := input.InputState{Up: true}
	fs.ConsumeFuel(player, inputState, 1.0)

	expectedFuel := fuelCapacity - FuelConsumptionMoving
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s flying up, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesMovingRateWhenDigging(t *testing.T) {
	fs := NewFuelSystem()
	player := entities.NewPlayer(0, 0)
	fuelCapacity := player.FuelTank.Capacity()

	// Digging should use movement rate (active work)
	inputState := input.InputState{Dig: true}
	fs.ConsumeFuel(player, inputState, 1.0)

	expectedFuel := fuelCapacity - FuelConsumptionMoving
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s digging, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesIdleRateWhenNoInput(t *testing.T) {
	fs := NewFuelSystem()
	player := entities.NewPlayer(0, 0)
	fuelCapacity := player.FuelTank.Capacity()

	// No input = idle state
	inputState := input.InputState{}
	fs.ConsumeFuel(player, inputState, 1.0)

	expectedFuel := fuelCapacity - FuelConsumptionIdle
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s idle, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesIdleRateWhenOnlySellingInput(t *testing.T) {
	fs := NewFuelSystem()
	player := entities.NewPlayer(0, 0)
	fuelCapacity := player.FuelTank.Capacity()

	// Sell input alone should use idle rate (not active movement)
	inputState := input.InputState{Sell: true}
	fs.ConsumeFuel(player, inputState, 1.0)

	expectedFuel := fuelCapacity - FuelConsumptionIdle
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel after 1s with sell input, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_ConsumesMovingRateWhenMovingAndSelling(t *testing.T) {
	fs := NewFuelSystem()
	player := entities.NewPlayer(0, 0)
	fuelCapacity := player.FuelTank.Capacity()

	// Moving + selling = use movement rate (movement takes priority)
	inputState := input.InputState{Up: true, Sell: true}
	fs.ConsumeFuel(player, inputState, 1.0)

	expectedFuel := fuelCapacity - FuelConsumptionMoving
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f fuel with movement + sell, got %.4f", expectedFuel, player.Fuel)
	}
}

func TestFuelSystem_FuelDoesNotGoBelowZero(t *testing.T) {
	fs := NewFuelSystem()
	player := entities.NewPlayer(0, 0)

	// Consume all fuel in one very large frame
	inputState := input.InputState{Left: true}
	fs.ConsumeFuel(player, inputState, 1000.0) // Way more than 10 liters

	if player.Fuel < 0 {
		t.Errorf("expected fuel >= 0, got %.4f", player.Fuel)
	}

	if player.Fuel != 0 {
		t.Errorf("expected fuel to clamp at 0, got %.4f", player.Fuel)
	}
}

func TestFuelSystem_FrameRateIndependence(t *testing.T) {
	// 60 seconds at two different frame rates should consume the same fuel
	fs := NewFuelSystem()

	// Test at 60 FPS
	player60 := entities.NewPlayer(0, 0)
	inputState := input.InputState{Up: true}
	frameTime60 := float32(1.0 / 60.0)
	for i := 0; i < 3600; i++ { // 60 frames/sec * 60 seconds
		fs.ConsumeFuel(player60, inputState, frameTime60)
	}

	// Test at 30 FPS
	player30 := entities.NewPlayer(0, 0)
	frameTime30 := float32(1.0 / 30.0)
	for i := 0; i < 1800; i++ { // 30 frames/sec * 60 seconds
		fs.ConsumeFuel(player30, inputState, frameTime30)
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
	fs := NewFuelSystem()
	player := entities.NewPlayer(0, 0)
	fuelCapacity := player.FuelTank.Capacity()

	inputState := input.InputState{Up: true}
	dt := float32(0.1) // Simulate 0.1 second frames
	remainingFuel := fuelCapacity
	frames := 0
	maxFrames := 500 // Safety limit

	for remainingFuel > 0 && frames < maxFrames {
		fs.ConsumeFuel(player, inputState, dt)
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
	fs := NewFuelSystem()
	player := entities.NewPlayer(0, 0)
	fuelCapacity := player.FuelTank.Capacity()

	inputState := input.InputState{} // No input = idle
	dt := float32(1.0)               // Simulate 1 second frames for speed
	remainingFuel := fuelCapacity
	frames := 0
	maxFrames := 150 // Safety limit

	for remainingFuel > 0 && frames < maxFrames {
		fs.ConsumeFuel(player, inputState, dt)
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
	fs := NewFuelSystem()
	player := entities.NewPlayer(0, 0)
	fuelCapacity := player.FuelTank.Capacity()

	// Consume fuel multiple times
	inputState1 := input.InputState{Left: true}
	fs.ConsumeFuel(player, inputState1, 1.0) // -0.0833L

	inputState2 := input.InputState{}        // Idle
	fs.ConsumeFuel(player, inputState2, 1.0) // -0.0167L

	inputState3 := input.InputState{Up: true}
	fs.ConsumeFuel(player, inputState3, 2.0) // -0.1667L (2 seconds moving)

	// Total should be: 10 - 0.0833 - 0.0167 - 0.1667 = 9.7333
	expectedFuel := fuelCapacity - FuelConsumptionMoving - FuelConsumptionIdle - (FuelConsumptionMoving * 2.0)
	if math.Abs(float64(player.Fuel-expectedFuel)) > 0.0001 {
		t.Errorf("expected %.4f after multiple consumptions, got %.4f", expectedFuel, player.Fuel)
	}
}
