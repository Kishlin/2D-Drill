# Physics System

This document covers physics implementation including AABB collision, movement, damage formulas, and camera. For high-level architecture, see [ARCHITECTURE.md](ARCHITECTURE.md).

---

## Overview

The physics system handles movement, collision detection, and damage calculation. All physics code is framework-independent (no Raylib imports) and lives in `internal/domain/physics/`.

**Key Files:**
- `physics/constants.go` — Physics tuning values
- `physics/movement.go` — Horizontal/vertical movement
- `physics/gravity.go` — Gravity application
- `physics/collision.go` — AABB collision detection and resolution
- `physics/fall_damage.go` — Fall damage calculation
- `physics/heat.go` — Heat damage calculation

---

## Domain Types

### Vec2

```go
type Vec2 struct {
    X float32
    Y float32
}

func (v Vec2) Add(other Vec2) Vec2
func (v Vec2) Sub(other Vec2) Vec2
func (v Vec2) Scale(scalar float32) Vec2
func (v Vec2) Magnitude() float32
func (v Vec2) Normalize() Vec2
```

### AABB

Axis-Aligned Bounding Box for collision:

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

**Design:** Value types (not pointers) for simplicity and cache locality. Small types (8-16 bytes) should be values per Go idiom.

---

## Physics Constants

Fixed constants in `physics/constants.go`:

```go
const (
    Gravity             = 800     // pixels/sec² - downward acceleration
    MoveDamping         = 1000    // pixels/sec² - how fast player slows down
    FlyDamping          = 300     // pixels/sec² - air resistance when flying
    FallDamageThreshold = 500.0   // pixels/sec - minimum speed for damage
    FallDamageDivisor   = 20.0    // damage scaling factor
)
```

Dynamic values come from player upgrades (Engine component):

| Stat | Base | Max (Mk5) |
|------|------|-----------|
| MaxMoveSpeed | 450 px/s | 600 px/s |
| MoveAcceleration | 2500 px/s² | 3500 px/s² |
| FlyAcceleration | 2500 px/s² | 3500 px/s² |
| MaxUpwardVelocity | -600 px/s | -775 px/s |

---

## Movement

### Horizontal Movement

```go
func ApplyHorizontalMovement(
    velocity Vec2,
    inputState InputState,
    dt float32,
    maxSpeed, acceleration float32,
) Vec2 {
    if inputState.Left {
        velocity.X -= acceleration * dt
    }
    if inputState.Right {
        velocity.X += acceleration * dt
    }

    // Apply damping when no input
    if !inputState.Left && !inputState.Right {
        if velocity.X > 0 {
            velocity.X -= MoveDamping * dt
            if velocity.X < 0 { velocity.X = 0 }
        } else if velocity.X < 0 {
            velocity.X += MoveDamping * dt
            if velocity.X > 0 { velocity.X = 0 }
        }
    }

    // Clamp to max speed
    if velocity.X > maxSpeed { velocity.X = maxSpeed }
    if velocity.X < -maxSpeed { velocity.X = -maxSpeed }

    return velocity
}
```

### Vertical Movement (Flying)

```go
func ApplyVerticalMovement(
    velocity Vec2,
    inputState InputState,
    dt float32,
    flyAcceleration, maxUpwardSpeed float32,
) Vec2 {
    if inputState.Up {
        velocity.Y -= flyAcceleration * dt  // Negative Y = upward
    }

    // Clamp upward velocity
    if velocity.Y < maxUpwardSpeed {
        velocity.Y = maxUpwardSpeed
    }

    return velocity
}
```

### Gravity

```go
func ApplyGravity(velocity Vec2, dt float32) Vec2 {
    velocity.Y += Gravity * dt
    return velocity
}
```

---

## AABB Collision System

The game uses Axis-Aligned Bounding Box collision with axis-separated resolution for precise 2D platformer physics.

### Why Axis-Separated?

**Without axis separation (naive):**
- Player moving diagonally into corner gets "stuck"
- Cannot slide along walls smoothly
- Ground detection is ambiguous

**With axis separation:**
- X collision resolved first, Y collision resolved second
- Player slides along walls naturally
- Clear ground/ceiling/wall detection

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

