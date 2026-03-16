#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/../.." && pwd))"
use_product_root="$root_dir/products/tool-cli"

run_validate() {
  if [[ -n "${USE_CLI_BIN:-}" ]]; then
    "$USE_CLI_BIN" validate "$@"
    return
  fi

  if [[ -f "$use_product_root/go.mod" ]] && [[ -f "$use_product_root/cmd/tool-cli/main.go" ]]; then
    (
      cd "$use_product_root"
      go run ./cmd/tool-cli validate "$@"
    )
    return
  fi

  if command -v tool-cli >/dev/null 2>&1; then
    tool-cli validate "$@"
    return
  fi

  echo "FAIL: tool-cli validate command unavailable" >&2
  exit 1
}

args=()
if [[ -n "${CAIRN_TOOL_REGISTRY:-}" ]]; then
  args+=(--registry "$CAIRN_TOOL_REGISTRY")
fi
if [[ "${CAIRN_TOOL_INCLUDE_LOCAL:-false}" == "true" ]]; then
  args+=(--include-local)
fi

if [[ ${#args[@]} -gt 0 ]]; then
  run_validate "${args[@]}"
else
  run_validate
fi
