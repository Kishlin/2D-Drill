# Boss Fight System - Architecture Overview

## System Flow Diagram

```
Game Engine (game.go)
    ↓
    └─→ NewGame()
        ├─ Creates World (with BossRoom if configured)
        ├─ Creates Player
        ├─ Creates ItemSystem
        ├─ Creates BossFightSystem (if BossRoom exists)
        │  ├─ Calls bosses.CreateBoss() via createBossByType()
        │  └─ Links ItemSystem → BossFightSystem
        └─ Stores all systems
    ↓
    └─→ Update(dt, inputState)
        ├─ Update chunks
        ├─ Shop UI (pause)
        ├─ Physics
        ├─ Fuel
        ├─ Drilling
        ├─ Item usage (calls BossFightSystem.DamageBoss on bombs)
        ├─ Market
        ├─ Fuel Station
        ├─ Hospital
        └─→ BossFightSystem.Update()  ← NEW
            ├─ Check if player in room
            ├─ Activate/deactivate boss
            ├─ Update boss AI
            ├─ Handle projectile collisions
            ├─ Handle floor damage
            └─ Return GameState (Playing/Victory/Defeat)
```

## Component Interactions

### 1. Configuration Layer
```
GameConfig
├─ Level
│  └─ BossRoom
│     ├─ BossType: "test_boss" (string key)
│     ├─ FloorType: "concrete" or "lava"
│     ├─ RoomHeight: 720.0 pixels
│     └─ FloorHeight: 2.0 tiles (multiplied by TileSize)
└─ World
   └─ Height: Total world height (includes boss room + floor)
```

### 2. World Generation
```
ChunkGenerator
├─ groundTileY: Y coordinate of ground level
├─ bossRoomStartY: Top of boss room area
├─ bossRoomEndY: Bottom of boss room (start of floor)
├─ floorStartY: Top of floor area
└─ floorEndY: Bottom of world

GenerateTile(tileX, tileY)
├─ if tileY < groundTileY → TileTypeEmpty (sky)
├─ if tileY == groundTileY → TileTypeDirt (always)
├─ if isBossRoomTile(tileY) → TileTypeEmpty (room interior)
├─ if isFloorTile(tileY) → TileTypeFloor (solid, indestructible)
└─ else → Normal generation (dirt, ore, hazards)
```

### 3. Boss System

#### Boss Interface Hierarchy
```
Boss (interface)
├─ Update(player, dt)
├─ GetHP() float32
├─ GetMaxHP() float32
├─ IsDefeated() bool
├─ IsActive() bool
├─ Activate()
├─ Deactivate()
└─ GetProjectiles() []*Projectile
    │
    └─→ PhysicalBoss (extends Boss)
        ├─ GetAABB() types.AABB
        └─ TakeDamage(damage float32)
```

#### TestBoss Implementation
```
TestBoss struct
├─ aabb: Positioned at center of boss room
├─ hp: Current health (starts at 100)
├─ active: Only updates when true
├─ defeated: Set to true when hp <= 0
└─ projectiles: Empty for test boss (could be used by other bosses)

Constants
├─ MaxHP: 100.0
├─ Width: 200.0
├─ Height: 200.0
├─ DamagePerBomb: 10.0
└─ DamagePerBigBomb: 25.0
```

### 4. Boss Fight System
```
BossFightSystem
├─ boss: Reference to Boss interface
├─ bossRoomStartY: Top of boss room (pixels)
├─ bossRoomEndY: Bottom of boss room (pixels)
├─ floorStartY: Top of floor (pixels)
├─ floorEndY: Bottom of floor (pixels)
├─ floorType: "concrete" or "lava"
└─ wasPlayerInRoom: Tracks entry/exit

Update(player, dt) → GameState
├─ Check if player center is in room bounds
├─ Activate boss on entry, deactivate on exit
├─ Call boss.Update(player, dt)
├─ Check projectile-player collisions
├─ Apply floor damage if standing on lava
└─ Return GameState (Playing/Victory/Defeat)
```

