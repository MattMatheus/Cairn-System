# memory-cli Skill

Use this skill when the task is about memory indexing, retrieval quality, embeddings, snapshots, or memory CLI behavior.

## Workspace Focus

- `products/memory-cli`

## V1 Product Intent

- markdown-first ingestion from work harness content
- `sqlite` as the default local backend
- Go CLI as the main integration surface
- OpenTelemetry is required for all memory-cli runtime paths

## Preferred V1 Commands

- `write`
- `retrieve`
- `bootstrap`
- `verify`
- `snapshot`

If a task starts drifting toward lightweight wrapper/proxy behavior, prefer `tool-cli` as the boundary for that work instead of expanding memory-cli's role.

## Typical Commands

```bash
CAIRN_HOME="${CAIRN_HOME:-$PWD/.cairn}"
CAIRN_MEMORY_ROOT="${CAIRN_MEMORY_ROOT:-$CAIRN_HOME/memory/core}"
(cd products/memory-cli && go run ./cmd/memory-cli write --root "$CAIRN_MEMORY_ROOT" --id example --title "Example" --type prompt --domain platform --body "example body" --stage planning --reviewer matt --decision approved --reason "seed" --risk "low" --notes "bootstrap")
(cd products/memory-cli && go run ./cmd/memory-cli retrieve --root "$CAIRN_MEMORY_ROOT" --query "memory lifecycle")
(cd products/memory-cli && go run ./cmd/memory-cli bootstrap --root "$CAIRN_MEMORY_ROOT" --repo Cairn --session-id local-bootstrap --scenario engineering)
(cd products/memory-cli && go run ./cmd/memory-cli verify embeddings --root "$CAIRN_MEMORY_ROOT")
(cd products/memory-cli && go run ./cmd/memory-cli verify health --root "$CAIRN_MEMORY_ROOT" --query "memory lifecycle")
```
