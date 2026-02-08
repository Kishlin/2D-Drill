# Boss Architecture Analysis

**Date:** January 2026
**Last Updated:** February 2026
**Goal:** Evaluate the boss system for Go idiomatics, 2D game dev best practices, and ease of adding new bosses.

---

## Overall Rating: 9/10

Excellent architecture with BaseBoss significantly reducing boilerplate. Adding a new boss requires creating a subpackage, embedding BaseBoss, and implementing handlers.

---

## What's Done Well

### 1. Return-Based Architecture
`Update()` returning spawn requests instead of mutating external state is clean and testable. No hidden side effects in the update loop.

### 2. Separation of Concerns
Movement, attacks, phases, and state machine are properly decoupled into separate packages/types.

### 3. Three-Box System
Flexible hitbox/hurtbox/collision box separation allows complex boss scenarios (invulnerability phases, damage-only zones, etc.).

### 4. Data-Driven Phases
`phases.Config` struct allows phase behavior to be defined declaratively rather than in code.

### 5. Pool-Based Projectile System
Pre-allocated projectile pool avoids GC pressure during combat.

### 6. BoxSet System
Pre-allocated boxes with position synchronization eliminate per-frame allocations.

### 7. Registration Pattern
Boss packages self-register via `init()`, enabling "add package only" workflow.

### 8. BaseBoss Struct (NEW)
Embeddable struct provides default implementations for 11+ interface methods, reducing ~80 lines of boilerplate per boss.

### 9. Package Organization (NEW)
- `bosses/phases/` - Phase management as separate package
- `boss_catalog/` - Boss implementations separate from infrastructure
- Clean import paths and separation of concerns

---

## Issues & Status

### 1. Boss Factory Violates Extensibility Goal ✅ RESOLVED

**Location:** `internal/domain/bosses/registry.go`

**Solution Implemented:** Registration pattern where each boss package registers itself:

```go
// In bosses/registry.go
var registry = make(map[string]BossConstructor)

func Register(bossType string, constructor BossConstructor) { ... }
func Create(bossType string, roomStartY, worldWidth float32) (Boss, error) { ... }

// In boss_catalog/test_boss/boss.go init()
func init() {
    bosses.Register("test_boss", func(roomStartY, worldWidth float32) bosses.Boss {
        return New(roomStartY, worldWidth)
    })
}
```

`game.go` now uses `bosses.Create()` instead of a switch statement.

---

### 2. String-Based IDs Are Fragile ✅ RESOLVED

**Location:** `internal/domain/bosses/statemachine/types.go`

**Solution Implemented:** Typed integer constants with iota:

```go
type StateID int
const StateIDNone StateID = -1

// In boss_catalog/test_boss/states.go
const (
    StatePatrol statemachine.StateID = iota
    StateWindup
    StateWindupBetween
    StateSlam
    StateVulnerable
)
```

Typos in state names are now caught at compile time.

---

### 3. Boss Interface Is Too Large (13+ methods) ✅ RESOLVED

**Location:** `internal/domain/bosses/boss.go`

**Solution Implemented:** `BaseBoss` struct provides default implementations for most interface methods. Concrete bosses only need to:
- Embed `*bosses.BaseBoss`
- Implement handlers (`PhaseChangeHandler`, `DamageReactionHandler`)
- Override `GetHurtboxes()` for custom vulnerability logic
- Implement `Update()` by calling `b.BaseUpdate(player, dt)`

---

### 4. No Base Boss / Composition Helper ✅ IMPLEMENTED

**Location:** `internal/domain/bosses/base_boss.go`

**Solution Implemented:** `BaseBoss` struct with:

