# AGENTS.md - work harness Agent Operating Instructions

## Mission
Run work harness as a supervised, claim-first system with two connected planes:
- Delivery Plane: ship work through explicit queue states.
- Knowledge Plane: promote claims only with evidence.

## First 10 Minutes (New Agent)
1. Read `workspace/docs/HUMANS.md`.
2. Read `workspace/docs/09-Agent-Instructions.md`.
3. Read `products/work-harness/operating-system-vault/README.md`.
4. Read `products/work-harness/operating-system-vault/STAGE_EXIT_GATES.md`.
5. Read `products/work-harness/operating-system-vault/UNIFIED_METADATA_CONTRACT.md`.
6. Read `workspace/docs/policies/Content Rejection Policy.md`.
7. Open `workspace/agents/maps/Agent Control Center.md`.
8. Use templates from `workspace/templates/` for new notes/items.

## Non-Negotiable Rules
- Do not create new unnumbered top-level folders.
- Respect write boundaries.
- Every state transition must have explicit rationale.
- Claim promotions to `supported` or `falsified` require linked artifacts.
- If safety is uncertain, route to quarantine and escalate.

## Canonical Structure
- `workspace/work`
- `workspace/docs`
- `workspace/templates`
- `workspace/agents`
- `workspace/research`
- `products/work-harness`

## Write Boundaries
- `Zone A (agent write)`:
  - `workspace/agents/queue`
  - `workspace/agents/delivery/*`
  - `workspace/agents/memory/*`
  - `workspace/agents/staging/*`
  - `workspace/agents/escape/*`
- `Zone B (agent suggest, human approves)`:
  - `workspace/work/*`
  - `workspace/research/*`

## State Models
- Knowledge lifecycle: `idea -> candidate -> active-test -> supported|falsified|parked`
- Delivery lifecycle: `intake -> active -> qa -> done` (or back to `active` with defects)

## Required Metadata
All operational notes must include:
```yaml
---
id: ath-<id>
type: concept|claim|artifact|decision|run|memory|report|map|task|story|bug
status: idea|candidate|active-test|supported|falsified|parked|intake|active|qa|done|blocked|archived
domain: work|research
updated: YYYY-MM-DD
owner: matt
source_of_truth: human|agent
sensitivity: public|internal|confidential|private_personal
---
```

Agent writes must also include:
```yaml
agent_last_touch: YYYY-MM-DDTHH:MM:SSZ
agent_intent: ingest|summarize|extract|classify|redact|draft|implement|verify
review_state: pending|approved|rejected
```

## Execution Pattern (Default)
1. Intake: create or refine queue item in `workspace/agents/queue` or a delivery lane.
2. Plan: map task to explicit next transition and gate criteria.
3. Execute: do work in the active lane and track artifacts.
4. Validate: check stage exit gate before transition.
5. Record: log run notes in an approved run-log location once the canonical run path is finalized.
6. Promote: move to staging or approved destination with links.

## Stage Exit Checks
- Planning: intake complete, next stage explicit.
- Architect: scope and downstream mapping explicit.
- Engineering: implementation and tests complete, handoff package present.
- QA: acceptance evidence and regression statement explicit.
- Knowledge promotion: state change explicit with artifact links.
- Cycle closure: observer/report artifact exists and cycle sync complete.

## Escape Path

When you cannot safely proceed on a work item, you MUST escape rather than stall, guess, or loop.

**Trigger escape when:**
- A hard dependency is missing (`blocked`)
- Acceptance criteria are too ambiguous to determine done (`ambiguous_requirements`)
- Two policies conflict and you cannot self-resolve (`policy_conflict`)
- Insufficient context to proceed safely (`missing_context`)
- Confidence is below threshold for the current transition (`low_confidence`)
- The same transition has been attempted 2+ times without progress (`loop_detected`)
- The task requires authorization beyond your write boundaries (`scope_exceeded`)

**Escape procedure:**
1. Set source item `status: blocked` in place (do not move it out of its lane folder).
2. Create an escape record in `workspace/agents/escape/` using `workspace/templates/TemplateEscape.md`.
   - Name: `ESC-<source-id>-<YYYY-MM-DD>.md`
3. Fill in: `escape_class`, `escape_summary`, `attempt_count`, what was attempted, what is needed.
4. Set `review_state: pending` and `agent_intent: escalate`.
5. Stop. Do not attempt further transitions on the blocked item.

Human resolves by setting `resolution_action` and returning the item to queue.

## Safety Pattern
- Apply `workspace/docs/policies/Content Rejection Policy.md`.
- For questionable or sensitive ingest, attach safety decision metadata.
- Escalate on policy conflicts, missing metadata, or ambiguous transitions.

## Core Techniques
- Small safe cycles: prefer incremental transitions over large jumps.
- Evidence-first notes: link artifacts before claim promotion.
- Queue discipline: work only from explicit queue states.
- Contract-first writing: frontmatter correctness before content polish.
- Human gatekeeping: treat `review_state` as mandatory control.

## Template Routing
- `type: task` -> start from `workspace/templates/TemplateTask.md`
- `type: claim` -> start from `workspace/templates/TemplateClaim.md`
- `type: artifact` -> start from `workspace/templates/TemplateArtifact.md`
- `type: concept` -> start from `workspace/templates/TemplateConcept.md`
- For daily/journal capture -> `workspace/templates/TemplateDaily.md`
- Prefer template-first creation over freeform note starts.

## User Content Removal Protocol
- Layout/source contract: `workspace/agents/system.json`.
- Use `tools/migration/prune_user_content.sh` when a user asks to remove personal work/research content.
- Always offer backup export before deletion:
  - `tools/migration/prune_user_content.sh --export-zip /writable/path/work-harness-user-content-backup.zip`
- Destructive execution requires explicit confirmation token unless `--yes` is provided.
- If export fails due sandbox/permissions, rerun with escalated permissions instead of skipping export.
- After prune: review diff, commit only prune-related changes, and push a review branch.

## Vendor & API Docs

Before implementing or researching any vendor integration, check the registry for the pinned version and canonical docs URL:

- `products/work-harness/operating-system-vault/README.md`

Do not assume `latest`. If a vendor or version is missing from the registry, escape with `escape_class: missing_context`.

## Primary References
- `workspace/docs/09-Agent-Instructions.md`
- `products/work-harness/operating-system-vault/PROGRAM_CONTROL_PLANE.md`
- `products/work-harness/operating-system-vault/STAGE_EXIT_GATES.md`
- `products/work-harness/operating-system-vault/DELIVERY_STATE_MODEL.md`
- `products/work-harness/operating-system-vault/KNOWLEDGE_STATE_MODEL.md`
- `products/work-harness/operating-system-vault/UNIFIED_METADATA_CONTRACT.md`
- `workspace/agents/policies/Plugin Baseline.md`
