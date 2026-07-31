#!/usr/bin/env bash
# Hexagonal architecture boundary audit.
#
# The domain layer must stay framework-free and must not depend on anything
# outward (adapters, cmd). Exits non-zero on any violation.
set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT" || exit 1

DOMAIN="internal/domain"
failures=0

check() {
	local label="$1" pattern="$2" hint="$3"
	local hits

	hits="$(grep -rnE "$pattern" "$DOMAIN" --include='*.go' 2>/dev/null)"

	if [ -n "$hits" ]; then
		printf 'FAIL  %s\n' "$label"
		printf '%s\n' "$hits" | sed 's/^/      /'
		printf '      -> %s\n\n' "$hint"
		failures=$((failures + 1))
	else
		printf 'ok    %s\n' "$label"
	fi
}

check "domain imports no Raylib" \
	'raylib-go' \
	'Move the Raylib call into internal/adapters/ and pass plain data across.'

check "domain imports no adapters" \
	'drill-game/internal/adapters' \
	'Wrong dependency direction. Define the interface in domain, implement it in adapters.'

check "domain imports no cmd" \
	'drill-game/cmd' \
	'Wiring belongs in cmd/game/main.go, not the other way round.'

check "domain uses no Raylib types" \
	'\brl\.[A-Z]' \
	'Use the domain types instead: types.Vec2, types.AABB.'

printf -- '---\n'

if [ "$failures" -eq 0 ]; then
	printf 'Architecture check passed. Domain layer is clean.\n'
	exit 0
fi

printf 'Architecture check failed: %d violation(s).\n' "$failures"
exit 1
