# Boss Architecture Analysis

**Date:** January 2026
**Last Updated:** February 2026
**Goal:** Evaluate the boss system for Go idiomatics, 2D game dev best practices, and ease of adding new bosses.

---

## Overall Rating: 8/10

Strong architecture with well-separated infrastructure and a clear extensibility story. BaseBoss and the double-registry pattern (domain + renderer) make adding new bosses straightforward. However, the single concrete boss makes some design decisions hard to validate, and several areas show signs of premature generalization or TestBoss-specific assumptions baked into shared infrastructure.

---

## Strengths

### 1. Return-Based Update Loop
`Update()` returning `[]projectiles.SpawnRequest` instead of mutating external state is clean and testable. The state machine propagates spawn requests upward through `StateResult`, keeping the entire update path free of hidden side effects.

### 2. Double-Registry Pattern
Two parallel registries, one in domain (`bosses.Register`) and one in adapters (`bossrenderers.Register`), both triggered by `init()`. Adding a new boss requires creating a catalog package and a renderer package — no modifications to `game.go` or core rendering. The renderer registry uses a `CanRender(boss) bool` pattern that cleanly decouples the adapter layer from concrete boss types.

### 3. Three-Box System with BoxSet Pre-Allocation
The collision/hitbox/hurtbox separation covers real 2D game scenarios (invulnerability phases, damage-only zones). `BoxSet` pre-allocates all boxes and updates positions in-place with `UpdatePositions()`, achieving zero per-frame allocations. `NewBodyBoxSet` provides a single-call shortcut for the common case where all three boxes share the same bounds.

### 4. BaseBoss Composition
Embedding `*BaseBoss` provides default implementations for 11+ interface methods, reducing ~80 lines of boilerplate per boss. The `PhaseChangeHandler` and `DamageReactionHandler` interfaces with no-op defaults let concrete bosses opt into only the hooks they need.

### 5. State Machine Design
Clean separation between infrastructure (`statemachine` package) and boss-specific states (`buildStates()` with closures). Closures capture the boss struct, giving states direct field access without callback indirection. Typed `StateID` constants with `iota` catch typos at compile time.

### 6. Data-Driven Phase Progression
`phases.Manager` tracks HP thresholds and advances phases automatically. `BaseUpdate` calls the phase manager on every frame and dispatches to the `PhaseChangeHandler` when a transition occurs. The boss never needs to manually check "am I in a new phase."

### 7. Projectile Movement Polymorphism
The `projectiles.Movement` interface with `Linear`, `Sinusoidal`, `Homing`, and `Orbital` implementations gives bosses a rich projectile vocabulary without any boss-specific code in the projectile system.

### 8. Package Organization
Clear separation: `bosses/` (infrastructure), `boss_catalog/` (implementations), `bosses/phases/`, `bosses/statemachine/`, `bosses/attacks/`, `bosses/movement/` (composable building blocks). Import paths are clean and each package has a focused responsibility.

---

## Current Issues

### ~~1. `phases.Config` Is TestBoss-Specific~~ (Resolved)

**Resolution:** `phases.Config` now contains only `HPThreshold`. Boss-specific phase parameters (movement speed, cooldowns, etc.) are stored in each boss's own `phaseConfig` struct. The `PhaseChangeHandler.OnPhaseChange` signature was simplified to `OnPhaseChange(phaseIndex int)` — each boss looks up its own config by index.

---

### ~~2. `BaseBoss.TakeDamageAt` Bypasses Vulnerability Override~~ (Resolved)

**Resolution:** `BaseBoss` now has a `Self Boss` field that enables virtual dispatch. When set, `TakeDamageAt` calls `Self.GetHurtboxes()` instead of accessing `BoxSet.Hurtboxes` directly. Concrete bosses set `b.Self = b` during construction (alongside the existing handler assignments). TestBoss's duplicated `TakeDamageAt` override was removed — the base implementation now correctly dispatches to `TestBoss.GetHurtboxes()` for phase/state-dependent vulnerability.

