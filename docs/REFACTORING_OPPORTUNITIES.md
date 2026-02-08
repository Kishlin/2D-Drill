# Refactoring Opportunities

**Date:** January 2026
**Last updated:** February 2026

This document catalogs refactoring opportunities found during a codebase audit. Items 1-5 have been completed. Items 6-8 remain open.

---

## Completed Refactors

### 1. Grid Navigation Duplication — Done

**Files changed:**
- `internal/domain/ui/grid_navigator.go` (new) — Shared `GridNavigator` struct with `NavigateUp/Down/Left/Right`, `GetSelectedRow/Col`
- `internal/domain/ui/state.go` — `UpgradeShopState` and `ItemShopState` embed `GridNavigator`

Extracted a reusable `GridNavigator` type that handles grid-based selection with customizable dimensions. Both shop states embed it. `ItemShopState` overrides navigation methods to skip empty cells.

---

### 2. Modal UI Pattern Duplication — Done

**Files changed:**
- `internal/domain/ui/modal_service.go` (new) — `ModalServiceProvider` interface + shared `processModalService` function
- `internal/domain/ui/state.go` — Unified `HospitalState`/`FuelStationState` into `ModalServiceState`
- `internal/domain/ui/hospital.go` — Implements `ModalServiceProvider`, delegates `Process()` to shared function
- `internal/domain/ui/fuel_station.go` — Same pattern as hospital

Created a `ModalServiceProvider` interface (`GetAmount`, `GetCost`, `BuildEffect`) and a shared `processModalService` function that handles the common control flow (close check, first frame skip, navigation, interaction).

---

### 3. Panic-Based Error Handling — Done

**Files changed:**
- `internal/domain/engine/game.go` — `NewGame` returns `(*Game, error)` instead of panicking on boss creation failure
- `internal/domain/systems/drilling.go` — Replaced panics with fallthrough to normal duration calculation when hazard config is missing
- `cmd/game/main.go` — Handles `NewGame` error return

`registry.go` panic left as-is — it runs at init time and is idiomatic Go for invalid program configuration.

---

### 4. Type Assertion Boilerplate — Done

**Files changed:**
- `internal/domain/ui/ui.go` — Added `ResetState()` to the `UI` interface
- `internal/domain/engine/game.go` — Collapsed 5-case type switch in `resetUIState` to a single `registeredUI.ResetState()` call

---

### 5. Excessive Getter Methods + Game Struct Cleanup — Done

**Files changed:**
- `internal/domain/engine/game.go` — Removed 12 getter methods, made externally-accessed fields public, kept internal fields private
- `internal/domain/systems/projectile_system.go` — Converted from struct to package-level functions (`SpawnProjectiles`, `UpdateProjectiles`)
- `internal/adapters/rendering/raylib.go` — Updated all `game.GetX()` calls to direct field access
- `cmd/game/main.go` — Updated `game.GetPlayer()` to `game.Player`

Changes made:
- **Public fields** (used by renderer/main): `World`, `Player`, `Buildings`, `Boss`, `GameState`, `UIManager`, `InventoryUI`, `Projectiles`
- **Private fields** (internal only): `drillingSystem`, `bossFightSystem`, `projectilePool`, `projectileBounds`, `effectProcessor`, `effectContext`, `config`
- **Removed unused fields**: `UpgradeCatalog`, `ItemCatalog`, `damageables`
- **Projectile system**: Removed `ProjectileSystem` struct. Game owns the pool and builds a `[]types.AABB` snapshot (`Projectiles`) for the renderer each frame.
- **Struct field ordering**: Grouped by concern (core state, UI, render data, systems, projectile internals, effects, config)

---

## Remaining Opportunities

### 6. First Frame Pattern

**Priority:** Low
**Effort:** Small
**Files:**
- `internal/domain/ui/market.go`
- `internal/domain/ui/modal_service.go`

**Issue:** Multiple UI types use the same "skip first frame" pattern with `IsFirstFrame()` and `ClearFirstFrame()` methods.

**Suggested Fix:** Extract a `FirstFrameSkipper` helper if this pattern spreads to more UI types. Currently manageable.

---

### 7. Optional Handler Nil Checks

**Priority:** Low
**Effort:** Small
**File:** `internal/domain/bosses/base_boss.go`

**Issue:** `PhaseChangeHandler` and `DamageReactionHandler` are optional (can be nil), requiring nil checks at call sites.

**Suggested Fix:** Use the null-object pattern with empty no-op implementations as defaults:
```go
type noOpPhaseHandler struct{}
func (noOpPhaseHandler) OnPhaseChange(int, phases.Config) {}
```

---

### 8. Config Validation Timing

**Priority:** Low
**Effort:** Small
**Files:** `internal/domain/systems/drilling.go`, `internal/domain/config/game_config.go`

**Issue:** A `Validate()` method exists on `GameConfig` and is called at startup. The defensive panics in `drilling.go` were replaced with safe fallthrough behavior (item 3), but the code still does redundant runtime checks for config that was already validated at load time.

**Remaining Fix:** Remove the redundant runtime checks in `drilling.go`, trusting the upfront validation.

---

## Summary Table

| # | Issue | Priority | Effort | Status |
|---|-------|----------|--------|--------|
| 1 | Grid navigation duplication | High | Medium | Done |
| 2 | Modal UI duplication | High | Medium | Done |
| 3 | Panic error handling | High | Medium | Done |
| 4 | Type assertion boilerplate | Medium | Small | Done |
| 5 | Getter methods + Game struct cleanup | Medium | Small | Done |
| 6 | First frame pattern | Low | Small | Open |
| 7 | Optional handler nil checks | Low | Small | Open |
| 8 | Config validation timing | Low | Small | Open |