### 5. Item System Integration
```
ItemSystem
├─ New field: bossFightSystem *BossFightSystem
├─ New method: SetBossFightSystem(bfs)
└─ Modified applyBomb():
   ├─ Destroy tiles (existing)
   └─ NEW: Check bomb-boss collision
       ├─ Calculate blast AABB (circular)
       ├─ Check PhysicalBoss.GetAABB().Intersects(bombAABB)
       ├─ Call bossFightSystem.DamageBoss(damage)
       │  ├─ Regular bomb: 10 HP
       │  └─ Big bomb: 25 HP
       └─ Deactivate projectile if hit
```

### 6. Game Engine Integration
```
Game struct
├─ ... existing fields ...
├─ boss: bosses.Boss (nil if no boss room)
├─ bossFightSystem: *systems.BossFightSystem (nil if no boss)
└─ gameState: entities.GameState (Playing/Victory/Defeat)

NewGame() initialization order:
1. Create world (with BossRoom if configured)
2. Create item system
3. If level has BossRoom:
   - Call createBossByType(bossType)
   - Create BossFightSystem
   - Link ItemSystem → BossFightSystem
4. Return Game with all systems

Update() call order:
1. Physics (includes fall/heat damage)
2. Fuel consumption
3. Drilling
4. Item usage (bomb-boss collision happens here)
5. Market/Fuel Station/Hospital
6. BossFightSystem.Update() ← LAST SYSTEM
   - Updates gameState
   - Checks victory/defeat conditions
```

## Data Flow: Bomb Damage to Boss

```
Player Input: UseBomb (B key)
    ↓
ItemSystem.ProcessItemUsage()
    ↓
ItemSystem.applyBomb(player, radius)
    ├─ Destroy tiles (existing)
    └─ NEW: Check bomb-boss collision
        ├─ Get boss from bossFightSystem
        ├─ Cast to PhysicalBoss
        ├─ Create bomb AABB (circular blast)
        ├─ Check boss.GetAABB().Intersects(bombAABB)
        ├─ bossFightSystem.DamageBoss(damage)
        │  └─ physicalBoss.TakeDamage(damage)
        │     └─ TestBoss.hp -= damage
        │        └─ if hp <= 0: defeated = true
        └─ Deactivate projectile (if hit)
    ↓
Next frame: BossFightSystem.Update()
    └─ if boss.IsDefeated() → return GameStateVictory
```

## Data Flow: Player Entry to Boss Room

```
Player moves down (OnGround triggers boss room entry)
    ↓
BossFightSystem.Update()
    ├─ playerCenterY = player.AABB.Y + player.AABB.Height/2
    ├─ if playerCenterY >= bossRoomStartY AND playerCenterY < bossRoomEndY
    │  └─ playerInRoom = true
    ├─ if playerInRoom && !wasPlayerInRoom
    │  └─ boss.Activate()
    │     └─ TestBoss.active = true
    ├─ wasPlayerInRoom = playerInRoom
    └─ boss.Update(player, dt)
```

## Tile Generation at Different Depths

```
World Height: 1488 pixels
Ground Level: 128 pixels
Boss Room: From Y=448 to Y=1168 (720px)
Floor: From Y=1168 to Y=1488 (320px / 5 tiles)

Normal depth ranges:
Y=128-448: Normal terrain (10 tiles)
Y=448-1168: Boss room (empty - no tiles)
Y=1168-1488: Floor (solid, indestructible)
```

## Game State Machine

```
GameStatePlaying (initial)
    ├─ Normal gameplay
    ├─ Boss can be active or inactive
    ├─ BossFightSystem.Update() called every frame
    └─ If boss.IsDefeated() → GameStateVictory
       If player.HP <= 0 → GameStateDefeat

GameStateVictory (terminal)
    ├─ Boss defeated
    ├─ No further gameplay updates
    └─ Rendering layer shows victory screen

GameStateDefeat (terminal)
    ├─ Player killed
    ├─ No further gameplay updates
    └─ Rendering layer shows defeat screen
```

## File Dependencies

