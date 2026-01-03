# Development Guide

This guide covers common development tasks, workflows, and how to extend the game.

## Table of Contents

1. [Setup](#setup)
2. [Running the Game](#running-the-game)
3. [Testing](#testing)
4. [Debugging](#debugging)
5. [Making Changes](#making-changes)
6. [Code Review Checklist](#code-review-checklist)
7. [Performance & Profiling](#performance--profiling)

---

## Setup

### Prerequisites

- **Go**: 1.21+ (or latest stable)
- **Raylib dependencies** (platform-specific):
  - **macOS**: `brew install raylib`
  - **Linux**: `libasound2-dev`, `mesa-common-dev`, `libx11-dev`, `libxrandr-dev`, `libxinerama-dev`, `libxcursor-dev`, `libxi-dev`
  - **Windows**: See [raylib-go installation](https://github.com/gen2brain/raylib-go#requirements)

### Installation

```bash
# Clone the repository
git clone https://github.com/Kishlin/drill-game.git
cd drill-game

# Download dependencies
go mod download

# Verify setup works
go run cmd/game/main.go
```

---

## Running the Game

### Development Mode

```bash
# Run with output to console (useful for debugging with logs)
go run cmd/game/main.go

# Run with environment variables for logging
LOGLEVEL=debug go run cmd/game/main.go
```

### Build Executable

```bash
# Build optimized binary
go build -o drill-game cmd/game/main.go

# Run the built executable
./drill-game

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o drill-game cmd/game/main.go
GOOS=darwin GOARCH=amd64 go build -o drill-game cmd/game/main.go
GOOS=windows GOARCH=amd64 go build -o drill-game.exe cmd/game/main.go
```

---

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test package
go test -v ./internal/domain/physics
go test -v ./internal/domain/world

# Run specific test
go test -v ./internal/domain/physics -run TestApplyGravity

# Run with test coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Organization

**Physics & Types** (`internal/domain/physics/` and `internal/domain/types/`):
- `collision_test.go` — AABB detection, axis-separated resolution, wall/ceiling/ground
- `movement_test.go` — Acceleration, damping, speed capping
- `gravity_test.go` — Gravity effects, position integration
- `types/aabb_test.go` — AABB intersection, penetration

**World Generation** (`internal/domain/world/`):
- `generator_test.go` — Gaussian distribution, determinism, ore selection
- `world_test.go` — Chunk loading, lazy loading, proactive loading
- `integration_test.go` — End-to-end generation and validation

**Systems** (`internal/domain/systems/`):
- `digging_test.go` — Digging mechanics and ore collection
- `fuel_test.go` — Fuel consumption (active vs idle rates)
- `fuel_station_test.go` — Refueling transactions, cost calculation, edge cases

### Writing New Tests

```go
// Standard pattern: pure function tests
func TestSomething(t *testing.T) {
    // Arrange
    input := someValue
    expected := expectedValue

    // Act
    result := FunctionUnderTest(input)

    // Assert
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Performance Benchmarks

```bash
# Run benchmarks for world generation
go test ./internal/domain/world -bench=. -benchmem

# Example output:
# BenchmarkChunkGeneration-8    500    2.2ms/op    500 B/op
# Chunk generation: ~2.2ms per 16×16 chunk
# Cached tile lookup: ~38ns per tile
```

---

## Debugging

### Logging Output

Add structured logging to debug physics or game logic:

```go
import "log/slog"

// In your code:
slog.Info("Player position", "x", player.AABB.X, "y", player.AABB.Y)
slog.Debug("Collision detected", "tile_x", tileX, "tile_y", tileY)
slog.Error("Physics error", "error", err)
```

Run the game and logs will output to console:

```bash
go run cmd/game/main.go 2>&1 | tee debug.log
```

### Debugging Physics

1. **Check position/velocity values** in `internal/domain/physics/`
2. **Add logs** to movement, gravity, and collision functions
3. **Write test cases** for specific scenarios
4. **Use breakpoints** with Delve debugger:

```bash
# Install Delve debugger (if not already installed)
go install github.com/go-delve/delve/cmd/dlv@latest

# Start debugging
dlv debug cmd/game/main.go

# In dlv: set breakpoints, step through code
(dlv) break main.main
(dlv) continue
(dlv) next
(dlv) print player.Velocity
```

### Common Issues

| Issue | Solution |
|-------|----------|
| Game won't compile | `go mod download` and ensure Go 1.21+ |
| Movement feels wrong | Check physics constants in `internal/domain/physics/constants.go` |
| Collision not working | Verify AABB calculation and CheckCollisions function |
| Tests failing | Run `go test ./...` to see full error messages |

---

## Making Changes

### Adding a New Entity

1. **Create entity file** in `internal/domain/entities/newentity.go`:

```go
package entities

import "github.com/Kishlin/drill-game/internal/domain/types"

type NewEntity struct {
    AABB     types.AABB
    Velocity types.Vec2
    // ... other fields
}

// Implement any required interfaces
```

2. **Add to Game struct** in `internal/domain/engine/game.go`
3. **Update physics system** if entity needs collision/physics
4. **Add rendering** in `internal/adapters/rendering/raylib.go`
5. **Write tests** for entity behavior

### Adding a New Game System

1. **Create system file** in `internal/domain/systems/newsystem.go`:

```go
package systems

type NewSystem struct {
    // ... system state
}

func NewNewSystem() *NewSystem {
    return &NewSystem{}
}

func (ns *NewSystem) Update(game *engine.Game, dt float32) error {
    // Pure domain logic
    return nil
}
```

2. **Wire into Game.Update()** in the correct order
3. **Write tests** in `internal/domain/systems/newsystem_test.go`
4. **Keep framework code in adapters** - no Raylib imports

### Modifying Physics

Physics tuning is in `internal/domain/physics/constants.go`:

```go
const (
    Gravity           = 800          // pixels/sec²
    MaxMoveSpeed      = 450          // pixels/sec
    MoveAcceleration  = 2500         // pixels/sec²
    MoveDamping       = 1000         // pixels/sec²
    FlyAcceleration   = 2500         // pixels/sec²
    MaxUpwardVelocity = -600         // pixels/sec (negative = upward)
    FlyDamping        = 300          // pixels/sec²
)
```

**Process:**
1. Adjust constants in `physics/constants.go`
2. Run game: `go run cmd/game/main.go`
3. Test feel and responsiveness
4. Write/update unit tests for expected behavior
5. Verify existing tests still pass: `go test ./internal/domain/physics -v`

---

## Code Review Checklist

When submitting changes, verify:

### Architecture Compliance
- [ ] No Raylib imports in `internal/domain/`
- [ ] New domain code is pure functions (testable without framework)
- [ ] Framework integration stays in `internal/adapters/`
- [ ] Data flow is clear (domain → adapters → application)

### Testing
- [ ] All domain changes have unit tests
- [ ] Tests run successfully: `go test ./...`
- [ ] New tests follow existing patterns (pure function tests)
- [ ] Coverage maintained or improved

### Code Quality
- [ ] Code follows Go idioms (gofmt, effective Go)
- [ ] Error handling is appropriate (minimal in game loop, errors in domain)
- [ ] Comments explain "why", not "what" (code is self-documenting)
- [ ] Variable names are clear and concise

### Performance
- [ ] No new allocations in hot paths (physics frame)
- [ ] Used pointers for large types (Player, World, Game)
- [ ] Value types for small types (Vec2, AABB)
- [ ] No unnecessary copies or indirection

### Documentation
- [ ] CLAUDE.md updated if behavior changes
- [ ] ARCHITECTURE.md updated for structural changes
- [ ] Comments added for non-obvious logic

---

## Performance & Profiling

### CPU Profiling

```bash
# Generate CPU profile
go test -cpuprofile=cpu.prof ./...

# Analyze with pprof
go tool pprof cpu.prof

# In pprof:
# top         - show top functions by CPU time
# list Func   - show source code for function
# web         - generate visualization (requires graphviz)
```

### Memory Profiling

```bash
# Generate memory profile
go test -memprofile=mem.prof ./...

# Analyze
go tool pprof mem.prof
```

### Live Profiling (while game runs)

```bash
# Add import to cmd/game/main.go:
import _ "net/http/pprof"

// Add in main():
go func() {
    slog.Info("Profiling server listening", "url", "http://localhost:6060/debug/pprof")
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

# Run game
go run cmd/game/main.go

# In another terminal:
go tool pprof http://localhost:6060/debug/pprof/profile
```

### Performance Targets

- **Frame Time**: < 16.6ms at 60 FPS
- **Chunk Generation**: < 5ms per 16×16 chunk
- **Tile Lookup**: < 50ns per tile (sparse map)
- **Memory**: < 100MB for full world with entities

---

## Git Workflow

### Creating a Feature Branch

```bash
git checkout -b feature/description-of-change

# Make changes, test
go test ./...

# Commit
git add .
git commit -m "Feature: description of change

Detailed explanation of what and why.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>"

# Push
git push -u origin feature/description-of-change
```

### Code Review

1. Push to branch
2. Create pull request with description
3. Ensure CI/tests pass
4. Address review feedback
5. Merge when approved

---

## Common Commands Reference

```bash
# Development
go run cmd/game/main.go              # Run game
go test ./...                         # Run all tests
go test -cover ./...                 # Test with coverage
go fmt ./...                          # Format code
go vet ./...                          # Lint code

# Building
go build -o drill-game cmd/game/main.go        # Build executable
go build -ldflags="-s -w" -o drill-game ...    # Optimized build (smaller)

# Specific tests
go test -v ./internal/domain/physics           # Physics tests
go test -v ./internal/domain/world             # World generation tests
go test -v -run TestName ./...                 # Single test

# Benchmarks
go test ./internal/domain/world -bench=. -benchmem

# Debugging
dlv debug cmd/game/main.go           # Start debugger
go test -cpuprofile=cpu.prof ./...   # CPU profile
```

---

## Resources

- [CLAUDE.md](../CLAUDE.md) - Quick reference for AI assistants
- [ARCHITECTURE.md](ARCHITECTURE.md) - Technical architecture details
- [GAME_DESIGN.md](GAME_DESIGN.md) - Game mechanics and progression
- [Go Effective Guide](https://go.dev/doc/effective_go)
- [Raylib Documentation](https://www.raylib.com/)
