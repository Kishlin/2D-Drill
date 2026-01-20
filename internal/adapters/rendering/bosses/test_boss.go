package bosses

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/Kishlin/drill-game/internal/domain/bosses"
	"github.com/Kishlin/drill-game/internal/domain/bosses/test_boss"
)

func init() {
	Register(&TestBossRenderer{})
}

// TestBossRenderer handles rendering for TestBoss
type TestBossRenderer struct{}

func (r *TestBossRenderer) CanRender(boss bosses.Boss) bool {
	_, ok := boss.(*test_boss.TestBoss)
	return ok
}

func (r *TestBossRenderer) Render(boss bosses.Boss) {
	tb, ok := boss.(*test_boss.TestBoss)
	if !ok {
		return
	}

	aabb := tb.GetAABB()
	drawX := aabb.X
	drawY := aabb.Y

	state := tb.GetState()
	stateTimer := tb.GetStateTimer()

	// Vibration during windup
	if state == test_boss.StateWindup {
		vibrationIntensity := float32(4.0)
		vibrationSpeed := float32(30.0)
		offset := vibrationIntensity*float32(int(stateTimer*vibrationSpeed*2)%2) - vibrationIntensity/2
		drawX += offset
	}

	// Determine color based on state
	bossColor := BossColor
	borderColor := BossBorderColor

	switch state {
	case test_boss.StateWindup:
		// Windup: flash red/orange warning
		if int(stateTimer*8)%2 == 0 {
			bossColor = rl.NewColor(255, 100, 0, 255) // Orange flash
		}
	case test_boss.StateSlam:
		// Slamming: bright red
		bossColor = rl.NewColor(255, 50, 50, 255)
	case test_boss.StateVulnerable:
		// Vulnerable: flash pink
		if int(stateTimer*4)%2 == 0 {
			bossColor = BossVulnerableColor
		}
	case test_boss.StatePatrol:
		// Patrol: check if always vulnerable (phase 1)
		if tb.IsVulnerable() {
			timer := tb.GetVulnerableTimer()
			if timer < 0 || int(timer*4)%2 == 0 {
				bossColor = BossVulnerableColor
			}
		} else {
			// Invulnerable during patrol in phases 2-3
			bossColor = BossInvulnerableColor
			borderColor = rl.DarkGray
		}
	}

	// Draw boss
	rl.DrawRectangle(
		int32(drawX),
		int32(drawY),
		int32(aabb.Width),
		int32(aabb.Height),
		bossColor,
	)

	rl.DrawRectangleLines(
		int32(drawX),
		int32(drawY),
		int32(aabb.Width),
		int32(aabb.Height),
		borderColor,
	)

	// Render AOE effects
	r.renderAOE(tb)
}

func (r *TestBossRenderer) renderAOE(tb *test_boss.TestBoss) {
	aoeInfo := tb.GetAOEInfo()
	if aoeInfo == nil {
		return
	}

	centerX := int32(aoeInfo.Position.X)
	centerY := int32(aoeInfo.Position.Y)
	radius := aoeInfo.Radius

	if aoeInfo.IsTelegraph {
		// Telegraph: pulsing yellow circle
		pulseScale := 0.8 + 0.2*float32(int(aoeInfo.StateTimer*4)%2)
		rl.DrawCircle(centerX, centerY, radius*pulseScale, AOETelegraphColor)
		rl.DrawCircleLines(centerX, centerY, radius, rl.Yellow)
	} else if aoeInfo.IsDamaging {
		// Damage: solid orange-red circle
		rl.DrawCircle(centerX, centerY, radius, AOEDamageColor)
		rl.DrawCircleLines(centerX, centerY, radius, rl.Red)
	}
}
