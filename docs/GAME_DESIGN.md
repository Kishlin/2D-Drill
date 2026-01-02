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
- **E**: Sell entire inventory at shop (when inside shop area)

### Vehicle Mechanics
- Gravity pulls vehicle downward
- Can move left/right freely
- **Directional Drilling**:
  - Dig **downward** with explicit S/Down key press (snaps player to grid)
  - Dig **left/right** by moving into a wall while grounded (automatic, no grid snap)
- **Ore Collection**: Ores are automatically collected into inventory when dug
  - Each tile dug = 1 ore collected (1:1 ratio)
  - Dirt tiles are destroyed but not collected
  - Inventory displays counts for all 7 ore types in real-time
- Taking damage from heat, collisions, or hazards
- Cargo fills up as ores are collected

### Directional Digging Details
- **Grounded Left/Right Digging**: When the player is on solid ground and presses left/right against a diggable wall, the tile is automatically destroyed, allowing the player to move through
- **Downward Digging**: Can be performed anytime with S/Down key, aligns player to tile grid horizontally
- **Mid-Air**: Left/Right digging is disabled while airborne; player will bounce off walls instead

### Ore Inventory System
- **Automatic Collection**: When any ore tile is dug, it's automatically added to the player's inventory
- **Storage**: Inventory tracks count of each ore type (Copper, Iron, Silver, Gold, Mythril, Platinum, Diamond)
- **Display**: Current ore counts shown in debug overlay at top-left of screen
- **Simple Economy**: 1 tile dug = 1 ore collected (no partial ores, no quantity variance)
- **Dirt Ignored**: Only ore tiles contribute to inventory (dirt is destroyed but not collected)

### Currency & Shop System
- **Shop Location**: Visible on the surface (green rectangle, ~3 tiles right of spawn)
- **Selling Ores**: Press E while overlapping with shop to sell entire inventory
- **Instant Transaction**: All ores converted to money immediately, inventory cleared
- **Ore Values**:
  | Ore      | Value |
  |----------|-------|
  | Copper   | $10   |
  | Iron     | $25   |
  | Silver   | $75   |
  | Gold     | $250  |
  | Mythril  | $1000 |
  | Platinum | $5000 |
  | Diamond  | $30000|
- **Money Display**: Current balance shown in debug overlay
- **No Carrying Limit**: Store unlimited ore, but must return to shop to convert to currency (future: cargo upgrades)

### Fuel System

Fuel is a limited resource that creates time pressure for each expedition. Tank capacity is 10 liters with consumption rates that vary based on activity level.

**Consumption Rates:**
- Active movement (moving/digging): 0.333 L/sec
- Idle (standing still): 0.0833 L/sec

**Future Mechanics** (not yet implemented):
- Game over or limitations when fuel reaches zero
- Refueling mechanic (shop or surface station)
- Fuel efficiency upgrades

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed fuel system implementation and configuration.

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

Seven ore types are distributed using Gaussian curves, creating depth-based progression. Copper appears near the surface, while Diamond only exists at extreme depths. Each ore type has specific value and rarity.

**Game Design:**
- Early game: Copper and Iron provide quick income and skill practice
- Mid game: Silver and Gold increase risk/reward as you venture deeper
- Late game: Mythril, Platinum, and Diamond are high-value targets for endgame progression

See [ARCHITECTURE.md](ARCHITECTURE.md) for ore types table, values, distribution parameters, and depth preferences.

## Environmental Hazards

### Heat
- Temperature increases with depth
- Causes damage over time if heat resistance is insufficient
- Formula: `temperature = baseTemp + (depth * heatMultiplier)`
- Visual feedback: Screen tint gets redder with heat

### Pressure (Future)
- Hull takes damage at extreme depths without upgrades
- Creates risk/reward for deep diving

### Hazardous Tiles (Future)
- Lava pockets
- Gas pockets (explosive)
- Underground water (slows drilling)

## Upgrade System

### Upgrade Categories

#### 1. Movement Upgrades
- **Engine Power**: Increase horizontal speed
- **Thruster Power**: Improve vertical movement (fight gravity)
- **Fuel Capacity**: Increase tank size beyond 10L (future upgrade)
- **Fuel Efficiency**: Reduce consumption rates (future upgrade)

#### 2. Drilling Upgrades
- **Drill Strength**: Drill through tiles faster
- **Drill Efficiency**: Reduced heat generation while drilling
- **Drill Radius**: Larger drilling area (future)

#### 3. Survivability Upgrades
- **Hull Plating**: Increase maximum health
- **Heat Shielding**: Resist higher temperatures
- **Auto-Repair**: Slowly regenerate health over time (expensive)

#### 4. Cargo Upgrades
- **Cargo Hold**: Carry more ore per trip
- **Ore Detector**: Highlight valuable ores on screen (future)
- **Magnetic Field**: Attract nearby ore (future)

### Upgrade Progression

Upgrades have multiple tiers with exponential pricing:

```
Tier 1: $100    (+10% improvement)
Tier 2: $250    (+20% improvement)
Tier 3: $600    (+30% improvement)
Tier 4: $1,500  (+40% improvement)
Tier 5: $4,000  (+50% improvement)
Tier 6: $10,000 (+60% improvement)
...exponential growth
```

## Progression Curve

### Early Game (0-500m)
- **Goal**: Learn mechanics, earn first upgrades
- **Ores**: Iron, Copper, Silver
- **Focus**: Speed and drill power upgrades
- **Challenge**: Learning to navigate, avoiding damage

### Mid Game (500-2000m)
- **Goal**: Optimize mining routes, build up funds
- **Ores**: Gold, Platinum, Emerald
- **Focus**: Heat resistance, cargo capacity
- **Challenge**: Managing heat, deeper dives

### Late Game (2000m+)
- **Goal**: Max out upgrades, find legendary ores
- **Ores**: Ruby, Diamond, Unobtainium
- **Focus**: Max-tier upgrades, efficiency
- **Challenge**: Extreme heat, long journeys, risk management

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
