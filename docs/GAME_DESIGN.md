# Game Design

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
- **Space / Mouse**: Activate drill
- **ESC**: Pause / Open shop

### Vehicle Mechanics
- Gravity pulls vehicle downward
- Can move left/right freely
- Drilling destroys tiles in front of vehicle
- Taking damage from heat, collisions, or hazards
- Cargo fills up as ores are collected

## World

### Structure
- Infinite vertical depth (procedurally generated)
- Fixed width (e.g., 200 tiles wide)
- Surface area with shop and landing pad
- Tiles become harder to drill with depth

### Tile Types
- **Empty**: No collision, can move through
- **Dirt**: Easy to drill, no value
- **Stone**: Harder to drill, no value
- **Ore Tiles**: Contain valuable ores when destroyed

### Ore Types & Rarity

Ores become progressively rarer and more valuable with depth:

| Ore Type    | Depth Range  | Value | Color      | Rarity |
|-------------|--------------|-------|------------|--------|
| Iron        | 0-100m       | $10   | Gray       | Common |
| Copper      | 50-200m      | $25   | Orange     | Common |
| Silver      | 100-400m     | $50   | Silver     | Uncommon |
| Gold        | 200-600m     | $100  | Gold       | Uncommon |
| Platinum    | 400-1000m    | $250  | White      | Rare |
| Emerald     | 600-1500m    | $500  | Green      | Rare |
| Ruby        | 1000-2500m   | $1000 | Red        | Very Rare |
| Diamond     | 2000-5000m   | $2500 | Blue       | Very Rare |
| Unobtainium | 5000m+       | $10000| Purple     | Legendary |

*Note: Depth ranges overlap to create smooth progression*

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
- **Fuel Capacity**: Longer operation time (future feature)

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
