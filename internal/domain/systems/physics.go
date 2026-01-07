package systems

import (
	"math"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/physics"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type PhysicsSystem struct {
	world *world.World
}

func NewPhysicsSystem(w *world.World) *PhysicsSystem {
	return &PhysicsSystem{world: w}
}

// UpdatePhysics applies physics using axis-separated AABB collision
func (ps *PhysicsSystem) UpdatePhysics(
	player *entities.Player,
	inputState input.InputState,
	dt float32,
) {
	// 1. Apply movement and gravity to velocity
	player.Velocity = physics.ApplyHorizontalMovement(
		player.Velocity, inputState, dt,
		player.Engine.MaxSpeed(), player.Engine.Acceleration(),
	)
	player.Velocity = physics.ApplyVerticalMovement(
		player.Velocity, inputState, dt,
		player.Engine.FlyAcceleration(), player.Engine.MaxUpwardSpeed(),
	)
	player.Velocity = physics.ApplyGravity(player.Velocity, dt)

	// 2. AXIS-SEPARATED COLLISION RESOLUTION

	// X-axis: integrate position → check → resolve
	player.AABB.X += player.Velocity.X * dt
	collisionsX := physics.CheckCollisions(player.AABB, ps.world)
	player.AABB, player.Velocity = physics.ResolveCollisionsX(player.AABB, player.Velocity, collisionsX)

	// Y-axis: integrate position → check → resolve
	player.AABB.Y += player.Velocity.Y * dt
	collisionsY := physics.CheckCollisions(player.AABB, ps.world)

	// Capture state before Y-resolution for fall damage calculation
	wasAirborne := !player.OnGround
	ySpeedBeforeLanding := player.Velocity.Y

	player.AABB, player.Velocity, player.OnGround = physics.ResolveCollisionsY(player.AABB, player.Velocity, collisionsY)

	// Apply fall damage on landing transition
	if wasAirborne && player.OnGround {
		ps.applyFallDamage(player, ySpeedBeforeLanding)
	}

	// Apply heat damage continuously based on depth
	ps.applyHeatDamage(player, dt)

	// 3. Enforce world boundary constraints (prevent player from leaving game area)
	ps.constrainPlayerToWorldBounds(player)
}

// constrainPlayerToWorldBounds clamps player position to world boundaries
func (ps *PhysicsSystem) constrainPlayerToWorldBounds(player *entities.Player) {
	// Horizontal bounds: player cannot go left of 0 or right of (worldWidth - playerWidth)
	minX := float32(0.0)
	maxX := ps.world.Width - float32(entities.PlayerWidth)

	if player.AABB.X < minX {
		player.AABB.X = minX
		player.Velocity.X = 0 // Stop horizontal velocity at boundary
	} else if player.AABB.X > maxX {
		player.AABB.X = maxX
		player.Velocity.X = 0 // Stop horizontal velocity at boundary
	}

	// Vertical bounds: player cannot go above y=0
	minY := float32(0.0)

	if player.AABB.Y < minY {
		player.AABB.Y = minY
		player.Velocity.Y = 0 // Stop vertical velocity at top boundary
	}

	// No maximum Y - player can dig infinitely deep
}

// applyFallDamage calculates and applies damage based on fall velocity.
// ySpeed is positive when falling downward (screen coordinates).
func (ps *PhysicsSystem) applyFallDamage(player *entities.Player, ySpeed float32) {
	// Only apply damage if falling fast enough
	if ySpeed < physics.FallDamageThreshold {
		return
	}

	// Calculate damage: (ySpeed - threshold) / divisor
	damage := (ySpeed - physics.FallDamageThreshold) / physics.FallDamageDivisor

	// Apply damage and clamp at zero (same pattern as fuel)
	player.HP -= damage
	if player.HP < 0 {
		player.HP = 0
	}
}

// applyHeatDamage calculates and applies damage based on depth-based temperature
func (ps *PhysicsSystem) applyHeatDamage(player *entities.Player, dt float32) {
	temperature := physics.CalculateTemperature(player.AABB.Y)

	excessHeat := temperature - player.HeatShield.HeatResistance()
	if excessHeat <= 0 {
		return // Player is within safe temperature range
	}

	// damage = baseDPS * (excessHeat / divisor)^exponent * dt
	damagePerSecond := float32(physics.HeatDamageBaseDPS) *
		float32(math.Pow(float64(excessHeat/float32(physics.HeatDamageDivisor)), float64(physics.HeatDamageExponent)))

	damage := damagePerSecond * dt

	// Apply damage and clamp at zero
	player.HP -= damage
	if player.HP < 0 {
		player.HP = 0
	}
}
