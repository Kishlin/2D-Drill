package systems

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

func TestVerticalDrilling_StartsAnimation(t *testing.T) {
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
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
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
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

	// Dirt at ground level should take 1.0 seconds (with base drill, no speedup)
	if drillingSystem.animation.Duration != 1.0 {
		t.Errorf("Dirt at ground level should take 1.0s, got %f", drillingSystem.animation.Duration)
	}
}

func TestOreDrilling_AppliesHardnessMultiplier(t *testing.T) {
	oreTests := []struct {
		oreID    string
		expected float32
	}{
		{"copper", 1.2},   // 1.0 * 1.2
		{"iron", 1.5},     // 1.0 * 1.5
		{"gold", 1.8},     // 1.0 * 1.8
		{"mythril", 2.1},  // 1.0 * 2.1
		{"platinum", 2.5}, // 1.0 * 2.5
		{"diamond", 3.0},  // 1.0 * 3.0
	}

	for _, test := range oreTests {
		// Reset for each ore type
		w2 := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
		player2 := entities.NewPlayer(100, 500)
		player2.OnGround = true
		ds := NewDrillingSystem(w2)

		playerCenterX := player2.AABB.X + player2.AABB.Width/2
		playerBottomY := player2.AABB.Y + player2.AABB.Height
		tileX := int(playerCenterX / world.TileSize)
		tileY := int(playerBottomY / world.TileSize)
		w2.SetTile(tileX, tileY, entities.NewOreTileByID(test.oreID))

		inputState := input.InputState{Drill: true}
		ds.ProcessDrilling(player2, inputState, 0.01)

		// Use tolerance-based comparison for floats
		const tolerance = 0.001
		if ds.animation.Duration < test.expected-tolerance || ds.animation.Duration > test.expected+tolerance {
			t.Errorf("Ore %s at ground level: expected ~%f seconds, got %f",
				test.oreID, test.expected, ds.animation.Duration)
		}
	}
}

func TestDrilling_DepthAffectsDuration(t *testing.T) {
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
	drillingSystem := NewDrillingSystem(w)

	depthTests := []struct {
		tileGridY int
		minExpect float32
		maxExpect float32
	}{
		{10, 0.9, 1.1},    // Near ground (Y=640): ~1.0s
		{500, 12.0, 13.0}, // Mid-depth (Y=32000): ~12.4s
		{990, 23.5, 24.5}, // Deep (Y=63360): ~24s
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
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place ore tile to the left
	playerCenterY := player.AABB.Y + player.AABB.Height/2
	tileX := int((player.AABB.X - 1) / world.TileSize)
	tileY := int(playerCenterY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTileByID("diamond"))

	// Drill left (start animation)
	inputState := input.InputState{Left: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	if !player.IsDrilling {
		t.Error("Drilling animation should be active")
	}

	// Verify animation duration is correct for diamond (1.0 * 3.0 = 3.0)
	if drillingSystem.animation.Duration != 3.0 {
		t.Errorf("Diamond ore should take 3.0s, got %f", drillingSystem.animation.Duration)
	}

	// Complete animation
	dt := drillingSystem.animation.Duration + 0.01
	drillingSystem.ProcessDrilling(player, inputState, dt)

	// Should collect diamond
	if player.OreInventory["diamond"] != 1 {
		t.Errorf("Expected 1 diamond collected, got %d", player.OreInventory["diamond"])
	}

	// Animation should be complete
	if player.IsDrilling {
		t.Error("IsDrilling should be false after animation completes")
	}
}

func TestDrilling_DoesNotStartOnNonDrillableTile(t *testing.T) {
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
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
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place ore to the right
	playerCenterY := player.AABB.Y + player.AABB.Height/2
	tileX := int((player.AABB.X + player.AABB.Width + 1) / world.TileSize)
	tileY := int(playerCenterY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTileByID("iron"))

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
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place gold ore below player
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTileByID("gold"))

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
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
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
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place ore below and to the right
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTileByID("iron"))

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

// === Hazard Tile Drilling Tests ===

func TestDrilling_LavaTileDrillsQuickly(t *testing.T) {
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
	drillingSystem := NewDrillingSystem(w)

	// Test at various depths
	depths := []int{50, 200, 500, 800}
	for _, depth := range depths {
		tileY := float32(depth) * world.TileSize
		lavaTile := entities.NewHazardTileByID("lava", nil)
		duration := drillingSystem.calculateDrillingDuration(tileY, lavaTile)

		// Lava should always drill in exactly 0.3 seconds (before any floor clamp)
		if duration != 0.3 {
			t.Errorf("Lava tile at depth %d should calculate as 0.3s, got %f", depth, duration)
		}
	}
}

func TestDrilling_LavaTileDrillsQuicklyAtAnyDepth(t *testing.T) {
	// Test lava drilling at multiple depths
	depthTests := []int{50, 200, 500, 800} // Various depths

	for _, tileGridY := range depthTests {
		w2 := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
		player := entities.NewPlayer(100, 500)
		player.OnGround = true
		ds := NewDrillingSystem(w2)

		tileY := float32(tileGridY) * world.TileSize
		lavaTile := entities.NewHazardTileByID("lava", nil)
		duration := ds.calculateDrillingDuration(tileY, lavaTile)

		if duration != 0.3 {
			t.Errorf("Lava at depth %d should drill in 0.3s, got %f", tileGridY, duration)
		}
	}
}

func TestDrilling_RockTileBlocksDrilling(t *testing.T) {
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place rock tile below player (rock is not drillable)
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	rockCfg := &config.HazardConfig{Drillable: false}
	w.SetTile(tileX, tileY, entities.NewHazardTileByID("rock", rockCfg))

	// Try to drill
	inputState := input.InputState{Drill: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	// Rock is not drillable, so animation should not start
	if player.IsDrilling {
		t.Error("Drilling should not start on rock tile (not drillable)")
	}
}

func TestDrilling_LavaDealsDamage(t *testing.T) {
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place lava tile below player
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewHazardTileByID("lava", nil))

	initialHP := player.HP

	// Start and complete drilling
	inputState := input.InputState{Drill: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)
	dt := drillingSystem.animation.Duration + 0.01
	drillingSystem.ProcessDrilling(player, inputState, dt)

	// Player should take damage
	if player.HP >= initialHP {
		t.Errorf("Player HP should decrease after drilling lava, was %f, now %f", initialHP, player.HP)
	}
}

func TestDrilling_LavaTileRemovedAfterDrilling(t *testing.T) {
	w := world.NewWorld(world.NewWorldConfigForTesting(7680, 64000, 640, 42))
	player := entities.NewPlayer(100, 500)
	player.OnGround = true
	drillingSystem := NewDrillingSystem(w)

	// Place lava tile below player
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewHazardTileByID("lava", nil))

	// Verify tile exists
	if w.GetTileAtGrid(tileX, tileY) == nil {
		t.Error("Lava tile should exist before drilling")
	}

	// Complete drilling
	inputState := input.InputState{Drill: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)
	dt := drillingSystem.animation.Duration + 0.01
	drillingSystem.ProcessDrilling(player, inputState, dt)

	// Tile should be removed (even though it deals damage)
	if w.GetTileAtGrid(tileX, tileY) != nil {
		t.Error("Lava tile should be removed after drilling")
	}
}
