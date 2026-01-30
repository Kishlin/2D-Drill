# Boss Architecture Analysis

**Date:** January 2026
**Goal:** Evaluate the boss system for Go idiomatics, 2D game dev best practices, and ease of adding new bosses.

---

## Overall Rating: 7/10

Solid foundations, but friction points for extensibility and some Go anti-patterns.

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

---

## Issues & Concerns

### 1. Boss Factory Violates Extensibility Goal

**Location:** `internal/domain/engine/game.go:318-325`

```go
func createBossByType(bossType string, ...) (bosses.Boss, error) {
    switch bossType {
    case "test_boss":
        return test_boss.New(roomStartY, worldWidth), nil
    // Every new boss requires modifying this file
    }
}
```

**Problem:** Adding a boss requires editing `game.go`, violating the "don't touch generic files" goal.

**Solution:** Registration pattern where each boss package registers itself:

```go
// In bosses/registry.go
var registry = make(map[string]BossFactory)

func Register(name string, factory BossFactory) {
    registry[name] = factory
}

func Create(name string, params BossParams) (Boss, error) {
    factory, ok := registry[name]
    if !ok {
        return nil, fmt.Errorf("unknown boss: %s", name)
    }
    return factory(params), nil
}

// In each boss package's init()
func init() {
    bosses.Register("lava_worm", New)
}
```

---

### 2. String-Based IDs Are Fragile

**Locations:**
- `internal/domain/bosses/statemachine/types.go` - `type StateID string`
- `internal/domain/bosses/boss.go` - `TakeDamageAt(hurtboxID string, ...)`

**Problem:** String IDs compile fine with typos but fail at runtime. No compile-time safety.

**Solution:** Use typed constants:

```go
type StateID int
const (
    StatePatrol StateID = iota
    StateWindup
    StateSlam
    StateVulnerable
)
```

For hurtbox IDs, consider using an index or a typed identifier.

---

### 3. Boss Interface Is Too Large (13+ methods)

**Location:** `internal/domain/bosses/boss.go`

```go
type Boss interface {
    // Lifecycle (5 methods)
    Update, Activate, Deactivate, IsActive, IsDefeated

    // Health (3 methods)
    GetHP, GetMaxHP, GetDamageable

    // Position (1 method)
    GetPosition

    // Boxes (3 methods)
    GetCollisionBoxes, GetHitboxes, GetHurtboxes

    // Damage (1 method)
    TakeDamageAt
}
```

**Problem:** Go prefers small interfaces. Every boss must implement all 13+ methods even for simple bosses.

**Solution:** Either split into smaller interfaces or provide a `BaseBoss` struct that handles boilerplate with sensible defaults.

---

### 4. No Base Boss / Composition Helper

**Problem:** Every boss must manually wire up:
- `components.Damageable`
- `PhaseManager`
- `StateMachine`
- `Movement` behavior
- `attacks.ProjectileAttack`
- Active flag tracking
- Position management

**Solution:** Create a `BaseBoss` struct:

```go
type BaseBoss struct {
    position     types.Vec2
    damageable   components.Damageable
    active       bool
    phaseManager *PhaseManager
    stateMachine *statemachine.StateMachine
    movement     movement.MovementBehavior

    // Cached boxes to avoid allocation
    collisionBoxes []CollisionBox
    hitboxes       []Hitbox
    hurtboxes      []Hurtbox
}

// Default implementations for common methods
func (b *BaseBoss) GetHP() float32 { return b.damageable.HP }
func (b *BaseBoss) IsActive() bool { return b.active }
// ... etc
```

Concrete bosses embed `BaseBoss` and override only what's unique.

---

### 5. StateBehaviors Callback Pattern Is Unusual for Go

**Location:** `internal/domain/bosses/test_boss/states.go`

```go
type StateBehaviors struct {
    GetPhase           func() int
    GetAOECooldown     func() float32
    SetAOECooldown     func(float32)
    FireProjectiles    func(player) []SpawnRequest
    // ... many callbacks
}
```

**Problem:** This pattern is more JavaScript/functional than idiomatic Go. It adds indirection and makes code harder to follow.

**Idiomatic Go alternatives:**
1. Pass boss pointer directly to state handlers
2. Use method receivers on state types that hold boss reference
3. Define states as methods on the boss itself

```go
// Option 1: Pass boss to state
func PatrolUpdate(boss *TestBoss, ctx *StateContext) StateResult { ... }

// Option 2: State holds boss reference
type PatrolState struct {
    boss *TestBoss
}
func (s *PatrolState) Update(ctx *StateContext) StateResult { ... }
```

---

### 6. Vulnerability Logic Is Scattered

**Problem:** To determine "is the boss vulnerable?" you must check:
1. `PhaseConfig.AlwaysVulnerable` (phase system)
2. Current state == `StateVulnerable` (state machine)
3. `GetHurtboxes()` returns non-empty (box system)

The logic lives in `GetHurtboxes()`, mixing "what are my hurtboxes" with "am I vulnerable".

**Solution:** Single source of truth:

```go
func (b *TestBoss) IsVulnerable() bool {
    phase := b.phaseManager.CurrentPhase()
    if phase.AlwaysVulnerable {
        return true
    }
    return b.stateMachine.CurrentState() == StateVulnerable
}

func (b *TestBoss) GetHurtboxes() []Hurtbox {
    if !b.IsVulnerable() {
        return nil
    }
    return b.hurtboxes // Pre-allocated
}
```

---

### 7. TakeDamageAt Has Hidden Side Effects

