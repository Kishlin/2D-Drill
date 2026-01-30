package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/world"
)

// BossFightResult contains the results of a boss fight update
type BossFightResult struct {
	GameState     entities.GameState
	SpawnRequests []projectiles.SpawnRequest
}

type BossFightSystem struct {
	boss            bosses.Boss
	bossRoomStartY  float32 // Top of boss room in pixels
	bossRoomEndY    float32 // Bottom of boss room (start of floor) in pixels
	floorStartY     float32 // Top of floor area in pixels
	floorEndY       float32 // Bottom of world in pixels
	floorType       config.FloorType
	bossRoomCfg     config.BossRoomConfig
	wasPlayerInRoom bool // Track if player was in room last frame
}

func NewBossFightSystem(boss bosses.Boss, bossRoomCfg config.BossRoomConfig, worldHeight float32) *BossFightSystem {
	if boss == nil {
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
// Returns game state and any projectile spawn requests from the boss
func (s *BossFightSystem) Update(player *entities.Player, dt float32) BossFightResult {
	playerInRoom := s.IsPlayerInBossRoom(player)

	if playerInRoom && s.wasPlayerInRoom == false {
		s.boss.Activate()
	}

	if playerInRoom == false && s.wasPlayerInRoom {
		s.boss.Deactivate()
	}

	s.wasPlayerInRoom = playerInRoom

	spawnRequests := s.boss.Update(player, dt)

	s.handleContactDamage(player, dt)

	s.handleFloorDamage(player)

	var gameState entities.GameState
	if s.boss.IsDefeated() {
		gameState = entities.GameStateVictory
	} else if player.HP <= 0 {
		gameState = entities.GameStateDefeat
	} else {
		gameState = entities.GameStatePlaying
	}

	return BossFightResult{
		GameState:     gameState,
		SpawnRequests: spawnRequests,
	}
}

func (s *BossFightSystem) IsPlayerInBossRoom(player *entities.Player) bool {
	playerCenterY := player.AABB.Y + player.AABB.Height/2
	return playerCenterY >= s.bossRoomStartY && playerCenterY < s.bossRoomEndY
}

func (s *BossFightSystem) DamageBoss(damage float32) {
	// Damage the first available hurtbox
	hurtboxes := s.boss.GetHurtboxes()
	if len(hurtboxes) > 0 {
		s.boss.TakeDamageAt(hurtboxes[0].ID, damage)
	}
}

func (s *BossFightSystem) GetBoss() bosses.Boss {
	return s.boss
}

func (s *BossFightSystem) handleContactDamage(player *entities.Player, dt float32) {
	for _, hitbox := range s.boss.GetHitboxes() {
		if hitbox.DamagePerSec <= 0 {
			continue
		}

		if player.AABB.Intersects(hitbox.AABB()) {
			player.DealDamage(hitbox.DamagePerSec * dt)
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
	return s.boss.IsActive()
}
