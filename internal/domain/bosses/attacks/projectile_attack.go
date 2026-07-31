package attacks

import (
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// MovementFactory creates a projectile movement from a velocity vector.
// When nil, Linear movement is used by default.
type MovementFactory func(velocity types.Vec2) projectiles.Movement

// ProjectileAttackConfig holds configuration for a projectile attack
type ProjectileAttackConfig struct {
	Cooldown        float32         // Time between attacks in seconds
	ProjectileCount int             // Number of projectiles per volley
	ProjectileSpeed float32         // Speed in pixels per second
	ProjectileSize  float32         // Size of projectiles (width and height)
	Damage          float32         // Damage per projectile
	MovementFactory MovementFactory // Optional: custom movement type (nil = Linear)
}

// ProjectileAttack fires projectiles at the player
type ProjectileAttack struct {
	config        ProjectileAttackConfig
	cooldownTimer float32
	spreadAngle   float32 // Spread angle in radians for multiple projectiles
}

// NewProjectileAttack creates a new projectile attack
func NewProjectileAttack(cfg ProjectileAttackConfig) *ProjectileAttack {
	return &ProjectileAttack{
		config:        cfg,
		cooldownTimer: 0,   // Ready immediately
		spreadAngle:   0.3, // ~17 degrees spread
	}
}

// Update updates the attack and returns spawn requests for projectiles
func (a *ProjectileAttack) Update(bossAABB, playerAABB types.AABB, dt float32) []projectiles.SpawnRequest {
	// Update cooldown
	if a.cooldownTimer > 0 {
		a.cooldownTimer -= dt
	}

	// Check if ready to fire
	if a.IsReady() == false {
		return nil
	}

	// Fire projectiles
	requests := a.fire(bossAABB, playerAABB)
	a.cooldownTimer = a.config.Cooldown

	return requests
}

// fire creates spawn requests for projectiles aimed at the player
func (a *ProjectileAttack) fire(bossAABB, playerAABB types.AABB) []projectiles.SpawnRequest {
	// Calculate boss center
	bossX := bossAABB.X + bossAABB.Width/2
	bossY := bossAABB.Y + bossAABB.Height/2

	// Calculate player center
	playerX := playerAABB.X + playerAABB.Width/2
	playerY := playerAABB.Y + playerAABB.Height/2

	// Calculate direction to player
	dx := playerX - bossX
	dy := playerY - bossY
	direction := types.NewVec2(dx, dy).Normalize()

	requests := make([]projectiles.SpawnRequest, 0, a.config.ProjectileCount)

	if a.config.ProjectileCount == 1 {
		// Single projectile aimed directly at player
		velocity := direction.Scale(a.config.ProjectileSpeed)
		req := projectiles.SpawnRequest{
			Position: types.NewVec2(bossX, bossY),
			Size:     a.config.ProjectileSize,
			Damage:   a.config.Damage,
			Movement: a.createMovement(velocity),
		}
		requests = append(requests, req)
	} else {
		// Multiple projectiles with spread
		halfCount := float32(a.config.ProjectileCount-1) / 2
		angleStep := a.spreadAngle / float32(a.config.ProjectileCount-1)

		for i := 0; i < a.config.ProjectileCount; i++ {
			// Calculate angle offset from center
			angleOffset := (float32(i) - halfCount) * angleStep

			// Rotate direction by angle offset
			rotatedDir := rotateVector(direction, angleOffset)
			velocity := rotatedDir.Scale(a.config.ProjectileSpeed)

			req := projectiles.SpawnRequest{
				Position: types.NewVec2(bossX, bossY),
				Size:     a.config.ProjectileSize,
				Damage:   a.config.Damage,
				Movement: a.createMovement(velocity),
			}
			requests = append(requests, req)
		}
	}

	return requests
}

// createMovement creates the appropriate movement for a projectile
func (a *ProjectileAttack) createMovement(velocity types.Vec2) projectiles.Movement {
	if a.config.MovementFactory != nil {
		return a.config.MovementFactory(velocity)
	}
	return projectiles.Linear{Velocity: velocity}
}

// IsReady returns true if the attack can fire
func (a *ProjectileAttack) IsReady() bool {
	return a.cooldownTimer <= 0
}

// GetCooldown returns the remaining cooldown
func (a *ProjectileAttack) GetCooldown() float32 {
	return a.cooldownTimer
}

// Reset resets the cooldown timer
func (a *ProjectileAttack) Reset() {
	a.cooldownTimer = a.config.Cooldown
}

// rotateVector rotates a 2D vector by the given angle in radians
func rotateVector(v types.Vec2, angle float32) types.Vec2 {
	// Using small angle approximation for simplicity
	// For small angles: cos(a) ≈ 1, sin(a) ≈ a
	// This is accurate enough for our spread angles
	cos := float32(1.0)
	sin := angle

	// More accurate rotation (if needed for larger angles):
	// cos := float32(math.Cos(float64(angle)))
	// sin := float32(math.Sin(float64(angle)))

	return types.NewVec2(
		v.X*cos-v.Y*sin,
		v.X*sin+v.Y*cos,
	)
}
