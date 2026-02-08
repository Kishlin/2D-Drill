package bosses

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/types"
)

func TestNewBoxSet_CreatesCorrectNumberOfBoxes(t *testing.T) {
	collisions := []BoxDef{
		{ID: "c1", Width: 10, Height: 20},
		{ID: "c2", Width: 30, Height: 40},
	}
	hitboxes := []HitboxDef{
		{BoxDef: BoxDef{ID: "h1", Width: 10, Height: 10}, DamagePerSec: 5},
	}
	hurtboxes := []HurtboxDef{
		{BoxDef: BoxDef{ID: "hb1", Width: 8, Height: 8}, DamageMultiplier: 1.0},
		{BoxDef: BoxDef{ID: "hb2", Width: 12, Height: 12}, DamageMultiplier: 2.0},
		{BoxDef: BoxDef{ID: "hb3", Width: 6, Height: 6}, DamageMultiplier: 0.5},
	}

	bs := NewBoxSet(collisions, hitboxes, hurtboxes)

	if len(bs.CollisionBoxes) != 2 {
		t.Errorf("Expected 2 collision boxes, got %d", len(bs.CollisionBoxes))
	}
	if len(bs.Hitboxes) != 1 {
		t.Errorf("Expected 1 hitbox, got %d", len(bs.Hitboxes))
	}
	if len(bs.Hurtboxes) != 3 {
		t.Errorf("Expected 3 hurtboxes, got %d", len(bs.Hurtboxes))
	}
}

func TestNewBoxSet_InitializesCollisionBoxProperties(t *testing.T) {
	collisions := []BoxDef{
		{ID: "wall", Width: 50, Height: 100},
	}

	bs := NewBoxSet(collisions, nil, nil)

	cb := bs.CollisionBoxes[0]
	if cb.ID != "wall" {
		t.Errorf("Expected ID 'wall', got %q", cb.ID)
	}
	if cb.Width != 50 {
		t.Errorf("Expected Width 50, got %f", cb.Width)
	}
	if cb.Height != 100 {
		t.Errorf("Expected Height 100, got %f", cb.Height)
	}
}

func TestNewBoxSet_InitializesHitboxProperties(t *testing.T) {
	hitboxes := []HitboxDef{
		{BoxDef: BoxDef{ID: "attack", Width: 30, Height: 40}, DamagePerSec: 15},
	}

	bs := NewBoxSet(nil, hitboxes, nil)

	hb := bs.Hitboxes[0]
	if hb.ID != "attack" {
		t.Errorf("Expected ID 'attack', got %q", hb.ID)
	}
	if hb.Width != 30 {
		t.Errorf("Expected Width 30, got %f", hb.Width)
	}
	if hb.Height != 40 {
		t.Errorf("Expected Height 40, got %f", hb.Height)
	}
	if hb.DamagePerSec != 15 {
		t.Errorf("Expected DamagePerSec 15, got %f", hb.DamagePerSec)
	}
}

func TestNewBoxSet_InitializesHurtboxProperties(t *testing.T) {
	hurtboxes := []HurtboxDef{
		{BoxDef: BoxDef{ID: "weak", Width: 20, Height: 25}, DamageMultiplier: 2.0},
	}

	bs := NewBoxSet(nil, nil, hurtboxes)

	hb := bs.Hurtboxes[0]
	if hb.ID != "weak" {
		t.Errorf("Expected ID 'weak', got %q", hb.ID)
	}
	if hb.Width != 20 {
		t.Errorf("Expected Width 20, got %f", hb.Width)
	}
	if hb.Height != 25 {
		t.Errorf("Expected Height 25, got %f", hb.Height)
	}
	if hb.DamageMultiplier != 2.0 {
		t.Errorf("Expected DamageMultiplier 2.0, got %f", hb.DamageMultiplier)
	}
}

func TestNewBoxSet_EmptyDefs(t *testing.T) {
	bs := NewBoxSet(nil, nil, nil)

	if len(bs.CollisionBoxes) != 0 {
		t.Errorf("Expected 0 collision boxes, got %d", len(bs.CollisionBoxes))
	}
	if len(bs.Hitboxes) != 0 {
		t.Errorf("Expected 0 hitboxes, got %d", len(bs.Hitboxes))
	}
	if len(bs.Hurtboxes) != 0 {
		t.Errorf("Expected 0 hurtboxes, got %d", len(bs.Hurtboxes))
	}
}

func TestBoxSet_UpdatePositions_CollisionBoxes(t *testing.T) {
	collisions := []BoxDef{
		{ID: "c1", OffsetX: 10, OffsetY: 20, Width: 50, Height: 60},
	}
	bs := NewBoxSet(collisions, nil, nil)

	bs.UpdatePositions(100, 200)

	cb := bs.CollisionBoxes[0]
	if cb.X != 110 {
		t.Errorf("Expected X 110, got %f", cb.X)
	}
	if cb.Y != 220 {
		t.Errorf("Expected Y 220, got %f", cb.Y)
	}
}

