---
tags:
  - Index
Publish: "False"
---

# HUMANS

Primary human entrypoint for the markdown workspace surface.

## System Concept

- The notebook or vault is the human review and knowledge surface.
- `memory-cli` stores durable promoted memory.
- `tool-cli` exposes governed tool context.
- `work-harness` runs the delivery and execution loop.

Use the vault for review and meaning.

Use the harness for execution.

For the shortest supported internal-beta path, start with:

1. `INTERNAL_BETA.md`
2. `PLATFORM_QUICKSTART.md`
3. `workspace/docs/HUMANS.md`

### Docs
- `workspace/docs/README.md`
- `workspace/docs/08-Human-Runbook.md`

### Templates
- `workspace/templates/TemplateTask.md`
- `workspace/templates/TemplateClaim.md`
- `workspace/templates/TemplateArtifact.md`
- `workspace/templates/TemplateResearch.md`

### Operating System
- `products/work-harness/HUMANS.md`
- `products/work-harness/DEVELOPMENT_CYCLE.md`
- `products/work-harness/operating-system-vault/README.md`
- `products/work-harness/operating-system-vault/PROGRAM_CONTROL_PLANE.md`
- `products/work-harness/operating-system-vault/STAGE_EXIT_GATES.md`
- When codegraph work is likely, run `products/work-harness/tools/check_gitnexus_readiness.sh`

### Work
- `workspace/work/`
- `workspace/agents/delivery/Engineering Active/README.md`
- `workspace/agents/delivery/Engineering QA/README.md`

### Research
- `workspace/research/`

### Agents
- `workspace/agents/maps/Agent Control Center.md`

### Follow Ups
```dataview
LIST
FROM #FollowUp
```
