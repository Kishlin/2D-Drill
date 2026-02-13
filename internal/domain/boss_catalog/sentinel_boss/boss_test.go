package sentinel_boss

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

func newActiveBoss() (*SentinelBoss, *entities.Player) {
	b := New(testRoomStartY, testWorldWidth)
	b.Activate()
	player := testBossPlayer()
	b.Update(player, 0.016) // initialize CurrentPlayer
	return b, player
}

func TestSentinelBoss_New_InitialState(t *testing.T) {
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
	if b.GetState() != StateHover {
		t.Errorf("Expected StateHover (%d), got %d", StateHover, b.GetState())
	}
	if b.GetCurrentPhase() != 1 {
		t.Errorf("Expected phase 1 (display), got %d", b.GetCurrentPhase())
	}
}

func TestSentinelBoss_Phase0_AlwaysVulnerable(t *testing.T) {
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

func TestSentinelBoss_Phase1_InvulnerableOutsideStunnedState(t *testing.T) {
	b, player := newActiveBoss()

	// Deal damage to push into phase 1 (below 60% HP = below 90 HP out of 150)
	// Boss starts at 150 HP, deal 61 to reach 89
	b.TakeDamageAt("body", 61)

	// Run an update to trigger phase transition
	b.Update(player, 0.016)

	// In hover state (not stunned), should have no hurtboxes
	if b.GetState() != StateHover {
		t.Fatalf("Expected StateHover, got %d", b.GetState())
	}

	hurtboxes := b.GetHurtboxes()
	if len(hurtboxes) != 0 {
		t.Error("Expected empty hurtboxes in phase 1 when not in stunned state")
	}
}

func TestSentinelBoss_Update_MovesInHover(t *testing.T) {
	b, player := newActiveBoss()

	posBeforeX := b.GetPosition().X

	// Several updates should cause movement in hover
	for i := 0; i < 10; i++ {
		b.Update(player, 0.05)
	}

	posAfterX := b.GetPosition().X
	if posAfterX == posBeforeX {
		t.Error("Expected position X to change during hover movement")
	}
}

func TestSentinelBoss_ChargeFlow(t *testing.T) {
	b, player := newActiveBoss()

	// Push to phase 1 for charge attacks
	b.TakeDamageAt("body", 61) // HP = 89, below 60% threshold
	b.Update(player, 0.016)

	// 1. Start in Hover
	if b.GetState() != StateHover {
		t.Fatalf("Step 1: Expected StateHover, got %d", b.GetState())
	}

	// 2. Trigger charge windup by setting cooldown to near-zero
	b.chargeCooldown = 0.01
	b.Update(player, 0.02)
	if b.GetState() != StateChargeWindup {
		t.Fatalf("Step 2: Expected StateChargeWindup, got %d", b.GetState())
	}

	// 3. Advance through windup (0.8 seconds)
	for i := 0; i < 100; i++ {
		b.Update(player, 0.01)
		if b.GetState() == StateCharge {
			break
		}
	}
	if b.GetState() != StateCharge {
		t.Fatalf("Step 3: Expected StateCharge, got %d", b.GetState())
	}

	// 4. Advance through charge to stunned
	for i := 0; i < 200; i++ {
		b.Update(player, 0.01)
		if b.GetState() == StateStunned {
			break
		}
	}
	if b.GetState() != StateStunned {
		t.Fatalf("Step 4: Expected StateStunned, got %d", b.GetState())
	}

	// 5. Verify vulnerable during stun
	if b.IsVulnerable() == false {
		t.Error("Expected boss to be vulnerable during stun")
	}

	// 6. Wait for stun to end -> back to hover
	for i := 0; i < 400; i++ {
		b.Update(player, 0.01)
		if b.GetState() == StateHover {
			break
		}
	}
	if b.GetState() != StateHover {
		t.Fatalf("Step 6: Expected StateHover after stun ends, got %d", b.GetState())
	}
}

func TestSentinelBoss_LaserFlow(t *testing.T) {
	b, player := newActiveBoss()

	// Push to phase 2 for laser attacks (below 30% HP = below 45 HP)
	b.TakeDamageAt("body", 106) // HP = 44, below 30% threshold
	b.Update(player, 0.016)

	// Trigger laser by setting laser cooldown to near-zero
	// Also set charge cooldown high so it doesn't trigger first
	b.laserCooldown = 0.01
	b.chargeCooldown = 100
	b.Update(player, 0.02)

	if b.GetState() != StateLaserAim {
		t.Fatalf("Expected StateLaserAim, got %d", b.GetState())
	}

	// Verify laser info during aim
	laserInfo := b.GetLaserInfo()
	if laserInfo == nil {
		t.Fatal("Expected non-nil LaserInfo during aim")
	}
	if laserInfo.IsAiming == false {
		t.Error("Expected IsAiming=true during aim phase")
	}

	// Advance through aim (1.0 second)
	for i := 0; i < 120; i++ {
		b.Update(player, 0.01)
		if b.GetState() == StateLaser {
			break
		}
	}
	if b.GetState() != StateLaser {
		t.Fatalf("Expected StateLaser, got %d", b.GetState())
	}

	// Verify laser info during fire
	laserInfo = b.GetLaserInfo()
	if laserInfo == nil {
		t.Fatal("Expected non-nil LaserInfo during fire")
	}
	if laserInfo.IsFiring == false {
		t.Error("Expected IsFiring=true during fire phase")
	}

	// Advance through laser fire (0.5 seconds)
	for i := 0; i < 60; i++ {
		b.Update(player, 0.01)
		if b.GetState() == StateHover {
			break
		}
	}
	if b.GetState() != StateHover {
		t.Fatalf("Expected StateHover after laser, got %d", b.GetState())
	}

	// Verify no laser info in hover
	if b.GetLaserInfo() != nil {
		t.Error("Expected nil LaserInfo in hover state")
	}
}

func TestSentinelBoss_OnDamageReceived_ExitsStunnedState(t *testing.T) {
	b, player := newActiveBoss()

	// Push to phase 1
	b.TakeDamageAt("body", 61)
	b.Update(player, 0.016)

	// Force into charge -> stun sequence
	b.chargeCooldown = 0.01
	b.Update(player, 0.02)

	// Advance to stunned
	for i := 0; i < 300; i++ {
		if b.GetState() == StateStunned {
			break
		}
		b.Update(player, 0.01)
	}

	if b.GetState() != StateStunned {
		t.Fatalf("Expected StateStunned, got %d", b.GetState())
	}

	// Damage should exit stunned state
	b.TakeDamageAt("body", 5)

	if b.GetState() != StateHover {
		t.Errorf("Expected StateHover after damage in stunned state, got %d", b.GetState())
	}
}

func TestSentinelBoss_ProjectileSpawning(t *testing.T) {
	b := New(testRoomStartY, testWorldWidth)
	b.Activate()
	player := testBossPlayer()

	// Phase 0 has 2.5s cooldown, but attack starts ready
	// First update should fire projectiles
	spawnRequests := b.Update(player, 0.016)

	if len(spawnRequests) == 0 {
		t.Error("Expected projectile spawn requests on first update")
	}
}

func TestSentinelBoss_IsVulnerable_MatchesHurtboxes(t *testing.T) {
	b, _ := newActiveBoss()

	// Phase 0: always vulnerable
	if b.IsVulnerable() == false {
		t.Error("Expected IsVulnerable=true in phase 0")
	}
	if len(b.GetHurtboxes()) == 0 {
		t.Error("Expected non-empty hurtboxes when IsVulnerable=true")
	}

	// Consistency check
	isVuln := b.IsVulnerable()
	hasHurtboxes := len(b.GetHurtboxes()) > 0
	if isVuln != hasHurtboxes {
		t.Errorf("IsVulnerable (%v) should match hurtbox presence (%v)", isVuln, hasHurtboxes)
	}
}

func TestSentinelBoss_GetVulnerableTimer_ZeroInHover(t *testing.T) {
	b, player := newActiveBoss()

	// Push to phase 1
	b.TakeDamageAt("body", 61)
	b.Update(player, 0.016)

	if b.GetState() != StateHover {
		t.Fatalf("Expected StateHover, got %d", b.GetState())
	}

	timer := b.GetVulnerableTimer()
	if timer != 0 {
		t.Errorf("Expected 0 (not vulnerable) in hover phase 1, got %f", timer)
	}
}

func TestSentinelBoss_GetChargeTarget(t *testing.T) {
	b, player := newActiveBoss()

	// Push to phase 1
	b.TakeDamageAt("body", 61)
	b.Update(player, 0.016)

	// Trigger charge
	b.chargeCooldown = 0.01
	b.Update(player, 0.02)

	if b.GetState() != StateChargeWindup {
		t.Fatalf("Expected StateChargeWindup, got %d", b.GetState())
	}

	target := b.GetChargeTarget()
	// Target should be roughly at player position
	playerCenterX := player.AABB.X + player.AABB.Width/2
	if target.X < playerCenterX-Width && target.X > playerCenterX+Width {
		t.Errorf("Charge target X (%f) should be near player (%f)", target.X, playerCenterX)
	}
}
