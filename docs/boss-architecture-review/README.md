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

**Current State:** 7/10 - Good foundations but friction points:
- Must modify `game.go` factory for each new boss
- Must implement 13+ interface methods
- Manual wiring of PhaseManager, StateMachine, Movement, Attacks
- StateBehaviors callback pattern adds boilerplate

**Top Recommendations:**
1. Registration pattern (boss packages self-register via `init()`)
2. BaseBoss struct with common boilerplate
3. Typed state IDs instead of strings
4. Config struct per boss for all timing values

## Key Files to Read

When working on the boss system, these are the core files:

```
internal/domain/
├── bosses/
│   ├── boss.go              # Boss interface (start here)
│   ├── boxes.go             # Hit/Hurt/Collision box types
│   ├── phase.go             # Phase management
│   └── statemachine/        # State machine implementation
│       ├── types.go
│       └── machine.go
├── systems/
│   ├── boss_fight.go        # Room detection, contact damage
│   └── projectile_system.go # Projectile pool
└── engine/
    └── game.go              # Boss creation factory (createBossByType)
```

## Reference Implementation

The `test_boss` package shows the current pattern:

```
internal/domain/bosses/test_boss/
├── boss.go    # Struct, Boss interface implementation
└── states.go  # StateBehaviors callbacks, state definitions
```

## Related Documentation

- `docs/BOSS.md` - User-facing boss fight documentation
