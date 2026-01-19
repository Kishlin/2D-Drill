# Boss Fight System Implementation Plan (Original)

## Overview
Add boss fights at the bottom of each level. Players dig through tiles, enter a boss room, and defeat the boss to win the level. Bosses are implemented as separate packages with common interfaces.

## Key Design Decisions
- **Boss Packages**: Each boss type gets its own package (`internal/domain/bosses/<boss_name>/`) with custom logic
- **Common Interfaces**: `Boss` interface for all bosses, `PhysicalBoss` for bosses with AABB that can be hit
- **Floor Types**: Configurable floor behavior (concrete, lava) instead of just visual color
- **Projectiles**: Shared projectile entity for spell-casting and bullet-hell bosses
- **Config Simplicity**: Level config only specifies boss type string; boss package owns its constants

---

## Implementation Steps

### 1. Configuration Layer

**New file: `internal/domain/config/boss_room_config.go`**
```go
type FloorType string
const (
    FloorConcrete FloorType = "concrete"  // Solid, walkable
    FloorLava     FloorType = "lava"      // Deals damage on contact
)

type BossRoomConfig struct {
    BossType    string     // e.g., "earth_guardian", "bullet_hell_1"
    FloorType   FloorType
    RoomHeight  float32    // Boss room height in pixels (~720)
    FloorHeight float32    // Floor tiles below room
}
```

**Modify: `internal/domain/config/game_config.go`**
- Add `BossRoom *BossRoomConfig` to `LevelConfig` (replace placeholder comment)
- Add validation: if BossRoom != nil, validate RoomHeight > 0, FloorHeight >= 1

### 2. Boss Package Structure

**New directory: `internal/domain/bosses/`**

**New file: `internal/domain/bosses/boss.go`** (common interfaces)
```go
type Boss interface {
    Update(player *entities.Player, dt float32)
    GetHP() float32
    GetMaxHP() float32
    IsDefeated() bool
    IsActive() bool
    Activate()
    Deactivate()
    GetProjectiles() []*Projectile  // Empty if none
}

// Physical bosses can be hit by bombs
type PhysicalBoss interface {
    Boss
    GetAABB() types.AABB
    TakeDamage(damage float32)
}
```

**New file: `internal/domain/bosses/projectile.go`**
```go
type Projectile struct {
    AABB     types.AABB
    Velocity types.Vec2
    Damage   float32
    Active   bool
}

func (p *Projectile) Update(dt float32)
func (p *Projectile) Intersects(playerAABB types.AABB) bool
```

**New file: `internal/domain/bosses/registry.go`**
```go
func CreateBoss(bossType string, roomStartY, worldWidth float32) (Boss, error) {
    switch bossType {
    case "test_boss":
        return test_boss.New(roomStartY, worldWidth), nil
    // Future bosses added here
    default:
        return nil, fmt.Errorf("unknown boss type: %s", bossType)
    }
}
```

**New directory: `internal/domain/bosses/test_boss/`** (simple physical boss for testing)
```go
// boss.go
const (
    MaxHP            = 100.0
    Width            = 200.0
    Height           = 200.0
    DamagePerBomb    = 10.0
    DamagePerBigBomb = 25.0
)

type TestBoss struct {
    aabb        types.AABB
    hp          float32
    active      bool
    defeated    bool
    projectiles []*bosses.Projectile
}

func New(roomStartY, worldWidth float32) *TestBoss
func (b *TestBoss) Update(player *entities.Player, dt float32)
func (b *TestBoss) TakeDamage(damage float32)
// ... implement Boss and PhysicalBoss interfaces
```

### 3. Entity Layer

**New file: `internal/domain/entities/game_state.go`**
```go
type GameState int
const (
    GameStatePlaying GameState = iota
    GameStateVictory
    GameStateDefeat
)
```

**Modify: `internal/domain/entities/tile.go`**
- Add `TileTypeFloor` (solid, not drillable, not nukeable)
- Update `IsDrillable()` to return false for floor

### 4. World Layer

**Modify: `internal/domain/world/generator.go`**
- Store `bossRoomConfig *config.BossRoomConfig` reference
- In `GenerateTile()`:
  - Return empty tile for boss room depth
  - Return floor tile for floor depth
- Helper methods: `isBossRoomTile(tileY)`, `isFloorTile(tileY)`

**Modify: `internal/domain/world/world.go`**
- Guard `NukeTileAtGrid()`: don't nuke floor tiles

### 5. System Layer

**New file: `internal/domain/systems/boss_fight.go`**
```go
type BossFightSystem struct {
    boss           bosses.Boss
    bossRoomStartY float32
    bossRoomEndY   float32
    floorType      config.FloorType
}

func NewBossFightSystem(boss bosses.Boss, bossRoomCfg *config.BossRoomConfig, worldHeight float32) *BossFightSystem

func (s *BossFightSystem) Update(player *entities.Player, inputState input.InputState, dt float32) entities.GameState {
    // Check if player in boss room -> activate/deactivate boss
    // Call boss.Update(player, dt)
    // Process projectile-player collisions -> player.DealDamage()
    // Apply floor damage if FloorLava and player touching floor
    // Check victory condition
}

func (s *BossFightSystem) DamageBoss(damage float32) {
    if physical, ok := s.boss.(bosses.PhysicalBoss); ok {
        physical.TakeDamage(damage)
    }
}

func (s *BossFightSystem) IsPlayerInBossRoom(player *entities.Player) bool
func (s *BossFightSystem) GetBoss() bosses.Boss
```

