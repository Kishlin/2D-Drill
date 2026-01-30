# Boss System Implementation Reference

This document provides a detailed reference of the current boss fight implementation. Use this to understand the existing system before making changes.

---

## File Structure

```
internal/domain/
├── bosses/
│   ├── boss.go              # Boss interface definition
│   ├── boxes.go             # CollisionBox, Hitbox, Hurtbox types
│   ├── phase.go             # PhaseManager and PhaseConfig
│   ├── phase_test.go
│   ├── statemachine/
│   │   ├── types.go         # State, StateID, StateContext, StateResult
│   │   └── machine.go       # StateMachine implementation
│   ├── attacks/
│   │   ├── projectile_attack.go  # Cooldown-based projectile volleys
│   │   └── aoe_attack.go         # Telegraph → Damage → Vulnerable cycle
│   ├── movement/
│   │   └── grounded.go      # Left-right patrol movement
│   └── test_boss/
│       ├── boss.go          # TestBoss struct and methods
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

### CollisionBox
```go
type CollisionBox struct {
    X, Y, Width, Height float32
}

func (b CollisionBox) AABB() types.AABB
```
- Blocks player movement
- Used for physical presence
- No damage mechanics

### Hitbox
```go
type Hitbox struct {
    ID           string
    X, Y, Width, Height float32
    DamagePerSec float32
}

func (b Hitbox) AABB() types.AABB
```
- Damages player on intersection
- `DamagePerSec` applied each frame while intersecting
- BossFightSystem handles: `player.DealDamage(hitbox.DamagePerSec * dt)`

### Hurtbox
```go
type Hurtbox struct {
    ID               string
    X, Y, Width, Height float32
    DamageMultiplier float32  // 1.0 = normal, 2.0+ = weak point
}

func (b Hurtbox) AABB() types.AABB
```
- Vulnerable zones where boss receives damage
- Empty slice = boss is invulnerable
- Damage applied: `baseDamage * DamageMultiplier`

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
type PhaseManager struct {
    phases       []PhaseConfig
    currentPhase int
    maxHP        float32
}

func NewPhaseManager(phases []PhaseConfig, maxHP float32) *PhaseManager
func (pm *PhaseManager) Update(currentHP float32) bool  // Returns true if phase changed
func (pm *PhaseManager) CurrentPhase() PhaseConfig
func (pm *PhaseManager) CurrentPhaseIndex() int
```

**Phase Progression:**
- Phases ordered from full HP to low HP
- Phase 1 at 100% HP, Phase 2 when HP drops below threshold, etc.
- `Update()` checks thresholds and advances phase

---

## State Machine

**Location:** `internal/domain/bosses/statemachine/`

### Types (`types.go`)

```go
type StateID string

type StateContext struct {
    Player  *entities.Player
    Dt      float32  // Delta time this frame
    Elapsed float32  // Time since entering current state
}

type StateResult struct {
    NextState     StateID  // Empty string = stay in current state
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
type StateMachine struct {
    states       map[StateID]*State
    currentState StateID
    elapsed      float32
}

func NewStateMachine(states []*State, initialState StateID) *StateMachine
func (sm *StateMachine) Update(ctx *StateContext) StateResult
func (sm *StateMachine) Transition(newState StateID)
func (sm *StateMachine) CurrentState() StateID
func (sm *StateMachine) CanMove() bool
func (sm *StateMachine) Elapsed() float32
```

**Update Flow:**
1. Set `ctx.Elapsed` to accumulated time
2. Call current state's `OnUpdate`
3. If `NextState` is set:
   - Call `OnExit` on current state
   - Reset elapsed to 0
   - Call `OnEnter` on new state
4. Otherwise increment elapsed by dt
5. Return `StateResult` (spawn requests, next state)

---

## Attack Systems

**Location:** `internal/domain/bosses/attacks/`

### ProjectileAttack (`projectile_attack.go`)

```go
type ProjectileAttackConfig struct {
    Cooldown        float32  // Time between attacks
    ProjectileCount int      // Projectiles per volley
    ProjectileSpeed float32  // Pixels per second
    ProjectileSize  float32
    Damage          float32
}

type ProjectileAttack struct {
    config    ProjectileAttackConfig
    cooldown  float32
}

func NewProjectileAttack(cfg ProjectileAttackConfig) *ProjectileAttack
func (pa *ProjectileAttack) IsReady() bool
func (pa *ProjectileAttack) Update(bossAABB, playerAABB types.AABB, dt float32) []projectiles.SpawnRequest
func (pa *ProjectileAttack) Reset()
func (pa *ProjectileAttack) GetCooldown() float32
func (pa *ProjectileAttack) SetCooldown(cooldown float32)
```

**Firing Behavior:**
- Single projectile: Aimed at player center
- Multiple projectiles: Spread pattern (~17 degree spacing)
- Returns empty slice if on cooldown

### AOEAttack (`aoe_attack.go`)

