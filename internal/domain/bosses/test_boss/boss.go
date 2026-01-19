package test_boss

import (
	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

const (
	MaxHP            = 100.0
	Width            = 200.0
	Height           = 200.0
	DamagePerBomb    = 10.0
	DamagePerBigBomb = 25.0
)

// TestBoss is a simple physical boss for testing
type TestBoss struct {
	aabb        types.AABB
	hp          float32
	active      bool
	defeated    bool
	projectiles []*bosses.Projectile
}

// New creates a new test boss positioned in the center of the boss room
func New(roomStartY, worldWidth float32) *TestBoss {
	centerX := (worldWidth - Width) / 2
	centerY := roomStartY + 100

	return &TestBoss{
		aabb: types.AABB{
			X:      centerX,
			Y:      centerY,
			Width:  Width,
			Height: Height,
		},
		hp:          MaxHP,
		active:      false,
		defeated:    false,
		projectiles: make([]*bosses.Projectile, 0),
	}
}

// Update runs boss AI
// For now, the test boss just stands there
func (b *TestBoss) Update(player *entities.Player, dt float32) {
	if !b.active || b.defeated {
		return
	}

	// Test boss doesn't have any AI or projectiles yet
	// This is just a simple target for bombs
}

// GetHP returns the boss's current health
func (b *TestBoss) GetHP() float32 {
	return b.hp
}

// GetMaxHP returns the boss's maximum health
func (b *TestBoss) GetMaxHP() float32 {
	return MaxHP
}

// IsDefeated returns true if the boss is dead
func (b *TestBoss) IsDefeated() bool {
	return b.defeated
}

// IsActive returns true if the boss fight is active
func (b *TestBoss) IsActive() bool {
	return b.active
}

// Activate starts the boss fight
func (b *TestBoss) Activate() {
	b.active = true
}

// Deactivate pauses the boss fight
func (b *TestBoss) Deactivate() {
	b.active = false
}

// GetProjectiles returns the boss's projectiles
func (b *TestBoss) GetProjectiles() []*bosses.Projectile {
	return b.projectiles
}

// GetAABB returns the boss's bounding box
func (b *TestBoss) GetAABB() types.AABB {
	return b.aabb
}

// TakeDamage applies damage to the boss
func (b *TestBoss) TakeDamage(damage float32) {
	b.hp -= damage
	if b.hp <= 0 {
		b.hp = 0
		b.defeated = true
	}
}