**Modify: `internal/domain/systems/item.go`**
- Add `bossFightSystem *BossFightSystem` field and setter
- In `applyBomb()`: if in boss room and boss is PhysicalBoss, check bomb-boss collision and call `DamageBoss()`

### 6. Game Loop Integration

**Modify: `internal/domain/engine/game.go`**
- Add fields: `boss bosses.Boss`, `bossFightSystem *systems.BossFightSystem`, `gameState entities.GameState`
- In `NewGame()`:
  - If config has BossRoom, call `bosses.CreateBoss()` and create system
  - Pass boss system reference to item system
- In `Update()`:
  - After physics, call `bossFightSystem.Update()` (handles boss AI, projectiles, floor damage)
  - Check returned GameState for victory
- New getters: `GetBoss()`, `GetGameState()`, `IsBossFightActive()`

### 7. Rendering Layer

**Modify: `internal/adapters/rendering/raylib.go`**

```go
func (r *RaylibRenderer) renderBoss(boss bosses.Boss) {
    // Render projectiles
    for _, proj := range boss.GetProjectiles() {
        if proj.Active {
            rl.DrawRectangle(...) // Projectile visual
        }
    }

    // Render physical boss model if applicable
    if physical, ok := boss.(bosses.PhysicalBoss); ok {
        aabb := physical.GetAABB()
        rl.DrawRectangle(...) // Boss rectangle (future: sprite)
    }
}

func (r *RaylibRenderer) renderBossHPBar(boss bosses.Boss) {
    // Centered at top of screen
    // Shows boss.GetHP() / boss.GetMaxHP()
}

func (r *RaylibRenderer) renderFloor(bossRoomCfg *config.BossRoomConfig, worldHeight float32) {
    // Draw floor based on FloorType
    // Concrete: gray
    // Lava: orange/red with possible animation
}

func (r *RaylibRenderer) renderVictoryScreen()
```

### 8. Boss Test Level

**New file: `internal/domain/levels/level_boss_test.go`**
```go
func GetBossTestLevelConfig() *config.GameConfig {
    cfg := GetTestLevelConfig()  // Start with dev player stats

    // Minimal world: 10 tiles to dig + boss room + floor
    cfg.World.Height = 1488  // ~23 tiles total
    cfg.World.GroundLevel = 128
    cfg.World.PlayerSpawn.Y = 58

    cfg.Level = config.LevelConfig{
        Number: -2,
        Name:   "Boss Test Level",
        BossRoom: &config.BossRoomConfig{
            BossType:    "test_boss",
            FloorType:   config.FloorConcrete,
            RoomHeight:  720.0,
            FloorHeight: 2.0,
        },
    }
    return cfg
}
```

**Modify: `internal/domain/levels/registry.go`**
- Add case `-2` returning `GetBossTestLevelConfig()`

---

## Files Summary

### New Files (9)
| File | Purpose |
|------|---------|
| `internal/domain/config/boss_room_config.go` | BossRoomConfig, FloorType |
| `internal/domain/entities/game_state.go` | GameState enum |
| `internal/domain/bosses/boss.go` | Boss, PhysicalBoss interfaces |
| `internal/domain/bosses/projectile.go` | Shared Projectile entity |
| `internal/domain/bosses/registry.go` | Boss factory by type string |
| `internal/domain/bosses/test_boss/boss.go` | Simple physical boss for testing |
| `internal/domain/systems/boss_fight.go` | Boss fight orchestration |
| `internal/domain/levels/level_boss_test.go` | Quick test level (-2) |

### Modified Files (6)
| File | Changes |
|------|---------|
| `internal/domain/config/game_config.go` | Add BossRoom to LevelConfig |
| `internal/domain/entities/tile.go` | Add TileTypeFloor |
| `internal/domain/world/generator.go` | Generate boss room and floor |
| `internal/domain/world/world.go` | Guard NukeTileAtGrid for floor |
| `internal/domain/systems/item.go` | Boss damage from bombs |
| `internal/domain/engine/game.go` | Integrate boss system |
| `internal/domain/levels/registry.go` | Add level -2 |
| `internal/adapters/rendering/raylib.go` | Boss, projectiles, HP bar, floor, victory |

---

## Future Boss Examples

Once the foundation is in place, adding new bosses follows this pattern:

**Bullet Hell Boss** (`internal/domain/bosses/bullet_hell_1/`):
- No AABB (doesn't implement PhysicalBoss)
- HP depletes over time (survival mode)
- `Update()` spawns projectile patterns
- Player must dodge for X seconds to win

**Spell Caster Boss** (`internal/domain/bosses/spell_caster/`):
- Has AABB (implements PhysicalBoss)
- Cycles through spell phases
- Each spell spawns different projectile patterns
- Vulnerable windows between spells

---

## Verification

1. **Build**: `go build ./...`
2. **Tests**: `go test ./...`
3. **Run boss test level**: `go run cmd/game/main.go -level -2`
4. **Manual testing**:
   - Dig down ~10 tiles to reach boss room
   - Verify boss room is empty (no tiles)
   - Verify floor is solid (can't drill or nuke)
   - Verify HP bar appears when entering boss room
   - Use bombs (B) and big bombs (G) to damage boss
   - Verify HP bar updates
   - Teleport out (T) - HP bar disappears
   - Return - HP bar shows same HP value
   - Defeat boss - victory screen appears
5. **Domain constraint**: `grep -r "raylib" internal/domain/` returns empty
