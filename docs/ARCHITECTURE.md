# Architecture

## Overview

Drill Game uses **Hexagonal Architecture** (Ports & Adapters) to achieve a clean separation between pure domain logic and framework-specific integration. This ensures the game logic is:

- **Testable**: Physics and game logic can be tested without initializing Raylib
- **Portable**: Domain code has zero framework dependencies
- **Maintainable**: Clear responsibilities and data flow
- **Extensible**: Easy to add new features, entities, or rendering backends

---

## Architecture Pattern: Hexagonal Architecture

```
┌─────────────────────────────────────────────────────┐
│                   APPLICATION LAYER                 │
│                  (cmd/game/main.go)                 │
│    Orchestrates: Input Reading → Game Update → Rendering
└────────────┬──────────────────────────────┬─────────┘
             │                              │
      ┌──────▼─────────┐        ┌──────────▼──────┐
      │ INPUT ADAPTER  │        │ RENDERING ADAPTER│
      │ (Raylib Keys)  │        │  (Raylib Drawing)
      └──────┬─────────┘        └──────────┬──────┘
             │                              │
             │      ┌──────────────────┐    │
             │      │   DOMAIN LAYER   │    │
             │      │  (Pure Business  │    │
             └─────▶│     Logic)       │◀───┘
                    │                  │
                    │ • Game (engine/) │
                    │ • Player (entities/)
                    │ • Physics (systems/)
                    │ • Types (types/)
                    │ • Input State (input/)
                    │ • World (world/)
                    │
                    └──────────────────┘
```

### The Three Layers

1. **Application Layer** (`cmd/game/main.go`)
   - Orchestrates the entire game
   - Reads input from adapters
   - Updates domain logic
   - Delegates rendering to adapters
   - Manages window lifecycle

2. **Adapter Layer** (`internal/adapters/`)
   - **Input Adapter**: Translates Raylib keyboard input to domain `InputState`
   - **Rendering Adapter**: Takes domain entities and renders them with Raylib
   - All Raylib integration lives here
   - No business logic

3. **Domain Layer** (`internal/domain/`)
   - Pure game logic, zero framework dependencies
   - Can be tested without Raylib initialization
   - Fully portable (could swap Raylib for SDL, canvas, etc.)

---

## Project Structure

```
drill-game/
├── cmd/
│   └── game/
│       └── main.go                          # Application orchestration
│
├── internal/
│   ├── adapters/                            # Framework Integration (Raylib)
│   │   ├── input/
│   │   │   └── raylib.go                    # RaylibInputAdapter
│   │   └── rendering/
│   │       └── raylib.go                    # RaylibRenderer
│   │
│   └── domain/                              # Pure Business Logic
│       ├── engine/
│       │   └── game.go                      # Game orchestration (domain)
│       ├── systems/
│       │   ├── physics.go                   # PhysicsSystem
│       │   ├── digging.go                   # DiggingSystem (ore collection)
│       │   ├── shop.go                      # ShopSystem (selling inventory)
│       │   ├── fuel.go                      # FuelSystem (consumption based on activity)
│       │   ├── fuel_station.go              # FuelStationSystem (refueling)
│       │   ├── hospital.go                  # HospitalSystem (healing HP)
│       │   ├── upgrade.go                   # UpgradeSystem (purchase upgrades at shops)
│       │   ├── digging_test.go              # Digging & ore collection tests
│       │   ├── fuel_test.go                 # Fuel consumption tests
│       │   ├── fuel_station_test.go         # Fuel station transaction tests
│       │   ├── hospital_test.go             # Hospital healing transaction tests
│       │   └── upgrade_test.go              # Upgrade purchase tests
│       ├── entities/
│       │   ├── player.go                    # Player aggregate root (AABB, inventory, money, fuel, HP, components)
│       │   ├── player_test.go               # Player inventory tests
│       │   ├── engine.go                    # Engine component (tier, name, speed/acceleration stats)
│       │   ├── hull.go                      # Hull component (tier, name, maxHP)
│       │   ├── fuel_tank.go                 # FuelTank component (tier, name, capacity)
│       │   ├── cargo_hold.go                # CargoHold component (tier, name, ore capacity)
│       │   ├── heat_shield.go               # HeatShield component (tier, name, heat resistance)
│       │   ├── tile.go                      # Tile entity (Empty, Dirt, Ore)
│       │   ├── shop.go                      # Shop entity (AABB-based interactable)
│       │   ├── fuel_station.go              # FuelStation entity (AABB-based interactable)
│       │   ├── hospital.go                  # Hospital entity (AABB-based interactable)
│       │   ├── upgrade_shop.go              # UpgradeShop types with catalogs (Engine/Hull/FuelTank/CargoHold/HeatShield)
│       │   └── ore_type.go                  # Ore types & values, Gaussian parameters
│       ├── physics/
│       │   ├── constants.go                 # Physics parameters
│       │   ├── movement.go                  # Movement functions
│       │   ├── gravity.go                   # Gravity + velocity integration
│       │   ├── collision.go                 # AABB collision detection/resolution
│       │   ├── damage.go                    # Fall damage calculations
│       │   ├── heat.go                      # Temperature calculation & heat damage
│       │   ├── movement_test.go             # Movement tests
│       │   ├── gravity_test.go              # Gravity tests
│       │   └── collision_test.go            # AABB collision tests
│       ├── types/
│       │   ├── vec2.go                      # Custom Vec2 (no Raylib types)
│       │   ├── aabb.go                      # AABB collision primitive
│       │   └── aabb_test.go                 # AABB unit tests
│       ├── input/
│       │   ├── input_state.go               # InputState struct (framework-agnostic)
│       │   └── input_state_test.go          # InputState helper method tests
│       └── world/
│           ├── world.go                     # World: chunk loading, sparse tile map
│           ├── generator.go                 # Procedural tile generation
│           ├── hash.go                      # Deterministic seeding (FNV-1a)
│           ├── generator_test.go            # Generator unit tests
│           ├── world_test.go                # Chunk loading tests
│           └── integration_test.go          # End-to-end world generation tests
│
├── docs/
│   ├── ARCHITECTURE.md                      # This file
│   └── GAME_DESIGN.md                       # Game mechanics
├── go.mod
├── go.sum
└── README.md
```

---

## Data Flow

### Single Frame Update

```
main.go Loop:
┌─────────────────────────────────────────┐
│ 1. Read Input from Adapter              │
│    inputState := inputAdapter.ReadInput()
│    (Converts Raylib keys → InputState)  │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│ 2. Update Domain Logic                  │
│    game.Update(dt, inputState)          │
│    • Load chunks around player (3×3)    │
│    • Apply physics & fall damage        │
│    • Consume fuel (active or idle)      │
│    • Drill downward or horizontal       │
│    • Animate player over 1 second       │
│    • Remove tile on animation complete  │
│    • Collect ore if available           │
│    • (Block other inputs during dig)    │
│    • Shop selling (E key + overlap)     │
│    • Fuel station refueling (E key)     │
│    • Hospital healing (E key)           │
│    • Upgrade purchases (E key)          │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│ 3. Render via Adapter                   │
│    renderer.Render(game)                │
│    • Extracts Player, World, Shop       │
│    • Renders tiles with ore colors      │
│    • Draws shop, player, entities       │
│    • Displays debug info (money, ore, fuel) │
│    • Camera follows player              │
└─────────────────────────────────────────┘
```

