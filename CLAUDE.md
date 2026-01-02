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

# Run single test file
go test -v ./internal/domain/physics -run TestApplyGravity

# Build and run
go build -o drill-game cmd/game/main.go && ./drill-game
```

## Architecture Overview

**Drill Game uses Hexagonal Architecture (Ports & Adapters)** with three clear layers:

### Domain Layer (`internal/domain/`)
**Pure business logic with ZERO Raylib dependencies** — fully testable without framework initialization.

- `engine/game.go` — Game orchestration, system coordination
- `systems/physics.go` — Physics system (movement, gravity, AABB collision)
- `systems/digging.go` — Tile destruction and player grid alignment
- `entities/player.go` — Player data (AABB-based) and behavior
- `entities/tile.go` — Tile types (Empty, Dirt, Ore) with solid/diggable state
- `entities/ore_type.go` — 7 ore types with Gaussian distribution parameters
- `physics/` — Pure physics functions (movement, gravity, AABB collision)
- `world/world.go` — Sparse tile map, chunk loading, terrain queries
- `world/generator.go` — Procedural tile generation with Gaussian ore distribution
- `world/hash.go` — Deterministic seeding for reproducible worlds
- `input/input_state.go` — Platform-agnostic input representation
- `types/vec2.go` — Custom 2D vector (no Raylib types)
- `types/aabb.go` — AABB collision primitive (no Raylib types)

### Adapter Layer (`internal/adapters/`)
**Framework integration — all Raylib code lives here, zero business logic.**

- `input/raylib.go` — Translates Raylib keyboard input to domain InputState
- `rendering/raylib.go` — Renders domain entities using Raylib graphics

### Application Layer (`cmd/game/main.go`)
**Simple orchestration: window management, main loop (Input → Update → Render)**

**Key Design Principles:**
- Domain layer has ZERO Raylib imports (checked at compile time)
- Adapters only translate, never contain business logic
- Pure functions in physics enable comprehensive testing
- Sparse tile map (only stores non-empty tiles) scales to large worlds

## Game Systems and Data Flow

### Per-Frame Flow (60 FPS)
```
1. renderer.GetFrameTime()            → delta time
2. inputAdapter.ReadInput()           → keyboard → InputState
3. game.Update(dt, inputState)        → chunk loading → digging → physics
4. renderer.Render(game)              → extract entities → draw with Raylib
```

### Game Update Order
1. **Chunk Loading**: Proactively load 3×3 grid of chunks around player (16×16 tile chunks)
2. **Digging**: Remove tiles, align player to grid (before physics)
3. **Physics System Order**:
   - Movement: Apply horizontal acceleration/damping (input-driven)
   - Vertical Movement: Apply upward thrust when jumping
   - Gravity: Apply downward acceleration every frame
   - **Axis-Separated Collision**: X-axis integration → collision → resolve, then Y-axis integration → collision → resolve

### World & AABB Collision
- **Tile Size**: 64×64 pixels
- **Sparse Map**: `map[[2]int]*entities.Tile` (only stores non-empty tiles)
- **Collision Type**: AABB (Axis-Aligned Bounding Box) for player and tiles
- **Detection**: CheckCollisions() finds all solid tiles intersecting player AABB (max 4 tiles)
- **Resolution**: Axis-separated (X then Y) to prevent corner-catching and enable wall sliding
- **Ground Detection**: Vertical collision with `dy > 0` (pushed upward) sets OnGround
- **Wall/Ceiling**: Horizontal collision stops X movement, ceiling collision stops Y movement

### Procedural World Generation
- **Chunk-Based**: 16×16 tile chunks generated on-demand when player enters area
- **Deterministic**: `seed + coordinates → unique RNG state` ensures reproducible worlds
- **Ore Distribution**: 7 ore types (Copper, Iron, Silver, Gold, Mythril, Platinum, Diamond)
  - Each ore has Gaussian distribution centered at depth (e.g., Gold peaks at depth 300)
  - Distribution formula: `weight = maxWeight × e^(-(depth - peak)² / (2σ²))`
- **Tile Mix Underground**:
  - 20% Empty (air pockets, caves)
  - 65% Dirt (filler, solid ground)
  - 15% Ore (procedurally distributed by depth using Gaussian weighting)
- **Above Ground**: Always Empty (sky)
- **Ground Level (Y=10)**: Always solid Dirt (safe spawning/landing)
- **Sparse Storage**: Empty tiles not stored (only solid tiles in map) — saves ~20% memory

## Testing Strategy

**31 unit tests + 7 integration tests, zero framework dependencies**:

**Physics & Types (17 tests)** — `internal/domain/physics/` and `internal/domain/types/`:
- `collision_test.go` — AABB collision detection, axis-separated resolution, wall/ceiling/ground
- `movement_test.go` — Acceleration, damping, speed capping
- `gravity_test.go` — Gravity effects, position integration
- `types/aabb_test.go` — AABB intersection, penetration calculation

**World Generation (14 tests)** — `internal/domain/world/`:
- `generator_test.go` — Gaussian distribution, determinism, ore selection, hash functions
- `world_test.go` — Chunk loading, lazy loading, proactive loading
- `integration_test.go` — End-to-end world generation, determinism, ore distribution validation

**Run specific test suite:**
```bash
go test ./internal/domain/physics -v
go test ./internal/domain/world -v
go test ./internal/domain/world -run "TestIntegration" -v
```

**Performance Benchmarks:**
```bash
go test ./internal/domain/world -bench=. -benchmem
# Chunk generation: ~2.2ms per 16×16 chunk
# Cached tile lookup: ~38ns per tile
```

**Key Testing Pattern**: Pure functions that take input values and return results. No Raylib initialization required.

## Physics Constants

Located in `internal/domain/physics/constants.go`:

```go
Gravity               = 800 pixels/sec²
MaxMoveSpeed          = 450 pixels/sec
MoveAcceleration      = 2500 pixels/sec²
MoveDamping           = 1000 pixels/sec²
FlyAcceleration       = 2500 pixels/sec²
MaxUpwardVelocity     = -600 pixels/sec (negative = upward)
FlyDamping            = 300 pixels/sec²
```

## Extending the System

### Adding a New Entity
1. Create `internal/domain/entities/newentity.go`
2. Implement `PhysicsEntity` interface (Position, Velocity, Height methods)
3. Add to `internal/domain/engine/game.go` Game struct
4. Update physics system to process it
5. Add rendering in `internal/adapters/rendering/raylib.go`

### Adding a New Game System
1. Create `internal/domain/systems/newsystem.go` package
2. Wire into `Game.Update()` in correct order
3. Test pure functions in `internal/domain/newsystem/` without framework
4. Keep framework code in adapters only

### Swapping Rendering Backend
1. Create new adapter: `internal/adapters/rendering/myframework.go`
2. Implement same interface as `RaylibRenderer` (Render, InitWindow, etc.)
3. Update `cmd/game/main.go` to use new adapter
4. Domain logic unchanged

## Input System

**InputState** (`internal/domain/input/input_state.go`):
```go
type InputState struct {
    Left  bool  // Move left
    Right bool  // Move right
    Up    bool  // Fly/jump
    Dig   bool  // Dig tile below
}
```

**Input Adapter Maps** (`internal/adapters/input/raylib.go`):
- `Left` ← Arrow Left OR A key
- `Right` ← Arrow Right OR D key
- `Up` ← Arrow Up OR W key
- `Dig` ← Down arrow OR S key (spacebar mapping can be updated here)

## Game Configuration

**Window & Display** (`cmd/game/main.go`):
- Screen: 1280×720 pixels
- Target FPS: 60
- Ground Level: 640.0 (10 × TileSize for clean grid alignment)

**Player** (`internal/domain/entities/player.go`):
- Size: 64×64 pixels (matches tile size)
- Start Position: Center X, just above ground level
- Max Move Speed: 300 px/sec
- Jump/Fly Speed: -300 px/sec (upward)

## Documentation Files

- **`docs/ARCHITECTURE.md`** (730 lines) — Detailed system design, data flow, testing strategy
- **`docs/GAME_DESIGN.md`** — Mechanics, progression, ore system, hazards, upgrades
- **`README.md`** — Project overview, roadmap, getting started

## Common Development Tasks

### Debugging Physics
1. Add `log/slog` output to `internal/domain/physics/` functions
2. Run via `go run cmd/game/main.go` (logs to stdout)
3. Check position/velocity values in movement, gravity, collision

### Testing Physics Changes
1. Write test in `internal/domain/physics/*_test.go`
2. Run: `go test ./internal/domain/physics -v -run TestName`
3. Verify builds and existing tests still pass: `go test ./...`

### Profiling Performance
Use Go's built-in profiling via `go test -cpuprofile=cpu.prof ./...` and analyze with `go tool pprof`.

## Code Style & Conventions

- **Domain Types**: Small types use values (Vec2), larger types use pointers (Player, World, Game)
- **Methods on Entities**: Data + behavior colocated (Player methods, Tile methods)
- **Pure Functions**: Physics functions are deterministic, no side effects
- **Interfaces**: Named, explicit contracts (PhysicsEntity)
- **Error Handling**: Minimal in game loop; domain functions can return errors
- **Constants**: Grouped logically (physics constants in one place, colors in adapter)

## Architecture Verification

To verify hexagonal architecture is maintained:
1. Domain layer imports: `internal/domain/*` only, no `internal/adapters`
2. Adapters import: `internal/domain/*` and framework code
3. Application layer: imports both adapters and domain
4. No Raylib imports in `internal/domain/`

Run: `grep -r "raylib" internal/domain/` — should return nothing (except comments).

## Roadmap Context

**Phase 1 (Complete)**: Core gameplay & world generation
- ✅ Game loop, player movement, physics
- ✅ Tile-based world, AABB collision system
- ✅ Digging system with player grid alignment
- ✅ Axis-separated collision (walls, ceiling, ground)
- ✅ Procedural chunk-based world generation
- ✅ 7 ore types with Gaussian depth distribution
- ✅ Deterministic seeding for reproducible worlds

**Phase 2**: Progression system (ore collection, inventory, upgrades, shop)
- Ore inventory system
- Mining duration per ore type
- Ore value/selling mechanics
- Upgrade system

**Phase 3**: Polish (particles, sound, UI)
- Particle effects for digging
- Sound effects and music
- HUD improvements
- Settings menu

**Phase 4**: Extended content (more ores, hazards, achievements)
- Biome variations
- Environmental hazards (lava, gas, etc.)
- Achievement system
- More ore types or special materials

When adding features, maintain domain purity and hexagonal architecture. New game systems go in `internal/domain/systems/`, framework integration in `internal/adapters/`.
