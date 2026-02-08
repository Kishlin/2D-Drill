package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/effects"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
	"github.com/Kishlin/drill-game/internal/domain/types"
	"github.com/Kishlin/drill-game/internal/domain/world"
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
	genCfg    config.GenerationConfig
	drillCfg  config.DrillingConfig
	animation DrillingAnimation
}

func NewDrillingSystemWithConfig(w *world.World, genCfg config.GenerationConfig, drillCfg config.DrillingConfig) *DrillingSystem {
	return &DrillingSystem{
		world:    w,
		genCfg:   genCfg,
		drillCfg: drillCfg,
	}
}

// ProcessDrilling handles vertical and horizontal drilling with animation
// Returns effects to be applied (e.g., damage from hazard tiles)
func (ds *DrillingSystem) ProcessDrilling(
	player *entities.Player,
	inputState input.InputState,
	dt float32,
) []effects.Effect {
	// Update animation if in progress
	if ds.animation.Active {
		return ds.updateDrillAnimation(player, dt)
	}

	// Handle vertical drilling (S/Down key)
	if inputState.Drill && player.OnGround {
		ds.processVerticalDrilling(player)
		return nil
	}

	// Handle horizontal drilling (Left/Right when grounded)
	if player.OnGround {
		ds.processHorizontalDrilling(player, inputState)
	}

	return nil
}