func TestBoxSet_UpdatePositions_Hitboxes(t *testing.T) {
	hitboxes := []HitboxDef{
		{BoxDef: BoxDef{ID: "h1", OffsetX: -5, OffsetY: 15, Width: 30, Height: 30}, DamagePerSec: 10},
	}
	bs := NewBoxSet(nil, hitboxes, nil)

	bs.UpdatePositions(50, 80)

	hb := bs.Hitboxes[0]
	if hb.X != 45 {
		t.Errorf("Expected X 45, got %f", hb.X)
	}
	if hb.Y != 95 {
		t.Errorf("Expected Y 95, got %f", hb.Y)
	}
}

func TestBoxSet_UpdatePositions_Hurtboxes(t *testing.T) {
	hurtboxes := []HurtboxDef{
		{BoxDef: BoxDef{ID: "hb1", OffsetX: 0, OffsetY: -10, Width: 20, Height: 20}, DamageMultiplier: 1.0},
	}
	bs := NewBoxSet(nil, nil, hurtboxes)

	bs.UpdatePositions(300, 400)

	hb := bs.Hurtboxes[0]
	if hb.X != 300 {
		t.Errorf("Expected X 300, got %f", hb.X)
	}
	if hb.Y != 390 {
		t.Errorf("Expected Y 390, got %f", hb.Y)
	}
}

func TestCollisionBox_AABB(t *testing.T) {
	cb := CollisionBox{ID: "test", X: 10, Y: 20, Width: 30, Height: 40}
	aabb := cb.AABB()

	expected := types.AABB{X: 10, Y: 20, Width: 30, Height: 40}
	if aabb != expected {
		t.Errorf("Expected AABB %+v, got %+v", expected, aabb)
	}
}

func TestHitbox_AABB(t *testing.T) {
	hb := Hitbox{ID: "test", X: 5, Y: 15, Width: 25, Height: 35, DamagePerSec: 10}
	aabb := hb.AABB()

	expected := types.AABB{X: 5, Y: 15, Width: 25, Height: 35}
	if aabb != expected {
		t.Errorf("Expected AABB %+v, got %+v", expected, aabb)
	}
}

func TestHurtbox_AABB(t *testing.T) {
	hb := Hurtbox{ID: "test", X: 8, Y: 12, Width: 16, Height: 24, DamageMultiplier: 1.5}
	aabb := hb.AABB()

	expected := types.AABB{X: 8, Y: 12, Width: 16, Height: 24}
	if aabb != expected {
		t.Errorf("Expected AABB %+v, got %+v", expected, aabb)
	}
}

func TestNewBodyBoxSet_CreatesMatchingBoxes(t *testing.T) {
	bs := NewBodyBoxSet(BodyBoxConfig{
		ID:               "body",
		Width:            100,
		Height:           80,
		OffsetX:          5,
		OffsetY:          10,
		DamagePerSec:     20,
		DamageMultiplier: 1.5,
	})

	if len(bs.CollisionBoxes) != 1 {
		t.Fatalf("Expected 1 collision box, got %d", len(bs.CollisionBoxes))
	}
	if len(bs.Hitboxes) != 1 {
		t.Fatalf("Expected 1 hitbox, got %d", len(bs.Hitboxes))
	}
	if len(bs.Hurtboxes) != 1 {
		t.Fatalf("Expected 1 hurtbox, got %d", len(bs.Hurtboxes))
	}

	// Check collision box
	cb := bs.CollisionBoxes[0]
	if cb.ID != "body" {
		t.Errorf("Collision box ID: expected 'body', got %q", cb.ID)
	}
	if cb.Width != 100 || cb.Height != 80 {
		t.Errorf("Collision box dims: expected 100x80, got %fx%f", cb.Width, cb.Height)
	}

	// Check hitbox
	hb := bs.Hitboxes[0]
	if hb.DamagePerSec != 20 {
		t.Errorf("Hitbox DamagePerSec: expected 20, got %f", hb.DamagePerSec)
	}

	// Check hurtbox
	hr := bs.Hurtboxes[0]
	if hr.DamageMultiplier != 1.5 {
		t.Errorf("Hurtbox DamageMultiplier: expected 1.5, got %f", hr.DamageMultiplier)
	}

	// Verify UpdatePositions works
	bs.UpdatePositions(200, 300)
	if bs.CollisionBoxes[0].X != 205 {
		t.Errorf("After UpdatePositions, collision box X: expected 205, got %f", bs.CollisionBoxes[0].X)
	}
	if bs.CollisionBoxes[0].Y != 310 {
		t.Errorf("After UpdatePositions, collision box Y: expected 310, got %f", bs.CollisionBoxes[0].Y)
	}
}
