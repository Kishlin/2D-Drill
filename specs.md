# World Generation with Ores - Specification

## Goals & Scope

**Phase Goal**: Transform the uniform dirt ground into a procedurally generated world rich with various ore types distributed by depth.

**In Scope:**
- Procedural chunk-based world generation with seeded RNG
- Multiple ore types with Gaussian depth distribution
- Visual distinction between ore types (colored tiles)
- Empty tiles (no collision, 20% of world)
- Dirt tiles (filler for remaining space)
- Solid dirt ground at surface level

**Out of Scope (Future Work):**
- Ore collection and inventory system
- Variable dig times per ore type
- Ore value/selling mechanics
- UI for seed selection
- Sprite-based ore visuals (using colored rectangles for now)

---

## Ore Types & Distribution Parameters

### Ore List

| Ore Type  | Peak Depth (tiles) | Sigma (spread) | Max Weight | Color (placeholder) |
|-----------|-------------------|----------------|------------|---------------------|
| Copper    | -50               | 150            | 10.0       | Orange (255, 140, 0) |
| Iron      | 0                 | 200            | 8.0        | Gray (128, 128, 128) |
| Silver    | 150               | 180            | 6.0        | Light Gray (192, 192, 192) |
| Gold      | 300               | 150            | 4.0        | Gold (255, 215, 0) |
| Mythril   | 500               | 200            | 3.0        | Cyan (0, 255, 255) |
| Platinum  | 700               | 180            | 2.0        | White (230, 230, 250) |
| Diamond   | 900               | 150            | 1.0        | Blue (0, 191, 255) |

**Notes:**
- Peak depth is in **tile coordinates** (not pixels). Ground level = tile Y of 10, so tile Y > 10 is underground.
- Negative peak depths mean the ore is already declining by ground level (e.g., copper very common at surface).
- Weights are relative - they'll be normalized when generating tiles.
- Colors are RGB values for rendering (adapter layer will use these).

### Distribution Formula

For a given ore at a given depth (tile Y coordinate):
```
weight(ore, tileY) = maxWeight × e^(-(tileY - peakDepth)² / (2 × sigma²))
```

Where:
- `tileY` is the tile's Y coordinate (10 = ground level, 11+ = underground)
- `peakDepth` is where this ore is most common
- `sigma` controls spread (larger = more gradual transitions)
- `maxWeight` is the ore's base weight at peak depth