// processVerticalDrilling handles downward drilling (starts animation)
func (ds *DrillingSystem) processVerticalDrilling(player *entities.Player) {
	// Calculate tile beneath player's center-bottom
	playerCenterX := player.AABB.X + player.AABB.Width/2
	playerBottomY := player.AABB.Y + player.AABB.Height

	// Check tile directly below player
	tile := ds.world.GetTileAt(playerCenterX, playerBottomY)
	if tile == nil || tile.IsDrillable() == false {
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
	// Calculate tile Y position and base drilling duration
	tileY := float32(tileGridY) * world.TileSize
	baseDuration := ds.calculateDrillingDuration(tileY, tile)

	// Calculate depth factor (0 at ground level, 1 at max depth)
	groundLevel := ds.world.GroundLevel
	maxDepth := ds.world.Height - groundLevel
	depthBelowGround := tileY - groundLevel
	depthFactor := depthBelowGround / maxDepth
	if depthFactor < 0 {
		depthFactor = 0
	}
	if depthFactor > 1 {
		depthFactor = 1
	}

	// Apply depth-interpolated drill speed
	// Speed interpolates from SpeedAtSurface (depthFactor=0) to SpeedAtMaxDepth (depthFactor=1)
	speedAtSurface := player.DrillSpeedAtSurface()
	speedAtMaxDepth := player.DrillSpeedAtMaxDepth()
	effectiveSpeed := speedAtSurface + (speedAtMaxDepth-speedAtSurface)*depthFactor
	duration := baseDuration / effectiveSpeed

	// Apply floor clamp
	if duration < ds.drillCfg.FloorDrillingDuration {
		duration = ds.drillCfg.FloorDrillingDuration
	}

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

func (ds *DrillingSystem) updateDrillAnimation(player *entities.Player, dt float32) []effects.Effect {
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
		return ds.finishDrillAnimation(player)
	}

	return nil
}

func (ds *DrillingSystem) finishDrillAnimation(player *entities.Player) []effects.Effect {
	var result []effects.Effect

	// Remove tile via grid coordinates
	if dugTile, success := ds.world.DrillTileAtGrid(ds.animation.TargetGridX, ds.animation.TargetGridY); success {
		ds.collectOreIfPresent(player, dugTile)

		// Create on-drill effect if drilling a hazard tile
		if dugTile.Type == entities.TileTypeHazard {
			if effect := ds.createOnDrillEffect(dugTile); effect != nil {
				result = append(result, effect)
			}
		}
	}

	// Reset animation state
	ds.animation = DrillingAnimation{}

	player.IsDrilling = false

	// Zero player velocity to prevent physics residue
	player.Velocity = types.Vec2{}

	return result
}

// collectOreIfPresent adds ore to player inventory if the dug tile is ore
// Ore is lost if cargo is full
func (ds *DrillingSystem) collectOreIfPresent(player *entities.Player, dugTile *entities.Tile) {
	if dugTile != nil && dugTile.Type == entities.TileTypeOre {
		player.AddOreByID(dugTile.OreID)
		// If AddOreByID returns false (cargo full), ore is silently lost
	}
}

// createOnDrillEffect creates the appropriate effect for a hazard tile based on config
func (ds *DrillingSystem) createOnDrillEffect(tile *entities.Tile) effects.Effect {
	hazardCfg := ds.genCfg.GetHazardByID(tile.HazardID)
	if hazardCfg == nil {
		return nil
	}

	switch hazardCfg.OnDrillEffect.Type {
	case config.HazardEffectDamage:
		return effects.HazardDamage{
			Damage: hazardCfg.OnDrillEffect.Damage,
		}
	case config.HazardEffectHeatDamage:
		return effects.HazardHeatDamage{
			BaseDamage:         hazardCfg.OnDrillEffect.BaseDamage,
			MaxHeatResistance:  hazardCfg.OnDrillEffect.MaxHeatResistance,
			MaxDamageReduction: hazardCfg.OnDrillEffect.MaxDamageReduction,
		}
	case config.HazardEffectMoney:
		return effects.HazardMoney{
			Amount: hazardCfg.OnDrillEffect.MoneyAmount,
		}
	case config.HazardEffectNone:
		return nil
	default:
		return nil
	}
}

// calculateDrillingDuration computes the time to drill a tile based on hardness and depth
func (ds *DrillingSystem) calculateDrillingDuration(tileY float32, tile *entities.Tile) float32 {
	// Hazard tiles use fixed duration (depth-independent) since effects are the penalty
	if tile.Type == entities.TileTypeHazard {
		hazardCfg := ds.genCfg.GetHazardByID(tile.HazardID)
		if hazardCfg != nil && hazardCfg.FixedDuration > 0 {
			return hazardCfg.FixedDuration
		}
		// Fall through to normal duration calculation if config is missing
	}

	hardness := ds.getHardness(tile)
	depthFactor := ds.calculateDepthFactor(tileY)

	return ds.drillCfg.MinDrillingDuration * hardness * depthFactor
}

// getHardness returns the hardness value for a tile based on its type
func (ds *DrillingSystem) getHardness(tile *entities.Tile) float32 {
	switch tile.Type {
	case entities.TileTypeOre:
		if oreCfg := ds.genCfg.GetOreByID(tile.OreID); oreCfg != nil {
			return oreCfg.Hardness
		}
		return 1.5 // Fallback for unknown ore types
	case entities.TileTypeDirt:
		return ds.genCfg.DirtHardness
	default:
		return ds.genCfg.DirtHardness // Fallback
	}
}

// calculateDepthFactor returns a multiplier based on depth (1.0 at surface → MaxDrillingDuration at max depth)
func (ds *DrillingSystem) calculateDepthFactor(tileY float32) float32 {
	groundLevel := ds.world.GroundLevel
	depthBelowGround := tileY - groundLevel

	// Above ground: use minimum depth factor
	if depthBelowGround <= 0 {
		return 1.0
	}

	maxDepth := ds.world.Height - groundLevel
	normalizedDepth := depthBelowGround / maxDepth

	// Clamp normalized depth to [0, 1] in case tile exceeds ds.world.Height
	if normalizedDepth > 1.0 {
		normalizedDepth = 1.0
	}

	// Linear interpolation from 1.0 to MaxDrillingDuration
	depthFactor := 1.0 + normalizedDepth*(ds.drillCfg.MaxDrillingDuration-1.0)

	return depthFactor
}
