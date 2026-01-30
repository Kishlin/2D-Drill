package attacks

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/projectiles"
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

func TestProjectileAttack_SpawnRequestProperties(t *testing.T) {
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

	requests := attack.Update(bossAABB, playerAABB, 0.016)

	if len(requests) != 1 {
		t.Fatalf("expected 1 spawn request, got %d", len(requests))
	}

	req := requests[0]

	if req.Damage != 5.0 {
		t.Errorf("expected damage 5.0, got %f", req.Damage)
	}

	if req.Size != 16.0 {
		t.Errorf("expected size 16.0, got %f", req.Size)
	}

	if req.Movement == nil {
		t.Error("expected movement to be set")
	}

	// Verify it's a Linear movement
	if _, ok := req.Movement.(projectiles.Linear); ok == false {
		t.Error("expected Linear movement type")
	}
}