```go
type BaseBoss struct {
    Position      types.Vec2
    Damageable    components.Damageable
    Active        bool
    BoxSet        *BoxSet
    StateMachine  *statemachine.StateMachine
    PhaseManager  *phases.Manager
    CurrentPlayer *entities.Player

    PhaseChangeHandler    PhaseChangeHandler
    DamageReactionHandler DamageReactionHandler
}

// Default implementations provided:
func (b *BaseBoss) Activate()
func (b *BaseBoss) Deactivate()
func (b *BaseBoss) IsActive() bool
func (b *BaseBoss) IsDefeated() bool
func (b *BaseBoss) GetHP() float32
func (b *BaseBoss) GetMaxHP() float32
func (b *BaseBoss) GetDamageable() *components.Damageable
func (b *BaseBoss) GetPosition() types.Vec2
func (b *BaseBoss) GetCollisionBoxes() []CollisionBox
func (b *BaseBoss) GetHitboxes() []Hitbox
func (b *BaseBoss) GetHurtboxes() []Hurtbox
func (b *BaseBoss) TakeDamageAt(hurtboxID string, baseDamage float32) float32
func (b *BaseBoss) BaseUpdate(player *entities.Player, dt float32) []projectiles.SpawnRequest
```

**Usage in TestBoss:**
```go
type TestBoss struct {
    *bosses.BaseBoss
    // Boss-specific fields only
    movement         *movement.Grounded
    projectileAttack *attacks.ProjectileAttack
    // ...
}

func (b *TestBoss) Update(player *entities.Player, dt float32) []projectiles.SpawnRequest {
    return b.BaseUpdate(player, dt)
}
```

---

### 5. StateBehaviors Callback Pattern Is Unusual for Go ✅ RESOLVED

**Location:** `internal/domain/boss_catalog/test_boss/boss.go`

**Solution Implemented:** Removed the callback pattern entirely. States are now defined in a `buildStates()` method on the boss struct with direct field access:

```go
func (b *TestBoss) buildStates() map[statemachine.StateID]*statemachine.State {
    return map[statemachine.StateID]*statemachine.State{
        StatePatrol: {
            ID:      StatePatrol,
            CanMove: true,
            OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
                b.Position = b.movement.Update(b.Position, ctx.Dt)
                // Direct field access - no callbacks needed
                if b.hasAOEAttack() {
                    b.aoeCooldown -= ctx.Dt
                    // ...
                }
                return statemachine.StateResult{NextState: statemachine.StateIDNone}
            },
        },
        // ... other states
    }
}
```

**Benefits:**
- No callback indirection - states directly access boss fields
- Simpler mental model - one file contains boss logic
- Better IDE support - direct field/method access enables navigation
- More idiomatic Go - methods instead of function fields
- `states.go` now only contains state ID constants (~12 lines)

---

### 6. Vulnerability Logic Is Scattered ✅ RESOLVED

**Location:** `internal/domain/boss_catalog/test_boss/boss.go`

**Solution Implemented:** `GetHurtboxes()` is the single source of truth. Each boss decides its own vulnerability rules:

```go
func (b *TestBoss) GetHurtboxes() []bosses.Hurtbox {
    // Phase 1: always has hurtboxes
    // Phase 2+: only has hurtboxes during StateVulnerable
    if b.PhaseManager.GetCurrentPhase() == 0 || b.StateMachine.CurrentState() == StateVulnerable {
        return b.BoxSet.Hurtboxes
    }
    return []bosses.Hurtbox{}
}

func (b *TestBoss) IsVulnerable() bool {
    return len(b.GetHurtboxes()) > 0
}
```

**Key improvement:** Vulnerability is boss-specific logic, not part of `phases.Config`. Phases only define HP thresholds and phase-specific parameters (speed, cooldowns).

---

### 7. TakeDamageAt Has Hidden Side Effects ⚠️ DOCUMENTED

**Location:** `internal/domain/boss_catalog/test_boss/boss.go`

**Status:** Side effect remains (transitioning out of vulnerable state on damage) but is now handled via `DamageReactionHandler`:

```go
func (b *TestBoss) OnDamageReceived(hurtboxID string, damage float32) {
    if b.StateMachine.CurrentState() == StateVulnerable {
        b.StateMachine.TransitionTo(StatePatrol, &statemachine.StateContext{})
    }
}
```

This is desired behavior: hitting the boss during vulnerability window immediately ends that window.

---

### 8. Hardcoded Duration Values in States ✅ RESOLVED