### Code Example

```go
// main.go - Application layer orchestration
for renderer.WindowShouldClose() == false {
    dt := renderer.GetFrameTime()

    // 1. Read input (adapter responsibility)
    inputState := inputAdapter.ReadInput()

    // 2. Update game (domain logic - pure, testable)
    err := game.Update(dt, inputState)
    if err != nil {
        return err
    }

    // 3. Render (adapter responsibility)
    renderer.Render(game)
}
```

---

## Core Concepts

### Domain Layer (Pure Game Logic)

#### Game (`domain/engine/game.go`)

Orchestrates domain systems without any framework knowledge:

```go
type Game struct {
    world             *world.World
    player            *entities.Player
    physicsSystem     *systems.PhysicsSystem
    diggingSystem     *systems.DiggingSystem
    shopSystem        *systems.ShopSystem
    fuelSystem        *systems.FuelSystem
    fuelStationSystem *systems.FuelStationSystem
    hospitalSystem    *systems.HospitalSystem
    upgradeSystem     *systems.UpgradeSystem
}

func (g *Game) Update(dt float32, inputState input.InputState) error {
    // Pure domain logic - zero Raylib

    // 0. Update chunks around player (proactive loading)
    playerX := g.player.AABB.X + g.player.AABB.Width/2
    playerY := g.player.AABB.Y + g.player.AABB.Height/2
    g.world.UpdateChunksAroundPlayer(playerX, playerY)

    // 1. Physics FIRST - handles landing/fall damage before dig can start
    //    Also applies heat damage and skips movement during dig animation
    g.physicsSystem.UpdatePhysics(g.player, inputState, dt)

    // 2. Always: fuel consumption (runs even during dig animation)
    g.fuelSystem.ConsumeFuel(g.player, inputState, dt)

    // 3. Handle digging (vertical + horizontal, with animation)
    g.diggingSystem.ProcessDigging(g.player, inputState, dt)

    // Skip interactions during dig animation
    if g.player.IsDigging {
        return nil
    }

    // 4. Handle shop selling
    g.shopSystem.ProcessSelling(g.player, inputState)

    // 5. Handle fuel station refueling
    g.fuelStationSystem.ProcessRefueling(g.player, inputState)

    // 6. Handle hospital healing
    g.hospitalSystem.ProcessHealing(g.player, inputState)

    // 7. Handle upgrade purchases
    g.upgradeSystem.ProcessUpgrade(g.player, inputState)

    return nil
}

// Adapters pull data from getters
func (g *Game) GetWorld() *world.World {
    return g.world
}

func (g *Game) GetPlayer() *entities.Player {
    return g.player
}
```

**Why this design:**
- Game doesn't import Raylib (fully testable)
- Game accepts InputState (not Raylib keys)
- Game provides getters for adapters to read state
- All rendering responsibility is external

#### Drilling System (`domain/systems/digging.go`)

Handles both vertical and horizontal drilling with smooth 1-second animations. When a dig is initiated, the player interpolates toward the tile center while the tile is progressively revealed. The tile is only removed when the animation completes.

**Core Concepts:**

```go
type DiggingSystem struct {
    world     *world.World
    animation DiggingAnimation
}

// Animation state tracks active digs
type DiggingAnimation struct {
    Active      bool
    Direction   DigDirection  // Down, Left, or Right
    StartX      float32       // Player position when animation started
    StartY      float32
    TargetX     float32       // Where player moves to during animation
    TargetY     float32
    TargetGridX int           // Tile coordinates for removal
    TargetGridY int
    Elapsed     float32       // Time elapsed in animation
    Duration    float32       // Total animation duration (1.0 second)
    Tile        *entities.Tile
}

const DigAnimationDuration float32 = 1.0 // seconds
```

**Game Loop Flow:**

The digging system receives input AFTER physics (which handles landing), ensuring:
1. Player lands on ground (physics)
2. Fall damage applied if landing from height (physics)
3. Digging animation can only start when grounded
4. Player is locked in animation for 1 second
5. Tile removed on completion, ore collected

```go
// Game loop order (internal/domain/engine/game.go)
g.physicsSystem.UpdatePhysics(g.player, inputState, dt)      // Physics FIRST
g.fuelSystem.ConsumeFuel(g.player, inputState, dt)           // Fuel runs always
g.diggingSystem.ProcessDigging(g.player, inputState, dt)     // Then digging
if g.player.IsDigging { return }                             // Block interactions during dig
// ... interaction systems (shop, upgrade, etc.)
```

**Vertical Digging (S/Down Key):**

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

// Start 1-second animation to target position
ds.startDigAnimation(player, DigDown, tileGridX, tileGridY, targetX, targetY, tile)
```

**Horizontal Digging (Left/Right When Grounded):**

```go
// Check tile beside player
playerCenterY := player.AABB.Y + player.AABB.Height/2

if inputState.Left {
    tile := ds.world.GetTileAt(player.AABB.X - 1, playerCenterY)
    if tile != nil && tile.IsDiggable() {
        // Move to tile center (X) but keep current Y
        tileCenterX := float32(tileGridX)*world.TileSize + world.TileSize/2
        targetX := tileCenterX - player.AABB.Width/2
        targetY := player.AABB.Y  // Stay at ground level

        ds.startDigAnimation(player, DigLeft, tileGridX, tileGridY, targetX, targetY, tile)
    }
}
// Similar for Right
```

**Animation Update (Each Frame):**

```go
// Lerp player position toward target
ds.animation.Elapsed += dt
progress := ds.animation.Elapsed / ds.animation.Duration
if progress > 1.0 { progress = 1.0 }

player.AABB.X = ds.animation.StartX + (ds.animation.TargetX - ds.animation.StartX) * progress
player.AABB.Y = ds.animation.StartY + (ds.animation.TargetY - ds.animation.StartY) * progress

// On completion (progress >= 1.0)
if dugTile, success := ds.world.DigTileAtGrid(ds.animation.TargetGridX, ds.animation.TargetGridY); success {
    ds.collectOreIfPresent(player, dugTile)
}
```

**Player State Flags:**

The player has two state flags set by systems:
- `player.OnGround` — Set by physics system on ground contact
- `player.IsDigging` — Set by digging system during animation

This allows both physics and rendering to query state directly:

```go
// Physics checks digging state to skip movement
if player.IsDigging {
    return  // Skip velocity/collision, but heat damage still applies
}

// Rendering displays state
fmt.Sprintf("IsDigging: %v", player.IsDigging)
```

**Why This Design:**

- **Clear State Machine**: Animation lifecycle (start → update → complete) is explicit
- **Lerp-Based Movement**: Smooth animation feels natural vs teleporting
- **Grounding Requirement**: Only dig when on ground (prevents mid-air exploits)
- **Continuous Effects**: Fuel consumption and heat damage run during animation
- **Fall Damage Protection**: Physics runs first, applies damage before dig starts
- **One System Responsibility**: All animation logic in DiggingSystem, not scattered
- **Testable**: Animation state tracked in struct, no external dependencies

#### Physics System (`domain/systems/physics.go`)

Orchestrates pure physics functions with axis-separated collision and fall damage:

```go
type PhysicsSystem struct {
    world *world.World
}

