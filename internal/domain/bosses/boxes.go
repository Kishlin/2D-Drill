package bosses

import "github.com/Kishlin/drill-game/internal/domain/types"

// CollisionBox represents a physical presence that blocks player movement
type CollisionBox struct {
	ID     string
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// AABB returns the bounding box for this collision box
func (c CollisionBox) AABB() types.AABB {
	return types.AABB{
		X:      c.X,
		Y:      c.Y,
		Width:  c.Width,
		Height: c.Height,
	}
}

// Hitbox represents an attack zone that damages the player on intersection
type Hitbox struct {
	ID           string
	X            float32
	Y            float32
	Width        float32
	Height       float32
	DamagePerSec float32
}

// AABB returns the bounding box for this hitbox
func (h Hitbox) AABB() types.AABB {
	return types.AABB{
		X:      h.X,
		Y:      h.Y,
		Width:  h.Width,
		Height: h.Height,
	}
}

// Hurtbox represents a vulnerable zone where the boss receives damage
type Hurtbox struct {
	ID               string
	X                float32
	Y                float32
	Width            float32
	Height           float32
	DamageMultiplier float32 // 1.0 = normal, 2.0 = weak point
}

// AABB returns the bounding box for this hurtbox
func (h Hurtbox) AABB() types.AABB {
	return types.AABB{
		X:      h.X,
		Y:      h.Y,
		Width:  h.Width,
		Height: h.Height,
	}
}

// BoxDef defines a box's static properties relative to boss position
type BoxDef struct {
	ID      string
	OffsetX float32 // Relative to boss position
	OffsetY float32
	Width   float32
	Height  float32
}

// HitboxDef defines a hitbox with damage properties
type HitboxDef struct {
	BoxDef
	DamagePerSec float32
}

// HurtboxDef defines a hurtbox with damage multiplier
type HurtboxDef struct {
	BoxDef
	DamageMultiplier float32
}

// BoxSet manages pre-allocated boxes with position synchronization
type BoxSet struct {
	// Definitions (static)
	collisionDefs []BoxDef
	hitboxDefs    []HitboxDef
	hurtboxDefs   []HurtboxDef

	// Runtime boxes (pre-allocated, positions updated)
	CollisionBoxes []CollisionBox
	Hitboxes       []Hitbox
	Hurtboxes      []Hurtbox
}

// NewBoxSet creates a BoxSet with the given definitions
func NewBoxSet(collisions []BoxDef, hitboxes []HitboxDef, hurtboxes []HurtboxDef) *BoxSet {
	bs := &BoxSet{
		collisionDefs:  collisions,
		hitboxDefs:     hitboxes,
		hurtboxDefs:    hurtboxes,
		CollisionBoxes: make([]CollisionBox, len(collisions)),
		Hitboxes:       make([]Hitbox, len(hitboxes)),
		Hurtboxes:      make([]Hurtbox, len(hurtboxes)),
	}

	// Initialize runtime boxes from definitions
	for i, def := range collisions {
		bs.CollisionBoxes[i] = CollisionBox{
			ID:     def.ID,
			Width:  def.Width,
			Height: def.Height,
		}
	}
	for i, def := range hitboxes {
		bs.Hitboxes[i] = Hitbox{
			ID:           def.BoxDef.ID,
			Width:        def.BoxDef.Width,
			Height:       def.BoxDef.Height,
			DamagePerSec: def.DamagePerSec,
		}
	}
	for i, def := range hurtboxes {
		bs.Hurtboxes[i] = Hurtbox{
			ID:               def.BoxDef.ID,
			Width:            def.BoxDef.Width,
			Height:           def.BoxDef.Height,
			DamageMultiplier: def.DamageMultiplier,
		}
	}

	return bs
}

// UpdatePositions updates all box positions based on boss position
func (bs *BoxSet) UpdatePositions(bossX, bossY float32) {
	for i, def := range bs.collisionDefs {
		bs.CollisionBoxes[i].X = bossX + def.OffsetX
		bs.CollisionBoxes[i].Y = bossY + def.OffsetY
	}
	for i, def := range bs.hitboxDefs {
		bs.Hitboxes[i].X = bossX + def.BoxDef.OffsetX
		bs.Hitboxes[i].Y = bossY + def.BoxDef.OffsetY
	}
	for i, def := range bs.hurtboxDefs {
		bs.Hurtboxes[i].X = bossX + def.BoxDef.OffsetX
		bs.Hurtboxes[i].Y = bossY + def.BoxDef.OffsetY
	}
}

// BodyBoxConfig provides a convenient way to define a single body box
type BodyBoxConfig struct {
	ID               string
	Width, Height    float32
	OffsetX, OffsetY float32
	DamagePerSec     float32
	DamageMultiplier float32
}

// NewBodyBoxSet creates a BoxSet for a simple body-only configuration
func NewBodyBoxSet(cfg BodyBoxConfig) *BoxSet {
	boxDef := BoxDef{
		ID:      cfg.ID,
		OffsetX: cfg.OffsetX,
		OffsetY: cfg.OffsetY,
		Width:   cfg.Width,
		Height:  cfg.Height,
	}

	return NewBoxSet(
		[]BoxDef{boxDef},
		[]HitboxDef{{BoxDef: boxDef, DamagePerSec: cfg.DamagePerSec}},
		[]HurtboxDef{{BoxDef: boxDef, DamageMultiplier: cfg.DamageMultiplier}},
	)
}
