# Boss System

This document covers the boss fight system, including interfaces, state machines, phases, and rendering. For high-level architecture, see [ARCHITECTURE.md](ARCHITECTURE.md).

---

## Overview

The boss system provides extensible end-of-level encounters. Bosses are implemented in a separate catalog package with common interfaces, enabling different boss types with varying mechanics. Each boss type has its own AI logic (domain) and rendering (adapter).

**Key Principle:** Boss-specific behavior stays in boss-specific files. No generic `IsSlamming()` interfaces—renderers type-assert to concrete types.

---

## Package Structure

```
internal/domain/bosses/              # Boss infrastructure
├── boss.go                          # Core Boss interface + shared types (AOEInfo)
├── base_boss.go                     # BaseBoss struct with default implementations
├── boxes.go                         # Box types (CollisionBox, Hitbox, Hurtbox)
├── registry.go                      # Boss registration: Register() and Create()
├── phases/                          # Phase management package
│   └── phase.go                     # phases.Config, phases.Manager
├── attacks/                         # Reusable attack patterns
│   └── projectile_attack.go         # Fires projectiles at player
├── movement/                        # Reusable movement behaviors
│   └── grounded.go                  # Left-right patrol on floor
└── statemachine/                    # Generic state machine framework
    ├── types.go                     # StateID (int), StateContext, StateResult, State
    └── machine.go                   # StateMachine with transitions and lifecycle

internal/domain/boss_catalog/        # Boss implementations
└── test_boss/                       # TestBoss implementation
    ├── boss.go                      # Boss struct (embeds BaseBoss), states, init() registration
    └── states.go                    # State ID constants (iota)

internal/adapters/rendering/bosses/
├── renderer.go                      # BossRenderer interface + registry
└── test_boss.go                     # TestBoss-specific rendering
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

### BaseBoss Struct

`BaseBoss` provides default implementations for most interface methods. Embed it in concrete bosses:

```go
type BaseBoss struct {
    Position      types.Vec2
    Damageable    components.Damageable
    Active        bool
    BoxSet        *BoxSet
    StateMachine  *statemachine.StateMachine
    PhaseManager  *phases.Manager
    CurrentPlayer *entities.Player

    Self Boss  // Enables virtual dispatch (set to concrete boss: b.Self = b)

    PhaseChangeHandler    PhaseChangeHandler
    DamageReactionHandler DamageReactionHandler
}

// Handler interfaces
type PhaseChangeHandler interface {
    OnPhaseChange(phaseIndex int)
}

type DamageReactionHandler interface {
    OnDamageReceived(hurtboxID string, damage float32)
}
```

**Default implementations provided:**
- `Activate()`, `Deactivate()`, `IsActive()`, `IsDefeated()`
- `GetHP()`, `GetMaxHP()`, `GetDamageable()`
- `GetPosition()`, `GetCollisionBoxes()`, `GetHitboxes()`, `GetHurtboxes()`
- `TakeDamageAt()` (with DamageReactionHandler hook)
- `BaseUpdate()` (handles phase transitions, state machine, box positions)

**No-op handler defaults:** `NewBaseBoss` initializes both `PhaseChangeHandler` and `DamageReactionHandler` with no-op defaults. Concrete bosses only need to set `b.PhaseChangeHandler = b` or `b.DamageReactionHandler = b` if they want to react to those events.

**Virtual dispatch via `Self`:** The `Self Boss` field enables `BaseBoss.TakeDamageAt` to call the concrete boss's `GetHurtboxes()` override. Concrete bosses set `b.Self = b` during construction so that vulnerability checks dispatch correctly through the interface.

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

**Key Design:** Vulnerability is controlled by `GetHurtboxes()` returning empty slice when invulnerable. Each boss controls this based on its state machine and phase.

### BoxSet (Position Management)

`BoxSet` handles pre-allocation and position synchronization for all box types:

```go
type BoxSet struct {
    CollisionBoxes []CollisionBox
    Hitboxes       []Hitbox
    Hurtboxes      []Hurtbox
}

func NewBoxSet(collisions []BoxDef, hitboxes []HitboxDef, hurtboxes []HurtboxDef) *BoxSet
func (bs *BoxSet) UpdatePositions(bossX, bossY float32)
```

**Simple body box helper:**

```go
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