---

### ~~3. Hardcoded Projectile Parameters in `OnPhaseChange`~~ (Resolved)

**Resolution:** Extracted `ProjectileCount`, `ProjectileSpeed`, `ProjectileSize`, and `ProjectileDamage` as named constants in the TestBoss package. Both `New()` and `OnPhaseChange()` now reference these constants instead of magic numbers.

---

### ~~4. `AOEAttack` Component Is Unused~~ (Resolved)

**Resolution:** Removed `aoe_attack.go` (161 lines) per YAGNI. No boss used it — TestBoss manages AOE through its own state machine states. Can be recreated when a boss actually needs a reusable AOE component, designed around the real use case.

---

### ~~5. Lava Floor Damage Is Hardcoded~~ (Resolved)

**Resolution:** Added `FloorDamage float32` field to `config.BossRoomConfig`. `handleFloorDamage` now reads `s.bossRoomCfg.FloorDamage` instead of using a hardcoded `10.0`. All level configs set `FloorDamage: 10.0`.

---

### 6. `GetAOEInfo` Allocates Per Frame

**Location:** `internal/domain/boss_catalog/test_boss/boss.go:367-399`

**Issue:** `GetAOEInfo()` returns `&bosses.AOEInfo{...}`, heap-allocating a new struct every frame it's called. This is inconsistent with the zero-allocation philosophy behind `BoxSet`.

**Suggestion:** Store `AOEInfo` as a field on TestBoss and return a pointer to it (same pattern as `BoxSet`), or return by value.

**Severity:** Low — one small allocation per frame is negligible, but inconsistent.

---

### 7. Boss Update Called When Inactive

**Location:** `internal/domain/systems/boss_fight.go:69`

**Issue:** `BossFightSystem.Update` calls `s.boss.Update(player, dt)` unconditionally, even when the boss is deactivated. `BaseUpdate` early-returns when inactive, so there's no bug, but the system is delegating a responsibility it should own.

```go
spawnRequests := s.boss.Update(player, dt) // Called even after Deactivate()
```

**Severity:** Low — no behavioral impact, but the system should guard this.

---

### 8. No Boss Reset Mechanism

**Issue:** There's no way to reset a boss to its initial state. If a player dies and the game needs to restart the level, the boss must be recreated entirely. This works for now but won't scale if boss creation becomes expensive (loading assets, complex initialization).

**Severity:** Low — not a problem until level restart is needed in gameplay.

---

### 9. Test Coverage Is Low (~25-30%)

**Issue:** The tested components (phase manager, state machine, projectile attack, grounded movement) have reasonable happy-path coverage. But the core infrastructure and the only concrete boss have zero tests:

| Component | Coverage |
|-----------|----------|
| `phases/` | ~70% (happy path) |
| `statemachine/` | ~75% (happy path) |
| `attacks/projectile_attack` | ~60% |
| `movement/grounded` | ~70% |
| `base_boss` | 0% |
| `boxes` | 0% |
| `registry` | 0% |
| `boss_fight` (system) | 0% |
| `test_boss` | 0% |
| ~~`aoe_attack`~~ | ~~0%~~ (Removed) |

The untested code includes `BaseBoss` (foundation for all bosses), `BossFightSystem` (game state transitions), and `TestBoss` (the only integration of all components).

**Severity:** Medium — the code works in-game, but regression risk is high when adding a second boss.

---

### 10. `State.CanMove` Is Purely Informational

**Location:** `internal/domain/bosses/statemachine/types.go:31`

**Issue:** The `CanMove` field on `State` is exposed via `StateMachine.CanMove()` but the state machine doesn't enforce anything with it. Movement is handled manually in the StatePatrol `OnUpdate` callback. No code outside TestBoss's `buildStates()` reads `CanMove`.

This is documentation masquerading as data. If movement were handled by BaseBoss (move when `CanMove` is true, skip when false), it would be useful infrastructure. Currently it's just a label.

