package physics_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/physics"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

func TestGetOccupiedTileRange(t *testing.T) {
	tests := []struct {
		name                          string
		aabb                          types.AABB
		tileSize                      float32
		expectedMinX, expectedMaxX    int
		expectedMinY, expectedMaxY    int
	}{
		{
			name:         "Single tile",
			aabb:         types.NewAABB(10, 10, 32, 32),
			tileSize:     64,
			expectedMinX: 0, expectedMaxX: 0,
			expectedMinY: 0, expectedMaxY: 0,
		},
		{
			name:         "Spanning 2x2 tiles",
			aabb:         types.NewAABB(32, 32, 64, 64),
			tileSize:     64,
			expectedMinX: 0, expectedMaxX: 1,
			expectedMinY: 0, expectedMaxY: 1,
		},
		{
			name:         "Exact tile boundary",
			aabb:         types.NewAABB(64, 64, 64, 64),
			tileSize:     64,
			expectedMinX: 1, expectedMaxX: 1,
			expectedMinY: 1, expectedMaxY: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minX, maxX, minY, maxY := physics.GetOccupiedTileRange(tt.aabb, tt.tileSize)

			if minX != tt.expectedMinX || maxX != tt.expectedMaxX ||
				minY != tt.expectedMinY || maxY != tt.expectedMaxY {
				t.Errorf("Expected range (%d-%d, %d-%d), got (%d-%d, %d-%d)",
					tt.expectedMinX, tt.expectedMaxX, tt.expectedMinY, tt.expectedMaxY,
					minX, maxX, minY, maxY)
			}
		})
	}
}

func TestCheckCollisions_NoCollisions(t *testing.T) {
	w := world.NewWorld(1280, 720, 640, 42)
	playerAABB := types.NewAABB(100, 100, 54, 54) // Above ground, in air

	collisions := physics.CheckCollisions(playerAABB, w)

	if len(collisions) != 0 {
		t.Errorf("Expected no collisions, got %d", len(collisions))
	}
}

func TestCheckCollisions_GroundCollision(t *testing.T) {
	w := world.NewWorld(1280, 720, 640, 42)
	// Player overlapping ground tiles (ground at Y=640, tiles start at grid Y=10)
	playerAABB := types.NewAABB(100, 620, 54, 54) // Bottom at 674

	collisions := physics.CheckCollisions(playerAABB, w)

	if len(collisions) == 0 {
		t.Error("Expected ground collision, got none")
	}
}

func TestResolveCollisionsY_GroundLanding(t *testing.T) {
	w := world.NewWorld(1280, 720, 640, 42)

	// Player falling onto ground (ground at 640)
	aabb := types.NewAABB(100, 620, 54, 54)
	velocity := types.Vec2{X: 0, Y: 100}

	collisions := physics.CheckCollisions(aabb, w)

	if len(collisions) == 0 {
		t.Fatalf("Expected to find collisions, but found none. Player AABB: %+v", aabb)
	}

	newAABB, newVel, onGround := physics.ResolveCollisionsY(aabb, velocity, collisions)

	if !onGround {
		t.Error("Expected OnGround to be true")
	}
	if newVel.Y != 0 {
		t.Errorf("Expected Y velocity zeroed, got %f", newVel.Y)
	}
	// Player should be pushed up to sit on top of tile
	expectedY := float32(640.0 - 54.0) // Ground at 640, player height 54
	if newAABB.Y != expectedY {
		t.Errorf("Expected Y=%f, got %f", expectedY, newAABB.Y)
	}
}

func TestResolveCollisionsX_NoCollisions(t *testing.T) {
	aabb := types.NewAABB(100, 100, 54, 54)
	velocity := types.Vec2{X: 50, Y: 0}
	var collisions []physics.TileCollision

	newAABB, newVel := physics.ResolveCollisionsX(aabb, velocity, collisions)

	// With no collisions, AABB and velocity should be unchanged
	if newAABB != aabb {
		t.Error("Expected AABB unchanged with no collisions")
	}
	if newVel != velocity {
		t.Error("Expected velocity unchanged with no collisions")
	}
}

func TestResolveCollisionsY_NoCollisions(t *testing.T) {
	aabb := types.NewAABB(100, 100, 54, 54)
	velocity := types.Vec2{X: 0, Y: 50}
	var collisions []physics.TileCollision

	newAABB, newVel, onGround := physics.ResolveCollisionsY(aabb, velocity, collisions)

	// With no collisions, AABB and velocity should be unchanged
	if newAABB != aabb {
		t.Error("Expected AABB unchanged with no collisions")
	}
	if newVel != velocity {
		t.Error("Expected velocity unchanged with no collisions")
	}
	if onGround {
		t.Error("Expected OnGround to be false with no collisions")
	}
}
