# Installation

## Summary

Install runtime prerequisites and verify the slim memory-cli distribution is healthy.

## Prerequisites

- macOS or Linux shell.
- Go 1.22+.
- Network access if pulling new Go modules.

## Install Steps

1. Choose install method:
   - Precompiled binary guidance: `docs/operator/memory-cli/getting-started/binaries.md`
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
cd products/memory-cli
go mod download
```
6. Run full tests:
```bash
cd products/memory-cli
go test ./...
```

## Optional Services

### Ollama embeddings
```bash
export CAIRN_EMBEDDING_ENDPOINT="http://127.0.0.1:11434"
export CAIRN_OLLAMA_EMBED_MODEL="mxbai-embed-large:latest"
```

Run the platform smoke path once the endpoint is available:

```bash
./tools/platform/smoke_v1.sh
```

## Memory Storage Defaults

- Default runtime boundary: repo-local `.cairn/`
- Recommended root for examples: `$CAIRN_HOME/memory/default`
- Recommended for real use: set explicit roots through `CAIRN_HOME` or `CAIRN_MEMORY_ROOT`

## Observability Defaults

- Local telemetry events: `<root>/telemetry/events.jsonl`
- Retrieval metrics: `<root>/telemetry/retrieval-metrics.jsonl`
- OTel tracing enabled in CLI runtime, with OTLP/collector config via env vars.

## Storage Note

The Cairn memory-cli product is sqlite-only.
If you want backend experiments beyond sqlite, use the separate research repo rather than this stripped personal product.

## Verify Runtime

```bash
export CAIRN_HOME="$PWD/.cairn"
export CAIRN_MEMORY_ROOT="$CAIRN_HOME/memory/default"
mkdir -p "$CAIRN_MEMORY_ROOT"
(cd products/memory-cli && go run ./cmd/memory-cli retrieve --root "$CAIRN_MEMORY_ROOT" --query "health check")
```

If index is empty, you should see a clear error indicating no entries exist yet.
