# Boss System Implementation Reference

This document provides a detailed reference of the current boss fight implementation. Use this to understand the existing system before making changes.

**Last Updated:** February 2026

---

## File Structure

```
internal/domain/
├── bosses/                          # Boss infrastructure
│   ├── boss.go                      # Boss interface + AOEInfo struct
│   ├── base_boss.go                 # BaseBoss struct with default implementations
│   ├── boxes.go                     # BoxSet, CollisionBox, Hitbox, Hurtbox, BoxDef types
│   ├── registry.go                  # Boss registration pattern
│   ├── phases/                      # Phase management package
│   │   ├── phase.go                 # phases.Config, phases.Manager
│   │   └── phase_test.go
│   ├── statemachine/
│   │   ├── types.go                 # State, StateID (int), StateContext, StateResult
│   │   ├── machine.go               # StateMachine implementation
│   │   └── machine_test.go
│   ├── attacks/
│   │   ├── attack.go                # Historical note (Attack interface removed)
│   │   ├── projectile_attack.go     # Cooldown-based projectile volleys
│   │   ├── projectile_attack_test.go
│   │   └── aoe_attack.go            # AOE attack with telegraph/damage/vulnerable phases
│   └── movement/
│       ├── movement.go              # MovementBehavior interface
│       ├── grounded.go              # Left-right patrol movement (implements MovementBehavior)
│       └── grounded_test.go
├── boss_catalog/                    # Boss implementations
│   └── test_boss/
│       ├── boss.go                  # TestBoss (embeds BaseBoss) + buildStates()
│       └── states.go                # State ID constants (iota)
├── systems/
│   ├── boss_fight.go                # BossFightSystem (room detection, contact damage)
│   └── projectile_system.go         # Projectile pool management
├── projectiles/
│   ├── spawn.go                     # SpawnRequest struct
│   └── movement.go                  # Movement interface + implementations
└── effects/
    ├── effect.go                    # Effect, EffectContext, DamageableEntity
    └── projectile.go                # ProjectileDamage effect
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

### AOEInfo Struct

**Location:** `internal/domain/bosses/boss.go`

Used by boss-specific renderers that type-assert to concrete boss types to render AOE effects.

```go
type AOEInfo struct {
    Position    types.Vec2
    Radius      float32
    IsTelegraph bool    // Warning phase
    IsDamaging  bool    // Damage phase
    StateTimer  float32
}
```

## BaseBoss Struct

**Location:** `internal/domain/bosses/base_boss.go`

Embeddable struct that provides default implementations for most Boss interface methods.

### Handler Interfaces

```go
// PhaseChangeHandler is called when the boss transitions to a new phase
type PhaseChangeHandler interface {
    OnPhaseChange(phaseIndex int, config phases.Config)
}

// DamageReactionHandler is called when the boss receives damage
type DamageReactionHandler interface {
    OnDamageReceived(hurtboxID string, damage float32)
}
```

### No-Op Default Handlers

Handlers default to no-ops so concrete bosses only override what they need:

```go
type noOpPhaseChangeHandler struct{}
func (noOpPhaseChangeHandler) OnPhaseChange(int, phases.Config) {}

type noOpDamageReactionHandler struct{}
func (noOpDamageReactionHandler) OnDamageReceived(string, float32) {}
```

### BaseBossConfig

```go
type BaseBossConfig struct {
    Position types.Vec2
    MaxHP    float32
    BoxSet   *BoxSet
    Phases   []phases.Config
}
```

### BaseBoss Struct

```go
type BaseBoss struct {
    Position      types.Vec2
    Damageable    components.Damageable
    Active        bool
    BoxSet        *BoxSet
    StateMachine  *statemachine.StateMachine
    PhaseManager  *phases.Manager
    CurrentPlayer *entities.Player

    // Optional handlers (nil = skip)
    PhaseChangeHandler    PhaseChangeHandler
    DamageReactionHandler DamageReactionHandler
}

