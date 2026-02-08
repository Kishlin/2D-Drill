# Refactoring Opportunities

**Date:** January 2026
**Status:** Identified, not yet implemented

This document catalogs refactoring opportunities found during a codebase audit. Use this as a reference when looking for improvements to tackle.

---

## High-Impact Opportunities

### 1. Grid Navigation Duplication

**Priority:** High
**Effort:** Medium
**Files:** `internal/domain/ui/state.go` (lines 38-88, 118-192)

**Issue:** `UpgradeShopState` and `ItemShopState` have identical grid navigation logic (NavigateUp, NavigateDown, NavigateLeft, NavigateRight, GetSelectedRow, GetSelectedCol) that calculates row/col from a flat index using the same pattern repeatedly.

**Example of duplicated code:**
```go
row := s.SelectedTier / UpgradeShopGridCols
col := s.SelectedTier % UpgradeShopGridCols
// ... same pattern repeated ~8 times across two types
```

**Suggested Fix:** Extract a shared `GridNavigator` type or mixin that handles grid-based selection with customizable grid dimensions. Both states could embed or use this component.

---

### 2. Modal UI Pattern Duplication

**Priority:** High
**Effort:** Medium
**Files:**
- `internal/domain/ui/hospital.go` (lines 21-70)
- `internal/domain/ui/fuel_station.go` (lines 21-70)

**Issue:** `HospitalUI` and `FuelStationUI` have nearly identical `Process()` methods with same control flow:
1. Close shop check
2. First frame skip
3. Navigation handling
4. Interaction handling

The only differences are the specific healing/refueling amounts and cost calculations.

**Duplicated pattern:**
```go
// Both follow identical pattern:
if inputState.CloseShop { return Close() }
if u.state.IsFirstFrame() { /* skip */ }
if inputState.NavUp/Down { /* navigate */ }
if inputState.Interact { /* get amount, cost, apply effects */ }
```

**Suggested Fix:** Create a base `ModalServiceUI` type that handles the common control flow and lets subclasses define the service-specific logic (GetAmount, GetCost, GetOptionLabel).

---

### 3. Panic-Based Error Handling

**Priority:** High
**Effort:** Medium
**Files:**
- `internal/domain/engine/game.go` (line 91)
- `internal/domain/systems/drilling.go` (lines 297, 300)
- `internal/domain/bosses/registry.go`

**Issue:** Three panic calls in domain layer for error handling. Domain code should return errors instead, allowing graceful error handling at the boundary (adapters).

**Suggested Fix:** Convert panics to error returns. Callers at the adapter layer can then decide how to handle errors (log, show message, etc.).

---

## Medium-Impact Opportunities

### 4. Type Assertion Boilerplate

**Priority:** Medium
**Effort:** Small
**File:** `internal/domain/engine/game.go` (lines 239-260)

**Issue:** `resetUIState()` uses multiple type assertions with the same pattern repeated 5 times:
```go
case ComponentType:
    if uiCast, ok := registeredUI.(*UIType); ok {
        uiCast.ResetState()
    }
```

**Suggested Fix:** Create a `Resettable` interface:
```go
type Resettable interface {
    ResetState()
}
```
Then call directly without type switching.

---

### 5. Excessive Getter Methods

**Priority:** Medium
**Effort:** Small (but needs design decision)
**File:** `internal/domain/engine/game.go` (lines 263-309)

**Issue:** The `Game` struct has 12 simple getter methods (GetWorld, GetPlayer, GetBuildings, etc.) that just return fields. This is verbose boilerplate.

**Options:**
1. Make fields public (Go idiom for simple structs)
2. Keep getters if encapsulation is important
3. Reduce if some getters are unused

---

## Lower-Priority Opportunities

### 6. First Frame Pattern

**Priority:** Low
**Effort:** Small
**Files:**
- `internal/domain/ui/market.go` (lines 28-31)
- `internal/domain/ui/hospital.go` (lines 28-30)
- `internal/domain/ui/fuel_station.go` (lines 28-30)

**Issue:** Multiple UI types use the same "skip first frame" pattern with `IsFirstFrame()` and `ClearFirstFrame()` methods.

**Suggested Fix:** Extract a `FirstFrameSkipper` helper if this pattern spreads to more UI types. Currently manageable.

---

### 7. Optional Handler Nil Checks

**Priority:** Low
**Effort:** Small
**File:** `internal/domain/bosses/base_boss.go` (lines 41-43, 127-129, 148-152)

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

**Issue:** A `Validate()` method now exists on `GameConfig` (covering hazard configs, drilling params, etc.), but `drilling.go` (lines 297, 300) still has two defensive `panic()` calls for missing/invalid hazard configs at runtime.

**Remaining Fix:** Either remove the panics in `drilling.go` (trusting the upfront validation) or convert them to error returns.

---

## Summary Table

| # | Issue | Priority | Effort | Status |
|---|-------|----------|--------|--------|
| 1 | Grid navigation duplication | High | Medium | Open |
| 2 | Modal UI duplication | High | Medium | Open |
| 3 | Panic error handling | High | Medium | Open |
| 4 | Type assertion boilerplate | Medium | Small | Open |
| 5 | Excessive getter methods | Medium | Small | Open |
| 6 | First frame pattern | Low | Small | Open |
| 7 | Optional handler nil checks | Low | Small | Open |
| 8 | Config validation timing | Low | Small | Partial |

---

## Quick Wins (Small Effort)

If looking for fast improvements:
1. **Type assertion boilerplate** (#4) - Add Resettable interface

## High-Impact Refactors (Medium Effort)

For significant code quality improvements:
1. **Grid navigation** (#1) - Reduces ~100 lines of duplication
2. **Modal UI pattern** (#2) - Reduces ~50 lines of duplication
3. **Panic handling** (#3) - Improves error handling robustness
