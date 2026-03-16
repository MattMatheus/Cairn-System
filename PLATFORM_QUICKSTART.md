# Cairn Quickstart

This quickstart shows the intended developer path through the unified platform.

## What You Need

- Go installed
- a local checkout of this repository
- optional: Podman and Ollama if you want local embedding service support

## 0. Bootstrap The Repo

From repo root:

```bash
./tools/dev/bootstrap_platform.sh
```

Optional module download during bootstrap:

```bash
./tools/dev/bootstrap_platform.sh --download-modules
```

## Local Runtime Boundary

Use a repo-local `.cairn/` directory for uncommitted Cairn runtime state.

Suggested local folders:

```text
.cairn/
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
export CAIRN_HOME="$PWD/.cairn"
export CAIRN_MEMORY_ROOT="$CAIRN_HOME/memory/default"
mkdir -p "$CAIRN_MEMORY_ROOT"
```

## 1. Verify memory-cli

From repo root:

```bash
cd products/memory-cli && go test ./...
cd products/tool-cli && go test ./...
```

## 2. Create One Local Memory Entry

```bash
(cd products/memory-cli && go run ./cmd/memory-cli write \
  --root "$CAIRN_MEMORY_ROOT" \
  --id quickstart-bootstrap \
  --title "Quickstart Bootstrap" \
  --type prompt \
  --domain platform \
  --body "Use work harness queue metadata and capture observer output before closing a cycle." \
  --stage planning \
  --reviewer matt \
  --decision approved \
  --reason "platform quickstart" \
  --risk "low" \
  --notes "example bootstrap")
```

## 3. Retrieve It

```bash
(cd products/memory-cli && go run ./cmd/memory-cli retrieve \
  --root "$CAIRN_MEMORY_ROOT" \
  --query "observer output before closing a cycle" \
  --domain platform)
```

Optional embedding-backed validation:

```bash
(cd products/memory-cli && \
  MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS=0 \
  CAIRN_OLLAMA_EMBED_MODEL="mxbai-embed-large:latest" \
  go run ./cmd/memory-cli write \
    --root "$CAIRN_MEMORY_ROOT" \
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

(cd products/memory-cli && \
  MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS=0 \
  CAIRN_OLLAMA_EMBED_MODEL="mxbai-embed-large:latest" \
  go run ./cmd/memory-cli verify health \
    --root "$CAIRN_MEMORY_ROOT" \
    --query "improve local memory selection with embeddings" \
    --domain platform \
    --embedding-endpoint http://192.168.1.35:11434)
```

## 4. Inspect The Example Workspace Flow

Read:

- `workspace/docs/EXAMPLE-WORKFLOW.md`
- `workspace/agents/queue/EXAMPLE-TASK-local-memory-bootstrap.md`
- `workspace/work/examples/EXAMPLE-PROJECT-OVERVIEW.md`

## 5. Launch work harness

```bash
./products/work-harness/tools/launch_stage.sh engineering
./products/work-harness/tools/launch_stage.sh qa
./products/work-harness/tools/run_observer_cycle.sh --cycle-id quickstart-example
```

The launcher emits stage instructions plus scoped approved tool context from tool-cli.

## 6. Navigate The Platform

- platform overview: `README.md`
- platform layout: `docs/platform-layout.md`
- memory-cli product: `products/memory-cli/README.md`
- work harness product: `products/work-harness/README.md`
- workspace operator surface: `workspace/docs/HUMANS.md`
- release artifacts: `./tools/dev/build_release_artifacts.sh --version <label>`

## Notes

- The default memory-cli path is local and sqlite-first.
- Example workspace content is intentionally minimal and safe.
