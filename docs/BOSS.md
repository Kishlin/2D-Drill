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
├── boss.go              # Core interfaces (Boss, PhysicalBoss) + shared types
├── projectile.go        # Projectile entity with AABB collision
├── phase.go             # PhaseManager for HP-threshold phases
├── attacks/             # Reusable attack patterns
│   ├── attack.go        # Attack interface
│   ├── projectile_attack.go  # Fires projectiles at player
│   └── aoe_attack.go    # Ground slam with telegraph
├── movement/            # Reusable movement behaviors
│   ├── movement.go      # MovementBehavior interface
│   └── grounded.go      # Left-right patrol on floor
└── test_boss/           # TestBoss implementation
    └── boss.go          # State machine, phases, attacks

internal/adapters/rendering/bosses/
├── renderer.go          # BossRenderer interface + registry
└── test_boss.go         # TestBoss-specific rendering
```

---

## Core Interfaces

### Boss Interface

All bosses must implement the basic `Boss` interface:

```go
type Boss interface {
    Update(player *entities.Player, dt float32)
    GetHP() float32
    GetMaxHP() float32
    IsDefeated() bool
    IsActive() bool
    Activate()
    Deactivate()
    GetProjectiles() []*Projectile
}
```

### PhysicalBoss Interface

Bosses that can be damaged by bombs implement `PhysicalBoss`:

```go
type PhysicalBoss interface {
    Boss
    GetAABB() types.AABB
    GetDamageable() *components.Damageable  // HP component
    TakeDamage(damage float32)              // Respects vulnerability
    IsVulnerable() bool                     // When can boss take damage?
    GetVulnerableTimer() float32            // For UI feedback
    GetContactDamage() float32              // Damage per second on player contact
}
```

**Key Design:** `TakeDamage()` checks vulnerability internally—callers don't need to check `IsVulnerable()` first.

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

**Important:** The Damageable component is pure HP data. It does NOT handle vulnerability logic—that's boss-specific and controlled by the boss's state machine.

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

---

## State Machine

Bosses use state machines for animation and behavior. States drive visual feedback and vulnerability windows.

**Single Source of Truth:** The boss's state machine is the single source of truth for vulnerability. The Damageable component only stores HP—vulnerability logic lives in the boss.

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
    ├─ More slams to do? → StateWindup (0.4s pause)
    v
StateVulnerable (immobile, can be bombed)
    │
    ├─ Timer expires OR bomb hit
    v
StatePatrol (cooldown reset)
```

### State Implementation

```go
type BossState int

const (
    StatePatrol BossState = iota
    StateWindup
    StateSlam
    StateVulnerable
)

type TestBoss struct {
    aabb       types.AABB
    damageable components.Damageable  // HP component
    state      BossState
    stateTimer float32
    // ...
}

func (b *TestBoss) Update(player *entities.Player, dt float32) {
    switch b.state {
    case StatePatrol:
        b.updatePatrol(player, dt)
    case StateWindup:
        b.updateWindup(dt)
    case StateSlam:
        b.updateSlam(player, dt)
    case StateVulnerable:
        b.updateVulnerable(dt)
    }
}

// Vulnerability is controlled by state machine + phase config
func (b *TestBoss) IsVulnerable() bool {
    phaseCfg := b.phaseManager.GetCurrentConfig()
    return phaseCfg.AlwaysVulnerable || b.state == StateVulnerable
}

// TakeDamage checks vulnerability, then applies damage
func (b *TestBoss) TakeDamage(damage float32) {
    if !b.IsVulnerable() {
        return
    }
    b.damageable.TakeDamage(damage)

    // Taking damage ends vulnerability window
    if b.state == StateVulnerable {
        b.endVulnerability()
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
func (bfs *BossFightSystem) Update(player *entities.Player, dt float32) entities.GameState {
    playerY := player.AABB.Y + player.AABB.Height
    inRoom := playerY >= bfs.bossRoomStartY

    if inRoom && !bfs.playerInBossRoom {
        bfs.boss.Activate()
    } else if !inRoom && bfs.playerInBossRoom {
        bfs.boss.Deactivate()
        // Clear projectiles when leaving
    }
    bfs.playerInBossRoom = inRoom
    // ...
}
```

