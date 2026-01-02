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
│       │   ├── player.go                    # Player entity (AABB-based)
│       │   ├── tile.go                      # Tile entity (Empty, Dirt, Ore)
│       │   └── ore_type.go                  # Ore types & Gaussian parameters
│       ├── physics/
│       │   ├── constants.go                 # Physics parameters
│       │   ├── movement.go                  # Movement functions
│       │   ├── gravity.go                   # Gravity + velocity integration
│       │   ├── collision.go                 # AABB collision detection/resolution
│       │   ├── movement_test.go             # Movement tests
│       │   ├── gravity_test.go              # Gravity tests
│       │   └── collision_test.go            # AABB collision tests
│       ├── types/
│       │   ├── vec2.go                      # Custom Vec2 (no Raylib types)
│       │   ├── aabb.go                      # AABB collision primitive
│       │   └── aabb_test.go                 # AABB unit tests
│       ├── input/
│       │   └── input_state.go               # InputState struct (framework-agnostic)
│       └── world/
│           ├── world.go                     # World: chunk loading, sparse tile map
│           ├── generator.go                 # Procedural tile generation
│           ├── hash.go                      # Deterministic seeding (FNV-1a)
│           ├── generator_test.go            # Generator unit tests
│           ├── world_test.go                # Chunk loading tests
│           └── integration_test.go          # End-to-end world generation tests
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
│    • Load chunks around player (3×3)    │
│    • Downward digging (S/Down key)      │
│    • Horizontal digging (L/R when grounded)
│    • Physics system applies forces      │
│    • Gravity, movement, collision       │
│    • Updates Player position/velocity   │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│ 3. Render via Adapter                   │
│    renderer.Render(game)                │
│    • Extracts Player, World from game   │
│    • Renders tiles with ore colors      │
│    • Draws player, entities             │
│    • Camera follows player              │
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
    world          *world.World
    player         *entities.Player
    physicsSystem  *systems.PhysicsSystem
    diggingSystem  *systems.DiggingSystem
}

