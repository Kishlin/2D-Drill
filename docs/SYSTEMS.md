# Game Systems

This document covers the domain-layer systems that orchestrate game logic. For high-level architecture, see [ARCHITECTURE.md](ARCHITECTURE.md). For physics details, see [PHYSICS.md](PHYSICS.md).

---

## Overview

Systems are stateful orchestrators that coordinate domain logic each frame. They live in `internal/domain/systems/` and operate on entities without any framework dependencies.

**System Update Order** (in `engine/game.go`):
1. Chunk loading (proactive around player)
2. UI Manager (modal pause if active)
3. Building interactions (E key detection)
4. Physics (movement, collision, damage)
5. Fuel consumption
6. Drilling animation
7. Item usage
8. Boss fight

---

## Game Orchestration (`engine/game.go`)

The `Game` struct orchestrates all systems:

```go
type Game struct {
    world           *world.World
    player          *entities.Player
    buildings       []*entities.Building
    uiManager       *ui.Manager
    effectProcessor *effects.Processor
    effectContext   *effects.EffectContext
    damageables     []effects.DamageableEntity
    physicsSystem   *systems.PhysicsSystem
    drillingSystem  *systems.DrillingSystem
    fuelSystem      *systems.FuelSystem
    bossFightSystem *systems.BossFightSystem
}
```

**Update Flow:**

```go
func (g *Game) Update(dt float32, inputState input.InputState) error {
    // 0. Update chunks around player (proactive loading)
    g.world.UpdateChunksAroundPlayer(playerX, playerY)

    // 1. Process active UI (modal pause if open)
    if g.uiManager.HasActiveUI() {
        result := g.uiManager.Process(g.player, inputState)
        g.effectProcessor.Apply(g.effectContext, result.Effects)
        if g.uiManager.HasActiveUI() {
            return nil  // Modal UI open - pause gameplay
        }
    }

    // 2. Detect building interactions (E key + overlap)
    if interactionType := systems.DetectInteraction(...); interactionType != nil {
        g.uiManager.OpenUI(*interactionType)
        // ... process and apply effects via effectContext
    }

    // 3. Physics - handles movement, collision, fall/heat damage
    g.physicsSystem.UpdatePhysics(g.player, inputState, dt)

    // 4. Fuel consumption (runs even during drilling)
    g.fuelSystem.ConsumeFuel(g.player, inputState, dt)

    // 5. Drilling animation
    g.drillingSystem.ProcessDrilling(g.player, inputState, dt)
    if g.player.IsDrilling {
        return nil  // Skip interactions during drill
    }

    // 6. Item usage (returns effects, applied via effectContext)
    itemEffects := systems.DetectItemUsage(g.player, inputState, g.config.Items)
    if len(itemEffects) > 0 {
        g.effectProcessor.Apply(g.effectContext, itemEffects)
    }

    // 7. Boss fight system
    if g.bossFightSystem != nil {
        g.gameState = g.bossFightSystem.Update(g.player, dt)
    }

    return nil
}
```

---

## Drilling System (`systems/drilling.go`)

Handles both vertical and horizontal drilling with variable animation duration based on tile hardness and depth.

### Duration Formula

```
duration = baseTime × hardness × depthFactor / drillSpeed
```

Where:
- `baseTime`: 1.0 second constant
- `hardness`: per-tile-type value (Dirt 1.0, Copper 1.2, Diamond 3.0, Lava 0.3)
- `depthFactor`: scales 1.0 at surface → 24.0 at max depth
- `drillSpeed`: from drill upgrades (1.0 → 6.0)

**Lava Exception:** Fixed 0.3s duration regardless of depth (damage is the penalty).

### Animation State

```go
type DrillingSystem struct {
    world     *world.World
    animation DrillingAnimation
}

type DrillingAnimation struct {
    Active      bool
    Direction   DrillDirection // Down, Left, or Right
    StartX      float32        // Player position when animation started
    StartY      float32
    TargetX     float32        // Where player moves to during animation
    TargetY     float32
    TargetGridX int            // Tile coordinates for removal
    TargetGridY int
    Elapsed     float32        // Time elapsed in animation
    Duration    float32        // Variable duration (1.0-24+ seconds)
    Tile        *entities.Tile
}
```

