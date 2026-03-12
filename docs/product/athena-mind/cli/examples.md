# CLI Examples

## Summary
Copy/paste examples for common CLI tasks.

Recommended local setup:

```bash
export ATHENA_HOME="$PWD/.athena"
export ATHENA_MEMORY_ROOT="$ATHENA_HOME/memory/default"
mkdir -p "$ATHENA_MEMORY_ROOT"
```

## Write Entry
```bash
go run ./cmd/memory-cli write \
  --root "$ATHENA_MEMORY_ROOT" \
  --id handoff-template \
  --title "Handoff Template" \
  --type instruction \
  --domain docs \
  --body "Always include risks, evidence, and next-state recommendation." \
  --stage pm \
  --reviewer maya \
  --decision approved \
  --reason "baseline docs quality" \
  --risk "low; reversible by git revert" \
  --notes "approved for docs baseline" \
  --embedding-endpoint http://localhost:11434
```

## Retrieve Entry
```bash
go run ./cmd/memory-cli retrieve \
  --root "$ATHENA_MEMORY_ROOT" \
  --query "handoff instruction template" \
  --mode hybrid \
  --top-k 10 \
  --retrieval-backend qdrant \
  --embedding-endpoint http://localhost:11434
```

## Run Evaluation
```bash
go run ./cmd/memory-cli evaluate \
  --root "$ATHENA_MEMORY_ROOT" \
  --query-file cmd/memory-cli/testdata/eval-query-set-v1.json \
  --corpus-id memory-corpus-v1 \
  --query-set-id query-set-v1 \
  --mode hybrid \
  --top-k 10 \
  --retrieval-backend sqlite \
  --embedding-endpoint http://localhost:11434
```

## Verify Embedding Coverage And Semantic Health
```bash
go run ./cmd/memory-cli verify embeddings \
  --root "$ATHENA_MEMORY_ROOT"
```

```bash
go run ./cmd/memory-cli verify health \
  --root "$ATHENA_MEMORY_ROOT" \
  --query "memory lifecycle" \
  --domain docs-crawl \
  --session-id docs-session-1
```

## Verify Optional MongoDB Contract
```bash
go run ./cmd/memory-cli verify mongodb \
  --mongodb-uri mongodb://127.0.0.1:27017 \
  --mongodb-database athenamind
```

## Use MongoDB-Backed Index Persistence
```bash
ATHENA_INDEX_BACKEND=mongodb \
ATHENA_MONGODB_URI='mongodb://admin:changeme@127.0.0.1:27017/?authSource=admin' \
ATHENA_MONGODB_DATABASE=athenamind \
go run ./cmd/memory-cli write \
  --root "$ATHENA_MEMORY_ROOT" \
  --id mongo-backed-entry \
  --title "Mongo Backed Entry" \
  --type prompt \
  --domain platform \
  --body "Keep markdown files local while persisting index state in MongoDB." \
  --stage planning \
  --reviewer maya \
  --decision approved \
  --reason "mongodb adapter example" \
  --risk "low" \
  --notes "approved"
```

## Sync Embeddings To Qdrant
```bash
go run ./cmd/memory-cli sync-qdrant \
  --root "$ATHENA_MEMORY_ROOT" \
  --qdrant-url http://localhost:6333 \
  --collection athena_memories \
  --batch-size 128
```

## Snapshot Create/List/Restore
```bash
go run ./cmd/memory-cli snapshot create \
  --root "$ATHENA_MEMORY_ROOT" \
  --created-by clara \
  --reason "pre-release docs freeze"
```

```bash
go run ./cmd/memory-cli snapshot list --root "$ATHENA_MEMORY_ROOT"
```

```bash
go run ./cmd/memory-cli snapshot restore \
  --root "$ATHENA_MEMORY_ROOT" \
  --snapshot-id <snapshot-id> \
  --stage pm \
  --reviewer maya \
  --decision approved \
  --reason "rollback to known good set" \
  --risk "low; manifest/checksum verified" \
  --notes "approved restore"
```

## API Retrieve With Fallback
```bash
go run ./cmd/memory-cli api-retrieve \
  --root "$ATHENA_MEMORY_ROOT" \
  --query "handoff instruction template" \
  --session-id docs-session-1 \
  --mode hybrid \
  --top-k 10 \
  --retrieval-backend qdrant \
  --gateway-url http://127.0.0.1:8788
```
