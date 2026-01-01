# Architecture

## Overview

Drill Game uses **Hexagonal Architecture** (Ports & Adapters) to achieve a clean separation between pure domain logic and framework-specific integration. This ensures the game logic is:

- **Testable**: Physics and game logic can be tested without initializing Raylib
- **Portable**: Domain code has zero framework dependencies
- **Maintainable**: Clear responsibilities and data flow
- **Extensible**: Easy to add new features, entities, or rendering backends

---

## Architecture Pattern: Hexagonal Architecture

```
┌─────────────────────────────────────────────────────┐
│                   APPLICATION LAYER                 │
│                  (cmd/game/main.go)                 │
│    Orchestrates: Input Reading → Game Update → Rendering
└────────────┬──────────────────────────────┬─────────┘
             │                              │
      ┌──────▼─────────┐        ┌──────────▼──────┐
      │ INPUT ADAPTER  │        │ RENDERING ADAPTER│
      │ (Raylib Keys)  │        │  (Raylib Drawing)
      └──────┬─────────┘        └──────────┬──────┘
             │                              │
             │      ┌──────────────────┐    │
             │      │   DOMAIN LAYER   │    │
             │      │  (Pure Business  │    │
             └─────▶│     Logic)       │◀───┘
                    │                  │
                    │ • Game (engine/) │
                    │ • Player (entities/)
                    │ • Physics (systems/)
                    │ • Types (types/)
                    │ • Input State (input/)
                    │ • World (world/)
                    │
                    └──────────────────┘
```

### The Three Layers

1. **Application Layer** (`cmd/game/main.go`)
   - Orchestrates the entire game
   - Reads input from adapters
   - Updates domain logic
   - Delegates rendering to adapters
   - Manages window lifecycle

2. **Adapter Layer** (`internal/adapters/`)
   - **Input Adapter**: Translates Raylib keyboard input to domain `InputState`
   - **Rendering Adapter**: Takes domain entities and renders them with Raylib
   - All Raylib integration lives here
   - No business logic

3. **Domain Layer** (`internal/domain/`)
   - Pure game logic, zero framework dependencies
   - Can be tested without Raylib initialization
   - Fully portable (could swap Raylib for SDL, canvas, etc.)

---

## Project Structure

```
drill-game/
├── cmd/
│   └── game/
│       └── main.go                          # Application orchestration
│
├── internal/
│   ├── adapters/                            # Framework Integration (Raylib)
│   │   ├── input/
│   │   │   └── raylib.go                    # RaylibInputAdapter
│   │   └── rendering/
│   │       └── raylib.go                    # RaylibRenderer
│   │
│   └── domain/                              # Pure Business Logic
│       ├── engine/
│       │   └── game.go                      # Game orchestration (domain)
│       ├── systems/
│       │   └── physics.go                   # PhysicsSystem
│       ├── entities/
│       │   ├── player.go                    # Player entity (data + behavior)
│       │   └── contracts.go                 # PhysicsEntity interface
│       ├── physics/
│       │   ├── constants.go                 # Physics parameters
│       │   ├── movement.go                  # Movement functions
│       │   ├── gravity.go                   # Gravity + velocity integration
│       │   ├── collision.go                 # Collision detection/resolution
│       │   ├── movement_test.go             # Movement tests
│       │   ├── gravity_test.go              # Gravity tests
│       │   └── collision_test.go            # Collision tests
│       ├── types/
│       │   └── vec2.go                      # Custom Vec2 (no Raylib types)
│       ├── input/
│       │   └── input_state.go               # InputState struct (framework-agnostic)
│       └── world/
│           └── world.go                     # World data and methods
│
├── docs/
│   ├── ARCHITECTURE.md                      # This file
│   └── GAME_DESIGN.md                       # Game mechanics
├── go.mod
├── go.sum
└── README.md
```

---

## Data Flow

### Single Frame Update

```
main.go Loop:
┌─────────────────────────────────────────┐
│ 1. Read Input from Adapter              │
│    inputState := inputAdapter.ReadInput()
│    (Converts Raylib keys → InputState)  │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│ 2. Update Domain Logic                  │
│    game.Update(dt, inputState)          │
│    • Physics system applies forces      │
│    • Gravity                            │
│    • Movement                           │
│    • Collision detection & resolution   │
│    • Updates Player position/velocity   │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│ 3. Render via Adapter                   │
│    renderer.Render(game)                │
│    • Extracts Player, World from game   │
│    • Converts to Raylib types           │
│    • Draws everything                   │
└─────────────────────────────────────────┘
```

### Code Example

