package systems

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

func TestDiggingSystem_CollectsOre(t *testing.T) {
	// Create world with known seed
	w := world.NewWorld(7680, 64000, 640, 42)
	player := entities.NewPlayer(100, 500)
	player.OnGround = true // Required to start dig animation
	diggingSystem := NewDiggingSystem(w)

	// Manually place an ore tile below player
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTile(entities.OreGold))

	// Initial inventory should be empty
	if player.OreInventory[entities.OreGold] != 0 {
		t.Error("Player should start with 0 gold")
	}

	// Trigger digging (start animation)
	inputState := input.InputState{Dig: true}
	dt := float32(0.01)
	diggingSystem.ProcessDigging(player, inputState, dt)

	// Animation started, tile not removed yet
	if !player.IsDigging {
		t.Error("Digging animation should be active")
	}
	if player.OreInventory[entities.OreGold] != 0 {
		t.Error("Ore should not be collected until animation completes")
	}

	// Complete animation
	dt = DigAnimationDuration + 0.01 // More than animation duration
	diggingSystem.ProcessDigging(player, inputState, dt)

	// Player should have collected 1 gold
	if player.OreInventory[entities.OreGold] != 1 {
		t.Errorf("Expected 1 gold after digging, got %d", player.OreInventory[entities.OreGold])
	}

	// Tile should be gone
	tileAfter := w.GetTileAt(playerCenterX, playerBottomY)
	if tileAfter != nil {
		t.Error("Tile should be removed after digging")
	}

	// Animation should be complete
	if player.IsDigging {
		t.Error("Digging animation should be complete")
	}
}

func TestDiggingSystem_DoesNotCollectDirt(t *testing.T) {
	w := world.NewWorld(7680, 64000, 640, 42)
	player := entities.NewPlayer(100, 500)
	player.OnGround = true // Required to start dig animation
	diggingSystem := NewDiggingSystem(w)

	// Manually place a dirt tile below player
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewTile(entities.TileTypeDirt))

	// Trigger digging and complete animation
	inputState := input.InputState{Dig: true}
	dt := DigAnimationDuration + 0.01
	diggingSystem.ProcessDigging(player, inputState, dt)

	// Player should have 0 ores (dirt not collected)
	totalOre := 0
	for _, count := range player.OreInventory {
		totalOre += count
	}
	if totalOre != 0 {
		t.Error("Dirt should not be added to inventory")
	}
}

func TestHorizontalDigging_CollectsOre(t *testing.T) {
	w := world.NewWorld(7680, 64000, 640, 42)
	player := entities.NewPlayer(100, 500)
	player.OnGround = true // Required for horizontal digging
	diggingSystem := NewDiggingSystem(w)

	// Place ore tile to the left
	playerCenterY := player.AABB.Y + player.AABB.Height/2
	tileX := int((player.AABB.X - 1) / world.TileSize)
	tileY := int(playerCenterY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTile(entities.OreDiamond))

	// Dig left (start animation)
	inputState := input.InputState{Left: true}
	dt := float32(0.01)
	diggingSystem.ProcessDigging(player, inputState, dt)

	// Animation should be active
	if !player.IsDigging {
		t.Error("Digging animation should be active")
	}

	// Complete animation
	dt = DigAnimationDuration + 0.01
	diggingSystem.ProcessDigging(player, inputState, dt)

	// Should collect diamond
	if player.OreInventory[entities.OreDiamond] != 1 {
		t.Errorf("Expected 1 diamond, got %d", player.OreInventory[entities.OreDiamond])
	}
}
