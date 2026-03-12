# AthenaWork Skill

Use this skill when the task is about stage workflow, backlog progression, QA/PM handoff, observer cycles, release operations, or workspace operating contracts.

## Workspace Focus

- `products/athena-work`
- `workspace`

## Product Intent

AthenaWork in AthenaPlatform combines:

- executable staged workflow machinery
- markdown-native workspace structure
- reusable prompts, templates, and contracts

## Typical Commands

```bash
./products/athena-work/tools/launch_stage.sh engineering
./products/athena-work/tools/launch_stage.sh qa
./products/athena-work/tools/run_observer_cycle.sh --cycle-id <cycle-id>
```

## Launch Behavior

`launch_stage.sh` now composes three context sources for stage work:

- stage seed prompt
- AthenaMind memory context when available through the CLI contract
- AthenaUse tool context from approved repo-backed tools