```go
func CheckCollisions(aabb AABB, world *World) []TileCollision {
    // Calculate which tiles the AABB might overlap
    minX, maxX, minY, maxY := GetOccupiedTileRange(aabb, TileSize)

    var collisions []TileCollision
    for x := minX; x <= maxX; x++ {
        for y := minY; y <= maxY; y++ {
            tile := world.GetTileAtGrid(x, y)
            if tile != nil && tile.IsSolid() {
                tileAABB := tile.GetAABB(x, y, TileSize)
                if aabb.Intersects(tileAABB) {
                    collisions = append(collisions, TileCollision{
                        TileX: x, TileY: y, TileAABB: tileAABB,
                    })
                }
            }
        }
    }
    return collisions
}
```

**Performance:** Player can overlap at most 4 tiles (2×2 grid), so maximum 4 intersection tests per frame.

### Collision Resolution

**X-Axis Resolution:**

```go
func ResolveCollisionsX(aabb AABB, velocity Vec2, collisions []TileCollision) (AABB, Vec2) {
    for _, col := range collisions {
        dx, _ := aabb.Penetration(col.TileAABB)
        aabb.X -= dx
        velocity.X = 0  // Stop horizontal movement
    }
    return aabb, velocity
}
```

**Y-Axis Resolution:**

```go
func ResolveCollisionsY(aabb AABB, velocity Vec2, collisions []TileCollision) (AABB, Vec2, bool) {
    onGround := false
    for _, col := range collisions {
        _, dy := aabb.Penetration(col.TileAABB)
        aabb.Y -= dy

        if dy > 0 {
            // Pushed up = landed on ground
            onGround = true
            velocity.Y = 0
        } else if dy < 0 {
            // Pushed down = hit ceiling
            velocity.Y = 0
        }
    }
    return aabb, velocity, onGround
}
```

### Penetration Calculation

```go
func (a AABB) Penetration(b AABB) (dx, dy float32) {
    // Calculate overlap on each axis
    overlapX := min(a.X+a.Width, b.X+b.Width) - max(a.X, b.X)
    overlapY := min(a.Y+a.Height, b.Y+b.Height) - max(a.Y, b.Y)

    // Determine push direction based on relative positions
    if a.X < b.X {
        dx = overlapX   // Push left
    } else {
        dx = -overlapX  // Push right
    }

    if a.Y < b.Y {
        dy = overlapY   // Push up
    } else {
        dy = -overlapY  // Push down
    }

    return dx, dy
}
```

**Key insight:** Signs are chosen so `position -= penetration` always pushes objects apart.

---

## World Boundary Constraints

Players cannot leave the game area:

```go
func (ps *PhysicsSystem) constrainPlayerToWorldBounds(player *entities.Player) {
    // Horizontal bounds
    minX := float32(0.0)
    maxX := ps.world.Width - float32(entities.PlayerWidth)

    if player.AABB.X < minX {
        player.AABB.X = minX
        player.Velocity.X = 0
    } else if player.AABB.X > maxX {
        player.AABB.X = maxX
        player.Velocity.X = 0
    }

    // Vertical upper bound only (can drill infinitely deep)
    if player.AABB.Y < 0 {
        player.AABB.Y = 0
        player.Velocity.Y = 0
    }
}
```

Called after collision resolution in the physics pipeline.

---

## Fall Damage

Damage calculated when landing (airborne → grounded transition):

```go
func ApplyFallDamage(player *entities.Player, ySpeed float32) {
    if ySpeed < FallDamageThreshold {
        return  // Below 500 px/sec - safe landing
    }

    damage := (ySpeed - FallDamageThreshold) / FallDamageDivisor
    player.DealDamage(damage)
}
```

**Formula:** `damage = (velocity - 500) / 20`

| Fall Speed | Damage |
|------------|--------|
| 500 px/s | 0 (safe) |
| 600 px/s | 5 HP |
| 700 px/s | 10 HP |
| 800 px/s | 15 HP |

### Landing Detection

```go
// In PhysicsSystem.UpdatePhysics()
wasAirborne := !player.OnGround
ySpeedBeforeLanding := player.Velocity.Y

// After Y collision resolution...
player.AABB, player.Velocity, player.OnGround = physics.ResolveCollisionsY(...)

// Apply fall damage on transition
if wasAirborne && player.OnGround {
    physics.ApplyFallDamage(player, ySpeedBeforeLanding)
}
```