```go
type AOEAttackConfig struct {
    Cooldown           float32
    TelegraphDuration  float32  // Warning before damage
    DamageDuration     float32  // Active damage time
    VulnerableDuration float32  // Vulnerability window after
    Radius             float32
    Damage             float32
}

type aoeState int
const (
    aoeStateIdle aoeState = iota
    aoeStateTelegraph
    aoeStateDamage
    aoeStateVulnerable
)

func NewAOEAttack(cfg AOEAttackConfig) *AOEAttack
func (a *AOEAttack) Update(dt float32)
func (a *AOEAttack) IsReady() bool
func (a *AOEAttack) StartAttack(bossAABB types.AABB)
func (a *AOEAttack) GetDamageToPlayer(playerAABB types.AABB) float32
func (a *AOEAttack) IsVulnerableWindow() bool
func (a *AOEAttack) GetState() aoeState
func (a *AOEAttack) GetPosition() types.Vec2
func (a *AOEAttack) GetRadius() float32
```

**AOE Cycle:**
1. **Idle**: Cooldown ticks down
2. **Telegraph**: Visual warning (no damage)
3. **Damage**: Applies damage if player in radius
4. **Vulnerable**: Boss vulnerable window
5. Back to **Idle**

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

### Grounded Movement (`grounded.go`)

```go
type Grounded struct {
    speed     float32
    minX      float32
    maxX      float32
    floorY    float32
    bossWidth float32
    direction float32  // 1.0 or -1.0
}

func NewGrounded(speed, minX, maxX, floorY, bossWidth float32) *Grounded
```

- Left-right patrol on floor
- Reverses at boundaries
- Speed adjustable per phase

---

## TestBoss Implementation

**Location:** `internal/domain/bosses/test_boss/`

### Boss Struct (`boss.go`)

```go
type TestBoss struct {
    position         types.Vec2
    damageable       components.Damageable
    active           bool
    movement         *movement.Grounded
    projectileAttack *attacks.ProjectileAttack
    phaseManager     *bosses.PhaseManager
    stateMachine     *statemachine.StateMachine

    // AOE attack data
    aoeCooldown float32
    slamCount   int
    maxSlams    int
    aoeRadius   float32
    aoeDamage   float32
    aoePosition types.Vec2
}
```

**Initialization (`New()`):**
- Position: Center-bottom of boss room
- Size: 100x100 pixels
- HP: 100
- 3 projectiles per volley, 200px/s speed
- Grounded movement with world boundaries
- 5-state state machine

### Phase Configuration

```go
phases := []bosses.PhaseConfig{
    {   // Phase 1: 100% - 66% HP
        HPThreshold:        0.66,
        MovementSpeed:      BaseSpeed,
        ProjectileCooldown: 3.0,
        AOECooldown:        0,      // No AOE
        AlwaysVulnerable:   true,
    },
    {   // Phase 2: 66% - 33% HP
        HPThreshold:        0.33,
        MovementSpeed:      BaseSpeed * 1.25,
        ProjectileCooldown: 2.0,
        AOECooldown:        6.0,
        AlwaysVulnerable:   false,
        VulnerableDuration: 3.0,
    },
    {   // Phase 3: 33% - 0% HP
        HPThreshold:        0,
        MovementSpeed:      BaseSpeed * 1.5,
        ProjectileCooldown: 1.0,
        AOECooldown:        4.0,
        AlwaysVulnerable:   false,
        VulnerableDuration: 2.0,
    },
}
```

### States (`states.go`)

**StateBehaviors Pattern:**
```go
type StateBehaviors struct {
    GetPhase             func() int
    GetAOECooldown       func() float32
    SetAOECooldown       func(float32)
    GetPhaseCooldown     func() float32
    GetVulnerableDuration func() float32
    FireProjectiles      func(*entities.Player) []projectiles.SpawnRequest
    SetSlamCount         func(int)
    GetSlamCount         func() int
    SetMaxSlams          func(int)
    GetMaxSlams          func() int
    SetAOEPosition       func(types.Vec2)
    GetAOEPosition       func() types.Vec2
    GetAOERadius         func() float32
    GetAOEDamage         func() float32
    TransitionTo         func(statemachine.StateID)
}
```

**Five States:**

1. **StatePatrol** (`CanMove: true`)
   - Fires projectiles on cooldown
   - Decrements AOE cooldown
   - Transitions to StateWindup when AOE ready (Phase 2+)

2. **StateWindup** (`CanMove: false`)
   - Duration: 1.0 second (hardcoded)
   - Sets slam count (1 normally, 2 in Phase 3)
   - Records AOE position at boss's feet
   - Transitions to StateSlam

3. **StateWindupBetween** (`CanMove: false`)
   - Duration: 0.4 seconds (hardcoded)
   - Pause between slams in double-slam
   - Transitions to StateSlam

