# World System

This document covers world configuration, procedural generation, and tile mechanics. For high-level architecture, see [ARCHITECTURE.md](ARCHITECTURE.md).

---

## Overview

The world is a chunk-based procedurally generated environment with depth-dependent tile distribution. Players mine through terrain that becomes progressively more challenging with depth.

**Key Files:**
- `internal/domain/world/world.go` — World struct, chunk management
- `internal/domain/world/generator.go` — Procedural generation
- `internal/domain/config/world_config.go` — World configuration
- `internal/domain/config/generation_config.go` — Ore/hazard distributions

---

## World Configuration

All world parameters are defined in `WorldConfig`:

```go
type WorldConfig struct {
    Width          float32         // World width in pixels
    Height         float32         // World height in pixels
    GroundLevel    float32         // Y coordinate of ground level
    Seed           int64           // Procedural generation seed
    PlayerSpawn    PlayerSpawn     // Player starting position (X, Y)
    BuildingLayout BuildingLayout  // Building X positions (Y auto-calculated)
}
```

---

## World Dimensions

The game world extends far beyond the screen:

| Dimension | Size | Tiles |
|-----------|------|-------|
| Width | 3072 pixels | 48 tiles wide |
| Height | 51200 pixels | 800 tiles deep |
| Ground Level | 640 pixels | 10 tiles from top |
| Tile Size | 64×64 pixels | Standard |

**Surface Layout:**
```
[480px pad] [Hospital] [50px] [FuelStation] [230px] [Market] [130px] [UpgradeShop] [50px] [ItemShop] [532px pad]
```

**Storage:** Sparse tile storage — only non-empty tiles are stored in memory, enabling efficient large worlds.

---

## Chunk System

The world uses lazy chunk loading for performance:

```go
type World struct {
    tiles          map[[2]int]*entities.Tile  // Sparse tile storage
    loadedChunks   map[[2]int]bool            // Loaded chunk tracking
    generator      *ChunkGenerator
    groundLevelY   int                        // Ground level in tiles
}
```

### Chunk Loading

Chunks are 16×16 tiles. Loading happens proactively around the player:

```go
func (w *World) UpdateChunksAroundPlayer(playerX, playerY float32) {
    // Calculate player's chunk coordinates
    chunkX := int(playerX / TileSize / ChunkSize)
    chunkY := int(playerY / TileSize / ChunkSize)

    // Load 3×3 grid of chunks around player
    for dy := -1; dy <= 1; dy++ {
        for dx := -1; dx <= 1; dx++ {
            w.loadChunkIfNeeded(chunkX+dx, chunkY+dy)
        }
    }
}
```

### Sparse Storage

Tiles are stored by grid coordinates:

```go
func (w *World) GetTileAtGrid(gridX, gridY int) *entities.Tile {
    return w.tiles[[2]int{gridX, gridY}]
}

func (w *World) SetTileAtGrid(gridX, gridY int, tile *entities.Tile) {
    w.tiles[[2]int{gridX, gridY}] = tile
}
```

---

## Tile Types

### Basic Types

```go
const (
    TileTypeEmpty  TileType = iota  // Air, no collision
    TileTypeDirt                     // Solid, drillable, no value
    TileTypeOre                      // Solid, drillable, valuable
    TileTypeRock                     // Solid, NOT drillable (bomb only)
    TileTypeLava                     // Solid, drillable, deals damage
    TileTypeFloor                    // Solid, indestructible (boss room)
)
```

### Tile Properties

| Type | Solid | Drillable | Collectible | Notes |
|------|-------|-----------|-------------|-------|
| Empty | No | N/A | No | Air pockets, caves |
| Dirt | Yes | Yes | No | Filler terrain |
| Ore | Yes | Yes | Yes | Contains valuable resources |
| Rock | Yes | **No** | No | Only destroyed by bombs |
| Lava | Yes | Yes | No | Deals damage on drill completion |
| Floor | Yes | **No** | No | Indestructible boss room floor |

### Tile Methods

```go
func (t *Tile) IsSolid() bool {
    return t.Type != TileTypeEmpty
}

func (t *Tile) IsDrillable() bool {
    switch t.Type {
    case TileTypeDirt, TileTypeOre, TileTypeLava:
        return true
    default:
        return false
    }
}
```

---

## Procedural Generation

### Generation Config