---

## Heat Damage

Temperature increases with depth, causing exponential damage when exceeding heat resistance:

```go
func ApplyHeatDamage(player *entities.Player, dt float32) {
    temperature := CalculateTemperature(player.AABB.Y)

    excessHeat := temperature - player.HeatShield.HeatResistance()
    if excessHeat <= 0 {
        return  // Within safe range
    }

    // Exponential damage formula
    damagePerSecond := float32(HeatDamageBaseDPS) *
        float32(math.Pow(float64(excessHeat/float32(HeatDamageDivisor)),
                         float64(HeatDamageExponent)))
    damage := damagePerSecond * dt

    player.DealDamage(damage)
}
```

### Temperature Calculation

```go
func CalculateTemperature(playerY float32) float32 {
    if playerY < GroundLevelY {
        return BaseTemperature  // 15°C above ground
    }

    depth := playerY - GroundLevelY
    maxDepth := MaxWorldDepth - GroundLevelY

    // Linear interpolation: 15°C at ground → 350°C at max depth
    return BaseTemperature + (MaxTemperature - BaseTemperature) * (depth / maxDepth)
}
```

### Heat Constants

```go
const (
    BaseTemperature     = 15.0   // °C at ground level
    MaxTemperature      = 350.0  // °C at max depth
    GroundLevelY        = 640.0  // pixels
    HeatDamageBaseDPS   = 0.5    // Base damage per second
    HeatDamageDivisor   = 10.0   // Scaling factor
    HeatDamageExponent  = 1.5    // Exponential curve
)
```

### Damage Examples

| Excess Heat | Damage/sec |
|-------------|------------|
| 10°C | 0.16 HP/s |
| 30°C | 0.82 HP/s |
| 40°C | 1.26 HP/s |
| 50°C | 1.77 HP/s |

### Heat Shield Upgrades

| Tier | Resistance | Safe Depth |
|------|------------|-----------|
| Base | 50°C | 0-6,600px |
| Mk1 | 90°C | 6,600-14,000px |
| Mk2 | 140°C | 14,000-23,500px |
| Mk3 | 190°C | 23,500-33,000px |
| Mk4 | 250°C | 33,000-44,500px |
| Mk5 | 320°C | 44,500-64,000px |

---

## Camera System

The camera follows the player and is implemented in the rendering adapter (not domain).

### Camera2D

```go
type RaylibRenderer struct {
    camera       rl.Camera2D
    screenWidth  float32
    screenHeight float32
    worldWidth   float32
}
```

### Camera Following

```go
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

    targetX := clamp(playerCenterX, minX, maxX)
    targetY := clamp(playerCenterY, minY, maxY)
    r.camera.Target = rl.Vector2{X: targetX, Y: targetY}
}
```

### Viewport Culling

Only visible tiles are rendered for performance:

```go
func (r *RaylibRenderer) renderTiles(w *world.World) {
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

### Render Spaces

```go
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

    // Screen space rendering (UI, always visible)
    r.renderDebugInfo(game.GetPlayer(), inputState)

    rl.EndDrawing()
}
```

---

## Input State

Platform-agnostic input representation:

```go
type InputState struct {
    // Continuous inputs (checked every frame while held)
    Left  bool  // A or Arrow Left
    Right bool  // D or Arrow Right
    Up    bool  // W or Arrow Up
    Drill bool  // S or Arrow Down

    // Discrete inputs (true only on first frame pressed)
    Interact    bool  // E
    UseTeleport bool  // T
    UseRepair   bool  // R
    UseRefuel   bool  // F
    UseBomb     bool  // B
    UseBigBomb  bool  // G

    // Shop navigation
    CloseShop bool  // Q or Escape
    PrevTab   bool  // Z
    NextTab   bool  // X
    NavLeft   bool  // Arrow Left
    NavRight  bool  // Arrow Right
    NavUp     bool  // Arrow Up
    NavDown   bool  // Arrow Down
}

func (is InputState) HasMovementInput() bool {
    return is.Left || is.Right || is.Up || is.Drill
}
```

**Continuous vs Discrete:**
- Continuous: `IsKeyDown()` — true every frame while held
- Discrete: `IsKeyPressed()` — true only on first frame
