# Example Workflow

This is a minimal example showing how the unified platform is intended to work.

Recommended local setup:

```bash
export ATHENA_HOME="$PWD/.athena"
export ATHENA_MEMORY_ROOT="$ATHENA_HOME/memory/default"
mkdir -p "$ATHENA_MEMORY_ROOT"
```

## Step 1

Use AthenaWork to define or pick a queue item:

- `workspace/agents/queue/EXAMPLE-TASK-local-memory-bootstrap.md`

## Step 2

Use AthenaMind to create and retrieve local context:

```bash
(cd products/athena-mind && go run ./cmd/memory-cli write \
  --root "$ATHENA_MEMORY_ROOT" \
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

(cd products/athena-mind && go run ./cmd/memory-cli retrieve \
  --root "$ATHENA_MEMORY_ROOT" \
  --query "observer output after cycle closure" \
  --domain platform)
```

## Step 3

Advance work through AthenaWork:

```bash
./products/athena-work/tools/launch_stage.sh engineering
./products/athena-work/tools/launch_stage.sh qa
./products/athena-work/tools/run_observer_cycle.sh --cycle-id example-local-memory-bootstrap
```

The stage launcher now emits:

- launch checklist
- approved `tool_context` from AthenaUse
- workspace adapter status and other control-plane output

## Step 4

Record evidence in research/example notes if the pattern proves useful.
