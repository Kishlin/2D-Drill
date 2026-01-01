package physics_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/physics"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

func TestApplyGravity_IncreasesDownwardVelocity(t *testing.T) {
	// Arrange
	velocity := types.Zero()
	dt := float32(0.016)

	// Act
	newVelocity := physics.ApplyGravity(velocity, dt)

	// Assert
	expectedY := physics.Gravity * dt
	if newVelocity.Y != expectedY {
		t.Errorf("Expected Y velocity = %f, got %f", expectedY, newVelocity.Y)
	}
}

func TestApplyGravity_PreservesHorizontalVelocity(t *testing.T) {
	// Arrange
	velocity := types.Vec2{X: 100, Y: 0}
	dt := float32(0.016)

	// Act
	newVelocity := physics.ApplyGravity(velocity, dt)

	// Assert
	if newVelocity.X != velocity.X {
		t.Errorf("Expected X velocity preserved at %f, got %f", velocity.X, newVelocity.X)
	}
}

func TestIntegrateVelocity_UpdatesPosition(t *testing.T) {
	// Arrange
	position := types.Vec2{X: 100, Y: 200}
	velocity := types.Vec2{X: 50, Y: -30}
	dt := float32(0.016)

	// Act
	newPosition := physics.IntegrateVelocity(position, velocity, dt)

	// Assert
	expectedX := position.X + velocity.X*dt
	expectedY := position.Y + velocity.Y*dt
	if newPosition.X != expectedX || newPosition.Y != expectedY {
		t.Errorf("Expected position (%f, %f), got (%f, %f)",
			expectedX, expectedY, newPosition.X, newPosition.Y)
	}
}