### Vertical Drilling (S/Down Key)

```go
// Check tile directly below player's center
playerCenterX := player.AABB.X + player.AABB.Width/2
playerBottomY := player.AABB.Y + player.AABB.Height
tile := ds.world.GetTileAt(playerCenterX, playerBottomY)

// Calculate animation targets
tileCenterX := float32(tileGridX)*world.TileSize + world.TileSize/2
targetX := tileCenterX - player.AABB.Width/2

tileBottomY := float32(tileGridY+1) * world.TileSize
targetY := tileBottomY - player.AABB.Height  // Align bottom edges

// Start variable-duration animation
ds.startDrillAnimation(player, DrillDown, tileGridX, tileGridY, targetX, targetY, tile)
```

### Horizontal Drilling (Left/Right When Grounded)

```go
// Check tile beside player
playerCenterY := player.AABB.Y + player.AABB.Height/2

if inputState.Left {
    tile := ds.world.GetTileAt(player.AABB.X - 1, playerCenterY)
    if tile != nil && tile.IsDrillable() {
        tileCenterX := float32(tileGridX)*world.TileSize + world.TileSize/2
        targetX := tileCenterX - player.AABB.Width/2
        targetY := player.AABB.Y  // Stay at ground level

        ds.startDrillAnimation(player, DrillLeft, tileGridX, tileGridY, targetX, targetY, tile)
    }
}
```

### Lava Tile Drilling

Lava tiles are special hazards with fixed duration and damage on completion:

```go
// Lava uses fixed duration from HazardHardness (0.3s), depth-independent
if tile.Type == entities.TileTypeLava {
    return entities.HazardHardness[entities.HazardLava]  // 0.3s fixed
}

// On completion, apply damage scaling with heat shield
if dugTile.Type == entities.TileTypeLava {
    baseDamage := 100.0
    damageReduction := (player.HeatResistance() / 320.0) * 50.0  // Stat facade
    finalDamage := baseDamage - damageReduction
    player.DealDamage(finalDamage)
}
```

### Drill Upgrade Scaling

Drill upgrades apply a depth-scaled divisor. At surface, only 10% applies; at max depth, 100% applies:

```go
depthFactor := depthBelowGround / maxDepth

// effectiveDivisor ranges from ~1 at surface to drillSpeed at max depth
drillSpeed := player.DrillSpeed()  // Stat facade method
effectiveDivisor := 1 + (drillSpeed-1)*(0.1+0.9*depthFactor)
duration := baseDuration / effectiveDivisor
```

### Animation Update (Each Frame)

```go
ds.animation.Elapsed += dt
progress := ds.animation.Elapsed / ds.animation.Duration
if progress > 1.0 { progress = 1.0 }

// Lerp player position toward target
player.AABB.X = ds.animation.StartX + (ds.animation.TargetX - ds.animation.StartX) * progress
player.AABB.Y = ds.animation.StartY + (ds.animation.TargetY - ds.animation.StartY) * progress

// On completion (progress >= 1.0)
if dugTile, success := ds.world.DrillTileAtGrid(ds.animation.TargetGridX, ds.animation.TargetGridY); success {
    ds.collectOreIfPresent(player, dugTile)
}
```

### Player State Flags

```go
// Physics checks drilling state to skip movement
if player.IsDrilling {
    return  // Skip velocity/collision, but heat damage still applies
}
```

---

## Physics System (`systems/physics.go`)

Orchestrates pure physics functions with axis-separated collision. For collision details, see [PHYSICS.md](PHYSICS.md).

