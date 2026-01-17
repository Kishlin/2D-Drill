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

- **Player as Aggregate Root** — `Engine`, `Hull`, `FuelTank`, `CargoHold`, `HeatShield`, `Drill` are exported component value objects. Access stats via `player.Engine.MaxSpeed()`, not through wrapper methods. Damage mutations go through `player.DealDamage(damage)`.
- **Named Constructors** — Components use `NewEngineBase()`, `NewEngineMk1()`, etc. Tier data lives in constructors.
- **Upgrade Shops Own Catalogs** — Each shop type (`EngineUpgradeShop`, `HullUpgradeShop`, `FuelTankUpgradeShop`, `CargoHoldUpgradeShop`, `HeatShieldUpgradeShop`, `DrillUpgradeShop`) holds its catalog with prices and component instances.
- **Cargo Capacity Limits** — `AddOre()` respects cargo hold capacity; ore is lost when full (intentional Motherload-style behavior).
- **Damage Application** — All damage sources (fall, heat, future hazards) call `player.DealDamage(damage)` which applies damage and clamps HP at zero. Physics package calculates damage; Player entity applies it.
- **Heat System** — Temperature increases with depth; players take exponential damage when temperature exceeds heat resistance. Heat shield is an upgradeable component enabling deeper mining.
- **Drilling Animation** — Both vertical and horizontal drilling is animated with variable duration based on tile hardness and depth. Formula: `duration = baseTime × hardness × depthFactor / drillSpeed`. Hardness values: `DirtHardness` (1.0) in tile.go, `OreHardness` map (1.2-3.0) in ore_type.go, `HazardHardness` map in hazard_type.go. Depth factor scales 1.0→24.0 from surface to max depth. Drill upgrades apply depth-scaled divisor (more effective at depth than surface). Lava tiles use fixed 0.3s duration (depth-independent) but deal 100 damage on completion (reduced by heat shield). Player progressively moves to tile center via lerp. Tile only removed on animation completion. All inputs blocked during drill except fuel/heat (continuous). See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md#drilling-system) for details.
- **Hazard Tiles** — Rock (impenetrable, blocks drilling/movement) and Lava (drillable but deals damage) tiles spawn at deep depths using Gaussian distributions. Hazards dominate terrain at 80%+ depths, creating natural progression gates. Bombs can destroy rocks via `NukeTileAtGrid()` which bypasses drillability checks. See [docs/GAME_DESIGN.md](docs/GAME_DESIGN.md#hazard-tiles) for mechanics.
- **Depth-Dependent Generation** — All tile types (empty, dirt, ore, hazards) use weighted random selection with weights that vary by depth. Surface (0%) generates mostly dirt/empty; deep (80%+) hazards dominate. Enables natural progression without hard level walls.
- **Item System** — Consumable items (Teleport, Repair, Refuel, Bomb, Big Bomb) stored in `ItemInventory [5]int`. Items used via T, R, F, B, G keys (discrete input, `IsKeyPressed()`). Bombs and Big Bombs use `NukeTileAtGrid()` to destroy all solid tiles including impenetrable rocks. Purchased at dedicated item shops with E key. See [docs/GAME_DESIGN.md](docs/GAME_DESIGN.md#items) for mechanics.

## Key Files

- `internal/domain/engine/game.go` — Game loop orchestration
- `internal/domain/entities/player.go` — Player aggregate root
- `internal/domain/entities/engine.go`, `hull.go`, `fuel_tank.go`, `cargo_hold.go`, `heat_shield.go`, `drill.go` — Component value objects
- `internal/domain/entities/upgrade_shop.go` — Six shop types with catalogs (Engine, Hull, FuelTank, CargoHold, HeatShield, Drill)
- `internal/domain/entities/item.go` — ItemType enum and item names
- `internal/domain/entities/item_shop.go` — ItemShop entity for purchasing consumables
- `internal/domain/systems/` — Physics, drilling, fuel, upgrades, items
- `internal/domain/systems/item_shop.go` — ItemShopSystem for item purchases
- `internal/domain/world/` — Chunk-based procedural world

## Documentation

- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** — Complete technical architecture
- **[docs/GAME_DESIGN.md](docs/GAME_DESIGN.md)** — Game mechanics, upgrades, progression
- **[docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)** — Development workflows, testing
