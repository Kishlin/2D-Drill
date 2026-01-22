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

- **Data-Driven Configuration** — All game parameters are defined in `GameConfig` struct within the `config` package. Levels define complete configurations including world, player, generation, upgrades, and items. No hardcoded values in domain code—everything flows from config.
- **Levels System** — Each level is defined in `internal/domain/levels/` as a function returning `*config.GameConfig`. Use `levels.GetLevelConfig(levelNum)` to load. Level -1 is a test level with advanced player stats for development.
- **Config-Driven Constructors** — Components use `NewEngineFromStats(stats)`, `NewPlayerFromConfig(cfg)`, etc. Stats come from config, not hardcoded in constructors. Legacy named constructors (`NewEngineBase()`) still exist for backward compatibility.
- **Player as Aggregate Root** — `Engine`, `Hull`, `FuelTank`, `CargoHold`, `HeatShield`, `Drill` are exported component value objects. Access stats via `player.Engine.MaxSpeed()`, not through wrapper methods. Damage mutations go through `player.DealDamage(damage)`.
- **Component-Based Architecture** — Buildings use components (`Position`, `Interactable`) rather than inheritance. Generic `Building` entity with `InteractableType` enum enables unified handling. Components live in `internal/domain/components/`.
- **Effects System** — All player state mutations go through `Effect` interface. Concrete effects (`TakeMoney`, `AddMoney`, `SetFuel`, `SetHP`, `SetEngine`, etc.) are processed by `Processor`. This decouples UI from state changes. Effects live in `internal/domain/effects/`.
- **UI Layer** — Unified UI system with `Manager` handling both modal (shops) and instant (market/hospital/fuel) interactions. `UI` interface returns `Result` with `ShouldClose` and `Effects`. Modal UIs stay open; instant UIs close immediately with effects. UI code lives in `internal/domain/ui/`.
- **Upgrade Catalog** — `UpgradeCatalog` (extracted from old shop) holds all 6 upgrade types with prices and component instances. Accessed via modal `UpgradeShopUI` with tab cycling, grid navigation, and no sequential purchase requirement—players can skip tiers with sufficient funds.
- **Cargo Capacity Limits** — `AddOre()` respects cargo hold capacity; ore is lost when full (intentional Motherload-style behavior).
- **Damage Application** — All damage sources (fall, heat, future hazards) call `player.DealDamage(damage)` which applies damage and clamps HP at zero. Physics package calculates damage; Player entity applies it.
- **Heat System** — Temperature increases with depth; players take exponential damage when temperature exceeds heat resistance. Heat shield is an upgradeable component enabling deeper mining.
- **Drilling Animation** — Both vertical and horizontal drilling is animated with variable duration based on tile hardness and depth. Formula: `duration = baseTime × hardness × depthFactor / drillSpeed`. Hardness values come from config: `DirtHardness` in `GenerationConfig`, ore hardness in each `OreConfig`, hazard properties in `HazardConfig`. Depth factor scales 1.0→24.0 from surface to max depth. Drill upgrades apply depth-scaled divisor (more effective at depth than surface). Lava tiles use `FixedDuration` from config (depth-independent) and deal `OnDrillDamage` on completion (reduced by heat shield). Player progressively moves to tile center via lerp. Tile only removed on animation completion. All inputs blocked during drill except fuel/heat (continuous). See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md#drilling-system) for details.
- **Hazard Tiles** — Rock (impenetrable, blocks drilling/movement) and Lava (drillable but deals damage) tiles spawn at deep depths using Gaussian distributions. Hazards dominate terrain at 80%+ depths, creating natural progression gates. Bombs can destroy rocks via `NukeTileAtGrid()` which bypasses drillability checks. See [docs/GAME_DESIGN.md](docs/GAME_DESIGN.md#hazard-tiles) for mechanics.
- **Depth-Dependent Generation** — All tile types (empty, dirt, ore, hazards) use weighted random selection with weights that vary by depth. Surface (0%) generates mostly dirt/empty; deep (80%+) hazards dominate. Enables natural progression without hard level walls.
- **Item Catalog** — `ItemCatalog` (extracted from old shop) holds all 5 consumable items (Teleport, Repair, Refuel, Bomb, Big Bomb) with prices. Items stored in `ItemInventory [5]int` and used via T, R, F, B, G keys (discrete input, `IsKeyPressed()`). Bombs and Big Bombs use `NukeTileAtGrid()` to destroy all solid tiles including impenetrable rocks. Purchased via modal `ItemShopUI` at item shop with E key. Grid navigation (arrows/WASD) selects items; E to purchase; Q/Escape to close. See [docs/GAME_DESIGN.md](docs/GAME_DESIGN.md#items) for mechanics.
- **Boss Fight System** — End-of-level boss encounters with configurable boss types and floor mechanics. Boss rooms are generated as empty spaces below the mining area, with solid, indestructible floor tiles. Bosses implement the `Boss` interface; bomb-vulnerable bosses implement `PhysicalBoss` for AABB collision detection (includes `IsVulnerable()`, `GetContactDamage()`). Game state tracks Playing/Victory/Defeat. Boss HP bar renders at screen top when boss is active. See [docs/GAME_DESIGN.md](docs/GAME_DESIGN.md#boss-fights) for full mechanics.
- **Boss State Machine** — Bosses use state machines for animation and behavior. TestBoss has states: `StatePatrol` (moving/shooting), `StateWindup` (vibrating warning), `StateSlam` (AOE damage), `StateVulnerable` (can be bombed). State transitions drive visual feedback and vulnerability windows.
- **Boss Phase System** — `PhaseManager` tracks HP-threshold based phases. Each phase can configure movement speed, attack cooldowns, vulnerability rules. TestBoss has 3 phases: Phase 1 (always vulnerable), Phase 2 (vulnerable after slam), Phase 3 (shorter windows, double slams).
- **Boss Rendering Architecture** — Each boss type has its own renderer in `internal/adapters/rendering/bosses/`. Renderers type-assert to concrete boss types for state access. No generic animation interfaces polluting domain—boss-specific logic stays in boss-specific files.

## Key Files

### Configuration & Levels
- `internal/domain/config/` — All configuration structs (`GameConfig`, `WorldConfig`, `PlayerConfig`, `GenerationConfig`, `UpgradeConfig`, `ItemConfig`)
- `internal/domain/levels/` — Level definitions (`level1.go`, `level_dev.go` for test level -1, `registry.go`)
- `cmd/game/main.go` — Application entry point; loads level config via `levels.GetLevelConfig(levelNum)`

### Core Domain
- `internal/domain/engine/game.go` — Game loop orchestration with update order: chunks → UI manager (modal pause) → physics → fuel → drilling → items → boss
- `internal/domain/entities/player.go` — Player aggregate root with `InShop` and `IsDrilling` pause flags; `NewPlayerFromConfig()` constructor
- `internal/domain/entities/engine.go`, `hull.go`, `fuel_tank.go`, `cargo_hold.go`, `heat_shield.go`, `drill.go` — Component value objects with `NewXFromStats()` config-driven constructors
- `internal/domain/entities/building.go` — Generic `Building` entity with `Position` and `Interactable` components; factory functions for all building types
- `internal/domain/entities/catalog.go` — `UpgradeCatalog` and `ItemCatalog` with prices and component instances

### Components
- `internal/domain/components/position.go` — `Position` component wrapping AABB for collision
- `internal/domain/components/interactable.go` — `InteractableType` enum (Market, FuelStation, Hospital, UpgradeShop, ItemShop)

### Effects
- `internal/domain/effects/effect.go` — `Effect` interface for player state mutations
- `internal/domain/effects/money.go` — `TakeMoney`, `AddMoney` effects
- `internal/domain/effects/stats.go` — `SetFuel`, `SetHP` effects
- `internal/domain/effects/upgrades.go` — `SetEngine`, `SetHull`, `SetFuelTank`, `SetCargoHold`, `SetHeatShield`, `SetDrill` effects
- `internal/domain/effects/inventory.go` — `ClearOreInventory`, `AddItem` effects
- `internal/domain/effects/processor.go` — `Processor` applies effects to player

### UI
- `internal/domain/ui/ui.go` — `UI` interface, `Result` type with `ShouldClose` and `Effects`
- `internal/domain/ui/manager.go` — `Manager` handles UI registration, opening, and processing
- `internal/domain/ui/state.go` — `UpgradeShopState`, `ItemShopState` with navigation methods
- `internal/domain/ui/upgrade_shop.go` — Modal `UpgradeShopUI` with tab cycling, grid navigation, purchase logic
- `internal/domain/ui/item_shop.go` — Modal `ItemShopUI` with grid navigation, purchase logic
- `internal/domain/ui/market.go` — Instant `MarketUI` (sells ore, returns effects immediately)
- `internal/domain/ui/hospital.go` — Instant `HospitalUI` (heals player, returns effects immediately)
- `internal/domain/ui/fuel_station.go` — Instant `FuelStationUI` (refuels player, returns effects immediately)

### Systems
- `internal/domain/systems/interaction.go` — `DetectInteraction()` function checks player-building overlap with E key
- `internal/domain/systems/drilling.go` — Drilling animations with depth-scaled durations; reads hazard config for lava damage
- `internal/domain/systems/fuel.go` — Fuel consumption based on input state
- `internal/domain/systems/physics.go` — Gravity, collision, damage (fall/heat)
- `internal/domain/systems/boss_fight.go` — BossFightSystem: orchestrates boss fights, tracks player entry/exit, handles projectile collisions, applies floor damage
- `internal/domain/world/` — Chunk-based procedural world; generator reads `GenerationConfig` for ore/hazard distributions

### Boss System
- `internal/domain/bosses/boss.go` — Core interfaces (`Boss`, `PhysicalBoss`) and shared types (`AOEInfo`)
- `internal/domain/bosses/projectile.go` — Projectile entity for boss attacks with AABB collision detection
- `internal/domain/bosses/phase.go` — `PhaseManager` for HP-threshold based phase transitions
- `internal/domain/bosses/attacks/` — Attack system (reusable): `Attack` interface, `ProjectileAttack`, `AOEAttack`
- `internal/domain/bosses/movement/` — Movement system (reusable): `MovementBehavior` interface, `Grounded` patrol
- `internal/domain/bosses/test_boss/boss.go` — TestBoss implementation with state machine, phases, movement, attacks
- `internal/domain/config/boss_room_config.go` — `BossRoomConfig` struct with boss type, floor type, room height, floor height
- `internal/domain/entities/game_state.go` — `GameState` enum: Playing, Victory, Defeat
- `internal/domain/levels/level_boss.go` — Boss test level (-2) configuration for development and testing
- `internal/adapters/rendering/bosses/` — Boss-specific renderers (each boss has its own rendering logic)

## Documentation

- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** — Complete technical architecture
- **[docs/GAME_DESIGN.md](docs/GAME_DESIGN.md)** — Game mechanics, upgrades, progression
- **[docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)** — Development workflows, testing
