# AthenaWork Quickstart

## Summary

AthenaWork is available directly in this repository for operator-driven workflows.

## Start Here

1. Read `workspace/docs/HUMANS.md` for workspace operator expectations.
2. Set repo-local runtime defaults for AthenaMind:

```bash
export ATHENA_HOME="$PWD/.athena"
export ATHENA_MEMORY_ROOT="$ATHENA_HOME/memory/default"
mkdir -p "$ATHENA_MEMORY_ROOT"
```

2. Launch planning or engineering stage.
3. Follow launcher checklist output.

## Minimal Cycle

From repo root:

```bash
./products/athena-work/tools/launch_stage.sh engineering
./products/athena-work/tools/launch_stage.sh qa
./products/athena-work/tools/run_observer_cycle.sh --cycle-id <story-id>
```

Then commit once for the cycle:

```bash
git commit -m "cycle-<cycle-id>"
```

## Non-Technical Operator Pattern

1. Pick stage goal in plain language.
2. Select or confirm specialist role from `products/athena-work/staff-personas/STAFF_DIRECTORY.md`.
3. Ask agent to execute only the stage checklist.
4. Require observer report before accepting completion.

## Related Docs

- `products/athena-work/README.md`
- `products/athena-work/HUMANS.md`
