# Rendering Layer Implementation Plan

**Status:** Pending
**Estimated Complexity:** Medium (Raylib integration)
**Dependencies:** All domain/system layers complete and tested

## Overview

The rendering layer is the final step to make the boss fight system visible. The domain layer is complete and tested; this task adds visual elements to Raylib.

## Implementation Steps

### Step 1: Add Boss Getter to Game

Already complete. Access boss via:
```go
boss := game.GetBoss()
bossFightSystem := game.GetBossFightSystem()
gameState := game.GetGameState()
```

### Step 2: Render Boss Rectangle

**File:** `internal/adapters/rendering/raylib.go`
**Method:** Add to `Render()` function or create `renderBoss()` helper

```go
func (r *RaylibRenderer) renderBoss(boss bosses.Boss) {
    if boss == nil {
        return
    }

    physicalBoss, ok := boss.(bosses.PhysicalBoss)
    if !ok {
        return // Skip non-physical bosses for now
    }

    aabb := physicalBoss.GetAABB()

    // Convert world coordinates to screen coordinates
    screenX := aabb.X
    screenY := aabb.Y - worldOffsetY // Adjust for camera if needed

    // Draw boss as colored rectangle (placeholder)
    rl.DrawRectangle(
        int32(screenX),
        int32(screenY),
        int32(aabb.Width),
        int32(aabb.Height),
        rl.Red,
    )

    // Draw boss outline
    rl.DrawRectangleLines(
        int32(screenX),
        int32(screenY),
        int32(aabb.Width),
        int32(aabb.Height),
        rl.DarkRed,
    )
}
```

### Step 3: Render Boss HP Bar

**Location:** Top center of screen
**Design:** Horizontal bar showing HP percentage

```go
func (r *RaylibRenderer) renderBossHPBar(boss bosses.Boss) {
    if boss == nil || !boss.IsActive() {
        return
    }

    barWidth := float32(200)
    barHeight := float32(20)
    barX := float32(screenWidth)/2 - barWidth/2  // Center horizontally
    barY := float32(20)                             // Top of screen

    // Background (empty bar)
    rl.DrawRectangle(
        int32(barX),
        int32(barY),
        int32(barWidth),
        int32(barHeight),
        rl.DarkGray,
    )

    // Health fill
    healthPercent := boss.GetHP() / boss.GetMaxHP()
    fillWidth := barWidth * healthPercent

    color := rl.Green
    if healthPercent < 0.33 {
        color = rl.Red
    } else if healthPercent < 0.66 {
        color = rl.Orange
    }

    rl.DrawRectangle(
        int32(barX),
        int32(barY),
        int32(fillWidth),
        int32(barHeight),
        color,
    )

    // Border
    rl.DrawRectangleLines(
        int32(barX),
        int32(barY),
        int32(barWidth),
        int32(barHeight),
        rl.Black,
    )

    // Text: "HP: X / Y"
    hpText := fmt.Sprintf("HP: %.0f / %.0f", boss.GetHP(), boss.GetMaxHP())
    rl.DrawText(
        hpText,
        int32(barX + 5),
        int32(barY + 2),
        12,
        rl.White,
    )
}
```

### Step 4: Render Projectiles

**Location:** World space (wherever projectiles exist)

```go
func (r *RaylibRenderer) renderProjectiles(boss bosses.Boss) {
    if boss == nil {
        return
    }

    projectiles := boss.GetProjectiles()
    for _, proj := range projectiles {
        if !proj.Active {
            continue
        }

        // Draw projectile as small circle or rectangle
        centerX := proj.AABB.X + proj.AABB.Width/2
        centerY := proj.AABB.Y + proj.AABB.Height/2

        rl.DrawCircle(
            int32(centerX),
            int32(centerY),
            proj.AABB.Width/2,
            rl.Yellow,
        )
    }
}
```

### Step 5: Render Floor Tiles

**File:** Modify existing tile rendering in `internal/adapters/rendering/raylib.go`
**Logic:** Check tile type and render accordingly

```go
// In existing tile rendering loop:
if tile.Type == entities.TileTypeFloor {
    // Render floor based on floor type from game config
    bossRoomCfg := game.GetConfig().Level.BossRoom
    if bossRoomCfg != nil {
        if bossRoomCfg.FloorType == config.FloorConcrete {
            rl.DrawRectangle(/* concrete gray color */)
        } else if bossRoomCfg.FloorType == config.FloorLava {
            rl.DrawRectangle(/* lava orange/red color */)
            // Optional: add animation (pulsing or flickering)
        }
    }
}
```

**Color Suggestions:**
- Concrete: `rl.DarkGray` or RGB(100, 100, 100)
- Lava: `rl.Red` or `rl.Orange` with optional pulsing animation

### Step 6: Render Victory/Defeat Screen

**Trigger:** Check `game.GetGameState()`

