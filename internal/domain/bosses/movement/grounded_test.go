package movement

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/types"
)

func TestGrounded_MovesRight(t *testing.T) {
	cfg := GroundedConfig{
		Speed:     100.0,
		MinX:      0,
		MaxX:      1000,
		FloorY:    500,
		BossWidth: 100,
	}
	movement := NewGrounded(cfg)

	startPos := types.NewVec2(400, 500)
	newPos := movement.Update(startPos, 1.0) // 1 second

	if newPos.X <= startPos.X {
		t.Errorf("expected X to increase, got %f -> %f", startPos.X, newPos.X)
	}

	expectedX := startPos.X + 100.0 // speed * dt
	if newPos.X != expectedX {
		t.Errorf("expected X = %f, got %f", expectedX, newPos.X)
	}
}

func TestGrounded_ReversesAtRightBoundary(t *testing.T) {
	cfg := GroundedConfig{
		Speed:     100.0,
		MinX:      0,
		MaxX:      500,
		FloorY:    500,
		BossWidth: 100,
	}
	movement := NewGrounded(cfg)

	// Start near right boundary
	startPos := types.NewVec2(390, 500)
	newPos := movement.Update(startPos, 1.0)

	// Should stop at boundary (MaxX - BossWidth = 400)
	if newPos.X > 400 {
		t.Errorf("expected X <= 400, got %f", newPos.X)
	}

	// Direction should have reversed
	if movement.GetDirection() != -1 {
		t.Error("expected direction to reverse to -1")
	}
}

func TestGrounded_ReversesAtLeftBoundary(t *testing.T) {
	cfg := GroundedConfig{
		Speed:     100.0,
		MinX:      50,
		MaxX:      500,
		FloorY:    500,
		BossWidth: 100,
	}
	movement := NewGrounded(cfg)

	// Manually set direction to left
	movement.reverseDirection()

	startPos := types.NewVec2(60, 500)
	newPos := movement.Update(startPos, 1.0)

	// Should stop at boundary (MinX = 50)
	if newPos.X < 50 {
		t.Errorf("expected X >= 50, got %f", newPos.X)
	}

	// Direction should have reversed back to right
	if movement.GetDirection() != 1 {
		t.Error("expected direction to reverse to 1")
	}
}

func TestGrounded_SetSpeed(t *testing.T) {
	cfg := GroundedConfig{
		Speed:     100.0,
		MinX:      0,
		MaxX:      1000,
		FloorY:    500,
		BossWidth: 100,
	}
	movement := NewGrounded(cfg)

	movement.SetSpeed(200.0)

	if movement.GetSpeed() != 200.0 {
		t.Errorf("expected speed 200, got %f", movement.GetSpeed())
	}

	// Velocity should also update
	vel := movement.GetVelocity()
	if vel.X != 200.0 {
		t.Errorf("expected velocity X = 200, got %f", vel.X)
	}
}

func TestGrounded_YPositionUnchanged(t *testing.T) {
	cfg := GroundedConfig{
		Speed:     100.0,
		MinX:      0,
		MaxX:      1000,
		FloorY:    500,
		BossWidth: 100,
	}
	movement := NewGrounded(cfg)

	startPos := types.NewVec2(400, 500)
	newPos := movement.Update(startPos, 1.0)

	if newPos.Y != startPos.Y {
		t.Errorf("expected Y unchanged, got %f -> %f", startPos.Y, newPos.Y)
	}
}
