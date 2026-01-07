package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/physics"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

const (
	minDrillingDuration = 0.8  // seconds (at ground level)
	maxDrillingDuration = 30.0 // seconds (at max depth)
)

type DrillDirection int

const (
	DrillDown DrillDirection = iota
	DrillLeft
	DrillRight
)

type DrillingAnimation struct {
	Active      bool
	Direction   DrillDirection // Which direction we're drilling
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

type DrillingSystem struct {
	world     *world.World
	animation DrillingAnimation
}

func NewDrillingSystem(w *world.World) *DrillingSystem {
	return &DrillingSystem{world: w}
}

// ProcessDrilling handles vertical and horizontal drilling with animation
func (ds *DrillingSystem) ProcessDrilling(
	player *entities.Player,
	inputState input.InputState,
	dt float32,
) {
	// Update animation if in progress
	if ds.animation.Active {
		ds.updateDrillAnimation(player, dt)
		return
	}

	// Handle vertical drilling (S/Down key)
	if inputState.Drill && player.OnGround {
		ds.processVerticalDrilling(player)
		return
	}

	// Handle horizontal drilling (Left/Right when grounded)
	if player.OnGround {
		ds.processHorizontalDrilling(player, inputState)
	}
}

// processVerticalDrilling handles downward drilling (starts animation)
func (ds *DrillingSystem) processVerticalDrilling(player *entities.Player) {
	// Calculate tile beneath player's center-bottom
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height

	// Check tile directly below player
	tile := ds.world.GetTileAt(playerCenterX, playerBottomY)
	if tile == nil || !tile.IsDrillable() {
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
	ds.startDrillAnimation(player, DrillDown, tileGridX, tileGridY, targetX, targetY, tile)
}

// processHorizontalDrilling handles left/right drilling (starts animation)
func (ds *DrillingSystem) processHorizontalDrilling(
	player *entities.Player,
	inputState input.InputState,
) {
	playerCenterY := player.AABB.Y + player.AABB.Height/2

	if inputState.Left {
		// Check tile just left of player's left edge
		tileX := player.AABB.X - 1
		tile := ds.world.GetTileAt(tileX, playerCenterY)
		if tile != nil && tile.IsDrillable() {
			tileGridX := int(tileX / world.TileSize)
			tileGridY := int(playerCenterY / world.TileSize)

			// Calculate target position (center of tile)
			tileCenterX := float32(tileGridX)*world.TileSize + world.TileSize/2
			targetX := tileCenterX - player.AABB.Width/2

			// Y stays at current ground level
			targetY := player.AABB.Y

			ds.startDrillAnimation(player, DrillLeft, tileGridX, tileGridY, targetX, targetY, tile)
			return
		}
	}

	if inputState.Right {
		// Check tile just right of player's right edge
		tileX := player.AABB.X + player.AABB.Width + 1
		tile := ds.world.GetTileAt(tileX, playerCenterY)
		if tile != nil && tile.IsDrillable() {
			tileGridX := int(tileX / world.TileSize)
			tileGridY := int(playerCenterY / world.TileSize)

			// Calculate target position (center of tile)
			tileCenterX := float32(tileGridX)*world.TileSize + world.TileSize/2
			targetX := tileCenterX - player.AABB.Width/2

			// Y stays at current ground level
			targetY := player.AABB.Y

			ds.startDrillAnimation(player, DrillRight, tileGridX, tileGridY, targetX, targetY, tile)
			return
		}
	}
}

func (ds *DrillingSystem) startDrillAnimation(
	player *entities.Player,
	direction DrillDirection,
	tileGridX, tileGridY int,
	targetX, targetY float32,
	tile *entities.Tile,
) {
	// Calculate tile Y position and drilling duration
	tileY := float32(tileGridY) * world.TileSize
	duration := ds.calculateDrillingDuration(tileY, tile)

	ds.animation = DrillingAnimation{
		Active:      true,
		Direction:   direction,
		StartX:      player.AABB.X,
		StartY:      player.AABB.Y,
		TargetX:     targetX,
		TargetY:     targetY,
		TargetGridX: tileGridX,
		TargetGridY: tileGridY,
		Elapsed:     0,
		Duration:    duration,
		Tile:        tile,
	}

	player.IsDrilling = true

	// Zero player velocity to prevent physics interference
	player.Velocity = types.Vec2{}
}

func (ds *DrillingSystem) updateDrillAnimation(player *entities.Player, dt float32) {
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
		ds.finishDrillAnimation(player)
	}
}

func (ds *DrillingSystem) finishDrillAnimation(player *entities.Player) {
	// Remove tile via grid coordinates
	if dugTile, success := ds.world.DrillTileAtGrid(ds.animation.TargetGridX, ds.animation.TargetGridY); success {
		ds.collectOreIfPresent(player, dugTile)
	}

	// Reset animation state
	ds.animation = DrillingAnimation{}

	player.IsDrilling = false

	// Zero player velocity to prevent physics residue
	player.Velocity = types.Vec2{}
}

// collectOreIfPresent adds ore to player inventory if the dug tile is ore
// Ore is lost if cargo is full
func (ds *DrillingSystem) collectOreIfPresent(player *entities.Player, dugTile *entities.Tile) {
	if dugTile != nil && dugTile.Type == entities.TileTypeOre {
		player.AddOre(dugTile.OreType)
		// If AddOre returns false (cargo full), ore is silently lost
	}
}

// calculateDrillingDuration computes the time to drill a tile based on depth and type
func (ds *DrillingSystem) calculateDrillingDuration(tileY float32, tile *entities.Tile) float32 {
	baseDuration := ds.calculateBaseDuration(tileY)

	// Apply ore hardness multiplier if applicable
	if tile.Type == entities.TileTypeOre {
		hardness, ok := entities.OreHardness[tile.OreType]
		if !ok {
			hardness = 1.5 // Fallback for unknown ore types
		}
		return baseDuration * hardness
	}

	return baseDuration
}

// calculateBaseDuration computes drilling time for dirt based on depth
// Linear interpolation: 1 second at ground level, 30 seconds at max depth
func (ds *DrillingSystem) calculateBaseDuration(tileY float32) float32 {
	groundLevel := ds.world.GroundLevel
	depthBelowGround := tileY - groundLevel

	// Above ground: use minimum duration
	if depthBelowGround <= 0 {
		return minDrillingDuration
	}

	maxDepth := physics.MaxUndergroundY - groundLevel
	normalizedDepth := depthBelowGround / maxDepth

	// Clamp normalized depth to [0, 1] in case tile exceeds MaxUndergroundY
	if normalizedDepth > 1.0 {
		normalizedDepth = 1.0
	}

	// Linear interpolation
	duration := minDrillingDuration + normalizedDepth*(maxDrillingDuration-minDrillingDuration)

	return duration
}
