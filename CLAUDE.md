# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Quick Start

```bash
go run cmd/game/main.go      # Run the game
go build -o drill-game cmd/game/main.go  # Build executable
go test ./...                # Run all tests
./scripts/architecture-check.sh  # Verify hexagonal boundaries
./scripts/style-check.sh         # Verify Go style rules
```

## Architecture

**Hexagonal Architecture** with three layers:
- `internal/domain/` — Pure business logic (zero Raylib imports)
- `internal/adapters/` — Framework integration (Raylib)
- `cmd/game/main.go` — Application orchestration

**Key constraint:** Domain layer CANNOT import Raylib. Verify with `./scripts/architecture-check.sh`,
which also enforces that domain imports no adapters, no `cmd/`, and uses no `rl.*` types. A
`PostToolUse` hook runs it automatically on every edit under `internal/domain/` and blocks on failure.

## Project Structure

```
cmd/game/main.go           # Entry point
scripts/                   # Dev scripts (architecture-check.sh, style-check.sh)
.claude/                   # Skills, slash commands, rules, hooks, permissions
internal/
├── adapters/              # Raylib integration (input, rendering)
└── domain/                # Pure game logic
    ├── engine/            # Game loop orchestration
    ├── entities/          # Player, Building, Tile, ItemCatalog
    ├── upgrades/          # Upgrade types, Catalog, UpgradeType enum
    ├── systems/           # Physics, Drilling, Fuel, Items, Boss
    ├── effects/           # Player state mutations
    ├── ui/                # Shop interfaces
    ├── world/             # Chunk-based procedural generation
    ├── physics/           # Movement, collision, damage
    ├── bosses/            # Boss infrastructure (interfaces, base types, state machine)
    ├── boss_catalog/      # Boss implementations (one package per boss)
    ├── projectiles/       # Projectile spawning and movement types
    ├── config/            # Configuration structs
    ├── levels/            # Level definitions
    ├── types/             # Vec2, AABB
    ├── input/             # InputState
    └── components/        # Position, Interactable, Damageable
```

## Documentation Index

Read these docs on-demand when you need details:

| Doc | When to read |
|-----|--------------|
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | Understanding layer responsibilities, data flow, design principles |
| [SYSTEMS.md](docs/SYSTEMS.md) | Working on drilling, physics system, fuel, effects, UI, or items |
| [PHYSICS.md](docs/PHYSICS.md) | Working on collision, movement, damage formulas, or camera |
| [WORLD.md](docs/WORLD.md) | Working on world generation, tiles, chunks, or ore/hazard distributions |
| [BOSS.md](docs/BOSS.md) | Working on boss fights, state machines, phases, or adding new bosses |
| [CONFIGURATION.md](docs/CONFIGURATION.md) | Working on config structs, levels, or need reference tables for values |
| [GAME_DESIGN.md](docs/GAME_DESIGN.md) | Understanding game mechanics, progression, or player-facing features |
| [DEVELOPMENT.md](docs/DEVELOPMENT.md) | Testing, debugging, or development workflows |

## Tooling

**Skills** (`.claude/skills/`) — invoked by name or automatically when relevant:

| Skill | Use when |
|-------|----------|
| `architecture-check` | Finishing work under `internal/domain/`, reviewing a branch, auditing layering |

**Slash commands** (`.claude/commands/`):

| Command | Does |
|---------|------|
| `/new-boss <name>` | Scaffolds a boss end to end — domain package, states, renderer, level wiring |

**Hooks** (`.claude/settings.json`, `PostToolUse` on Write/Edit):

| Hook | Trigger | Effect |
|------|---------|--------|
| `gofmt` | any `*.go` | Formats in place |
| `hexagonal boundary check` | `internal/domain/**.go` | Runs `architecture-check.sh` — **blocks the edit** |
| `go style check` | any `*.go` | Runs `style-check.sh` on that file — **blocks the edit** |

Rules that must always hold live in `.claude/rules/` and are backed by a script
in `scripts/`, so they are enforced by the harness rather than by remembering.

**Permissions** — `.claude/settings.json` (checked in) allows **only commands
that are harmless in every case**: reads, queries, and the two audit scripts.
Anything that writes, deletes, stages, or executes repo code — `go build`,
`go test`, `gofmt -w`, `git add`, `find`, `chmod` — belongs in each developer's
own `.claude/settings.local.json` (gitignored). The versioned `deny` list is a
project-wide floor that local settings cannot override.

## Key Patterns

- **Data-driven config** — All parameters in `config/` structs, loaded via `levels.GetLevelConfig(n)`
- **Effects system** — All state mutations via `Effect` interface with `EffectContext` (player, world, damageables)
- **Component-based entities** — Buildings use `Position` + `Interactable`, bosses use `Damageable` for HP
- **Upgrade facade** — Player stats via `player.MaxSpeed()`, upgrades via `player.GetUpgrade(type)`
- **Unified upgrade catalog** — Single `Catalog` type handles all upgrade types
- **Damage through entity** — Player uses `player.DealDamage()`, bosses control own vulnerability via state machine
- **Public Game fields** — Renderer accesses `game.World`, `game.Player`, etc. directly (no getters)
- **UI composition** — Shop states embed `GridNavigator`; Market/Service/Inventory states embed `FirstFrameTracker`; Hospital/FuelStation share `ModalServiceProvider` interface
- **Null-object handlers** — Boss handlers (`PhaseChangeHandler`, `DamageReactionHandler`) default to no-ops; concrete bosses override only when needed

## Communication Style

- **Be direct and challenge me** — Push back when you think I'm wrong. No flattery, no sugarcoating, no sycophancy.

## Code Style

Full rules: **[.claude/rules/go-style.md](.claude/rules/go-style.md)** — mechanically enforced, read it before writing Go.

- **No `!` negation** — Use `myVar == false`, never `!myVar`. Applies to `if`, `for`, assignments, and `&&`/`||` operands. A `PostToolUse` hook **rejects the write** on violation.

## Special Levels

- Level 1: Production level
- Level -1: Test level (advanced stats for development)
- Level -2: Boss test level (TestBoss)
- Level -3: Sentinel boss test level
