package physics_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/physics"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// Base engine stats for testing (from upgrade tier 0)
var (
	baseMaxSpeed        = entities.EngineTiers[0].MaxSpeed
	baseAcceleration    = entities.EngineTiers[0].Acceleration
	baseFlyAcceleration = entities.EngineTiers[0].FlyAcceleration
	baseMaxUpwardSpeed  = entities.EngineTiers[0].MaxUpwardSpeed
)

func TestApplyHorizontalMovement_Acceleration(t *testing.T) {
	// Arrange
	velocity := types.Zero()
	inputState := input.InputState{Right: true}
	dt := float32(0.016) // ~60fps

	// Act
	newVelocity := physics.ApplyHorizontalMovement(velocity, inputState, dt, baseMaxSpeed, baseAcceleration)

	// Assert
	expected := baseAcceleration * dt
	if newVelocity.X != expected {
		t.Errorf("Expected velocity X = %f, got %f", expected, newVelocity.X)
	}
}

func TestApplyHorizontalMovement_MaxSpeed(t *testing.T) {
	// Arrange: velocity already at max
	velocity := types.Vec2{X: baseMaxSpeed, Y: 0}
	inputState := input.InputState{Right: true}
	dt := float32(0.016)

	// Act
	newVelocity := physics.ApplyHorizontalMovement(velocity, inputState, dt, baseMaxSpeed, baseAcceleration)

	// Assert: should cap at max speed
	if newVelocity.X != baseMaxSpeed {
		t.Errorf("Expected velocity capped at %f, got %f", baseMaxSpeed, newVelocity.X)
	}
}

func TestApplyHorizontalMovement_Damping(t *testing.T) {
	// Arrange: moving right, no input
	velocity := types.Vec2{X: 100.0, Y: 0}
	inputState := input.InputState{} // no input
	dt := float32(0.016)

	// Act
	newVelocity := physics.ApplyHorizontalMovement(velocity, inputState, dt, baseMaxSpeed, baseAcceleration)

	// Assert: should slow down
	if newVelocity.X >= velocity.X {
		t.Errorf("Expected velocity to decrease due to damping")
	}
}

func TestApplyVerticalMovement_FlyAcceleration(t *testing.T) {
	// Arrange
	velocity := types.Zero()
	inputState := input.InputState{Up: true}
	dt := float32(0.016)

	// Act
	newVelocity := physics.ApplyVerticalMovement(velocity, inputState, dt, baseFlyAcceleration, baseMaxUpwardSpeed)

	// Assert: velocity should be negative (upward)
	expected := -baseFlyAcceleration * dt
	if newVelocity.Y != expected {
		t.Errorf("Expected velocity Y = %f, got %f", expected, newVelocity.Y)
	}
}

func TestApplyVerticalMovement_MaxUpwardSpeed(t *testing.T) {
	// Arrange: already at max upward velocity
	velocity := types.Vec2{X: 0, Y: baseMaxUpwardSpeed}
	inputState := input.InputState{Up: true}
	dt := float32(0.016)

	// Act
	newVelocity := physics.ApplyVerticalMovement(velocity, inputState, dt, baseFlyAcceleration, baseMaxUpwardSpeed)

	// Assert: should cap at max upward velocity
	if newVelocity.Y != baseMaxUpwardSpeed {
		t.Errorf("Expected velocity capped at %f, got %f", baseMaxUpwardSpeed, newVelocity.Y)
	}
}
