# Installation

## Summary

Install runtime prerequisites and verify the slim AthenaMind distribution is healthy.

## Prerequisites

- macOS or Linux shell.
- Go 1.22+.
- Network access if pulling new Go modules.

## Install Steps

1. Choose install method:
   - Precompiled binary guidance: `docs/operator/athena-mind/getting-started/binaries.md`
   - Build from source: continue below
2. Clone and enter repo.
3. Bootstrap the repo-local runtime and verify the toolchain:
```bash
./tools/dev/bootstrap_platform.sh
```
4. Verify Go:
```bash
go version
```
5. Download modules:
```bash
cd products/athena-mind
go mod download
```
6. Run full tests:
```bash
cd products/athena-mind
go test ./...
```

## Optional Services

### Ollama embeddings
```bash
export ATHENA_EMBEDDING_ENDPOINT="http://127.0.0.1:11434"
export ATHENA_OLLAMA_EMBED_MODEL="mxbai-embed-large:latest"
```

Run the platform smoke path once the endpoint is available:

```bash
./tools/platform/smoke_v1.sh
```

### Optional Mongo-backed persistence
```bash
export ATHENA_INDEX_BACKEND="mongodb"
export ATHENA_MONGODB_URI="mongodb://127.0.0.1:27017"
export ATHENA_MONGODB_DATABASE="athenamind"
```

Validate the optional Mongo path:

```bash
./tools/dev/check_mongodb_local.sh
./tools/platform/smoke_mongodb.sh
```

## Memory Storage Defaults

- Default runtime boundary: repo-local `.athena/`
- Recommended root for examples: `$ATHENA_HOME/memory/default`
- Recommended for real use: set explicit roots through `ATHENA_HOME` or `ATHENA_MEMORY_ROOT`

## Observability Defaults

- Local telemetry events: `<root>/telemetry/events.jsonl`
- Retrieval metrics: `<root>/telemetry/retrieval-metrics.jsonl`
- OTel tracing enabled in CLI runtime, with OTLP/collector config via env vars.

## Verify Runtime

```bash
export ATHENA_HOME="$PWD/.athena"
export ATHENA_MEMORY_ROOT="$ATHENA_HOME/memory/default"
mkdir -p "$ATHENA_MEMORY_ROOT"
(cd products/athena-mind && go run ./cmd/memory-cli retrieve --root "$ATHENA_MEMORY_ROOT" --query "health check")
```

If index is empty, you should see a clear error indicating no entries exist yet.
