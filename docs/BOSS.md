# Boss System

This document covers the boss fight system, including interfaces, state machines, phases, and rendering. For high-level architecture, see [ARCHITECTURE.md](ARCHITECTURE.md).

---

## Overview

The boss system provides extensible end-of-level encounters. Bosses are implemented as a separate package with common interfaces, enabling different boss types with varying mechanics. Each boss type has its own AI logic (domain) and rendering (adapter).

**Key Principle:** Boss-specific behavior stays in boss-specific files. No generic `IsSlamming()` interfaces—renderers type-assert to concrete types.

---

## Package Structure

```
internal/domain/bosses/
├── boss.go              # Core Boss interface + shared types (AOEInfo)
├── boxes.go             # Box types (CollisionBox, Hitbox, Hurtbox)
├── registry.go          # Boss registration: Register() and Create()
├── phase.go             # PhaseManager for HP-threshold phases
├── attacks/             # Reusable attack patterns
│   ├── attack.go        # Attack interface
│   ├── projectile_attack.go  # Fires projectiles at player
│   └── aoe_attack.go    # Ground slam with telegraph
├── movement/            # Reusable movement behaviors
│   ├── movement.go      # MovementBehavior interface
│   └── grounded.go      # Left-right patrol on floor
├── statemachine/        # Generic state machine framework
│   ├── types.go         # StateID (int), StateContext, StateResult, State
│   └── machine.go       # StateMachine with transitions and lifecycle
└── test_boss/           # TestBoss implementation
    ├── boss.go          # Boss struct + init() registration
    └── states.go        # State IDs (iota) + state definitions

internal/adapters/rendering/bosses/
├── renderer.go          # BossRenderer interface + registry
└── test_boss.go         # TestBoss-specific rendering
```

---

## Core Interfaces

### Boss Interface

All bosses implement the `Boss` interface with three box types for collision/damage:

```go
type Boss interface {
    // Lifecycle
    Update(player *entities.Player, dt float32) []projectiles.SpawnRequest
    Activate()
    Deactivate()
    IsActive() bool
    IsDefeated() bool

    // Health
    GetHP() float32
    GetMaxHP() float32
    GetDamageable() *components.Damageable

    // Position (origin for box offsets)
    GetPosition() types.Vec2

    // Three box types (state-dependent)
    GetCollisionBoxes() []CollisionBox  // Blocks player movement
    GetHitboxes() []Hitbox              // Damages player (contact damage)
    GetHurtboxes() []Hurtbox            // Receives damage (empty = invulnerable)

    // Damage (only works if hurtbox exists)
    TakeDamageAt(hurtboxID string, baseDamage float32) float32
}
```

### Box Types

```go
// CollisionBox blocks player movement
type CollisionBox struct {
    ID            string
    X, Y          float32
    Width, Height float32
}

// Hitbox deals damage to player on contact
type Hitbox struct {
    ID            string
    X, Y          float32
    Width, Height float32
    DamagePerSec  float32
}

// Hurtbox receives damage (empty slice = invulnerable)
type Hurtbox struct {
    ID               string
    X, Y             float32
    Width, Height    float32
    DamageMultiplier float32  // 1.0 = normal, 2.0 = weak point
}
```

**Key Design:** Vulnerability is controlled by `GetHurtboxes()` returning empty slice when invulnerable. The boss controls this based on its state machine and phase config.

### BoxSet (Position Management)

`BoxSet` handles pre-allocation and position synchronization for all box types:

```go
// Definition types (static, relative to boss position)
type BoxDef struct {
    ID      string
    OffsetX float32  // Relative to boss position
    OffsetY float32
    Width   float32
    Height  float32
}

type HitboxDef struct {
    BoxDef
    DamagePerSec float32
}

type HurtboxDef struct {
    BoxDef
    DamageMultiplier float32
}

// BoxSet manages pre-allocated boxes with position synchronization
type BoxSet struct {
    CollisionBoxes []CollisionBox  // Runtime boxes (positions updated)
    Hitboxes       []Hitbox
    Hurtboxes      []Hurtbox
}

func NewBoxSet(collisions []BoxDef, hitboxes []HitboxDef, hurtboxes []HurtboxDef) *BoxSet
func (bs *BoxSet) UpdatePositions(bossX, bossY float32)
```