**Location:** `internal/domain/boss_catalog/test_boss/boss.go`

**Solution Implemented:** All timing values are package-level constants:

```go
const (
    MaxHP                    = 100.0
    WindupDuration           = 1.0
    SlamDuration             = 0.3
    DoubleSlamPause          = 0.4
    Phase2VulnerableDuration = 3.0
    Phase3VulnerableDuration = 2.0
)
```

States reference these constants:
```go
if ctx.Elapsed >= WindupDuration { ... }
if ctx.Elapsed >= SlamDuration { ... }
```

---

### 9. Potential GC Pressure from Slice Returns ✅ RESOLVED

**Location:** `internal/domain/bosses/boxes.go`

**Solution Implemented:** `BoxSet` pre-allocates all boxes and updates positions in-place:

```go
type BoxSet struct {
    CollisionBoxes []CollisionBox
    Hitboxes       []Hitbox
    Hurtboxes      []Hurtbox
}

func (bs *BoxSet) UpdatePositions(bossX, bossY float32) {
    for i, def := range bs.collisionDefs {
        bs.CollisionBoxes[i].X = bossX + def.OffsetX
        bs.CollisionBoxes[i].Y = bossY + def.OffsetY
    }
    // ... hitboxes and hurtboxes similarly
}
```

Zero allocations per frame.

---

### 10. Nested State Machines (AOEAttack) ✅ RESOLVED

**Status:** The TestBoss no longer uses a separate AOE state machine. AOE phases are handled as boss states directly:
- `StateWindup` - Telegraph before slam
- `StateWindupBetween` - Pause between slams
- `StateSlam` - Active damage
- `StateVulnerable` - Vulnerability window

One state machine, one source of truth.

**Note:** A standalone `AOEAttack` component still exists in `attacks/aoe_attack.go` as a reusable building block with its own internal state machine (Idle → Telegraph → Damage → Vulnerable). Future bosses can choose either approach: manage AOE through boss states (like TestBoss) or use the self-contained `AOEAttack` component.

---

### 11. Three-Box System May Be Overkill ⚠️ DOCUMENTED

**Status:** System retained with helper for common case:

```go
// NewBodyBoxSet for bosses where collision = hitbox = hurtbox
func NewBodyBoxSet(cfg BodyBoxConfig) *BoxSet {
    // Creates all three box types from single config
}
```

Full flexibility available when needed, simple helper for common case.

---

### 12. No Animation/Sound Hooks ⚠️ NOT YET NEEDED

**Status:** Not implemented. Will address when adding visual/audio polish.

Current rendering uses `GetState()` and `GetStateTimer()` for basic visual feedback.

---

## Summary Table

| Issue | Status | Notes |
|-------|--------|-------|
| Boss factory switch statement | ✅ Resolved | Registration pattern via `init()` |
| String-based state IDs | ✅ Resolved | `StateID int` with iota constants |
| Large Boss interface | ✅ Resolved | BaseBoss provides defaults |
| No BaseBoss struct | ✅ Implemented | ~80 lines boilerplate reduction |
| StateBehaviors callbacks | ✅ Resolved | Removed, direct field access |
| Scattered vulnerability logic | ✅ Resolved | `GetHurtboxes()` as single source |
| TakeDamageAt side effects | ⚠️ Documented | Via DamageReactionHandler |
| Hardcoded durations | ✅ Resolved | Package constants |
| GC pressure from slices | ✅ Resolved | BoxSet pre-allocation |
| Nested state machines | ✅ Resolved | Single state machine |
| Three-box complexity | ⚠️ Documented | Helper for common case |
| No animation/sound hooks | ⚠️ Deferred | Add when needed |

---

## Architecture Quality Progression

| Version | Rating | Key Changes |
|---------|--------|-------------|
| Initial | 7/10 | Good foundations, friction points |
| Previous | 8.5/10 | Registry, BoxSet, typed IDs, config constants |
| Current | 9/10 | BaseBoss, phases package, boss_catalog separation |

The system now achieves the goal: **adding a new boss requires only creating a subpackage and embedding BaseBoss**.