func NewBaseBoss(cfg BaseBossConfig) *BaseBoss
func (b *BaseBoss) SetStateMachine(sm *statemachine.StateMachine)
```

### Default Implementations

```go
func (b *BaseBoss) Activate()
func (b *BaseBoss) Deactivate()
func (b *BaseBoss) IsActive() bool
func (b *BaseBoss) IsDefeated() bool
func (b *BaseBoss) GetHP() float32
func (b *BaseBoss) GetMaxHP() float32
func (b *BaseBoss) GetDamageable() *components.Damageable
func (b *BaseBoss) GetPosition() types.Vec2
func (b *BaseBoss) GetCollisionBoxes() []CollisionBox
func (b *BaseBoss) GetHitboxes() []Hitbox
func (b *BaseBoss) GetHurtboxes() []Hurtbox  // Override for custom vulnerability
func (b *BaseBoss) TakeDamageAt(hurtboxID string, baseDamage float32) float32
func (b *BaseBoss) BaseUpdate(player *entities.Player, dt float32) []projectiles.SpawnRequest
func (b *BaseBoss) GetCurrentPhase() int      // 1-indexed for display
func (b *BaseBoss) GetState() statemachine.StateID
func (b *BaseBoss) GetStateElapsed() float32
```

### BaseUpdate Flow

1. Check if active and not defeated
2. Store player reference
3. Check for phase transitions (calls `PhaseChangeHandler` if set)
4. Update state machine
5. Update box positions
6. Return spawn requests

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

Each box type has an `AABB()` method returning a `types.AABB` for intersection checks:

```go
func (c CollisionBox) AABB() types.AABB
func (h Hitbox) AABB() types.AABB
func (h Hurtbox) AABB() types.AABB
```

### Box Definition Types

Definitions describe boxes relative to boss position (static offsets). Used by `BoxSet` to create and update runtime boxes.

```go
type BoxDef struct {
    ID      string
    OffsetX float32 // Relative to boss position
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
```

### BoxSet (Pre-allocated Box Management)

```go
type BoxSet struct {
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

---

## Phase System

**Location:** `internal/domain/bosses/phases/phase.go`

### phases.Config

```go
type Config struct {
    HPThreshold        float32 // % HP where phase ends (0.66 = 66%)
    MovementSpeed      float32 // Speed in this phase
    ProjectileCooldown float32 // Time between projectile attacks
    AOECooldown        float32 // Time between AOE attacks (0 = disabled)
}
```

**Note:** Vulnerability is boss-specific logic, not part of Config. Each boss decides its own vulnerability rules based on phase index and state.

### phases.Manager

```go
func NewManager(maxHP float32, phases []Config) *Manager
func (pm *Manager) Update(currentHP float32) bool  // Returns true if phase changed
func (pm *Manager) GetCurrentConfig() Config
func (pm *Manager) GetCurrentPhase() int  // 0-indexed
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

---

## Movement System

**Location:** `internal/domain/bosses/movement/`

### MovementBehavior Interface (`movement.go`)

```go
type MovementBehavior interface {
    Update(currentPos types.Vec2, dt float32) types.Vec2
    GetVelocity() types.Vec2
    SetSpeed(speed float32)
    GetSpeed() float32
}
```

### Grounded (`grounded.go`)

Left-right patrol movement along the floor. Implements `MovementBehavior`.

```go
type GroundedConfig struct {
    Speed     float32 // Movement speed in pixels per second
    MinX      float32 // Left boundary
    MaxX      float32 // Right boundary
    FloorY    float32 // Y position of the floor
    BossWidth float32 // Width of the boss (for boundary calculation)
}

type Grounded struct { ... }

func NewGrounded(cfg GroundedConfig) *Grounded
func (g *Grounded) Update(currentPos types.Vec2, dt float32) types.Vec2
func (g *Grounded) GetVelocity() types.Vec2
func (g *Grounded) SetSpeed(speed float32)
func (g *Grounded) GetSpeed() float32
func (g *Grounded) GetDirection() float32  // 1 = right, -1 = left
```

---

## Attack Components

**Location:** `internal/domain/bosses/attacks/`

### ProjectileAttack (`projectile_attack.go`)

Cooldown-based projectile volleys aimed at the player with configurable spread.

```go
type ProjectileAttackConfig struct {
    Cooldown        float32
    ProjectileCount int
    ProjectileSpeed float32
    ProjectileSize  float32
    Damage          float32
}

type ProjectileAttack struct { ... }

func NewProjectileAttack(cfg ProjectileAttackConfig) *ProjectileAttack
func (a *ProjectileAttack) Update(bossAABB, playerAABB types.AABB, dt float32) []projectiles.SpawnRequest
func (a *ProjectileAttack) IsReady() bool
func (a *ProjectileAttack) GetCooldown() float32
func (a *ProjectileAttack) Reset()
```

### AOEAttack (`aoe_attack.go`)

Standalone AOE attack helper with its own state machine (telegraph → damage → vulnerable). Available for bosses that want a self-contained AOE component rather than handling AOE through the boss state machine.

**Note:** The TestBoss does NOT use this component — it manages AOE phases directly through its boss state machine states instead. This component exists as a reusable building block for future bosses.

```go
type AOEState int
const (
    AOEStateIdle AOEState = iota
    AOEStateTelegraph
    AOEStateDamage
    AOEStateVulnerable
)

type AOEAttackConfig struct {
    Cooldown           float32
    TelegraphDuration  float32
    DamageDuration     float32
    VulnerableDuration float32
    Radius             float32
    Damage             float32
}

type AOEAttack struct { ... }

func NewAOEAttack(cfg AOEAttackConfig) *AOEAttack
func (a *AOEAttack) Update(dt float32)
func (a *AOEAttack) StartAttack(bossAABB types.AABB)
func (a *AOEAttack) GetDamageToPlayer(playerAABB types.AABB) float32
func (a *AOEAttack) IsReady() bool
func (a *AOEAttack) GetCooldown() float32
func (a *AOEAttack) Reset()
func (a *AOEAttack) GetState() AOEState
func (a *AOEAttack) IsVulnerableWindow() bool
func (a *AOEAttack) IsTelegraphing() bool
func (a *AOEAttack) IsDamaging() bool
func (a *AOEAttack) GetPosition() types.Vec2
func (a *AOEAttack) GetRadius() float32
func (a *AOEAttack) GetStateTimer() float32
```

---

## Effects System

**Location:** `internal/domain/effects/`

### Core Types (`effect.go`)

```go
type DamageableEntity interface {
    GetHurtboxes() []bosses.Hurtbox
    GetDamageable() *components.Damageable
    TakeDamageAt(hurtboxID string, baseDamage float32) float32
}

type EffectContext struct {
    Player      *entities.Player
    World       *world.World
    Damageables []DamageableEntity // Boss + future enemies
}

type Effect interface {
    Apply(ctx *EffectContext)
}
```

### ProjectileDamage (`projectile.go`)

```go
type ProjectileDamage struct {
    Damage float32
}

func (e ProjectileDamage) Apply(ctx *EffectContext)  // Calls ctx.Player.DealDamage()
```

---

## TestBoss Implementation

**Location:** `internal/domain/boss_catalog/test_boss/`

### Configuration Constants (`boss.go`)

```go
const (
    MaxHP         = 100.0
    Width         = 100.0
    Height        = 100.0
    BaseSpeed     = 80.0
    ContactDamage = 20.0

    WindupDuration           = 1.0
    SlamDuration             = 0.3
    DoubleSlamPause          = 0.4
    Phase2VulnerableDuration = 3.0
    Phase3VulnerableDuration = 2.0
)

var phaseConfigs = []phases.Config{
    // Phase 1: 100% - 66% HP
    {HPThreshold: 0.66, MovementSpeed: BaseSpeed, ProjectileCooldown: 3.0, AOECooldown: 0},
    // Phase 2: 66% - 33% HP
    {HPThreshold: 0.33, MovementSpeed: BaseSpeed * 1.25, ProjectileCooldown: 2.0, AOECooldown: 6.0},
    // Phase 3: 33% - 0% HP
    {HPThreshold: 0.0, MovementSpeed: BaseSpeed * 1.5, ProjectileCooldown: 1.0, AOECooldown: 4.0},
}
```

### Struct (Embeds BaseBoss)

```go
type TestBoss struct {
    *bosses.BaseBoss

    // Boss-specific components
    movement         *movement.Grounded
    projectileAttack *attacks.ProjectileAttack
    worldWidth       float32
    floorY           float32

    // Boss-specific data
    aoeCooldown float32
    slamCount   int
    maxSlams    int
    aoeRadius   float32
    aoeDamage   float32
    aoePosition types.Vec2
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

### Handler Implementations

```go
// OnPhaseChange implements PhaseChangeHandler
func (b *TestBoss) OnPhaseChange(phaseIndex int, phaseCfg phases.Config) {
    b.movement.SetSpeed(phaseCfg.MovementSpeed)
    // Update projectile attack cooldown
    // Reset slam cooldown if in patrol state
}

// OnDamageReceived implements DamageReactionHandler
func (b *TestBoss) OnDamageReceived(hurtboxID string, damage float32) {
    if b.StateMachine.CurrentState() == StateVulnerable {
        b.StateMachine.TransitionTo(StatePatrol, &statemachine.StateContext{})
    }
}
```

### Vulnerability Implementation

```go
func (b *TestBoss) GetHurtboxes() []bosses.Hurtbox {
    // Phase 1: always has hurtboxes
    // Phase 2+: only has hurtboxes during StateVulnerable
    if b.PhaseManager.GetCurrentPhase() == 0 || b.StateMachine.CurrentState() == StateVulnerable {
        return b.BoxSet.Hurtboxes
    }
    return []bosses.Hurtbox{}
}

func (b *TestBoss) IsVulnerable() bool {
    return len(b.GetHurtboxes()) > 0
}
```

### Helper Methods

```go
// Internal logic
func (b *TestBoss) vulnerableDuration() float32          // Phase-dependent vulnerable window duration
func (b *TestBoss) updateProjectileAttack(dt float32) []projectiles.SpawnRequest
func (b *TestBoss) hasAOEAttack() bool                   // True if current phase has AOE
func (b *TestBoss) determineMaxSlams()                   // Sets maxSlams (1, or 50% chance of 2 in phase 3)
func (b *TestBoss) dealAOEDamage(dt float32)             // Applies damage to player if within AOE radius
```

### Rendering Helpers

These methods are used by the adapter layer (renderer) via type-assertion to the concrete `TestBoss`:

```go
func (b *TestBoss) IsVulnerable() bool              // True if GetHurtboxes() is non-empty
func (b *TestBoss) GetVulnerableTimer() float32      // -1 in phase 1 (always vulnerable), 0 when not, or remaining time
func (b *TestBoss) GetAOEInfo() *bosses.AOEInfo      // Returns AOE state for rendering, nil when no AOE active
func (b *TestBoss) GetState() statemachine.StateID   // Current state ID
func (b *TestBoss) GetStateTimer() float32           // Remaining time in current state (0 for StatePatrol)
```

### TakeDamageAt Override

TestBoss overrides `BaseBoss.TakeDamageAt` to check vulnerability via its own `GetHurtboxes()` first (which is phase/state-dependent), then delegates damage application and handler notification:

```go
func (b *TestBoss) TakeDamageAt(hurtboxID string, baseDamage float32) float32
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

### Constructor

```go
func NewBossFightSystem(boss bosses.Boss, bossRoomCfg config.BossRoomConfig, worldHeight float32) *BossFightSystem
```

Returns `nil` if `boss` is `nil`.

### Methods

```go
func (s *BossFightSystem) Update(player *entities.Player, dt float32) BossFightResult
func (s *BossFightSystem) IsPlayerInBossRoom(player *entities.Player) bool
func (s *BossFightSystem) DamageBoss(damage float32)           // Damages first available hurtbox
func (s *BossFightSystem) GetBoss() bosses.Boss
func (s *BossFightSystem) IsBossFightActive() bool             // Delegates to boss.IsActive()
```

### Update Flow

1. Check if player is in boss room
2. Activate/deactivate boss based on room entry/exit
3. Call `boss.Update(player, dt)`
4. Handle contact damage from hitboxes (via `hitbox.AABB()` intersection)
5. Handle lava floor damage (if `floorType == config.FloorLava`)
6. Determine game state (Playing/Victory/Defeat)

---

## Adding a New Boss

1. **Create package:** `internal/domain/boss_catalog/my_boss/`

2. **Define constants and phases** at top of `boss.go`:
   ```go
   const (
       MaxHP  = 150.0
       Width  = 120.0
       Height = 80.0
       // timing constants...
   )

   var phaseConfigs = []phases.Config{...}
   ```

3. **Create struct embedding BaseBoss:**
   ```go
   type MyBoss struct {
       *bosses.BaseBoss
       // Boss-specific fields
   }
   ```

4. **Register via init():**
   ```go
   func init() {
       bosses.Register("my_boss", func(roomStartY, worldWidth float32) bosses.Boss {
           return New(roomStartY, worldWidth)
       })
   }
   ```

5. **Create using NewBaseBoss:**
   ```go
   baseBoss := bosses.NewBaseBoss(bosses.BaseBossConfig{
       Position: types.NewVec2(centerX, floorY),
       MaxHP:    MaxHP,
       BoxSet:   bosses.NewBodyBoxSet(...),
       Phases:   phaseConfigs,
   })

   b := &MyBoss{BaseBoss: baseBoss, ...}
   b.PhaseChangeHandler = b
   b.DamageReactionHandler = b
   b.SetStateMachine(statemachine.NewStateMachine(b.buildStates(), StateIdle))
   ```

6. **Define buildStates() method** with direct field access:
   ```go
   func (b *MyBoss) buildStates() map[statemachine.StateID]*statemachine.State {
       return map[statemachine.StateID]*statemachine.State{
           StateIdle: {
               ID:      StateIdle,
               CanMove: true,
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

7. **Implement handlers:**
   ```go
   func (b *MyBoss) OnPhaseChange(phaseIndex int, cfg phases.Config) { ... }
   func (b *MyBoss) OnDamageReceived(hurtboxID string, damage float32) { ... }
   ```

8. **Override GetHurtboxes for custom vulnerability:**
   ```go
   func (b *MyBoss) GetHurtboxes() []bosses.Hurtbox {
       if /* vulnerability condition */ {
           return b.BoxSet.Hurtboxes
       }
       return []bosses.Hurtbox{}
   }
   ```

9. **Implement Update:**
   ```go
   func (b *MyBoss) Update(player *entities.Player, dt float32) []projectiles.SpawnRequest {
       return b.BaseUpdate(player, dt)
   }
   ```

10. **Import in game.go:**
    ```go
    import _ "github.com/Kishlin/drill-game/internal/domain/boss_catalog/my_boss"
    ```

**No modifications to core files required besides the import.**