**Location:** `internal/domain/bosses/test_boss/boss.go`

```go
func (b *TestBoss) TakeDamageAt(hurtboxID string, baseDamage float32) float32 {
    // Applies damage AND...
    if b.stateMachine.CurrentState() == StateVulnerable {
        b.stateMachine.Transition(StatePatrol) // ...triggers state transition
    }
}
```

**Problem:** Callers don't expect `TakeDamage` to change boss behavior state. This couples damage and state machine.

**Solution:** Return information about what happened, let caller decide:

```go
type DamageResult struct {
    DamageDealt     float32
    InterruptedState bool
}

func (b *TestBoss) TakeDamageAt(...) DamageResult {
    // Only apply damage, don't transition
}

// Caller or boss update loop handles state transitions
```

Or handle the interrupt in the state's `OnUpdate` by checking if damage was received.

---

### 8. Hardcoded Duration Values in States

**Location:** `internal/domain/bosses/test_boss/states.go`

```go
if ctx.Elapsed >= 1.0 {  // Windup duration - hardcoded
if ctx.Elapsed >= 0.3 {  // Slam duration - hardcoded
if ctx.Elapsed >= 0.4 {  // Between-slam pause - hardcoded
```

**Problem:** Timing values buried in code, not configurable.

**Solution:** Pass timing config through state context or behaviors:

```go
type BossTimingConfig struct {
    WindupDuration      float32
    SlamDuration        float32
    BetweenSlamsPause   float32
}

// Access via behaviors or context
if ctx.Elapsed >= behaviors.GetWindupDuration() {
```

---

### 9. Potential GC Pressure from Slice Returns

**Location:** All `Get*Boxes()` methods

```go
func (b *TestBoss) GetHitboxes() []Hitbox {
    return []Hitbox{{
        X: b.position.X,
        // ... new allocation every frame
    }}
}
```

**Problem:** If called every frame, creates garbage for GC.

**Solution:** Pre-allocate and update in place:

```go
type TestBoss struct {
    // Pre-allocated slices
    hitboxCache []Hitbox
}

func (b *TestBoss) updateBoxes() {
    b.hitboxCache[0].X = b.position.X
    // Update positions without reallocating
}

func (b *TestBoss) GetHitboxes() []Hitbox {
    return b.hitboxCache
}
```

---

### 10. Nested State Machines (AOEAttack)

**Location:** `internal/domain/bosses/attacks/aoe_attack.go`

```go
type aoeState int
const (
    aoeStateIdle aoeState = iota
    aoeStateTelegraph
    aoeStateDamage
    aoeStateVulnerable
)
```

**Problem:** Now there are two state machines (boss + AOE) that must stay synchronized. This adds cognitive overhead and potential for bugs.

**Alternative:** Handle AOE phases as boss states directly:
- `StateAOETelegraph`
- `StateAOEDamage`
- `StateVulnerable`

One state machine, one source of truth.

---

### 11. Three-Box System May Be Overkill

**Problem:** Most 2D platformer bosses use:
- Collision box = hitbox (same shape)
- One or more hurtboxes

Having `CollisionBox` separate from `Hitbox` is unusual if they're always identical.

**When it's useful:** Bosses with projectile-only damage (collision but no contact damage) or damage auras larger than physical body.

**Recommendation:** Keep the system but document when to use each. Consider a helper for the common case:

```go
// For bosses where collision = hitbox
func (b *BaseBoss) SetBodyBox(box types.AABB, contactDamage float32) {
    b.collisionBoxes = []CollisionBox{{...}}
    b.hitboxes = []Hitbox{{..., DamagePerSec: contactDamage}}
}
```

---

### 12. No Animation/Sound Hooks

**Problem:** States don't communicate which animation or sound should play. Real boss fights need audiovisual feedback.

**Solution:** Extend `StateResult`:

```go
type StateResult struct {
    NextState     StateID
    SpawnRequests []projectiles.SpawnRequest
    Animation     AnimationID  // What animation to play
    Sound         SoundID      // What sound to trigger
}
```

Or have states set animation on the boss directly (if boss holds animation state).

---

## Recommendations Summary

### To Achieve "Add Boss = Define Subpackage Only"

1. **Registration pattern** - Boss packages register themselves in `init()`
2. **BaseBoss struct** - Handles common fields and default method implementations
3. **Typed state IDs** - Compile-time safety for state names
4. **Config struct per boss** - All timing/damage values in one place
5. **Remove StateBehaviors callbacks** - Pass boss pointer to states directly
6. **Single vulnerability source** - `IsVulnerable()` method as single source of truth
7. **Pre-allocated boxes** - Avoid per-frame allocations

### Ideal New Boss Structure

```
bosses/
  lava_worm/
    boss.go      # Embeds BaseBoss, defines config, registers via init()
    states.go    # State definitions only (the unique part)
    config.go    # All timing/damage values
```

Minimal boilerplate - define only what's unique about this boss.

---

## Performance Concerns

| Concern | Severity | Notes |
|---------|----------|-------|
| Slice allocations in `Get*Boxes()` | Minor | Fix with pre-allocation |
| String comparisons for IDs | Minor | Fix with typed constants |
| Projectile pool | Good | Already optimized |
| State machine overhead | Negligible | Simple switch/map lookups |

---

## Next Steps

1. Read `IMPLEMENTATION_REFERENCE.md` for detailed current implementation docs
2. Decide which issues to address (see priority in each section)
3. Consider creating `BaseBoss` first - biggest impact on extensibility
4. Add registration pattern second - enables "add package only" workflow
