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
- `systems/physics.go` — Physics system (movement, gravity, collision)
- `systems/digging.go` — Tile destruction and player grid alignment
- `entities/player.go` — Player data and behavior
- `entities/tile.go` — Tile types (Empty, Dirt) with solid/diggable state
- `physics/` — Pure physics functions (movement, gravity, collision)
- `world/world.go` — Sparse tile map, grid-based terrain, collision queries
- `input/input_state.go` — Platform-agnostic input representation
- `types/vec2.go` — Custom 2D vector (no Raylib types)

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
1. renderer.GetFrameTime()       → delta time
2. inputAdapter.ReadInput()      → keyboard → InputState
3. game.Update(dt, inputState)   → digging system → physics system
4. renderer.Render(game)         → extract entities → draw with Raylib
```

### Physics System Order
1. Digging: Remove tiles, align player to grid (before physics)
2. Movement: Apply horizontal acceleration/damping (input-driven)
3. Vertical Movement: Apply upward thrust when jumping
4. Gravity: Apply downward acceleration every frame
5. Integration: Update position from velocity
6. Collision: Resolve tile collisions, snap to tile tops

### World & Collision
- **Tile Size**: 64×64 pixels
- **Sparse Map**: `map[[2]int]*entities.Tile` (only stores non-empty tiles)
- **Collision Detection**: Full-width check (left, center, right edges of player)
- **Tile-Based Ground**: Player stands on tiles, falls through dug tiles
- **Ground Level**: Must align to tile boundary (e.g., 640 = 10 × TileSize)

## Testing Strategy

**11 unit tests, zero framework dependencies** — all tests in `internal/domain/physics/`:

- `movement_test.go` — Acceleration, damping, speed capping
- `gravity_test.go` — Gravity effects, position integration
- `collision_test.go` — Ground collision, tile standing, takeoff

**Run tests:**
```bash
go test ./internal/domain/physics -v
```

**Key Testing Pattern**: Pure functions that take input values and return results. No Raylib initialization required.

## Physics Constants

Located in `internal/domain/physics/constants.go`:

```go
Gravity               = 800 pixels/sec²
MaxMoveSpeed          = 300 pixels/sec
MoveAcceleration      = 2500 pixels/sec²
MoveDamping           = 1000 pixels/sec²
FlyAcceleration       = 2000 pixels/sec²
MaxUpwardVelocity     = -300 pixels/sec (negative = upward)
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

**Phase 1 (Current)**: Core gameplay
- ✅ Game loop, player movement, physics
- ✅ Tile-based world, collision
- ✅ Digging system with player grid alignment

**Phase 2**: Progression system (ore types, upgrades, shop)
**Phase 3**: Polish (particles, sound, UI)
**Phase 4**: Extended content (more ores, hazards, achievements)

When adding features, maintain domain purity and hexagonal architecture. New game systems go in `internal/domain/systems/`, framework integration in `internal/adapters/`.
