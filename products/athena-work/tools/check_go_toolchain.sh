#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
go_mods=(
  "$root_dir/products/athena-mind/go.mod"
  "$root_dir/products/athena-use/go.mod"
)

fail() {
  echo "FAIL: $1" >&2
  exit 1
}

extract_major_minor() {
  local version="$1"
  local major="${version%%.*}"
  local rest="${version#*.}"
  local minor="${rest%%.*}"
  echo "${major:-0} ${minor:-0}"
}

required_version=""
for go_mod in "${go_mods[@]}"; do
  if [[ ! -f "$go_mod" ]]; then
    fail "go.mod not found at $go_mod"
  fi
  parsed_version="$(awk '/^go[[:space:]]+[0-9]+\.[0-9]+/{print $2; exit}' "$go_mod")"
  if [[ -z "$parsed_version" ]]; then
    fail "could not parse required Go version from $go_mod"
  fi
  if [[ -z "$required_version" || "$parsed_version" > "$required_version" ]]; then
    required_version="$parsed_version"
  fi
done

if ! command -v go >/dev/null 2>&1; then
  cat >&2 <<EOF
FAIL: Go toolchain not installed.
Required: Go >= $required_version (from products/athena-mind/go.mod and products/athena-use/go.mod)
Install: brew install go
Then verify with: go version
EOF
  exit 1
fi

installed_raw="$(go version | awk '{print $3}')"
installed_version="${installed_raw#go}"

read -r req_major req_minor <<<"$(extract_major_minor "$required_version")"
read -r inst_major inst_minor <<<"$(extract_major_minor "$installed_version")"

if (( inst_major < req_major )) || (( inst_major == req_major && inst_minor < req_minor )); then
  fail "Go $installed_version is below required $required_version from product go.mod files"
fi

echo "PASS: Go toolchain available ($installed_raw), requirement satisfied (>= $required_version)"
