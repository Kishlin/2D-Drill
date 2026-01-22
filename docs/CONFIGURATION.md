# Configuration System

This document covers the data-driven configuration architecture, config structs, levels system, and reference tables. For high-level architecture, see [ARCHITECTURE.md](ARCHITECTURE.md).

---

## Overview

The game uses a **data-driven configuration architecture** where all game parameters flow from config structs. This enables per-level customization without code changes.

**Key Principle:** No hardcoded values in domain code—everything flows from config.

---

## Config Package Structure

```
internal/domain/config/
├── game_config.go       # GameConfig aggregate with Validate()
├── world_config.go      # WorldConfig (dimensions, spawn, buildings)
├── player_config.go     # PlayerConfig (starting money, items, upgrades)
├── generation_config.go # GenerationConfig (ores, hazards, distributions)
├── upgrade_config.go    # UpgradeConfig (6 upgrade type tiers)
├── item_config.go       # ItemConfig (5 items with prices/effects)
├── boss_room_config.go  # BossRoomConfig (boss type, floor, dimensions)
└── component_stats.go   # EngineStats, HullStats, etc.
```

---

## GameConfig Aggregate

The root configuration struct:

```go
type GameConfig struct {
    World      WorldConfig      // Dimensions, spawn, building positions
    Player     PlayerConfig     // Starting money, items, upgrade tiers
    Generation GenerationConfig // Ore/hazard distributions, colors, values
    Upgrades   UpgradeConfig    // 6 upgrade types with price/stats per tier
    Items      ItemConfig       // 5 items with prices and effects
    Level      LevelConfig      // Level number and name
    BossRoom   *BossRoomConfig  // Optional boss room configuration
}

func (c *GameConfig) Validate() error {
    // Validates world config
    // Checks starting tiers don't exceed available tiers
    // Ensures at least one ore exists
    // Validates ore/hazard ID uniqueness
}
```

---

## WorldConfig

```go
type WorldConfig struct {
    Width          float32         // World width in pixels
    Height         float32         // World height in pixels
    GroundLevel    float32         // Y coordinate of ground level
    Seed           int64           // Procedural generation seed
    PlayerSpawn    PlayerSpawn     // Player starting position
    BuildingLayout BuildingLayout  // Building X positions
}

type PlayerSpawn struct {
    X float32
    Y float32
}

type BuildingLayout struct {
    HospitalX    float32
    FuelStationX float32
    MarketX      float32
    UpgradeShopX float32
    ItemShopX    float32
}
```

---

## PlayerConfig

```go
type PlayerConfig struct {
    StartingMoney    int              // Initial currency
    StartingItems    [5]int           // Initial item counts
    StartingUpgrades StartingUpgrades // Initial upgrade tiers
}

type StartingUpgrades struct {
    Engine     int  // 0 = Base, 1 = Mk1, etc.
    Hull       int
    FuelTank   int
    CargoHold  int
    HeatShield int
    Drill      int
}
```

---

## GenerationConfig

```go
type GenerationConfig struct {
    Empty        TileDistribution  // Air pocket distribution
    Dirt         TileDistribution  // Dirt distribution
    DirtHardness float32           // Drilling time multiplier
    Ores         []OreConfig       // Dynamic list of ores
    Hazards      []HazardConfig    // Dynamic list of hazards
}
```

### OreConfig

```go
type OreConfig struct {
    ID           string           // Unique identifier (e.g., "copper")
    Name         string           // Display name
    Value        int              // Sell price at market
    Hardness     float32          // Drilling time multiplier
    Distribution TileDistribution // Gaussian spawn parameters
    Color        [4]uint8         // RGBA for rendering
}
```

### HazardConfig

```go
type HazardConfig struct {
    ID            string           // Unique identifier
    Name          string           // Display name
    Drillable     bool             // false = impenetrable
    FixedDuration float32          // Fixed drill time (0 = use depth formula)
    OnDrillDamage float32          // Damage on drill completion
    Distribution  TileDistribution // Gaussian spawn parameters
    Color         [4]uint8         // RGBA for rendering
}
```

### TileDistribution

```go
type TileDistribution struct {
    PeakDepth     float32  // Depth where most common (Gaussian center)
    Sigma         float32  // Distribution spread
    MaxWeight     float32  // Maximum spawn weight
    SurfaceWeight float32  // Weight at surface (for Empty/Dirt)
    DeepWeight    float32  // Weight at max depth (for Empty/Dirt)
}
```

---

## UpgradeConfig

Generic upgrade tiers using Go generics:

