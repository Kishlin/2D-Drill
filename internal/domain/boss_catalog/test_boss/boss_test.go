package test_boss

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
)

const (
	testRoomStartY = 8000.0
	testWorldWidth = 800.0
)

func testBossPlayer() *entities.Player {
	playerCfg := config.PlayerConfig{
		StartingMoney:    0,
		StartingItems:    [5]int{},
		StartingUpgrades: config.StartingUpgrades{},
	}
	upgradeCfg := config.UpgradeConfig{
		Engines:     []config.UpgradeTier[config.EngineStats]{{Stats: config.EngineStats{MaxSpeed: 450, Acceleration: 2500, FlyAcceleration: 2500, MaxUpwardSpeed: -600}}},
		Hulls:       []config.UpgradeTier[config.HullStats]{{Stats: config.HullStats{MaxHP: 200}}},
		FuelTanks:   []config.UpgradeTier[config.FuelTankStats]{{Stats: config.FuelTankStats{Capacity: 10}}},
		CargoHolds:  []config.UpgradeTier[config.CargoHoldStats]{{Stats: config.CargoHoldStats{Capacity: 10}}},
		HeatShields: []config.UpgradeTier[config.HeatShieldStats]{{Stats: config.HeatShieldStats{HeatResistance: 50}}},
		Drills:      []config.UpgradeTier[config.DrillStats]{{Stats: config.DrillStats{SpeedAtSurface: 1.0, SpeedAtMaxDepth: 1.0}}},
	}
	return entities.NewPlayerFromConfig(350, testRoomStartY+300, playerCfg, upgradeCfg)
}

func newActiveBoss() (*TestBoss, *entities.Player) {
	b := New(testRoomStartY, testWorldWidth)
	b.Activate()
	player := testBossPlayer()
	b.Update(player, 0.016) // initialize CurrentPlayer
	return b, player
}

func TestTestBoss_New_InitialState(t *testing.T) {
	b := New(testRoomStartY, testWorldWidth)

	if b.IsActive() {
		t.Error("Expected boss to start inactive")
	}
	if b.IsDefeated() {
		t.Error("Expected boss to start undefeated")
	}
	if b.GetHP() != MaxHP {
		t.Errorf("Expected HP %f, got %f", float32(MaxHP), b.GetHP())
	}
	if b.GetState() != StatePatrol {
		t.Errorf("Expected StatePatrol (%d), got %d", StatePatrol, b.GetState())
	}
	if b.GetCurrentPhase() != 1 {
		t.Errorf("Expected phase 1 (display), got %d", b.GetCurrentPhase())
	}
}

func TestTestBoss_Phase0_AlwaysVulnerable(t *testing.T) {
	b, _ := newActiveBoss()

	hurtboxes := b.GetHurtboxes()
	if len(hurtboxes) == 0 {
		t.Error("Expected non-empty hurtboxes in phase 0 (always vulnerable)")
	}

	timer := b.GetVulnerableTimer()
	if timer != -1 {
		t.Errorf("Expected vulnerable timer -1 in phase 0, got %f", timer)
	}
}

func TestTestBoss_Phase1_InvulnerableOutsideVulnerableState(t *testing.T) {
	b, player := newActiveBoss()

	// Deal damage to push into phase 1 (below 66% HP = below 66 HP)
	// Boss starts at 100 HP, deal 35 to reach 65
	b.TakeDamageAt("body", 35)

	// Run an update to trigger phase transition
	b.Update(player, 0.016)

	// In patrol state (not vulnerable), should have no hurtboxes
	if b.GetState() != StatePatrol {
		t.Fatalf("Expected StatePatrol, got %d", b.GetState())
	}

	hurtboxes := b.GetHurtboxes()
	if len(hurtboxes) != 0 {
		t.Error("Expected empty hurtboxes in phase 1 when not in vulnerable state")
	}
}

func TestTestBoss_OnDamageReceived_ExitsVulnerableState(t *testing.T) {
	b, player := newActiveBoss()

	// Push to phase 1 (need AOE to reach vulnerable state)
	b.TakeDamageAt("body", 35) // HP = 65, below 66% threshold
	b.Update(player, 0.016)

	// Force transition to vulnerable state via slam cycle
	// Set aoeCooldown to near-zero to trigger windup
	b.aoeCooldown = 0.01
	b.Update(player, 0.02) // trigger windup transition

	if b.GetState() != StateWindup {
		t.Fatalf("Expected StateWindup, got %d", b.GetState())
	}

	// Advance through windup (1.0 seconds)
	for i := 0; i < 100; i++ {
		b.Update(player, 0.012)
	}

	// Should now be in Slam or Vulnerable
	state := b.GetState()
	if state == StateWindup {
		t.Fatal("Expected to have left StateWindup after 1.2 seconds")
	}

	// Keep advancing until vulnerable
	for i := 0; i < 50; i++ {
		if b.GetState() == StateVulnerable {
			break
		}
		b.Update(player, 0.02)
	}

	if b.GetState() != StateVulnerable {
		t.Fatalf("Expected StateVulnerable, got %d", b.GetState())
	}

	// Now damage should exit vulnerable state
	b.TakeDamageAt("body", 5)

	if b.GetState() != StatePatrol {
		t.Errorf("Expected StatePatrol after damage in vulnerable state, got %d", b.GetState())
	}
}

func TestTestBoss_GetAOEInfo_NilInPatrol(t *testing.T) {
	b, _ := newActiveBoss()

	info := b.GetAOEInfo()
	if info != nil {
		t.Error("Expected nil AOEInfo in patrol state")
	}
}

