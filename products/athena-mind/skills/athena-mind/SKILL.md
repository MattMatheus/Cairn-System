# AthenaMind Skill

Use this skill when the task is about memory indexing, retrieval quality, embeddings, snapshots, or memory CLI behavior.

## Workspace Focus

- `products/athena-mind`

## V1 Product Intent

- markdown-first ingestion from AthenaWork content
- `sqlite` as the default local backend
- optional Mongo support is available for developers who want stronger local persistence, but it is not the default operating path
- Go CLI as the main integration surface

## Preferred V1 Commands

- `write`
- `retrieve`
- `bootstrap`
- `verify`
- `snapshot`

## Typical Commands

```bash
ATHENA_HOME="${ATHENA_HOME:-$PWD/.athena}"
ATHENA_MEMORY_ROOT="${ATHENA_MEMORY_ROOT:-$ATHENA_HOME/memory/core}"
(cd products/athena-mind && go run ./cmd/memory-cli write --root "$ATHENA_MEMORY_ROOT" --id example --title "Example" --type prompt --domain platform --body "example body" --stage planning --reviewer matt --decision approved --reason "seed" --risk "low" --notes "bootstrap")
(cd products/athena-mind && go run ./cmd/memory-cli retrieve --root "$ATHENA_MEMORY_ROOT" --query "memory lifecycle")
(cd products/athena-mind && go run ./cmd/memory-cli bootstrap --root "$ATHENA_MEMORY_ROOT" --repo AthenaPlatform --scenario engineering)
(cd products/athena-mind && go run ./cmd/memory-cli verify embeddings --root "$ATHENA_MEMORY_ROOT")
(cd products/athena-mind && go run ./cmd/memory-cli verify health --root "$ATHENA_MEMORY_ROOT" --query "memory lifecycle")
```
