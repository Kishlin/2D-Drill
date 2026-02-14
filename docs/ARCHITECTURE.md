# Architecture

## Overview

Drill Game uses **Hexagonal Architecture** (Ports & Adapters) to achieve a clean separation between pure domain logic and framework-specific integration.

**Benefits:**
- **Testable**: Physics and game logic can be tested without initializing Raylib
- **Portable**: Domain code has zero framework dependencies
- **Maintainable**: Clear responsibilities and data flow
- **Extensible**: Easy to add new features, entities, or rendering backends

**Related Documentation:**
- [SYSTEMS.md](SYSTEMS.md) — Game systems (drilling, physics, fuel, effects, UI, items)
- [PHYSICS.md](PHYSICS.md) — Collision detection, damage formulas, camera
- [WORLD.md](WORLD.md) — World generation, tiles, ore/hazard distributions
- [BOSS.md](BOSS.md) — Boss system, state machines, phases
- [CONFIGURATION.md](CONFIGURATION.md) — Config structs, levels, reference tables
- [GAME_DESIGN.md](GAME_DESIGN.md) — Game mechanics and progression
- [DEVELOPMENT.md](DEVELOPMENT.md) — Development workflows and testing

---

## The Three Layers

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
                    │ • World (world/)
                    │ • Effects (effects/)
                    │ • UI (ui/)
                    └──────────────────┘
```

### 1. Application Layer (`cmd/game/main.go`)

Orchestrates the entire game:
- Creates adapters and game instance
- Runs the main loop: Input → Update → Render
- Manages window lifecycle

```go
func main() {
    renderer := rendering.NewRaylibRenderer()
    inputAdapter := input.NewRaylibInputAdapter()

    renderer.InitWindow(screenWidth, screenHeight, "Drill Game")
    defer renderer.CloseWindow()

    game, err := engine.NewGame(gameCfg)
    if err != nil {
        slog.Error("Failed to create game", "error", err)
        return
    }

    for renderer.WindowShouldClose() == false {
        dt := renderer.GetFrameTime()
        inputState := inputAdapter.ReadInput()
        game.Update(dt, inputState)
        renderer.Render(game)
    }
}
```

### 2. Adapter Layer (`internal/adapters/`)

Framework integration with zero business logic:

**Input Adapter** — Translates Raylib input to domain `InputState`:
```go
func (a *RaylibInputAdapter) ReadInput() input.InputState {
    return input.InputState{
        Left:  rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA),
        Right: rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD),
        // ...
    }
}
```

**Rendering Adapter** — Takes domain entities and renders with Raylib:
```go
func (r *RaylibRenderer) Render(game *engine.Game) {
    rl.BeginDrawing()
    r.renderWorld(game.World)
    r.renderPlayer(game.Player)
    rl.EndDrawing()
}
```

### 3. Domain Layer (`internal/domain/`)

Pure game logic with zero framework dependencies:

| Package | Purpose |
|---------|---------|
| `engine/` | Game orchestration |
| `entities/` | Player, Building, Tile, ItemCatalog |
| `upgrades/` | Upgrade types, Catalog, UpgradeType enum |
| `systems/` | Physics, Drilling, Fuel, Items, Boss fights |
| `effects/` | State mutations via EffectContext |
| `ui/` | Shop interfaces |
| `world/` | Chunk-based procedural world |
| `physics/` | Movement, collision, damage |
| `types/` | Vec2, AABB |
| `input/` | InputState struct |
| `config/` | Configuration structs |
| `levels/` | Level definitions |
| `bosses/` | Boss infrastructure (interfaces, base types, state machine) |
| `boss_catalog/` | Boss implementations (one package per boss) |
| `projectiles/` | Projectile spawning and movement types |
| `components/` | Position, Interactable, Damageable |

---

## Project Structure

```
drill-game/
├── cmd/game/main.go              # Application orchestration
│
├── internal/
│   ├── adapters/                 # Framework Integration (Raylib)
│   │   ├── input/raylib.go       # RaylibInputAdapter
│   │   └── rendering/
│   │       ├── raylib.go         # RaylibRenderer
│   │       └── bosses/           # Boss-specific renderers
│   │
│   └── domain/                   # Pure Business Logic (NO RAYLIB)
│       ├── engine/game.go        # Game loop orchestration
│       ├── entities/             # Player, Building, Tile, ItemCatalog
│       ├── upgrades/             # Upgrade types, Catalog, UpgradeType
│       ├── systems/              # Physics, Drilling, Fuel, Items, Boss
│       ├── effects/              # Effect interface and implementations
│       ├── ui/                   # UI interface, Manager, shops
│       ├── world/                # World, chunk generator
│       ├── physics/              # Movement, collision, damage
│       ├── types/                # Vec2, AABB
│       ├── input/                # InputState
│       ├── config/               # Configuration structs
│       ├── levels/               # Level definitions
│       ├── bosses/               # Boss infrastructure
│       ├── boss_catalog/         # Boss implementations
│       └── components/           # Position, Interactable, Damageable
│
└── docs/                         # Documentation
```

---

## Data Flow

### Single Frame Update

```
1. Read Input from Adapter
   inputState := inputAdapter.ReadInput()
   (Converts Raylib keys → InputState)
        │
        ▼
2. Update Domain Logic
   game.Update(dt, inputState)
   • Load chunks around player
   • Process active UI (modal pause)
   • Detect building interactions
   • Apply physics & damage
   • Consume fuel
   • Process drilling animation
   • Handle item usage
   • Update boss fight
        │
        ▼
