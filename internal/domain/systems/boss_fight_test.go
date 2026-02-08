package systems

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/components"
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// stubBoss implements bosses.Boss with controllable fields for testing.
type stubBoss struct {
	active          bool
	defeated        bool
	hp              float32
	maxHP           float32
	position        types.Vec2
	damageable      components.Damageable
	collisionBoxes  []bosses.CollisionBox
	hitboxes        []bosses.Hitbox
	hurtboxes       []bosses.Hurtbox
	updateCalled    bool
	activateCalled  int
	spawnRequests   []projectiles.SpawnRequest
	lastDamageAt    string
	lastDamageValue float32
}

func newStubBoss() *stubBoss {
	return &stubBoss{
		hp:         100,
		maxHP:      100,
		damageable: components.NewDamageable(100, 100),
	}
}

func (s *stubBoss) Update(_ *entities.Player, _ float32) []projectiles.SpawnRequest {
	s.updateCalled = true
	return s.spawnRequests
}

func (s *stubBoss) Activate()                             { s.active = true; s.activateCalled++ }
func (s *stubBoss) Deactivate()                           { s.active = false }
func (s *stubBoss) IsActive() bool                        { return s.active }
func (s *stubBoss) IsDefeated() bool                      { return s.defeated }
func (s *stubBoss) GetHP() float32                        { return s.hp }
func (s *stubBoss) GetMaxHP() float32                     { return s.maxHP }
func (s *stubBoss) GetPosition() types.Vec2               { return s.position }
func (s *stubBoss) GetDamageable() *components.Damageable { return &s.damageable }

func (s *stubBoss) GetCollisionBoxes() []bosses.CollisionBox { return s.collisionBoxes }
func (s *stubBoss) GetHitboxes() []bosses.Hitbox             { return s.hitboxes }
func (s *stubBoss) GetHurtboxes() []bosses.Hurtbox           { return s.hurtboxes }

func (s *stubBoss) TakeDamageAt(hurtboxID string, baseDamage float32) float32 {
	s.lastDamageAt = hurtboxID
	s.lastDamageValue = baseDamage
	s.hp -= baseDamage
	return baseDamage
}

func testBossFightPlayer(x, y, hp float32) *entities.Player {
	playerCfg := config.PlayerConfig{
		StartingMoney:    0,
		StartingItems:    [5]int{},
		StartingUpgrades: config.StartingUpgrades{},
	}
	upgradeCfg := config.UpgradeConfig{
		Engines:     []config.UpgradeTier[config.EngineStats]{{Stats: config.EngineStats{MaxSpeed: 450, Acceleration: 2500, FlyAcceleration: 2500, MaxUpwardSpeed: -600}}},
		Hulls:       []config.UpgradeTier[config.HullStats]{{Stats: config.HullStats{MaxHP: hp}}},
		FuelTanks:   []config.UpgradeTier[config.FuelTankStats]{{Stats: config.FuelTankStats{Capacity: 10}}},
		CargoHolds:  []config.UpgradeTier[config.CargoHoldStats]{{Stats: config.CargoHoldStats{Capacity: 10}}},
		HeatShields: []config.UpgradeTier[config.HeatShieldStats]{{Stats: config.HeatShieldStats{HeatResistance: 50}}},
		Drills:      []config.UpgradeTier[config.DrillStats]{{Stats: config.DrillStats{SpeedAtSurface: 1.0, SpeedAtMaxDepth: 1.0}}},
	}
	p := entities.NewPlayerFromConfig(x, y, playerCfg, upgradeCfg)
	return p
}

func testBossRoomCfg() config.BossRoomConfig {
	return config.BossRoomConfig{
		BossType:    "test",
		FloorType:   config.FloorConcrete,
		FloorDamage: 5,
		RoomHeight:  680,
		FloorHeight: 6,
	}
}

// With worldHeight=10000, FloorHeight=6 tiles (TileSize=64 => 384 pixels):
//
//	floorEndY   = 10000
//	floorStartY = 10000 - 384 = 9616
//	bossRoomEndY   = 9616
//	bossRoomStartY = 9616 - 680 = 8936
const testBossWorldHeight float32 = 10000

func TestNewBossFightSystem_ReturnsNilForNilBoss(t *testing.T) {
	sys := NewBossFightSystem(nil, testBossRoomCfg(), testBossWorldHeight)
	if sys != nil {
		t.Error("Expected nil system for nil boss")
	}
}

