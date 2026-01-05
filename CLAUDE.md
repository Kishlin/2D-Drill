# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Quick Start Commands

```bash
# Run the game (60 FPS, 1280x720 window)
go run cmd/game/main.go

# Build executable
go build -o drill-game cmd/game/main.go

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run physics tests
go test -v ./internal/domain/physics

# Run world tests
go test -v ./internal/domain/world

# Build and run
go build -o drill-game cmd/game/main.go && ./drill-game
```

## Architecture Overview

**Hexagonal Architecture (Ports & Adapters)** — Three layers:
- `internal/domain/` — Pure business logic, zero framework dependencies
- `internal/adapters/` — Framework integration (Raylib), no business logic
- `cmd/game/main.go` — Application orchestration

**Key Constraint:** Domain layer CANNOT import Raylib or adapters. Verify with:
```bash
grep -r "raylib" internal/domain/
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for complete technical details.

## Key Files by Task

**Game Logic:**
- `internal/domain/engine/game.go` — Game loop orchestration, system coordination
- `internal/domain/systems/physics.go` — Movement, gravity, AABB collision
- `internal/domain/systems/digging.go` — Tile destruction, ore collection, grid alignment
- `internal/domain/systems/fuel.go` — Fuel consumption based on activity level
- `internal/domain/systems/upgrade.go` — Upgrade purchasing at upgrade shops

**Entities:**
- `internal/domain/entities/player.go` — Player state (AABB, inventory, money, fuel, upgrades)
- `internal/domain/entities/tile.go` — Tile types (Empty, Dirt, Ore)
- `internal/domain/entities/shop.go` — Shop interaction (AABB-based)
- `internal/domain/entities/upgrade_shop.go` — Upgrade shop entity
- `internal/domain/entities/upgrade_tiers.go` — Upgrade tier definitions and stats

**World:**
- `internal/domain/world/world.go` — Sparse tile map, chunk loading
- `internal/domain/world/generator.go` — Procedural generation, Gaussian ore distribution
- `internal/domain/world/hash.go` — Deterministic seeding

**Adapters:**
- `internal/adapters/input/raylib.go` — Keyboard input translation
- `internal/adapters/rendering/raylib.go` — Raylib rendering

**Physics:**
- `internal/domain/physics/constants.go` — Physics tuning values
- `internal/domain/physics/` — Movement, gravity, collision functions

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for complete system descriptions.

## Game Configuration

### Physics Constants (`internal/domain/physics/constants.go`)

```
Gravity             = 800 pixels/sec²
MoveDamping         = 1000 pixels/sec²
FlyDamping          = 300 pixels/sec²
FallDamageThreshold = 500 pixels/sec (min speed to take damage)
FallDamageDivisor   = 20.0 (damage scaling)
```

**Note:** Movement speeds and accelerations are defined in `entities/upgrade_tiers.go` as part of the upgrade system. See [docs/GAME_DESIGN.md](docs/GAME_DESIGN.md) for upgrade details.

### Player Configuration

- **Size**: 54×54 pixels
- **Inventory**: 7 ore types (Copper, Iron, Silver, Gold, Mythril, Platinum, Diamond)
- **Upgrades**: Engine, Hull, Fuel Tank (each with Mk1-Mk5 tiers)

Base stats (upgradeable):
- **Max Move Speed**: 450 px/sec (Engine)
- **Max Upward Speed**: -600 px/sec (Engine)
- **Fuel Capacity**: 10.0 liters (Fuel Tank)
- **Max HP**: 10.0 (Hull)

### Shop Configuration

- **Position**: (960, 576) — 3 tiles right of player spawn, ground level
- **Size**: 320×192 pixels (5 tiles wide × 3 tiles tall)
- **Interaction**: E key to sell entire inventory when overlapping

### Ore Values & Distribution

| Ore | Value | Depth |
|-----|-------|-------|
| Copper | $10 | 50-100 tiles |
| Iron | $25 | 100-150 tiles |
| Silver | $75 | 150-250 tiles |
| Gold | $250 | 250-350 tiles |
| Mythril | $1000 | 350-450 tiles |
| Platinum | $5000 | 450-550 tiles |
| Diamond | $30000 | 550+ tiles |

**Distribution**: Gaussian around depth preference, rarer ores appear deeper.

### Fuel System

- **Base Tank Capacity**: 10.0 liters (upgradeable via Fuel Tank)
- **Active Consumption**: 0.333 L/sec (Left, Right, Up, or Dig inputs)
- **Idle Consumption**: 0.0833 L/sec (no movement inputs)

### Controls & Input

| Input | Action |
|-------|--------|
| **Left** (A or ←) | Move left / Dig left (when grounded) |
| **Right** (D or →) | Move right / Dig right (when grounded) |
| **Up** (W or ↑) | Jump/Fly (hold to fly) |
| **Dig** (S or ↓) | Dig downward, snap to grid |
| **Sell** (E) | Sell inventory at shop |

## Development Quick Reference

**Running/Testing:**
- Run game: `go run cmd/game/main.go`
- All tests: `go test ./...`
- Physics tests: `go test ./internal/domain/physics -v`
- World tests: `go test ./internal/domain/world -v`
- Single test: `go test -v ./path -run TestName`

**Debugging:**
- Add logs via `log/slog`: `slog.Info("message", "key", value)`
- Raylib check: `grep -r "raylib" internal/domain/` (should be empty)
- Breakpoint debug: `dlv debug cmd/game/main.go`

**Performance:**
- Benchmarks: `go test ./internal/domain/world -bench=. -benchmem`
- CPU profile: `go test -cpuprofile=cpu.prof ./...`
- Analyze: `go tool pprof cpu.prof`

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for detailed workflows.

## Design Principles (Condensed)

1. **Framework Independence** — Domain layer has zero Raylib dependencies
2. **Testability** — All business logic testable without framework initialization
3. **Separation of Concerns** — Domain (logic) → Adapters (framework) → Application (orchestration)
4. **Value Types** — Small types (Vec2, AABB) as values, large types (Player, Game) as pointers
5. **Pure Functions** — Physics functions are deterministic, no side effects
6. **Clear Data Flow** — Unidirectional: domain → adapters → application

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for full rationale and examples.

## Current Implementation Status

**Phase 1 (Complete):**
- ✅ Game loop, player movement, physics
- ✅ AABB collision system, axis-separated resolution
- ✅ Tile-based world with sparse storage
- ✅ Procedural chunk-based generation (deterministic seeding)
- ✅ 7 ore types with Gaussian depth distribution

**Phase 2 (In Progress):**
- ✅ Ore inventory system with real-time display
- ✅ Shop entity and selling mechanics
- ✅ Money system
- ✅ Fuel system (consumption based on activity)
- ✅ Fuel station (refueling with cost)
- ✅ Fall damage system (10 HP with 500 px/sec threshold)
- ✅ Hospital (healing HP for $2 per HP)
- ✅ Upgrade system (Engine, Hull, Fuel Tank with Mk1-Mk5 tiers)
- ⏳ Mining duration per ore type
- ⏳ Game-over behavior (fuel depletion, HP reaching 0)

**Phase 3+:** Polish, hazards, more content

See [README.md](README.md) for full roadmap.

## Documentation Guide

- **[CLAUDE.md](CLAUDE.md)** — This file, quick reference for AI assistants
- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** — Complete technical architecture and design
- **[docs/GAME_DESIGN.md](docs/GAME_DESIGN.md)** — Game mechanics, progression, upgrades
- **[docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)** — Development workflows, testing, debugging
- **[README.md](README.md)** — Project overview, installation, roadmap

## Key Design Decision

**Single Source of Truth:** Each piece of information lives in exactly one documentation file. CLAUDE.md is a reference that links to detailed docs, not a copy. This maintains consistency and reduces context bloat.
