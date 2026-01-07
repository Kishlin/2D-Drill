package physics_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/physics"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

var (
	testMaxSpeed        = float32(450.0)
	testAcceleration    = float32(2500.0)
	testFlyAcceleration = float32(2500.0)
	testMaxUpwardSpeed  = float32(-600.0)
)

func TestApplyHorizontalMovement_Acceleration(t *testing.T) {
	// Arrange
	velocity := types.Zero()
	inputState := input.InputState{Right: true}
	dt := float32(0.016) // ~60fps

	// Act
	newVelocity := physics.ApplyHorizontalMovement(velocity, inputState, dt, testMaxSpeed, testAcceleration)

	// Assert
	expected := testAcceleration * dt
	if newVelocity.X != expected {
		t.Errorf("Expected velocity X = %f, got %f", expected, newVelocity.X)
	}
}

func TestApplyHorizontalMovement_MaxSpeed(t *testing.T) {
	// Arrange: velocity already at max
	velocity := types.Vec2{X: testMaxSpeed, Y: 0}
	inputState := input.InputState{Right: true}
	dt := float32(0.016)

	// Act
	newVelocity := physics.ApplyHorizontalMovement(velocity, inputState, dt, testMaxSpeed, testAcceleration)

	// Assert: should cap at max speed
	if newVelocity.X != testMaxSpeed {
		t.Errorf("Expected velocity capped at %f, got %f", testMaxSpeed, newVelocity.X)
	}
}

func TestApplyHorizontalMovement_Damping(t *testing.T) {
	// Arrange: moving right, no input
	velocity := types.Vec2{X: 100.0, Y: 0}
	inputState := input.InputState{} // no input
	dt := float32(0.016)

	// Act
	newVelocity := physics.ApplyHorizontalMovement(velocity, inputState, dt, testMaxSpeed, testAcceleration)

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
	newVelocity := physics.ApplyVerticalMovement(velocity, inputState, dt, testFlyAcceleration, testMaxUpwardSpeed)

	// Assert: velocity should be negative (upward)
	expected := -testFlyAcceleration * dt
	if newVelocity.Y != expected {
		t.Errorf("Expected velocity Y = %f, got %f", expected, newVelocity.Y)
	}
}

func TestApplyVerticalMovement_MaxUpwardSpeed(t *testing.T) {
	// Arrange: already at max upward velocity
	velocity := types.Vec2{X: 0, Y: testMaxUpwardSpeed}
	inputState := input.InputState{Up: true}
	dt := float32(0.016)

	// Act
	newVelocity := physics.ApplyVerticalMovement(velocity, inputState, dt, testFlyAcceleration, testMaxUpwardSpeed)

	// Assert: should cap at max upward velocity
	if newVelocity.Y != testMaxUpwardSpeed {
		t.Errorf("Expected velocity capped at %f, got %f", testMaxUpwardSpeed, newVelocity.Y)
	}
}
