# Claude Skills for Drill Game

This directory contains custom skills for Claude Code to assist with game development workflows.

## Available Skills

### 🧪 test-runner
**Purpose**: Run Go tests with coverage reporting  
**Use cases**:
- Quick test runs after development
- Check domain layer test coverage
- Verify physics tests pass

**Parameters**:
- `scope`: Test scope (`all`, `domain`, `physics`)
- `coverage`: Generate coverage report (default: `true`)
- `verbose`: Show all test names (default: `false`)

**Examples**:
```
"Run all tests"
"Run physics tests with coverage"
"Run domain tests in verbose mode"
```

---

### 🔨 build-game
**Purpose**: Build the game binary with optional execution  
**Use cases**:
- Quick iteration: build and test changes
- Create optimized release builds
- Verify compilation succeeds

**Parameters**:
- `run`: Run the game after building (default: `false`)
- `optimize`: Build with optimizations (default: `false`)

**Examples**:
```
"Build the game"
"Build and run the game"
"Build optimized game"
```

---

### 🏗️ architecture-check
**Purpose**: Verify hexagonal architecture compliance  
**Use cases**:
- Ensure domain layer remains framework-independent
- Catch architecture violations early
- Maintain clean separation of concerns

**Checks**:
1. ✅ Domain layer has no Raylib imports
2. ✅ Domain doesn't import adapters (wrong dependency direction)
3. ✅ Domain uses custom types (not Raylib types like `rl.Vector2`)

**Examples**:
```
"Check architecture compliance"
"Verify clean architecture"
```

---

### 📊 profile-game
**Purpose**: Performance profiling for optimization  
**Use cases**:
- Identify CPU bottlenecks in game loop
- Analyze memory allocation patterns
- Optimize frame rendering

**Parameters**:
- `type`: Profile type (`cpu`, `memory`)
- `duration`: How long to profile (default: `30s`)

**Note**: This skill provides guidance on adding profiling to the game. For full profiling, you'll need to add pprof hooks to `main.go`.

**Examples**:
```
"Profile CPU usage"
"Profile memory for 1 minute"
```

---

## Why These Skills?

Based on the Drill Game's architecture and development needs:

1. **test-runner**: Your project has comprehensive physics tests (11+ tests). This skill makes TDD workflows faster and shows coverage at a glance.

2. **build-game**: Quick build/run cycles are essential for game dev. This skill streamlines the build process and provides feedback on binary size.

3. **architecture-check**: Your hexagonal architecture is a key design principle. This skill acts as a guard rail to prevent violations (like accidentally importing Raylib in domain code).

4. **profile-game**: Game performance is critical. This skill guides you through profiling setup and analysis, crucial for optimizing physics and rendering.

## How to Use

Claude Code will automatically discover these skills. Simply ask Claude naturally:

- "Run the tests" → Uses `test-runner`
- "Build and run" → Uses `build-game`
- "Check if architecture is clean" → Uses `architecture-check`
- "Profile the game's performance" → Uses `profile-game`

## Adding More Skills

To add new skills:

1. Create a new JSON file in `.claude/skills/`
2. Follow the existing skill format
3. Add permission to `.claude/settings.local.json`:
   ```json
   "Skill(your-skill-name:*)"
   ```
4. Document it in this README

## Skill Ideas for Future

- **asset-validator**: Check for missing/invalid assets
- **benchmark-runner**: Run performance benchmarks
- **entity-generator**: Scaffold new game entities
- **system-generator**: Create new game systems following patterns
- **integration-test**: Run full integration tests with mocked Raylib

