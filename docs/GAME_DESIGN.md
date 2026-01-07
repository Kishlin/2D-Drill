# Game Design

> For implementation details, physics constants, and configuration values, see [ARCHITECTURE.md](ARCHITECTURE.md).

## Core Concept

A 2D vertical mining game inspired by Motherload. Players control a small drilling vehicle on a planet rich with ores. The core loop is simple but addictive: mine ores, return to surface, sell for currency, upgrade your vehicle, and venture deeper for rarer treasures.

## Game Loop

```
┌─────────────┐
│   Descend   │ ──> Dig deeper into the planet
└──────┬──────┘
       │
       v
┌─────────────┐
│  Mine Ores  │ ──> Collect valuable resources
└──────┬──────┘
       │
       v
┌─────────────┐
│   Ascend    │ ──> Return to surface (don't die!)
└──────┬──────┘
       │
       v
┌─────────────┐
│  Sell Ores  │ ──> Convert to currency
└──────┬──────┘
       │
       v
┌─────────────┐
│   Upgrade   │ ──> Improve your vehicle
└──────┬──────┘
       │
       └──────> Repeat (descend deeper)
```

## Player Vehicle

### Base Stats
- **Speed**: Horizontal movement speed
- **Drill Power**: How fast you can drill through tiles
- **Hull Strength**: How much damage you can take
- **Heat Resistance**: Protection from increasing temperature with depth
- **Cargo Capacity**: How much ore you can carry

### Controls
- **Arrow Keys / WASD**: Move vehicle
  - **Left (A) / Right (D)**: When grounded against a wall, automatically dig through the blocking tile
  - **Up (W)**: Fly/jump upward
- **Down (S) / Down Arrow**: Dig downward (with player grid alignment)
- **E**: Interact
  - At market: Sell entire inventory
  - At fuel station: Refuel tank (if affordable)
  - At hospital: Heal to full HP (if affordable)
  - At upgrade shop: Buy next upgrade tier (if affordable)

### Vehicle Mechanics
- Gravity pulls vehicle downward
- Can move left/right freely
- **Directional Drilling**:
  - Dig **downward** with explicit S/Down key press (snaps player to grid)
  - Dig **left/right** by moving into a wall while grounded (automatic, no grid snap)
- **Ore Collection**: Ores are automatically collected into inventory when dug
  - Each tile dug = 1 ore collected (1:1 ratio)
  - Dirt tiles are destroyed but not collected
  - Inventory displays counts for all 6 ore types in real-time
  - Collection respects cargo hold capacity (ore lost if cargo full)
- **Cargo Management**: Cargo hold limits total ore you can carry per trip
  - Base capacity: 10 ore (upgradeable to 75)
  - When full, newly dug ore is lost
  - Must return to surface and sell to make room
- Taking damage from heat, collisions, or hazards

### Directional Drilling & Animation

Both vertical and horizontal drilling feature smooth variable-duration animations based on depth and ore type:
- **Dirt at ground level**: 0.8 seconds
- **Dirt at max depth**: 30 seconds (linear scaling with depth)
- **Ore multipliers**: Copper 1.2x, Iron 1.5x, Gold 1.8x, Mythril 2.1x, Platinum 2.5x, Diamond 3.0x

The player moves toward the tile's center during the animation. The tile is only removed when the animation completes, then ore is collected.

**Downward Drilling (S/Down Key):**
- **Availability**: Can start anytime (must be grounded)
- **Animation**: Player moves to tile center (X-axis) and bottom edge (Y-axis) over variable duration (0.8-30+ seconds)
- **Completion**: Tile removed, ore collected if cargo permits
- **Effect**: Player is locked in animation; no other inputs processed

**Left/Right Drilling (A/D or Arrow Keys when Grounded):**
- **Availability**: Only when player is on solid ground (grounded against wall)
- **Animation**: Player moves to tile center (X-axis) while staying at ground level (Y-axis) over variable duration (0.8-30+ seconds)
- **Completion**: Tile removed, ore collected if cargo permits
- **Effect**: Player is locked in animation; no other inputs processed
- **Mid-Air Disabled**: Left/Right drilling blocked while airborne; player bounces off walls instead

**During Animation:**
- Fuel consumption continues (active rate if drilling, idle otherwise)
- Heat damage continues (based on depth and resistance)
- Fall damage does not apply (physics movement skipped)
- All other interactions blocked (market, upgrade, healing)
- Animation cannot be interrupted or cancelled