```go
type PhysicsSystem struct {
    world *world.World
}

func (ps *PhysicsSystem) UpdatePhysics(
    player *entities.Player,
    inputState input.InputState,
    dt float32,
) {
    // 1. Apply movement and gravity to velocity
    player.Velocity = physics.ApplyHorizontalMovement(
        player.Velocity, inputState, dt,
        player.MaxSpeed(), player.Acceleration(),  // Stat facade methods
    )
    player.Velocity = physics.ApplyVerticalMovement(...)
    player.Velocity = physics.ApplyGravity(player.Velocity, dt)

    // 2. AXIS-SEPARATED COLLISION RESOLUTION
    // X-axis: integrate → check → resolve
    player.AABB.X += player.Velocity.X * dt
    collisionsX := physics.CheckCollisions(player.AABB, ps.world)
    player.AABB, player.Velocity = physics.ResolveCollisionsX(...)

    // Y-axis: integrate → check → resolve
    player.AABB.Y += player.Velocity.Y * dt
    collisionsY := physics.CheckCollisions(player.AABB, ps.world)

    // Capture state for fall damage
    wasAirborne := !player.OnGround
    ySpeedBeforeLanding := player.Velocity.Y

    player.AABB, player.Velocity, player.OnGround = physics.ResolveCollisionsY(...)

    // Apply fall damage on landing transition
    if wasAirborne && player.OnGround {
        ps.applyFallDamage(player, ySpeedBeforeLanding)
    }

    // Apply heat damage (continuous)
    physics.ApplyHeatDamage(player, dt)
}
```

---

## Fuel System (`systems/fuel.go`)

Manages fuel consumption based on player activity:

```go
type FuelSystem struct {
    // No state - purely functional
}

func (fs *FuelSystem) ConsumeFuel(
    player *entities.Player,
    inputState input.InputState,
    dt float32,
) {
    var rate float32
    if inputState.HasMovementInput() {
        rate = FuelConsumptionMoving  // 0.333 L/s
    } else {
        rate = FuelConsumptionIdle    // 0.0833 L/s
    }

    player.Fuel -= rate * dt
    if player.Fuel < 0 {
        player.Fuel = 0
    }
}
```

**Consumption Rates:**
- **Active Input** (Left, Right, Up, Drill): 10L in 30 seconds = 0.333 L/s
- **Idle** (no movement/drilling): 10L in 120 seconds = 0.0833 L/s
- **Interact Input** (E key): Uses idle rate

---

## Upgrades System (`upgrades/`)

The upgrades package provides a unified type system for player upgrades with a clean interface and catalog.

### UpgradeType Enum

```go
type UpgradeType int

const (
    TypeEngine UpgradeType = iota
    TypeHull
    TypeFuelTank
    TypeCargoHold
    TypeHeatShield
    TypeDrill
    TypeCount
)

// Array-based String/ShortName (no switch statements)
var typeNames = [TypeCount]string{"Engine", "Hull", "Fuel Tank", "Cargo Hold", "Heat Shield", "Drill"}
func (t UpgradeType) String() string { return typeNames[t] }
```

### Upgrade Interface

```go
type Upgrade interface {
    Tier() int
    Name() string
    Type() UpgradeType
}
```

All upgrade types (Engine, Hull, FuelTank, CargoHold, HeatShield, Drill) implement this interface plus their stat-specific methods.

### Unified Catalog

```go
type CatalogEntry struct {
    Price   int
    Upgrade Upgrade
}

type Catalog struct {
    entries [TypeCount][]CatalogEntry
}

func (c *Catalog) GetEntry(t UpgradeType, tier int) *CatalogEntry
func (c *Catalog) GetPrice(t UpgradeType, tier int) int
func (c *Catalog) GetName(t UpgradeType, tier int) string
func (c *Catalog) TierCount(t UpgradeType) int
```

### Player Integration

Upgrades are accessed through the Player's stat facade methods:

```go
// Read stats directly (preferred)
player.MaxSpeed()       // from Engine
player.MaxHP()          // from Hull
player.FuelCapacity()   // from FuelTank
player.CargoCapacity()  // from CargoHold
player.HeatResistance() // from HeatShield
player.DrillSpeed()     // from Drill

// Generic upgrade access (for shop/catalog)
player.GetUpgrade(upgrades.TypeEngine)
player.SetUpgrade(newEngine)
player.GetUpgradeTier(upgrades.TypeEngine)
```

---

## Effects System (`effects/`)

All state mutations go through effects with an `EffectContext`, decoupling UI from state changes and enabling effects that interact with the world and damageable entities.

### EffectContext

```go
// DamageableEntity is what effects need to damage entities in area
type DamageableEntity interface {
    GetAABB() types.AABB
    GetDamageable() *components.Damageable
    TakeDamage(amount float32)
}

// EffectContext provides access to game state for effects
type EffectContext struct {
    Player      *entities.Player
    World       *world.World
    Damageables []DamageableEntity  // Boss + future enemies
}
```

### Effect Interface