**Simple body box helper:**

```go
type BodyBoxConfig struct {
    ID               string
    Width, Height    float32
    OffsetX, OffsetY float32
    DamagePerSec     float32
    DamageMultiplier float32
}

func NewBodyBoxSet(cfg BodyBoxConfig) *BoxSet
```

**Usage:**

```go
// In constructor:
boxSet: bosses.NewBodyBoxSet(bosses.BodyBoxConfig{
    ID:               "body",
    Width:            100,
    Height:           100,
    DamagePerSec:     20.0,
    DamageMultiplier: 1.0,
}),

// In Update():
b.boxSet.UpdatePositions(b.position.X, b.position.Y)

// In interface methods:
func (b *MyBoss) GetCollisionBoxes() []bosses.CollisionBox { return b.boxSet.CollisionBoxes }
func (b *MyBoss) GetHitboxes() []bosses.Hitbox             { return b.boxSet.Hitboxes }
func (b *MyBoss) GetHurtboxes() []bosses.Hurtbox           { return b.boxSet.Hurtboxes }
```

**Multi-box example (boss with weak point):**

```go
boxSet: bosses.NewBoxSet(
    []bosses.BoxDef{{ID: "body", Width: 100, Height: 120}},
    []bosses.HitboxDef{{BoxDef: bosses.BoxDef{ID: "body", Width: 100, Height: 120}, DamagePerSec: 15}},
    []bosses.HurtboxDef{
        {BoxDef: bosses.BoxDef{ID: "body", OffsetX: 10, OffsetY: 40, Width: 80, Height: 80}, DamageMultiplier: 1.0},
        {BoxDef: bosses.BoxDef{ID: "head", OffsetX: 25, OffsetY: -30, Width: 50, Height: 40}, DamageMultiplier: 2.0},
    },
),
```

### Damageable Component

Bosses use the `components.Damageable` struct for HP storage:

```go
// components/damageable.go
type Damageable struct {
    HP    float32
    MaxHP float32
}

func (d *Damageable) TakeDamage(amount float32)  // Just decrements HP
func (d *Damageable) IsDefeated() bool           // HP <= 0
```

**Important:** The Damageable component is pure HP data. Vulnerability logic is boss-specific, controlled by the state machine.

### Shared Types

```go
// AOEInfo for rendering AOE effects
type AOEInfo struct {
    Position    types.Vec2
    Radius      float32
    IsTelegraph bool    // Warning phase
    IsDamaging  bool    // Damage phase
    StateTimer  float32
}
```

### Boss Registry

Bosses self-register via `init()`:

```go
// bosses/registry.go
func Register(bossType string, constructor BossConstructor)
func Create(bossType string, roomStartY, worldWidth float32) (Boss, error)

// Usage in boss package
func init() {
    bosses.Register("my_boss", func(roomStartY, worldWidth float32) bosses.Boss {
        return New(roomStartY, worldWidth)
    })
}
```

---

## State Machine

Bosses use a generic state machine framework (`bosses/statemachine/`) for animation and behavior. States drive visual feedback and vulnerability windows.

**Single Source of Truth:** The boss's state machine is the single source of truth for vulnerability. The Damageable component only stores HP—vulnerability logic lives in the boss.

### StateID Type

State IDs use typed integers for compile-time safety:

```go
// statemachine/types.go
type StateID int
const StateIDNone StateID = -1  // Stay in current state

// test_boss/states.go - each boss defines its own states
const (
    StatePatrol statemachine.StateID = iota
    StateWindup
    StateWindupBetween
    StateSlam
    StateVulnerable
)
```

### State Definition

States are declarative structs with lifecycle hooks:

```go
// statemachine/types.go
type State struct {
    ID      StateID
    CanMove bool  // Movement behavior active in this state

    OnEnter  func(ctx *StateContext)
    OnUpdate func(ctx *StateContext) StateResult
    OnExit   func(ctx *StateContext)
}

type StateResult struct {
    NextState     StateID  // StateIDNone = stay in current state
    SpawnRequests []projectiles.SpawnRequest
}
```

### TestBoss State Flow

```
StatePatrol (moving + shooting)
    │
    ├─ AOE cooldown expires
    v
StateWindup (stopped, vibrating warning, 1 second)
    │
    v
StateSlam (AOE damage zone, 0.3 seconds)
    │
    ├─ More slams to do? → StateWindupBetween (0.4s pause) → StateSlam
    v
StateVulnerable (immobile, can be bombed)
    │
    ├─ Timer expires OR bomb hit
    v
StatePatrol (cooldown reset)
```

### State Machine Usage

```go
// Boss creates state machine with state definitions
behaviors := b.buildStateBehaviors()
states := BuildStates(behaviors)
b.stateMachine = statemachine.NewStateMachine(states, StatePatrol)

// In Update(), run state machine
ctx := &statemachine.StateContext{Player: player, Dt: dt}
result := b.stateMachine.Update(ctx)
return result.SpawnRequests

// Query current state for rendering/logic
if b.stateMachine.CurrentState() == StateVulnerable {
    // ...
}
```

### StateBehaviors Pattern

States access boss data through a behaviors struct (avoids circular dependencies):

```go
type StateBehaviors struct {
    GetAOECooldown    func() float32
    DecrementCooldown func(dt float32)
    UpdateMovement    func(dt float32)
    // ... other callbacks
}

func BuildStates(behaviors *StateBehaviors) map[statemachine.StateID]*statemachine.State {
    return map[statemachine.StateID]*statemachine.State{
        StatePatrol: {
            ID:      StatePatrol,
            CanMove: true,
            OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
                behaviors.UpdateMovement(ctx.Dt)
                // ...
            },
        },
        // ... other states
    }
}
```

---

## Phase System

`PhaseManager` tracks HP-based phase transitions with configurable behaviors:

```go
type PhaseConfig struct {
    HPThreshold        float32  // Phase ends when HP% drops below this
    MovementSpeed      float32
    ProjectileCooldown float32
    AOECooldown        float32  // 0 = disabled
    AlwaysVulnerable   bool
    VulnerableDuration float32
}

type PhaseManager struct {
    phases       []PhaseConfig
    currentPhase int
}

func (pm *PhaseManager) Update(hpPercent float32) bool {
    if pm.currentPhase < len(pm.phases)-1 {
        threshold := pm.phases[pm.currentPhase].HPThreshold
        if hpPercent < threshold {
            pm.currentPhase++
            return true  // Phase changed
        }
    }
    return false
}
```

### TestBoss Phases

| Phase | HP Range | Movement | Projectiles | Slam | Vulnerability |
|-------|----------|----------|-------------|------|---------------|
| 1 | 100-66% | 80 px/s | Every 3s | None | Always |
| 2 | 66-33% | 100 px/s | Every 2s | Every 6s | 3s after slam |
| 3 | 33-0% | 120 px/s | Every 1s | Every 4s | 2s after slam |

**Phase 3 Special:** 50% chance of double slam (slam → 0.4s pause → slam → vulnerable)

---

## Attack System

Reusable attack patterns in `bosses/attacks/`:

### Attack Interface

```go
type Attack interface {
    Update(dt float32)
    CanFire() bool
    Fire(bossPos, playerPos types.Vec2) []*Projectile
    GetCooldown() float32
    SetCooldown(cooldown float32)
}
```

### ProjectileAttack

Fires projectiles at the player:

```go
type ProjectileAttack struct {
    cooldown     float32
    maxCooldown  float32
    speed        float32
    damage       float32
    count        int     // Projectiles per volley
}

func (pa *ProjectileAttack) Fire(bossPos, playerPos types.Vec2) []*Projectile {
    direction := playerPos.Sub(bossPos).Normalize()
    projectiles := make([]*Projectile, pa.count)
    for i := 0; i < pa.count; i++ {
        projectiles[i] = NewProjectile(bossPos, direction.Scale(pa.speed), pa.damage)
    }
    return projectiles
}
```

