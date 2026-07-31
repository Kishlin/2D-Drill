---
name: architecture-check
description: Verify hexagonal architecture compliance across internal/domain/ — checks that the domain layer imports no Raylib, no adapters, no cmd, and uses no Raylib types. Use before finishing work that touched internal/domain/, when reviewing a branch, or when the user asks about architecture or layering violations.
---

# Architecture Check

Runs the four-part hexagonal boundary audit over `internal/domain/`.

```bash
./scripts/architecture-check.sh
```

Runs on the host, needs nothing but `grep`. Exits non-zero on any violation.

## What it checks

| # | Rule | Why |
|---|------|-----|
| 1 | No `raylib-go` import | The domain layer must compile and test without a graphics stack |
| 2 | No `internal/adapters` import | Dependency direction — adapters depend on domain, never the reverse |
| 3 | No `cmd/` import | Wiring flows inward from `main.go` only |
| 4 | No `rl.*` types | Domain speaks `types.Vec2` / `types.AABB`, not `rl.Vector2` |

`_test.go` files are **included**. Domain tests are pure Go here — there is no
legitimate reason for one to reach for Raylib or an adapter. If a test needs a
fake, define the interface in the domain package and write the fake there.

## Reading a failure

Output names the offending `file:line` and its import. The fix is essentially
never to loosen the check:

- **Domain needs to draw something** → it does not. Return data (positions,
  colors as domain values, box sets) and let `internal/adapters/rendering/`
  translate it. See how `BossRenderer` consumes `bosses.Boss`.
- **Domain needs input** → it consumes `input.InputState`, a plain domain
  struct that `internal/adapters/input/` populates from Raylib.
- **Domain needs a Raylib vector** → use `types.Vec2`. The adapter converts at
  the boundary.

## Relationship to the automatic hook

A `PostToolUse` hook runs this script whenever a file under `internal/domain/`
is written or edited, and blocks the edit on failure. Invoke this skill for a
**full sweep** — before a commit, after a refactor that moved packages, or when
auditing work you did not do yourself.