```go
type UpgradeTier[T any] struct {
    Name  string  // Tier name (e.g., "Base", "Mk1")
    Price int     // Purchase price (0 for base tier)
    Stats T       // Type-specific stats
}

type UpgradeConfig struct {
    Engines     []UpgradeTier[EngineStats]
    Hulls       []UpgradeTier[HullStats]
    FuelTanks   []UpgradeTier[FuelTankStats]
    CargoHolds  []UpgradeTier[CargoHoldStats]
    HeatShields []UpgradeTier[HeatShieldStats]
    Drills      []UpgradeTier[DrillStats]
}
```

### Component Stats

```go
type EngineStats struct {
    MaxSpeed        float32
    Acceleration    float32
    FlyAcceleration float32
    MaxUpwardSpeed  float32
}

type HullStats struct {
    MaxHP float32
}

type FuelTankStats struct {
    Capacity float32
}

type CargoHoldStats struct {
    Capacity int
}

type HeatShieldStats struct {
    HeatResistance float32
}

type DrillStats struct {
    DrillSpeed float32
}
```

---

## ItemConfig

```go
type ItemConfig struct {
    TeleportPrice int
    RepairPrice   int
    RefuelPrice   int
    BombPrice     int
    BigBombPrice  int
    BombRadius    int  // Tiles
    BigBombRadius int  // Tiles
}
```

---

## BossRoomConfig

```go
type BossRoomConfig struct {
    BossType    string    // Boss type identifier
    FloorType   FloorType // Concrete or Lava
    RoomHeight  float32   // Height of boss room in pixels
    FloorHeight float32   // Height of floor in tiles
}

type FloorType int

const (
    FloorConcrete FloorType = iota
    FloorLava
)
```

---

## Levels System

Levels are defined in `internal/domain/levels/` as functions returning complete `GameConfig`:

### Registry

```go
// levels/registry.go
func GetLevelConfig(levelNum int) (*config.GameConfig, error) {
    switch levelNum {
    case -2:
        return GetBossTestLevelConfig(), nil  // Boss testing
    case -1:
        return GetTestLevelConfig(), nil      // Development
    case 1:
        return GetLevel1Config(), nil         // Production
    default:
        return nil, fmt.Errorf("level %d not found", levelNum)
    }
}
```

### Level Definition Example

```go
// levels/level1.go
func GetLevel1Config() *config.GameConfig {
    return &config.GameConfig{
        World: config.WorldConfig{
            Width:       3072,
            Height:      64 * 800,
            GroundLevel: 640.0,
            Seed:        12345,
            PlayerSpawn: config.PlayerSpawn{X: 1536.0, Y: 570.0},
            BuildingLayout: config.BuildingLayout{
                HospitalX:    480.0,
                FuelStationX: 850.0,
                MarketX:      1400.0,
                UpgradeShopX: 1850.0,
                ItemShopX:    2220.0,
            },
        },
        Player: config.PlayerConfig{
            StartingMoney: 500,
            StartingItems: [5]int{0, 0, 0, 0, 0},
            StartingUpgrades: config.StartingUpgrades{
                Engine: 0, Hull: 0, FuelTank: 0,
                CargoHold: 0, HeatShield: 0, Drill: 0,
            },
        },
        Generation: config.GenerationConfig{...},
        Upgrades:   config.UpgradeConfig{...},
        Items:      config.ItemConfig{...},
        Level:      config.LevelConfig{Number: 1, Name: "Level 1"},
    }
}
```

### Special Levels

**Test Level (-1):** Development level with advanced player stats:
- Starting money: $100,000
- Starting items: 3 Teleports, 5 Repairs, 5 Refuels, 10 Bombs, 20 Big Bombs
- Max tier Engine and Drill
- Mid tier (Mk3) Hull, FuelTank, CargoHold, HeatShield

**Boss Test Level (-2):** For boss development and testing.

---

## Config Flow

```
main.go                    → levels.GetLevelConfig(levelNum)
                          ↓
GameConfig                → gameCfg.Validate()
                          ↓
engine.NewGame            → receives GameConfig, creates world/systems
                          ↓
entities.NewPlayerFromConfig → receives PlayerConfig, UpgradeConfig
                          ↓
renderer.NewWithConfig    → receives GenerationConfig (ore/hazard colors)
```

---

## Reference Tables

### Window & Display

| Setting | Value |
|---------|-------|
| Screen Width | 1280 pixels |
| Screen Height | 720 pixels |
| Target FPS | 60 |
| Ground Level | 640.0 pixels |

### Player

| Property | Value |
|----------|-------|
| Size | 64×64 pixels |
| Start Position | (1536, 570) |
| Initial Money | $500 (Level 1) |
| Initial Fuel | Full tank |
| Initial HP | Full HP |

### Ore Values (Level 1)

