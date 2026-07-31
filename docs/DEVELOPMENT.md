# Development Guide

This guide covers common development tasks, workflows, and how to extend the game.

## Table of Contents

1. [Setup](#setup)
2. [Running the Game](#running-the-game)
3. [Testing](#testing)
4. [Architecture Check](#architecture-check)
5. [Style Check](#style-check)
6. [Claude Code Tooling](#claude-code-tooling)
7. [Debugging](#debugging)
8. [Making Changes](#making-changes)
9. [Code Review Checklist](#code-review-checklist)
10. [Performance & Profiling](#performance--profiling)

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

### Test Level (Level -1)

A special development level is available with an advanced player configuration for easier testing:

```bash
# Edit cmd/game/main.go to use test level
# Change: levels.GetLevelConfig(1) → levels.GetLevelConfig(-1)
# Or: levels.GetLevelConfig(-2)  for TestBoss testing
# Or: levels.GetLevelConfig(-3)  for SentinelBoss testing
```

**Test Level Player Stats:**
- Starting money: $100,000
- Starting items: 3 Teleports, 5 Repairs, 5 Refuels, 10 Bombs, 20 Big Bombs
- Engine: Max tier (Mk5)
- Drill: Max tier (Mk5)
- Hull/FuelTank/CargoHold/HeatShield: Mid tier (Mk3)

This allows testing deep mining, bomb mechanics, and shop interactions without grinding through early game.

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
- `drilling_test.go` — Drilling mechanics and ore collection
- `fuel_test.go` — Fuel consumption (active vs idle rates)

**Effects** (`internal/domain/effects/`):
- `money_test.go` — TakeMoney, AddMoney effects
- `stats_test.go` — SetFuel, SetHP effects
- `upgrades_test.go` — SetEngine, SetHull, etc. effects
- `inventory_test.go` — ClearOreInventory, AddItem effects
- `processor_test.go` — Processor applying effects to player

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

## Architecture Check

The hexagonal boundary is enforced by a script, not by discipline:

```bash
./scripts/architecture-check.sh
```

It audits every `.go` file under `internal/domain/` (tests included) for four
violations and exits non-zero on any of them:

| # | Rule | Why |
|---|------|-----|
| 1 | No `raylib-go` import | Domain must compile and test without a graphics stack |
| 2 | No `internal/adapters` import | Dependency direction — adapters depend on domain, never the reverse |
| 3 | No `cmd/` import | Wiring flows inward from `main.go` only |
| 4 | No `rl.*` types | Domain speaks `types.Vec2` / `types.AABB` |

Failures print `file:line` plus the fix. The fix is essentially never to loosen
the check — define the interface in the domain package and implement it in
adapters instead.

Run it before every commit that touched `internal/domain/`, and after any
refactor that moved packages around.

---

## Style Check

The explicit-false-boolean rule is enforced by a script, not by review:

```bash
./scripts/style-check.sh                 # sweep internal/ and cmd/
./scripts/style-check.sh path/to/file.go # check specific files
```

**This codebase never uses the `!` operator.** Write `x == false` instead of
`!x` — in `if`, in `for`, in assignments, and in `&&` / `||` operands. `!=` is
unaffected. Full rationale in [.claude/rules/go-style.md](../.claude/rules/go-style.md).

```go
// NO                                  // YES
if !ok { ... }                         if ok == false { ... }
wasAirborne := !player.OnGround        wasAirborne := player.OnGround == false
if a && !b { ... }                     if a && b == false { ... }
```

The check strips string literals, rune literals, raw strings, and comments
before matching, so a `!` inside those never trips it. Multi-line block comments
are the one unhandled case; the codebase does not use them inside function
bodies, which is the only place `!` can appear.

Note that grepping for `if !` is **not** equivalent — a negation can sit far
from the `if`, as in `if _, ok := m[k]; !ok {`. Use the script.

---

## Claude Code Tooling

Shared configuration lives in `.claude/settings.json` (checked in).
Machine-specific permissions belong in `.claude/settings.local.json`, which is
gitignored — do not put project-wide settings there.

### Permission policy

The versioned allow-list contains **only commands that are harmless in every
case**. The test is not "would I usually be fine with this?" but "is there any
invocation of this that could damage the working tree?" If yes, it does not go
in the versioned file — it is each developer's own call.

| Goes in `settings.json` (versioned) | Goes in `settings.local.json` (yours) |
|---|---|
| Reads and queries: `ls`, `grep`, `head`, `wc`, `tree`, `jq`, `command -v` | Anything that writes: `gofmt -w`, `go fmt`, `go mod tidy` |
| Read-only Go: `go doc`, `go list`, `go vet`, `gofmt -l`, `gofmt -d` | Anything that builds or executes repo code: `go build`, `go test`, `go run` |
| Read-only git: `status`, `diff`, `log`, `show`, `check-ignore` | Anything that mutates git state: `git add`, `git reset`, `git rm` |
| The audit scripts — both are pure `grep` | Filesystem mutation: `chmod`, `rm`, and `find` (which takes `-delete` and `-exec`) |

Two entries worth explaining:

- **`gofmt -l` and `gofmt -d` are versioned, plain `gofmt` is not.** The first
  two only list and diff; the glob `gofmt:*` would also permit `gofmt -w`.
  Splitting the flag out keeps the read-only forms frictionless.
- **`find` is local, not versioned.** It reads — until someone appends
  `-delete` or `-exec rm {} \;`. It fails the "any invocation" test.

### Deny list

`settings.json` also carries a `deny` list — `rm -rf`, `git reset --hard`,
`git clean -fd`, `git push --force`. **Deny beats allow**, including a local
allow, so this is a project-wide floor no individual config can lower. Note
that `git reset:*` is allowed locally while `git reset --hard:*` stays denied;
that layering is intentional.

### Automatic hooks

`PostToolUse` hooks fire whenever Claude writes or edits a file:

| Hook | Trigger | Effect |
|------|---------|--------|
| `gofmt` | any `*.go` | Formats the file in place |
| `hexagonal boundary check` | any `internal/domain/**.go` | Runs `architecture-check.sh`, **blocks the edit** on failure |
| `go style check` | any `*.go` | Runs `style-check.sh` on that file, **blocks the edit** on failure |

The blocking hooks mean a layering or style violation is rejected at write time
with the diagnostic fed back — it cannot land silently and be discovered later.

### Rules

`.claude/rules/` holds rules that must always hold, kept short and separate from
CLAUDE.md so they are not diluted by surrounding prose.

| File | Covers |
|------|--------|
| [go-style.md](../.claude/rules/go-style.md) | No `!` negation |

A rule belongs here only if a script in `scripts/` can enforce it and a hook can
block on it. A rule that relies on the model remembering it belongs in CLAUDE.md
and should be understood as a strong preference, not a guarantee.

### Skills

| Skill | Use when |
|-------|----------|
| `architecture-check` | Full boundary sweep — before a commit, after a refactor, or auditing someone else's work |

### Slash commands

| Command | Does |
|---------|------|
| `/new-boss <boss_name>` | Scaffolds a boss end to end: domain package, state machine, renderer, level wiring, tests |

Adding a skill: create `.claude/skills/<name>/SKILL.md` with `name` and
`description` in YAML frontmatter. The description is what Claude matches
against, so state *when* to use it, not just what it does. Back anything
non-trivial with a script in `scripts/` so a human can run it too.

Adding a command: create `.claude/commands/<name>.md` with a `description` and
optional `argument-hint` in frontmatter; `$1` interpolates the argument.

---

## Testing Fall Damage

To manually test the fall damage system:

1. **Run the game**: `go run cmd/game/main.go`
2. **Small jump**: Jump from ground level - should take no damage
3. **Medium fall**: Find a high platform and jump/fall from it - should take damage based on formula
4. **Check HP display**: Look at top-left debug overlay showing `HP: X.X` value
5. **Lethal fall**: Fall from extreme height - HP should clamp at 0

**Formula reminder:**
- Threshold: 500 px/sec downward velocity
- Below threshold: No damage
- At threshold: 0 damage
- Above threshold: `damage = (velocity - 500) / 20`

**Examples:**
- 600 px/sec fall → 5 damage
- 700 px/sec fall → 10 damage (lethal)
- 800 px/sec fall → 15 damage (clamped at 0)

---

## Testing Drilling Animation

To manually test the drilling animation system:

1. **Run the game**: `go run cmd/game/main.go`
2. **Vertical Drilling** (S/Down key):
   - Stand above a dirt or ore tile
   - Press S/Down key
   - Observe player animate smoothly to tile center over ~1 second
   - Tile disappears when animation completes
   - Ore collected if it was an ore tile
3. **Horizontal Drilling** (A/D when grounded against wall):
   - Stand on ground next to a drillable wall
   - Press A/D to move into the wall
   - Animation starts: player moves to tile center while staying grounded
   - After ~1 second: wall breaks, player moves into space
   - Continue left/right to drill through walls
4. **Mid-Air Behavior**:
   - Jump up and press A/D
   - Player should NOT drill mid-air
   - Player bounces off solid tiles instead
5. **Animation Interruption**:
   - Start drilling (animation begins)
   - Try pressing other keys (E, I, etc.)
   - Inputs should be ignored until animation completes
   - Fuel should still consume during animation
6. **Animation Duration**:
   - Drill a tile at surface and time it
   - Should take approximately 1.0 second for dirt at ground level
   - Duration increases with depth (up to 24 seconds at max depth)
   - Ore multipliers increase duration (Copper 1.2x → Diamond 3.0x)
   - Drill upgrades reduce duration via depth-scaled divisor

**Debug Display:**
- Top-left corner shows: `IsDrilling: true/false`
- Verify IsDrilling is true during animation, false otherwise

**Formula (variable-duration animation):**
```
position = start + (target - start) * (elapsed / duration)
```
- `elapsed = 0` → position = start
- `elapsed = duration/2` → position = midpoint
- `elapsed = duration` → position = target (tile removed, ore collected)

Duration is calculated based on depth, ore type, and drill upgrades:
- Base duration: 1.0s (surface) to 24s (max depth)
- Ore hardness multiplier: 1.2x-3.0x
- Drill upgrade divisor: depth-scaled (more effective at depth)

---

## Testing Heat Damage & Temperature

To manually test the heat system:

1. **Run the game**: `go run cmd/game/main.go`
2. **Surface (safe)**: Start at ground level - no heat damage (15°C)
3. **Shallow depth**: Drill down 1-2 screens - temperature rises but below 50°C resistance, no damage
4. **Medium depth**: Drill down 10+ screens - temperature exceeds base resistance (50°C), take gradual damage
5. **Deep dive** (before upgrades): Attempt to go very deep without heat upgrades - will take increasing exponential damage
6. **Buy Heat Shield Mk1**: Return to surface, get $200+, purchase heat shield upgrade (orange-red shop)
7. **Descend further**: With Mk1 (90°C resistance), can drill deeper before heat kills you
8. **Repeat progression**: Each upgrade unlocks approximately 8,000px more safe depth

**Debug Display:**
- Top-left corner shows: `Temperature: X.X°C (Resistance: Y.Y°C)`
- Current HP displayed: `HP: Z.Z`
- Heat damage only occurs when Temperature > Resistance

**Formula verification:**
- At 60°C resistance, 100°C temperature (40°C excess):
  - Damage/sec = 0.5 * (40/10)^1.5 = 0.5 * 8 ≈ 4 HP/sec
  - Full 10 HP depleted in ~2.5 seconds of continuous exposure
- At 140°C resistance, 150°C temperature (10°C excess):
  - Damage/sec = 0.5 * (10/10)^1.5 = 0.5 * 1 = 0.5 HP/sec
  - Takes 20 seconds to reach lethal damage

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

### Creating a New Level

Levels are defined in `internal/domain/levels/` and return complete `GameConfig` structs.

1. **Create a new level file** (e.g., `level2.go`):

```go
package levels

import "github.com/Kishlin/drill-game/internal/domain/config"

func GetLevel2Config() *config.GameConfig {
    return &config.GameConfig{
        World: config.WorldConfig{
            Width:       3072,
            Height:      64 * 800,
            GroundLevel: 640.0,
            Seed:        123,
            PlayerSpawn: config.PlayerSpawn{X: 1536.0, Y: 570.0},
            BuildingLayout: config.BuildingLayout{
                HospitalX:    480.0,
                FuelStationX: 850.0,
                MarketX:      1400.0,
                UpgradeShopX: 1850.0,
                ItemShopX:    2220.0,
            },
        },
        Player: config.PlayerConfig{
            StartingMoney: 500,
            StartingItems: [5]int{1, 0, 0, 0, 0}, // 1 Teleport
            StartingUpgrades: config.StartingUpgrades{
                Engine: 1, Hull: 1, // Start with Mk1
            },
        },
        Generation: config.GenerationConfig{
            // Define ores, hazards, distributions...
        },
        Upgrades: config.UpgradeConfig{
            // Define upgrade tiers...
        },
        Items: config.ItemConfig{
            // Define item prices...
        },
        Level: config.LevelConfig{Number: 2, Name: "Level 2"},
    }
}
```

2. **Register the level** in `registry.go`:

```go
func GetLevelConfig(levelNum int) (*config.GameConfig, error) {
    switch levelNum {
    case -1:
        return GetTestLevelConfig(), nil
    case 1:
        return GetLevel1Config(), nil
    case 2:
        return GetLevel2Config(), nil  // Add new case
    default:
        return nil, fmt.Errorf("level %d not found", levelNum)
    }
}
```

3. **Use the level** in `main.go`:

```go
gameCfg, err := levels.GetLevelConfig(2)
```

**Tips for Level Design:**
- Copy an existing level as a starting point
- Each level is independent—define all values explicitly
- Use different ore sets per level (level 1 might have copper/iron, level 2 might have uranium)
- Adjust upgrade prices for difficulty scaling
- Test with validation: `gameCfg.Validate()` catches common errors

---

### Configuring the World

World dimensions, building positions, player spawn, and seed are all defined in level config files (`internal/domain/levels/levelN.go`).

**Key WorldConfig Fields:**
- `Width`, `Height` — World dimensions in pixels
- `GroundLevel` — Y position of ground (align to 64px tile boundary)
- `PlayerSpawn.X`, `PlayerSpawn.Y` — Starting position
- `BuildingLayout` — Hospital, FuelStation, Market, UpgradeShop, ItemShop X positions
- `Seed` — Procedural generation seed

**Validation Rules (enforced by `WorldConfig.Validate()`):**
- Dimensions must be positive
- Player spawn must be within bounds
- Buildings must not be completely off-screen (partial off-screen is allowed)

**Design Tip:** Keep world width as a multiple of 64 (tile size) for clean alignment.

---

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

Fixed physics constants are in `internal/domain/physics/constants.go`:

```go
const (
    Gravity             = 800.0    // pixels/sec² - downward acceleration
    MoveDamping         = 1000.0   // pixels/sec² - horizontal deceleration
    FlyDamping          = 300.0    // pixels/sec² - vertical deceleration
    FallDamageThreshold = 500.0    // px/sec - minimum speed for damage
    FallDamageDivisor   = 20.0     // damage scaling factor
)
```

Dynamic movement values (MaxSpeed, Acceleration, FlyAcceleration, MaxUpwardSpeed) come from player upgrades — see [CONFIGURATION.md](docs/CONFIGURATION.md) for the Engine upgrade table.

**Process:**
1. Adjust constants in `physics/constants.go` (or upgrade tiers in level config for dynamic values)
2. Run game: `go run cmd/game/main.go`
3. Test feel and responsiveness
4. Write/update unit tests for expected behavior
5. Verify existing tests still pass: `go test ./internal/domain/physics -v`

### Adding a New Damage Source

The damage system follows a clear pattern: **Physics calculates → Player applies**.

**Pattern:**
```go
// 1. Create damage calculation in physics package
// File: internal/domain/physics/newdamage.go
package physics

func ApplyNewDamage(player *entities.Player, param float32) {
    // Calculate damage based on physics
    damage := someCalculation(param)

    // Apply through Player.DealDamage() (handles HP clamping)
    player.DealDamage(damage)
}

// 2. Call from PhysicsSystem or other systems
physics.ApplyNewDamage(player, param)

// 3. Write tests for damage calculation
// File: internal/domain/physics/newdamage_test.go
func TestApplyNewDamage_Example(t *testing.T) {
    player := NewPlayer(0, 0)
    ApplyNewDamage(player, 5.0)
    // Verify HP changed as expected
}
```

**Advantages of this pattern:**
- Physics package calculates damage (pure, testable functions)
- Player entity applies damage via `DealDamage()` (centralizes HP mutations)
- Future damage sources (gas, lava, pressure) follow the same pattern
- No duplicate HP clamping logic

**Example damage sources:**
- Fall damage: `domain/physics/fall_damage.go:ApplyFallDamage()`
- Heat damage: `domain/physics/heat.go:ApplyHeatDamage()`
- Pressure (future): `domain/physics/pressure.go:ApplyPressureDamage()`

---

## Code Review Checklist

When submitting changes, verify:

### Architecture Compliance
- [ ] `./scripts/architecture-check.sh` passes
- [ ] New domain code is pure functions (testable without framework)
- [ ] Framework integration stays in `internal/adapters/`
- [ ] Data flow is clear (domain → adapters → application)

### Testing
- [ ] All domain changes have unit tests
- [ ] Tests run successfully: `go test ./...`
- [ ] New tests follow existing patterns (pure function tests)
- [ ] Coverage maintained or improved

### Code Quality
- [ ] `./scripts/style-check.sh` passes (no `!` negation)
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
./scripts/architecture-check.sh
./scripts/style-check.sh

# Commit
git add .
git commit -m "Feature: description of change

Detailed explanation of what and why."

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
./scripts/architecture-check.sh      # Verify hexagonal boundaries
./scripts/style-check.sh             # Verify Go style rules

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
- [ARCHITECTURE.md](ARCHITECTURE.md) - High-level architecture overview
- [SYSTEMS.md](SYSTEMS.md) - Game systems implementation
- [PHYSICS.md](PHYSICS.md) - Physics and collision details
- [WORLD.md](WORLD.md) - World generation
- [BOSS.md](BOSS.md) - Boss system
- [CONFIGURATION.md](CONFIGURATION.md) - Config structs and reference tables
- [GAME_DESIGN.md](GAME_DESIGN.md) - Game mechanics and progression
- [Go Effective Guide](https://go.dev/doc/effective_go)
- [Raylib Documentation](https://www.raylib.com/)