// In Update() - handled by BaseUpdate():
b.BoxSet.UpdatePositions(b.Position.X, b.Position.Y)
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

// boss_catalog/test_boss/states.go - each boss defines its own states
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
type State struct {
    ID StateID

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
// Boss creates state machine with inline state definitions
b.SetStateMachine(statemachine.NewStateMachine(b.buildStates(), StatePatrol))

// In Update(), BaseUpdate() handles state machine
func (b *TestBoss) Update(player *entities.Player, dt float32) []projectiles.SpawnRequest {
    return b.BaseUpdate(player, dt)
}

// Query current state for rendering/logic
if b.StateMachine.CurrentState() == StateVulnerable {
    // ...
}
```

### State Definitions Pattern

States are defined in a `buildStates()` method on the boss, with direct access to boss fields:

```go
func (b *TestBoss) buildStates() map[statemachine.StateID]*statemachine.State {
    return map[statemachine.StateID]*statemachine.State{
        StatePatrol: {
            ID: StatePatrol,
            OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
                b.Position = b.movement.Update(b.Position, ctx.Dt)
                // Direct field access - no callbacks needed
                if b.hasAOEAttack() {
                    b.aoeCooldown -= ctx.Dt
                    // ...
                }
                return statemachine.StateResult{NextState: statemachine.StateIDNone}
            },
        },
        // ... other states
    }
}
```

---

## Phase System

`phases.Manager` tracks HP-based phase transitions with configurable behaviors:

```go
// phases/phase.go
type Config struct {
    HPThreshold float32  // Phase ends when HP% drops below this
}

type Manager struct {
    phases       []Config
    currentPhase int
}

func (pm *Manager) Update(currentHP float32) bool  // Returns true if phase changed
func (pm *Manager) GetCurrentPhase() int           // 0-indexed
func (pm *Manager) GetCurrentConfig() Config
```

**Note:** `phases.Config` only contains the HP threshold. Boss-specific phase parameters (speeds, cooldowns, etc.) are stored in the concrete boss's own phase configuration struct. Vulnerability is also boss-specific logic — each boss decides its own vulnerability rules based on phase index and state.

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

func (pa *ProjectileAttack) Update(bossAABB, playerAABB types.AABB, dt float32) []SpawnRequest
```

---

## Movement System

Reusable movement behaviors in `bosses/movement/`:

### Grounded Movement

Left-right patrol on floor:

```go
type Grounded struct {
    position  types.Vec2
    direction float32  // 1 or -1
    minX      float32
    maxX      float32
}

func (g *Grounded) Update(position types.Vec2, dt float32) types.Vec2
```

---

## Boss Fight System Integration

The `BossFightSystem` orchestrates encounters:

```go
type BossFightSystem struct {
    boss            bosses.Boss
    bossRoomStartY  float32
    bossRoomEndY    float32
    floorStartY     float32
    floorEndY       float32
    floorType       config.FloorType
    bossRoomCfg     config.BossRoomConfig
    wasPlayerInRoom bool
}

type BossFightResult struct {
    GameState     entities.GameState
    SpawnRequests []projectiles.SpawnRequest
}
```

### Update Flow

```go
func (s *BossFightSystem) Update(player *entities.Player, dt float32) BossFightResult {
    // 1. Track player entry/exit
    // 2. Activate/deactivate boss
    // 3. If active: boss.Update(), handle contact damage, handle floor damage
    // 4. Return game state + spawn requests
}
```

### Contact Damage (Hitboxes)

Hitboxes deal damage per second on player contact:

```go
for _, hitbox := range boss.GetHitboxes() {
    if player.AABB.Intersects(hitbox.AABB()) {
        player.DealDamage(hitbox.DamagePerSec * dt)
    }
}
```

### Floor Damage

Lava floors deal configurable damage from `BossRoomConfig.FloorDamage`:

```go
if s.floorType == config.FloorLava {
    player.DealDamage(s.bossRoomCfg.FloorDamage * dt)
}
```

### Bomb-Boss Interaction (Hurtboxes)

