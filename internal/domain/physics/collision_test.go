package physics_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/physics"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

func TestResolveGroundCollision_PlayerAboveGround(t *testing.T) {
	// Arrange
	w := world.NewWorld(1280, 720, 600)
	position := types.Vec2{X: 100, Y: 400} // above ground (400 + 64 = 464 < 600)
	velocity := types.Vec2{X: 0, Y: 50}
	height := float32(64)

	// Act
	result := physics.ResolveGroundCollision(position, velocity, height, w)

	// Assert
	if result.OnGround {
		t.Error("Player should not be on ground when above it")
	}
	if result.Position != position {
		t.Error("Position should not change when above ground")
	}
	if result.Velocity != velocity {
		t.Error("Velocity should not change when above ground")
	}
}

func TestResolveGroundCollision_PlayerBelowGround(t *testing.T) {
	// Arrange
	w := world.NewWorld(1280, 720, 600)
	position := types.Vec2{X: 100, Y: 550} // 550 + 64 = 614 > 600 (below ground)
	velocity := types.Vec2{X: 0, Y: 100}
	height := float32(64)

	// Act
	result := physics.ResolveGroundCollision(position, velocity, height, w)

	// Assert
	if !result.OnGround {
		t.Error("Player should be on ground")
	}
	expectedY := w.GetGroundLevel() - height
	if result.Position.Y != expectedY {
		t.Errorf("Expected Y = %f, got %f", expectedY, result.Position.Y)
	}
	if result.Velocity.Y != 0 {
		t.Error("Vertical velocity should be zeroed on ground collision")
	}
}

func TestResolveGroundCollision_AllowsTakeoff(t *testing.T) {
	// Arrange: on ground but moving upward
	w := world.NewWorld(1280, 720, 600)
	position := types.Vec2{X: 100, Y: 536} // 536 + 64 = 600 (exactly on ground)
	velocity := types.Vec2{X: 0, Y: -100}  // moving upward
	height := float32(64)

	// Act
	result := physics.ResolveGroundCollision(position, velocity, height, w)

	// Assert
	if result.Velocity.Y != velocity.Y {
		t.Errorf("Upward velocity should be preserved, expected %f, got %f", velocity.Y, result.Velocity.Y)
	}
}
