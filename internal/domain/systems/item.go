package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type ItemSystem struct {
	world         *world.World
	spawnX        float32
	spawnY        float32
	bombRadius    int
	bigBombRadius int
	bossFightSystem *BossFightSystem
}

func NewItemSystemWithConfig(w *world.World, spawnX, spawnY float32, itemCfg config.ItemConfig) *ItemSystem {
	return &ItemSystem{
		world:         w,
		spawnX:        spawnX,
		spawnY:        spawnY,
		bombRadius:    itemCfg.Bomb.Radius,
		bigBombRadius: itemCfg.BigBomb.Radius,
		bossFightSystem: nil,
	}
}

// SetBossFightSystem sets the boss fight system reference for bomb-boss collisions
func (is *ItemSystem) SetBossFightSystem(bfs *BossFightSystem) {
	is.bossFightSystem = bfs
}

// ProcessItemUsage checks for item inputs and applies effects
func (is *ItemSystem) ProcessItemUsage(player *entities.Player, inputState input.InputState) {
	if inputState.UseTeleport && player.UseItem(entities.ItemTeleport) {
		is.applyTeleport(player)
	}
	if inputState.UseRepair && player.UseItem(entities.ItemRepair) {
		is.applyRepair(player)
	}
	if inputState.UseRefuel && player.UseItem(entities.ItemRefuel) {
		is.applyRefuel(player)
	}
	if inputState.UseBomb && player.UseItem(entities.ItemBomb) {
		is.applyBomb(player, is.bombRadius)
	}
	if inputState.UseBigBomb && player.UseItem(entities.ItemBigBomb) {
		is.applyBomb(player, is.bigBombRadius)
	}
}

func (is *ItemSystem) applyTeleport(player *entities.Player) {
	player.AABB.X = is.spawnX
	player.AABB.Y = is.spawnY
	player.Velocity = types.Zero()
	player.OnGround = false
}

func (is *ItemSystem) applyRepair(player *entities.Player) {
	player.HP = player.Hull.MaxHP()
}

func (is *ItemSystem) applyRefuel(player *entities.Player) {
	player.Fuel = player.FuelTank.Capacity()
}

func (is *ItemSystem) applyBomb(player *entities.Player, radius int) {
	// Calculate player center in grid coordinates
	centerX := int((player.AABB.X + player.AABB.Width/2) / world.TileSize)
	centerY := int((player.AABB.Y + player.AABB.Height/2) / world.TileSize)

	// Check for bomb-boss collision
	if is.bossFightSystem != nil {
		damage := float32(0)
		if radius == is.bombRadius {
			damage = 10.0 // Bomb damage
		} else {
			damage = 25.0 // Big bomb damage
		}

		// Calculate bomb AABB (circular blast)
		blastRadius := float32(radius) * world.TileSize
		bombAABB := types.AABB{
			X:      player.AABB.X + player.AABB.Width/2 - blastRadius,
			Y:      player.AABB.Y + player.AABB.Height/2 - blastRadius,
			Width:  blastRadius * 2,
			Height: blastRadius * 2,
		}

		// Check for physical boss collision
		physicalBoss, ok := is.bossFightSystem.GetBoss().(bosses.PhysicalBoss)
		if ok && physicalBoss.GetAABB().Intersects(bombAABB) {
			is.bossFightSystem.DamageBoss(damage)
		}
	}

	// Destroy tiles in circular radius (ore is lost, not collected)
	// Use NukeTileAtGrid to bypass drillability check, allowing bombs to destroy rocks
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			// Circular blast check
			if dx*dx+dy*dy <= radius*radius {
				gridX, gridY := centerX+dx, centerY+dy
				is.world.NukeTileAtGrid(gridX, gridY)
			}
		}
	}
}
