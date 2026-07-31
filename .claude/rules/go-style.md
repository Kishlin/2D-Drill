# Go Style Rules

Non-negotiable. Enforced by `scripts/style-check.sh` and a blocking
`PostToolUse` hook — a write that violates these is rejected, not warned about.

## No `!` negation

This codebase does not use the `!` operator. Compare against `false` explicitly.

```go
// NO
if !ok { ... }
if err != nil && !player.OnGround { ... }
wasAirborne := !player.OnGround
for !done { ... }

// YES
if ok == false { ... }
if err != nil && player.OnGround == false { ... }
wasAirborne := player.OnGround == false
for done == false { ... }
```

Applies everywhere `!` would appear as negation — `if`, `for`, assignments,
and `&&` / `||` operands. `!=` is unaffected.

The rule exists because `!` is a single character that inverts the meaning of a
condition and is easy to miss when skimming. `== false` cannot be misread.
