package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type ItemSystem struct {
	world  *world.World
	spawnX float32
	spawnY float32
}

func NewItemSystem(w *world.World, spawnX, spawnY float32) *ItemSystem {
	return &ItemSystem{
		world:  w,
		spawnX: spawnX,
		spawnY: spawnY,
	}
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
		is.applyBomb(player, 2)
	}
	if inputState.UseBigBomb && player.UseItem(entities.ItemBigBomb) {
		is.applyBomb(player, 4)
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

	// Destroy tiles in circular radius (ore is lost, not collected)
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			// Circular blast check
			if dx*dx+dy*dy <= radius*radius {
				gridX, gridY := centerX+dx, centerY+dy
				is.world.DrillTileAtGrid(gridX, gridY)
			}
		}
	}
}
