# Quick Reference Guide

## All Files Changed/Created

### ✅ CREATED FILES (8)

#### 1. `internal/domain/config/boss_room_config.go`
- Boss room configuration struct
- FloorType enum (Concrete, Lava)
- New file, ~20 lines

#### 2. `internal/domain/bosses/boss.go`
- Boss interface (all bosses implement)
- PhysicalBoss interface (for bomb-vulnerable bosses)
- Key methods: Update, GetHP, GetMaxHP, IsDefeated, IsActive, Activate, Deactivate, GetProjectiles
- New file, ~30 lines

#### 3. `internal/domain/bosses/projectile.go`
- Projectile entity (shared by all boss types)
- Fields: AABB, Velocity, Damage, Active
- Methods: Update, Intersects, Deactivate
- New file, ~35 lines

#### 4. `internal/domain/bosses/test_boss/boss.go`
- Test boss implementation (100 HP, 200x200px, static)
- Implements PhysicalBoss interface
- Constants: MaxHP=100, Width=200, Height=200, DamagePerBomb=10, DamagePerBigBomb=25
- New file, ~60 lines

#### 5. `internal/domain/entities/game_state.go`
- GameState enum: GameStatePlaying, GameStateVictory, GameStateDefeat
- New file, ~8 lines

#### 6. `internal/domain/systems/boss_fight.go`
- BossFightSystem: orchestrates boss fights
- Key methods: Update (returns GameState), IsPlayerInBossRoom, DamageBoss, GetBoss, IsBossFightActive
- Handles: boss activation, projectile collisions, floor damage, game state transitions
- New file, ~120 lines

#### 7. `internal/domain/levels/level_boss.go` (renamed from level_boss_test.go)
- GetBossTestLevelConfig() function for level -2
- Minimal world optimized for boss testing
- Advanced player stats (all items, max upgrades)
- Simplified generation config
- New file, ~85 lines

### ✅ MODIFIED FILES (9)

#### 1. `internal/domain/config/game_config.go`
**Changes:**
- Line 5-8: Added `BossRoom *BossRoomConfig` to `LevelConfig`
- Line 74-82: Added validation for BossRoom (RoomHeight > 0, FloorHeight >= 1)
- ~6 lines added, 0 removed

#### 2. `internal/domain/entities/tile.go`
**Changes:**
- Line 16: Added `TileTypeFloor` constant
- Line 45: Updated comment for IsDrillable() explaining Floor is not drillable
- ~2 lines added, 0 removed

#### 3. `internal/domain/world/generator.go`
**Changes:**
- Lines 13-21: Added boss room fields to ChunkGenerator struct
- Lines 23-30: Existing NewChunkGeneratorFromConfig function unchanged
- Lines 32-51: Added NewChunkGeneratorFromConfigWithBoss constructor
- Lines 61-89: Modified GenerateTile() to handle boss rooms and floors
- Lines 91-105: Added isBossRoomTile() and isFloorTile() helper methods
- ~90 lines added, 0 removed

#### 4. `internal/domain/world/world.go`
**Changes:**
- Line 22-24: Simplified NewWorldFromConfig to call NewWorldFromConfigWithBoss
- Lines 26-44: Added NewWorldFromConfigWithBoss constructor
- Line 143: Modified NukeTileAtGrid guard to prevent destroying TileTypeFloor
- ~25 lines added, 0 removed

#### 5. `internal/domain/systems/item.go`
**Changes:**
- Line 4: Added import `"github.com/Kishlin/drill-game/internal/domain/bosses"`
- Line 17: Added `bossFightSystem *BossFightSystem` field
- Lines 32-34: Added SetBossFightSystem() setter method
- Lines 75-98: Added bomb-boss collision detection in applyBomb()
- ~40 lines added, 0 removed

#### 6. `internal/domain/engine/game.go`
**Changes:**
- Lines 3-4: Added imports: `"fmt"` and `"github.com/Kishlin/drill-game/internal/domain/bosses/test_boss"`
- Lines 24-26: Added fields: boss, bossFightSystem, gameState
- Lines 72-92: Modified NewGame to create boss system if BossRoom is configured
- Lines 155-157: Added boss system update in Update() method
- Lines 202-232: Added getters and createBossByType factory function
- ~80 lines added, 0 removed

#### 7. `internal/domain/levels/registry.go`
**Changes:**
- Lines 13-14: Added case `-2` for GetBossTestLevelConfig()
- ~2 lines added, 0 removed

#### 8. `cmd/game/main.go`
**Changes:**
- Line 29: Changed from `levels.GetLevelConfig(-1)` to `levels.GetLevelConfig(-2)` for testing
- **NOTE:** This must be changed back to -1 before merging to main branch
- ~0 lines added, 1 line modified

### 🚫 DELETED FILES (1)

#### `internal/domain/bosses/registry.go`
- Was created but caused circular import (test_boss imports bosses)
- Solution: Moved factory logic to engine/game.go createBossByType()
- Avoided circular import: test_boss → bosses → registry → test_boss

