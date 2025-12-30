# Drill Game - Copilot Instructions

## Project Overview
Drill Game is a 2D mining progression game built with Golang and Raylib. Players control a drilling vehicle, mining ores of increasing rarity as they dig deeper, then sell and upgrade to progress further.

**Tech Stack**: Go + Raylib (via raylib-go bindings)

## Architecture Principles

### 1. Golang Structure
- **Game-Oriented Architecture**: Entity-Component patterns where appropriate, but pragmatic over dogmatic
- **Directory Structure**: 
  - `cmd/game/`: Application entry point with main()
  - `internal/engine/`: Core game loop, rendering, resource management
  - `internal/entities/`: Game objects (Player, Ore, Tile, etc.)
  - `internal/systems/`: Game systems (Camera, Physics, Mining, Upgrades, etc.)
  - `internal/ui/`: User interface, menus, HUD
  - `assets/`: Sprites, sounds, fonts (not in version control if large)
  - Tests are colocated with the code they test
- **Naming Conventions**:
  - Types are clear and descriptive (e.g., `Player`, `MiningSystem`, `CameraController`)
  - Interfaces describe behavior (e.g., `Renderable`, `Updatable`, `Collidable`)
  - Use Go's short variable names in tight scopes, descriptive names in wider scopes

### 2. Testing Strategy
- **Unit tests**: `*_test.go` files in same package, suffix `TestSuite` for suite-based tests
- **Integration tests**: Test system interactions, suffix `IntegrationTestSuite`
- Use `testify/suite` for test organization when beneficial
- Simple tests can use standard `testing.T` without suites
- **Test Lifecycle**:
  - Use `SetupSuite()` for suite-wide setup (runs once before all tests)
  - Use `SetupSubTest()` for cleanup/reset between test cases
  - Prefer `s.Run()` over `s.T().Run()` when you need per-test cleanup

### 3. Game Loop Patterns
- **Fixed timestep** for physics updates (typically 60 FPS)
- **Variable timestep** for rendering
- Separate Update() and Render() phases
- Use delta time for frame-independent movement

### 4. Error Handling
- Always wrap errors with context: `fmt.Errorf("loading texture %s: %w", path, err)`
- Use structured logging with `slog` package
- Panic is acceptable for unrecoverable game initialization errors (missing critical assets)
- Log at appropriate levels: Debug for frame operations, Info for state changes, Warn for recoverable issues, Error for failures

### 5. Resource Management
- Load assets once during initialization, not per-frame
- Unload resources when no longer needed (textures, sounds)
- Use defer for cleanup: `defer rl.UnloadTexture(texture)`
- Consider resource pooling for frequently created/destroyed objects

## Code Style

### Golang
```go
// Use clear, explicit boolean checks for readability
if isGameOver == false {
    // Continue game logic
}

// Always wrap errors with context
if err != nil {
    return fmt.Errorf("initializing player: %w", err)
}

// Prefer early returns to reduce nesting
func Update(dt float32) error {
    if isPaused {
        return nil
    }
    
    if err := updatePhysics(dt); err != nil {
        return fmt.Errorf("updating physics: %w", err)
    }
    
    return nil
}

// Use struct embedding for component-like behavior
type Player struct {
    Position  Vector2
    Velocity  Vector2
    Stats     PlayerStats
    Inventory Inventory
}

// Constants for game values (not magic numbers)
const (
    PlayerSpeed      = 100.0
    DrillSpeedBase   = 1.0
    MaxDepth         = 10000
    TileSize         = 16
)
```

### Test Naming
```go
type PlayerTestSuite struct {
    suite.Suite
    player *Player
}

func (s *PlayerTestSuite) SetupSuite() {
    // Suite-wide setup
}

func (s *PlayerTestSuite) SetupSubTest() {
    // Reset state between test cases
    s.player = NewPlayer()
}

func (s *PlayerTestSuite) TestMovement() {
    s.Run("moves right when right key pressed", func() {
        s.player.MoveRight(1.0)
        s.Greater(s.player.Position.X, 0.0)
    })
    
    s.Run("stops at max speed", func() {
        for i := 0; i < 100; i++ {
            s.player.MoveRight(1.0)
        }
        s.LessOrEqual(s.player.Velocity.X, s.player.Stats.MaxSpeed)
    })
}
```

## Common Patterns

### Game State Management
```go
type GameState int

const (
    StateMenu GameState = iota
    StatePlaying
    StatePaused
    StateGameOver
)

type Game struct {
    state GameState
    // ... other fields
}

func (g *Game) Update() error {
    switch g.state {
    case StateMenu:
        return g.updateMenu()
    case StatePlaying:
        return g.updateGameplay()
    // ... etc
    }
    return nil
}
```

### Component-like Entities
```go
// Simple component pattern - no need for complex ECS
type Renderable interface {
    Render()
}

type Updatable interface {
    Update(dt float32) error
}

type Entity struct {
    Position Vector2
    Sprite   rl.Texture2D
}

func (e *Entity) Render() {
    rl.DrawTextureV(e.Sprite, e.Position, rl.White)
}
```

### Resource Loading
```go
type AssetManager struct {
    textures map[string]rl.Texture2D
    sounds   map[string]rl.Sound
}

func (am *AssetManager) LoadTexture(name, path string) error {
    texture := rl.LoadTexture(path)
    if texture.ID == 0 {
        return fmt.Errorf("failed to load texture: %s", path)
    }
    am.textures[name] = texture
    return nil
}
```

## File Locations

- **Entry point**: `cmd/game/main.go`
- **Core engine**: `internal/engine/`
- **Game entities**: `internal/entities/`
- **Game systems**: `internal/systems/`
- **UI components**: `internal/ui/`
- **Assets**: `assets/` (sprites, sounds, fonts)

## Performance Considerations

1. **Avoid allocations in game loop** - reuse objects, use object pools if needed
2. **Batch rendering** - group similar draw calls together
3. **Spatial partitioning** - use quad trees or grid for collision detection at scale
4. **Profile regularly** - use Go's pprof to identify bottlenecks
5. **Pre-allocate slices** - `make([]Ore, 0, expectedCapacity)`

## When Suggesting Code

1. **Follow game development patterns** - favor simplicity and performance
2. **Include error handling** - never ignore errors, especially for resource loading
3. **Add tests** for game logic (not rendering code)
4. **Use constants** instead of magic numbers
5. **Comment non-obvious game logic** - especially physics/math calculations
6. **Consider frame budget** - avoid expensive operations in Update/Render loops
7. **Use Raylib idiomatically** - follow raylib-go examples and patterns
