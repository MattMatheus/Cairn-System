# memory-cli Quickstart

## Goal

Run a practical local write -> retrieve -> verify -> snapshot cycle.

If you installed a precompiled binary, replace `go run ./cmd/memory-cli` with `memory-cli` in each command.

Recommended local setup:

```bash
export CAIRN_HOME="$PWD/.cairn"
export CAIRN_MEMORY_ROOT="$CAIRN_HOME/memory/default"
mkdir -p "$CAIRN_MEMORY_ROOT"
```

## 1. Write One Entry

```bash
(cd products/memory-cli && go run ./cmd/memory-cli write \
  --root "$CAIRN_MEMORY_ROOT" \
  --id onboarding-guide \
  --title "Onboarding Guide" \
  --type prompt \
  --domain engineering \
  --body "Use deterministic fallbacks and always include source_path in retrieval outputs." \
  --stage planning \
  --reviewer maya \
  --decision approved \
  --reason "bootstrap baseline" \
  --risk "low" \
  --notes "approved")
```

## 2. Retrieve It

```bash
(cd products/memory-cli && go run ./cmd/memory-cli retrieve \
  --root "$CAIRN_MEMORY_ROOT" \
  --query "fallback and source path policy" \
  --domain engineering)
```

Confirm response includes `selected_id`, `selection_mode`, and `source_path`.

## 3. Verify Embeddings And Retrieval Health

```bash
(cd products/memory-cli && go run ./cmd/memory-cli verify embeddings --root "$CAIRN_MEMORY_ROOT")
(cd products/memory-cli && go run ./cmd/memory-cli verify health --root "$CAIRN_MEMORY_ROOT" --query "fallback and source path policy")
```

## 4. Create Snapshot

```bash
(cd products/memory-cli && go run ./cmd/memory-cli snapshot create \
  --root "$CAIRN_MEMORY_ROOT" \
  --created-by operator \
  --reason "post-quickstart checkpoint")
```

## 5. (Optional) Export OTel Traces

Telemetry and OTLP setup docs still need a dedicated Cairn import pass.