Bombs damage bosses through hurtboxes. Empty hurtbox list = invulnerable:

```go
for _, hurtbox := range boss.GetHurtboxes() {
    if blastAABB.Intersects(hurtbox.AABB()) {
        boss.TakeDamageAt(hurtbox.ID, damage)
        break  // Only damage once per blast
    }
}
```

**Damage Values:**
- Regular bomb: 10 HP
- Big bomb: 25 HP

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

---

## Adding a New Boss

Adding a boss = create package + embed BaseBoss + call `bosses.Register()` in `init()`.

### 1. Create Boss Package

`internal/domain/boss_catalog/my_boss/boss.go`:

```go
package my_boss

import (
    "github.com/Kishlin/drill-game/internal/domain/bosses"
    "github.com/Kishlin/drill-game/internal/domain/bosses/phases"
)

// Register boss type on package import
func init() {
    bosses.Register("my_boss", func(roomStartY, worldWidth float32) bosses.Boss {
        return New(roomStartY, worldWidth)
    })
}

// HP thresholds for the phase manager (generic infrastructure)
var phaseThresholds = []phases.Config{
    {HPThreshold: 0.5},
    {HPThreshold: 0.0},
}

// Boss-specific parameters per phase
type myBossPhaseConfig struct {
    MovementSpeed      float32
    ProjectileCooldown float32
}

var myBossPhaseConfigs = []myBossPhaseConfig{
    {MovementSpeed: 100, ProjectileCooldown: 2.0},
    {MovementSpeed: 150, ProjectileCooldown: 1.0},
}

type MyBoss struct {
    *bosses.BaseBoss
    // Boss-specific fields
}

func New(roomStartY, worldWidth float32) *MyBoss {
    baseBoss := bosses.NewBaseBoss(bosses.BaseBossConfig{
        Position: types.NewVec2(centerX, floorY),
        MaxHP:    200,
        BoxSet: bosses.NewBodyBoxSet(bosses.BodyBoxConfig{
            ID:               "body",
            Width:            100,
            Height:           100,
            DamagePerSec:     20.0,
            DamageMultiplier: 1.0,
        }),
        Phases: phaseThresholds,
    })

    b := &MyBoss{BaseBoss: baseBoss}
    b.Self = b  // Required: enables virtual dispatch for GetHurtboxes/TakeDamageAt
    // Optional: override no-op defaults only if the boss needs to react to these events
    b.PhaseChangeHandler = b
    b.DamageReactionHandler = b
    b.SetStateMachine(statemachine.NewStateMachine(b.buildStates(), StateIdle))

    return b
}

// Implement handlers (only needed if handlers are set above)
func (b *MyBoss) OnPhaseChange(phaseIndex int) { ... }
func (b *MyBoss) OnDamageReceived(hurtboxID string, damage float32) { ... }

// Override for custom vulnerability
func (b *MyBoss) GetHurtboxes() []bosses.Hurtbox {
    if b.isVulnerable() {
        return b.BoxSet.Hurtboxes
    }
    return []bosses.Hurtbox{}
}

func (b *MyBoss) Update(player *entities.Player, dt float32) []projectiles.SpawnRequest {
    return b.BaseUpdate(player, dt)
}
```

`internal/domain/boss_catalog/my_boss/states.go`:

```go
package my_boss

import "github.com/Kishlin/drill-game/internal/domain/bosses/statemachine"

// State IDs (typed integers for compile-time safety)
const (
    StateIdle statemachine.StateID = iota
    StateAttack
    StateVulnerable
)
```

State definitions go in `boss.go` as a method with direct field access:

```go
func (b *MyBoss) buildStates() map[statemachine.StateID]*statemachine.State {
    return map[statemachine.StateID]*statemachine.State{
        StateIdle: {
            ID: StateIdle,
            OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
                // Direct field access
                b.someField = someValue
                return statemachine.StateResult{NextState: statemachine.StateIDNone}
            },
        },
        // ... other states
    }
}
```

### 2. Registration via Renderer Import

The boss registers automatically when its renderer package imports the concrete type. No changes to `engine/game.go` required — `bosses.Create()` finds the registered constructor at runtime.

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
    FloorDamage: 10.0,
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
