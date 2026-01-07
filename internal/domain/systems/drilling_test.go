package systems

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

func TestVerticalDrilling_StartsAnimation(t *testing.T) {
	w := world.NewWorld(7680, 64000, 640, 42)
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place dirt tile below player
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewTile(entities.TileTypeDirt))

	// Start drilling
	inputState := input.InputState{Drill: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	// Animation should be active
	if !player.IsDrilling {
		t.Error("Drilling animation should be active after ProcessDrilling")
	}
	if !drillingSystem.animation.Active {
		t.Error("Internal animation state should be active")
	}
	if drillingSystem.animation.Duration <= 0 {
		t.Error("Animation duration should be positive")
	}
}

func TestVerticalDrilling_DirtDuration(t *testing.T) {
	w := world.NewWorld(7680, 64000, 640, 42)
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place dirt at ground level
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewTile(entities.TileTypeDirt))

	inputState := input.InputState{Drill: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	// Dirt at ground level should take 0.8 seconds
	if drillingSystem.animation.Duration != 0.8 {
		t.Errorf("Dirt at ground level should take 0.8s, got %f", drillingSystem.animation.Duration)
	}
}

func TestOreDrilling_AppliesHardnessMultiplier(t *testing.T) {
	oreTests := []struct {
		oreType  entities.OreType
		expected float32
	}{
		{entities.OreCopper, 0.96},   // 0.8 * 1.2
		{entities.OreIron, 1.2},      // 0.8 * 1.5
		{entities.OreGold, 1.44},     // 0.8 * 1.8
		{entities.OreMythril, 1.68},  // 0.8 * 2.1
		{entities.OrePlatinum, 2.0},  // 0.8 * 2.5
		{entities.OreDiamond, 2.4},   // 0.8 * 3.0
	}

	for _, test := range oreTests {
		// Reset for each ore type
		w2 := world.NewWorld(7680, 64000, 640, 42)
		player2 := entities.NewPlayer(100, 500)
		player2.OnGround = true
		ds := NewDrillingSystem(w2)

		playerCenterX := player2.AABB.X + player2.AABB.Width/2
		playerBottomY := player2.AABB.Y + player2.AABB.Height
		tileX := int(playerCenterX / world.TileSize)
		tileY := int(playerBottomY / world.TileSize)
		w2.SetTile(tileX, tileY, entities.NewOreTile(test.oreType))

		inputState := input.InputState{Drill: true}
		ds.ProcessDrilling(player2, inputState, 0.01)

		// Use tolerance-based comparison for floats
		const tolerance = 0.001
		if ds.animation.Duration < test.expected-tolerance || ds.animation.Duration > test.expected+tolerance {
			t.Errorf("Ore %v at ground level: expected ~%f seconds, got %f",
				test.oreType, test.expected, ds.animation.Duration)
		}
	}
}

func TestDrilling_DepthAffectsDuration(t *testing.T) {
	w := world.NewWorld(7680, 64000, 640, 42)
	drillingSystem := NewDrillingSystem(w)

	depthTests := []struct {
		tileGridY int
		minExpect float32
		maxExpect float32
	}{
		{10, 0.7, 0.9},       // Near ground (Y=640): ~0.8s
		{500, 15.0, 15.8},    // Mid-depth (Y=32000): ~15.4s
		{990, 29.5, 30.5},    // Deep (Y=63360): ~30s
	}

	for _, test := range depthTests {
		tileY := float32(test.tileGridY) * world.TileSize
		tile := entities.NewTile(entities.TileTypeDirt)
		duration := drillingSystem.calculateDrillingDuration(tileY, tile)

		if duration < test.minExpect || duration > test.maxExpect {
			t.Errorf("Grid Y=%d (pixel Y=%f): expected ~[%f, %f], got %f",
				test.tileGridY, tileY, test.minExpect, test.maxExpect, duration)
		}
	}
}

