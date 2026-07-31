---
description: Scaffold a new boss end to end — domain package, state machine, phases, renderer, and level wiring.
argument-hint: <boss_name> (snake_case, e.g. crawler, magma_worm)
---

Scaffold a new boss: **$1**

If `$1` is empty, ask which boss to create before doing anything else.

Read [docs/BOSS.md](docs/BOSS.md) first for the interfaces. Work through all five
steps — a boss that is missing its renderer compiles fine and fails at **runtime**
with `unknown boss type: $1`, because the renderer import is what triggers
registration (see step 3).

## Layer placement (verified against the codebase)

| File | Package | Path |
|---|---|---|
| Boss + states builder | `$1` | `internal/domain/boss_catalog/$1/boss.go` |
| State IDs | `$1` | `internal/domain/boss_catalog/$1/states.go` |
| Tests | `$1` | `internal/domain/boss_catalog/$1/boss_test.go` |
| Renderer | `bosses` | `internal/adapters/rendering/bosses/$1.go` |
| Level config | `levels` | `internal/domain/levels/level_*.go` |

The domain half must import **no Raylib** — a `PostToolUse` hook blocks the edit
if it does. All visuals live in the renderer.

## 1. Domain package

`boss.go` embeds `*bosses.BaseBoss` and registers the constructor in `init()`:

```go
func init() {
    bosses.Register("$1", func(roomStartY, worldWidth float32) bosses.Boss {
        return New(roomStartY, worldWidth)
    })
}
```

`bosses.Register` **panics on a duplicate type name**, so pick a name no other
boss in `boss_catalog/` uses.

In `New`, after building the `BaseBoss`:

```go
b := &Boss{BaseBoss: baseBoss}
b.Self = b // REQUIRED — without it GetHurtboxes/TakeDamageAt dispatch to the base, and the boss is unkillable
b.SetStateMachine(statemachine.NewStateMachine(b.buildStates(), StateIdle))
```

Set `b.PhaseChangeHandler = b` / `b.DamageReactionHandler = b` **only** if the
boss actually reacts to those events — they default to no-ops, and wiring them
without implementing meaningful behaviour is noise.

Keep the split the existing bosses use: `phases.Config` holds only HP
thresholds (generic infrastructure), while a boss-local
`<name>PhaseConfig` slice holds the per-phase tuning (speeds, cooldowns).

## 2. States

State IDs go in `states.go` as typed constants:

```go
const (
    StateIdle statemachine.StateID = iota
    StateAttack
    StateVulnerable
)
```

The state map goes in `boss.go` as `func (b *Boss) buildStates() map[statemachine.StateID]*statemachine.State`
— a method, so handlers get direct field access instead of threading state
through the context. Return `statemachine.StateResult{NextState: statemachine.StateIDNone}`
to stay put.

Model vulnerability as a state, not a flag: return an empty hurtbox slice from
`GetHurtboxes()` while invulnerable. Compare `sentinel_boss/states.go`.

## 3. Renderer — this is what registers the boss

`internal/adapters/rendering/bosses/$1.go`, `package bosses`:

```go
func init() {
    Register(&BossRenderer{})
}

func (r *BossRenderer) CanRender(boss bosses.Boss) bool {
    _, ok := boss.(*$1.Boss)
    return ok
}
```

Two things that trip people up here:

- The package is **`bosses`**, not `bossrenderers`. That is only the import alias
  used at `internal/adapters/rendering/raylib.go:8`. Inside the package, the
  unaliased `bosses` identifier refers to the *domain* `bosses` package.
- This file's import of `boss_catalog/$1` is the **only** reference to the boss
  package in the whole binary. It is what runs the domain `init()` and populates
  the registry. Skip the renderer and `bosses.Create("$1", ...)` fails at runtime.

Use `ok == false` for the failed type assertion in `Render` — this codebase
never uses `!`, and a hook will reject the write. See
[.claude/rules/go-style.md](.claude/rules/go-style.md).

## 4. Level wiring

Add a boss test level next to `level_boss.go` / `level_boss2.go`, register the
level number in `internal/domain/levels/registry.go`, and point the room at the
new type:

```go
BossRoom: config.BossRoomConfig{
    BossType:    "$1",
    FloorType:   config.FloorConcrete,
    FloorDamage: 10.0,
    RoomHeight:  680.0,
    FloorHeight: 6.0,
},
```

Levels -2 and -3 are already taken by TestBoss and SentinelBoss.

## 5. Tests

Colocate `boss_test.go` in the boss package, following
`boss_catalog/sentinel_boss/boss_test.go`. Cover at minimum: phase transitions
fire at the right HP thresholds, hurtboxes are empty while invulnerable, and
each state reaches its intended successor. Assert on effects, not merely that a
call did not panic.

## 6. Verify

```bash
go build ./...
go test ./internal/domain/...
./scripts/architecture-check.sh
./scripts/style-check.sh
```

Then confirm the boss actually resolves at runtime — registration failures do
not show up in any of the above:

```bash
# point cmd/game/main.go at the new level, then:
go run cmd/game/main.go
```

Finally, document the boss in `docs/BOSS.md` alongside the TestBoss and
SentinelBoss reference sections.
