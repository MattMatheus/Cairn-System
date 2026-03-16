# work harness Skill

Use this skill when the task is about stage workflow, backlog progression, QA/PM handoff, observer cycles, release operations, or workspace operating contracts.

## Workspace Focus

- `products/work-harness`
- `workspace`

## Product Intent

work harness in Cairn combines:

- executable staged workflow machinery
- markdown-native workspace structure
- reusable prompts, templates, and contracts

## Typical Commands

```bash
./products/work-harness/tools/launch_stage.sh engineering
./products/work-harness/tools/launch_stage.sh qa
./products/work-harness/tools/run_observer_cycle.sh --cycle-id <cycle-id>
```

## Launch Behavior

`launch_stage.sh` now composes three context sources for stage work:

- stage seed prompt
- memory-cli memory context when available through the CLI contract
- tool-cli tool context from approved repo-backed tools
