#!/usr/bin/env bash
# Go style audit.
#
# Enforces the project's explicit-false-boolean rule: this codebase never uses
# the `!` negation operator. Write `x == false`, not `!x`.
#
#   ./scripts/style-check.sh                 # sweep internal/ and cmd/
#   ./scripts/style-check.sh path/to/file.go # check specific files
#
# Exits non-zero on any violation.
set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT" || exit 1

# Strip anything where a `!` is not an operator, so comments and string
# literals never trigger a false positive. Order matters: strings go first, so
# that a `//` inside a literal (e.g. a URL) does not start a comment.
#
# Not handled: block comments spanning multiple lines. This codebase does not
# use them inside function bodies, which is the only place `!` can appear.
strip_noise() {
	sed -e 's/`[^`]*`/``/g' \
		-e 's/"\([^"\\]\|\\.\)*"/""/g' \
		-e "s/'\([^'\\]\|\\\\.\)*'/''/g" \
		-e 's:/\*.*\*/::g' \
		-e 's://.*::'
}

# Matches a `!` used as negation: not part of `!=`, and not already stripped.
NEGATION='![^=]'

collect() {
	local file="$1"
	strip_noise <"$file" | grep -nE "$NEGATION" | sed "s|^|${file}:|"
}

targets=()
if [ "$#" -gt 0 ]; then
	targets=("$@")
else
	while IFS= read -r f; do
		targets+=("$f")
	done < <(find internal cmd -name '*.go' -type f 2>/dev/null | sort)
fi

violations=""
for file in "${targets[@]}"; do
	case "$file" in
	*.go) ;;
	*) continue ;;
	esac
	[ -f "$file" ] || continue

	hits="$(collect "$file")"
	if [ -n "$hits" ]; then
		violations+="$hits"$'\n'
	fi
done

violations="$(printf '%s' "$violations" | sed '/^$/d')"

if [ -z "$violations" ]; then
	printf 'ok    no `!` negation (%d file(s) checked)\n' "${#targets[@]}"
	exit 0
fi

printf 'FAIL  explicit-false-boolean rule\n'
printf '%s\n' "$violations" | sed 's/^/      /'
printf '      -> This codebase does not use `!`. Write `x == false` instead of `!x`.\n'
printf '         See .claude/rules/go-style.md\n'
exit 1