### Projectile Collision

```go
func (bfs *BossFightSystem) handleProjectileCollisions(player *entities.Player) {
    for _, projectile := range bfs.boss.GetProjectiles() {
        if projectile.Active && projectile.Intersects(player.AABB) {
            player.DealDamage(projectile.Damage)
            projectile.Active = false
        }
    }
}
```

### Contact Damage

```go
func (bfs *BossFightSystem) handleContactDamage(player *entities.Player, dt float32) {
    if physicalBoss, ok := bfs.boss.(bosses.PhysicalBoss); ok {
        contactDamage := physicalBoss.GetContactDamage()
        if contactDamage > 0 && player.AABB.Intersects(physicalBoss.GetAABB()) {
            player.DealDamage(contactDamage * dt)
        }
    }
}
```

### Bomb-Boss Interaction

Bombs damage entities through the `DamageableEntity` interface (defined in `effects/`):

```go
// effects/effect.go
type DamageableEntity interface {
    GetAABB() types.AABB
    GetDamageable() *components.Damageable
    TakeDamage(amount float32)  // Entity controls vulnerability check
}
```

Bomb effect iterates damageables from EffectContext:

```go
// effects/items.go
func (e Bomb) Apply(ctx *EffectContext) {
    // Calculate blast AABB...

    // Damage all damageable entities in range
    for _, entity := range ctx.Damageables {
        if entity.GetAABB().Intersects(blastAABB) {
            entity.TakeDamage(e.Damage)  // Boss checks vulnerability internally
        }
    }

    // Destroy tiles...
}
```

**Key Design:** The bomb calls `entity.TakeDamage()`, not `entity.GetDamageable().TakeDamage()`. This lets the boss control its own vulnerability check and side effects (like ending the vulnerability window).

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

### 1. Create Boss Implementation

`internal/domain/bosses/my_boss/boss.go`:

```go
package my_boss

type MyBoss struct {
    aabb       types.AABB
    damageable components.Damageable  // HP component
    active     bool
    state      MyBossState
    movement   *movement.Grounded
    attacks    []attacks.Attack
    phases     *bosses.PhaseManager
}

func New(roomStartY, worldWidth float32) *MyBoss {
    return &MyBoss{
        aabb:       types.AABB{X: centerX, Y: floorY, Width: 100, Height: 100},
        damageable: components.NewDamageable(200, 200),
        movement:   movement.NewGrounded(startPos, minX, maxX),
        attacks:    []attacks.Attack{attacks.NewProjectileAttack(...)},
        phases:     bosses.NewPhaseManager(myPhaseConfigs),
    }
}

// Boss interface - delegate to Damageable
func (b *MyBoss) GetHP() float32    { return b.damageable.HP }
func (b *MyBoss) GetMaxHP() float32 { return b.damageable.MaxHP }
func (b *MyBoss) IsDefeated() bool  { return b.damageable.IsDefeated() }

// PhysicalBoss interface
func (b *MyBoss) GetAABB() types.AABB                   { return b.aabb }
func (b *MyBoss) GetDamageable() *components.Damageable { return &b.damageable }

// TakeDamage - boss controls its own vulnerability
func (b *MyBoss) TakeDamage(damage float32) {
    if !b.IsVulnerable() {
        return
    }
    b.damageable.TakeDamage(damage)
    // Add boss-specific side effects here
}

// IsVulnerable - defined by state machine or phase config
func (b *MyBoss) IsVulnerable() bool {
    // Example: vulnerable during StateStunned or if AlwaysVulnerable phase
    return b.phases.GetCurrentConfig().AlwaysVulnerable || b.state == StateStunned
}
```

### 2. Create Boss Renderer

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

### 3. Register in Game

`internal/domain/engine/game.go`:

```go
func createBossByType(bossType string, roomStartY, worldWidth float32) (bosses.Boss, error) {
    switch bossType {
    case "test_boss":
        return test_boss.New(roomStartY, worldWidth), nil
    case "my_boss":
        return my_boss.New(roomStartY, worldWidth), nil
    default:
        return nil, fmt.Errorf("unknown boss type: %s", bossType)
    }
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
