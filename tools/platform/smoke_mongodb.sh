#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/../.." && pwd))"
athena_home="${ATHENA_HOME:-$root_dir/.athena}"

memory_root="${ATHENA_SMOKE_MEMORY_ROOT:-${ATHENA_MEMORY_ROOT:-$athena_home/memory/mongodb-smoke}}"
mkdir -p "$memory_root"

mongodb_uri="${ATHENA_MONGODB_URI:-mongodb://127.0.0.1:27017}"
mongodb_database="${ATHENA_MONGODB_DATABASE:-athenamind}"
embedding_endpoint="${ATHENA_EMBEDDING_ENDPOINT:-}"
embedding_model="${ATHENA_OLLAMA_EMBED_MODEL:-}"
latency_threshold="${MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS:-}"

run_memory_cli() {
  (
    cd "$root_dir/products/athena-mind"
    export ATHENA_INDEX_BACKEND="mongodb"
    export ATHENA_MONGODB_URI="$mongodb_uri"
    export ATHENA_MONGODB_DATABASE="$mongodb_database"
    if [[ -n "$embedding_model" ]]; then
      export ATHENA_OLLAMA_EMBED_MODEL="$embedding_model"
    fi
    go run ./cmd/memory-cli "$@"
  )
}

echo "== AthenaMind tests =="
(
  cd "$root_dir/products/athena-mind"
  go test ./...
)

echo "== MongoDB readiness =="
run_memory_cli verify mongodb >/dev/null

echo "== AthenaMind Mongo write =="
write_args=(write
  --root "$memory_root"
  --id mongo-smoke-bootstrap
  --title "Mongo Smoke Bootstrap"
  --type prompt
  --domain platform
  --body "Persist AthenaMind index state in MongoDB while keeping markdown files local."
  --stage planning
  --reviewer smoke
  --decision approved
  --reason "mongodb smoke bootstrap"
  --risk "low"
  --notes "smoke")
if [[ -n "$embedding_endpoint" ]]; then
  write_args+=(--embedding-endpoint "$embedding_endpoint")
fi
run_memory_cli "${write_args[@]}" >/dev/null

echo "== AthenaMind Mongo retrieve =="
retrieve_args=(retrieve
  --root "$memory_root"
  --query "persist athenamind index state in mongodb"
  --domain platform)
if [[ -n "$embedding_endpoint" ]]; then
  retrieve_args+=(--embedding-endpoint "$embedding_endpoint")
fi
run_memory_cli "${retrieve_args[@]}" >/dev/null

if [[ -n "$embedding_endpoint" ]]; then
  echo "== AthenaMind Mongo semantic health =="
  (
    cd "$root_dir/products/athena-mind"
    export ATHENA_INDEX_BACKEND="mongodb"
    export ATHENA_MONGODB_URI="$mongodb_uri"
    export ATHENA_MONGODB_DATABASE="$mongodb_database"
    export MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS="${latency_threshold:-0}"
    if [[ -n "$embedding_model" ]]; then
      export ATHENA_OLLAMA_EMBED_MODEL="$embedding_model"
    fi
    go run ./cmd/memory-cli verify health \
      --root "$memory_root" \
      --query "persist athenamind index state in mongodb" \
      --domain platform \
      --session-id mongo-smoke \
      --embedding-endpoint "$embedding_endpoint" >/dev/null
  )
fi

echo "smoke_mongodb: PASS"
echo "athena_home: $athena_home"
echo "memory_root: $memory_root"
echo "mongodb_uri: $mongodb_uri"
echo "mongodb_database: $mongodb_database"
if [[ -n "$embedding_endpoint" ]]; then
  echo "embedding_endpoint: $embedding_endpoint"
fi
if [[ -n "$embedding_model" ]]; then
  echo "embedding_model: $embedding_model"
fi
if [[ -n "$embedding_endpoint" ]]; then
  echo "latency_threshold_ms: ${latency_threshold:-0}"
fi