func TestIsPlayerInBossRoom_InsideRoom(t *testing.T) {
	boss := newStubBoss()
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	// Player center Y at 9000, which is between 8936 and 9616
	player := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)

	if sys.IsPlayerInBossRoom(player) == false {
		t.Error("Expected player to be in boss room at center Y=9000")
	}
}

func TestIsPlayerInBossRoom_OutsideRoom(t *testing.T) {
	boss := newStubBoss()
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	// Player center Y at 8000, well above bossRoomStartY (8936)
	player := testBossFightPlayer(100, 8000-entities.PlayerHeight/2, 100)

	if sys.IsPlayerInBossRoom(player) {
		t.Error("Expected player to be outside boss room at center Y=8000")
	}
}

func TestIsPlayerInBossRoom_AtBoundaries(t *testing.T) {
	boss := newStubBoss()
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	// Exactly at bossRoomStartY (8936) - should be in room (>=)
	playerAtStart := testBossFightPlayer(100, 8936-entities.PlayerHeight/2, 100)
	if sys.IsPlayerInBossRoom(playerAtStart) == false {
		t.Error("Expected player at bossRoomStartY boundary to be in room")
	}

	// Exactly at bossRoomEndY (9616) - should be out (<)
	playerAtEnd := testBossFightPlayer(100, 9616-entities.PlayerHeight/2, 100)
	if sys.IsPlayerInBossRoom(playerAtEnd) {
		t.Error("Expected player at bossRoomEndY boundary to be outside room")
	}
}

func TestBossFight_ActivatesBossOnRoomEntry(t *testing.T) {
	boss := newStubBoss()
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	player := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)

	sys.Update(player, 0.016)

	if boss.active == false {
		t.Error("Expected boss to be activated on room entry")
	}
}

func TestBossFight_DeactivatesBossOnRoomExit(t *testing.T) {
	boss := newStubBoss()
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	// First, enter the room
	playerIn := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)
	sys.Update(playerIn, 0.016)

	// Then leave
	playerOut := testBossFightPlayer(100, 8000-entities.PlayerHeight/2, 100)
	sys.Update(playerOut, 0.016)

	if boss.active {
		t.Error("Expected boss to be deactivated on room exit")
	}
}

func TestBossFight_DoesNotReactivateOnConsecutiveFrames(t *testing.T) {
	boss := newStubBoss()
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	player := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)

	sys.Update(player, 0.016)
	sys.Update(player, 0.016)
	sys.Update(player, 0.016)

	if boss.activateCalled != 1 {
		t.Errorf("Expected Activate() called once, got %d", boss.activateCalled)
	}
}

func TestBossFight_CallsBossUpdateWhenActive(t *testing.T) {
	boss := newStubBoss()
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	player := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)
	sys.Update(player, 0.016) // activates boss

	boss.updateCalled = false // reset after activation frame
	sys.Update(player, 0.016)

	if boss.updateCalled == false {
		t.Error("Expected boss.Update() to be called when active")
	}
}

func TestBossFight_DoesNotCallBossUpdateWhenInactive(t *testing.T) {
	boss := newStubBoss()
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	// Player outside room
	player := testBossFightPlayer(100, 8000-entities.PlayerHeight/2, 100)
	sys.Update(player, 0.016)

	if boss.updateCalled {
		t.Error("Expected boss.Update() not to be called when inactive")
	}
}

func TestBossFight_PropagatesSpawnRequests(t *testing.T) {
	boss := newStubBoss()
	boss.spawnRequests = []projectiles.SpawnRequest{
		{Position: types.NewVec2(1, 2), Damage: 5},
		{Position: types.NewVec2(3, 4), Damage: 10},
	}
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	player := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)
	result := sys.Update(player, 0.016)

	if len(result.SpawnRequests) != 2 {
		t.Fatalf("Expected 2 spawn requests, got %d", len(result.SpawnRequests))
	}
	if result.SpawnRequests[0].Damage != 5 {
		t.Errorf("Expected first spawn request damage 5, got %f", result.SpawnRequests[0].Damage)
	}
}

func TestBossFight_ContactDamage_AppliesHitboxDamage(t *testing.T) {
	boss := newStubBoss()
	// Place hitbox overlapping with player
	boss.hitboxes = []bosses.Hitbox{
		{ID: "body", X: 90, Y: 8970, Width: 100, Height: 100, DamagePerSec: 20},
	}
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	// Player overlapping hitbox
	player := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)
	initialHP := player.HP

	sys.Update(player, 0.5) // 0.5 seconds

	expectedDamage := float32(20 * 0.5) // DamagePerSec * dt
	expectedHP := initialHP - expectedDamage
	if player.HP != expectedHP {
		t.Errorf("Expected player HP %f, got %f", expectedHP, player.HP)
	}
}

