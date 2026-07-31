# Drill Game 🚜⛏️

A 2D mining progression game built with Go and Raylib. Drill deep, mine ores, upgrade your vehicle, and venture even deeper!

## Overview

Control a small drilling vehicle on a planet rich with valuable ores. The deeper you drill, the rarer the treasures - but also the more dangerous the environment. Mine, sell, upgrade, and repeat in this addictive progression loop inspired by classics like Motherload.

**Current Status**: 🚧 Early Development

## Features

### Implemented ✅
- ⛏️ **Procedurally Generated Worlds** — Seeded chunk-based generation for infinite depth
- 💎 **7 Ore Types** — Copper, Iron, Silver, Gold, Mythril, Platinum, Diamond with Gaussian depth distribution
- 🎮 **Smooth Gameplay** — 60 FPS movement, physics, directional drilling with AABB collision
- 🗺️ **Chunk Loading** — Lazy 16×16 chunks around player
- 📦 **Ore Inventory & Market System** — Automatic collection, sell for currency
- ⛽ **Fuel System** — Limited tank with activity-based consumption

See [CLAUDE.md](CLAUDE.md) for current feature status and configuration.

### Planned (Phase 2+)
See [GAME_DESIGN.md](docs/GAME_DESIGN.md) for detailed game mechanics and progression system.

## Tech Stack

- **Language**: Go (latest stable)
- **Graphics/Audio**: Raylib (via [raylib-go](https://github.com/gen2brain/raylib-go))
- **Testing**: [testify](https://github.com/stretchr/testify)

## Getting Started

### Prerequisites

- Go 1.21+ (or latest stable)
- A C compiler and the Raylib system dependencies for your platform. raylib-go builds Raylib from source via cgo, so these **must** be installed before `go run` — follow [raylib-go § Requirements](https://github.com/gen2brain/raylib-go#requirements), which lists the packages per OS.

  Missing them shows up as a cgo compile error rather than a Go one, e.g. `fatal error: wayland-client-core.h: No such file or directory`.

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

See [CLAUDE.md](CLAUDE.md) for quick commands to run tests and build the game.

## Project Structure

```
drill-game/
├── cmd/game/                    # Application entry point
├── internal/
│   ├── domain/                  # Pure business logic
│   │   ├── engine/              # Game loop orchestration
│   │   ├── entities/            # Game objects (player, tiles, ores)
│   │   ├── systems/             # Game systems (physics, drilling, fuel)
│   │   ├── world/               # Procedural generation, chunk loading
│   │   └── physics/             # Physics functions, collision
│   └── adapters/                # Framework integration (Raylib only)
│       ├── input/               # Keyboard input mapping
│       └── rendering/           # Raylib rendering
├── docs/                        # Documentation
└── .github/                     # GitHub configuration
```

## Documentation Guide

Start here based on what you need:

- **[CLAUDE.md](CLAUDE.md)** — Quick reference for AI assistants, commands, and configuration
- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** — Technical design, hexagonal architecture, complete system reference
- **[docs/GAME_DESIGN.md](docs/GAME_DESIGN.md)** — Game mechanics, progression system, future features
- **[docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)** — Development workflows, testing, debugging, how to extend the game

## Roadmap

### Phase 1: Core Gameplay & World Generation ✅ Complete
- [x] Game loop and window management
- [x] Player movement, controls, and physics
- [x] Procedurally generated worlds with seeded RNG
- [x] 7 ore types with Gaussian depth distribution
- [x] Tile-based collision (AABB) and axis-separated resolution
- [x] Directional drilling system (downward with grid alignment, left/right while grounded)
- [x] Chunk loading (16×16 chunks, 3×3 proactive grid)
- [x] Ore inventory system (automatic collection on drilling)
- [x] 38 unit tests + 10 integration tests
- [x] Deterministic world generation

### Phase 2: Progression System
- [x] Ore inventory system
- [x] Currency and market system
- [x] Ore selling mechanics
- [x] Fuel consumption system
- [x] Fuel station (refueling mechanic)
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

- **Developer**: Pierre-Louis Legrand
- **Inspired by**: Motherload (XGen Studios), Steamworld Dig Series, Terraria
- **Built with**: [Raylib](https://www.raylib.com/) and [Go](https://go.dev/)

---

*Drill deep, upgrade hard, repeat! 🚀*