**Severity:** Low — harmless, but could mislead future boss authors into thinking movement is handled automatically.

---

### 11. `MovementBehavior` Interface Is Not Used Polymorphically

**Location:** `internal/domain/bosses/movement/movement.go`

**Issue:** The `MovementBehavior` interface exists but TestBoss uses `*movement.Grounded` directly, not the interface. No code accepts `MovementBehavior` as a parameter or field type. The interface is speculative generalization.

**Severity:** Low — doesn't hurt, but premature. Will be validated when a boss needs a different movement type.

---

## Historical Issues (Previously Resolved)

These issues were identified in earlier reviews and have been addressed:

| Issue | Resolution |
|-------|-----------|
| Boss factory switch statement | Registration pattern via `init()` |
| String-based state IDs | `StateID int` with iota constants |
| Large Boss interface (13+ methods) | BaseBoss provides defaults |
| No BaseBoss struct | Implemented, ~80 lines boilerplate reduction |
| `phases.Config` is TestBoss-specific | `Config` reduced to `HPThreshold` only; boss-specific params in concrete boss |
| `TakeDamageAt` bypasses overrides | `Self Boss` field on BaseBoss enables virtual dispatch; concrete bosses set `b.Self = b` |
| Hardcoded projectile params in `OnPhaseChange` | Extracted as named constants (`ProjectileCount`, `ProjectileSpeed`, `ProjectileSize`, `ProjectileDamage`) |
| StateBehaviors callback pattern | Removed — closures with direct field access |
| Scattered vulnerability logic | `GetHurtboxes()` as single source of truth |
| Hardcoded duration values | Package-level constants |
| GC pressure from slice returns | BoxSet pre-allocation |
| Nested state machines (AOE) | Single state machine per boss |
| `TakeDamageAt` side effects | Channeled through `DamageReactionHandler` |

---

## Summary Table (Current Issues)

| # | Issue | Severity | Category |
|---|-------|----------|----------|
| ~~1~~ | ~~`phases.Config` is TestBoss-specific~~ | ~~Medium~~ | ~~Extensibility~~ (Resolved) |
| ~~2~~ | ~~`BaseBoss.TakeDamageAt` bypasses overrides~~ | ~~Medium~~ | ~~Go composition~~ (Resolved) |
| ~~3~~ | ~~Hardcoded projectile params in `OnPhaseChange`~~ | ~~Low~~ | ~~Data-driven~~ (Resolved) |
| ~~4~~ | ~~`AOEAttack` component unused and untested~~ | ~~Low~~ | ~~Dead code~~ (Resolved) |
| ~~5~~ | ~~Lava floor damage hardcoded~~ | ~~Low~~ | ~~Data-driven~~ (Resolved) |
| 6 | `GetAOEInfo` allocates per frame | Low | Allocation |
| 7 | Boss update called when inactive | Low | Layering |
| 8 | No boss reset mechanism | Low | Lifecycle |
| 9 | Test coverage ~25-30% | Medium | Testing |
| 10 | `State.CanMove` purely informational | Low | API clarity |
| 11 | `MovementBehavior` not used polymorphically | Low | Premature abstraction |

---

## Architecture Quality Progression

| Version | Rating | Key Changes |
|---------|--------|-------------|
| Initial | 7/10 | Good foundations, friction points |
| +Registry/BoxSet/TypedIDs | 8.5/10 | Extensibility and performance |
| +BaseBoss/Phases/Catalog | 9/10 | Boilerplate reduction, package organization |
| Current reassessment | 8/10 | Strong with one boss; some abstractions unvalidated by a second |

The previous 9/10 was fair given the trajectory of improvements. With a fresh look, the rating accounts for the fact that several design decisions (MovementBehavior, AOEAttack, CanMove) are speculative — they look right for TestBoss but haven't been stress-tested by a second boss with different needs. The true extensibility score will be known when boss #2 arrives.