**Implementation note:** Use `math.Exp()` in Go. If weight < 0.01, treat as 0 (ore doesn't spawn).

---

## World Generation Algorithm

### Chunk System

**Chunk Size:** 16×16 tiles (256 tiles per chunk)

**Chunk Coordinates:**
- Chunk X = `floor(tileX / 16)`
- Chunk Y = `floor(tileY / 16)`
- Example: Tile (35, 100) is in chunk (2, 6)

**Generation Trigger:**
- Generate chunk when player enters it or gets within 1 chunk distance
- Cache generated chunks in memory (never regenerate a chunk)
- Use sparse storage (only store non-empty tiles)

### Per-Tile Generation Logic

For each tile at `(tileX, tileY)` in a chunk:
```
1. Check if at ground level (tileY == 10):
   → Tile is Dirt (solid surface)
   → DONE, skip to next tile

2. Initialize seeded RNG with: hash(worldSeed, chunkX, chunkY, tileX, tileY)

3. Roll for Empty (20% chance):
   - If random(0, 1) < 0.20:
     → Tile is Empty (no collision, no ore)
     → DONE, skip to next tile

4. Calculate ore weights at this depth:
   - For each ore type:
     - weight[ore] = calculateGaussianWeight(ore, tileY)
     - If weight < 0.01: weight = 0 (ignore this ore)
   - totalOreWeight = sum of all ore weights

5. If totalOreWeight == 0:
   → Tile is Dirt (solid, diggable, no ore)
   → DONE

6. Roll for which ore spawns:
   - Generate random value: r = random(0, totalOreWeight)
   - Iterate through ores, subtracting their weights from r
   - When r <= 0, that ore spawns
   → Tile is Ore(oreType) (solid, diggable, contains ore)

7. If somehow nothing spawned (shouldn't happen):
   → Tile is Dirt (fallback)
```

**Algorithm Properties:**
- **Deterministic**: Same seed + coordinates = same tile every time
- **Independent**: Each tile generated independently (parallel-friendly)
- **Seeded**: World seed + chunk/tile coords → unique RNG state per tile
- **Solid Ground**: Ground level (Y=10) always generates as dirt

---

## Seeding System

### Seed Structure
```go
type WorldSeed struct {
    Value int64  // e.g., 12345
}
```

**Phase 1 Implementation:**
- Hard-coded seed in `main.go` (e.g., `seed := int64(42)`)
- Passed to `world.NewWorld(width, height, groundLevel, seed)`

**RNG Per Tile:**
```go
// Deterministic hash of seed + coordinates
rngSeed := hash(worldSeed, chunkX, chunkY, localX, localY)
rng := rand.New(rand.NewSource(rngSeed))
```

Use Go's `hash/fnv` or similar to create deterministic hashes:
```go
func hashCoordinates(worldSeed int64, chunkX, chunkY, tileX, tileY int) int64 {
    h := fnv.New64a()
    binary.Write(h, binary.LittleEndian, worldSeed)
    binary.Write(h, binary.LittleEndian, int64(chunkX))
    binary.Write(h, binary.LittleEndian, int64(chunkY))
    binary.Write(h, binary.LittleEndian, int64(tileX))
    binary.Write(h, binary.LittleEndian, int64(tileY))
    return int64(h.Sum64())
}
```

**Future Work:** UI for entering/displaying seed, random seed generation.

---

## Domain Layer Changes

### Tile Types (Enum)

`internal/domain/entities/tile.go`:
```go
type TileType int

const (
    TileEmpty TileType = iota  // No collision, player passes through
    TileDirt                    // Solid, diggable, no ore
    TileOre                     // Solid, diggable, contains ore
)
```

### Ore Types (Enum)

`internal/domain/entities/ore_type.go` (new file):
```go
type OreType int

const (
    OreCopper OreType = iota
    OreIron
    OreSilver
    OreGold
    OreMythril
    OrePlatinum
    OreDiamond
)

// Metadata for each ore
type OreMetadata struct {
    PeakDepth  float32  // Tile Y coordinate of peak
    Sigma      float32  // Standard deviation (spread)
    MaxWeight  float32  // Weight at peak depth
}

var OreDistributions = map[OreType]OreMetadata{
    OreCopper:   {PeakDepth: -50, Sigma: 150, MaxWeight: 10.0},
    OreIron:     {PeakDepth: 0, Sigma: 200, MaxWeight: 8.0},
    OreSilver:   {PeakDepth: 150, Sigma: 180, MaxWeight: 6.0},
    OreGold:     {PeakDepth: 300, Sigma: 150, MaxWeight: 4.0},
    OreMythril:  {PeakDepth: 500, Sigma: 200, MaxWeight: 3.0},
    OrePlatinum: {PeakDepth: 700, Sigma: 180, MaxWeight: 2.0},
    OreDiamond:  {PeakDepth: 900, Sigma: 150, MaxWeight: 1.0},
}
```

### Updated Tile Entity

`internal/domain/entities/tile.go`:
```go
type Tile struct {
    Type    TileType
    OreType OreType  // Only meaningful if Type == TileOre
}

func NewEmptyTile() *Tile {
    return &Tile{Type: TileEmpty}
}

func NewDirtTile() *Tile {
    return &Tile{Type: TileDirt}
}

func NewOreTile(oreType OreType) *Tile {
    return &Tile{Type: TileOre, OreType: oreType}
}

func (t *Tile) IsSolid() bool {
    return t.Type != TileEmpty
}

func (t *Tile) IsDiggable() bool {
    return t.Type != TileEmpty  // Both Dirt and Ore are diggable
}

// Getter for ore type (safe even if not an ore tile)
func (t *Tile) GetOreType() OreType {
    return t.OreType
}
```

### World Generation System

`internal/domain/world/generator.go` (new file):
```go
package world

type ChunkGenerator struct {
    seed        int64
    emptyRate   float32  // 0.20 for 20%
    groundLevel int      // Tile Y coordinate of ground (e.g., 10)
}

func NewChunkGenerator(seed int64, groundLevel float32) *ChunkGenerator {
    return &ChunkGenerator{
        seed:        seed,
        emptyRate:   0.20,
        groundLevel: int(groundLevel / TileSize),
    }
}

// Generate a single tile at given coordinates
func (cg *ChunkGenerator) GenerateTile(tileX, tileY int) *entities.Tile {
    // 1. Ground level is always dirt
    if tileY == cg.groundLevel {
        return entities.NewDirtTile()
    }
    
    // 2. Seed RNG deterministically
    rng := cg.seedRNG(tileX, tileY)
    
    // 3. Roll for empty
    if rng.Float32() < cg.emptyRate {
        return entities.NewEmptyTile()
    }
    
    // 4. Calculate ore weights
    weights := cg.calculateOreWeights(tileY)
    totalWeight := sumWeights(weights)
    
    // 5. If no ore weight, return dirt
    if totalWeight < 0.01 {
        return entities.NewDirtTile()
    }
    
    // 6. Roll for ore type
    oreType := cg.selectOreByWeight(rng, weights, totalWeight)
    if oreType == nil {
        return entities.NewDirtTile()  // Fallback
    }
    
    return entities.NewOreTile(*oreType)
}

// Calculate Gaussian weight for each ore at given depth
func (cg *ChunkGenerator) calculateOreWeights(tileY int) map[entities.OreType]float32 {
    weights := make(map[entities.OreType]float32)
    
    for oreType, meta := range entities.OreDistributions {
        weight := cg.gaussianWeight(float32(tileY), meta.PeakDepth, meta.Sigma, meta.MaxWeight)
        if weight >= 0.01 {
            weights[oreType] = weight
        }
    }
    
    return weights
}

// Gaussian distribution formula
func (cg *ChunkGenerator) gaussianWeight(depth, peak, sigma, maxWeight float32) float32 {
    exponent := -math.Pow(float64(depth-peak), 2) / (2 * math.Pow(float64(sigma), 2))
    return maxWeight * float32(math.Exp(exponent))
}

// Weighted random selection
func (cg *ChunkGenerator) selectOreByWeight(
    rng *rand.Rand, 
    weights map[entities.OreType]float32, 
    totalWeight float32,
) *entities.OreType {
    r := rng.Float32() * totalWeight
    
    for oreType, weight := range weights {
        r -= weight
        if r <= 0 {
            return &oreType
        }
    }
    
    return nil  // Shouldn't happen
}

// Deterministic RNG seeding
func (cg *ChunkGenerator) seedRNG(tileX, tileY int) *rand.Rand {
    seed := hashCoordinates(cg.seed, tileX/16, tileY/16, tileX%16, tileY%16)
    return rand.New(rand.NewSource(seed))
}

// Helper to sum weights
func sumWeights(weights map[entities.OreType]float32) float32 {
    total := float32(0)
    for _, w := range weights {
        total += w
    }
    return total
}
```

### World Chunk Management

`internal/domain/world/world.go` updates:
```go
type World struct {
    Width       float32
    Height      float32
    GroundLevel float32
    
    tiles        map[[2]int]*entities.Tile  // Sparse storage
    generator    *ChunkGenerator
    loadedChunks map[[2]int]bool            // Track generated chunks
}

func NewWorld(width, height, groundLevel float32, seed int64) *World {
    return &World{
        Width:        width,
        Height:       height,
        GroundLevel:  groundLevel,
        tiles:        make(map[[2]int]*entities.Tile),
        generator:    NewChunkGenerator(seed, groundLevel),
        loadedChunks: make(map[[2]int]bool),
    }
}

// Generate chunk if not already generated
func (w *World) EnsureChunkLoaded(chunkX, chunkY int) {
    key := [2]int{chunkX, chunkY}
    if w.loadedChunks[key] {
        return  // Already generated
    }
    
    // Generate all 16×16 tiles in chunk
    for localX := 0; localX < 16; localX++ {
        for localY := 0; localY < 16; localY++ {
            tileX := chunkX*16 + localX
            tileY := chunkY*16 + localY
            
            tile := w.generator.GenerateTile(tileX, tileY)
            
            // Only store non-empty tiles (sparse storage)
            if tile.Type != entities.TileEmpty {
                w.tiles[[2]int{tileX, tileY}] = tile
            }
        }
    }
    
    w.loadedChunks[key] = true
}

// Get tile at grid position (generate chunk if needed)
func (w *World) GetTileAtGrid(x, y int) *entities.Tile {
    chunkX := x / 16
    chunkY := y / 16
    w.EnsureChunkLoaded(chunkX, chunkY)
    
    return w.tiles[[2]int{x, y}]  // Returns nil if empty (sparse storage)
}

// Update to load chunks around player
func (w *World) UpdateChunksAroundPlayer(playerX, playerY float32) {
    playerChunkX := int(playerX/TileSize) / 16
    playerChunkY := int(playerY/TileSize) / 16
    
    // Load 3×3 grid of chunks around player (current + 1 in each direction)
    for dx := -1; dx <= 1; dx++ {
        for dy := -1; dy <= 1; dy++ {
            w.EnsureChunkLoaded(playerChunkX+dx, playerChunkY+dy)
        }
    }
}
```

### Game Loop Integration

`internal/domain/engine/game.go`:
```go
func (g *Game) Update(dt float32, inputState input.InputState) error {
    // Update chunks around player (before physics/digging)
    player := g.player
    g.world.UpdateChunksAroundPlayer(player.AABB.X, player.AABB.Y)
    
    // Existing systems (digging, physics, etc.)
    g.diggingSystem.Update(g.player, g.world, inputState, dt)
    g.physicsSystem.UpdatePhysics(g.player, inputState, dt)
    
    return nil
}
```

### Main Entry Point

`cmd/game/main.go`:
```go
const (
    screenWidth  = 1280
    screenHeight = 720
    targetFPS    = 60

    worldWidth  = screenWidth * 6
    worldHeight = 64 * 1000
    groundLevel = 640.0
    
    worldSeed = int64(42)  // Hard-coded for Phase 1
)

func main() {
    // ... existing setup ...
    
    gameWorld := world.NewWorld(worldWidth, worldHeight, groundLevel, worldSeed)
    game := engine.NewGame(gameWorld)
    
    // ... main loop ...
}
```

---

## Adapter Layer Changes

### Rendering Ore Colors

`internal/adapters/rendering/raylib.go`:
```go
// Map ore types to colors
var oreColors = map[entities.OreType]rl.Color{
    entities.OreCopper:   rl.NewColor(255, 140, 0, 255),    // Orange
    entities.OreIron:     rl.NewColor(128, 128, 128, 255),  // Gray
    entities.OreSilver:   rl.NewColor(192, 192, 192, 255),  // Light Gray
    entities.OreGold:     rl.NewColor(255, 215, 0, 255),    // Gold
    entities.OreMythril:  rl.NewColor(0, 255, 255, 255),    // Cyan
    entities.OrePlatinum: rl.NewColor(230, 230, 250, 255),  // White-ish
    entities.OreDiamond:  rl.NewColor(0, 191, 255, 255),    // Blue
}

func (r *RaylibRenderer) renderTiles(w *world.World) {
    tiles := w.GetAllTiles()
    
    // Viewport culling (existing logic)
    minVisibleX := int((r.camera.Target.X - r.screenWidth/2) / world.TileSize) - 1
    maxVisibleX := int((r.camera.Target.X + r.screenWidth/2) / world.TileSize) + 1
    minVisibleY := int((r.camera.Target.Y - r.screenHeight/2) / world.TileSize) - 1
    maxVisibleY := int((r.camera.Target.Y + r.screenHeight/2) / world.TileSize) + 1
    
    for coord, tile := range tiles {
        gridX, gridY := coord[0], coord[1]
        
        // Skip if outside viewport
        if gridX < minVisibleX || gridX > maxVisibleX ||
           gridY < minVisibleY || gridY > maxVisibleY {
            continue
        }
        
        // Determine color based on tile type
        var color rl.Color
        switch tile.Type {
        case entities.TileEmpty:
            continue  // Don't render empty tiles
        case entities.TileDirt:
            color = rl.Brown
        case entities.TileOre:
            color = oreColors[tile.GetOreType()]
        default:
            color = rl.Magenta  // Error color
        }
        
        // Draw tile rectangle
        x := float32(gridX) * world.TileSize
        y := float32(gridY) * world.TileSize
        rl.DrawRectangle(int32(x), int32(y), int32(world.TileSize), int32(world.TileSize), color)
    }
}
```

---

## Testing Strategy

### Unit Tests (Domain Layer)

**Gaussian Weight Calculation** (`internal/domain/world/generator_test.go`):
```go
func TestGaussianWeight(t *testing.T) {
    gen := NewChunkGenerator(42, 640)
    
    // Copper should have high weight at shallow depths
    copperShallow := gen.gaussianWeight(20, -50, 150, 10.0)
    assert.Greater(t, copperShallow, 5.0)
    
    // Diamond should have high weight at deep depths
    diamondDeep := gen.gaussianWeight(900, 900, 150, 1.0)
    assert.Greater(t, diamondDeep, 0.5)
    
    // Weights should decrease away from peak
    goldAtPeak := gen.gaussianWeight(300, 300, 150, 4.0)
    goldOffPeak := gen.gaussianWeight(500, 300, 150, 4.0)
    assert.Greater(t, goldAtPeak, goldOffPeak)
}
```

**Tile Generation Determinism** (`internal/domain/world/generator_test.go`):
```go
func TestGenerateTile_Deterministic(t *testing.T) {
    gen1 := NewChunkGenerator(12345, 640)
    gen2 := NewChunkGenerator(12345, 640)
    
    // Same seed + coords = same tile
    tile1 := gen1.GenerateTile(10, 50)
    tile2 := gen2.GenerateTile(10, 50)
    
    assert.Equal(t, tile1.Type, tile2.Type)
    if tile1.Type == entities.TileOre {
        assert.Equal(t, tile1.OreType, tile2.OreType)
    }
}
```

**Ground Level Always Dirt** (`internal/domain/world/generator_test.go`):
```go
func TestGenerateTile_GroundLevelIsDirt(t *testing.T) {
    gen := NewChunkGenerator(42, 640)
    groundTileY := 10  // 640 / 64 = 10
    
    // Test multiple X coordinates at ground level
    for x := 0; x < 100; x++ {
        tile := gen.GenerateTile(x, groundTileY)
        assert.Equal(t, entities.TileDirt, tile.Type, 
            "Ground level tile at X=%d should be dirt", x)
    }
}
```

**Empty Tile Rate** (`internal/domain/world/generator_test.go`):
```go
func TestGenerateTile_EmptyRate(t *testing.T) {
    gen := NewChunkGenerator(42, 640)
    
    emptyCount := 0
    totalTiles := 1000
    
    // Test underground tiles (not ground level)
    for i := 0; i < totalTiles; i++ {
        tile := gen.GenerateTile(i, 50)  // Below ground
        if tile.Type == entities.TileEmpty {
            emptyCount++
        }
    }
    
    // Should be ~20% (allow 15-25% margin)
    emptyRate := float32(emptyCount) / float32(totalTiles)
    assert.InDelta(t, 0.20, emptyRate, 0.05)
}
```

**Chunk Loading** (`internal/domain/world/world_test.go`):
```go
func TestEnsureChunkLoaded_OnlyOnce(t *testing.T) {
    world := NewWorld(7680, 64000, 640, 42)
    
    // First load
    world.EnsureChunkLoaded(0, 0)
    assert.True(t, world.loadedChunks[[2]int{0, 0}])
    
    // Second load should not regenerate
    tilesBefore := len(world.tiles)
    world.EnsureChunkLoaded(0, 0)
    tilesAfter := len(world.tiles)
    
    assert.Equal(t, tilesBefore, tilesAfter)
}
```

### Manual Testing

1. **Visual Verification**: Run game, dig down, observe:
    - Ground level is solid dirt (no gaps)
    - Empty tiles below ground (no collision, can fall through)
    - Dirt tiles (brown, solid)
    - Ore distribution changes with depth (more gold deeper, etc.)
    - Colors match ore types

2. **Seed Testing**: Change seed in `main.go`, verify world changes

3. **Chunk Loading**: Fly around world, verify no lag spikes when entering new areas

4. **Performance**: Profile chunk generation time (should be <1ms per chunk)

---

## Implementation Order

1. **Tile & Ore Types**: Add enums and metadata (`tile.go`, `ore_type.go`)
2. **Generator Logic**: Implement Gaussian weights and tile generation (`generator.go`)
3. **Chunk System**: Chunk loading and caching in World (`world.go`)
4. **Game Loop Integration**: Call `UpdateChunksAroundPlayer()` in `game.go`
5. **Main Entry Point**: Pass seed to world constructor (`main.go`)
6. **Rendering**: Ore colors in adapter (`raylib.go`)
7. **Testing**: Unit tests for generation logic
8. **Tuning**: Adjust Gaussian parameters based on visual results

---

## Future Enhancements (Out of Scope)

- Ore collection → inventory system
- Variable dig times per ore type
- Ore value and selling mechanics
- UI for seed input/display
- Biomes or horizontal variation (caves, veins)
- Sprite-based ore visuals (replace colored rectangles)
- Chunk unloading (memory management for very deep worlds)
- Save/load generated world state

---

## Success Criteria

✅ World generates deterministically with same seed  
✅ Ground level (Y=10) is always solid dirt  
✅ ~20% of underground tiles are empty (verified via test)  
✅ Ore distribution follows Gaussian curves (visual verification)  
✅ Multiple ore types appear at same depths (mixed)  
✅ Chunk loading is seamless (no lag when flying around)  
✅ Domain layer has zero Raylib dependencies  
✅ Rendering adapter displays ores with distinct colors