func TestHorizontalDrilling_CollectsOre(t *testing.T) {
	w := world.NewWorld(7680, 64000, 640, 42)
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place ore tile to the left
	playerCenterY := player.AABB.Y + player.AABB.Height/2
	tileX := int((player.AABB.X - 1) / world.TileSize)
	tileY := int(playerCenterY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTile(entities.OreDiamond))

	// Drill left (start animation)
	inputState := input.InputState{Left: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	if !player.IsDrilling {
		t.Error("Drilling animation should be active")
	}

	// Verify animation duration is correct for diamond (0.8 * 3.0 = 2.4)
	if drillingSystem.animation.Duration != 2.4 {
		t.Errorf("Diamond ore should take 2.4s, got %f", drillingSystem.animation.Duration)
	}

	// Complete animation
	dt := drillingSystem.animation.Duration + 0.01
	drillingSystem.ProcessDrilling(player, inputState, dt)

	// Should collect diamond
	if player.OreInventory[entities.OreDiamond] != 1 {
		t.Errorf("Expected 1 diamond collected, got %d", player.OreInventory[entities.OreDiamond])
	}

	// Animation should be complete
	if player.IsDrilling {
		t.Error("IsDrilling should be false after animation completes")
	}
}

func TestDrilling_DoesNotStartOnNonDrillableTile(t *testing.T) {
	w := world.NewWorld(7680, 64000, 640, 42)
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place empty tile below player (no tile at all)
	// This should prevent drilling from starting

	inputState := input.InputState{Drill: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	if player.IsDrilling {
		t.Error("Drilling should not start on empty/non-drillable tile")
	}
}

func TestDrilling_AnimationProgress(t *testing.T) {
	w := world.NewWorld(7680, 64000, 640, 42)
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place ore to the right
	playerCenterY := player.AABB.Y + player.AABB.Height/2
	tileX := int((player.AABB.X + player.AABB.Width + 1) / world.TileSize)
	tileY := int(playerCenterY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTile(entities.OreIron))

	startX := player.AABB.X

	// Start drilling right
	inputState := input.InputState{Right: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	if !player.IsDrilling {
		t.Error("Drilling animation should be active")
	}

	duration := drillingSystem.animation.Duration

	// Advance animation halfway
	drillingSystem.ProcessDrilling(player, inputState, duration/2)

	// Player should have moved toward target
	if player.AABB.X <= startX {
		t.Error("Player X position should increase during rightward drill")
	}

	// Complete the animation
	remainingTime := duration/2 + 0.01
	drillingSystem.ProcessDrilling(player, inputState, remainingTime)

	// Should be at target position now
	if player.IsDrilling {
		t.Error("Animation should be complete")
	}
}

func TestDrilling_TileRemovedOnCompletion(t *testing.T) {
	w := world.NewWorld(7680, 64000, 640, 42)
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place gold ore below player
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTile(entities.OreGold))

	// Verify tile exists before drilling
	tileBeforeDrilling := w.GetTileAtGrid(tileX, tileY)
	if tileBeforeDrilling == nil {
		t.Error("Tile should exist before drilling")
	}

	// Start and complete drilling
	inputState := input.InputState{Drill: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)
	dt := drillingSystem.animation.Duration + 0.01
	drillingSystem.ProcessDrilling(player, inputState, dt)

	// Tile should be removed
	tileAfterDrilling := w.GetTileAtGrid(tileX, tileY)
	if tileAfterDrilling != nil {
		t.Error("Tile should be removed after drilling completes")
	}
}

func TestDrilling_DoesNotCollectDirt(t *testing.T) {
	w := world.NewWorld(7680, 64000, 640, 42)
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place dirt below player
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewTile(entities.TileTypeDirt))

	// Record initial ore count
	initialTotal := 0
	for _, count := range player.OreInventory {
		initialTotal += count
	}

	// Start and complete drilling
	inputState := input.InputState{Drill: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)
	dt := drillingSystem.animation.Duration + 0.01
	drillingSystem.ProcessDrilling(player, inputState, dt)

	// Check inventory - should not have changed (dirt not collected)
	finalTotal := 0
	for _, count := range player.OreInventory {
		finalTotal += count
	}

	if finalTotal != initialTotal {
		t.Errorf("Dirt should not be collected, but inventory changed from %d to %d", initialTotal, finalTotal)
	}

	// But tile should still be removed
	if w.GetTileAtGrid(tileX, tileY) != nil {
		t.Error("Dirt tile should still be removed from world")
	}
}

func TestDrilling_SkipsInputWhileAnimating(t *testing.T) {
	w := world.NewWorld(7680, 64000, 640, 42)
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place ore below and to the right
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTile(entities.OreIron))

	// Start vertical drilling
	inputState := input.InputState{Drill: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	if !player.IsDrilling {
		t.Error("Should start drilling")
	}

	// While drilling is active, try to start a different drill direction
	// (right drilling) - it should be ignored
	inputState = input.InputState{Right: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	// The original animation should still be progressing (vertical)
	if drillingSystem.animation.Direction != DrillDown {
		t.Error("Direction should remain DrillDown while animation is active")
	}
}