func TestBossFight_ContactDamage_NoOverlapNoDamage(t *testing.T) {
	boss := newStubBoss()
	// Hitbox far from player
	boss.hitboxes = []bosses.Hitbox{
		{ID: "body", X: 5000, Y: 5000, Width: 100, Height: 100, DamagePerSec: 20},
	}
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	player := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)
	initialHP := player.HP

	sys.Update(player, 0.5)

	if player.HP != initialHP {
		t.Errorf("Expected no damage (no overlap), HP should be %f, got %f", initialHP, player.HP)
	}
}

func TestBossFight_FloorDamage_LavaAppliesDamage(t *testing.T) {
	boss := newStubBoss()
	lavaCfg := testBossRoomCfg()
	lavaCfg.FloorType = config.FloorLava
	lavaCfg.FloorDamage = 5
	sys := NewBossFightSystem(boss, lavaCfg, testBossWorldHeight)

	// floorStartY=9616, floorEndY=10000
	// Player standing on floor: bottom Y within floor range, on ground
	player := testBossFightPlayer(100, 9616-entities.PlayerHeight, 100)
	player.OnGround = true
	// Player bottom = 9616 - PlayerHeight + PlayerHeight = 9616, which is >= floorStartY
	initialHP := player.HP

	// Enter room first to activate boss
	player.AABB.Y = 9000 - entities.PlayerHeight/2
	sys.Update(player, 0.016) // activate

	// Move to floor
	player.AABB.Y = 9616 - entities.PlayerHeight
	player.OnGround = true
	// Player bottom = 9616, which is >= floorStartY(9616) and <= floorEndY(10000)
	// But player center is at 9616 - PlayerHeight/2 = 9589, which is in the room
	sys.Update(player, 0.016)

	if player.HP >= initialHP {
		t.Error("Expected lava floor damage to reduce player HP")
	}
}

func TestBossFight_FloorDamage_ConcreteNoDamage(t *testing.T) {
	boss := newStubBoss()
	concreteCfg := testBossRoomCfg()
	concreteCfg.FloorType = config.FloorConcrete
	sys := NewBossFightSystem(boss, concreteCfg, testBossWorldHeight)

	// Player on floor area, but concrete floor
	player := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)
	player.OnGround = true
	initialHP := player.HP

	sys.Update(player, 0.016)

	if player.HP != initialHP {
		t.Errorf("Expected no damage on concrete floor, HP should be %f, got %f", initialHP, player.HP)
	}
}

func TestBossFight_GameState_Victory(t *testing.T) {
	boss := newStubBoss()
	boss.defeated = true
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	player := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)
	result := sys.Update(player, 0.016)

	if result.GameState != entities.GameStateVictory {
		t.Errorf("Expected GameStateVictory, got %d", result.GameState)
	}
}

func TestBossFight_GameState_Defeat(t *testing.T) {
	boss := newStubBoss()
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	player := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)
	player.HP = 0

	result := sys.Update(player, 0.016)

	if result.GameState != entities.GameStateDefeat {
		t.Errorf("Expected GameStateDefeat, got %d", result.GameState)
	}
}

func TestBossFight_GameState_Playing(t *testing.T) {
	boss := newStubBoss()
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	player := testBossFightPlayer(100, 9000-entities.PlayerHeight/2, 100)
	result := sys.Update(player, 0.016)

	if result.GameState != entities.GameStatePlaying {
		t.Errorf("Expected GameStatePlaying, got %d", result.GameState)
	}
}

func TestBossFight_DamageBoss(t *testing.T) {
	boss := newStubBoss()
	boss.hurtboxes = []bosses.Hurtbox{
		{ID: "weak_point", X: 0, Y: 0, Width: 50, Height: 50, DamageMultiplier: 1.0},
	}
	sys := NewBossFightSystem(boss, testBossRoomCfg(), testBossWorldHeight)

	sys.DamageBoss(25)

	if boss.lastDamageAt != "weak_point" {
		t.Errorf("Expected damage at 'weak_point', got %q", boss.lastDamageAt)
	}
	if boss.lastDamageValue != 25 {
		t.Errorf("Expected damage value 25, got %f", boss.lastDamageValue)
	}
}
