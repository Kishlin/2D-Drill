# Boss System Implementation Reference

This document provides a detailed reference of the current boss fight implementation. Use this to understand the existing system before making changes.

**Last Updated:** January 2026

---

## File Structure

```
internal/domain/
├── bosses/
│   ├── boss.go              # Boss interface definition
│   ├── boxes.go             # BoxSet, CollisionBox, Hitbox, Hurtbox types
│   ├── registry.go          # Boss registration pattern
│   ├── phase.go             # PhaseManager and PhaseConfig
│   ├── phase_test.go
│   ├── statemachine/
│   │   ├── types.go         # State, StateID (int), StateContext, StateResult
│   │   └── machine.go       # StateMachine implementation
│   ├── attacks/
│   │   └── projectile_attack.go  # Cooldown-based projectile volleys
│   ├── movement/
│   │   └── grounded.go      # Left-right patrol movement
│   └── test_boss/
│       ├── boss.go          # TestBoss struct, init() registration, config constants
│       └── states.go        # State definitions and StateBehaviors
├── systems/
│   ├── boss_fight.go        # BossFightSystem (room detection, contact damage)
│   └── projectile_system.go # Projectile pool management
├── projectiles/
│   ├── spawn.go             # SpawnRequest struct
│   └── movement.go          # Movement interface + implementations
└── effects/
    ├── effect.go            # Effect interface, EffectContext, DamageableEntity
    └── projectile.go        # ProjectileDamage effect
```

---

## Boss Registry

**Location:** `internal/domain/bosses/registry.go`

```go
// BossConstructor creates a boss instance given room parameters
type BossConstructor func(roomStartY, worldWidth float32) Boss

// registry holds all registered boss constructors
var registry = make(map[string]BossConstructor)

// Register adds a boss constructor to the registry.
// Call this in your boss package's init() function.
func Register(bossType string, constructor BossConstructor)

// Create instantiates a boss by type name using the registry.
func Create(bossType string, roomStartY, worldWidth float32) (Boss, error)
```

**Usage in boss packages:**
```go
func init() {
    bosses.Register("test_boss", func(roomStartY, worldWidth float32) bosses.Boss {
        return New(roomStartY, worldWidth)
    })
}
```

---

## Core Interfaces

### Boss Interface

**Location:** `internal/domain/bosses/boss.go`

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
    GetHitboxes() []Hitbox              // Damages player on contact
    GetHurtboxes() []Hurtbox            // Receives damage (empty = invulnerable)

    // Damage (only works if hurtbox exists)
    TakeDamageAt(hurtboxID string, baseDamage float32) float32
}
```

### DamageableEntity Interface

**Location:** `internal/domain/effects/effect.go`

```go
type DamageableEntity interface {
    GetHurtboxes() []bosses.Hurtbox
    GetDamageable() *components.Damageable
    TakeDamageAt(hurtboxID string, baseDamage float32) float32
}
```

Used by the effect system to damage any entity (boss, future enemies).

---

## Box System

**Location:** `internal/domain/bosses/boxes.go`

### Box Types

```go
// CollisionBox - Blocks player movement
type CollisionBox struct {
    ID     string
    X, Y, Width, Height float32
}

// Hitbox - Damages player on intersection
type Hitbox struct {
    ID           string
    X, Y, Width, Height float32
    DamagePerSec float32
}

// Hurtbox - Vulnerable zone where boss receives damage
type Hurtbox struct {
    ID               string
    X, Y, Width, Height float32
    DamageMultiplier float32  // 1.0 = normal, 2.0+ = weak point
}
```

### BoxSet (Pre-allocated Box Management)

```go
// BoxDef defines static properties relative to boss position
type BoxDef struct {
    ID      string
    OffsetX float32
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
    // Definitions (static)
    collisionDefs []BoxDef
    hitboxDefs    []HitboxDef
    hurtboxDefs   []HurtboxDef

    // Runtime boxes (pre-allocated, positions updated)
    CollisionBoxes []CollisionBox
    Hitboxes       []Hitbox
    Hurtboxes      []Hurtbox
}

func NewBoxSet(collisions []BoxDef, hitboxes []HitboxDef, hurtboxes []HurtboxDef) *BoxSet
func (bs *BoxSet) UpdatePositions(bossX, bossY float32)
```

### BodyBoxConfig Helper

For bosses with a single body box (collision = hitbox = hurtbox):

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
boxSet: bosses.NewBodyBoxSet(bosses.BodyBoxConfig{
    ID:               "body",
    Width:            100,
    Height:           100,
    DamagePerSec:     20,
    DamageMultiplier: 1.0,
}),
```

---

## Phase System

**Location:** `internal/domain/bosses/phase.go`

### PhaseConfig
```go
type PhaseConfig struct {
    HPThreshold        float32 // % HP where phase ends (0.66 = 66%)
    MovementSpeed      float32 // Speed in this phase
    ProjectileCooldown float32 // Time between projectile attacks
    AOECooldown        float32 // Time between AOE attacks (0 = disabled)
    AlwaysVulnerable   bool    // If true, always damageable
    VulnerableDuration float32 // Duration after AOE where vulnerable
}
```

