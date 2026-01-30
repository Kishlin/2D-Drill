# Boss System Refactor

Tracking document for the boss system architectural overhaul.

## Goals

- Declarative state machines instead of implicit switch statements
- Centralized projectile system decoupled from entity lifecycle
- Hitbox/hurtbox composition instead of Boss/PhysicalBoss hierarchy
- All damage through the effects system
- Zero-allocation game loops where possible

## Phases

### Phase 1: Central ProjectileSystem
- [ ] Create `systems/projectile_system.go` with pooled projectiles
- [ ] Move projectile update/collision logic out of bosses
- [ ] Fix allocation issues (in-place filtering)
- [ ] Update BossFightSystem to use central system
- [ ] Remove projectile management from TestBoss

### Phase 2: Declarative State Machine
- [ ] Design `State` and `StateMachine` types
- [ ] Vulnerability as state property
- [ ] Reimplement TestBoss using declarative states
- [ ] Remove scattered vulnerability checks

### Phase 3: Hitbox/Hurtbox Composition
- [ ] Replace Boss/PhysicalBoss interfaces
- [ ] States define active boxes
- [ ] Update collision detection

### Phase 4: Effects Integration
- [ ] Boss damage emits effects
- [ ] Projectile damage emits effects
- [ ] Contact damage emits effects

### Phase 5: Projectile Patterns
- [ ] Direction injection (decouple from "aim at player")
- [ ] Pattern generators (spiral, spread, wave)
- [ ] Built on central system from Phase 1

## Current Status

**Phase 1: In Progress**
