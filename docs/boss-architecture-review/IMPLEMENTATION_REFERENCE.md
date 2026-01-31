# Boss System Implementation Reference

This document provides a detailed reference of the current boss fight implementation. Use this to understand the existing system before making changes.

**Last Updated:** January 2026

---

## File Structure

```
internal/domain/
├── bosses/                      # Boss infrastructure
│   ├── boss.go                  # Boss interface definition
│   ├── base_boss.go             # BaseBoss struct with default implementations
│   ├── boxes.go                 # BoxSet, CollisionBox, Hitbox, Hurtbox types
│   ├── registry.go              # Boss registration pattern
│   ├── phases/                  # Phase management package
│   │   ├── phase.go             # phases.Config, phases.Manager
│   │   └── phase_test.go
│   ├── statemachine/
│   │   ├── types.go             # State, StateID (int), StateContext, StateResult
│   │   └── machine.go           # StateMachine implementation
│   ├── attacks/
│   │   └── projectile_attack.go # Cooldown-based projectile volleys
│   └── movement/
│       └── grounded.go          # Left-right patrol movement
├── boss_catalog/                # Boss implementations
│   └── test_boss/
│       ├── boss.go              # TestBoss (embeds BaseBoss)
│       └── states.go            # State definitions and StateBehaviors
├── systems/
│   ├── boss_fight.go            # BossFightSystem (room detection, contact damage)
│   └── projectile_system.go     # Projectile pool management
├── projectiles/
│   ├── spawn.go                 # SpawnRequest struct
│   └── movement.go              # Movement interface + implementations
└── effects/
    ├── effect.go                # Effect interface, EffectContext, DamageableEntity
    └── projectile.go            # ProjectileDamage effect
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

**Update Flow:**
1. Check if player is in boss room
2. Activate/deactivate boss based on room entry/exit
3. Call `boss.Update(player, dt)`
4. Handle contact damage from hitboxes
5. Handle lava floor damage
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
   b.SetStateMachine(statemachine.NewStateMachine(states, StateIdle))
   ```

6. **Implement handlers:**
   ```go
   func (b *MyBoss) OnPhaseChange(phaseIndex int, cfg phases.Config) { ... }
   func (b *MyBoss) OnDamageReceived(hurtboxID string, damage float32) { ... }
   ```

7. **Override GetHurtboxes for custom vulnerability:**
   ```go
   func (b *MyBoss) GetHurtboxes() []bosses.Hurtbox {
       if /* vulnerability condition */ {
           return b.BoxSet.Hurtboxes
       }
       return []bosses.Hurtbox{}
   }
   ```

8. **Implement Update:**
   ```go
   func (b *MyBoss) Update(player *entities.Player, dt float32) []projectiles.SpawnRequest {
       return b.BaseUpdate(player, dt)
   }
   ```

9. **Import in game.go:**
   ```go
   import _ "github.com/Kishlin/drill-game/internal/domain/boss_catalog/my_boss"
   ```

**No modifications to core files required besides the import.**