func (ps *PhysicsSystem) UpdatePhysics(
    player *entities.Player,
    inputState input.InputState,
    dt float32,
) {
    // 1. Apply movement and gravity to velocity (using player's component stats)
    player.Velocity = physics.ApplyHorizontalMovement(
        player.Velocity, inputState, dt,
        player.Engine.MaxSpeed(), player.Engine.Acceleration(),
    )
    player.Velocity = physics.ApplyVerticalMovement(
        player.Velocity, inputState, dt,
        player.Engine.FlyAcceleration(), player.Engine.MaxUpwardSpeed(),
    )
    player.Velocity = physics.ApplyGravity(player.Velocity, dt)

    // 2. AXIS-SEPARATED COLLISION RESOLUTION

    // X-axis: integrate position → check → resolve
    player.AABB.X += player.Velocity.X * dt
    collisionsX := physics.CheckCollisions(player.AABB, ps.world)
    player.AABB, player.Velocity = physics.ResolveCollisionsX(player.AABB, player.Velocity, collisionsX)

    // Y-axis: integrate position → check → resolve
    player.AABB.Y += player.Velocity.Y * dt
    collisionsY := physics.CheckCollisions(player.AABB, ps.world)

    // Capture state before Y-resolution for fall damage calculation
    wasAirborne := !player.OnGround
    ySpeedBeforeLanding := player.Velocity.Y

    player.AABB, player.Velocity, player.OnGround = physics.ResolveCollisionsY(player.AABB, player.Velocity, collisionsY)

    // Apply fall damage on landing transition
    if wasAirborne && player.OnGround {
        ps.applyFallDamage(player, ySpeedBeforeLanding)
    }
}
```

**Why this design:**
- Direct field access (no getters/setters) for simplicity
- Axis-separated collision prevents corner-catching
- Pure physics functions fully testable without framework
- Accepts `InputState` (not Raylib types)
- Works with `*Player` directly (no interface needed)

**Fall Damage Implementation:**

Damage is calculated when a landing transition occurs (airborne → grounded). The physics package computes damage, and Player.DealDamage() applies it:

```go
// domain/physics/fall_damage.go - Pure calculation function
func ApplyFallDamage(player *entities.Player, ySpeed float32) {
    if ySpeed < FallDamageThreshold {
        return  // Below 500 px/sec threshold - safe landing
    }

    // Calculate damage: (ySpeed - threshold) / divisor
    damage := (ySpeed - FallDamageThreshold) / FallDamageDivisor

    // Apply through Player.DealDamage() which clamps HP at 0
    player.DealDamage(damage)
}

// domain/systems/physics.go - Called from PhysicsSystem.UpdatePhysics()
if wasAirborne && player.OnGround {
    physics.ApplyFallDamage(player, ySpeedBeforeLanding)
}
```

**Landing Detection:**
- `wasAirborne = !player.OnGround` captures state before Y resolution
- `player.OnGround` is set to true by `ResolveCollisionsY()` on ground contact
- Condition `wasAirborne && player.OnGround` ensures damage only on transition
- Prevents repeated damage when already grounded

**Damage Application Pattern:**
- Physics calculates damage: `damage = (ySpeed - threshold) / divisor`
- Player applies damage: `player.DealDamage(damage)` clamps at zero
- Centralizes HP mutation logic in Player entity (aggregate root)

#### Heat System (`domain/physics/heat.go` & `domain/systems/physics.go`)

Temperature increases with depth and deals exponential damage when exceeding heat resistance:

```go
// domain/physics/heat.go - Pure calculation + damage application
func ApplyHeatDamage(player *entities.Player, dt float32) {
    // Calculate temperature at player depth (15°C at surface, 350°C at max depth)
    temperature := CalculateTemperature(player.AABB.Y)

    // Check excess heat beyond resistance
    excessHeat := temperature - player.HeatShield.HeatResistance()
    if excessHeat <= 0 {
        return  // Within safe temperature range
    }

    // Apply exponential damage: baseDPS * (excessHeat / divisor)^exponent * dt
    damagePerSecond := float32(HeatDamageBaseDPS) *
        float32(math.Pow(float64(excessHeat/float32(HeatDamageDivisor)),
                         float64(HeatDamageExponent)))
    damage := damagePerSecond * dt

    // Apply through Player.DealDamage() which clamps HP at 0
    player.DealDamage(damage)
}

// domain/systems/physics.go - Called every frame from PhysicsSystem.UpdatePhysics()
physics.ApplyHeatDamage(player, dt)
```

**Temperature Calculation:**
- **Ground Level** (Y=640): 15°C base temperature
- **Max Depth** (Y=64,000): 350°C maximum temperature
- **Formula**: Linear interpolation based on depth below ground
- **No Damage**: Above ground level (Y < 640)

**Damage Constants:**
- `HeatDamageBaseDPS = 0.5` — Base damage per second
- `HeatDamageDivisor = 10.0` — Scaling factor
- `HeatDamageExponent = 1.5` — Exponential curve steepness

**Heat Shield Component:**
6 tiers enabling safe mining at progressively deeper zones:

| Tier | Resistance | Price | Safe Depth |
|------|------------|-------|-----------|
| Base | 50°C | - | 0-6,600px |
| Mk1 | 90°C | $200 | 6,600-14,000px |
| Mk2 | 140°C | $500 | 14,000-23,500px |
| Mk3 | 190°C | $1,200 | 23,500-33,000px |
| Mk4 | 250°C | $3,000 | 33,000-44,500px |
| Mk5 | 320°C | $7,500 | 44,500-64,000px |

**Why this design:**
- Called every frame continuously (unlike fall damage on landing only)
- Exponential scaling creates meaningful progression gates
- Heat becomes limiting factor for deep mining
- Upgrades enable deeper exploration without level caps

#### Fuel System (`domain/systems/fuel.go`)

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
    // Determine consumption rate based on input
    var rate float32
    if inputState.HasMovementInput() {
        rate = FuelConsumptionMoving  // 0.333 L/s
    } else {
        rate = FuelConsumptionIdle    // 0.0833 L/s
    }

    // Consume fuel this frame
    player.Fuel -= rate * dt
    if player.Fuel < 0 {
        player.Fuel = 0
    }
}
```

**Consumption Rates:**
- **Active Input** (Left, Right, Up, Dig): 10L in 30 seconds = 0.333 L/s
- **Idle** (no movement/digging): 10L in 120 seconds = 0.0833 L/s
- **Sell Input** (E key): Uses idle rate (not active activity)

**Why this design:**
- Called after physics to ensure movement is fully resolved
- Uses `HasMovementInput()` helper to distinguish movement from shop interaction
- Simple rate-based consumption with delta-time independence
- Fuel clamped at zero (no negative values)
- Direct field mutation for simplicity (no getters/setters)
- Pure logic - could be replaced without affecting game structure

#### Fuel Station System (`domain/systems/fuel_station.go`)

Manages refueling transactions at the fuel station:

```go
type FuelStationSystem struct {
    fuelStation *entities.FuelStation
}

func (fss *FuelStationSystem) ProcessRefueling(
    player *entities.Player,
    inputState input.InputState,
) {
    if !inputState.Sell {
        return
    }

    if !fss.fuelStation.IsPlayerInRange(player) {
        return
    }

    // Calculate liters needed to fill tank
    litersNeeded := entities.FuelCapacity - player.Fuel
    
    // Calculate cost (1 money per liter, rounded up)
    cost := int(math.Ceil(float64(litersNeeded)))

    // Check if player has enough money
    if player.Money < cost {
        return // Cannot afford refueling
    }

    // Deduct money and refuel
    player.Money -= cost
    player.Fuel = entities.FuelCapacity
}
```

**Transaction Rules:**
- **Cost**: $1 per liter needed (rounded up using `math.Ceil`)
- **Examples**:
  - Full tank (0L needed) = $0 (no transaction)
  - Empty tank (10L needed) = $10
  - 6.8L current (3.2L needed) = $4 (ceil of 3.2)
  - 9.9L current (0.1L needed) = $1 (ceil of 0.1)
- **Rejection**: Transaction rejected if player has insufficient money
- **Instant Fill**: Fuel immediately set to `FuelCapacity` (10L) on success

**Why this design:**
- Mirrors ShopSystem pattern (AABB + interaction key)
- Uses same E key as shop (no spatial conflict due to separation)
- Called before physics (consistent with shop processing)
- Simple rounding up ensures player always pays at least $1 if not full
- Direct field mutation for clarity
- Fully testable without framework (6 comprehensive unit tests)

#### Hospital System (`domain/systems/hospital.go`)

Manages healing transactions at the hospital:

```go
type HospitalSystem struct {
    hospital *entities.Hospital
}

func (hs *HospitalSystem) ProcessHealing(
    player *entities.Player,
    inputState input.InputState,
) {
    if !inputState.Sell {
        return
    }

    if !hs.hospital.IsPlayerInRange(player) {
        return
    }

    // Calculate HP needed to reach max HP
    hpNeeded := entities.MaxHP - player.HP

    // Early exit: Already at max HP (no healing needed)
    if hpNeeded <= 0 {
        return
    }

    // Calculate cost: $2 per HP, rounded up
    cost := int(math.Ceil(float64(hpNeeded) * 2.0))

    // Early exit: Insufficient money
    if player.Money < cost {
        return // Cannot afford healing
    }

    // Execute transaction: Deduct money and restore HP to max
    player.Money -= cost
    player.HP = entities.MaxHP
}
```

**Transaction Rules:**
- **Cost**: $2 per HP needed (rounded up using `math.Ceil`)
- **Examples**:
  - Max HP (0 HP needed) = $0 (no transaction)
  - Zero HP (10 HP needed) = $20
  - 7.2 HP (2.8 HP needed) = $6 (ceil of 5.6)
  - 9.9 HP (0.1 HP needed) = $1 (ceil of 0.2)
- **Rejection**: Transaction rejected if player has insufficient money
- **Instant Heal**: HP immediately set to `MaxHP` (10.0) on success

**Why this design:**
- Mirrors FuelStationSystem pattern (AABB + interaction key)
- Uses same E key as shop and fuel station (no spatial conflict due to separation)
- Called before physics (consistent with shop/fuel station processing)
- Simple rounding up ensures player always pays at least $1 if not at max HP
- Direct field mutation for clarity
- Fully testable without framework (7 comprehensive unit tests)

#### Upgrade System (`domain/systems/upgrade.go`)

Manages upgrade purchases at dedicated upgrade shops:

```go
type UpgradeSystem struct {
    engineShop    *entities.EngineUpgradeShop
    hullShop      *entities.HullUpgradeShop
    fuelTankShop  *entities.FuelTankUpgradeShop
    cargoHoldShop *entities.CargoHoldUpgradeShop
}

func (us *UpgradeSystem) ProcessUpgrade(
    player *entities.Player,
    inputState input.InputState,
) {
    if !inputState.Sell {
        return
    }

    // Check each shop and attempt upgrade (only one can be in range at a time)
    if us.engineShop.IsPlayerInRange(player) {
        us.tryUpgradeEngine(player)
        return
    }
    // ... similar for hullShop and fuelTankShop
}

func (us *UpgradeSystem) tryUpgradeEngine(player *entities.Player) {
    entry := us.engineShop.GetNextEngine(player.Engine.Tier())
    if entry == nil {
        return // Already at max level
    }
    if !player.CanAfford(entry.Price) {
        return // Cannot afford
    }
    player.BuyEngine(entry.Engine, entry.Price)
}
```

**Upgrade Rules:**
- **Sequential**: Must buy Mk1 before Mk2, etc.
- **Permanent**: Upgrades cannot be undone
- **No Auto-Restore**: Hull/Tank upgrades don't restore HP/Fuel to new max

**Why this design:**
- Mirrors Hospital/FuelStation pattern (AABB + E key interaction)
- Three separate shops prevent spatial conflict
- Called before physics (consistent with other interactions)
- Each shop owns its catalog (DDD: shop knows what it sells)
- Player is aggregate root (mutations go through Player methods)
- Fully testable without framework (5 comprehensive unit tests)

#### Pure Physics Functions (`domain/physics/`)

Framework-independent mathematical functions:

```go
// movement.go - Pure functions, no Raylib, fully testable
// Parameters for max speed/acceleration come from player upgrades
func ApplyHorizontalMovement(velocity Vec2, inputState InputState, dt float32, maxSpeed, acceleration float32) Vec2
func ApplyVerticalMovement(velocity Vec2, inputState InputState, dt float32, flyAcceleration, maxUpwardVelocity float32) Vec2

// gravity.go - Pure functions
func ApplyGravity(velocity Vec2, dt float32) Vec2

// collision.go - AABB-based collision functions
func CheckCollisions(aabb AABB, world *World) []TileCollision
func ResolveCollisionsX(aabb AABB, velocity Vec2, collisions []TileCollision) (AABB, Vec2)
func ResolveCollisionsY(aabb AABB, velocity Vec2, collisions []TileCollision) (AABB, Vec2, bool)
func GetOccupiedTileRange(aabb AABB, tileSize float32) (minX, maxX, minY, maxY int)
```

**Why this design:**
- Zero Raylib imports
- Can be tested standalone
- Input/output are domain types (AABB, Vec2, etc.)
- Pure functions enable unit testing without framework
- Value-based (no pointer mutations in function signatures)

#### Player Entity (`domain/entities/player.go`)

Player is the **aggregate root** with exported component value objects (Engine, Hull, FuelTank, CargoHold). Stats are accessed via components; mutations go through Player methods.

