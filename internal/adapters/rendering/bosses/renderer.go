package bosses

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/Kishlin/drill-game/internal/domain/bosses"
)

// Colors shared across boss renderers
var (
	BossColor             = rl.NewColor(255, 0, 0, 255)     // Red
	BossBorderColor       = rl.NewColor(139, 0, 0, 255)     // Dark Red
	BossVulnerableColor   = rl.NewColor(255, 200, 200, 255) // Light pink
	BossInvulnerableColor = rl.NewColor(150, 150, 150, 255) // Gray
	AOETelegraphColor     = rl.NewColor(255, 255, 0, 128)   // Semi-transparent yellow
	AOEDamageColor        = rl.NewColor(255, 100, 0, 200)   // Orange-red
)

// Renderer handles rendering for a specific boss type
type Renderer interface {
	// CanRender returns true if this renderer handles the given boss
	CanRender(boss bosses.Boss) bool

	// Render draws the boss
	Render(boss bosses.Boss)
}

// registry holds all registered boss renderers
var registry []Renderer

// Register adds a boss renderer to the registry
func Register(renderer Renderer) {
	registry = append(registry, renderer)
}

// RenderBoss finds and uses the appropriate renderer for the boss
// Returns true if a specific renderer was found and used
func RenderBoss(boss bosses.Boss) bool {
	for _, renderer := range registry {
		if renderer.CanRender(boss) {
			renderer.Render(boss)
			return true
		}
	}
	return false
}

// RenderGeneric provides fallback rendering for any Boss using collision boxes
func RenderGeneric(boss bosses.Boss) {
	collisionBoxes := boss.GetCollisionBoxes()

	if len(collisionBoxes) == 0 {
		return
	}

	// Check vulnerability via hurtbox presence
	color := BossColor
	if len(boss.GetHurtboxes()) > 0 {
		color = BossVulnerableColor
	}

	for _, box := range collisionBoxes {
		aabb := box.AABB()

		rl.DrawRectangle(
			int32(aabb.X),
			int32(aabb.Y),
			int32(aabb.Width),
			int32(aabb.Height),
			color,
		)

		rl.DrawRectangleLines(
			int32(aabb.X),
			int32(aabb.Y),
			int32(aabb.Width),
			int32(aabb.Height),
			BossBorderColor,
		)
	}
}