### PhaseManager
```go
func NewPhaseManager(maxHP float32, phases []PhaseConfig) *PhaseManager
func (pm *PhaseManager) Update(currentHP float32) bool  // Returns true if phase changed
func (pm *PhaseManager) GetCurrentConfig() PhaseConfig
func (pm *PhaseManager) GetCurrentPhase() int
func (pm *PhaseManager) IsAlwaysVulnerable() bool
func (pm *PhaseManager) GetVulnerableDuration() float32
```

**Phase Progression:**
- Phases ordered from full HP to low HP
- Phase 0 at 100% HP, Phase 1 when HP drops below threshold, etc.
- `Update()` checks thresholds and advances phase

---

## State Machine

**Location:** `internal/domain/bosses/statemachine/`

### Types (`types.go`)

```go
// StateID - typed integer for compile-time safety
type StateID int
const StateIDNone StateID = -1  // Stay in current state

type StateContext struct {
    Player  *entities.Player
    Dt      float32  // Delta time this frame
    Elapsed float32  // Time since entering current state
}

type StateResult struct {
    NextState     StateID  // StateIDNone = stay in current state
    SpawnRequests []projectiles.SpawnRequest
}

type State struct {
    ID      StateID
    CanMove bool  // Whether movement behavior runs in this state

    OnEnter  func(ctx *StateContext)
    OnUpdate func(ctx *StateContext) StateResult
    OnExit   func(ctx *StateContext)
}
```

### StateMachine (`machine.go`)

```go
func NewStateMachine(states map[StateID]*State, initialState StateID) *StateMachine
func (sm *StateMachine) Update(ctx *StateContext) StateResult
func (sm *StateMachine) TransitionTo(newState StateID, ctx *StateContext)
func (sm *StateMachine) CurrentState() StateID
func (sm *StateMachine) CanMove() bool
func (sm *StateMachine) Elapsed() float32
```

**Update Flow:**
1. Set `ctx.Elapsed` to accumulated time
2. Call current state's `OnUpdate`
3. If `NextState` is set (not `StateIDNone`):
   - Call `OnExit` on current state
   - Reset elapsed to 0
   - Call `OnEnter` on new state
4. Otherwise increment elapsed by dt
5. Return `StateResult` (spawn requests, next state)

---

## TestBoss Implementation

**Location:** `internal/domain/bosses/test_boss/`

### Configuration Constants (`boss.go`)

```go
const (
    MaxHP         = 100.0
    Width         = 100.0
    Height        = 100.0
    BaseSpeed     = 80.0
    ContactDamage = 20.0

    WindupDuration  = 1.0
    SlamDuration    = 0.3
    DoubleSlamPause = 0.4
)

var phases = []bosses.PhaseConfig{
    // Phase 1: 100% - 66% HP - Always vulnerable
    {HPThreshold: 0.66, MovementSpeed: BaseSpeed, ProjectileCooldown: 3.0, AlwaysVulnerable: true},
    // Phase 2: 66% - 33% HP - AOE attacks
    {HPThreshold: 0.33, MovementSpeed: BaseSpeed * 1.25, ProjectileCooldown: 2.0, AOECooldown: 6.0, VulnerableDuration: 3.0},
    // Phase 3: 33% - 0% HP - Faster, double slams
    {HPThreshold: 0.0, MovementSpeed: BaseSpeed * 1.5, ProjectileCooldown: 1.0, AOECooldown: 4.0, VulnerableDuration: 2.0},
}
```

### State IDs (`states.go`)

```go
const (
    StatePatrol statemachine.StateID = iota
    StateWindup
    StateWindupBetween
    StateSlam
    StateVulnerable
)
```

### StateBehaviors Pattern

```go
type StateBehaviors struct {
    // Cooldown management
    GetAOECooldown    func() float32
    SetAOECooldown    func(float32)
    DecrementCooldown func(dt float32)

    // Slam management
    GetSlamCount      func() int
    IncrementSlam     func()
    ResetSlamCount    func()
    GetMaxSlams       func() int
    SetMaxSlams       func(int)
    DetermineMaxSlams func()

    // Movement and attacks
    UpdateMovement         func(dt float32)
    UpdateProjectileAttack func(dt float32) []projectiles.SpawnRequest

    // Phase info
    GetVulnerableDuration func() float32
    HasAOEAttack          func() bool

    // Damage and vulnerability
    DealAOEDamage    func(dt float32)
    EndVulnerability func()
}
```

### Five States

1. **StatePatrol** (`CanMove: true`)
   - Fires projectiles on cooldown
   - Decrements AOE cooldown
   - Transitions to StateWindup when AOE ready (Phase 2+)