func (g *Game) Update(dt float32, inputState input.InputState) error {
    // Pure domain logic - zero Raylib
    // Update chunks around player
    playerX := g.player.AABB.X + g.player.AABB.Width/2
    playerY := g.player.AABB.Y + g.player.AABB.Height/2
    g.world.UpdateChunksAroundPlayer(playerX, playerY)

    // Downward digging (before physics, so grid alignment works first)
    g.diggingSystem.ProcessDigging(g.player, inputState)

    // Horizontal digging (auto-dig left/right when grounded against walls)
    g.diggingSystem.ProcessHorizontalDigging(g.player, inputState)

    // Physics with axis-separated collision
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

#### Digging System (`domain/systems/digging.go`)

Handles tile destruction for both downward and horizontal digging:

```go
type DiggingSystem struct {
    world *world.World
}

// Downward digging with grid alignment (S/Down key)
func (ds *DiggingSystem) ProcessDigging(
    player *entities.Player,
    inputState input.InputState,
) {
    if !inputState.Dig {
        return
    }

    playerCenterX := player.AABB.X + player.AABB.Width/2
    playerBottomY := player.AABB.Y + player.AABB.Height

    tile := ds.world.GetTileAt(playerCenterX, playerBottomY)
    if tile != nil && tile.IsDiggable() {
        // Snap to grid and dig
        tileGridX := int(playerCenterX / world.TileSize)
        tileCenterX := float32(tileGridX)*world.TileSize + world.TileSize/2
        player.AABB.X = tileCenterX - player.AABB.Width/2

        ds.world.DigTile(playerCenterX, playerBottomY)
    }
}

// Horizontal digging when grounded (A/D or Left/Right keys)
func (ds *DiggingSystem) ProcessHorizontalDigging(
    player *entities.Player,
    inputState input.InputState,
) {
    // Only dig when grounded
    if !player.OnGround {
        return
    }

    playerCenterY := player.AABB.Y + player.AABB.Height/2

    if inputState.Left {
        // Check tile just left of player
        tile := ds.world.GetTileAt(player.AABB.X-1, playerCenterY)
        if tile != nil && tile.IsDiggable() {
            ds.world.DigTile(player.AABB.X-1, playerCenterY)
            return
        }
    }

    if inputState.Right {
        // Check tile just right of player
        tile := ds.world.GetTileAt(player.AABB.X+player.AABB.Width+1, playerCenterY)
        if tile != nil && tile.IsDiggable() {
            ds.world.DigTile(player.AABB.X+player.AABB.Width+1, playerCenterY)
        }
    }
}
```

**Why this design:**
- Called before physics (so blocked tiles can be removed before collision resolution)
- Uses direct world coordinate sampling (like downward digging)
- Horizontal digging only when `player.OnGround == true` (prevents mid-air digging)
- No collision detection needed (simple world queries)
- Direct field access for simplicity

#### Physics System (`domain/systems/physics.go`)

Orchestrates pure physics functions with axis-separated collision:

```go
type PhysicsSystem struct {
    world *world.World
}

func (ps *PhysicsSystem) UpdatePhysics(
    player *entities.Player,
    inputState input.InputState,
    dt float32,
) {
    // 1. Apply movement and gravity to velocity
    player.Velocity = physics.ApplyHorizontalMovement(player.Velocity, inputState, dt)
    player.Velocity = physics.ApplyVerticalMovement(player.Velocity, inputState, dt)
    player.Velocity = physics.ApplyGravity(player.Velocity, dt)

    // 2. AXIS-SEPARATED COLLISION RESOLUTION

    // X-axis: integrate position → check → resolve
    player.AABB.X += player.Velocity.X * dt
    collisionsX := physics.CheckCollisions(player.AABB, ps.world)
    player.AABB, player.Velocity = physics.ResolveCollisionsX(player.AABB, player.Velocity, collisionsX)

    // Y-axis: integrate position → check → resolve
    player.AABB.Y += player.Velocity.Y * dt
    collisionsY := physics.CheckCollisions(player.AABB, ps.world)
    player.AABB, player.Velocity, player.OnGround = physics.ResolveCollisionsY(player.AABB, player.Velocity, collisionsY)
}
```

**Why this design:**
- Direct field access (no getters/setters) for simplicity
- Axis-separated collision prevents corner-catching
- Pure physics functions fully testable without framework
- Accepts `InputState` (not Raylib types)
- Works with `*Player` directly (no interface needed)

#### Pure Physics Functions (`domain/physics/`)

Framework-independent mathematical functions:

```go
// movement.go - Pure functions, no Raylib, fully testable
func ApplyHorizontalMovement(velocity Vec2, inputState InputState, dt float32) Vec2
func ApplyVerticalMovement(velocity Vec2, inputState InputState, dt float32) Vec2

// gravity.go - Pure functions
func ApplyGravity(velocity Vec2, dt float32) Vec2

// collision.go - AABB-based collision functions
func CheckCollisions(aabb AABB, world *World) []TileCollision
func ResolveCollisionsX(aabb AABB, velocity Vec2, collisions []TileCollision) (AABB, Vec2)
func ResolveCollisionsY(aabb AABB, velocity Vec2, collisions []TileCollision) (AABB, Vec2, bool)
func GetOccupiedTileRange(aabb AABB, tileSize float32) (minX, maxX, minY, maxY int)
```

**Why this design:**
- Zero Raylib imports
- Can be tested standalone
- Input/output are domain types (AABB, Vec2, etc.)
- Pure functions enable unit testing without framework
- Value-based (no pointer mutations in function signatures)

#### Player Entity (`domain/entities/player.go`)

Pure data entity with AABB collision primitive:

```go
type Player struct {
    AABB     types.AABB  // Position and dimensions - direct access
    Velocity types.Vec2  // Pixels per second - direct access
    OnGround bool        // Collision state - direct access
}

func NewPlayer(startX, startY float32) *Player {
    return &Player{
        AABB:     types.NewAABB(startX, startY, PlayerWidth, PlayerHeight),
        Velocity: types.Zero(),
        OnGround: false,
    }
}
```

**Why this design:**
- AABB eliminates redundant Position storage (X, Y already in AABB)
- No Render() method (rendering is adapter responsibility)
- Uses domain types (AABB, Vec2, not rl.Vector2)
- Zero Raylib dependency
- Direct field access (no getters/setters) for simplicity
- AABB enables proper collision detection (not just ground)

#### Types (`domain/types/`)

Custom math types independent of framework:

**Vec2** (`vec2.go`):
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

**AABB** (`aabb.go`):
```go
type AABB struct {
    X, Y          float32 // Top-left corner position
    Width, Height float32 // Dimensions
}

func (a AABB) Intersects(b AABB) bool
func (a AABB) Penetration(b AABB) (dx, dy float32)
func (a AABB) Min() Vec2
func (a AABB) Max() Vec2
```

**Why this design:**
- No Raylib dependency
- Physics can use its own types
- Conversion to Raylib types only happens in rendering adapter
- AABB provides proper collision detection (not just point-based)
- Value types (not pointers) for simplicity and performance

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

type AABB struct {
    X, Y          float32  // 16 bytes total
    Width, Height float32
}

// ✓ Good: Values for small types
player.Velocity = types.Vec2{X: 100, Y: 200}
player.AABB = types.NewAABB(0, 0, 64, 64)

// ✗ Bad: Pointers for small types
player.Velocity = &types.Vec2{X: 100, Y: 200}
```

**Why:** Small types (8-16 bytes) should be values:
- Faster on stack than heap allocation
- Better cache locality
- No nil pointer issues
- Go idiom (see time.Time, image.Point)
- Cheaper to copy than pointer indirection on modern CPUs

### 5. Direct Field Access for Simplicity

Use direct field access instead of getters/setters when appropriate:

```go
// ✓ Good: Direct field access
player.AABB.X += player.Velocity.X * dt
player.Velocity.Y += gravity * dt
player.OnGround = true

// ✗ Overly complex: Unnecessary indirection
player.SetPosition(player.GetPosition().Add(player.GetVelocity().Scale(dt)))
player.SetVelocity(player.GetVelocity().Add(Vec2{Y: gravity * dt}))
player.SetOnGround(true)
```

**Why:**
- Simpler code, easier to read
- Less boilerplate (no getter/setter methods)
- Better performance (no function call overhead)
- Go idiom: exported fields for simple data structures
- Still maintains encapsulation at package boundaries

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

## AABB Collision System

The game uses **Axis-Aligned Bounding Box (AABB) collision detection** with axis-separated resolution for precise 2D platformer physics.

### Core Concepts

**AABB Primitive:**
- Rectangular collision box defined by position (X, Y) and dimensions (Width, Height)
- Axis-aligned (no rotation) for fast intersection tests
- Used for both player and tiles

**Axis-Separated Resolution:**
- X-axis movement and collision resolved first
- Y-axis movement and collision resolved second
- Prevents corner-catching and enables natural wall sliding

### Collision Pipeline

```go
// 1. Apply movement and gravity
player.Velocity = ApplyHorizontalMovement(player.Velocity, input, dt)
player.Velocity = ApplyGravity(player.Velocity, dt)

// 2. X-axis: integrate → detect → resolve
player.AABB.X += player.Velocity.X * dt
collisionsX := CheckCollisions(player.AABB, world)
player.AABB, player.Velocity = ResolveCollisionsX(player.AABB, player.Velocity, collisionsX)

// 3. Y-axis: integrate → detect → resolve
player.AABB.Y += player.Velocity.Y * dt
collisionsY := CheckCollisions(player.AABB, world)
player.AABB, player.Velocity, player.OnGround = ResolveCollisionsY(player.AABB, player.Velocity, collisionsY)
```

### Collision Detection

**CheckCollisions()** finds all solid tiles overlapping the player:

```go
func CheckCollisions(aabb AABB, world *World) []TileCollision {
    // 1. Calculate which tiles the AABB might overlap
    minX, maxX, minY, maxY := GetOccupiedTileRange(aabb, TileSize)

    // 2. Check each potentially overlapping tile
    for x := minX; x <= maxX; x++ {
        for y := minY; y <= maxY; y++ {
            tile := world.GetTileAtGrid(x, y)
            if tile != nil && tile.IsSolid() && aabb.Intersects(tile.GetAABB(x, y, TileSize)) {
                // Found collision!
            }
        }
    }
}
```

**Performance:** Player can overlap at most 4 tiles (2×2 grid), so maximum 4 intersection tests per frame.

### Collision Resolution

**ResolveCollisionsX()** pushes player out horizontally:
- Calculates penetration depth using `AABB.Penetration()`
- Adjusts position: `aabb.X -= dx`
- Zeros horizontal velocity on wall hit

**ResolveCollisionsY()** pushes player out vertically:
- Calculates penetration depth
- Adjusts position: `aabb.Y -= dy`
- Detects ground: if pushed up (`dy > 0`), set `OnGround = true`
- Detects ceiling: if pushed down (`dy < 0`), zero upward velocity

### Why Axis-Separated?

**Without axis separation (naive AABB):**
- Player moving diagonally into corner gets "stuck"
- Cannot slide along walls smoothly
- Ground detection is ambiguous

**With axis separation:**
- X collision resolved first, Y collision resolved second
- Player slides along walls naturally during diagonal movement
- Clear ground/ceiling/wall detection based on which axis had collision

### Penetration Calculation

```go
func (a AABB) Penetration(b AABB) (dx, dy float32) {
    // Calculate overlap on each axis
    overlapX := min(a.X+a.Width, b.X+b.Width) - max(a.X, b.X)
    overlapY := min(a.Y+a.Height, b.Y+b.Height) - max(a.Y, b.Y)

    // Determine push direction based on relative positions
    if a.X < b.X {
        dx = overlapX  // Push left (subtract to move right)
    } else {
        dx = -overlapX // Push right (subtract to move left)
    }

    // Same for Y axis
    // ...
}
```

**Key insight:** Signs are chosen so `position -= penetration` always pushes objects apart.

---

## Camera System

The game uses Raylib's `Camera2D` for viewport management, allowing the player to explore a world much larger than the screen.

### Camera Implementation

**Camera2D** lives in the rendering adapter (not domain):

```go
type RaylibRenderer struct {
    camera       rl.Camera2D
    screenWidth  float32
    screenHeight float32
    worldWidth   float32
}

func (r *RaylibRenderer) updateCamera(player *entities.Player, w *world.World) {
    // Camera target follows player center
    playerCenterX := player.AABB.X + player.AABB.Width/2
    playerCenterY := player.AABB.Y + player.AABB.Height/2

    // Clamp camera to world bounds
    halfScreenW := r.screenWidth / 2
    halfScreenH := r.screenHeight / 2

    minX := halfScreenW
    maxX := r.worldWidth - halfScreenW
    minY := w.GetGroundLevel() - halfScreenH

    // Clamp and assign to camera target
    targetX := clamp(playerCenterX, minX, maxX)
    targetY := clamp(playerCenterY, minY, maxY)
    r.camera.Target = rl.Vector2{X: targetX, Y: targetY}
}

func (r *RaylibRenderer) Render(game *engine.Game, inputState input.InputState) {
    r.updateCamera(game.GetPlayer(), game.GetWorld())

    rl.BeginDrawing()
    rl.ClearBackground(rl.RayWhite)

    // World space rendering (camera applied)
    rl.BeginMode2D(r.camera)
    r.renderWorld(game.GetWorld())
    r.renderTiles(game.GetWorld())
    r.renderPlayer(game.GetPlayer())
    rl.EndMode2D()

    // Screen space rendering (no camera, UI always visible)
    r.renderDebugInfo(game.GetPlayer(), inputState)

    rl.EndDrawing()
}
```

### Viewport Culling

For performance with large worlds, only tiles visible in the camera viewport are rendered:

```go
func (r *RaylibRenderer) renderTiles(w *world.World) {
    tiles := w.GetAllTiles()

    // Calculate visible tile range
    minVisibleX := int((r.camera.Target.X - r.screenWidth/2) / world.TileSize) - 1
    maxVisibleX := int((r.camera.Target.X + r.screenWidth/2) / world.TileSize) + 1
    minVisibleY := int((r.camera.Target.Y - r.screenHeight/2) / world.TileSize) - 1
    maxVisibleY := int((r.camera.Target.Y + r.screenHeight/2) / world.TileSize) + 1

    for coord, tile := range tiles {
        gridX, gridY := coord[0], coord[1]

        // Skip tiles outside viewport
        if gridX < minVisibleX || gridX > maxVisibleX ||
           gridY < minVisibleY || gridY > maxVisibleY {
            continue
        }

        // Render visible tile...
    }
}
```

**Performance:** Reduces tiles rendered from ~94,000 to ~300 (~300× improvement).

### Why in Adapter, Not Domain?

- Camera is a **rendering concern**, not game logic
- Tightly coupled to Raylib's `Camera2D` struct
- Player position and world bounds already in domain (no new logic)
- Follows pattern: adapters translate domain state to visual representation

---

## World Boundary Constraints

Players cannot leave the game area. World boundaries are enforced by the physics system:

```go
func (ps *PhysicsSystem) constrainPlayerToWorldBounds(player *entities.Player) {
    // Horizontal: player.X must be in [0, worldWidth - playerWidth]
    minX := float32(0.0)
    maxX := ps.world.Width - float32(entities.PlayerWidth)

    if player.AABB.X < minX {
        player.AABB.X = minX
        player.Velocity.X = 0
    } else if player.AABB.X > maxX {
        player.AABB.X = maxX
        player.Velocity.X = 0
    }

    // Vertical: player.Y must be >= 0
    minY := float32(0.0)

    if player.AABB.Y < minY {
        player.AABB.Y = minY
        player.Velocity.Y = 0
    }
    // No maximum Y - player can dig infinitely deep
}
```

Called after collision resolution in the physics pipeline:

```go
// 2. Axis-separated collision resolution
player.AABB.X += player.Velocity.X * dt
collisionsX := physics.CheckCollisions(player.AABB, ps.world)
player.AABB, player.Velocity = physics.ResolveCollisionsX(player.AABB, player.Velocity, collisionsX)

player.AABB.Y += player.Velocity.Y * dt
collisionsY := physics.CheckCollisions(player.AABB, ps.world)
player.AABB, player.Velocity, player.OnGround = physics.ResolveCollisionsY(player.AABB, player.Velocity, collisionsY)

// 3. Enforce world boundary constraints
ps.constrainPlayerToWorldBounds(player)
```

**Design note:** Boundary constraints are purely domain-level (physics system), not rendering. Camera clamping happens independently in the adapter, preventing the camera from showing off-screen areas.

---

## World Dimensions

The game world extends far beyond the screen:

| Dimension | Size | Tiles |
|-----------|------|-------|
| Width | 7680 pixels | 120 tiles wide (6× screen width) |
| Height | 64000 pixels | 1000 tiles deep |
| Ground Level | 640 pixels | 10 tiles up from bottom |
| Tile Size | 64×64 pixels | Standard |

**Sparse tile storage:** Only non-empty tiles are stored in memory, enabling efficient large worlds.

---

## Future Architecture Considerations

### Adding New Entities

To add an Enemy that also uses physics:

```go
// 1. Create entity with AABB
type Enemy struct {
    AABB     types.AABB
    Velocity types.Vec2
    Health   float32
    AI       AIState
    // ...
}

// 2. Create a separate UpdateEnemyPhysics method or generalize UpdatePhysics
func (ps *PhysicsSystem) UpdateEnemyPhysics(enemy *entities.Enemy, dt float32) {
    // Same collision logic as player
    enemy.Velocity = physics.ApplyGravity(enemy.Velocity, dt)

    enemy.AABB.X += enemy.Velocity.X * dt
    collisionsX := physics.CheckCollisions(enemy.AABB, ps.world)
    enemy.AABB, enemy.Velocity = physics.ResolveCollisionsX(enemy.AABB, enemy.Velocity, collisionsX)

    enemy.AABB.Y += enemy.Velocity.Y * dt
    collisionsY := physics.CheckCollisions(enemy.AABB, ps.world)
    enemy.AABB, enemy.Velocity, _ = physics.ResolveCollisionsY(enemy.AABB, enemy.Velocity, collisionsY)
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
# - Collision (AABB detection, axis-separated resolution, wall/ceiling/ground)
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