```go
// main.go - Application layer orchestration
for renderer.WindowShouldClose() == false {
    dt := renderer.GetFrameTime()

    // 1. Read input (adapter responsibility)
    inputState := inputAdapter.ReadInput()

    // 2. Update game (domain logic - pure, testable)
    err := game.Update(dt, inputState)
    if err != nil {
        return err
    }

    // 3. Render (adapter responsibility)
    renderer.Render(game)
}
```

---

## Core Concepts

### Domain Layer (Pure Game Logic)

#### Game (`domain/engine/game.go`)

Orchestrates domain systems without any framework knowledge:

```go
type Game struct {
    world         *world.World
    player        *entities.Player
    physicsSystem *systems.PhysicsSystem
}

func (g *Game) Update(dt float32, inputState input.InputState) error {
    // Pure domain logic - zero Raylib
    g.physicsSystem.UpdatePhysics(g.player, inputState, dt)
    return nil
}

// Adapters pull data from getters
func (g *Game) GetWorld() *world.World {
    return g.world
}

func (g *Game) GetPlayer() *entities.Player {
    return g.player
}
```

**Why this design:**
- Game doesn't import Raylib (fully testable)
- Game accepts InputState (not Raylib keys)
- Game provides getters for adapters to read state
- All rendering responsibility is external

#### Physics System (`domain/systems/physics.go`)

Orchestrates pure physics functions:

```go
type PhysicsSystem struct {
    world *world.World
}

func (ps *PhysicsSystem) UpdatePhysics(
    entity entities.PhysicsEntity,
    inputState input.InputState,
    dt float32,
) {
    // Orchestrate pure functions (all testable, no Raylib)
    velocity := physics.ApplyHorizontalMovement(velocity, inputState, dt)
    velocity = physics.ApplyVerticalMovement(velocity, inputState, dt)
    velocity = physics.ApplyGravity(velocity, dt)
    position = physics.IntegrateVelocity(position, velocity, dt)

    result := physics.ResolveGroundCollision(position, velocity, entity.GetHeight(), ps.world)

    entity.SetPosition(result.Position)
    entity.SetVelocity(result.Velocity)
    entity.SetOnGround(result.OnGround)
}
```

**Why this design:**
- Uses explicit `PhysicsEntity` interface (not anonymous inline interfaces)
- Pure physics functions fully testable without framework
- Accepts `InputState` (not Raylib types)
- IDEs can resolve all symbols correctly

#### Pure Physics Functions (`domain/physics/`)

Framework-independent mathematical functions:

```go
// movement.go - Pure functions, no Raylib, fully testable
func ApplyHorizontalMovement(velocity Vec2, inputState InputState, dt float32) Vec2
func ApplyVerticalMovement(velocity Vec2, inputState InputState, dt float32) Vec2

// gravity.go - Pure functions
func ApplyGravity(velocity Vec2, dt float32) Vec2
func IntegrateVelocity(position, velocity Vec2, dt float32) Vec2

// collision.go - Pure functions
func ResolveGroundCollision(position, velocity Vec2, height float32, world *World) CollisionResult
```

**Why this design:**
- Zero Raylib imports
- Can be tested standalone
- Input/output are domain types (Vec2, InputState, etc.)
- Pure functions enable unit testing without framework

#### Player Entity (`domain/entities/player.go`)

Pure data entity with domain behavior:

```go
type Player struct {
    Position types.Vec2  // Domain type (not rl.Vector2)
    Velocity types.Vec2
    OnGround bool
}

// Domain interface for physics systems
type PhysicsEntity interface {
    GetPositionVec() *types.Vec2
    GetVelocityVec() *types.Vec2
    SetPosition(types.Vec2)
    SetVelocity(types.Vec2)
    SetOnGround(bool)
    GetHeight() float32
}
```

**Why this design:**
- No Render() method (rendering is adapter responsibility)
- Uses domain types (Vec2, not rl.Vector2)
- Zero Raylib dependency
- Clear contract via `PhysicsEntity` interface for extensibility
- Easy to add new entities (Enemy, NPC) that implement the same interface

#### Types (`domain/types/vec2.go`)

Custom math types independent of framework:

```go
type Vec2 struct {
    X float32
    Y float32
}

func (v Vec2) Add(other Vec2) Vec2
func (v Vec2) Scale(scalar float32) Vec2
func (v Vec2) Magnitude() float32
// ... other operations
```

**Why this design:**
- No Raylib dependency
- Physics can use its own types
- Conversion to Raylib types only happens in rendering adapter
- Provides abstraction layer

#### InputState (`domain/input/input_state.go`)

Platform-agnostic input representation:

```go
type InputState struct {
    Left  bool
    Right bool
    Up    bool
    Down  bool
}
```

**Why this design:**
- Not a Raylib type
- Physics and game logic receive this, not raw keyboard input
- Easy to swap input sources (file playback, network, AI)
- Domain logic decoupled from input mechanism

#### World (`domain/world/world.go`)

World data structure:

```go
type World struct {
    GroundLevel float32
    Width       float32
    Height      float32
}

func (w *World) GetGroundLevel() float32
func (w *World) IsInBounds(x, y float32) bool
```

**Why this design:**
- Centralizes world parameters
- No rendering data (colors, textures)
- Extensible for future terrain, tiles, etc.
- Used by physics system for collision checks

---

### Adapter Layer (Framework Integration)

#### Input Adapter (`internal/adapters/input/raylib.go`)

Translates Raylib input to domain InputState:

```go
type RaylibInputAdapter struct{}

func (a *RaylibInputAdapter) ReadInput() input.InputState {
    return input.InputState{
        Left:  rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA),
        Right: rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD),
        Up:    rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW),
        Down:  rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS),
    }
}
```

**Why this design:**
- Single responsibility: Read Raylib keys, output platform-agnostic InputState
- All Raylib input code in one place
- Easy to add new input sources (just create a new adapter)
- Can be mocked for testing

#### Rendering Adapter (`internal/adapters/rendering/raylib.go`)

Takes domain entities and renders them with Raylib:

```go
type RaylibRenderer struct{}

func (r *RaylibRenderer) Render(game *engine.Game) {
    rl.BeginDrawing()
    rl.ClearBackground(rl.RayWhite)

    r.renderWorld(game.GetWorld())
    r.renderPlayer(game.GetPlayer())

    rl.EndDrawing()
}

// Window lifecycle
func (r *RaylibRenderer) InitWindow(width, height int32, title string)
func (r *RaylibRenderer) CloseWindow()
func (r *RaylibRenderer) WindowShouldClose() bool
func (r *RaylibRenderer) SetTargetFPS(fps int32)
func (r *RaylibRenderer) GetFrameTime() float32
```

**Why this design:**
- Complete abstraction of Raylib
- Takes domain entities (Game, Player, World) and renders them
- All window lifecycle management in one place
- No business logic (pure rendering)
- Easy to swap for different renderer (SDL, canvas, headless)

---

### Application Layer (Orchestration)

#### Main (`cmd/game/main.go`)

Orchestrates the game loop:

```go
func main() {
    // Create adapters
    renderer := rendering.NewRaylibRenderer()
    inputAdapter := input.NewRaylibInputAdapter()

    // Initialize Raylib (adapter responsibility)
    renderer.InitWindow(screenWidth, screenHeight, "Drill Game")
    defer renderer.CloseWindow()
    renderer.SetTargetFPS(targetFPS)

    // Create domain objects
    gameWorld := world.NewWorld(screenWidth, screenHeight, groundLevel)
    game := engine.NewGame(gameWorld)

    // Main loop
    for renderer.WindowShouldClose() == false {
        dt := renderer.GetFrameTime()

        // 1. Read input from adapter
        inputState := inputAdapter.ReadInput()

        // 2. Update domain
        err := game.Update(dt, inputState)
        if err != nil {
            slog.Error("Error during update", "error", err)
            break
        }

        // 3. Render via adapter
        renderer.Render(game)
    }
}
```

**Why this design:**
- Clear, linear flow: Input → Update → Render
- No game logic (orchestration only)
- Dependencies are explicit
- Easy to understand at a glance

---

## Design Principles

### 1. Separation of Concerns

Each layer has a single responsibility:

- **Domain**: Business logic (physics, game rules)
- **Adapters**: Framework integration (Raylib)
- **Application**: Orchestration (wiring it together)

### 2. Framework Independence

The domain layer has **zero framework dependencies**:

- No `import rl "github.com/gen2brain/raylib-go/raylib"`
- No Raylib types (rl.Vector2, rl.Color, etc.)
- All conversions happen in adapters

### 3. Testability

Core logic is fully testable without Raylib:

```go
// This runs WITHOUT rl.InitWindow()
go test ./internal/domain/physics/...

// 11 physics tests pass with zero Raylib dependency
TestApplyHorizontalMovement_Acceleration
TestApplyHorizontalMovement_MaxSpeed
TestApplyGravity_IncreasesDownwardVelocity
// ... etc
```

### 4. Value Types for Small Objects

Small types are values, not pointers:

```go
type Vec2 struct {
    X float32  // 8 bytes total
    Y float32
}

// ✓ Good: Values for small types
player.Position = types.Vec2{X: 100, Y: 200}

// ✗ Bad: Pointers for small types
player.Position = &types.Vec2{X: 100, Y: 200}
```

**Why:** Small types (8-16 bytes) should be values:
- Faster on stack than heap allocation
- Better cache locality
- No nil pointer issues
- Go idiom (see time.Time, image.Point)

### 5. Explicit Contracts

Use explicit interfaces, not anonymous inline interfaces:

```go
// ✓ Good: Named interface in domain
type PhysicsEntity interface {
    GetPositionVec() *types.Vec2
    GetVelocityVec() *types.Vec2
    SetPosition(types.Vec2)
    SetVelocity(types.Vec2)
    SetOnGround(bool)
    GetHeight() float32
}

// ✗ Bad: Anonymous inline interface (confuses IDEs)
func UpdatePhysics(
    player interface {
        GetPositionVec() *types.Vec2
        // ... etc
    },
) {
}
```

**Why:**
- IDEs can resolve symbols
- "Find all implementations" works
- Refactoring is safe
- Documents the contract clearly

### 6. Getters for External Access

Adapters access domain data through getters:

```go
// Domain provides getters
func (g *Game) GetWorld() *world.World
func (g *Game) GetPlayer() *entities.Player

// Adapter reads via getters (doesn't modify)
func (r *RaylibRenderer) Render(game *engine.Game) {
    world := game.GetWorld()
    player := game.GetPlayer()
    // ... render
}
```

**Why:**
- Encapsulation (adapter can't modify game state)
- Clear data flow (one-way: domain → adapter)
- Easy to add state management later

---

## Future Architecture Considerations

### Adding New Entities

To add an Enemy that also uses physics:

```go
// 1. Implement PhysicsEntity interface
type Enemy struct {
    Position types.Vec2
    Velocity types.Vec2
    Health   float32
    // ...
}

func (e *Enemy) GetPositionVec() *types.Vec2 { return &e.Position }
func (e *Enemy) GetVelocityVec() *types.Vec2 { return &e.Velocity }
// ... implement interface

// 2. Use in physics system (zero changes needed!)
func (ps *PhysicsSystem) UpdatePhysics(
    entity entities.PhysicsEntity,  // Works with Enemy too
    inputState input.InputState,
    dt float32,
) {
    // Same logic works for all PhysicsEntity implementations
}
```

### Swapping Renderers

To use SDL instead of Raylib:

```go
// Create SDL adapter with same interface
type SDLRenderer struct {
    window *sdl.Window
    // ...
}

func (r *SDLRenderer) Render(game *engine.Game) {
    // SDL drawing logic
}

// Swap in main.go
renderer := sdl.NewSDLRenderer()  // Instead of Raylib
// Rest of main.go works unchanged!
```

### Adding Input Sources

To add file-based replay input:

```go
// New adapter with same interface
type FileInputAdapter struct {
    frames []InputState
    index  int
}

func (a *FileInputAdapter) ReadInput() input.InputState {
    state := a.frames[a.index]
    a.index++
    return state
}

// Swap in main.go
inputAdapter := file.NewFileInputAdapter("replay.bin")
// Game loop works unchanged!
```

---

## Testing Strategy

### Unit Tests (Domain Logic)

Test pure functions without framework:

```bash
# Physics tests (11 tests, no Raylib required)
go test ./internal/domain/physics/...

# Tests cover:
# - Movement (acceleration, damping, max speed)
# - Gravity (falling, velocity integration)
# - Collision (ground detection, resolution)
```

### Integration Tests (Systems)

Test systems working together:

```go
// Example: Test player movement with physics
world := domain.NewWorld(1280, 720, 600)
player := domain.NewPlayer(640, 500)
physics := domain.NewPhysicsSystem(world)

inputState := domain.InputState{Right: true}
physics.UpdatePhysics(player, inputState, 0.016)

// Assert player moved right
assert.True(player.GetVelocityVec().X > 0)
```

### Manual Testing

Rendering and feel testing requires manual play:

- Gameplay feel (movement responsiveness)
- Visual polish (animations, particles)
- Performance (frame rates, memory)

---

## Performance Considerations

### Current (Phase 1)

- Simple physics (position + velocity)
- No spatial partitioning yet
- Direct collision checks
- Frame-independent movement via delta time

### Future Optimizations

- **Spatial Partitioning**: Grid or quadtree for collision queries
- **Object Pooling**: Reuse frequently created objects
- **Batch Rendering**: Group draw calls
- **Chunk Loading**: Only simulate/render visible area

---

## Dependencies

Minimal, intentional dependencies:

- `github.com/gen2brain/raylib-go/raylib` - Graphics/audio (adapters only)
- `github.com/stretchr/testify` - Testing utilities (optional, use as needed)
- Go standard library - Everything else

**Strict rule:** No domain code imports Raylib.

---

## Summary

This architecture achieves:

- ✅ **Testability**: Domain logic 100% testable without framework
- ✅ **Portability**: Could swap Raylib for any renderer
- ✅ **Maintainability**: Clear responsibilities, linear data flow
- ✅ **Extensibility**: Easy to add new entities, systems, input sources
- ✅ **Clarity**: Folder structure = architecture diagram

The core principle: **Domain stays pure, framework stays outside.**