```
Main Dependencies:

bosses/boss.go ← bosses/projectile.go
    ↑
bosses/test_boss/boss.go (implements Boss, PhysicalBoss)
    ↑
engine/game.go (createBossByType factory, uses Boss interface)
    ↑
systems/boss_fight.go (uses Boss, PhysicalBoss interfaces)
    ↑
systems/item.go (uses BossFightSystem)
    ↑
world/generator.go (uses config.BossRoomConfig)
    ↑
config/boss_room_config.go + config/game_config.go
    ↑
levels/level_boss.go (creates BossRoomConfig in LevelConfig)
```

## Hexagonal Architecture Compliance

**Domain Layer (internal/domain/):**
- ✅ No Raylib imports
- ✅ Pure business logic
- ✅ Interfaces for extensibility (Boss, PhysicalBoss)
- ✅ Config-driven behavior

**Adapter Layer (internal/adapters/):**
- Rendering (NOT YET IMPLEMENTED)
- Input handling (unchanged)
- Framework integration (Raylib)

## Extensibility Points

### 1. Adding New Boss Types
```go
// In internal/domain/bosses/my_boss/boss.go
type MyBoss struct { ... }
func New(roomStartY, worldWidth float32) *MyBoss { ... }
func (b *MyBoss) Update(player *entities.Player, dt float32) { ... }
// Implement Boss or PhysicalBoss interface

// In engine/game.go createBossByType()
case "my_boss":
    return my_boss.New(roomStartY, worldWidth), nil
```

### 2. Boss with Projectiles
```go
// In boss Update() method:
projectile := bosses.NewProjectile(
    x, y, width, height,
    types.Vec2{X: vx, Y: vy},
    damage,
)
b.projectiles = append(b.projectiles, projectile)

// In BossFightSystem.Update():
// Already handles projectile collision detection
// and calls player.DealDamage()
```

### 3. Different Floor Types
```go
// Already supported in config
BossRoom: &config.BossRoomConfig{
    FloorType: config.FloorLava,  // vs FloorConcrete
    ...
}

// In BossFightSystem.Update():
if s.floorType == config.FloorLava {
    // Apply damage to standing player
}
```

## Key Design Patterns

### 1. Strategy Pattern (Boss Types)
- Boss interface is strategy
- Each boss implementation is a concrete strategy
- Engine selects strategy based on config string

### 2. Composite Pattern (Game Systems)
- Game aggregates all systems
- Each system is independent
- BossFightSystem is one system among many

### 3. Observer Pattern (State Changes)
- Player health changes trigger damage
- Boss health changes trigger defeat condition
- GameState tracks win/lose conditions

### 4. Factory Pattern (Boss Creation)
- createBossByType() in engine/game.go
- Takes config string, returns Boss interface
- Avoids circular imports by living in engine layer

## Performance Considerations

1. **Boss Rendering (pending):**
   - Single rectangle draw
   - One HP bar draw
   - Projectiles: typically 0-5 per frame (low count)
   - Floor tiles: rendered by existing tile loop

2. **Collision Detection:**
   - Bomb-boss: One AABB.Intersects() call per bomb
   - Projectile-player: Few projectiles, simple AABB checks
   - Floor-player: Simple Y-coordinate comparison

3. **Memory:**
   - One boss instance per game
   - Projectiles: Reusable, deactivated instead of deleted
   - No significant overhead

## Testing

All 145 tests pass:
- Config validation tests ✅
- World generation tests ✅
- Physics tests ✅
- Systems tests ✅
- Drilling tests ✅
- No boss-specific tests yet (domain logic only)

To add boss tests:
```go
// internal/domain/bosses/test_boss/boss_test.go
func TestTestBoss_TakeDamage(t *testing.T) { ... }
func TestTestBoss_Defeated(t *testing.T) { ... }
```

## Next Phase: Rendering

The rendering layer will:
1. Query game state via Game.GetBoss(), Game.GetGameState()
2. Render boss AABB as rectangle (or sprite texture)
3. Render HP bar at top of screen
4. Render projectiles in world space
5. Style floor tiles based on FloorType
6. Display victory/defeat screens based on GameState

No domain changes needed for rendering - it's purely visual integration with Raylib.