```go
func (r *RaylibRenderer) renderGameStateOverlay(game *engine.Game) {
    state := game.GetGameState()

    switch state {
    case entities.GameStateVictory:
        r.renderVictoryScreen()
    case entities.GameStateDefeat:
        r.renderDefeatScreen()
    }
}

func (r *RaylibRenderer) renderVictoryScreen() {
    // Semi-transparent overlay
    rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.NewColor(0, 0, 0, 200))

    // Victory text
    rl.DrawText(
        "VICTORY!",
        screenWidth/2 - 80,
        screenHeight/2 - 40,
        40,
        rl.Green,
    )

    rl.DrawText(
        "Press ESC to return to main menu",
        screenWidth/2 - 150,
        screenHeight/2 + 40,
        20,
        rl.White,
    )
}

func (r *RaylibRenderer) renderDefeatScreen() {
    // Semi-transparent overlay
    rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.NewColor(0, 0, 0, 200))

    // Defeat text
    rl.DrawText(
        "DEFEATED!",
        screenWidth/2 - 80,
        screenHeight/2 - 40,
        40,
        rl.Red,
    )

    rl.DrawText(
        "Press ESC to try again",
        screenWidth/2 - 120,
        screenHeight/2 + 40,
        20,
        rl.White,
    )
}
```

### Step 7: Integration Points

**In `Render()` method of RaylibRenderer:**

```go
func (r *RaylibRenderer) Render(game *engine.Game, inputState input.InputState) {
    // ... existing rendering code ...

    // Boss and boss fight rendering (new)
    if game.GetBoss() != nil {
        r.renderBoss(game.GetBoss())
        r.renderBossHPBar(game.GetBoss())
        r.renderProjectiles(game.GetBoss())
    }

    // Game state overlay (new)
    r.renderGameStateOverlay(game)

    rl.EndDrawing()
}
```

**Important:** Add these calls AFTER main game rendering but BEFORE `rl.EndDrawing()`.

## Testing Checklist

- [ ] Boss renders as colored rectangle in boss room
- [ ] HP bar appears at top when boss is active
- [ ] HP bar updates when boss takes damage (bombs)
- [ ] HP bar disappears when boss is not active (player outside room)
- [ ] Projectiles render (if any boss shoots them)
- [ ] Floor tiles render correctly (concrete=gray, lava=orange)
- [ ] Victory screen appears on boss defeat
- [ ] Defeat screen appears on player defeat
- [ ] Game state overlay doesn't block gameplay when not needed

## Styling Notes

### Boss Rectangle
- Current: Red with dark red outline
- Alternative: Use sprite if available
- Size: Matches AABB from test_boss (200x200px)

### HP Bar
- Position: Top center of screen (fixed screen coordinates)
- Width: 200px
- Height: 20px
- Color gradient:
  - Green (> 66% HP)
  - Orange (33%-66% HP)
  - Red (< 33% HP)

### Floor Tiles
- Concrete: Solid gray (no animation)
- Lava: Orange/red with optional pulsing effect
  - Optional: `sin(time) * 0.5 + 0.5` for brightness variation

### Screens
- Semi-transparent black overlay: RGBA(0, 0, 0, 200)
- Victory: Green text "VICTORY!"
- Defeat: Red text "DEFEATED!"
- Instructions: White text

## Optional Enhancements (Future)

1. **Particle Effects:** Explosions on bomb hits
2. **Boss Sprite Animation:** Replace rectangle with animated sprite
3. **Projectile Trails:** Trails following projectiles
4. **Screen Shake:** Shake on impact/explosion
5. **Sound Effects:** Boss sounds, impact sounds (requires audio system)
6. **Lava Animation:** Pulsing or flowing lava effect

## Important Notes

1. **Imports Required:**
   - Boss and system types already available in game engine
   - May need: `"fmt"` for text formatting

2. **Screen Coordinates:**
   - HP bar uses screen space (fixed position)
   - Boss/projectiles use world space (need camera adjustment)

3. **Game State Flow:**
   - Victory/defeat screens should persist until user action
   - Consider pausing normal game rendering when showing overlay
   - May need to add input handling for ESC to exit screen

4. **Performance:**
   - Boss rendering is minimal (one rectangle + bar)
   - Projectiles are dynamic but typically few in count
   - Floor rendering uses existing tile loop

## Commit Message

```
Feature: Add boss fight rendering (HP bar, boss sprite, projectiles, floor)

- Render boss as colored rectangle with outline
- Display HP bar at top of screen with health percentage
- Render active projectiles in world space
- Style floor tiles (concrete gray, lava orange/red)
- Add victory and defeat screen overlays
- Integrate all rendering into main game loop

All rendering in adapter layer; domain layer unchanged.
```

## Resume Instructions

1. Open `internal/adapters/rendering/raylib.go`
2. Locate the main `Render()` function
3. Add rendering calls before `rl.EndDrawing()`
4. Implement helper methods for each visual element
5. Test with `go run cmd/game/main.go` (level -2)
6. Check main.go line 29 - change back to -1 if this was test level

See COMPLETION_STATUS.md for file locations and current state.
