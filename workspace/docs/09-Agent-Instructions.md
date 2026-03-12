# Agent Instructions

## Primary Directives
1. Follow `AGENTS.md` first.
2. Apply operating contracts from `products/athena-work/operating-system-vault/*`.
3. Use queue lanes for all delivery state movement.
4. Keep knowledge promotions evidence-backed.
5. Do not bypass human approval boundaries.

## Write Boundaries
- Allowed writes: queue, runs, memory candidates, staging.
- Suggested updates only for canonical work and research zones.
- Do not write private-personal zones unless explicitly requested.

## Minimum Agent Output Standard
- Correct frontmatter and status.
- Clear transition rationale.
- Linked evidence for claim promotions.
- Safety decision metadata when required.

## Template-First Rule
- For new work items, start from `workspace/templates/` instead of freeform files.
- Default mappings:
  - `task` -> `TemplateTask.md`
  - `claim` -> `TemplateClaim.md`
  - `artifact` -> `TemplateArtifact.md`

## Escalation Conditions
- Conflicting policy references.
- Ambiguous state transitions.
- Missing required metadata fields.
- Any potential sensitive-content violations.

## User Content Prune Workflow
- Use `tools/migration/prune_user_content.sh` for user-requested removal of personal work/research content.
- Use `workspace/agents/system.json` as the layout and prune-target contract.
- Offer zip export first with `--export-zip`.
- If export path permissions fail, rerun with escalated permissions rather than abandoning export.