### AOEAttack

Ground slam with telegraph warning:

```go
type AOEAttack struct {
    cooldown        float32
    maxCooldown     float32
    radius          float32
    damage          float32
    telegraphTime   float32
    damageTime      float32
}

func (aoe *AOEAttack) GetAOEInfo() AOEInfo {
    return AOEInfo{
        Position:    aoe.position,
        Radius:      aoe.radius,
        IsTelegraph: aoe.phase == TelegraphPhase,
        IsDamaging:  aoe.phase == DamagePhase,
        StateTimer:  aoe.timer,
    }
}
```

---

## Movement System

Reusable movement behaviors in `bosses/movement/`:

### MovementBehavior Interface

```go
type MovementBehavior interface {
    Update(dt float32, speed float32) types.Vec2
    GetPosition() types.Vec2
    SetBounds(minX, maxX float32)
}
```

### Grounded Movement

Left-right patrol on floor:

```go
type Grounded struct {
    position  types.Vec2
    direction float32  // 1 or -1
    minX      float32
    maxX      float32
}

func (g *Grounded) Update(dt float32, speed float32) types.Vec2 {
    g.position.X += g.direction * speed * dt

    // Reverse at boundaries
    if g.position.X <= g.minX {
        g.position.X = g.minX
        g.direction = 1
    } else if g.position.X >= g.maxX {
        g.position.X = g.maxX
        g.direction = -1
    }

    return g.position
}
```

---

## Projectile Entity

Projectiles are simple entities with AABB collision:

```go
type Projectile struct {
    AABB     types.AABB
    Velocity types.Vec2
    Damage   float32
    Active   bool
}

func (p *Projectile) Update(dt float32) {
    p.AABB.X += p.Velocity.X * dt
    p.AABB.Y += p.Velocity.Y * dt
}

func (p *Projectile) Intersects(other types.AABB) bool {
    return p.AABB.Intersects(other)
}
```

---

## Boss Fight System Integration

The `BossFightSystem` orchestrates encounters:

```go
type BossFightSystem struct {
    boss              bosses.Boss
    bossRoomStartY    float32
    playerInBossRoom  bool
}
```

### Player Entry/Exit

```go
func (bfs *BossFightSystem) Update(player *entities.Player, dt float32) BossFightResult {
    playerY := player.AABB.Y + player.AABB.Height
    inRoom := playerY >= bfs.bossRoomStartY

    if inRoom && !bfs.playerInBossRoom {
        bfs.boss.Activate()
    } else if !inRoom && bfs.playerInBossRoom {
        bfs.boss.Deactivate()
    }
    bfs.playerInBossRoom = inRoom

    // Boss Update returns projectile spawn requests
    spawnRequests := bfs.boss.Update(player, dt)
    // ...
}
```

### Contact Damage (Hitboxes)

```go
func (bfs *BossFightSystem) handleHitboxDamage(player *entities.Player, dt float32) {
    for _, hitbox := range bfs.boss.GetHitboxes() {
        hitboxAABB := types.AABB{X: hitbox.X, Y: hitbox.Y, Width: hitbox.Width, Height: hitbox.Height}
        if player.AABB.Intersects(hitboxAABB) {
            player.DealDamage(hitbox.DamagePerSec * dt)
        }
    }
}
```

### Bomb-Boss Interaction (Hurtboxes)

Bombs damage bosses through hurtboxes. Empty hurtbox list = invulnerable:

```go
func (bfs *BossFightSystem) handleBombDamage(blastAABB types.AABB, damage float32) {
    for _, hurtbox := range bfs.boss.GetHurtboxes() {
        hurtboxAABB := types.AABB{X: hurtbox.X, Y: hurtbox.Y, Width: hurtbox.Width, Height: hurtbox.Height}
        if blastAABB.Intersects(hurtboxAABB) {
            bfs.boss.TakeDamageAt(hurtbox.ID, damage)
            break  // Only damage once per blast
        }
    }
}
```

