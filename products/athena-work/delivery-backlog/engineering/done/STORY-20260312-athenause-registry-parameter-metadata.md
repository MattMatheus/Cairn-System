# Story: Enrich AthenaUse approved registry parameter metadata

## Metadata
- `id`: STORY-20260312-athenause-registry-parameter-metadata
- `owner_persona`: Product Manager - AthenaWork
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0005, ADR-0006, ADR-0007]
- `success_metric`: approved AthenaUse tools with input parameters expose useful descriptions and enums in context output for operator-facing consumption
- `release_checkpoint`: deferred

## Problem Statement
- AthenaUse context output now supports structured parameter metadata, but the approved registry still leaves most descriptions and enums empty, so the richer context is only partially useful.

## Scope
- In:
  - add parameter descriptions and enum metadata where appropriate in `products/athena-use/registry/approved-tools.yaml`
  - update tests to verify the approved registry exposes richer parameter metadata
  - confirm `use-cli context` surfaces the enriched metadata cleanly
- Out:
  - new tool entries
  - tool execution support
  - full schema-system redesign

## Acceptance Criteria
1. Approved tools with parameters expose meaningful descriptions in the registry and context output.
2. Enum metadata is added where the allowed values are known and stable.
3. Tests verify the shipped approved registry still validates and the richer metadata appears in context output.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `products/athena-use/registry/approved-tools.yaml`
- `use-cli context` structured parameter output

## Notes
- Product-first follow-on: this makes the bounded context/schema slice materially more useful without reopening interface-scope questions.