```go
type GenerationConfig struct {
    Empty        TileDistribution  // Air pocket distribution
    Dirt         TileDistribution  // Dirt distribution
    DirtHardness float32           // Drilling time multiplier for dirt
    Ores         []OreConfig       // Dynamic list of ores
    Hazards      []HazardConfig    // Dynamic list of hazards
}
```

### Tile Distribution

Each tile type uses weighted random selection with depth-dependent weights:

```go
type TileDistribution struct {
    SurfaceWeight float32  // Weight at ground level (0% depth)
    DeepWeight    float32  // Weight at max depth (100% depth)
}
```

**Depth-Based Weight Formula:**
```go
weight := surfaceWeight + (deepWeight - surfaceWeight) * depthFactor
// depthFactor = (tileY - groundTileY) / maxDepth  (0.0 to 1.0)
```

### Default Distributions

| Type | Surface Weight | Deep Weight | Trend |
|------|---------------|-------------|-------|
| Empty | 8.0 | 0.5 | Decreases (fewer caves deep) |
| Dirt | 20.0 | 2.0 | Decreases (more hazards deep) |
| Hazards | 0.0 | 15.0+ | Increases (dominates deep) |

---

## Ore Configuration

Each ore type is defined with Gaussian distribution parameters:

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

### Gaussian Distribution

Ores use Gaussian (bell curve) distributions centered at specific depths:

```go
type TileDistribution struct {
    PeakDepth float32  // Depth (in tiles) where ore is most common
    Sigma     float32  // Spread of distribution (standard deviation)
    MaxWeight float32  // Maximum spawn weight at peak
}
```

**Weight Formula:**
```
weight = maxWeight × e^(-(depth - peakDepth)² / (2σ²))
```

### Level 1 Ores

| Ore | Peak Depth | Sigma | Max Weight | Value | Hardness |
|-----|------------|-------|------------|-------|----------|
| Copper | -75px | 120 | 8.0 | $25 | 1.2 |
| Iron | 70px | 90 | 5.0 | $75 | 1.5 |
| Gold | 230px | 80 | 3.0 | $300 | 1.8 |
| Mythril | 360px | 70 | 2.2 | $1500 | 2.1 |
| Platinum | 500px | 80 | 1.8 | $10000 | 2.5 |
| Diamond | 600px | 180 | 0.15 | $30000 | 3.0 |

**Design Notes:**
- Tight sigma values (70-120) create distinct depth bands
- Diamond's wider sigma (180) spreads it but keeps it extremely rare (0.15 weight)
- Copper peaks slightly above ground for easy early access

---

## Hazard Configuration

Hazard tiles create obstacles and damage zones at deep depths:

```go
type HazardConfig struct {
    ID            string           // Unique identifier (e.g., "rock", "lava")
    Name          string           // Display name
    Drillable     bool             // false = impenetrable (rock)
    FixedDuration float32          // If drillable: fixed drill time (0 = depth formula)
    OnDrillDamage float32          // Damage dealt when drilling completes
    Distribution  TileDistribution // Gaussian spawn parameters
    Color         [4]uint8         // RGBA for rendering
}
```

### Level 1 Hazards

| Hazard | Peak Depth | Sigma | Max Weight | Drillable | Damage |
|--------|------------|-------|------------|-----------|--------|
| Rock | 650 tiles (~80%) | 200 | 15.0 | No | N/A |
| Lava | 750 tiles (~85%) | 150 | 12.0 | Yes (0.3s) | 100 HP |

### Rock Mechanics

- **Impenetrable**: Cannot be drilled
- **Block Movement**: Prevents player and drilling
- **Bomb Destruction**: Only destroyed by bombs (both sizes)
- **First Appears**: ~40% depth
- **Dominates**: 80%+ depth

### Lava Mechanics

- **Fast Drill**: Fixed 0.3 second duration (depth-independent)
- **Damage on Completion**: 100 HP base damage
- **Heat Shield Reduction**: `damage = 100 - (resistance / 320.0 * 50)`
  - 0°C resistance: 100 damage
  - 160°C (Mk2): ~75 damage
  - 320°C (Mk5): 50 damage
- **First Appears**: ~60% depth
- **More Common**: 80%+ depth

---

## Generation Algorithm

### Chunk Generation

When a chunk is loaded, the generator creates tiles:

```go
func (g *ChunkGenerator) GenerateChunk(chunkX, chunkY int) map[[2]int]*entities.Tile {
    tiles := make(map[[2]int]*entities.Tile)

    for localY := 0; localY < ChunkSize; localY++ {
        for localX := 0; localX < ChunkSize; localX++ {
            gridX := chunkX*ChunkSize + localX
            gridY := chunkY*ChunkSize + localY

            // Skip above ground
            if gridY < g.groundLevelY {
                continue
            }

            tile := g.generateTile(gridX, gridY)
            if tile != nil {
                tiles[[2]int{gridX, gridY}] = tile
            }
        }
    }

    return tiles
}
```

### Tile Selection

Tiles are selected via weighted random:

```go
func (g *ChunkGenerator) generateTile(gridX, gridY int) *entities.Tile {
    depthFactor := g.calculateDepthFactor(gridY)

    // Calculate weights for all tile types
    emptyWeight := g.interpolateWeight(g.config.Empty, depthFactor)
    dirtWeight := g.interpolateWeight(g.config.Dirt, depthFactor)
    oreWeights := g.calculateOreWeights(gridY)
    hazardWeights := g.calculateHazardWeights(gridY)

    totalWeight := emptyWeight + dirtWeight + sum(oreWeights) + sum(hazardWeights)

    // Random selection
    roll := g.rand.Float32() * totalWeight

    if roll < emptyWeight {
        return nil  // Empty tile (not stored)
    }
    roll -= emptyWeight

    if roll < dirtWeight {
        return &entities.Tile{Type: entities.TileTypeDirt}
    }
    roll -= dirtWeight

    // Check each ore type...
    // Check each hazard type...
}
```

### Deterministic Generation

The generator uses seeded random for reproducibility:

```go
func NewChunkGeneratorFromConfig(config *GenerationConfig, seed int64) *ChunkGenerator {
    return &ChunkGenerator{
        config: config,
        rand:   rand.New(rand.NewSource(seed)),
    }
}
```

Same seed + same chunk coordinates = same tiles every time.

---

## Boss Room Generation

Boss rooms are special areas at world bottom:

```go
type BossRoomConfig struct {
    BossType    string      // Boss type identifier
    FloorType   FloorType   // Concrete or Lava floor
    RoomHeight  float32     // Height of boss room in pixels
    FloorHeight float32     // Height of floor in tiles
}
```

### Layout

```
[Mining Area - normal terrain]
         |
[Boss Room Start Y] ─────────────
         |
    [Empty Space - boss room, ~680px]
         |
[Floor Start Y] ─────────────────
    [Floor Tiles - indestructible]
[World Bottom] ──────────────────
```

### Floor Types

**Concrete Floor:**
- Solid, walkable surface
- Gray appearance
- Safe to stand on

**Lava Floor:**
- Solid but deals damage while standing
- Orange/red appearance
- Encourages movement

---

## World Methods

### Tile Access

```go
// Get tile at world coordinates
func (w *World) GetTileAt(worldX, worldY float32) *entities.Tile

// Get tile at grid coordinates
func (w *World) GetTileAtGrid(gridX, gridY int) *entities.Tile

// Get all tiles (for rendering)
func (w *World) GetAllTiles() map[[2]int]*entities.Tile
```

### Tile Removal

```go
// Standard drilling (respects drillability)
func (w *World) DrillTileAtGrid(gridX, gridY int) (*entities.Tile, bool)

// Bomb destruction (bypasses drillability)
func (w *World) NukeTileAtGrid(gridX, gridY int) (*entities.Tile, bool)
```

### Coordinate Conversion

```go
// World coordinates to grid coordinates
gridX := int(worldX / TileSize)
gridY := int(worldY / TileSize)

// Grid coordinates to world coordinates (tile center)
worldX := float32(gridX)*TileSize + TileSize/2
worldY := float32(gridY)*TileSize + TileSize/2
```

---

## Performance

### Sparse Storage Benefits

- Only ~30% of underground is solid tiles
- Empty tiles not stored (nil in map)
- Memory scales with actual content, not world size

### Chunk Loading Strategy

- 3×3 chunks loaded around player (9 chunks)
- Each chunk: 16×16 = 256 potential tiles
- Chunks cached once generated
- No unloading (memory acceptable for current world size)

### Benchmarks

```
Chunk generation: ~2.2ms per 16×16 chunk
Tile lookup: ~38ns per tile (sparse map)
Memory: ~100MB for full 48×800 tile world
```
