# AthenaPlatform Quickstart

This quickstart shows the intended developer path through the unified platform.

## What You Need

- Go installed
- a local checkout of this repository
- optional: Podman and Ollama if you want local embedding service support

## Local Runtime Boundary

Use a repo-local `.athena/` directory for uncommitted Athena runtime state.

Suggested local folders:

```text
.athena/
  workspace/
  memory/
  artifacts/
  cache/
  runs/
  config/
```

This keeps product code and docs committed while local runtime state stays isolated and disposable.

Recommended shell setup:

```bash
export ATHENA_HOME="$PWD/.athena"
export ATHENA_MEMORY_ROOT="$ATHENA_HOME/memory/default"
mkdir -p "$ATHENA_MEMORY_ROOT"
```

## 1. Verify AthenaMind

From repo root:

```bash
cd products/athena-mind
go test ./...
cd ../..
```

## 2. Create One Local Memory Entry

```bash
(cd products/athena-mind && go run ./cmd/memory-cli write \
  --root "$ATHENA_MEMORY_ROOT" \
  --id quickstart-bootstrap \
  --title "Quickstart Bootstrap" \
  --type prompt \
  --domain platform \
  --body "Use AthenaWork queue metadata and capture observer output before closing a cycle." \
  --stage planning \
  --reviewer matt \
  --decision approved \
  --reason "platform quickstart" \
  --risk "low" \
  --notes "example bootstrap")
```

## 3. Retrieve It

```bash
(cd products/athena-mind && go run ./cmd/memory-cli retrieve \
  --root "$ATHENA_MEMORY_ROOT" \
  --query "observer output before closing a cycle" \
  --domain platform)
```

Optional embedding-backed validation:

```bash
(cd products/athena-mind && \
  MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS=0 \
  ATHENA_OLLAMA_EMBED_MODEL="mxbai-embed-large:latest" \
  go run ./cmd/memory-cli write \
    --root "$ATHENA_MEMORY_ROOT" \
    --id quickstart-embed \
    --title "Quickstart Embed" \
    --type prompt \
    --domain platform \
    --body "Use embedding-backed retrieval to improve local memory selection." \
    --stage planning \
    --reviewer matt \
    --decision approved \
    --reason "embedding quickstart" \
    --risk "low" \
    --notes "embedding example" \
    --embedding-endpoint http://192.168.1.35:11434)

(cd products/athena-mind && \
  MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS=0 \
  ATHENA_OLLAMA_EMBED_MODEL="mxbai-embed-large:latest" \
  go run ./cmd/memory-cli verify health \
    --root "$ATHENA_MEMORY_ROOT" \
    --query "improve local memory selection with embeddings" \
    --domain platform \
    --embedding-endpoint http://192.168.1.35:11434)
```

## 4. Inspect The Example Workspace Flow

Read:

- `workspace/docs/EXAMPLE-WORKFLOW.md`
- `workspace/agents/queue/EXAMPLE-TASK-local-memory-bootstrap.md`
- `workspace/work/examples/EXAMPLE-PROJECT-OVERVIEW.md`

## 5. Launch AthenaWork

```bash
./products/athena-work/tools/launch_stage.sh engineering
./products/athena-work/tools/launch_stage.sh qa
./products/athena-work/tools/run_observer_cycle.sh --cycle-id quickstart-example
```

The launcher emits stage instructions plus scoped approved tool context from AthenaUse.

## 6. Navigate The Platform

- platform overview: `README.md`
- platform layout: `docs/platform-layout.md`
- AthenaMind product: `products/athena-mind/README.md`
- AthenaWork product: `products/athena-work/README.md`
- workspace operator surface: `workspace/docs/HUMANS.md`

## Notes

- The default AthenaMind path is local and sqlite-first.
- Mongo is available as an optional stronger local backend for AthenaMind index and embedding persistence.
- Example workspace content is intentionally minimal and safe.

Optional Mongo-backed smoke:

```bash
ATHENA_MONGODB_URI='mongodb://admin:changeme@127.0.0.1:27017/?authSource=admin' \
ATHENA_MONGODB_DATABASE=athenamind \
./tools/platform/smoke_mongodb.sh
```
