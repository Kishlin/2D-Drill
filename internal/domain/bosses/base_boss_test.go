package bosses

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/bosses/phases"
	"github.com/Kishlin/drill-game/internal/domain/bosses/statemachine"
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// State IDs for base boss tests
const (
	testBossStateIdle   statemachine.StateID = 0
	testBossStateActive statemachine.StateID = 1
)

// baseBossWrapper embeds *BaseBoss and implements the Boss interface
// so Self can be set for virtual dispatch.
type baseBossWrapper struct {
	*BaseBoss
	damageReactions []damageReaction
}

type damageReaction struct {
	hurtboxID string
	damage    float32
}

func (w *baseBossWrapper) Update(player *entities.Player, dt float32) []projectiles.SpawnRequest {
	return w.BaseUpdate(player, dt)
}

// recordingDamageReactionHandler records damage reactions for test assertions.
type recordingDamageReactionHandler struct {
	wrapper *baseBossWrapper
}

func (h *recordingDamageReactionHandler) OnDamageReceived(hurtboxID string, damage float32) {
	h.wrapper.damageReactions = append(h.wrapper.damageReactions, damageReaction{hurtboxID, damage})
}

func testBaseBossConfig() BaseBossConfig {
	return BaseBossConfig{
		Position: types.NewVec2(100, 200),
		MaxHP:    50,
		BoxSet: NewBodyBoxSet(BodyBoxConfig{
			ID:               "body",
			Width:            80,
			Height:           60,
			OffsetX:          0,
			OffsetY:          0,
			DamagePerSec:     10,
			DamageMultiplier: 1.0,
		}),
		Phases: []phases.Config{{HPThreshold: 0}},
	}
}

func newTestBaseBoss() *baseBossWrapper {
	base := NewBaseBoss(testBaseBossConfig())
	w := &baseBossWrapper{BaseBoss: base}
	w.Self = w
	w.SetStateMachine(statemachine.NewStateMachine(
		map[statemachine.StateID]*statemachine.State{
			testBossStateIdle: {
				ID: testBossStateIdle,
				OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
					return statemachine.StateResult{NextState: statemachine.StateIDNone}
				},
			},
		},
		testBossStateIdle,
	))
	return w
}

func testBaseBossPlayer() *entities.Player {
	playerCfg := config.PlayerConfig{
		StartingMoney:    0,
		StartingItems:    [5]int{},
		StartingUpgrades: config.StartingUpgrades{},
	}
	upgradeCfg := config.UpgradeConfig{
		Engines:     []config.UpgradeTier[config.EngineStats]{{Stats: config.EngineStats{MaxSpeed: 450, Acceleration: 2500, FlyAcceleration: 2500, MaxUpwardSpeed: -600}}},
		Hulls:       []config.UpgradeTier[config.HullStats]{{Stats: config.HullStats{MaxHP: 100}}},
		FuelTanks:   []config.UpgradeTier[config.FuelTankStats]{{Stats: config.FuelTankStats{Capacity: 10}}},
		CargoHolds:  []config.UpgradeTier[config.CargoHoldStats]{{Stats: config.CargoHoldStats{Capacity: 10}}},
		HeatShields: []config.UpgradeTier[config.HeatShieldStats]{{Stats: config.HeatShieldStats{HeatResistance: 50}}},
		Drills:      []config.UpgradeTier[config.DrillStats]{{Stats: config.DrillStats{SpeedAtSurface: 1.0, SpeedAtMaxDepth: 1.0}}},
	}
	return entities.NewPlayerFromConfig(0, 0, playerCfg, upgradeCfg)
}

func TestNewBaseBoss_InitializesHP(t *testing.T) {
	b := newTestBaseBoss()

	if b.GetHP() != 50 {
		t.Errorf("Expected HP 50, got %f", b.GetHP())
	}
	if b.GetMaxHP() != 50 {
		t.Errorf("Expected MaxHP 50, got %f", b.GetMaxHP())
	}
}

func TestNewBaseBoss_StartsInactive(t *testing.T) {
	b := newTestBaseBoss()

	if b.IsActive() {
		t.Error("Expected boss to start inactive")
	}
}

func TestNewBaseBoss_StartsUndefeated(t *testing.T) {
	b := newTestBaseBoss()

	if b.IsDefeated() {
		t.Error("Expected boss to start undefeated")
	}
}

func TestBaseBoss_ActivateDeactivate(t *testing.T) {
	b := newTestBaseBoss()

	b.Activate()
	if b.IsActive() == false {
		t.Error("Expected boss to be active after Activate()")
	}

	b.Deactivate()
	if b.IsActive() {
		t.Error("Expected boss to be inactive after Deactivate()")
	}
}

func TestBaseBoss_IsDefeated_DelegatesToDamageable(t *testing.T) {
	b := newTestBaseBoss()

	b.Damageable.TakeDamage(50)

	if b.IsDefeated() == false {
		t.Error("Expected boss to be defeated after taking full damage")
	}
}

func TestBaseBoss_GetHP_AfterDamage(t *testing.T) {
	b := newTestBaseBoss()

	b.Damageable.TakeDamage(20)

	if b.GetHP() != 30 {
		t.Errorf("Expected HP 30, got %f", b.GetHP())
	}
}

func TestBaseBoss_GetPosition(t *testing.T) {
	b := newTestBaseBoss()
	pos := b.GetPosition()

	if pos.X != 100 || pos.Y != 200 {
		t.Errorf("Expected position (100, 200), got (%f, %f)", pos.X, pos.Y)
	}
}

