# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Quick Start

```bash
go run cmd/game/main.go      # Run the game
go build -o drill-game cmd/game/main.go  # Build executable
go test ./...                # Run all tests
```

## Architecture

**Hexagonal Architecture** with three layers:
- `internal/domain/` — Pure business logic (zero Raylib imports)
- `internal/adapters/` — Framework integration (Raylib)
- `cmd/game/main.go` — Application orchestration

**Key constraint:** Domain layer CANNOT import Raylib. Verify: `grep -r "raylib" internal/domain/`

## Project Structure

```
cmd/game/main.go           # Entry point
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
    ├── bosses/            # Boss implementations
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

## Code Style

- **Explicit false booleans** — Use `if myVar == false` instead of `if !myVar`

## Special Levels

- Level 1: Production level
- Level -1: Test level (advanced stats for development)
- Level -2: Boss test level
