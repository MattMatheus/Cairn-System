# CLI Examples

## Summary
Copy/paste examples for common CLI tasks.

These examples cover the AthenaMind commands intentionally exposed in AthenaPlatform.

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

## Promote Curated Note
```bash
cd ../athena-use && go run ./cmd/promote-cli note \
  --root "$ATHENA_MEMORY_ROOT" \
  --reviewer matt \
  --reason "durable operating guidance worth retrieving later" \
  --risk "low; curated manually from vault note" \
  --notes "promoted after human review" \
  "/home/matt/Workspace/Athena/80 Research/40 Decisions/Cairn North Star.md"
```

## Retrieve Entry
```bash
go run ./cmd/memory-cli retrieve \
  --root "$ATHENA_MEMORY_ROOT" \
  --query "handoff instruction template" \
  --mode classic \
  --top-k 5 \
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
  --domain docs \
  --session-id docs-session-1
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