3. Render via Adapter
   renderer.Render(game)
   • Extract Player, World, Buildings
   • Render visible tiles
   • Draw entities
   • Render UI overlays
   • Camera follows player
```

---

## Design Principles

### 1. Framework Independence

The domain layer has **zero framework dependencies**:
- No `import rl "github.com/gen2brain/raylib-go/raylib"`
- No Raylib types (rl.Vector2, rl.Color, etc.)
- All conversions happen in adapters

**Verify:** `grep -r "raylib" internal/domain/`

### 2. Testability

Core logic is fully testable without Raylib:

```bash
go test ./internal/domain/...  # Runs without rl.InitWindow()
```

### 3. Value Types for Small Objects

Small types are values, not pointers:

```go
type Vec2 struct { X, Y float32 }      // 8 bytes
type AABB struct { X, Y, W, H float32 } // 16 bytes

// Good: Values for small types
player.Velocity = types.Vec2{X: 100, Y: 200}

// Bad: Pointers for small types
player.Velocity = &types.Vec2{X: 100, Y: 200}
```

### 4. Direct Field Access

Use direct field access for simplicity:

```go
// Good: Direct access
player.AABB.X += player.Velocity.X * dt
player.OnGround = true

// Overly complex
player.SetPosition(player.GetPosition().Add(...))
```

### 5. Public Fields for External Access

Externally-accessed Game fields are public for direct access by adapters. Internal-only fields remain private:

```go
// Public fields (used by renderer/main)
game.World       // *world.World
game.Player      // *entities.Player
game.Buildings   // []*entities.Building
game.Boss        // bosses.Boss
game.GameState   // entities.GameState
game.UIManager   // *ui.Manager
game.InventoryUI // *ui.InventoryUI
game.Projectiles // []types.AABB

// Adapter reads via direct field access
func (r *RaylibRenderer) Render(game *engine.Game) {
    r.renderWorld(game.World)
    r.renderPlayer(game.Player)
}
```

---

## Key Entities

### Player (Aggregate Root)

```go
type Player struct {
    AABB          types.AABB     // Position and dimensions
    Velocity      types.Vec2     // Pixels per second
    OnGround      bool           // Collision state
    IsDrilling    bool           // Animation state
    OreInventory  map[string]int // Ore counts by ID
    ItemInventory [5]int         // Item counts
    Money         int            // Currency
    Fuel          float32        // Current fuel
    HP            float32        // Hit points
    SpawnX        float32        // Spawn position for teleport
    SpawnY        float32        // Spawn position for teleport

    // Upgrades (unexported - access via methods)
    engine     upgrades.Engine
    hull       upgrades.Hull
    fuelTank   upgrades.FuelTank
    cargoHold  upgrades.CargoHold
    heatShield upgrades.HeatShield
    drill      upgrades.Drill
}
```

**Stat Accessors** — Clean facade methods for stats:
```go
// Movement stats (from Engine)
player.MaxSpeed()        // float32
player.Acceleration()    // float32
player.FlyAcceleration() // float32
player.MaxUpwardSpeed()  // float32

// Defense stats (from Hull)
player.MaxHP()           // float32

// Resource stats
player.FuelCapacity()    // float32
player.CargoCapacity()   // int

// Heat/Environment stats
player.HeatResistance()  // float32

// Drilling stats
player.DrillSpeedAtSurface()  // float32
player.DrillSpeedAtMaxDepth() // float32
```

**Upgrade Management** — Generic accessors for shop/catalog:
```go
player.GetUpgrade(upgrades.TypeEngine)     // returns upgrades.Upgrade
player.SetUpgrade(newEngine)               // accepts any upgrades.Upgrade
player.GetUpgradeTier(upgrades.TypeEngine) // returns int
```

### Building (Component-Based)

```go
type Building struct {
    Position    components.Position     // AABB wrapper
    Interactable components.Interactable // Type enum
}
```

Types: Market, FuelStation, Hospital, UpgradeShop, ItemShop

---

## Dependencies

Minimal, intentional dependencies:

- `github.com/gen2brain/raylib-go/raylib` — Graphics/audio (adapters only)
- `github.com/stretchr/testify` — Testing utilities
- Go standard library — Everything else

**Strict rule:** No domain code imports Raylib.

---

## Future Extensibility

### Adding New Entities

1. Create entity in `internal/domain/entities/`
2. Add to Game struct in `engine/game.go`
3. Update physics if needed
4. Add rendering in adapter
5. Write tests

### Swapping Renderers

```go
// Create SDL adapter with same interface
type SDLRenderer struct { ... }
func (r *SDLRenderer) Render(game *engine.Game) { ... }

// Swap in main.go
renderer := sdl.NewSDLRenderer()  // Game loop unchanged!
```

### Adding Input Sources

```go
// File-based replay input
type FileInputAdapter struct { frames []InputState }
func (a *FileInputAdapter) ReadInput() input.InputState { ... }

// Swap in main.go
inputAdapter := file.NewFileInputAdapter("replay.bin")
```

---

## Summary

This architecture achieves:

- **Testability**: Domain logic 100% testable without framework
- **Portability**: Could swap Raylib for any renderer
- **Maintainability**: Clear responsibilities, linear data flow
- **Extensibility**: Easy to add new entities, systems, input sources
- **Clarity**: Folder structure = architecture diagram

**Core principle: Domain stays pure, framework stays outside.**