| Ore | Value | Peak Depth | Sigma | Max Weight | Hardness |
|-----|-------|------------|-------|------------|----------|
| Copper | $25 | -75px | 120 | 8.0 | 1.2 |
| Iron | $75 | 70px | 90 | 5.0 | 1.5 |
| Gold | $300 | 230px | 80 | 3.0 | 1.8 |
| Mythril | $1500 | 360px | 70 | 2.2 | 2.1 |
| Platinum | $10000 | 500px | 80 | 1.8 | 2.5 |
| Diamond | $30000 | 600px | 180 | 0.15 | 3.0 |

### Engine Upgrades

| Tier | Max Speed | Acceleration | Fly Accel | Max Upward | Cost |
|------|-----------|--------------|-----------|------------|------|
| Base | 450 px/s | 2500 px/s² | 2500 px/s² | 600 px/s | - |
| Mk1 | 475 px/s | 2667 px/s² | 2667 px/s² | 635 px/s | $100 |
| Mk2 | 500 px/s | 2833 px/s² | 2833 px/s² | 670 px/s | $300 |
| Mk3 | 525 px/s | 3000 px/s² | 3000 px/s² | 705 px/s | $750 |
| Mk4 | 562 px/s | 3250 px/s² | 3250 px/s² | 740 px/s | $1,500 |
| Mk5 | 600 px/s | 3500 px/s² | 3500 px/s² | 775 px/s | $5,000 |

### Hull Upgrades

| Tier | Max HP | Cost |
|------|--------|------|
| Base | 10 | - |
| Mk1 | 15 | $150 |
| Mk2 | 20 | $400 |
| Mk3 | 30 | $1,000 |
| Mk4 | 45 | $2,500 |
| Mk5 | 75 | $8,000 |

### Fuel Tank Upgrades

| Tier | Capacity | Cost |
|------|----------|------|
| Base | 10L | - |
| Mk1 | 15L | $100 |
| Mk2 | 22L | $250 |
| Mk3 | 32L | $600 |
| Mk4 | 45L | $1,500 |
| Mk5 | 65L | $4,000 |

### Cargo Hold Upgrades

| Tier | Capacity | Cost |
|------|----------|------|
| Base | 10 ore | - |
| Mk1 | 14 ore | $125 |
| Mk2 | 18 ore | $350 |
| Mk3 | 24 ore | $800 |
| Mk4 | 31 ore | $2,000 |
| Mk5 | 40 ore | $6,000 |

### Heat Shield Upgrades

| Tier | Resistance | Cost | Safe Depth |
|------|------------|------|-----------|
| Base | 50°C | - | 0-6,600px |
| Mk1 | 90°C | $200 | 6,600-14,000px |
| Mk2 | 140°C | $500 | 14,000-23,500px |
| Mk3 | 190°C | $1,200 | 23,500-33,000px |
| Mk4 | 250°C | $3,000 | 33,000-44,500px |
| Mk5 | 320°C | $7,500 | 44,500-64,000px |

### Drill Upgrades

| Tier | Drill Speed | Cost | Effect at Surface | Effect at Max Depth |
|------|-------------|------|-------------------|---------------------|
| Base | 1.0x | - | 1.0s → 1.0s | 24s → 24s |
| Mk1 | 2.0x | $125 | 1.0s → 0.91s | 24s → 12s |
| Mk2 | 3.0x | $350 | 1.0s → 0.83s | 24s → 8s |
| Mk3 | 4.0x | $875 | 1.0s → 0.77s | 24s → 6s |
| Mk4 | 5.0x | $2,000 | 1.0s → 0.71s | 24s → 4.8s |
| Mk5 | 6.0x | $6,500 | 1.0s → 0.67s | 24s → 4s |

### Item Prices

| Item | Key | Price | Effect |
|------|-----|-------|--------|
| Teleport | T | $500 | Return to spawn |
| Repair Kit | R | $200 | Restore HP to max |
| Fuel Can | F | $100 | Fill fuel to max |
| Bomb | B | $300 | Destroy tiles (2-tile radius) |
| Big Bomb | G | $800 | Destroy tiles (4-tile radius) |

### Fuel Consumption

| Activity | Rate | Tank Duration (Base) |
|----------|------|---------------------|
| Active (moving/drilling) | 0.333 L/s | 30 seconds |
| Idle | 0.0833 L/s | 120 seconds |

### Controls

| Input | Type | Action |
|-------|------|--------|
| A / ← | Continuous | Move left / Drill left |
| D / → | Continuous | Move right / Drill right |
| W / ↑ | Continuous | Jump/Fly |
| S / ↓ | Continuous | Drill down |
| E | Discrete | Interact with building |
| T | Discrete | Use Teleport |
| R | Discrete | Use Repair |
| F | Discrete | Use Refuel |
| B | Discrete | Use Bomb |
| G | Discrete | Use Big Bomb |
| Z | Discrete | Previous tab (shop) |
| X | Discrete | Next tab (shop) |
| Q / Esc | Discrete | Close shop |
