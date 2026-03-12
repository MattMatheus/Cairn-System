# Story: Strengthen AthenaUse context shaping and schema output

## Metadata
- `id`: STORY-20260312-athenause-context-schema-output
- `owner_persona`: Product Manager - AthenaWork
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0005, ADR-0006, ADR-0007]
- `success_metric`: AthenaUse context output carries enough schema detail for downstream stage launch consumers without ad hoc parsing
- `release_checkpoint`: deferred

## Problem Statement
- AthenaUse context output currently summarizes approved tools, but schema detail is shallow and may not be sufficient once more tools are registered and downstream consumers depend on structured context.

## Scope
- In:
  - Improve `use-cli context` output shape for richer schema/context emission
  - Keep output contract small and stage-launch friendly
  - Add tests for the emitted shape
- Out:
  - Tool execution support
  - New trust tiers
  - Full model/tool interface specification work

## Acceptance Criteria
1. `use-cli context` emits stable schema-oriented fields that cover required parameters for registered tools.
2. Output format changes are covered by tests.
3. Stage launch consumption remains compatible or is updated in the same cycle.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- Current `use-cli context` output contract
- Outcome of the tool interface specification decision

## Notes
- Promoted after `ADR-0007` selected bounded context/schema shaping as the next approved AthenaUse slice.