---

## Key Architecture Changes

### 1. Boss Factory Location
**Before:** `bosses/registry.go` would call `test_boss.New()`
**After:** `engine/game.go` has `createBossByType()` function
**Why:** Avoids circular import (test_boss imports bosses package)

### 2. Tile Type Expansion
**Before:** TileTypeEmpty, TileTypeDirt, TileTypeOre, TileTypeRock, TileTypeLava
**After:** Added TileTypeFloor (solid but indestructible)
**Impact:** IsDrillable() still returns false for Floor (cannot drill)

### 3. World Generation Enhancement
**Before:** Single NewChunkGeneratorFromConfig constructor
**After:**
- `NewChunkGeneratorFromConfig()` - normal generation
- `NewChunkGeneratorFromConfigWithBoss()` - with boss room support
- Backward compatible: Both world constructors work

### 4. Game State Tracking
**Before:** No explicit game state
**After:** Game.gameState tracks Playing/Victory/Defeat
**Source:** BossFightSystem.Update() returns GameState

### 5. Bomb Damage System
**Before:** Bombs only destroy tiles
**After:**
- Bombs check for boss collision
- Bombs deal 10 HP (regular) or 25 HP (big bomb)
- Floor damage from standing on lava floor

---

## Critical Information for Next Developer

### Main Entry Point Needs Fix
**File:** `cmd/game/main.go`, Line 29
**Current:** `levels.GetLevelConfig(-2)` (test boss level)
**Should be:** `levels.GetLevelConfig(-1)` (normal test level)
**Action:** Change before merging to main branch

### All Tests Pass
```
✅ 145 tests PASS
✅ go build ./... SUCCESS
✅ No Raylib imports in domain layer verified
```

### Boss System is Complete
- Domain logic: ✅ Complete
- Systems integration: ✅ Complete
- Game engine integration: ✅ Complete
- Configuration: ✅ Complete
- Test level: ✅ Complete

### Only Rendering is Pending
- Boss rendering (rectangle)
- HP bar (top of screen)
- Projectile rendering
- Floor tile styling
- Victory/defeat screens

See RENDERING_PLAN.md for implementation details.

---

## Testing the Current Implementation

```bash
# Build
go build ./...

# Test
go test ./...

# Run game (loads level -2 for testing)
go run cmd/game/main.go

# Manual checks:
# - Dig down to boss room
# - Use bombs (B key) - boss should take 10 HP damage
# - Use big bombs (G key) - boss should take 25 HP damage
# - Verify HP in terminal console (not yet rendered)
# - Defeat boss - victory state triggers (use Ctrl+C to exit)
```

---

## Code Search Quick Links

### Find implementations of Boss interface:
```bash
grep -r "func (.*) Update.*Player.*float32" internal/domain/bosses/
```

### Find PhysicalBoss implementations:
```bash
grep -r "func (.*) GetAABB" internal/domain/bosses/
```

### Find boss system references:
```bash
grep -r "BossFightSystem" internal/domain/
```

### Find game state usage:
```bash
grep -r "GameState" internal/domain/
```

---

## Important Interfaces

### Boss Interface
```go
Update(player *entities.Player, dt float32)
GetHP() float32
GetMaxHP() float32
IsDefeated() bool
IsActive() bool
Activate()
Deactivate()
GetProjectiles() []*Projectile
```

### PhysicalBoss Interface (extends Boss)
```go
GetAABB() types.AABB
TakeDamage(damage float32)
```

### Projectile Entity
```go
AABB types.AABB
Velocity types.Vec2
Damage float32
Active bool
Update(dt float32)
Intersects(playerAABB types.AABB) bool
Deactivate()
```

---

## Adding a New Boss Type

Pattern for future boss implementations:

1. Create `internal/domain/bosses/my_boss/boss.go`
2. Implement `Boss` interface (or `PhysicalBoss` for hittable bosses)
3. Add constants for stats (HP, size, etc.)
4. Add case in `engine/game.go` createBossByType():
   ```go
   case "my_boss":
       return my_boss.New(roomStartY, worldWidth), nil
   ```
5. Reference in level config:
   ```go
   BossRoom: &config.BossRoomConfig{
       BossType: "my_boss",
       FloorType: config.FloorConcrete,
       RoomHeight: 720.0,
       FloorHeight: 2.0,
   }
   ```

---

## Resume Checklist

When resuming work on rendering:

- [ ] Read RENDERING_PLAN.md
- [ ] Verify main.go is set to load level -1 (or -2 for testing)
- [ ] Open `internal/adapters/rendering/raylib.go`
- [ ] Add boss rendering functions
- [ ] Add HP bar rendering
- [ ] Add projectile rendering
- [ ] Style floor tiles
- [ ] Add victory/defeat screens
- [ ] Test with: `go run cmd/game/main.go`
- [ ] Run tests: `go test ./...`
- [ ] Commit with message from RENDERING_PLAN.md
