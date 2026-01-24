package physics_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/physics"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// Test helper - creates minimal world for collision tests
func testWorld() *world.World {
	worldCfg := config.WorldConfig{
		Width:       1280,
		Height:      720,
		GroundLevel: 640,
		Seed:        42,
		PlayerSpawn: config.PlayerSpawn{X: 640, Y: 630},
		BuildingLayout: config.BuildingLayout{
			HospitalX: 0, FuelStationX: 0, MarketX: 0, UpgradeShopX: 0, ItemShopX: 0,
		},
	}
	genCfg := config.GenerationConfig{
		Empty:        config.TileDistribution{PeakDepth: 0, Sigma: 1000, MaxWeight: 20},
		Dirt:         config.TileDistribution{PeakDepth: 0, Sigma: 500, MaxWeight: 100},
		DirtHardness: 1.0,
		Ores:         []config.OreConfig{{ID: "copper", Name: "Copper", Value: 25, Hardness: 1.2, Distribution: config.TileDistribution{PeakDepth: -75, Sigma: 120, MaxWeight: 8}, Color: [4]uint8{184, 115, 51, 255}}},
		Hazards:      []config.HazardConfig{},
	}
	bossRoomCfg := config.BossRoomConfig{
		BossType:    "test_boss",
		FloorType:   config.FloorConcrete,
		RoomHeight:  680.0,
		FloorHeight: 6.0,
	}
	return world.NewWorldFromConfig(worldCfg, genCfg, bossRoomCfg)
}

func TestGetOccupiedTileRange(t *testing.T) {
	tests := []struct {
		name                       string
		aabb                       types.AABB
		tileSize                   float32
		expectedMinX, expectedMaxX int
		expectedMinY, expectedMaxY int
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
	w := testWorld()
	playerAABB := types.NewAABB(100, 100, 54, 54) // Above ground, in air

	collisions := physics.CheckCollisions(playerAABB, w)

	if len(collisions) != 0 {
		t.Errorf("Expected no collisions, got %d", len(collisions))
	}
}

func TestCheckCollisions_GroundCollision(t *testing.T) {
	w := testWorld()
	// Player overlapping ground tiles (ground at Y=640, tiles start at grid Y=10)
	playerAABB := types.NewAABB(100, 620, 54, 54) // Bottom at 674

	collisions := physics.CheckCollisions(playerAABB, w)

	if len(collisions) == 0 {
		t.Error("Expected ground collision, got none")
	}
}

func TestResolveCollisionsY_GroundLanding(t *testing.T) {
	w := testWorld()

	// Player falling onto ground (ground at 640)
	aabb := types.NewAABB(100, 620, 54, 54)
	velocity := types.Vec2{X: 0, Y: 100}

	collisions := physics.CheckCollisions(aabb, w)

	if len(collisions) == 0 {
		t.Fatalf("Expected to find collisions, but found none. Player AABB: %+v", aabb)
	}

	newAABB, newVel, onGround := physics.ResolveCollisionsY(aabb, velocity, collisions)

	if onGround == false {
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
