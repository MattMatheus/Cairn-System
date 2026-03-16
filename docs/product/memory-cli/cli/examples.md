# CLI Examples

## Summary
Copy/paste examples for common CLI tasks.

These examples cover the memory-cli commands intentionally exposed in Cairn.

Recommended local setup:

```bash
export CAIRN_HOME="$PWD/.cairn"
export CAIRN_MEMORY_ROOT="$CAIRN_HOME/memory/default"
mkdir -p "$CAIRN_MEMORY_ROOT"
```

## Write Entry
```bash
go run ./cmd/memory-cli write \
  --root "$CAIRN_MEMORY_ROOT" \
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
cd ../tool-cli && go run ./cmd/promote-cli note \
  --root "$CAIRN_MEMORY_ROOT" \
  --reviewer matt \
  --reason "durable operating guidance worth retrieving later" \
  --risk "low; curated manually from vault note" \
  --notes "promoted after human review" \
  "/home/matt/Workspace/Cairn/80 Research/40 Decisions/Cairn North Star.md"
```

## Retrieve Entry
```bash
go run ./cmd/memory-cli retrieve \
  --root "$CAIRN_MEMORY_ROOT" \
  --query "handoff instruction template" \
  --mode classic \
  --top-k 5 \
  --embedding-endpoint http://localhost:11434
```

## Verify Embedding Coverage And Semantic Health
```bash
go run ./cmd/memory-cli verify embeddings \
  --root "$CAIRN_MEMORY_ROOT"
```

```bash
go run ./cmd/memory-cli verify health \
  --root "$CAIRN_MEMORY_ROOT" \
  --query "memory lifecycle" \
  --domain docs \
  --session-id docs-session-1
```

## Snapshot Create/List/Restore
```bash
go run ./cmd/memory-cli snapshot create \
  --root "$CAIRN_MEMORY_ROOT" \
  --created-by clara \
  --reason "pre-release docs freeze"
```

```bash
go run ./cmd/memory-cli snapshot list --root "$CAIRN_MEMORY_ROOT"
```

```bash
go run ./cmd/memory-cli snapshot restore \
  --root "$CAIRN_MEMORY_ROOT" \
  --snapshot-id <snapshot-id> \
  --stage pm \
  --reviewer maya \
  --decision approved \
  --reason "rollback to known good set" \
  --risk "low; manifest/checksum verified" \
  --notes "approved restore"
```