**Key Design:** `GetHurtboxes()` returns empty slice when invulnerable. `TakeDamageAt()` applies damage multiplier and triggers boss-specific side effects (like ending vulnerability window).

**Damage Values:**
- Regular bomb: 10 HP
- Big bomb: 25 HP

---

## Game State Transitions

```
GameStatePlaying (initial)
    ├─ Normal gameplay
    ├─ Boss can be active or inactive
    └─ If boss.IsDefeated() → GameStateVictory
       If player.HP <= 0 → GameStateDefeat

GameStateVictory (terminal)
    └─ Boss defeated, victory screen displayed

GameStateDefeat (terminal)
    └─ Player defeated, defeat screen displayed
```

---

## Rendering

### BossRenderer Interface

Each boss type has its own renderer in `adapters/rendering/bosses/`:

```go
type Renderer interface {
    CanRender(boss bosses.Boss) bool
    Render(boss bosses.Boss)
}

// Registry dispatches to appropriate renderer
func RenderBoss(boss bosses.Boss) bool
func RenderGeneric(boss bosses.Boss)  // Fallback
```

### TestBossRenderer

Type-asserts to access boss-specific state:

```go
type TestBossRenderer struct{}

func (r *TestBossRenderer) CanRender(boss bosses.Boss) bool {
    _, ok := boss.(*test_boss.TestBoss)
    return ok
}

func (r *TestBossRenderer) Render(boss bosses.Boss) {
    tb := boss.(*test_boss.TestBoss)

    // Access state for rendering
    state := tb.GetState()
    stateTimer := tb.GetStateTimer()
    aoeInfo := tb.GetAOEInfo()

    // Render based on state
    switch state {
    case test_boss.StatePatrol:
        if tb.IsVulnerable() {
            r.renderVulnerable(tb)  // Pink flashing
        } else {
            r.renderInvulnerable(tb)  // Gray tint
        }
    case test_boss.StateWindup:
        r.renderWindup(tb)  // Orange flashing + vibration
    case test_boss.StateSlam:
        r.renderSlam(tb, aoeInfo)  // Bright red + AOE circle
    case test_boss.StateVulnerable:
        r.renderVulnerable(tb)  // Pink flashing
    }
}
```

### Visual Feedback

| State | Visual Effect |
|-------|---------------|
| Patrol (Vulnerable) | Pink flashing |
| Patrol (Invulnerable) | Gray tint |
| Windup | Orange flashing + horizontal vibration |
| Slam | Bright red + AOE damage circle |
| Vulnerable | Pink flashing |

**AOE Effects:**
- **Telegraph**: Pulsing yellow circle (warning)
- **Damage**: Solid orange-red circle

### Main Renderer Integration

```go
// In raylib.go
func (r *RaylibRenderer) renderBoss(game *engine.Game) {
    boss := game.GetBoss()
    if boss == nil || !boss.IsActive() {
        return
    }

    // Try boss-specific renderer first
    if !bossrenderers.RenderBoss(boss) {
        bossrenderers.RenderGeneric(boss)  // Fallback
    }
}

func (r *RaylibRenderer) renderBossHPBar(boss bosses.Boss) {
    // HP bar at screen top with health gradient
}

func (r *RaylibRenderer) renderProjectiles(boss bosses.Boss) {
    // Active projectiles in world space
}

func (r *RaylibRenderer) renderGameStateOverlay(state entities.GameState) {
    // Victory/defeat screens
}
```

---

## Adding a New Boss

Adding a boss = create package + call `bosses.Register()` in `init()`.

### 1. Create Boss Package

`internal/domain/bosses/my_boss/boss.go`:

```go
package my_boss

import "github.com/Kishlin/drill-game/internal/domain/bosses"

// Register boss type on package import
func init() {
    bosses.Register("my_boss", func(roomStartY, worldWidth float32) bosses.Boss {
        return New(roomStartY, worldWidth)
    })
}

type MyBoss struct {
    position     types.Vec2
    damageable   components.Damageable
    active       bool
    stateMachine *statemachine.StateMachine
    movement     *movement.Grounded
    phaseManager *bosses.PhaseManager
    boxSet       *bosses.BoxSet
}

func New(roomStartY, worldWidth float32) *MyBoss {
    b := &MyBoss{
        position:   types.NewVec2(centerX, floorY),
        damageable: components.NewDamageable(200, 200),
        movement:   movement.NewGrounded(...),
        boxSet: bosses.NewBodyBoxSet(bosses.BodyBoxConfig{
            ID:               "body",
            Width:            100,
            Height:           100,
            DamagePerSec:     20.0,
            DamageMultiplier: 1.0,
        }),
    }

    behaviors := b.buildStateBehaviors()
    states := BuildStates(behaviors)
    b.stateMachine = statemachine.NewStateMachine(states, StateIdle)

    return b
}

// Implement Boss interface...
func (b *MyBoss) GetCollisionBoxes() []bosses.CollisionBox { return b.boxSet.CollisionBoxes }
func (b *MyBoss) GetHitboxes() []bosses.Hitbox             { return b.boxSet.Hitboxes }
func (b *MyBoss) GetHurtboxes() []bosses.Hurtbox {
    if b.isVulnerable() {
        return b.boxSet.Hurtboxes
    }
    return nil  // Invulnerable
}

func (b *MyBoss) Update(player *entities.Player, dt float32) []projectiles.SpawnRequest {
    // ... state machine update ...
    b.boxSet.UpdatePositions(b.position.X, b.position.Y)
    return result.SpawnRequests
}
```

`internal/domain/bosses/my_boss/states.go`:

```go
package my_boss

import "github.com/Kishlin/drill-game/internal/domain/bosses/statemachine"

// State IDs (typed integers for compile-time safety)
const (
    StateIdle statemachine.StateID = iota
    StateAttack
    StateVulnerable
)

func BuildStates(behaviors *StateBehaviors) map[statemachine.StateID]*statemachine.State {
    return map[statemachine.StateID]*statemachine.State{
        StateIdle: {
            ID:      StateIdle,
            CanMove: true,
            OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
                // ...
                return statemachine.StateResult{NextState: statemachine.StateIDNone}
            },
        },
        // ... other states
    }
}
```

### 2. Import in Game Engine

`internal/domain/engine/game.go`:

```go
import (
    "github.com/Kishlin/drill-game/internal/domain/bosses"
    _ "github.com/Kishlin/drill-game/internal/domain/bosses/my_boss"  // Register my_boss
)

// Boss creation uses registry (no switch statement needed)
boss, err := bosses.Create(bossType, roomStartY, worldWidth)
```

### 3. Create Boss Renderer

`internal/adapters/rendering/bosses/my_boss.go`:

```go
package bossrenderers

type MyBossRenderer struct{}

func (r *MyBossRenderer) CanRender(boss bosses.Boss) bool {
    _, ok := boss.(*my_boss.MyBoss)
    return ok
}

func (r *MyBossRenderer) Render(boss bosses.Boss) {
    mb := boss.(*my_boss.MyBoss)
    // Boss-specific rendering logic
}

func init() {
    Register(&MyBossRenderer{})
}
```

### 4. Configure in Level

```go
BossRoom: config.BossRoomConfig{
    BossType:    "my_boss",
    FloorType:   config.FloorConcrete,
    RoomHeight:  680.0,
    FloorHeight: 6.0,
}
```

---

## TestBoss Reference

### Stats

- HP: 100
- Size: 100×100 pixels
- Base Movement: 80 px/s
- Contact Damage: 20 HP/sec

### Attacks

**Projectile Volley:**
- Count: 3 projectiles
- Speed: 200 px/s
- Damage: 5 HP each

**Ground Slam:**
- Radius: 150 pixels
- Damage: 15 HP
- Telegraph: 1 second warning
- Damage Zone: 0.3 seconds
