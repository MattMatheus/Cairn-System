# Story: Strengthen tool-cli context shaping and schema output

## Metadata
- `id`: STORY-20260312-toolcli-context-schema-output
- `owner_persona`: Product Manager - work harness
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0005, ADR-0006, ADR-0007]
- `success_metric`: tool-cli context output carries enough schema detail for downstream stage launch consumers without ad hoc parsing
- `release_checkpoint`: deferred

## Problem Statement
- tool-cli context output currently summarizes approved tools, but schema detail is shallow and may not be sufficient once more tools are registered and downstream consumers depend on structured context.

## Scope
- In:
  - Improve `tool-cli context` output shape for richer schema/context emission
  - Keep output contract small and stage-launch friendly
  - Add tests for the emitted shape
- Out:
  - Tool execution support
  - New trust tiers
  - Full model/tool interface specification work

## Acceptance Criteria
1. `tool-cli context` emits stable schema-oriented fields that cover required parameters for registered tools.
2. Output format changes are covered by tests.
3. Stage launch consumption remains compatible or is updated in the same cycle.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- Current `tool-cli context` output contract
- Outcome of the tool interface specification decision

## Notes
- Promoted after `ADR-0007` selected bounded context/schema shaping as the next approved tool-cli slice.
