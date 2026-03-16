#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/../../.." && pwd))"

lookup_root=""
format="text"

fail() {
  echo "FAIL: $1" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --root)
      [[ $# -ge 2 ]] || fail "--root requires a value"
      lookup_root="$2"
      shift 2
      ;;
    --format)
      [[ $# -ge 2 ]] || fail "--format requires a value"
      format="$2"
      shift 2
      ;;
    *)
      fail "unknown argument: $1"
      ;;
  esac
done

if [[ -z "$lookup_root" ]]; then
  lookup_root="$root_dir"
fi

discover_gitnexus_root() {
  local start="$1"
  local current="$start"
  while [[ -n "$current" ]]; do
    local candidate="$current/repos/untrusted/GitNexus/gitnexus"
    if [[ -d "$candidate" ]]; then
      echo "$candidate"
      return 0
    fi
    local parent
    parent="$(dirname "$current")"
    if [[ "$parent" == "$current" ]]; then
      break
    fi
    current="$parent"
  done
  return 1
}

gitnexus_bin="${CAIRN_GITNEXUS_BIN:-}"
gitnexus_root="${CAIRN_GITNEXUS_ROOT:-}"

if [[ -n "$gitnexus_bin" ]]; then
  if [[ ! -x "$gitnexus_bin" ]]; then
    fail "CAIRN_GITNEXUS_BIN is set but not executable: $gitnexus_bin"
  fi
  if [[ "$format" == "json" ]]; then
    printf '{\n  "status": "ready",\n  "mode": "binary",\n  "gitnexus_bin": "%s"\n}\n' "$gitnexus_bin"
  else
    echo "PASS: GitNexus ready via CAIRN_GITNEXUS_BIN"
    echo "CAIRN_GITNEXUS_BIN=$gitnexus_bin"
  fi
  exit 0
fi

if command -v gitnexus >/dev/null 2>&1; then
  resolved_bin="$(command -v gitnexus)"
  if [[ "$format" == "json" ]]; then
    printf '{\n  "status": "ready",\n  "mode": "binary",\n  "gitnexus_bin": "%s"\n}\n' "$resolved_bin"
  else
    echo "PASS: GitNexus ready via PATH"
    echo "CAIRN_GITNEXUS_BIN=$resolved_bin"
  fi
  exit 0
fi

if [[ -z "$gitnexus_root" ]]; then
  gitnexus_root="$(discover_gitnexus_root "$lookup_root" || true)"
fi

if [[ -z "$gitnexus_root" || ! -d "$gitnexus_root" ]]; then
  fail "GitNexus checkout not found at $gitnexus_root"
fi

if ! command -v node >/dev/null 2>&1; then
  fail "node is required to run a built local GitNexus checkout"
fi

entry_point="$gitnexus_root/dist/cli/index.js"
if [[ ! -f "$entry_point" ]]; then
  fail "built GitNexus entry point missing: $entry_point"
fi

node_bin="$(command -v node)"
if [[ "$format" == "json" ]]; then
  printf '{\n  "status": "ready",\n  "mode": "node",\n  "gitnexus_root": "%s",\n  "node_bin": "%s",\n  "entry_point": "%s"\n}\n' "$gitnexus_root" "$node_bin" "$entry_point"
else
  echo "PASS: GitNexus ready via built local checkout"
  echo "CAIRN_GITNEXUS_ROOT=$gitnexus_root"
  echo "NODE_BIN=$node_bin"
  echo "ENTRY_POINT=$entry_point"
fi
