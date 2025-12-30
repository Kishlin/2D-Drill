# Architecture

## Overview

Drill Game is a 2D mining progression game built with Go and Raylib. The architecture follows pragmatic game development patterns, prioritizing clarity, performance, and maintainability.

## Technology Stack

- **Language**: Go (latest stable version)
- **Graphics/Audio**: Raylib (via raylib-go bindings)
- **Testing**: testify/suite
- **Logging**: slog (standard library)

## Project Structure

```
drill-game/
├── cmd/
│   └── game/           # Application entry point
│       └── main.go     # Initialization, game loop
├── internal/
│   ├── engine/         # Core game systems
│   │   ├── game.go           # Main game loop, state management
│   │   ├── assets.go         # Asset loading and management
│   │   └── window.go         # Window configuration
│   ├── entities/       # Game objects
│   │   ├── player.go         # Player/vehicle entity
│   │   ├── tile.go           # World tile (dirt, ore, etc.)
│   │   ├── ore.go            # Ore definitions and instances
│   │   └── upgrade.go        # Upgrade definitions
│   ├── systems/        # Game systems
│   │   ├── camera.go         # Camera controller (following player)
│   │   ├── physics.go        # Basic 2D physics (gravity, collision)
│   │   ├── mining.go         # Mining mechanics and logic
│   │   ├── world.go          # World generation and management
│   │   └── progression.go    # Upgrade system, currency, shop
│   └── ui/             # User interface
│       ├── hud.go            # In-game HUD (depth, heat, inventory)
│       ├── menu.go           # Main menu
│       └── shop.go           # Upgrade shop UI
├── assets/
│   ├── sprites/        # PNG textures
│   ├── sounds/         # Audio files
│   └── fonts/          # Font files
├── docs/
│   ├── ARCHITECTURE.md       # This file
│   └── GAME_DESIGN.md        # Game mechanics documentation
├── .github/
│   └── copilot-instructions.md
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## Core Concepts

### Game Loop

Standard game loop with fixed timestep for physics and variable timestep for rendering:

```go
func (g *Game) Run() error {
    const targetFPS = 60
    const dt = 1.0 / targetFPS
    
    for !rl.WindowShouldClose() {
        // Update (fixed timestep)
        if err := g.Update(dt); err != nil {
            return err
        }
        
        // Render (variable timestep)
        rl.BeginDrawing()
        g.Render()
        rl.EndDrawing()
    }
    
    return nil
}
```

### State Management

Game uses a simple state machine:
- **Menu**: Main menu, settings
- **Playing**: Active gameplay
- **Paused**: Game paused, can access shop
- **GameOver**: Death/retirement screen (future: leaderboard submission)

### Entity Model

Entities are simple structs with behavior methods. No complex ECS framework - keep it simple:

```go
type Player struct {
    Position  rl.Vector2
    Velocity  rl.Vector2
    Stats     PlayerStats  // Speed, drill power, heat resistance, etc.
    Inventory map[OreType]int
    Health    float32
    Heat      float32
}

func (p *Player) Update(dt float32) error {
    // Physics, collision, mining logic
}

func (p *Player) Render(camera Camera) {
    // Draw sprite at position relative to camera
}
```

### World Representation

2D grid-based world with procedural generation:
- Each tile has a type (empty, dirt, ore)
- Ores have rarity based on depth
- Tiles can be destroyed by drilling
- World generates as player descends (infinite depth)

```go
type World struct {
    tiles  map[Vector2i]*Tile  // Sparse map for memory efficiency
    depth  int                 // Current deepest generated depth
}
```

### Camera System

2D camera follows the player with smooth scrolling:
- Vertical scrolling as player descends
- Horizontal tracking with some elasticity
- Viewport shows portion of the world

### Progression System

Core gameplay loop:
1. Mine ores by drilling through tiles
2. Return to surface (or die trying)
3. Sell ores for currency
4. Purchase upgrades (speed, drill power, heat resistance, hull strength)
5. Descend deeper for rarer ores
6. Repeat

## Performance Considerations

### Memory Management
- Use object pools for frequently created/destroyed entities (particles, debris)
- Sparse world representation (only store non-empty tiles)
- Unload off-screen chunks (future optimization)

### Rendering Optimization
- Only render visible tiles (camera viewport culling)
- Batch similar draw calls
- Use texture atlases to reduce texture switches

### Update Optimization
- Spatial partitioning for collision detection (grid or quad-tree)
- Skip updates for entities far from player
- Profile regularly to identify bottlenecks

## Cross-Platform Considerations

Raylib supports multiple platforms out of the box:
- **Desktop**: Windows, macOS, Linux (primary target)
- **Web**: WebAssembly (future)
- **Mobile**: iOS, Android (future)

Current focus: Desktop development first, then expand to other platforms.

## Future Architecture Considerations

As the game grows, consider:
- **Save system**: JSON or binary serialization of game state
- **Mod support**: Lua scripting for custom ores, upgrades
- **Online features**: REST API for leaderboards, events (Go backend)
- **Level editor**: Separate tool for creating custom scenarios
- **Multiplayer**: WebSocket-based co-op or competitive modes

Keep these in mind but don't over-engineer early. Build what's needed now.

## Testing Strategy

- **Unit tests**: Core game logic (progression calculations, mining mechanics)
- **Integration tests**: System interactions (physics + collision)
- **Manual testing**: Rendering, feel, game balance
- No automated rendering tests (visual testing is manual)

## Dependencies

Keep dependencies minimal:
- `github.com/gen2brain/raylib-go/raylib` - Graphics and audio
- `github.com/stretchr/testify` - Testing utilities
- Standard library for everything else

Avoid premature dependency addition. Prefer standard library solutions.