```go
type Player struct {
    AABB         types.AABB  // Position and dimensions
    Velocity     types.Vec2  // Pixels per second
    OnGround     bool        // Collision state
    OreInventory [7]int      // Ore counts indexed by OreType
    Money        int         // Currency from ore sales
    Fuel         float32     // Current fuel in liters
    HP           float32     // Hit points
    Engine       Engine      // Engine component (exported)
    Hull         Hull        // Hull component (exported)
    FuelTank     FuelTank    // FuelTank component (exported)
    CargoHold    CargoHold   // CargoHold component (exported)
    HeatShield   HeatShield  // HeatShield component (exported)
}

// Component types are value objects with named constructors
type Engine struct {
    tier, name, maxSpeed, acceleration, flyAcceleration, maxUpwardSpeed
}
func NewEngineBase() Engine  // tier 0, 450 px/s max speed, etc.
func NewEngineMk1() Engine   // tier 1, 500 px/s max speed, etc.
// ... through NewEngineMk5()

// Stats accessed via components
player.Engine.MaxSpeed()      // 450.0 for base engine
player.Engine.Tier()          // 0 for base engine
player.Hull.MaxHP()           // 10.0 for base hull
player.FuelTank.Capacity()    // 10.0 for base tank
player.CargoHold.Capacity()   // 10 for base cargo hold
player.HeatShield.HeatResistance() // 50.0 for base heat shield
player.GetTotalOreCount()     // Sum of all ore in inventory

// Purchase methods enforce invariants
func (p *Player) CanAfford(cost int) bool
func (p *Player) BuyEngine(e Engine, cost int)
func (p *Player) BuyHull(h Hull, cost int)
func (p *Player) BuyFuelTank(ft FuelTank, cost int)
func (p *Player) BuyCargoHold(ch CargoHold, cost int)
func (p *Player) BuyHeatShield(hs HeatShield, cost int)
func (p *Player) Refuel() bool  // checks money, fills tank
func (p *Player) Heal() bool    // checks money, restores HP
func (p *Player) AddOre(oreType OreType) bool  // returns false if cargo full

// Damage application (called by physics damage sources)
func (p *Player) DealDamage(damage float32)  // applies damage, clamps HP at 0

func NewPlayer(startX, startY float32) *Player {
    engine := NewEngineBase()
    hull := NewHullBase()
    fuelTank := NewFuelTankBase()
    cargoHold := NewCargoHoldBase()
    return &Player{
        AABB:      types.NewAABB(startX, startY, PlayerWidth, PlayerHeight),
        Velocity:  types.Zero(),
        Fuel:      fuelTank.Capacity(),
        HP:        hull.MaxHP(),
        Engine:    engine,
        Hull:      hull,
        FuelTank:  fuelTank,
        CargoHold: cargoHold,
    }
}

// AddOre increments ore count for given type
func (p *Player) AddOre(oreType OreType, amount int) {
    if oreType >= 0 && oreType < 7 {
        p.OreInventory[oreType] += amount
    }
}
```

**Why this design:**
- AABB eliminates redundant Position storage (X, Y already in AABB)
- No Render() method (rendering is adapter responsibility)
- Uses domain types (AABB, Vec2, not rl.Vector2)
- Zero Raylib dependency
- Direct field access (no getters/setters) for simplicity
- AABB enables proper collision detection (not just ground)
- `OreInventory [7]int` stores counts for all 7 ore types efficiently
- `AddOre()` is the only ore collection method (simple, one-purpose)

#### Types (`domain/types/`)

Custom math types independent of framework:

**Vec2** (`vec2.go`):
```go
type Vec2 struct {
    X float32
    Y float32
}

func (v Vec2) Add(other Vec2) Vec2
func (v Vec2) Scale(scalar float32) Vec2
func (v Vec2) Magnitude() float32
// ... other operations
```

**AABB** (`aabb.go`):
```go
type AABB struct {
    X, Y          float32 // Top-left corner position
    Width, Height float32 // Dimensions
}

func (a AABB) Intersects(b AABB) bool
func (a AABB) Penetration(b AABB) (dx, dy float32)
func (a AABB) Min() Vec2
func (a AABB) Max() Vec2
```

**Why this design:**
- No Raylib dependency
- Physics can use its own types
- Conversion to Raylib types only happens in rendering adapter
- AABB provides proper collision detection (not just point-based)
- Value types (not pointers) for simplicity and performance

#### InputState (`domain/input/input_state.go`)

Platform-agnostic input representation:

```go
type InputState struct {
    Left  bool  // A or Arrow Left - move left
    Right bool  // D or Arrow Right - move right
    Up    bool  // W or Arrow Up - jump/fly
    Dig   bool  // S or Arrow Down - dig downward
    Sell  bool  // E - sell inventory at shop
}

// HasMovementInput returns true if player is actively moving or digging
func (is InputState) HasMovementInput() bool {
    return is.Left || is.Right || is.Up || is.Dig
}

// HasHorizontalInput returns true if player is moving left or right
func (is InputState) HasHorizontalInput() bool {
    return is.Left || is.Right
}

// HasVerticalInput returns true if player is jumping/flying
func (is InputState) HasVerticalInput() bool {
    return is.Up
}
```