func TestTestBoss_GetAOEInfo_TelegraphInWindup(t *testing.T) {
	b, player := newActiveBoss()

	// Push to phase 1 for AOE attacks
	b.TakeDamageAt("body", 35)
	b.Update(player, 0.016)

	// Trigger windup
	b.aoeCooldown = 0.01
	b.Update(player, 0.02)

	if b.GetState() != StateWindup {
		t.Fatalf("Expected StateWindup, got %d", b.GetState())
	}

	info := b.GetAOEInfo()
	if info == nil {
		t.Fatal("Expected non-nil AOEInfo in windup state")
	}
	if info.IsTelegraph == false {
		t.Error("Expected IsTelegraph=true in windup")
	}
	if info.IsDamaging {
		t.Error("Expected IsDamaging=false in windup")
	}
}

func TestTestBoss_GetAOEInfo_DamagingInSlam(t *testing.T) {
	b, player := newActiveBoss()

	// Push to phase 1
	b.TakeDamageAt("body", 35)
	b.Update(player, 0.016)

	// Trigger windup
	b.aoeCooldown = 0.01
	b.Update(player, 0.02)

	// Advance through windup to slam
	for i := 0; i < 100; i++ {
		b.Update(player, 0.012)
		if b.GetState() == StateSlam {
			break
		}
	}

	if b.GetState() != StateSlam {
		t.Fatalf("Expected StateSlam, got %d", b.GetState())
	}

	info := b.GetAOEInfo()
	if info == nil {
		t.Fatal("Expected non-nil AOEInfo in slam state")
	}
	if info.IsTelegraph {
		t.Error("Expected IsTelegraph=false in slam")
	}
	if info.IsDamaging == false {
		t.Error("Expected IsDamaging=true in slam")
	}
}

func TestTestBoss_GetVulnerableTimer_NegativeInPhase0(t *testing.T) {
	b, _ := newActiveBoss()

	timer := b.GetVulnerableTimer()
	if timer != -1 {
		t.Errorf("Expected -1 (always vulnerable) in phase 0, got %f", timer)
	}
}

func TestTestBoss_GetVulnerableTimer_ZeroInPatrol(t *testing.T) {
	b, player := newActiveBoss()

	// Push to phase 1
	b.TakeDamageAt("body", 35)
	b.Update(player, 0.016)

	if b.GetState() != StatePatrol {
		t.Fatalf("Expected StatePatrol, got %d", b.GetState())
	}

	timer := b.GetVulnerableTimer()
	if timer != 0 {
		t.Errorf("Expected 0 (not vulnerable) in patrol phase 1, got %f", timer)
	}
}

func TestTestBoss_Update_MovesInPatrol(t *testing.T) {
	b, player := newActiveBoss()

	posBeforeX := b.GetPosition().X

	// Several updates should cause movement in patrol
	for i := 0; i < 10; i++ {
		b.Update(player, 0.05)
	}

	posAfterX := b.GetPosition().X
	if posAfterX == posBeforeX {
		t.Error("Expected position X to change during patrol movement")
	}
}

func TestTestBoss_StateCycle(t *testing.T) {
	b, player := newActiveBoss()

	// Push to phase 1 for AOE cycle
	b.TakeDamageAt("body", 35)
	b.Update(player, 0.016)

	// 1. Start in Patrol
	if b.GetState() != StatePatrol {
		t.Fatalf("Step 1: Expected StatePatrol, got %d", b.GetState())
	}

	// 2. Trigger windup
	b.aoeCooldown = 0.01
	b.Update(player, 0.02)
	if b.GetState() != StateWindup {
		t.Fatalf("Step 2: Expected StateWindup, got %d", b.GetState())
	}

	// 3. Advance through windup to slam
	for i := 0; i < 100; i++ {
		b.Update(player, 0.012)
		if b.GetState() == StateSlam {
			break
		}
	}
	if b.GetState() != StateSlam {
		t.Fatalf("Step 3: Expected StateSlam, got %d", b.GetState())
	}

	// 4. Advance through slam to vulnerable (or windupBetween for double slam)
	for i := 0; i < 50; i++ {
		b.Update(player, 0.02)
		if b.GetState() == StateVulnerable {
			break
		}
	}
	if b.GetState() != StateVulnerable {
		t.Fatalf("Step 4: Expected StateVulnerable, got %d", b.GetState())
	}

	// 5. Wait for vulnerable to end -> back to patrol
	for i := 0; i < 200; i++ {
		b.Update(player, 0.02)
		if b.GetState() == StatePatrol {
			break
		}
	}
	if b.GetState() != StatePatrol {
		t.Fatalf("Step 5: Expected StatePatrol after vulnerability ends, got %d", b.GetState())
	}
}

func TestTestBoss_IsVulnerable_MatchesHurtboxes(t *testing.T) {
	b, _ := newActiveBoss()

	// Phase 0: always vulnerable
	if b.IsVulnerable() == false {
		t.Error("Expected IsVulnerable=true in phase 0")
	}
	if len(b.GetHurtboxes()) == 0 {
		t.Error("Expected non-empty hurtboxes when IsVulnerable=true")
	}

	// Consistency check: IsVulnerable should match hurtbox presence
	isVuln := b.IsVulnerable()
	hasHurtboxes := len(b.GetHurtboxes()) > 0
	if isVuln != hasHurtboxes {
		t.Errorf("IsVulnerable (%v) should match hurtbox presence (%v)", isVuln, hasHurtboxes)
	}
}
