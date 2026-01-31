# Boss Architecture Analysis

**Date:** January 2026
**Last Updated:** January 2026
**Goal:** Evaluate the boss system for Go idiomatics, 2D game dev best practices, and ease of adding new bosses.

---

## Overall Rating: 8.5/10

Strong foundations with most friction points resolved. Adding a new boss now requires only creating a subpackage with no modifications to core files.

---

## What's Done Well

### 1. Return-Based Architecture
`Update()` returning spawn requests instead of mutating external state is clean and testable. No hidden side effects in the update loop.

### 2. Separation of Concerns
Movement, attacks, phases, and state machine are properly decoupled into separate packages/types.

### 3. Three-Box System
Flexible hitbox/hurtbox/collision box separation allows complex boss scenarios (invulnerability phases, damage-only zones, etc.).

### 4. Data-Driven Phases
`PhaseConfig` struct allows phase behavior to be defined declaratively rather than in code.

### 5. Pool-Based Projectile System
Pre-allocated projectile pool avoids GC pressure during combat.

### 6. BoxSet System (NEW)
Pre-allocated boxes with position synchronization eliminate per-frame allocations.

### 7. Registration Pattern (NEW)
Boss packages self-register via `init()`, enabling "add package only" workflow.

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

// In test_boss/boss.go init()
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

// In test_boss/states.go
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

### 3. Boss Interface Is Too Large (13+ methods) ⚠️ ACCEPTABLE

**Location:** `internal/domain/bosses/boss.go`

**Status:** Interface size unchanged, but `BoxSet` reduces implementation burden. Each method has clear purpose and the pattern is well-documented.

**Alternative considered:** BaseBoss struct could provide default implementations. Not yet implemented but remains an option for future optimization.

---

### 4. No Base Boss / Composition Helper ⚠️ NOT IMPLEMENTED

**Status:** Not implemented, but `BoxSet` significantly reduces boilerplate.

**Current pattern:** Bosses manually wire up components but the code is straightforward:

```go
b := &TestBoss{
    position:     types.NewVec2(centerX, floorY),
    damageable:   components.NewDamageable(MaxHP, MaxHP),
    boxSet:       bosses.NewBodyBoxSet(bosses.BodyBoxConfig{...}),
    // ...
}
```

**Recommendation:** Consider `BaseBoss` if adding many bosses reveals repetitive patterns.

---

### 5. StateBehaviors Callback Pattern Is Unusual for Go ⚠️ RETAINED (ACCEPTABLE)

**Location:** `internal/domain/bosses/test_boss/states.go`

**Status:** Pattern retained as an acceptable trade-off.

**Trade-offs documented:**
- Avoids circular dependencies between boss and states
- Provides clean encapsulation
- States don't need to know concrete boss type
- Well-documented in `docs/BOSS.md`

Alternative (passing boss pointer) would require either:
- Circular imports, or
- Defining states in same package as boss

Current approach is pragmatic.

---

### 6. Vulnerability Logic Is Scattered ✅ RESOLVED

**Location:** `internal/domain/bosses/test_boss/boss.go:308-336`

**Solution Implemented:** `GetHurtboxes()` is the single source of truth:

```go
func (b *TestBoss) GetHurtboxes() []bosses.Hurtbox {
    if b.phaseManager.IsAlwaysVulnerable() || b.stateMachine.CurrentState() == StateVulnerable {
        return b.boxSet.Hurtboxes
    }
    return nil // Invulnerable
}

func (b *TestBoss) IsVulnerable() bool {
    return len(b.GetHurtboxes()) > 0
}
```

Empty hurtbox slice = invulnerable. Simple, clear, single source of truth.

---

### 7. TakeDamageAt Has Hidden Side Effects ⚠️ DOCUMENTED

**Location:** `internal/domain/bosses/test_boss/boss.go:316-331`

**Status:** Side effect remains (transitioning out of vulnerable state on damage) but is now well-documented and intentional game design.

```go
func (b *TestBoss) TakeDamageAt(hurtboxID string, baseDamage float32) float32 {
    // ...
    if b.stateMachine.CurrentState() == StateVulnerable {
        b.stateMachine.TransitionTo(StatePatrol, &statemachine.StateContext{})
    }
    return actual
}
```

This is desired behavior: hitting the boss during vulnerability window immediately ends that window.

---

### 8. Hardcoded Duration Values in States ✅ RESOLVED

**Location:** `internal/domain/bosses/test_boss/boss.go:22-33`

**Solution Implemented:** All timing values are package-level constants:

```go
const (
    MaxHP           = 100.0
    WindupDuration  = 1.0
    SlamDuration    = 0.3
    DoubleSlamPause = 0.4
    // ...
)
```

States reference these constants:
```go
if ctx.Elapsed >= WindupDuration { ... }
if ctx.Elapsed >= SlamDuration { ... }
```

Phase-dependent values accessed via `behaviors.GetVulnerableDuration()`.

---

### 9. Potential GC Pressure from Slice Returns ✅ RESOLVED

**Location:** `internal/domain/bosses/boxes.go:85-151`

**Solution Implemented:** `BoxSet` pre-allocates all boxes and updates positions in-place:

```go
type BoxSet struct {
    // Definitions (static)
    collisionDefs []BoxDef
    hitboxDefs    []HitboxDef
    hurtboxDefs   []HurtboxDef

    // Runtime boxes (pre-allocated, positions updated)
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

**Status:** AOE attack no longer uses a separate state machine. AOE phases are handled as boss states directly:
- `StateWindup` - Telegraph before slam
- `StateWindupBetween` - Pause between slams
- `StateSlam` - Active damage
- `StateVulnerable` - Vulnerability window

One state machine, one source of truth.

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
| Large Boss interface | ⚠️ Acceptable | Well-documented, BoxSet helps |
| No BaseBoss struct | ⚠️ Not needed yet | BoxSet reduces boilerplate |
| StateBehaviors callbacks | ⚠️ Retained | Acceptable trade-off, documented |
| Scattered vulnerability logic | ✅ Resolved | `GetHurtboxes()` as single source |
| TakeDamageAt side effects | ⚠️ Documented | Intentional game design |
| Hardcoded durations | ✅ Resolved | Package constants + PhaseConfig |
| GC pressure from slices | ✅ Resolved | BoxSet pre-allocation |
| Nested state machines | ✅ Resolved | Single state machine |
| Three-box complexity | ⚠️ Documented | Helper for common case |
| No animation/sound hooks | ⚠️ Deferred | Add when needed |

---

## Remaining Opportunities

1. **BaseBoss struct** - Would reduce per-boss boilerplate by ~10-15%. Consider if adding many bosses.

2. **Animation/Sound hooks** - Extend `StateResult` when adding polish:
   ```go
   type StateResult struct {
       NextState     StateID
       SpawnRequests []projectiles.SpawnRequest
       Animation     AnimationID
       Sound         SoundID
   }
   ```

---

## Architecture Quality Progression

| Version | Rating | Key Changes |
|---------|--------|-------------|
| Initial | 7/10 | Good foundations, friction points |
| Current | 8.5/10 | Registry, BoxSet, typed IDs, config constants |

The system now achieves the goal: **adding a new boss requires only creating a subpackage**.