func TestBaseBoss_GetCollisionBoxes(t *testing.T) {
	b := newTestBaseBoss()
	boxes := b.GetCollisionBoxes()

	if len(boxes) != 1 {
		t.Fatalf("Expected 1 collision box, got %d", len(boxes))
	}
	if boxes[0].ID != "body" {
		t.Errorf("Expected collision box ID 'body', got %q", boxes[0].ID)
	}
}

func TestBaseBoss_GetHitboxes(t *testing.T) {
	b := newTestBaseBoss()
	boxes := b.GetHitboxes()

	if len(boxes) != 1 {
		t.Fatalf("Expected 1 hitbox, got %d", len(boxes))
	}
	if boxes[0].DamagePerSec != 10 {
		t.Errorf("Expected DamagePerSec 10, got %f", boxes[0].DamagePerSec)
	}
}

func TestBaseBoss_GetHurtboxes(t *testing.T) {
	b := newTestBaseBoss()
	boxes := b.GetHurtboxes()

	if len(boxes) != 1 {
		t.Fatalf("Expected 1 hurtbox, got %d", len(boxes))
	}
	if boxes[0].DamageMultiplier != 1.0 {
		t.Errorf("Expected DamageMultiplier 1.0, got %f", boxes[0].DamageMultiplier)
	}
}

func TestBaseBoss_TakeDamageAt_AppliesDamageWithMultiplier(t *testing.T) {
	cfg := BaseBossConfig{
		Position: types.NewVec2(0, 0),
		MaxHP:    100,
		BoxSet: NewBoxSet(
			nil, nil,
			[]HurtboxDef{
				{BoxDef: BoxDef{ID: "weak"}, DamageMultiplier: 2.0},
			},
		),
		Phases: []phases.Config{{HPThreshold: 0}},
	}
	base := NewBaseBoss(cfg)
	w := &baseBossWrapper{BaseBoss: base}
	w.Self = w
	w.SetStateMachine(statemachine.NewStateMachine(
		map[statemachine.StateID]*statemachine.State{
			0: {ID: 0, OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			}},
		}, 0,
	))

	actual := w.TakeDamageAt("weak", 10)

	if actual != 20 {
		t.Errorf("Expected actual damage 20 (10 * 2.0), got %f", actual)
	}
	if w.GetHP() != 80 {
		t.Errorf("Expected HP 80, got %f", w.GetHP())
	}
}

func TestBaseBoss_TakeDamageAt_UnknownHurtboxReturnsZero(t *testing.T) {
	b := newTestBaseBoss()
	initialHP := b.GetHP()

	actual := b.TakeDamageAt("nonexistent", 10)

	if actual != 0 {
		t.Errorf("Expected 0 damage for unknown hurtbox, got %f", actual)
	}
	if b.GetHP() != initialHP {
		t.Errorf("HP should be unchanged, expected %f, got %f", initialHP, b.GetHP())
	}
}

func TestBaseBoss_TakeDamageAt_CallsDamageReactionHandler(t *testing.T) {
	b := newTestBaseBoss()
	handler := &recordingDamageReactionHandler{wrapper: b}
	b.DamageReactionHandler = handler

	b.TakeDamageAt("body", 15)

	if len(b.damageReactions) != 1 {
		t.Fatalf("Expected 1 damage reaction, got %d", len(b.damageReactions))
	}
	if b.damageReactions[0].hurtboxID != "body" {
		t.Errorf("Expected hurtboxID 'body', got %q", b.damageReactions[0].hurtboxID)
	}
	if b.damageReactions[0].damage != 15 {
		t.Errorf("Expected damage 15 (15 * 1.0), got %f", b.damageReactions[0].damage)
	}
}

func TestBaseBoss_BaseUpdate_ReturnsNilWhenInactive(t *testing.T) {
	b := newTestBaseBoss()
	player := testBaseBossPlayer()

	result := b.BaseUpdate(player, 0.016)

	if result != nil {
		t.Error("Expected nil result when boss is inactive")
	}
}

func TestBaseBoss_BaseUpdate_ReturnsNilWhenDefeated(t *testing.T) {
	b := newTestBaseBoss()
	b.Activate()
	b.Damageable.TakeDamage(50)
	player := testBaseBossPlayer()

	result := b.BaseUpdate(player, 0.016)

	if result != nil {
		t.Error("Expected nil result when boss is defeated")
	}
}

func TestBaseBoss_BaseUpdate_UpdatesBoxPositions(t *testing.T) {
	b := newTestBaseBoss()
	b.Activate()
	player := testBaseBossPlayer()

	b.BaseUpdate(player, 0.016)

	cb := b.GetCollisionBoxes()[0]
	if cb.X != 100 {
		t.Errorf("Expected collision box X 100 (boss pos + offset 0), got %f", cb.X)
	}
	if cb.Y != 200 {
		t.Errorf("Expected collision box Y 200 (boss pos + offset 0), got %f", cb.Y)
	}
}

func TestBaseBoss_DefaultHandlers_DoNotPanic(t *testing.T) {
	base := NewBaseBoss(testBaseBossConfig())
	w := &baseBossWrapper{BaseBoss: base}
	w.Self = w
	w.SetStateMachine(statemachine.NewStateMachine(
		map[statemachine.StateID]*statemachine.State{
			0: {ID: 0, OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			}},
		}, 0,
	))

	// Should not panic: default no-op handlers
	w.TakeDamageAt("body", 5)

	if w.GetHP() != 45 {
		t.Errorf("Expected HP 45, got %f", w.GetHP())
	}
}