4. **StateSlam** (`CanMove: false`)
   - Duration: 0.3 seconds (hardcoded)
   - Checks player distance from AOE position
   - Applies damage if in radius
   - Increments slam count
   - If more slams: → StateWindupBetween
   - Else: → StateVulnerable

5. **StateVulnerable** (`CanMove: false`)
   - Duration: Phase-dependent (3.0, 3.0, 2.0 seconds)
   - Boss has hurtbox, can take damage
   - On damage: transitions to StatePatrol immediately
   - On timeout: transitions to StatePatrol, resets AOE cooldown

### Box Implementation

**GetCollisionBoxes():**
```go
return []CollisionBox{{
    X: b.position.X,
    Y: b.position.Y,
    Width: 100, Height: 100,
}}
```

**GetHitboxes():**
```go
return []Hitbox{{
    ID: "body",
    X: b.position.X,
    Y: b.position.Y,
    Width: 100, Height: 100,
    DamagePerSec: 20,
}}
```

**GetHurtboxes():**
```go
// Only vulnerable in Phase 1 OR StateVulnerable
phase := b.phaseManager.CurrentPhase()
if phase.AlwaysVulnerable || b.stateMachine.CurrentState() == StateVulnerable {
    return []Hurtbox{{
        ID: "body",
        X: b.position.X,
        Y: b.position.Y,
        Width: 100, Height: 100,
        DamageMultiplier: 1.0,
    }}
}
return nil  // Invulnerable
```

---

## BossFightSystem

**Location:** `internal/domain/systems/boss_fight.go`

```go
type BossFightSystem struct {
    boss            bosses.Boss
    bossRoomStartY  float32  // Top of boss room
    bossRoomEndY    float32  // Bottom (start of floor)
    floorStartY     float32  // Top of floor tiles
    floorEndY       float32  // World bottom
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
```go
func (s *BossFightSystem) Update(player, dt) BossFightResult {
    playerInRoom := s.IsPlayerInBossRoom(player)

    // Activation tracking
    if playerInRoom && !s.wasPlayerInRoom {
        s.boss.Activate()
    }
    if !playerInRoom && s.wasPlayerInRoom {
        s.boss.Deactivate()
    }
    s.wasPlayerInRoom = playerInRoom

    // Boss update
    spawnRequests := s.boss.Update(player, dt)

    // Contact damage from hitboxes
    s.handleContactDamage(player, dt)

    // Lava floor damage
    s.handleFloorDamage(player)

    // Determine game state
    var gameState entities.GameState
    if s.boss.IsDefeated() {
        gameState = entities.GameStateVictory
    } else if player.HP <= 0 {
        gameState = entities.GameStateDefeat
    } else {
        gameState = entities.GameStatePlaying
    }

    return BossFightResult{gameState, spawnRequests}
}
```

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

**Update Flow:**
1. For each active projectile:
   - Move via `movement.Update()`
   - Cull if out of bounds
   - Check collision with targets
   - If hit: add `ProjectileDamage` effect, deactivate projectile
2. Return accumulated effects

---

## Projectile Movement Types

**Location:** `internal/domain/projectiles/movement.go`

```go
type Movement interface {
    Update(currentPos types.Vec2, dt float32) types.Vec2
}
```

**Implementations:**
- `Linear` - Constant velocity
- `Sinusoidal` - Wave pattern perpendicular to velocity
- `Homing` - Tracks target position (pointer)
- `Orbital` - Orbits around center point

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

**Boss Creation:**
```go
func createBossByType(bossType string, roomStartY, worldWidth float32) (Boss, error) {
    switch bossType {
    case "test_boss":
        return test_boss.New(roomStartY, worldWidth), nil
    default:
        return nil, fmt.Errorf("unknown boss type: %s", bossType)
    }
}
```

---

## Key Constants and Values

### TestBoss Defaults
| Property | Value |
|----------|-------|
| Size | 100x100 |
| HP | 100 |
| Base Speed | configurable |
| Projectile Count | 3 |
| Projectile Speed | 200 px/s |
| Contact Damage | 20/sec |
| AOE Radius | configurable |
| AOE Damage | configurable |

### State Durations (Hardcoded)
| State | Duration |
|-------|----------|
| Windup | 1.0s |
| Slam | 0.3s |
| WindupBetween | 0.4s |
| Vulnerable | Phase-dependent (2-3s) |

### Projectile Pool
- Default size: 64
- Bounds: World bounds ± 100px

---

## Adding a New Boss (Current Process)

1. Create new package: `internal/domain/bosses/new_boss/`

2. Create `boss.go`:
   - Define struct with required fields
   - Implement all 13+ Boss interface methods
   - Wire up PhaseManager, StateMachine, Movement, Attacks

3. Create `states.go`:
   - Define StateID constants
   - Create StateBehaviors callbacks
   - Define State structs with OnEnter/OnUpdate/OnExit

4. Modify `internal/domain/engine/game.go`:
   - Add import for new boss package
   - Add case to `createBossByType` switch

5. Update config to use new boss type
