#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/../.." && pwd))"
cairn_home="${CAIRN_HOME:-$root_dir/.cairn}"

memory_root="${CAIRN_SMOKE_MEMORY_ROOT:-${CAIRN_MEMORY_ROOT:-$cairn_home/memory/smoke-v1}}"
mkdir -p "$memory_root"
embedding_endpoint="${CAIRN_EMBEDDING_ENDPOINT:-}"
embedding_model="${CAIRN_OLLAMA_EMBED_MODEL:-}"
latency_threshold="${MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS:-}"
runs_root="${CAIRN_RUNS_ROOT:-$cairn_home/runs}"
mkdir -p "$runs_root"

current_branch="$(git -C "$root_dir" branch --show-current)"
if [[ -z "$current_branch" ]]; then
  current_branch="master"
fi

echo "== tool-cli registry validation =="
"$root_dir/tools/platform/validate_tool_registry.sh" >/dev/null

echo "== memory-cli tests =="
(
  cd "$root_dir/products/memory-cli"
  go test ./...
)

echo "== memory-cli write =="
(
  cd "$root_dir/products/memory-cli"
  if [[ -n "$embedding_model" ]]; then
    export CAIRN_OLLAMA_EMBED_MODEL="$embedding_model"
  fi
  args=(write
    --root "$memory_root"
    --id smoke-bootstrap
    --title "Smoke Bootstrap"
    --type prompt
    --domain platform
    --body "Capture observer output before closing a cycle."
    --stage planning
    --reviewer smoke
    --decision approved
    --reason "smoke test bootstrap"
    --risk "low"
    --notes "smoke")
  if [[ -n "$embedding_endpoint" ]]; then
    args+=(--embedding-endpoint "$embedding_endpoint")
  fi
  go run ./cmd/memory-cli "${args[@]}"
)

echo "== memory-cli retrieve =="
(
  cd "$root_dir/products/memory-cli"
  if [[ -n "$embedding_model" ]]; then
    export CAIRN_OLLAMA_EMBED_MODEL="$embedding_model"
  fi
  args=(--root "$memory_root" --query "observer output before closing a cycle" --domain platform)
  if [[ -n "$embedding_endpoint" ]]; then
    args+=(--embedding-endpoint "$embedding_endpoint")
  fi
  go run ./cmd/memory-cli retrieve \
    "${args[@]}" >/dev/null
)

if [[ -n "$embedding_endpoint" ]]; then
  echo "== memory-cli semantic health =="
  (
    cd "$root_dir/products/memory-cli"
    if [[ -n "$embedding_model" ]]; then
      export CAIRN_OLLAMA_EMBED_MODEL="$embedding_model"
    fi
    export MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS="${latency_threshold:-0}"
    go run ./cmd/memory-cli verify health \
      --root "$memory_root" \
      --query "observer output before closing a cycle" \
      --domain platform \
      --embedding-endpoint "$embedding_endpoint" >/dev/null
  )
fi

echo "== work harness launch engineering =="
CAIRN_REQUIRED_BRANCH="$current_branch" \
  "$root_dir/products/work-harness/tools/launch_stage.sh" engineering >/dev/null || true

echo "== work harness launch qa =="
CAIRN_REQUIRED_BRANCH="$current_branch" \
  "$root_dir/products/work-harness/tools/launch_stage.sh" qa >/dev/null

echo "== work harness observer =="
observer_out="${CAIRN_OBSERVER_OUTPUT:-$runs_root/observer-smoke-report.md}"
mkdir -p "$(dirname "$observer_out")"
CAIRN_REQUIRED_BRANCH="$current_branch" \
  "$root_dir/products/work-harness/tools/run_observer_cycle.sh" \
  --cycle-id smoke-v1 \
  --output "$observer_out" >/dev/null

echo "smoke_v1: PASS"
echo "cairn_home: $cairn_home"
echo "memory_root: $memory_root"
echo "observer_report: $observer_out"
if [[ -n "$embedding_endpoint" ]]; then
  echo "embedding_endpoint: $embedding_endpoint"
fi
if [[ -n "$embedding_model" ]]; then
  echo "embedding_model: $embedding_model"
fi
if [[ -n "$embedding_endpoint" ]]; then
  echo "latency_threshold_ms: ${latency_threshold:-0}"
fi
