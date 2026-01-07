# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Quick Start Commands

```bash
go run cmd/game/main.go      # Run the game
go build -o drill-game cmd/game/main.go  # Build executable
go test ./...                # Run all tests
```

## Architecture Overview

**Hexagonal Architecture** — Three layers:
- `internal/domain/` — Pure business logic, zero framework dependencies
- `internal/adapters/` — Framework integration (Raylib)
- `cmd/game/main.go` — Application orchestration

**Key Constraint:** Domain layer CANNOT import Raylib. Verify: `grep -r "raylib" internal/domain/`

## Key Design Decisions

- **Player as Aggregate Root** — `Engine`, `Hull`, `FuelTank` are exported component value objects. Access stats via `player.Engine.MaxSpeed()`, not through wrapper methods.
- **Named Constructors** — Components use `NewEngineBase()`, `NewEngineMk1()`, etc. Tier data lives in constructors.
- **Upgrade Shops Own Catalogs** — Each shop type (`EngineUpgradeShop`, etc.) holds its catalog with prices and component instances.

## Key Files

- `internal/domain/engine/game.go` — Game loop orchestration
- `internal/domain/entities/player.go` — Player aggregate root
- `internal/domain/entities/engine.go`, `hull.go`, `fuel_tank.go` — Component value objects
- `internal/domain/entities/upgrade_shop.go` — Three shop types with catalogs
- `internal/domain/systems/` — Physics, digging, fuel, upgrades
- `internal/domain/world/` — Chunk-based procedural world

## Documentation

- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** — Complete technical architecture
- **[docs/GAME_DESIGN.md](docs/GAME_DESIGN.md)** — Game mechanics, upgrades, progression
- **[docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)** — Development workflows, testing
