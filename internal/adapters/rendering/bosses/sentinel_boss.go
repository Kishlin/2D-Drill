package bosses

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/Kishlin/drill-game/internal/domain/boss_catalog/sentinel_boss"
	"github.com/Kishlin/drill-game/internal/domain/bosses"
)

func init() {
	Register(&SentinelBossRenderer{})
}

// SentinelBossRenderer handles rendering for SentinelBoss
type SentinelBossRenderer struct{}

func (r *SentinelBossRenderer) CanRender(boss bosses.Boss) bool {
	_, ok := boss.(*sentinel_boss.SentinelBoss)
	return ok
}

func (r *SentinelBossRenderer) Render(boss bosses.Boss) {
	sb, ok := boss.(*sentinel_boss.SentinelBoss)
	if ok == false {
		return
	}

	pos := sb.GetPosition()
	drawX := pos.X
	drawY := pos.Y

	state := sb.GetState()
	stateTimer := sb.GetStateTimer()

	// Vibration during charge windup
	if state == sentinel_boss.StateChargeWindup {
		vibrationIntensity := float32(4.0)
		vibrationSpeed := float32(30.0)
		offset := vibrationIntensity*float32(int(stateTimer*vibrationSpeed*2)%2) - vibrationIntensity/2
		drawX += offset
	}

	// Determine color based on state
	bossColor := rl.NewColor(100, 100, 200, 255) // Blue-gray for sentinel
	borderColor := rl.NewColor(50, 50, 120, 255)  // Dark blue border

	switch state {
	case sentinel_boss.StateChargeWindup:
		// Flash red warning
		if int(stateTimer*8)%2 == 0 {
			bossColor = rl.NewColor(255, 100, 0, 255) // Orange flash
		}
	case sentinel_boss.StateCharge:
		// Charging: bright red
		bossColor = rl.NewColor(255, 50, 50, 255)
	case sentinel_boss.StateStunned:
		// Stunned/vulnerable: flash pink
		if int(stateTimer*4)%2 == 0 {
			bossColor = BossVulnerableColor
		}
	case sentinel_boss.StateLaserAim:
		// Aiming laser: yellow tint
		if int(stateTimer*6)%2 == 0 {
			bossColor = rl.NewColor(255, 255, 100, 255)
		}
	case sentinel_boss.StateLaser:
		// Firing laser: bright yellow
		bossColor = rl.NewColor(255, 255, 50, 255)
	case sentinel_boss.StateHover:
		// Hover: check if always vulnerable (phase 1)
		if sb.IsVulnerable() {
			timer := sb.GetVulnerableTimer()
			if timer < 0 || int(timer*4)%2 == 0 {
				bossColor = BossVulnerableColor
			}
		} else {
			bossColor = BossInvulnerableColor
			borderColor = rl.DarkGray
		}
	}

	// Draw boss body
	collisionBoxes := sb.GetCollisionBoxes()
	if len(collisionBoxes) > 0 {
		body := collisionBoxes[0]
		rl.DrawRectangle(
			int32(drawX),
			int32(drawY),
			int32(body.Width),
			int32(body.Height),
			bossColor,
		)

		rl.DrawRectangleLines(
			int32(drawX),
			int32(drawY),
			int32(body.Width),
			int32(body.Height),
			borderColor,
		)
	}

	// Render charge telegraph line
	if state == sentinel_boss.StateChargeWindup {
		target := sb.GetChargeTarget()
		rl.DrawLine(
			int32(drawX+sentinel_boss.Width/2),
			int32(drawY+sentinel_boss.Height/2),
			int32(target.X+sentinel_boss.Width/2),
			int32(target.Y+sentinel_boss.Height/2),
			rl.NewColor(255, 0, 0, 128),
		)
	}

	// Render laser effects
	r.renderLaser(sb)

	// Debug: render box outlines
	r.renderDebugBoxes(sb)
}

func (r *SentinelBossRenderer) renderLaser(sb *sentinel_boss.SentinelBoss) {
	laserInfo := sb.GetLaserInfo()
	if laserInfo == nil {
		return
	}

	startX := int32(laserInfo.StartX)
	startY := int32(laserInfo.StartY)
	endX := int32(laserInfo.EndX)
	endY := int32(laserInfo.EndY)

	if laserInfo.IsAiming {
		// Telegraph: thin yellow line
		rl.DrawLine(startX, startY, endX, endY, rl.NewColor(255, 255, 0, 180))
		rl.DrawLine(startX+1, startY, endX+1, endY, rl.NewColor(255, 255, 0, 100))
	} else if laserInfo.IsFiring {
		// Firing: thick bright red beam
		for offset := -int32(laserInfo.Width / 2); offset <= int32(laserInfo.Width/2); offset++ {
			rl.DrawLine(startX+offset, startY, endX+offset, endY, rl.NewColor(255, 50, 50, 200))
		}
		// Bright center line
		rl.DrawLine(startX, startY, endX, endY, rl.NewColor(255, 255, 255, 255))
	}
}

func (r *SentinelBossRenderer) renderDebugBoxes(sb *sentinel_boss.SentinelBoss) {
	// Collision boxes - blue
	for _, box := range sb.GetCollisionBoxes() {
		aabb := box.AABB()
		rl.DrawRectangleLinesEx(
			rl.NewRectangle(aabb.X, aabb.Y, aabb.Width, aabb.Height),
			2,
			rl.NewColor(0, 100, 255, 128),
		)
	}

	// Hitboxes - red
	for _, box := range sb.GetHitboxes() {
		aabb := box.AABB()
		rl.DrawRectangleLinesEx(
			rl.NewRectangle(aabb.X+2, aabb.Y+2, aabb.Width-4, aabb.Height-4),
			2,
			rl.NewColor(255, 50, 50, 128),
		)
	}

	// Hurtboxes - green
	for _, box := range sb.GetHurtboxes() {
		aabb := box.AABB()
		rl.DrawRectangleLinesEx(
			rl.NewRectangle(aabb.X+4, aabb.Y+4, aabb.Width-8, aabb.Height-8),
			2,
			rl.NewColor(50, 255, 50, 128),
		)
	}
}
