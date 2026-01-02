# Drill Game 🚜⛏️

A 2D mining progression game built with Go and Raylib. Dig deep, mine ores, upgrade your vehicle, and venture even deeper!

## Overview

Control a small drilling vehicle on a planet rich with valuable ores. The deeper you dig, the rarer the treasures - but also the more dangerous the environment. Mine, sell, upgrade, and repeat in this addictive progression loop inspired by classics like Motherload.

**Current Status**: 🚧 Early Development

## Features

### Implemented ✅
- ⛏️ **Procedurally Generated Worlds**: Seeded chunk-based generation for infinite depth
- 💎 **7 Ore Types**: Copper, Iron, Silver, Gold, Mythril, Platinum, Diamond with Gaussian depth distribution
- 🎮 **Smooth Gameplay**: 60 FPS player movement, physics, and directional digging with AABB collision
- 🗺️ **Chunk Loading**: Lazy 16×16 chunk loading around player (3×3 grid)
- 📐 **Deterministic Generation**: Same seed = reproducible worlds
- ⛏️ **Directional Digging**: Dig downward (S) or through left/right walls (A/D) when grounded
- 📦 **Ore Inventory System**: Automatic ore collection with real-time inventory display
- 🏪 **Shop System**: Visible shop on map, sell entire inventory for currency (E key)
- ⛽ **Fuel System**: Limited fuel tank (10L) with activity-based consumption (active: 0.333 L/s, idle: 0.0833 L/s)

### Planned (Phase 2+)
- ⛽ **Fuel Mechanics**: Game over behavior at zero fuel, refueling system, fuel efficiency upgrades
- 🔧 **Comprehensive upgrade system**: Speed, drilling, survivability, cargo, fuel capacity/efficiency
- 🌡️ **Environmental hazards**: Heat, pressure, lava, gas pockets, underwater areas
- 🎨 **Polish & Content**: Particle effects, sound effects, UI improvements, more ores/biomes, achievements

## Tech Stack

- **Language**: Go (latest stable)
- **Graphics/Audio**: Raylib (via [raylib-go](https://github.com/gen2brain/raylib-go))
- **Testing**: [testify](https://github.com/stretchr/testify)

## Getting Started

### Prerequisites

- Go 1.21+ (or latest stable)
- Raylib dependencies for your platform:
  - **macOS**: `brew install raylib`
  - **Linux**: Install raylib dev packages (`libasound2-dev`, `mesa-common-dev`, `libx11-dev`, etc.)
  - **Windows**: See [raylib-go installation](https://github.com/gen2brain/raylib-go#requirements)

### Installation

```bash
# Clone the repository
git clone https://github.com/Kishlin/drill-game.git
cd drill-game

# Download dependencies
go mod download

# Run the game
go run cmd/game/main.go
```

### Development

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Build executable
go build -o drill-game cmd/game/main.go
```

## Project Structure

```
drill-game/
├── cmd/game/           # Application entry point
├── internal/
│   ├── engine/         # Core game loop, asset management
│   ├── entities/       # Game objects (player, tiles, ores)
│   ├── systems/        # Game systems (camera, physics, mining)
│   └── ui/             # User interface components
├── assets/             # Sprites, sounds, fonts
├── docs/               # Documentation
└── .github/            # GitHub configuration
```

## Documentation

- [Architecture](docs/ARCHITECTURE.md) - Technical design and structure
- [Game Design](docs/GAME_DESIGN.md) - Gameplay mechanics and progression
- [Copilot Instructions](.github/copilot-instructions.md) - AI-assisted development guidelines

## Roadmap

### Phase 1: Core Gameplay & World Generation ✅ Complete
- [x] Game loop and window management
- [x] Player movement, controls, and physics
- [x] Procedurally generated worlds with seeded RNG
- [x] 7 ore types with Gaussian depth distribution
- [x] Tile-based collision (AABB) and axis-separated resolution
- [x] Directional digging system (downward with grid alignment, left/right while grounded)
- [x] Chunk loading (16×16 chunks, 3×3 proactive grid)
- [x] Ore inventory system (automatic collection on digging)
- [x] 38 unit tests + 10 integration tests
- [x] Deterministic world generation

### Phase 2: Progression System
- [x] Ore inventory system
- [x] Currency and shop system
- [x] Ore selling mechanics
- [ ] Upgrade mechanics (speed, drilling, capacity)
- [ ] Mining duration per ore type
- [ ] Save/load functionality

### Phase 3: Polish & Effects (Planned)
- [ ] Particle effects and juice
- [ ] Sound effects and music
- [ ] UI/UX improvements
- [ ] Visual feedback and polish

### Phase 4: Extended Content
- [ ] More ore types and upgrades
- [ ] Additional hazards
- [ ] Achievement system
- [ ] Challenge modes

### Future Vision
- Cross-platform release (Steam, Mobile)
- Online leaderboards
- Daily challenges and events
- Workshop/mod support

## Contributing

This is currently a personal project, but feedback and suggestions are welcome! Feel free to open issues for bugs or feature ideas.

## License

TBD (will be decided before public release)

## Credits

- **Developer**: [Your Name]
- **Inspired by**: Motherload (XGen Studios), Steamworld Dig, Terraria
- **Built with**: [Raylib](https://www.raylib.com/) and [Go](https://go.dev/)

---

*Dig deep, upgrade hard, repeat! 🚀*
