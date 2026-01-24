package attacks

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/types"
)

func TestProjectileAttack_IsReadyInitially(t *testing.T) {
	cfg := ProjectileAttackConfig{
		Cooldown:        2.0,
		ProjectileCount: 3,
		ProjectileSpeed: 200.0,
		ProjectileSize:  16.0,
		Damage:          5.0,
	}
	attack := NewProjectileAttack(cfg)

	if attack.IsReady() == false {
		t.Error("expected attack to be ready initially")
	}
}

func TestProjectileAttack_FiresProjectiles(t *testing.T) {
	cfg := ProjectileAttackConfig{
		Cooldown:        2.0,
		ProjectileCount: 3,
		ProjectileSpeed: 200.0,
		ProjectileSize:  16.0,
		Damage:          5.0,
	}
	attack := NewProjectileAttack(cfg)

	bossAABB := types.NewAABB(100, 100, 100, 100)
	playerAABB := types.NewAABB(300, 100, 54, 54)

	projectiles := attack.Update(bossAABB, playerAABB, 0.016)

	if len(projectiles) != 3 {
		t.Errorf("expected 3 projectiles, got %d", len(projectiles))
	}

	// Check that attack is on cooldown
	if attack.IsReady() {
		t.Error("expected attack to be on cooldown after firing")
	}
}

func TestProjectileAttack_CooldownDecreases(t *testing.T) {
	cfg := ProjectileAttackConfig{
		Cooldown:        2.0,
		ProjectileCount: 1,
		ProjectileSpeed: 200.0,
		ProjectileSize:  16.0,
		Damage:          5.0,
	}
	attack := NewProjectileAttack(cfg)

	bossAABB := types.NewAABB(100, 100, 100, 100)
	playerAABB := types.NewAABB(300, 100, 54, 54)

	// Fire once to start cooldown
	attack.Update(bossAABB, playerAABB, 0.016)

	initialCooldown := attack.GetCooldown()

	// Update without player (just to pass dt)
	attack.Update(bossAABB, playerAABB, 0.5)

	if attack.GetCooldown() >= initialCooldown {
		t.Error("expected cooldown to decrease")
	}
}

func TestProjectileAttack_ProjectileProperties(t *testing.T) {
	cfg := ProjectileAttackConfig{
		Cooldown:        2.0,
		ProjectileCount: 1,
		ProjectileSpeed: 200.0,
		ProjectileSize:  16.0,
		Damage:          5.0,
	}
	attack := NewProjectileAttack(cfg)

	bossAABB := types.NewAABB(100, 100, 100, 100)
	playerAABB := types.NewAABB(300, 100, 54, 54)

	projectiles := attack.Update(bossAABB, playerAABB, 0.016)

	if len(projectiles) != 1 {
		t.Fatalf("expected 1 projectile, got %d", len(projectiles))
	}

	proj := projectiles[0]

	if proj.Damage != 5.0 {
		t.Errorf("expected damage 5.0, got %f", proj.Damage)
	}

	if proj.AABB.Width != 16.0 || proj.AABB.Height != 16.0 {
		t.Errorf("expected size 16x16, got %fx%f", proj.AABB.Width, proj.AABB.Height)
	}

	if proj.Active == false {
		t.Error("expected projectile to be active")
	}
}
