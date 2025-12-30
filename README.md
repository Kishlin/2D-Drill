# Drill Game 🚜⛏️

A 2D mining progression game built with Go and Raylib. Dig deep, mine ores, upgrade your vehicle, and venture even deeper!

## Overview

Control a small drilling vehicle on a planet rich with valuable ores. The deeper you dig, the rarer the treasures - but also the more dangerous the environment. Mine, sell, upgrade, and repeat in this addictive progression loop inspired by classics like Motherload.

**Current Status**: 🚧 Early Development

## Features (Planned)

- ⛏️ Deep mining gameplay with procedurally generated worlds
- 💎 Diverse ore types with rarity based on depth
- 🔧 Comprehensive upgrade system (speed, drilling, survivability, cargo)
- 🌡️ Environmental hazards (heat, pressure)
- 📊 Progression curve from surface to legendary depths
- 🎮 Smooth controls and satisfying game feel

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

### Phase 1: Core Gameplay (Current)
- [ ] Basic game loop and window
- [ ] Player movement and controls
- [ ] World generation and rendering
- [ ] Tile-based collision and drilling
- [ ] Simple ore collection

### Phase 2: Progression
- [ ] Ore variety and rarity system
- [ ] Currency and shop system
- [ ] Upgrade mechanics
- [ ] Heat hazard system
- [ ] Save/load functionality

### Phase 3: Polish
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
