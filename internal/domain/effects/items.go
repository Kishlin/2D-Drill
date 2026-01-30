package effects

import (
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// Teleport moves player to spawn position
type Teleport struct{}

func (e Teleport) Apply(ctx *EffectContext) {
	ctx.Player.AABB.X = ctx.Player.SpawnX
	ctx.Player.AABB.Y = ctx.Player.SpawnY
	ctx.Player.Velocity = types.Zero()
	ctx.Player.OnGround = false
}

// Repair restores HP to maximum
type Repair struct{}

func (e Repair) Apply(ctx *EffectContext) {
	ctx.Player.HP = ctx.Player.MaxHP()
}

// Refuel restores fuel to capacity
type Refuel struct{}

func (e Refuel) Apply(ctx *EffectContext) {
	ctx.Player.Fuel = ctx.Player.FuelCapacity()
}

// Bomb destroys tiles and damages entities in radius
type Bomb struct {
	Radius int
	Damage float32
}

func (e Bomb) Apply(ctx *EffectContext) {
	player := ctx.Player
	centerX := player.AABB.X + player.AABB.Width/2
	centerY := player.AABB.Y + player.AABB.Height/2

	// Calculate blast AABB
	blastRadius := float32(e.Radius) * world.TileSize
	blastAABB := types.AABB{
		X:      centerX - blastRadius,
		Y:      centerY - blastRadius,
		Width:  blastRadius * 2,
		Height: blastRadius * 2,
	}

	// Damage all damageable entities in range via their hurtboxes
	for _, entity := range ctx.Damageables {
		for _, hurtbox := range entity.GetHurtboxes() {
			if hurtbox.AABB().Intersects(blastAABB) {
				entity.TakeDamageAt(hurtbox.ID, e.Damage)
				break // Only damage once per entity
			}
		}
	}

	// Destroy tiles in circular radius
	gridCenterX := int(centerX / world.TileSize)
	gridCenterY := int(centerY / world.TileSize)
	for dy := -e.Radius; dy <= e.Radius; dy++ {
		for dx := -e.Radius; dx <= e.Radius; dx++ {
			if dx*dx+dy*dy <= e.Radius*e.Radius {
				ctx.World.NukeTileAtGrid(gridCenterX+dx, gridCenterY+dy)
			}
		}
	}
}