```go
type Effect interface {
    Apply(ctx *EffectContext)
}
```

### Effect Types

**Money Effects:**
```go
type TakeMoney struct { Amount int }
func (e TakeMoney) Apply(ctx *EffectContext) { ctx.Player.Money -= e.Amount }

type AddMoney struct { Amount int }
func (e AddMoney) Apply(ctx *EffectContext) { ctx.Player.Money += e.Amount }
```

**Stat Effects:**
```go
type SetFuel struct { Amount float32 }
func (e SetFuel) Apply(ctx *EffectContext) { ctx.Player.Fuel = e.Amount }

type SetHP struct { Amount float32 }
func (e SetHP) Apply(ctx *EffectContext) { ctx.Player.HP = e.Amount }
```

**Upgrade Effects:**
```go
type SetUpgrade struct { Upgrade upgrades.Upgrade }
func (e SetUpgrade) Apply(ctx *EffectContext) {
    ctx.Player.SetUpgrade(e.Upgrade)
}
```

**Inventory Effects:**
```go
type ClearOreInventory struct{}
func (e ClearOreInventory) Apply(ctx *EffectContext) { ctx.Player.OreInventory = make(map[string]int) }

type AddItem struct { ItemType entities.ItemType }
func (e AddItem) Apply(ctx *EffectContext) { ctx.Player.AddItem(e.ItemType) }
```

**Item Effects** (`effects/items.go`):
```go
type Teleport struct{}
func (e Teleport) Apply(ctx *EffectContext) {
    ctx.Player.AABB.X = ctx.Player.SpawnX
    ctx.Player.AABB.Y = ctx.Player.SpawnY
    ctx.Player.Velocity = types.Zero()
}

type Repair struct{}
func (e Repair) Apply(ctx *EffectContext) { ctx.Player.HP = ctx.Player.MaxHP() }

type Refuel struct{}
func (e Refuel) Apply(ctx *EffectContext) { ctx.Player.Fuel = ctx.Player.FuelCapacity() }

type Bomb struct { Radius int; Damage float32 }
func (e Bomb) Apply(ctx *EffectContext) {
    // Damage damageable entities in range
    for _, entity := range ctx.Damageables {
        if entity.GetAABB().Intersects(blastAABB) {
            entity.TakeDamage(e.Damage)
        }
    }
    // Destroy tiles in circular radius
    ctx.World.NukeTileAtGrid(...)
}
```

### Processor

```go
type Processor struct{}

func (p *Processor) Apply(ctx *EffectContext, effects []Effect) {
    for _, effect := range effects {
        effect.Apply(ctx)
    }
}
```

---

## UI System (`ui/`)

Unified UI layer handling modal (shops) and instant (market/hospital/fuel) interactions.

### Core Types

```go
type Result struct {
    ShouldClose bool
    Effects     []effects.Effect
}

func NoChange() Result { return Result{} }
func Close() Result { return Result{ShouldClose: true} }
func WithEffects(effs ...effects.Effect) Result { return Result{Effects: effs} }
func CloseWithEffects(effs ...effects.Effect) Result { return Result{ShouldClose: true, Effects: effs} }

type UI interface {
    Process(player *entities.Player, inputState input.InputState) Result
    GetRenderState() interface{}
}
```

### Manager

```go
type Manager struct {
    uis        map[components.InteractableType]UI
    activeUI   UI
    activeType components.InteractableType
}

func (m *Manager) Register(t components.InteractableType, ui UI)
func (m *Manager) OpenUI(t components.InteractableType) bool
func (m *Manager) Process(player *entities.Player, inputState input.InputState) Result
func (m *Manager) HasActiveUI() bool
func (m *Manager) GetActiveUI() UI
```

### Modal UIs (UpgradeShopUI, ItemShopUI)

- Return `NoChange()` to stay open
- Return `WithEffects(...)` on purchase (stay open)
- Return `Close()` on Q/Escape
- Have render state for display

### Instant UIs (MarketUI, HospitalUI, FuelStationUI)

- Return `CloseWithEffects(...)` immediately
- Close on first process (no modal)
- No render state needed

---

## Interaction System (`systems/interaction.go`)

Simple function to detect player-building overlap with E key:

```go
func DetectInteraction(
    player *entities.Player,
    buildings []*entities.Building,
    inputState input.InputState,
) *components.InteractableType {
    if !inputState.Interact {
        return nil
    }

    playerPos := components.Position{AABB: player.AABB}
    for _, b := range buildings {
        if b.Position.Intersects(playerPos) {
            interactableType := b.Interactable.Type
            return &interactableType
        }
    }
    return nil
}
```

---

## Item Usage (`systems/items.go`)

Item usage detection is a pure function that returns effects to be applied via `EffectContext`:

```go
func DetectItemUsage(player *entities.Player, inputState input.InputState, itemCfg config.ItemConfig) []effects.Effect {
    var result []effects.Effect

    if inputState.UseTeleport && player.UseItem(entities.ItemTeleport) {
        result = append(result, effects.Teleport{})
    }
    if inputState.UseRepair && player.UseItem(entities.ItemRepair) {
        result = append(result, effects.Repair{})
    }
    if inputState.UseRefuel && player.UseItem(entities.ItemRefuel) {
        result = append(result, effects.Refuel{})
    }
    if inputState.UseBomb && player.UseItem(entities.ItemBomb) {
        result = append(result, effects.Bomb{Radius: itemCfg.Bomb.Radius, Damage: 10.0})
    }
    if inputState.UseBigBomb && player.UseItem(entities.ItemBigBomb) {
        result = append(result, effects.Bomb{Radius: itemCfg.BigBomb.Radius, Damage: 25.0})
    }

    return result
}
```

Item effects are defined in `effects/items.go` (see Effects System above).

### Bomb Damage

Bombs damage both tiles and damageable entities (bosses):
- **Regular Bomb**: Radius 2 tiles, 10 HP damage to entities
- **Big Bomb**: Radius 4 tiles, 25 HP damage to entities

Entities control their own vulnerability—the boss's `TakeDamage` method checks `IsVulnerable()` before applying damage.

### World Methods for Tile Removal

**DrillTileAtGrid** — Standard drilling (respects drillability):
```go
func (w *World) DrillTileAtGrid(gridX, gridY int) (*entities.Tile, bool) {
    tile := w.tiles[[2]int{gridX, gridY}]
    if tile != nil && tile.IsDrillable() {
        delete(w.tiles, [2]int{gridX, gridY})
        return tile, true
    }
    return nil, false
}
```

**NukeTileAtGrid** — Bomb destruction (bypasses drillability):
```go
func (w *World) NukeTileAtGrid(gridX, gridY int) (*entities.Tile, bool) {
    tile := w.tiles[[2]int{gridX, gridY}]
    if tile != nil && tile.IsSolid() {
        delete(w.tiles, [2]int{gridX, gridY})
        return tile, true
    }
    return nil, false
}
```

**Key Distinction:**
- Rock: `IsSolid() = true`, `IsDrillable() = false` → DrillTile fails, NukeTile succeeds
- Lava: `IsSolid() = true`, `IsDrillable() = true` → Both succeed
- Dirt: `IsSolid() = true`, `IsDrillable() = true` → Both succeed

---

## Boss Fight System (`systems/boss_fight.go`)

Orchestrates boss encounters. For detailed boss mechanics, see [BOSS.md](BOSS.md).

```go
type BossFightSystem struct {
    boss              bosses.Boss
    bossRoomStartY    float32
    playerInBossRoom  bool
}

func (bfs *BossFightSystem) Update(player *entities.Player, dt float32) entities.GameState {
    // Track player entry/exit
    playerY := player.AABB.Y + player.AABB.Height
    inRoom := playerY >= bfs.bossRoomStartY

    if inRoom && !bfs.playerInBossRoom {
        bfs.boss.Activate()
    } else if !inRoom && bfs.playerInBossRoom {
        bfs.boss.Deactivate()
    }
    bfs.playerInBossRoom = inRoom

    // Update boss if active
    if bfs.boss.IsActive() {
        bfs.boss.Update(player, dt)
        bfs.handleProjectileCollisions(player)
        bfs.handleContactDamage(player, dt)
    }

    // Check win/lose conditions
    if bfs.boss.IsDefeated() {
        return entities.GameStateVictory
    }
    if player.HP <= 0 {
        return entities.GameStateDefeat
    }
    return entities.GameStatePlaying
}
```
