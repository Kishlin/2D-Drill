package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

type BossFightSystem struct {
	boss            bosses.Boss
	bossRoomStartY  float32 // Top of boss room in pixels
	bossRoomEndY    float32 // Bottom of boss room (start of floor) in pixels
	floorStartY     float32 // Top of floor area in pixels
	floorEndY       float32 // Bottom of world in pixels
	floorType       config.FloorType
	bossRoomCfg     *config.BossRoomConfig
	wasPlayerInRoom bool // Track if player was in room last frame
}

func NewBossFightSystem(boss bosses.Boss, bossRoomCfg *config.BossRoomConfig, worldHeight float32) *BossFightSystem {
	if boss == nil || bossRoomCfg == nil {
		return nil
	}

	worldHeightPixels := worldHeight
	floorHeightPixels := bossRoomCfg.FloorHeight * world.TileSize
	roomHeightPixels := bossRoomCfg.RoomHeight

	floorEndY := worldHeightPixels
	floorStartY := floorEndY - floorHeightPixels
	bossRoomEndY := floorStartY
	bossRoomStartY := bossRoomEndY - roomHeightPixels

	return &BossFightSystem{
		boss:            boss,
		bossRoomStartY:  bossRoomStartY,
		bossRoomEndY:    bossRoomEndY,
		floorStartY:     floorStartY,
		floorEndY:       floorEndY,
		floorType:       bossRoomCfg.FloorType,
		bossRoomCfg:     bossRoomCfg,
		wasPlayerInRoom: false,
	}
}

// Update handles boss fight logic
// Returns the current game state (Playing, Victory, or Defeat)
func (s *BossFightSystem) Update(player *entities.Player, dt float32) entities.GameState {
	if s == nil || s.boss == nil {
		return entities.GameStatePlaying
	}

	playerInRoom := s.IsPlayerInBossRoom(player)

	if playerInRoom && !s.wasPlayerInRoom {
		s.boss.Activate()
	}

	if !playerInRoom && s.wasPlayerInRoom {
		s.boss.Deactivate()
	}

	s.wasPlayerInRoom = playerInRoom

	s.boss.Update(player, dt)

	s.handleProjectileCollisions(player)

	s.handleFloorDamage(player)

	if s.boss.IsDefeated() {
		return entities.GameStateVictory
	}

	if player.HP <= 0 {
		return entities.GameStateDefeat
	}

	return entities.GameStatePlaying
}

func (s *BossFightSystem) IsPlayerInBossRoom(player *entities.Player) bool {
	if s == nil {
		return false
	}

	playerCenterY := player.AABB.Y + player.AABB.Height/2
	return playerCenterY >= s.bossRoomStartY && playerCenterY < s.bossRoomEndY
}

func (s *BossFightSystem) DamageBoss(damage float32) {
	if s == nil || s.boss == nil {
		return
	}

	physicalBoss, ok := s.boss.(bosses.PhysicalBoss)
	if ok {
		physicalBoss.TakeDamage(damage)
	}
}

func (s *BossFightSystem) GetBoss() bosses.Boss {
	if s == nil {
		return nil
	}
	return s.boss
}

func (s *BossFightSystem) handleProjectileCollisions(player *entities.Player) {
	projectiles := s.boss.GetProjectiles()
	for _, proj := range projectiles {
		if proj.Active && proj.Intersects(player.AABB) {
			player.DealDamage(proj.Damage)
			proj.Deactivate()
		}
	}

	// Update projectiles
	for _, proj := range projectiles {
		if proj.Active {
			proj.Update(0.016) // Approximate dt if not passed
		}
	}
}

func (s *BossFightSystem) handleFloorDamage(player *entities.Player) {
	if s.floorType != config.FloorLava {
		return
	}

	// Check if player is on floor (standing on floor tiles)
	playerBottomY := player.AABB.Y + player.AABB.Height

	// If player is on the floor tiles, apply damage
	if playerBottomY >= s.floorStartY && playerBottomY <= s.floorEndY && player.OnGround {
		// Apply damage per second (rough estimate)
		player.DealDamage(10.0) // Can be made configurable
	}
}

func (s *BossFightSystem) IsBossFightActive() bool {
	if s == nil || s.boss == nil {
		return false
	}
	return s.boss.IsActive()
}