**Why this design:**
- Not a Raylib type
- Physics and game logic receive this, not raw keyboard input
- Easy to swap input sources (file playback, network, AI)
- Domain logic decoupled from input mechanism
- Helper methods (`HasMovementInput()`) avoid repeating conditionals
- `Sell` is separate from movement (doesn't count as active input)

#### World (`domain/world/world.go`)

World data structure:

```go
type World struct {
    GroundLevel float32
    Width       float32
    Height      float32
}

func (w *World) GetGroundLevel() float32
func (w *World) IsInBounds(x, y float32) bool
```

**Why this design:**
- Centralizes world parameters
- No rendering data (colors, textures)
- Extensible for future terrain, tiles, etc.
- Used by physics system for collision checks

---

### Adapter Layer (Framework Integration)

#### Input Adapter (`internal/adapters/input/raylib.go`)

Translates Raylib input to domain InputState:

```go
type RaylibInputAdapter struct{}

func (a *RaylibInputAdapter) ReadInput() input.InputState {
    return input.InputState{
        Left:  rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA),
        Right: rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD),
        Up:    rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW),
        Down:  rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS),
    }
}
```

**Why this design:**
- Single responsibility: Read Raylib keys, output platform-agnostic InputState
- All Raylib input code in one place
- Easy to add new input sources (just create a new adapter)
- Can be mocked for testing

#### Rendering Adapter (`internal/adapters/rendering/raylib.go`)

Takes domain entities and renders them with Raylib:

```go
type RaylibRenderer struct{}

func (r *RaylibRenderer) Render(game *engine.Game) {
    rl.BeginDrawing()
    rl.ClearBackground(rl.RayWhite)

    r.renderWorld(game.GetWorld())
    r.renderPlayer(game.GetPlayer())

    rl.EndDrawing()
}

// Window lifecycle
func (r *RaylibRenderer) InitWindow(width, height int32, title string)
func (r *RaylibRenderer) CloseWindow()
func (r *RaylibRenderer) WindowShouldClose() bool
func (r *RaylibRenderer) SetTargetFPS(fps int32)
func (r *RaylibRenderer) GetFrameTime() float32
```

**Why this design:**
- Complete abstraction of Raylib
- Takes domain entities (Game, Player, World) and renders them
- All window lifecycle management in one place
- No business logic (pure rendering)
- Easy to swap for different renderer (SDL, canvas, headless)

---

### Application Layer (Orchestration)

#### Main (`cmd/game/main.go`)

Orchestrates the game loop:

```go
func main() {
    // Create adapters
    renderer := rendering.NewRaylibRenderer()
    inputAdapter := input.NewRaylibInputAdapter()

    // Initialize Raylib (adapter responsibility)
    renderer.InitWindow(screenWidth, screenHeight, "Drill Game")
    defer renderer.CloseWindow()
    renderer.SetTargetFPS(targetFPS)

    // Create domain objects
    gameWorld := world.NewWorld(screenWidth, screenHeight, groundLevel)
    game := engine.NewGame(gameWorld)

    // Main loop
    for renderer.WindowShouldClose() == false {
        dt := renderer.GetFrameTime()

        // 1. Read input from adapter
        inputState := inputAdapter.ReadInput()

        // 2. Update domain
        err := game.Update(dt, inputState)
        if err != nil {
            slog.Error("Error during update", "error", err)
            break
        }

        // 3. Render via adapter
        renderer.Render(game)
    }
}
```

**Why this design:**
- Clear, linear flow: Input → Update → Render
- No game logic (orchestration only)
- Dependencies are explicit
- Easy to understand at a glance

---

## Design Principles

### 1. Separation of Concerns

Each layer has a single responsibility:

- **Domain**: Business logic (physics, game rules)
- **Adapters**: Framework integration (Raylib)
- **Application**: Orchestration (wiring it together)

### 2. Framework Independence

The domain layer has **zero framework dependencies**:

- No `import rl "github.com/gen2brain/raylib-go/raylib"`
- No Raylib types (rl.Vector2, rl.Color, etc.)
- All conversions happen in adapters

### 3. Testability

Core logic is fully testable without Raylib:

```go
// This runs WITHOUT rl.InitWindow()
go test ./internal/domain/physics/...

// 11 physics tests pass with zero Raylib dependency
TestApplyHorizontalMovement_Acceleration
TestApplyHorizontalMovement_MaxSpeed
TestApplyGravity_IncreasesDownwardVelocity
// ... etc
```

### 4. Value Types for Small Objects

Small types are values, not pointers:

```go
type Vec2 struct {
    X float32  // 8 bytes total
    Y float32
}

type AABB struct {
    X, Y          float32  // 16 bytes total
    Width, Height float32
}

// ✓ Good: Values for small types
player.Velocity = types.Vec2{X: 100, Y: 200}
player.AABB = types.NewAABB(0, 0, 64, 64)

// ✗ Bad: Pointers for small types
player.Velocity = &types.Vec2{X: 100, Y: 200}
```

**Why:** Small types (8-16 bytes) should be values:
- Faster on stack than heap allocation
- Better cache locality
- No nil pointer issues
- Go idiom (see time.Time, image.Point)
- Cheaper to copy than pointer indirection on modern CPUs

### 5. Direct Field Access for Simplicity

Use direct field access instead of getters/setters when appropriate:

```go
// ✓ Good: Direct field access
player.AABB.X += player.Velocity.X * dt
player.Velocity.Y += gravity * dt
player.OnGround = true

// ✗ Overly complex: Unnecessary indirection
player.SetPosition(player.GetPosition().Add(player.GetVelocity().Scale(dt)))
player.SetVelocity(player.GetVelocity().Add(Vec2{Y: gravity * dt}))
player.SetOnGround(true)
```

**Why:**
- Simpler code, easier to read
- Less boilerplate (no getter/setter methods)
- Better performance (no function call overhead)
- Go idiom: exported fields for simple data structures
- Still maintains encapsulation at package boundaries

### 6. Getters for External Access

Adapters access domain data through getters:

```go
// Domain provides getters
func (g *Game) GetWorld() *world.World
func (g *Game) GetPlayer() *entities.Player

// Adapter reads via getters (doesn't modify)
func (r *RaylibRenderer) Render(game *engine.Game) {
    world := game.GetWorld()
    player := game.GetPlayer()
    // ... render
}
```

**Why:**
- Encapsulation (adapter can't modify game state)
- Clear data flow (one-way: domain → adapter)
- Easy to add state management later

---

## AABB Collision System

The game uses **Axis-Aligned Bounding Box (AABB) collision detection** with axis-separated resolution for precise 2D platformer physics.

### Core Concepts

**AABB Primitive:**
- Rectangular collision box defined by position (X, Y) and dimensions (Width, Height)
- Axis-aligned (no rotation) for fast intersection tests
- Used for both player and tiles

**Axis-Separated Resolution:**
- X-axis movement and collision resolved first
- Y-axis movement and collision resolved second
- Prevents corner-catching and enables natural wall sliding

### Collision Pipeline

```go
// 1. Apply movement and gravity
player.Velocity = ApplyHorizontalMovement(player.Velocity, input, dt)
player.Velocity = ApplyGravity(player.Velocity, dt)

// 2. X-axis: integrate → detect → resolve
player.AABB.X += player.Velocity.X * dt
collisionsX := CheckCollisions(player.AABB, world)
player.AABB, player.Velocity = ResolveCollisionsX(player.AABB, player.Velocity, collisionsX)

// 3. Y-axis: integrate → detect → resolve
player.AABB.Y += player.Velocity.Y * dt
collisionsY := CheckCollisions(player.AABB, world)
player.AABB, player.Velocity, player.OnGround = ResolveCollisionsY(player.AABB, player.Velocity, collisionsY)
```

### Collision Detection

**CheckCollisions()** finds all solid tiles overlapping the player:

```go
func CheckCollisions(aabb AABB, world *World) []TileCollision {
    // 1. Calculate which tiles the AABB might overlap
    minX, maxX, minY, maxY := GetOccupiedTileRange(aabb, TileSize)

    // 2. Check each potentially overlapping tile
    for x := minX; x <= maxX; x++ {
        for y := minY; y <= maxY; y++ {
            tile := world.GetTileAtGrid(x, y)
            if tile != nil && tile.IsSolid() && aabb.Intersects(tile.GetAABB(x, y, TileSize)) {
                // Found collision!
            }
        }
    }
}
```

**Performance:** Player can overlap at most 4 tiles (2×2 grid), so maximum 4 intersection tests per frame.

### Collision Resolution

**ResolveCollisionsX()** pushes player out horizontally:
- Calculates penetration depth using `AABB.Penetration()`
- Adjusts position: `aabb.X -= dx`
- Zeros horizontal velocity on wall hit

**ResolveCollisionsY()** pushes player out vertically:
- Calculates penetration depth
- Adjusts position: `aabb.Y -= dy`
- Detects ground: if pushed up (`dy > 0`), set `OnGround = true`
- Detects ceiling: if pushed down (`dy < 0`), zero upward velocity

### Why Axis-Separated?

**Without axis separation (naive AABB):**
- Player moving diagonally into corner gets "stuck"
- Cannot slide along walls smoothly
- Ground detection is ambiguous

**With axis separation:**
- X collision resolved first, Y collision resolved second
- Player slides along walls naturally during diagonal movement
- Clear ground/ceiling/wall detection based on which axis had collision

### Penetration Calculation

```go
func (a AABB) Penetration(b AABB) (dx, dy float32) {
    // Calculate overlap on each axis
    overlapX := min(a.X+a.Width, b.X+b.Width) - max(a.X, b.X)
    overlapY := min(a.Y+a.Height, b.Y+b.Height) - max(a.Y, b.Y)

    // Determine push direction based on relative positions
    if a.X < b.X {
        dx = overlapX  // Push left (subtract to move right)
    } else {
        dx = -overlapX // Push right (subtract to move left)
    }

    // Same for Y axis
    // ...
}
```

**Key insight:** Signs are chosen so `position -= penetration` always pushes objects apart.

---

## Camera System

The game uses Raylib's `Camera2D` for viewport management, allowing the player to explore a world much larger than the screen.

### Camera Implementation

**Camera2D** lives in the rendering adapter (not domain):

```go
type RaylibRenderer struct {
    camera       rl.Camera2D
    screenWidth  float32
    screenHeight float32
    worldWidth   float32
}

func (r *RaylibRenderer) updateCamera(player *entities.Player, w *world.World) {
    // Camera target follows player center
    playerCenterX := player.AABB.X + player.AABB.Width/2
    playerCenterY := player.AABB.Y + player.AABB.Height/2

    // Clamp camera to world bounds
    halfScreenW := r.screenWidth / 2
    halfScreenH := r.screenHeight / 2

    minX := halfScreenW
    maxX := r.worldWidth - halfScreenW
    minY := w.GetGroundLevel() - halfScreenH

    // Clamp and assign to camera target
    targetX := clamp(playerCenterX, minX, maxX)
    targetY := clamp(playerCenterY, minY, maxY)
    r.camera.Target = rl.Vector2{X: targetX, Y: targetY}
}

func (r *RaylibRenderer) Render(game *engine.Game, inputState input.InputState) {
    r.updateCamera(game.GetPlayer(), game.GetWorld())

    rl.BeginDrawing()
    rl.ClearBackground(rl.RayWhite)

    // World space rendering (camera applied)
    rl.BeginMode2D(r.camera)
    r.renderWorld(game.GetWorld())
    r.renderTiles(game.GetWorld())
    r.renderPlayer(game.GetPlayer())
    rl.EndMode2D()

    // Screen space rendering (no camera, UI always visible)
    r.renderDebugInfo(game.GetPlayer(), inputState)

    rl.EndDrawing()
}
```

### Viewport Culling

For performance with large worlds, only tiles visible in the camera viewport are rendered:

```go
func (r *RaylibRenderer) renderTiles(w *world.World) {
    tiles := w.GetAllTiles()

    // Calculate visible tile range
    minVisibleX := int((r.camera.Target.X - r.screenWidth/2) / world.TileSize) - 1
    maxVisibleX := int((r.camera.Target.X + r.screenWidth/2) / world.TileSize) + 1
    minVisibleY := int((r.camera.Target.Y - r.screenHeight/2) / world.TileSize) - 1
    maxVisibleY := int((r.camera.Target.Y + r.screenHeight/2) / world.TileSize) + 1

    for coord, tile := range tiles {
        gridX, gridY := coord[0], coord[1]

        // Skip tiles outside viewport
        if gridX < minVisibleX || gridX > maxVisibleX ||
           gridY < minVisibleY || gridY > maxVisibleY {
            continue
        }

        // Render visible tile...
    }
}
```

**Performance:** Reduces tiles rendered from ~94,000 to ~300 (~300× improvement).

### Why in Adapter, Not Domain?

- Camera is a **rendering concern**, not game logic
- Tightly coupled to Raylib's `Camera2D` struct
- Player position and world bounds already in domain (no new logic)
- Follows pattern: adapters translate domain state to visual representation

---

## World Boundary Constraints

Players cannot leave the game area. World boundaries are enforced by the physics system:

```go
func (ps *PhysicsSystem) constrainPlayerToWorldBounds(player *entities.Player) {
    // Horizontal: player.X must be in [0, worldWidth - playerWidth]
    minX := float32(0.0)
    maxX := ps.world.Width - float32(entities.PlayerWidth)

    if player.AABB.X < minX {
        player.AABB.X = minX
        player.Velocity.X = 0
    } else if player.AABB.X > maxX {
        player.AABB.X = maxX
        player.Velocity.X = 0
    }

    // Vertical: player.Y must be >= 0
    minY := float32(0.0)

    if player.AABB.Y < minY {
        player.AABB.Y = minY
        player.Velocity.Y = 0
    }
    // No maximum Y - player can dig infinitely deep
}
```

Called after collision resolution in the physics pipeline:

```go
// 2. Axis-separated collision resolution
player.AABB.X += player.Velocity.X * dt
collisionsX := physics.CheckCollisions(player.AABB, ps.world)
player.AABB, player.Velocity = physics.ResolveCollisionsX(player.AABB, player.Velocity, collisionsX)

player.AABB.Y += player.Velocity.Y * dt
collisionsY := physics.CheckCollisions(player.AABB, ps.world)
player.AABB, player.Velocity, player.OnGround = physics.ResolveCollisionsY(player.AABB, player.Velocity, collisionsY)

// 3. Enforce world boundary constraints
ps.constrainPlayerToWorldBounds(player)
```

**Design note:** Boundary constraints are purely domain-level (physics system), not rendering. Camera clamping happens independently in the adapter, preventing the camera from showing off-screen areas.

---

## World Dimensions

The game world extends far beyond the screen:

| Dimension | Size | Tiles |
|-----------|------|-------|
| Width | 7680 pixels | 120 tiles wide (6× screen width) |
| Height | 64000 pixels | 1000 tiles deep |
| Ground Level | 640 pixels | 10 tiles up from bottom |
| Tile Size | 64×64 pixels | Standard |

**Sparse tile storage:** Only non-empty tiles are stored in memory, enabling efficient large worlds.

---

## Future Architecture Considerations

### Adding New Entities

To add an Enemy that also uses physics:

```go
// 1. Create entity with AABB
type Enemy struct {
    AABB     types.AABB
    Velocity types.Vec2
    Health   float32
    AI       AIState
    // ...
}

// 2. Create a separate UpdateEnemyPhysics method or generalize UpdatePhysics
func (ps *PhysicsSystem) UpdateEnemyPhysics(enemy *entities.Enemy, dt float32) {
    // Same collision logic as player
    enemy.Velocity = physics.ApplyGravity(enemy.Velocity, dt)

    enemy.AABB.X += enemy.Velocity.X * dt
    collisionsX := physics.CheckCollisions(enemy.AABB, ps.world)
    enemy.AABB, enemy.Velocity = physics.ResolveCollisionsX(enemy.AABB, enemy.Velocity, collisionsX)

    enemy.AABB.Y += enemy.Velocity.Y * dt
    collisionsY := physics.CheckCollisions(enemy.AABB, ps.world)
    enemy.AABB, enemy.Velocity, _ = physics.ResolveCollisionsY(enemy.AABB, enemy.Velocity, collisionsY)
}
```

### Swapping Renderers

To use SDL instead of Raylib:

```go
// Create SDL adapter with same interface
type SDLRenderer struct {
    window *sdl.Window
    // ...
}

func (r *SDLRenderer) Render(game *engine.Game) {
    // SDL drawing logic
}

// Swap in main.go
renderer := sdl.NewSDLRenderer()  // Instead of Raylib
// Rest of main.go works unchanged!
```

### Adding Input Sources

To add file-based replay input:

```go
// New adapter with same interface
type FileInputAdapter struct {
    frames []InputState
    index  int
}

func (a *FileInputAdapter) ReadInput() input.InputState {
    state := a.frames[a.index]
    a.index++
    return state
}

// Swap in main.go
inputAdapter := file.NewFileInputAdapter("replay.bin")
// Game loop works unchanged!
```

---

## Testing Strategy

### Unit Tests (Domain Logic)

Test pure functions without framework:

```bash
# Physics tests (11 tests, no Raylib required)
go test ./internal/domain/physics/...

# Tests cover:
# - Movement (acceleration, damping, max speed)
# - Gravity (falling, velocity integration)
# - Collision (AABB detection, axis-separated resolution, wall/ceiling/ground)
```

### Integration Tests (Systems)

Test systems working together:

```go
// Example: Test player movement with physics
world := domain.NewWorld(1280, 720, 600)
player := domain.NewPlayer(640, 500)
physics := domain.NewPhysicsSystem(world)

inputState := domain.InputState{Right: true}
physics.UpdatePhysics(player, inputState, 0.016)

// Assert player moved right
assert.True(player.GetVelocityVec().X > 0)
```

### Manual Testing

Rendering and feel testing requires manual play:

- Gameplay feel (movement responsiveness)
- Visual polish (animations, particles)
- Performance (frame rates, memory)

---

## Performance Considerations

### Current (Phase 1)

- Simple physics (position + velocity)
- No spatial partitioning yet
- Direct collision checks
- Frame-independent movement via delta time

### Future Optimizations

- **Spatial Partitioning**: Grid or quadtree for collision queries
- **Object Pooling**: Reuse frequently created objects
- **Batch Rendering**: Group draw calls
- **Chunk Loading**: Only simulate/render visible area

---

## Physics Constants

Physics tuning values are split between two locations:

**Fixed constants** (`internal/domain/physics/constants.go`):

```go
const (
    Gravity             = 800     // pixels/sec² - downward acceleration
    MoveDamping         = 1000    // pixels/sec² - how fast player slows down
    FlyDamping          = 300     // pixels/sec² - air resistance when flying
    FallDamageThreshold = 500.0   // pixels/sec - minimum downward speed for damage
    FallDamageDivisor   = 20.0    // damage scaling: (speed - threshold) / divisor
)
```

**Dynamic values** (`internal/domain/entities/engine.go`):

Movement stats are defined per engine upgrade tier via named constructors:

| Stat | Base | Mk5 (Max) |
|------|------|-----------|
| MaxMoveSpeed | 450 px/s | 750 px/s |
| MoveAcceleration | 2500 px/s² | 5000 px/s² |
| FlyAcceleration | 2500 px/s² | 5000 px/s² |
| MaxUpwardVelocity | -600 px/s | -1000 px/s |

**How these affect gameplay:**
- **Gravity=800**: Heavy downward pull (1.25× Earth gravity) makes falling quick
- **MoveDamping=1000**: Tight ground control (stops in 0.45s)
- **FallDamageThreshold=500**: Small falls (under 500 px/sec) are safe
- **FallDamageDivisor=20**: Scales impact speed into damage (500+ px/sec → 0+ damage points)
- **Engine upgrades**: Better engines allow faster movement and climbing

See `internal/domain/physics/constants.go` and component files (`engine.go`, `hull.go`, `fuel_tank.go`) for source of truth.

---

## Game Configuration Reference

### Window & Display (`cmd/game/main.go`)

| Setting | Value | Purpose |
|---------|-------|---------|
| Screen Width | 1280 pixels | Horizontal viewport |
| Screen Height | 720 pixels | Vertical viewport |
| Target FPS | 60 | Frame rate cap |
| Ground Level | 640.0 pixels | Safe spawning elevation (10 tiles up) |

### Player Configuration (`internal/domain/entities/player.go`)

| Property | Value | Notes |
|----------|-------|-------|
| Size | 64×64 pixels | Matches tile size |
| Start Position | (640, 576) | Center X, just above ground |
| Inventory Types | 7 ore types | Copper, Iron, Silver, Gold, Mythril, Platinum, Diamond |
| Initial Cargo Capacity | 10 ore | Upgradeable to 75 ore |
| Initial Money | $1,000,000 | For development; earned by selling ores in game |
| Initial Fuel | 10.0 liters | Full tank (upgradeable to 65L) |
| Initial Health | 10.0 HP | Full HP (upgradeable to 75 HP) |
| Initial Upgrades | All Base | Engine, Hull, FuelTank, CargoHold at level 0 |

**Movement stats** (upgradeable via Engine):
| Stat | Base | Max (Mk5) |
|------|------|-----------|
| Max Move Speed | 450 px/sec | 750 px/sec |
| Fly Speed | 600 px/sec | 1000 px/sec |

### Shop Configuration (`internal/domain/entities/shop.go`)

| Property | Value | Purpose |
|----------|-------|---------|
| Position | (960, 576) | 3 tiles right of player spawn, ground level |
| Size | 320×192 pixels | 5 tiles wide × 3 tiles tall |
| Appearance | Forest green rect with dark border | Visual identification |
| Interaction | E key to sell | Triggers inventory sale |

### Ore Values

| Ore | Value | Depth Preference |
|-----|-------|------------------|
| Copper | $10 | 50-100 tiles |
| Iron | $25 | 100-150 tiles |
| Silver | $75 | 150-250 tiles |
| Gold | $250 | 250-350 tiles |
| Mythril | $1000 | 350-450 tiles |
| Platinum | $5000 | 450-550 tiles |
| Diamond | $30000 | 550+ tiles |

**Distribution:** Each ore type uses Gaussian distribution centered at depth preference. Rarer ores appear deeper and are worth more.

### Fuel System (`internal/domain/systems/fuel.go`)

| Setting | Value | Formula |
|---------|-------|---------|
| Base Tank Capacity | 10.0 liters | Upgradeable to 65L via FuelTank |
| Active Consumption | 0.33333 L/s | Base tank depletes in 30s with active input |
| Idle Consumption | 0.08333 L/s | Base tank depletes in 120s with no input |
| Active Input Triggers | Left, Right, Up, Dig | Movement/digging inputs only |

**Consumption Behavior:**
- Holding movement keys (Left/Right/Up) or digging (Down/S) = active mode
- Pressing Sell (E) does NOT trigger active consumption
- No movement for 1+ frame = idle consumption applies

### Controls & Input Mapping

| Input | Action | Notes |
|-------|--------|-------|
| **Left** (A or ←) | Move left / Dig left | Dig only when grounded against wall |
| **Right** (D or →) | Move right / Dig right | Dig only when grounded against wall |
| **Up** (W or ↑) | Jump/Fly | Hold to fly continuously |
| **Dig** (S or ↓) | Dig downward | Always available, snaps to grid |
| **Interact** (E) | Sell / Refuel | Sell at shop, refuel at station (AABB overlap) |

**Digging Behavior:**
- **Downward (S/Down)**: Always available, auto-aligns player to tile grid
- **Horizontal (A/D or Left/Right)**: Only when grounded, auto-digs blocking tiles

---

## Dependencies

Minimal, intentional dependencies:

- `github.com/gen2brain/raylib-go/raylib` - Graphics/audio (adapters only)
- `github.com/stretchr/testify` - Testing utilities (optional, use as needed)
- Go standard library - Everything else

**Strict rule:** No domain code imports Raylib.

---

## Summary

This architecture achieves:

- ✅ **Testability**: Domain logic 100% testable without framework
- ✅ **Portability**: Could swap Raylib for any renderer
- ✅ **Maintainability**: Clear responsibilities, linear data flow
- ✅ **Extensibility**: Easy to add new entities, systems, input sources
- ✅ **Clarity**: Folder structure = architecture diagram

The core principle: **Domain stays pure, framework stays outside.**