2. **StateWindup** (`CanMove: false`)
   - Duration: `WindupDuration` (1.0s)
   - Determines slam count (1 or 2 in Phase 3)
   - Records AOE position at boss's feet
   - Transitions to StateSlam

3. **StateWindupBetween** (`CanMove: false`)
   - Duration: `DoubleSlamPause` (0.4s)
   - Pause between slams in double-slam
   - Transitions to StateSlam

4. **StateSlam** (`CanMove: false`)
   - Duration: `SlamDuration` (0.3s)
   - Damages player if in AOE radius
   - If more slams: → StateWindupBetween
   - Else: → StateVulnerable

5. **StateVulnerable** (`CanMove: false`)
   - Duration: Phase-dependent (2-3s)
   - Boss has hurtbox, can take damage
   - On damage: transitions to StatePatrol immediately
   - On timeout: transitions to StatePatrol, resets AOE cooldown

### Vulnerability Implementation

```go
func (b *TestBoss) GetHurtboxes() []bosses.Hurtbox {
    if b.phaseManager.IsAlwaysVulnerable() || b.stateMachine.CurrentState() == StateVulnerable {
        return b.boxSet.Hurtboxes
    }
    return nil // Invulnerable
}

func (b *TestBoss) IsVulnerable() bool {
    return len(b.GetHurtboxes()) > 0
}
```

---

## BossFightSystem

**Location:** `internal/domain/systems/boss_fight.go`

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

**Update Flow:**
1. Check if player is in boss room
2. Activate/deactivate boss based on room entry/exit
3. Call `boss.Update(player, dt)`
4. Handle contact damage from hitboxes
5. Handle lava floor damage
6. Determine game state (Playing/Victory/Defeat)

---

## Projectile System

**Location:** `internal/domain/systems/projectile_system.go`

```go
type Projectile struct {
    aabb     types.AABB
    movement projectiles.Movement
    damage   float32
    active   bool
}

type ProjectileSystem struct {
    pool         []Projectile  // Pre-allocated pool (64 default)
    bounds       ProjectileBounds
    renderBuffer []ProjectileRenderData
}

func NewProjectileSystem(bounds ProjectileBounds) *ProjectileSystem
func (ps *ProjectileSystem) SpawnAll(requests []SpawnRequest)
func (ps *ProjectileSystem) Update(dt float32, targets []CollisionTarget) []effects.Effect
func (ps *ProjectileSystem) GetActiveProjectiles() []ProjectileRenderData
func (ps *ProjectileSystem) Clear()
```

---

## Game Loop Integration

**Location:** `internal/domain/engine/game.go`

Boss fight is step 8 in the update sequence:

```go
func (g *Game) Update(dt, inputState) error {
    // 1. Chunk loading
    // 2. UI processing
    // 3. Heat damage
    // 4. Physics
    // 5. Fuel consumption
    // 6. Drilling
    // 7. Item usage

    // 8. Boss fight
    if g.bossFightSystem != nil {
        result := g.bossFightSystem.Update(g.player, dt)
        g.gameState = result.GameState
        g.projectileSystem.SpawnAll(result.SpawnRequests)
    }

    // 9. Projectile system
    projectileEffects := g.projectileSystem.Update(dt, []CollisionTarget{g.player})
    g.effectProcessor.Apply(g.effectContext, projectileEffects)

    return nil
}
```

**Boss Creation (via registry):**
```go
boss, err := bosses.Create(bossType, roomStartY, worldWidth)
```

---

## Adding a New Boss

1. **Create package:** `internal/domain/bosses/my_boss/`

2. **Define constants and phases** at top of `boss.go`:
   ```go
   const (
       MaxHP         = 150.0
       Width         = 120.0
       Height        = 80.0
       // timing constants...
   )

   var phases = []bosses.PhaseConfig{...}
   ```

3. **Register via init():**
   ```go
   func init() {
       bosses.Register("my_boss", func(roomStartY, worldWidth float32) bosses.Boss {
           return New(roomStartY, worldWidth)
       })
   }
   ```

4. **Use BoxSet for boxes:**
   ```go
   boxSet: bosses.NewBodyBoxSet(bosses.BodyBoxConfig{...}),
   // or for complex bosses:
   boxSet: bosses.NewBoxSet(collisionDefs, hitboxDefs, hurtboxDefs),
   ```

5. **Define typed state IDs** in `states.go`:
   ```go
   const (
       StateIdle statemachine.StateID = iota
       StateAttack
       StateVulnerable
   )
   ```

6. **Build states via BuildStates(behaviors):**
   ```go
   behaviors := b.buildStateBehaviors()
   states := BuildStates(behaviors)
   b.stateMachine = statemachine.NewStateMachine(states, StateIdle)
   ```

7. **Update box positions each frame:**
   ```go
   func (b *MyBoss) Update(...) {
       // ... state machine update ...
       b.boxSet.UpdatePositions(b.position.X, b.position.Y)
   }
   ```

**No modifications to `game.go` or other core files required.**
