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

**Current State:** 8.5/10 - Strong foundations with most friction points resolved.

**What's Been Implemented:**
- Registration pattern (boss packages self-register via `init()`)
- Typed state IDs (iota constants instead of strings)
- BoxSet system (pre-allocated boxes, zero GC pressure)
- Config constants per boss (all timing values at top of file)
- Single vulnerability source (`GetHurtboxes()` as source of truth)

**Remaining Opportunities:**
- BaseBoss struct could reduce boilerplate further
- StateBehaviors callback pattern retained (acceptable trade-off)

## Key Files to Read

When working on the boss system, these are the core files:

```
internal/domain/
├── bosses/
│   ├── boss.go              # Boss interface (start here)
│   ├── boxes.go             # BoxSet, Hit/Hurt/Collision box types
│   ├── registry.go          # Boss registration pattern
│   ├── phase.go             # Phase management
│   └── statemachine/        # State machine implementation
│       ├── types.go         # StateID (int), StateContext, StateResult
│       └── machine.go
├── systems/
│   ├── boss_fight.go        # Room detection, contact damage
│   └── projectile_system.go # Projectile pool
└── engine/
    └── game.go              # Uses bosses.Create() from registry
```

## Reference Implementation

The `test_boss` package shows the current pattern:

```
internal/domain/bosses/test_boss/
├── boss.go    # Struct, Boss interface, init() registration, config constants
└── states.go  # StateID iota constants, StateBehaviors, BuildStates()
```

## Adding a New Boss

1. Create package: `internal/domain/bosses/my_boss/`
2. Define config constants and phase configs at top of `boss.go`
3. Register via `init()`: `bosses.Register("my_boss", func(...) { return New(...) })`
4. Use `BoxSet` for pre-allocated hit/hurt/collision boxes
5. Define typed state IDs with `iota` in `states.go`
6. Build states via `BuildStates(behaviors)` pattern

No modifications to `game.go` or other core files required.

## Related Documentation

- `docs/BOSS.md` - User-facing boss fight documentation
