package systems

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

func TestVerticalDrilling_StartsAnimation(t *testing.T) {
	w := testWorld()
	player := testPlayer()
	player.OnGround = true
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

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
	if player.IsDrilling == false {
		t.Error("Drilling animation should be active after ProcessDrilling")
	}
	if drillingSystem.animation.Active == false {
		t.Error("Internal animation state should be active")
	}
	if drillingSystem.animation.Duration <= 0 {
		t.Error("Animation duration should be positive")
	}
}

func TestVerticalDrilling_DirtDuration(t *testing.T) {
	w := testWorld()
	player := testPlayer()
	player.OnGround = true
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

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
		w2 := testWorld()
		player2 := testPlayer()
		player2.OnGround = true
		genCfg := testGenerationConfig()
		ds := NewDrillingSystemWithConfig(w2, genCfg, testDrillingConfig())

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
	w := testWorld()
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

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
	w := testWorld()
	player := testPlayer()
	player.OnGround = true
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

	// Place ore tile to the left
	playerCenterY := player.AABB.Y + player.AABB.Height/2
	tileX := int((player.AABB.X - 1) / world.TileSize)
	tileY := int(playerCenterY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTileByID("diamond"))

	// Drill left (start animation)
	inputState := input.InputState{Left: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	if player.IsDrilling == false {
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
	w := testWorld()
	player := testPlayer()
	player.OnGround = true
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

	// Place empty tile below player (no tile at all)
	// This should prevent drilling from starting

	inputState := input.InputState{Drill: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	if player.IsDrilling {
		t.Error("Drilling should not start on empty/non-drillable tile")
	}
}

func TestDrilling_AnimationProgress(t *testing.T) {
	w := testWorld()
	player := testPlayer()
	player.OnGround = true
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

	// Place ore to the right
	playerCenterY := player.AABB.Y + player.AABB.Height/2
	tileX := int((player.AABB.X + player.AABB.Width + 1) / world.TileSize)
	tileY := int(playerCenterY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTileByID("iron"))

	startX := player.AABB.X

	// Start drilling right
	inputState := input.InputState{Right: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	if player.IsDrilling == false {
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
	w := testWorld()
	player := testPlayer()
	player.OnGround = true
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

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
	w := testWorld()
	player := testPlayer()
	player.OnGround = true
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

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
	w := testWorld()
	player := testPlayer()
	player.OnGround = true
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

	// Place ore below and to the right
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height
	tileX := int(playerCenterX / world.TileSize)
	tileY := int(playerBottomY / world.TileSize)
	w.SetTile(tileX, tileY, entities.NewOreTileByID("iron"))

	// Start vertical drilling
	inputState := input.InputState{Drill: true}
	drillingSystem.ProcessDrilling(player, inputState, 0.01)

	if player.IsDrilling == false {
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
	w := testWorld()
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

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
		w2 := testWorld()
		player := testPlayer()
		player.OnGround = true
		genCfg := testGenerationConfig()
		ds := NewDrillingSystemWithConfig(w2, genCfg, testDrillingConfig())

		tileY := float32(tileGridY) * world.TileSize
		lavaTile := entities.NewHazardTileByID("lava", nil)
		duration := ds.calculateDrillingDuration(tileY, lavaTile)

		if duration != 0.3 {
			t.Errorf("Lava at depth %d should drill in 0.3s, got %f", tileGridY, duration)
		}
	}
}

func TestDrilling_RockTileBlocksDrilling(t *testing.T) {
	w := testWorld()
	player := testPlayer()
	player.OnGround = true
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

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
	w := testWorld()
	player := testPlayer()
	player.OnGround = true
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

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
	drillEffects := drillingSystem.ProcessDrilling(player, inputState, dt)

	// Apply the returned effects
	ctx := &effects.EffectContext{Player: player, World: w}
	for _, effect := range drillEffects {
		effect.Apply(ctx)
	}

	// Player should take damage
	if player.HP >= initialHP {
		t.Errorf("Player HP should decrease after drilling lava, was %f, now %f", initialHP, player.HP)
	}
}

func TestDrilling_LavaTileRemovedAfterDrilling(t *testing.T) {
	w := testWorld()
	player := testPlayer()
	player.OnGround = true
	genCfg := testGenerationConfig()
	drillingSystem := NewDrillingSystemWithConfig(w, genCfg, testDrillingConfig())

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

// Test helpers

func testWorldConfig() config.WorldConfig {
	return config.WorldConfig{
		Width:       7680,
		Height:      64000,
		GroundLevel: 640,
		Seed:        42,
		PlayerSpawn: config.PlayerSpawn{X: 100, Y: 500},
		BuildingLayout: config.BuildingLayout{
			HospitalX: 0, FuelStationX: 0, MarketX: 0, UpgradeShopX: 0, ItemShopX: 0,
		},
	}
}

func testBossRoomConfig() config.BossRoomConfig {
	return config.BossRoomConfig{
		BossType:    "test_boss",
		FloorType:   config.FloorConcrete,
		RoomHeight:  680.0,
		FloorHeight: 6.0,
	}
}

func testDrillingConfig() config.DrillingConfig {
	return config.DrillingConfig{
		MinDrillingDuration:   1.0,
		MaxDrillingDuration:   24.0,
		FloorDrillingDuration: 0.5,
	}
}

func testGenerationConfig() config.GenerationConfig {
	return config.GenerationConfig{
		Empty:        config.TileDistribution{PeakDepth: 0, Sigma: 1000, MaxWeight: 20},
		Dirt:         config.TileDistribution{PeakDepth: 0, Sigma: 500, MaxWeight: 100},
		DirtHardness: 1.0,
		Ores: []config.OreConfig{
			{ID: "copper", Name: "Copper", Value: 25, Hardness: 1.2, Distribution: config.TileDistribution{PeakDepth: -75, Sigma: 120, MaxWeight: 8}, Color: [4]uint8{184, 115, 51, 255}},
			{ID: "iron", Name: "Iron", Value: 75, Hardness: 1.5, Distribution: config.TileDistribution{PeakDepth: 70, Sigma: 90, MaxWeight: 5}, Color: [4]uint8{112, 128, 144, 255}},
			{ID: "gold", Name: "Gold", Value: 300, Hardness: 1.8, Distribution: config.TileDistribution{PeakDepth: 230, Sigma: 80, MaxWeight: 3}, Color: [4]uint8{255, 215, 0, 255}},
			{ID: "mythril", Name: "Mythril", Value: 1500, Hardness: 2.1, Distribution: config.TileDistribution{PeakDepth: 360, Sigma: 70, MaxWeight: 2.2}, Color: [4]uint8{0, 191, 255, 255}},
			{ID: "platinum", Name: "Platinum", Value: 10000, Hardness: 2.5, Distribution: config.TileDistribution{PeakDepth: 500, Sigma: 80, MaxWeight: 1.8}, Color: [4]uint8{229, 228, 226, 255}},
			{ID: "diamond", Name: "Diamond", Value: 30000, Hardness: 3.0, Distribution: config.TileDistribution{PeakDepth: 600, Sigma: 180, MaxWeight: 0.15}, Color: [4]uint8{185, 242, 255, 255}},
		},
		Hazards: []config.HazardConfig{
			{ID: "rock", Name: "Rock", Drillable: false, Distribution: config.TileDistribution{PeakDepth: 650, Sigma: 200, MaxWeight: 15}, Color: [4]uint8{80, 80, 80, 255}},
			{ID: "lava", Name: "Lava", Drillable: true, FixedDuration: 0.3, OnDrillEffect: config.HazardEffectConfig{Type: config.HazardEffectHeatDamage, BaseDamage: 100, MaxHeatResistance: 320, MaxDamageReduction: 0.5}, Distribution: config.TileDistribution{PeakDepth: 750, Sigma: 150, MaxWeight: 12}, Color: [4]uint8{255, 100, 0, 255}},
		},
	}
}

func testWorld() *world.World {
	return world.NewWorldFromConfig(testWorldConfig(), testGenerationConfig(), testBossRoomConfig())
}

func testPlayer() *entities.Player {
	playerCfg := config.PlayerConfig{
		StartingMoney:    0,
		StartingItems:    [5]int{0, 0, 0, 0, 0},
		StartingUpgrades: config.StartingUpgrades{},
	}
	upgradeCfg := config.UpgradeConfig{
		Engines:     []config.UpgradeTier[config.EngineStats]{{Name: "Base", Price: 0, Stats: config.EngineStats{MaxSpeed: 450, Acceleration: 2500, FlyAcceleration: 2500, MaxUpwardSpeed: -600}}},
		Hulls:       []config.UpgradeTier[config.HullStats]{{Name: "Base", Price: 0, Stats: config.HullStats{MaxHP: 10}}},
		FuelTanks:   []config.UpgradeTier[config.FuelTankStats]{{Name: "Base", Price: 0, Stats: config.FuelTankStats{Capacity: 10}}},
		CargoHolds:  []config.UpgradeTier[config.CargoHoldStats]{{Name: "Base", Price: 0, Stats: config.CargoHoldStats{Capacity: 10}}},
		HeatShields: []config.UpgradeTier[config.HeatShieldStats]{{Name: "Base", Price: 0, Stats: config.HeatShieldStats{HeatResistance: 50}}},
		Drills:      []config.UpgradeTier[config.DrillStats]{{Name: "Base", Price: 0, Stats: config.DrillStats{SpeedAtSurface: 1.0, SpeedAtMaxDepth: 1.0}}},
	}
	return entities.NewPlayerFromConfig(100, 500, playerCfg, upgradeCfg)
}
