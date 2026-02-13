package movement

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/types"
)

func TestHovering_MovesHorizontally(t *testing.T) {
	h := NewHovering(HoveringConfig{
		Speed:        100,
		MinX:         0,
		MaxX:         800,
		HoverY:       200,
		BossWidth:    80,
		BobAmplitude: 0,
		BobFrequency: 0,
	})

	pos := types.NewVec2(400, 200)
	newPos := h.Update(pos, 0.1)

	if newPos.X <= pos.X {
		t.Errorf("Expected X to increase, got %f -> %f", pos.X, newPos.X)
	}
}

func TestHovering_BouncesAtBoundaries(t *testing.T) {
	h := NewHovering(HoveringConfig{
		Speed:        100,
		MinX:         0,
		MaxX:         800,
		HoverY:       200,
		BossWidth:    80,
		BobAmplitude: 0,
		BobFrequency: 0,
	})

	// Place boss near right boundary
	pos := types.NewVec2(750, 200)
	newPos := h.Update(pos, 1.0)

	// Should have bounced
	if h.GetDirection() != -1 {
		t.Errorf("Expected direction -1 after right boundary bounce, got %f", h.GetDirection())
	}

	if newPos.X > 720 {
		t.Errorf("Expected X clamped at right boundary, got %f", newPos.X)
	}
}

func TestHovering_BouncesAtLeftBoundary(t *testing.T) {
	h := NewHovering(HoveringConfig{
		Speed:        100,
		MinX:         0,
		MaxX:         800,
		HoverY:       200,
		BossWidth:    80,
		BobAmplitude: 0,
		BobFrequency: 0,
	})

	// Force left direction
	h.direction = -1

	pos := types.NewVec2(5, 200)
	newPos := h.Update(pos, 1.0)

	if h.GetDirection() != 1 {
		t.Errorf("Expected direction 1 after left boundary bounce, got %f", h.GetDirection())
	}

	if newPos.X < 0 {
		t.Errorf("Expected X clamped at left boundary, got %f", newPos.X)
	}
}

func TestHovering_BobsVertically(t *testing.T) {
	h := NewHovering(HoveringConfig{
		Speed:        0,
		MinX:         0,
		MaxX:         800,
		HoverY:       200,
		BossWidth:    80,
		BobAmplitude: 10,
		BobFrequency: 6.28, // ~1 Hz
	})

	pos := types.NewVec2(400, 200)

	// Collect Y positions over time
	var yPositions []float32
	for i := 0; i < 20; i++ {
		pos = h.Update(pos, 0.05)
		yPositions = append(yPositions, pos.Y)
	}

	// Check that Y varies around hoverY
	hasAbove := false
	hasBelow := false
	for _, y := range yPositions {
		if y > 200 {
			hasAbove = true
		}
		if y < 200 {
			hasBelow = true
		}
	}

	if hasAbove == false || hasBelow == false {
		t.Error("Expected vertical bobbing above and below hoverY")
	}
}

func TestHovering_PauseStopsHorizontalMovement(t *testing.T) {
	h := NewHovering(HoveringConfig{
		Speed:        100,
		MinX:         0,
		MaxX:         800,
		HoverY:       200,
		BossWidth:    80,
		BobAmplitude: 0,
		BobFrequency: 0,
	})

	h.Pause()

	pos := types.NewVec2(400, 200)
	newPos := h.Update(pos, 0.1)

	if newPos.X != pos.X {
		t.Errorf("Expected X unchanged when paused, got %f -> %f", pos.X, newPos.X)
	}

	if h.IsPaused() == false {
		t.Error("Expected IsPaused to be true")
	}
}

func TestHovering_ResumeRestoresMovement(t *testing.T) {
	h := NewHovering(HoveringConfig{
		Speed:        100,
		MinX:         0,
		MaxX:         800,
		HoverY:       200,
		BossWidth:    80,
		BobAmplitude: 0,
		BobFrequency: 0,
	})

	h.Pause()
	h.Resume()

	pos := types.NewVec2(400, 200)
	newPos := h.Update(pos, 0.1)

	if newPos.X == pos.X {
		t.Error("Expected X to change after resume")
	}

	if h.IsPaused() {
		t.Error("Expected IsPaused to be false after resume")
	}
}

func TestHovering_SetSpeed(t *testing.T) {
	h := NewHovering(HoveringConfig{
		Speed:        100,
		MinX:         0,
		MaxX:         800,
		HoverY:       200,
		BossWidth:    80,
		BobAmplitude: 0,
		BobFrequency: 0,
	})

	h.SetSpeed(200)

	if h.GetSpeed() != 200 {
		t.Errorf("Expected speed 200, got %f", h.GetSpeed())
	}
}