### Ore Inventory System
- **Automatic Collection**: When any ore tile is dug, it's automatically added to the player's inventory
- **Storage**: Inventory tracks count of each ore type (Copper, Iron, Gold, Mythril, Platinum, Diamond)
- **Display**: Current ore counts shown in debug overlay at top-left of screen
- **Simple Economy**: 1 tile dug = 1 ore collected (no partial ores, no quantity variance)
- **Dirt Ignored**: Only ore tiles contribute to inventory (dirt is destroyed but not collected)

### Currency & Market System
- **Market Location**: Visible on the surface (green rectangle, ~3 tiles right of spawn)
- **Selling Ores**: Press E while overlapping with market to sell entire inventory
- **Instant Transaction**: All ores converted to money immediately, inventory cleared
- **Ore Values**:
  | Ore      | Value |
  |----------|-------|
  | Copper   | $25   |
  | Iron     | $75   |
  | Gold     | $300  |
  | Mythril  | $1500 |
  | Platinum | $10000|
  | Diamond  | $30000|
- **Money Display**: Current balance shown in debug overlay
- **Cargo Limit**: Carry capacity determined by cargo hold upgrades (base: 10 ore, max: 75 ore)

### Fuel System

Fuel is a limited resource that creates time pressure for each expedition. Base tank capacity is 10 liters (upgradeable to 65L via Fuel Tank upgrades) with consumption rates that vary based on activity level.

**Consumption Rates:**
- Active movement (moving/digging): 0.333 L/sec
- Idle (standing still): 0.0833 L/sec

**Future Mechanics** (not yet implemented):
- Game over or limitations when fuel reaches zero
- Fuel efficiency upgrades

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed fuel system implementation and configuration.

### Damage & Health System

Players start with 10 hit points (upgradeable to 75 HP via Hull upgrades). Taking damage from falls creates risk when mining deep, but hospitals provide recovery.

**Fall Damage:**
- **Threshold**: 500 px/sec downward velocity (small falls are safe)
- **Formula**: `damage = (velocity - 500) / 20`
- **Examples**:
  - 500 px/sec fall → 0 damage (safe landing)
  - 600 px/sec fall → 5 damage
  - 700 px/sec fall → 10 damage (lethal)
- **Clamping**: HP never goes below 0 (no negative health)

**Healing System:**
- **Hospital Location**: Visible on the surface (crimson rectangle, 5 tiles left of fuel station)
- **Interaction**: Press E while overlapping hospital to heal
- **Healing Cost**: $2 per hit point needed (rounded up)
  - 0 HP needed (already max) = $0 (no transaction)
  - 10 HP needed (at 0 HP) = $20
  - 2.8 HP needed (at 7.2 HP) = $6 (ceil of 5.6)
  - 0.1 HP needed (at 9.9 HP) = $1 (ceil of 0.2)
- **Instant Heal**: HP immediately restored to max (10.0) on successful transaction
- **Rejection**: Cannot heal if insufficient money (healing prevented, no partial transaction)

**Future Mechanics** (not yet implemented):
- Game over when HP reaches 0
- Invulnerability frames after respawn
- Multiple healing tiers (partial vs full healing)
- Healing over time consumables

## World

### Structure
- Infinite vertical depth (procedurally generated)
- Fixed width (e.g., 200 tiles wide)
- Surface area with shop and landing pad
- Tiles become harder to drill with depth

### Tile Types
- **Empty**: No collision, can move through (air pockets, caves)
- **Dirt**: Solid, diggable, no value (filler)
- **Ore**: Solid, diggable, contains valuable resources

### Ore Types & Distribution

Six ore types are distributed using Gaussian curves, creating depth-based progression. Copper appears near the surface, while Diamond is found at mid-to-deep depths but remains extremely rare. Each ore type has specific value and rarity.

**Game Design:**
- Early game: Copper and Iron provide quick income and skill practice (shallower, tighter distributions)
- Mid game: Gold and Mythril increase risk/reward as you venture deeper
- Late game: Platinum and Diamond are high-value targets requiring deeper exploration

See [ARCHITECTURE.md](ARCHITECTURE.md) for ore types table, values, distribution parameters, and depth preferences.

## Environmental Hazards

### Heat

Temperature increases linearly with depth, causing exponential damage when heat resistance is exceeded. Heat is the primary limiting factor for deep mining.

**Temperature System:**
- **Ground Level** (Y=640 pixels): 15°C base temperature (safe)
- **Max Depth** (Y=64,000 pixels): 350°C maximum temperature
- **Formula**: Linear interpolation between ground and max depth
- **No Damage**: Temperature never rises above ground level (Y < 640)

**Damage Formula:**
- **Threshold**: When temperature exceeds player's heat resistance
- **Scaling**: Exponential `damage = 0.5 * (excessHeat / 10)^1.5 * dt`
- **Example**: At 60°C resistance with 100°C temperature:
  - Excess heat = 40°C
  - Damage/sec = 0.5 * (40/10)^1.5 = 0.5 * 8 ≈ 4 HP/sec
  - At 30°C excess: 1 HP/sec
  - At 10°C excess: 0.125 HP/sec

