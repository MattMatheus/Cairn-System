# Example Workflow

This is a minimal example showing how the unified platform is intended to work.

Recommended local setup:

```bash
export CAIRN_HOME="$PWD/.cairn"
export CAIRN_MEMORY_ROOT="$CAIRN_HOME/memory/default"
mkdir -p "$CAIRN_MEMORY_ROOT"
```

## Step 1

Use work harness to define or pick a queue item:

- `workspace/agents/queue/EXAMPLE-TASK-local-memory-bootstrap.md`

## Step 2

Use memory-cli to create and retrieve local context:

```bash
(cd products/memory-cli && go run ./cmd/memory-cli write \
  --root "$CAIRN_MEMORY_ROOT" \
  --id example-bootstrap \
  --title "Example Bootstrap" \
  --type prompt \
  --domain platform \
  --body "Use workspace queue metadata and record observer output after cycle closure." \
  --stage planning \
  --reviewer matt \
  --decision approved \
  --reason "example bootstrap" \
  --risk "low" \
  --notes "example")

(cd products/memory-cli && go run ./cmd/memory-cli retrieve \
  --root "$CAIRN_MEMORY_ROOT" \
  --query "observer output after cycle closure" \
  --domain platform)
```

## Step 3

Advance work through work harness:

```bash
./products/work-harness/tools/launch_stage.sh engineering
./products/work-harness/tools/launch_stage.sh qa
./products/work-harness/tools/run_observer_cycle.sh --cycle-id example-local-memory-bootstrap
```

The stage launcher now emits:

- launch checklist
- approved `tool_context` from tool-cli
- workspace adapter status and other control-plane output

## Step 4

Record evidence in research/example notes if the pattern proves useful.
