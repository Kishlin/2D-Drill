# Game Design

> **Related Documentation:**
> - [ARCHITECTURE.md](ARCHITECTURE.md) — High-level architecture overview
> - [SYSTEMS.md](SYSTEMS.md) — Game systems implementation details
> - [PHYSICS.md](PHYSICS.md) — Physics constants and damage formulas
> - [WORLD.md](WORLD.md) — World generation and tile distributions
> - [CONFIGURATION.md](CONFIGURATION.md) — Config structs and reference tables
>
> **Note:** All game values (ore prices, upgrade costs, drill speeds, hazard damage, etc.) are defined in level configuration files (`internal/domain/levels/`). Values shown in this document reflect Level 1 defaults and may vary per level.

## Core Concept

A 2D vertical mining game inspired by Motherload. Players control a small drilling vehicle on a planet rich with ores. The core loop is simple but addictive: mine ores, return to surface, sell for currency, upgrade your vehicle, and venture deeper for rarer treasures.

## Game Loop

```
┌─────────────┐
│   Descend   │ ──> Drill deeper into the planet
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
  - **Left (A) / Right (D)**: When grounded against a wall, automatically drill through the blocking tile
  - **Up (W)**: Fly/jump upward
- **Down (S) / Down Arrow**: Drill downward (with player grid alignment)
- **E**: Interact
  - At market: Open sell modal (E again to sell, Q to cancel)
  - At fuel station: Refuel tank (if affordable)
  - At hospital: Open heal modal (W/S to navigate, E to heal, Q to close)
  - At upgrade shop: Open modal, navigate with arrows, E to buy, Q to close
  - At item shop: Open modal, navigate with arrows, E to buy, Q to close
- **Item Keys** (press once to use if you have items):
  - **T**: Teleport to spawn point
  - **R**: Repair (restore HP to max)
  - **F**: Refuel (fill fuel tank to max)
  - **B**: Bomb (destroy tiles in small radius)
  - **G**: Big Bomb (destroy tiles in larger radius)

### Vehicle Mechanics
- Gravity pulls vehicle downward
- Can move left/right freely
- **Directional Drilling**:
  - Drill **downward** with explicit S/Down key press (snaps player to grid)
  - Drill **left/right** by moving into a wall while grounded (automatic, no grid snap)
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

Both vertical and horizontal drilling feature smooth variable-duration animations based on tile hardness and depth:

**Duration Formula:** `duration = baseTime × hardness × depthFactor / drillSpeed`

Where:
- **baseTime**: 1.0 second constant
- **hardness**: per-tile value (Dirt 1.0, Copper 1.2, Iron 1.5, Gold 1.8, Mythril 2.1, Platinum 2.5, Diamond 3.0)
- **depthFactor**: scales 1.0 at surface → 24.0 at max depth
- **drillSpeed**: from drill upgrades (1.0 → 6.0)

**Lava Exception:** Fixed 0.3s duration regardless of depth (damage is the penalty).

The player moves toward the tile's center during the animation. The tile is only removed when the animation completes, then ore is collected.

**Downward Drilling (S/Down Key):**
- **Availability**: Can start anytime (must be grounded)
- **Animation**: Player moves to tile center (X-axis) and bottom edge (Y-axis) over variable duration (1.0-24+ seconds based on depth/ore/drill)
- **Completion**: Tile removed, ore collected if cargo permits
- **Effect**: Player is locked in animation; no other inputs processed

**Left/Right Drilling (A/D or Arrow Keys when Grounded):**
- **Availability**: Only when player is on solid ground (grounded against wall)
- **Animation**: Player moves to tile center (X-axis) while staying at ground level (Y-axis) over variable duration (1.0-24+ seconds based on depth/ore/drill)
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
- **Selling Ores**: Press E while overlapping market to open modal UI
  - Modal displays ore inventory with names, counts, unit prices, and line totals
  - Shows grand total at bottom
  - Press E to sell all ore and close modal
  - Press Q/Escape to close without selling
  - If inventory is empty, displays "No ore to sell" message
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
- Active movement (moving/drilling): 0.333 L/sec
- Idle (standing still): 0.0833 L/sec

**Future Mechanics** (not yet implemented):
- Game over or limitations when fuel reaches zero
- Fuel efficiency upgrades

See [SYSTEMS.md](SYSTEMS.md#fuel-system-systemsfuelgo) for implementation details.

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
- **Interaction**: Press E while overlapping hospital to open healing modal
- **Modal Navigation**: W/S to navigate options, E to heal, Q to close
- **Healing Cost**: $2 per hit point (rounded up)

**Healing Options:**
| Option | HP Restored | Cost | Behavior |
|--------|-------------|------|----------|
| Restore 1 HP | min(1, hpNeeded) | $2 | Stays open (repeatable) |
| Restore 10 HP | min(10, hpNeeded) | $20 max | Stays open (repeatable) |
| Restore All HP | hpNeeded | ceil(hpNeeded × 2) | Closes after |
| Max Affordable | min(hpNeeded, floor(money/2)) | varies | Closes after |

**Display Details:**
- Fixed options (1 HP, 10 HP) show whole numbers unless capped below nominal
- Variable options (All, Max) show one decimal place (e.g., "+11.3 HP")
- Max Affordable shows "+0.0 HP for $0" when broke (greyed out)
- Unaffordable options are greyed with red cost text
- Selected affordable option shows green cost text

**Edge Cases:**
- At full HP: Shows "Already at full health" message
- Insufficient funds: Options greyed out, can still close with Q

**Future Mechanics** (not yet implemented):
- Game over when HP reaches 0
- Invulnerability frames after respawn
- Healing over time consumables

## World

### Structure
- Infinite vertical depth (procedurally generated)
- Fixed width (e.g., 200 tiles wide)
- Surface area with shop and landing pad
- Tiles become harder to drill with depth

### Tile Types
- **Empty**: No collision, can move through (air pockets, caves)
- **Dirt**: Solid, drillable, no value (filler)
- **Ore**: Solid, drillable, contains valuable resources
- **Rock**: Solid, impenetrable obstacle (cannot be drilled, only bombed)
- **Lava**: Solid, drillable hazard (fast 0.3s animation, deals damage on completion)

### Ore Types & Distribution

Ore types are defined per level in `GenerationConfig.Ores`. Each ore has an ID, display name, value, hardness, Gaussian distribution parameters, and color. Level 1 includes six ore types distributed using Gaussian curves, creating depth-based progression.

**Level 1 Ores:**
- **Copper** ($25) — Near surface, very common
- **Iron** ($75) — Shallow depth, common
- **Gold** ($300) — Mid-shallow, uncommon
- **Mythril** ($1500) — Mid-depth, rare
- **Platinum** ($10000) — Deep, very rare
- **Diamond** ($30000) — Mid-deep, extremely rare

**Game Design:**
- Early game: Copper and Iron provide quick income and skill practice (shallower, tighter distributions)
- Mid game: Gold and Mythril increase risk/reward as you venture deeper
- Late game: Platinum and Diamond are high-value targets requiring deeper exploration

**Customization:** Different levels can define entirely different ore sets (e.g., Level 2 might have Uranium instead of Copper). See [CONFIGURATION.md](CONFIGURATION.md) for config structure.

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

### Hazard Tiles

Hazard tiles are impenetrable and drillable obstacles that appear at deeper depths, creating new challenges and encouraging strategic bomb usage.

**Rock Tiles (Impenetrable Obstacles):**
- **Appearance**: Dark gray blocks
- **Drillability**: Cannot be drilled at all
- **Depth**: Appear starting ~40% depth, become common at 80%+
- **Interaction**: Block player movement and drilling attempts
- **Bomb Effect**: Destroyed by bombs (both sizes)
- **Strategy**: Use bombs to bypass or navigate around rock formations

**Lava Tiles (Drillable Hazards):**
- **Appearance**: Red-orange blocks
- **Drillability**: Drillable with fixed 0.3-second duration (from `HazardHardness` map, depth-independent)
- **Depth**: Appear starting ~60-65% depth, more common at 80%+
- **Damage**: Deals 100 damage on drilling completion (reduced to 50 with Mk5 heat shield)
- **Damage Formula**: `damage = 100 - (currentHeatResistance / 320.0 * 50)`
- **Strategy**: Use heat shield upgrades to reduce damage before drilling lava-rich areas

**Why Hazards Matter:**
- Hazards dominate terrain at 80%+ depths (creating natural progression gates)
- Rocks force strategic bomb usage for passages
- Lava incentivizes heat shield upgrades
- Both tile types make deep mining more challenging and tactical

### Pressure (Future)
- Hull takes damage at extreme depths without upgrades
- Creates risk/reward for deep diving

## Upgrade System

### Overview

Six upgrade types are available, each defined in `UpgradeConfig` with configurable tiers (prices and stats). Level 1 provides 6 tiers per type (Base + Mk1 through Mk5). All upgrades are purchased from a single **Unified Upgrade Shop** with a modal UI. Press **E** while overlapping the shop to open the modal. Use **Z/X** to cycle between upgrade categories (tabs), **arrows/WASD** to navigate the grid, and **E** to purchase. Press **Q** or **Escape** to close the shop.

**Key Features:**
- **No Sequential Requirement**: Buy any tier directly if you have the money (skip tiers if affordable)
- **Modal UI**: Opens a screen overlay that pauses all other gameplay
- **Grid Display**: Shows all tiers for comparison before purchasing
- **Visual Feedback**: Light cells = affordable, dark cells = too expensive
- **Tab Cycling**: Easy navigation between all 6 upgrade types
- **Config-Driven**: Prices and stats defined in level config, enabling per-level balancing

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

### Drill Upgrades

Improves drilling speed via a depth-scaled divisor. At surface, only 10% of the upgrade applies; at max depth, 100% applies. This design ensures upgrades feel impactful deep underground without trivializing surface drilling.

| Tier | Drill Speed | Cost | Effect at Surface | Effect at Max Depth |
|------|-------------|------|-------------------|---------------------|
| Base | 1.0x | - | 1.0s → 1.0s | 24s → 24s |
| Mk1 | 2.0x | $125 | 1.0s → 0.91s | 24s → 12s |
| Mk2 | 3.0x | $350 | 1.0s → 0.83s | 24s → 8s |
| Mk3 | 4.0x | $875 | 1.0s → 0.77s | 24s → 6s |
| Mk4 | 5.0x | $2,000 | 1.0s → 0.71s | 24s → 4.8s |
| Mk5 | 6.0x | $6,500 | 1.0s → 0.67s | 24s → 4s |

**Formula:** `effectiveDivisor = 1 + (drillSpeed - 1) * (0.1 + 0.9 * depthFactor)`

**Design Rationale:**
- Depth-scaled divisor prevents trivializing surface drilling
- Strong improvement at depth (24s → 4s with Mk5) makes deep mining viable
- Gentle improvement at surface (1.0s → 0.67s) maintains early game challenge
- Pricing balanced between Engine and Hull tiers

### Upgrade Shops

Six separate upgrade shops are located on the surface (right of the ore market), spaced 360 pixels apart:
- **Engine Shop** (Steel Blue): Engine upgrades
- **Hull Shop** (Dim Gray): Hull upgrades
- **Fuel Tank Shop** (Tomato): Fuel tank upgrades
- **Cargo Hold Shop** (Dark Violet): Cargo hold upgrades
- **Heat Shield Shop** (Orange Red): Heat shield upgrades
- **Drill Shop** (Dark Goldenrod): Drill upgrades

### Future Upgrades (Not Yet Implemented)

#### Survivability Upgrades
- **Auto-Repair**: Slowly regenerate health over time

#### Quality of Life Upgrades
- **Ore Detector**: Highlight valuable ores on screen
- **Auto-Seller**: Automatically sell when inventory is full

## Items

### Overview

Consumable items provide tactical advantages during deep mining expeditions. Item prices and effects (bomb radius) are defined in `ItemConfig`. Items are purchased from the unified Item Shop and used via single key presses. Items are one-time use and must be repurchased.

### Item Types & Effects

| Item | Key | Cost | Effect | Strategic Use |
|------|-----|------|--------|----------------|
| **Teleport** | T | $500 | Instantly return to spawn point at ground level | Emergency escape from danger (deep heat, low fuel) |
| **Repair Kit** | R | $200 | Instantly restore HP to maximum | Heal without visiting hospital |
| **Fuel Can** | F | $100 | Instantly fill fuel tank to maximum | Extend expedition range |
| **Bomb** | B | $300 | Destroy tiles in 2-tile radius circle (~13 tiles) | Quickly excavate around obstacle |
| **Big Bomb** | G | $800 | Destroy tiles in 4-tile radius circle (~49 tiles) | Clear large areas, create escape routes |

### Item Mechanics

**Using Items:**
- Press the item key (T, R, F, B, G) once to use
- Key must be pressed once per frame (held keys don't trigger repeats)
- Player must have at least 1 item of that type
- Item count decrements on successful use
- Display shows current item counts in debug overlay

**Blocked During Drilling:**
- Items cannot be used during drilling animations (consistent with other interactions)
- Finish drilling before using items

**Bomb Effects:**
- Bombs destroy all drillable tiles in blast radius
- Ore is lost (not collected) when destroyed by bombs
- Bombs ignore ore value—purely a terrain-clearing tool
- Useful for: bypassing obstacles, creating shortcuts, emergency escapes

**Player-Affecting Items (Teleport, Repair, Refuel):**
- Instantly apply effect (no cost, no confirmation)
- Effect bypasses normal systems (no hospital visit needed for repair, etc.)
- Teleport resets velocity and places player at spawn

### Item Shop (Unified Modal UI)

A single unified item shop is located on the surface, to the right of the upgrade shop. Press **E** while overlapping the shop to open a modal interface. Use **arrows or WASD** to navigate a 2×3 grid of 5 items (last cell empty), and **E** to purchase. Press **Q** or **Escape** to close and resume gameplay.

**Shop Location & Size:**
- Position: Right of upgrade shop with 40px gap (130px void between market and shops)
- Size: 320×192 pixels (same as all other buildings)
- Appearance: Blue-tinted building matching standard building size

**Modal Grid Layout:**
```
[Teleport]  [Repair]    [Refuel]
[Bomb]      [Big Bomb]  [Empty]
```

**Navigation & Interaction:**
- **Arrows/WASD**: Navigate through grid with wrapping
- **Empty Cell**: Automatically skipped during navigation
- **E Key**: Purchase currently selected item (if affordable)
- **Q or Escape**: Close modal and resume gameplay

**Purchasing:**
- If you have enough money, item count increases by 1
- If insufficient funds, purchase is blocked (no feedback)
- All 5 items visible and accessible from single shop

### Economy Implications

Items provide flexible power progression:
- **Early Game**: Affordable items (Refuel $100, Repair $200) help manage basic risk
- **Mid Game**: Bombs ($300) become cost-effective for clearing obstacles vs drilling
- **Late Game**: Teleport ($500) provides insurance against extreme depth penalties

Strategic Item Use:
- Bombs reduce drilling time in congested areas
- Teleport enables aggressive deep dives with escape plan
- Repair/Refuel items extend expedition length vs making return trips
- Mix of item usage and traditional methods balances spending

## Boss Fights

### Overview

Boss encounters are end-of-level challenges that occur at the bottom of the mineable world. Players dig through normal terrain to reach a boss room, then must defeat the boss to complete the level. Boss rooms feature solid, indestructible floor tiles that prevent digging further.

### Boss Room Layout

Each boss room has three layers:
1. **Mining Area** — Normal terrain with ores and hazards (top section)
2. **Boss Room** — Empty space where the boss waits (middle section, ~680 pixels)
3. **Floor** — Solid concrete or lava tiles (bottom section, indestructible)

The camera clamps at the world bottom to prevent viewing or nuking below the floor.

### Boss Mechanics

**Activation & Deactivation:**
- Boss activates when player enters the boss room (crosses threshold)
- Boss deactivates when player leaves the room (projectiles cleared)
- HP is preserved if player exits and re-enters

**Vulnerability System:**
- Bosses have `IsVulnerable()` method determining when damage is accepted
- Phase 1: Always vulnerable (easy mode)
- Phases 2-3: Only vulnerable during specific windows (after slam attacks)
- One bomb hit per vulnerability window (prevents spam)

**Contact Damage:**
- Physical bosses can deal contact damage (configurable per boss)
- TestBoss deals 20 HP/sec on player contact
- Some future bosses may have 0 contact damage (passable)

**Interaction:**
- Bombs deal damage to vulnerable bosses (10 HP per bomb, 25 HP per big bomb)
- Damage only applies during vulnerability windows
- Boss dies when HP reaches 0
- Victory screen appears on boss defeat
- Current game state (Playing/Victory/Defeat) is tracked and rendered

### TestBoss Behavior

**Stats:**
- 100 HP, 100×100 pixel sprite
- 80 px/s base movement speed (increases per phase)
- 20 HP/sec contact damage

**State Machine:**
1. **Patrol** — Moving left-right, shooting projectiles at player
2. **Windup** — Stopped, vibrating (1 second warning)
3. **Slam** — AOE damage zone (0.3 seconds, 150px radius, 15 damage)
4. **Vulnerable** — Immobile, can be damaged by bombs

**Attacks:**
- **Projectile Volley**: 3 projectiles aimed at player (200 px/s, 5 damage each)
- **Ground Slam**: AOE attack with telegraph warning, creates vulnerability window

**Phase System:**

| Phase | HP Range | Movement | Projectiles | Slam | Vulnerability |
|-------|----------|----------|-------------|------|---------------|
| 1 | 100-66% | 80 px/s | Every 3s | None | Always |
| 2 | 66-33% | 100 px/s | Every 2s | Every 6s | 3s after slam |
| 3 | 33-0% | 120 px/s | Every 1s | Every 4s | 2s after slam |

**Phase 3 Special**: 50% chance of double slam (slam → 0.4s pause → slam → vulnerable)

### Visual Feedback

**Boss States:**
- **Patrol (Vulnerable)**: Pink flashing
- **Patrol (Invulnerable)**: Gray tint
- **Windup**: Orange flashing + horizontal vibration
- **Slam**: Bright red + AOE damage circle
- **Vulnerable**: Pink flashing

**AOE Effects:**
- **Telegraph**: Pulsing yellow circle (warning)
- **Damage**: Solid orange-red circle

### Floor Types

Boss rooms can have different floor mechanics configured per level:

**Concrete Floor:**
- Solid, walkable surface
- Appears gray
- Safe to stand on indefinitely
- No special effects

**Lava Floor:**
- Solid, but deals damage while standing on it
- Appears orange/red
- Encourages quick movement or escape
- Damage rate configurable per level

### Design Notes

- Boss fights serve as level conclusion points
- Bomb weapons become tactical in boss encounters (timing vulnerability windows)
- Extensible via `Boss` interface (non-physical bosses possible, e.g., bullet-hell)
- Floor types allow varied boss mechanics (safe vs hazardous arenas)
- Game state tracking enables future pause/menu features
- Each boss has its own rendering logic (no generic animation interfaces)

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
- **Steamworld Dig Series**: Polish, upgrade satisfaction
- **Terraria**: Mining feel, ore variety
- **Cookie Clicker**: Exponential progression, "one more run" appeal

---

*This is a living document. Design will evolve based on playtesting and feedback.*