**Heat Shield Upgrades:**

Each heat shield tier enables safe mining at progressively deeper zones. Must be purchased sequentially.

| Tier | Resistance | Cost | Safe Depth (px) | Damage At +50°C |
|------|------------|------|-----------------|-----------------|
| Base | 50°C | - | 0-6,600 | ~8 HP/sec |
| Mk1 | 90°C | $200 | 6,600-14,000 | ~8 HP/sec |
| Mk2 | 140°C | $500 | 14,000-23,500 | ~8 HP/sec |
| Mk3 | 190°C | $1,200 | 23,500-33,000 | ~8 HP/sec |
| Mk4 | 250°C | $3,000 | 33,000-44,500 | ~8 HP/sec |
| Mk5 | 320°C | $7,500 | 44,500-64,000 | Full depth safe |

**Shop Appearance:**
- **Location**: Right of cargo hold shop (360px spacing pattern)
- **Fill Color**: Orange Red `(255, 69, 0)`
- **Border Color**: Red `(255, 0, 0)`
- **Interaction**: Press E while overlapping to purchase next tier

**Design Rationale:**
- Exponential damage creates meaningful progression gates (can't skip upgrades)
- Each tier unlocks approximately 8,000px of safe mining depth
- Pricing balanced between Hull ($150-$8,000) and FuelTank ($100-$4,000)
- Heat becomes the limiting factor for endgame progression
- Temperature display in debug overlay shows current and safe resistance

### Pressure (Future)
- Hull takes damage at extreme depths without upgrades
- Creates risk/reward for deep diving

### Hazardous Tiles (Future)
- Lava pockets
- Gas pockets (explosive)
- Underground water (slows drilling)

## Upgrade System

### Overview

Three upgrade types are available, each with 6 tiers (Base + Mk1 through Mk5). Upgrades must be purchased in order at dedicated upgrade shops on the surface. Press E while overlapping an upgrade shop to purchase the next tier.

### Engine Upgrades

Improves movement speed, acceleration, and flying capability.

| Tier | Max Speed | Acceleration | Fly Accel | Max Upward | Cost |
|------|-----------|--------------|-----------|------------|------|
| Base | 450 px/s | 2500 px/s² | 2500 px/s² | 600 px/s | - |
| Mk1 | 475 px/s | 2667 px/s² | 2667 px/s² | 635 px/s | $100 |
| Mk2 | 500 px/s | 2833 px/s² | 2833 px/s² | 670 px/s | $300 |
| Mk3 | 525 px/s | 3000 px/s² | 3000 px/s² | 705 px/s | $750 |
| Mk4 | 562 px/s | 3250 px/s² | 3250 px/s² | 740 px/s | $1,500 |
| Mk5 | 600 px/s | 3500 px/s² | 3500 px/s² | 775 px/s | $5,000 |

### Hull Upgrades

Increases maximum hit points.

| Tier | Max HP | Cost |
|------|--------|------|
| Base | 10 | - |
| Mk1 | 15 | $150 |
| Mk2 | 20 | $400 |
| Mk3 | 30 | $1,000 |
| Mk4 | 45 | $2,500 |
| Mk5 | 75 | $8,000 |

**Note:** Upgrading hull does NOT auto-heal. Visit the hospital to restore HP to new maximum.

### Fuel Tank Upgrades

Increases fuel tank capacity.

| Tier | Capacity | Cost |
|------|----------|------|
| Base | 10L | - |
| Mk1 | 15L | $100 |
| Mk2 | 22L | $250 |
| Mk3 | 32L | $600 |
| Mk4 | 45L | $1,500 |
| Mk5 | 65L | $4,000 |

**Note:** Upgrading tank does NOT auto-refuel. Visit the fuel station to fill to new capacity.

### Cargo Hold Upgrades

Increases ore cargo capacity (maximum amount of ore you can carry per trip).

| Tier | Capacity | Cost |
|------|----------|------|
| Base | 10 ore | - |
| Mk1 | 14 ore | $125 |
| Mk2 | 18 ore | $350 |
| Mk3 | 24 ore | $800 |
| Mk4 | 31 ore | $2,000 |
| Mk5 | 40 ore | $6,000 |

**Behavior:**
- When cargo is full, newly dug ore is lost (no auto-drop mechanic)
- Player must return to surface and sell inventory to make room
- Encourages strategic trip planning based on current cargo capacity

### Heat Shield Upgrades

Increases heat resistance, allowing safe mining at greater depths. Heat shield is essential for endgame progression since temperature damage increases exponentially with depth.

| Tier | Resistance | Cost | Safe Depth |
|------|------------|------|-----------|
| Base | 50°C | - | 0-6,600px |
| Mk1 | 90°C | $200 | 6,600-14,000px |
| Mk2 | 140°C | $500 | 14,000-23,500px |
| Mk3 | 190°C | $1,200 | 23,500-33,000px |
| Mk4 | 250°C | $3,000 | 33,000-44,500px |
| Mk5 | 320°C | $7,500 | 44,500-64,000px |

**Mechanics:**
- Each upgrade increases heat resistance by 40-70°C
- Temperature increases by ~335°C from surface to max depth
- Exponential damage formula ensures upgrades are mandatory (not optional)
- Late-game resource bottleneck (requires income to progress deeper)

**Note:** Unlike fuel tank and hull upgrades, heat shield doesn't auto-apply to new max. Resistance immediately applies on purchase.

### Upgrade Shops

Five separate upgrade shops are located on the surface (right of the ore market), spaced 360 pixels apart:
- **Engine Shop** (Steel Blue): Engine upgrades
- **Hull Shop** (Dim Gray): Hull upgrades
- **Fuel Tank Shop** (Tomato): Fuel tank upgrades
- **Cargo Hold Shop** (Dark Violet): Cargo hold upgrades
- **Heat Shield Shop** (Orange Red): Heat shield upgrades

### Future Upgrades (Not Yet Implemented)

#### Drilling Upgrades
- **Drill Strength**: Drill through tiles faster
- **Drill Efficiency**: Reduced heat generation while drilling

#### Survivability Upgrades
- **Auto-Repair**: Slowly regenerate health over time

#### Quality of Life Upgrades
- **Ore Detector**: Highlight valuable ores on screen
- **Auto-Seller**: Automatically sell when inventory is full

## Progression Curve

### Early Game (Surface to 5000px / ~78 tiles)
- **Goal**: Learn mechanics, earn first upgrades
- **Primary Ores**: Copper, Iron
- **Focus**: Speed and drill power upgrades
- **Challenge**: Learning to navigate, fuel management

### Mid Game (5000-20000px / 78-312 tiles)
- **Goal**: Build up funds, explore efficiently
- **Primary Ores**: Gold, Mythril
- **Focus**: Heat resistance, cargo capacity
- **Challenge**: Deeper dives, temperature management

### Late Game (20000px+ / 312+ tiles)
- **Goal**: Max out upgrades, hunt for rare ores
- **Primary Ores**: Platinum, Diamond
- **Focus**: Max-tier heat shield, fuel efficiency
- **Challenge**: Extreme heat, finding rare Diamond deposits, long journeys

## UI/UX

### HUD Elements
- **Top-left**: Health bar, Heat meter
- **Top-right**: Depth indicator, Currency
- **Bottom**: Cargo capacity/inventory preview
- **Minimap**: (future) Small overview of nearby area

### Shop Interface
- Grid of upgrade cards
- Shows current tier, next tier cost
- Preview of stat improvements
- "Repair" button (restore health for cost)
- "Sell All Ores" button

### Visual Feedback
- Screen shake on collisions
- Particle effects when drilling
- Heat distortion/tint at high temperatures
- Ore sparkle effects
- Damage flash on vehicle

## Future Features

### Short-term
- Sound effects and music
- Particle systems for polish
- More ore varieties
- Achievement system

### Medium-term
- Save/load game state
- Multiple vehicle types (trade-offs)
- Random events (cave-ins, ore veins, etc.)
- Challenge modes (time attack, depth race)

### Long-term (Steam/Mobile Vision)
- **Daily Challenges**: Fixed seed, compete on leaderboard
- **Events**: Limited-time special ores or modifiers
- **Leaderboards**: Deepest dive, most earnings, fastest time
- **Cloud Saves**: Play across devices
- **Workshop Support**: Custom ore mods, vehicle skins
- **Multiplayer**: Co-op drilling or competitive races

## Balancing Philosophy

- **Risk vs Reward**: Deeper = more valuable, but more dangerous
- **Meaningful Choices**: Each upgrade tier should feel impactful
- **Smooth Progression**: Avoid hard walls or grinding
- **Skill Expression**: Good routing and heat management rewarded
- **Replayability**: Random generation, multiple valid strategies

## Inspirations

- **Motherload** (Flash): Core loop, depth-based progression
- **Steamworld Dig**: Polish, upgrade satisfaction
- **Terraria**: Mining feel, ore variety
- **Cookie Clicker**: Exponential progression, "one more run" appeal

---

*This is a living document. Design will evolve based on playtesting and feedback.*
