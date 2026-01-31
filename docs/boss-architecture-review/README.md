# Boss Architecture Review

This folder contains a comprehensive analysis of the boss fight system, created to preserve context across Claude sessions.

## Contents

| File | Description |
|------|-------------|
| `ANALYSIS.md` | Critical review of the architecture with ratings, issues, and recommendations |
| `IMPLEMENTATION_REFERENCE.md` | Detailed documentation of current implementation (interfaces, structs, flow) |
| `README.md` | This file |

## Quick Context for New Sessions

**Goal:** Make adding new bosses as simple as "create a subpackage with states and transitions."

**Current State:** 9/10 - Strong foundations with BaseBoss reducing boilerplate significantly.

**What's Been Implemented:**
- Registration pattern (boss packages self-register via `init()`)
- Typed state IDs (iota constants instead of strings)
- BoxSet system (pre-allocated boxes, zero GC pressure)
- Config constants per boss (all timing values at top of file)
- Single vulnerability source (`GetHurtboxes()` as source of truth)
- **BaseBoss struct** (reduces ~80 lines of boilerplate per boss)
- **Phases package** (clean separation of phase management)
- **Boss catalog** (boss implementations separate from infrastructure)

## Key Files to Read

When working on the boss system, these are the core files:

```
internal/domain/
├── bosses/                  # Boss infrastructure
│   ├── boss.go              # Boss interface (start here)
│   ├── base_boss.go         # BaseBoss struct with default implementations
│   ├── boxes.go             # BoxSet, Hit/Hurt/Collision box types
│   ├── registry.go          # Boss registration pattern
│   ├── phases/              # Phase management package
│   │   └── phase.go         # phases.Config, phases.Manager
│   └── statemachine/        # State machine implementation
│       ├── types.go         # StateID (int), StateContext, StateResult
│       └── machine.go
├── boss_catalog/            # Boss implementations
│   └── test_boss/
│       ├── boss.go          # TestBoss (embeds BaseBoss)
│       └── states.go        # State IDs + BuildStates()
├── systems/
│   ├── boss_fight.go        # Room detection, contact damage
│   └── projectile_system.go # Projectile pool
└── engine/
    └── game.go              # Uses bosses.Create() from registry
```

## Reference Implementation

The `test_boss` package shows the current pattern:

```
internal/domain/boss_catalog/test_boss/
├── boss.go    # Embeds *bosses.BaseBoss, implements handlers
└── states.go  # StateID iota constants, StateBehaviors, BuildStates()
```

## Adding a New Boss

1. Create package: `internal/domain/boss_catalog/my_boss/`
2. Define config constants and phase configs at top of `boss.go`
3. Create boss struct embedding `*bosses.BaseBoss`
4. Register via `init()`: `bosses.Register("my_boss", func(...) { return New(...) })`
5. Implement `PhaseChangeHandler` and `DamageReactionHandler` interfaces
6. Override `GetHurtboxes()` for custom vulnerability logic
7. Define typed state IDs with `iota` in `states.go`
8. Build states via `BuildStates(behaviors)` pattern

No modifications to `game.go` or other core files required.

## Related Documentation

- `docs/BOSS.md` - User-facing boss fight documentation
