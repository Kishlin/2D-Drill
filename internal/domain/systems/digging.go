package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

const (
	// DigAnimationDuration is the base time to complete a dig.
	// Future: replace with computed value based on ore type and drill upgrades.
	DigAnimationDuration float32 = 1 // seconds
)

type DigDirection int

const (
	DigDown DigDirection = iota
	DigLeft
	DigRight
)

type DiggingAnimation struct {
	Active      bool
	Direction   DigDirection // Which direction we're digging
	StartX      float32
	StartY      float32
	TargetX     float32 // Tile-aligned X
	TargetY     float32 // Depends on direction
	TargetGridX int
	TargetGridY int
	Elapsed     float32
	Duration    float32
	Tile        *entities.Tile // For ore collection on completion
}

type DiggingSystem struct {
	world     *world.World
	animation DiggingAnimation
}

func NewDiggingSystem(w *world.World) *DiggingSystem {
	return &DiggingSystem{world: w}
}

// ProcessDigging handles vertical and horizontal digging with animation
func (ds *DiggingSystem) ProcessDigging(
	player *entities.Player,
	inputState input.InputState,
	dt float32,
) {
	// Update animation if in progress
	if ds.animation.Active {
		ds.updateDigAnimation(player, dt)
		return
	}

	// Handle vertical digging (S/Down key)
	if inputState.Dig && player.OnGround {
		ds.processVerticalDigging(player)
		return
	}

	// Handle horizontal digging (Left/Right when grounded)
	if player.OnGround {
		ds.processHorizontalDigging(player, inputState)
	}
}

// processVerticalDigging handles downward digging (starts animation)
func (ds *DiggingSystem) processVerticalDigging(player *entities.Player) {
	// Calculate tile beneath player's center-bottom
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height

	// Check tile directly below player
	tile := ds.world.GetTileAt(playerCenterX, playerBottomY)
	if tile == nil || !tile.IsDiggable() {
		return
	}

	// Get grid coordinates
	tileGridX := int(playerCenterX / world.TileSize)
	tileGridY := int(playerBottomY / world.TileSize)

	// Calculate target position
	tileCenterX := float32(tileGridX)*world.TileSize + world.TileSize/2
	targetX := tileCenterX - player.AABB.Width/2

	// Target Y: player bottom edge aligns with tile bottom edge
	tileBottomY := float32(tileGridY+1) * world.TileSize
	targetY := tileBottomY - player.AABB.Height

	// Start animation
	ds.startDigAnimation(player, DigDown, tileGridX, tileGridY, targetX, targetY, tile)
}

// processHorizontalDigging handles left/right digging (starts animation)
func (ds *DiggingSystem) processHorizontalDigging(
	player *entities.Player,
	inputState input.InputState,
) {
	playerCenterY := player.AABB.Y + player.AABB.Height/2

	if inputState.Left {
		// Check tile just left of player's left edge
		tileX := player.AABB.X - 1
		tile := ds.world.GetTileAt(tileX, playerCenterY)
		if tile != nil && tile.IsDiggable() {
			tileGridX := int(tileX / world.TileSize)
			tileGridY := int(playerCenterY / world.TileSize)

			// Calculate target position (center of tile)
			tileCenterX := float32(tileGridX)*world.TileSize + world.TileSize/2
			targetX := tileCenterX - player.AABB.Width/2

			// Y stays at current ground level
			targetY := player.AABB.Y

			ds.startDigAnimation(player, DigLeft, tileGridX, tileGridY, targetX, targetY, tile)
			return
		}
	}

	if inputState.Right {
		// Check tile just right of player's right edge
		tileX := player.AABB.X + player.AABB.Width + 1
		tile := ds.world.GetTileAt(tileX, playerCenterY)
		if tile != nil && tile.IsDiggable() {
			tileGridX := int(tileX / world.TileSize)
			tileGridY := int(playerCenterY / world.TileSize)

			// Calculate target position (center of tile)
			tileCenterX := float32(tileGridX)*world.TileSize + world.TileSize/2
			targetX := tileCenterX - player.AABB.Width/2

			// Y stays at current ground level
			targetY := player.AABB.Y

			ds.startDigAnimation(player, DigRight, tileGridX, tileGridY, targetX, targetY, tile)
			return
		}
	}
}

func (ds *DiggingSystem) startDigAnimation(
	player *entities.Player,
	direction DigDirection,
	tileGridX, tileGridY int,
	targetX, targetY float32,
	tile *entities.Tile,
) {
	ds.animation = DiggingAnimation{
		Active:      true,
		Direction:   direction,
		StartX:      player.AABB.X,
		StartY:      player.AABB.Y,
		TargetX:     targetX,
		TargetY:     targetY,
		TargetGridX: tileGridX,
		TargetGridY: tileGridY,
		Elapsed:     0,
		Duration:    DigAnimationDuration,
		Tile:        tile,
	}

	player.IsDigging = true

	// Zero player velocity to prevent physics interference
	player.Velocity = types.Vec2{}
}

func (ds *DiggingSystem) updateDigAnimation(player *entities.Player, dt float32) {
	ds.animation.Elapsed += dt

	// Calculate progress (0.0 to 1.0)
	progress := ds.animation.Elapsed / ds.animation.Duration
	if progress > 1.0 {
		progress = 1.0
	}

	// Lerp player position toward target
	player.AABB.X = ds.animation.StartX + (ds.animation.TargetX-ds.animation.StartX)*progress
	player.AABB.Y = ds.animation.StartY + (ds.animation.TargetY-ds.animation.StartY)*progress

	// On completion
	if progress >= 1.0 {
		ds.finishDigAnimation(player)
	}
}

func (ds *DiggingSystem) finishDigAnimation(player *entities.Player) {
	// Remove tile via grid coordinates
	if dugTile, success := ds.world.DigTileAtGrid(ds.animation.TargetGridX, ds.animation.TargetGridY); success {
		ds.collectOreIfPresent(player, dugTile)
	}

	// Reset animation state
	ds.animation = DiggingAnimation{}

	player.IsDigging = false

	// Zero player velocity to prevent physics residue
	player.Velocity = types.Vec2{}
}

// collectOreIfPresent adds ore to player inventory if the dug tile is ore
// Ore is lost if cargo is full
func (ds *DiggingSystem) collectOreIfPresent(player *entities.Player, dugTile *entities.Tile) {
	if dugTile != nil && dugTile.Type == entities.TileTypeOre {
		player.AddOre(dugTile.OreType)
		// If AddOre returns false (cargo full), ore is silently lost
	}
}
